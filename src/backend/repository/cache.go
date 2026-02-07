package repository

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/pkg/errors"

	"r2manager/domain"
)

type CacheRepository struct {
	db           *sql.DB
	cacheDir     string
	ttl          time.Duration
	maxCacheSize int64
}

var (
	once sync.Once
	repo *CacheRepository
)

func NewCacheRepository(db *sql.DB, cacheDir string, ttl time.Duration, opts ...CacheOption) *CacheRepository {
	once.Do(func() {
		repo = &CacheRepository{
			db:       db,
			cacheDir: cacheDir,
			ttl:      ttl,
		}
		for _, opt := range opts {
			opt(repo)
		}
	})
	return repo
}

type CacheOption func(*CacheRepository)

func WithMaxCacheSize(size int64) CacheOption {
	return func(r *CacheRepository) {
		r.maxCacheSize = size
	}
}

func (r *CacheRepository) Lookup(ctx context.Context, bucketName, objectKey string) (*domain.CacheEntry, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT bucket_name, object_key, content_type, size, etag, cache_path, cached_at, expires_at
		 FROM cache_entries
		 WHERE bucket_name = ? AND object_key = ? AND expires_at > ?`,
		bucketName, objectKey, time.Now().UTC(),
	)

	var entry domain.CacheEntry
	err := row.Scan(
		&entry.BucketName,
		&entry.ObjectKey,
		&entry.ContentType,
		&entry.Size,
		&entry.ETag,
		&entry.CachePath,
		&entry.CachedAt,
		&entry.ExpiresAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup cache entry")
	}

	return &entry, nil
}

func (r *CacheRepository) Store(ctx context.Context, bucketName, objectKey string, body io.Reader, contentType string, size int64, etag string) (*domain.CacheEntry, error) {
	cachePath := r.cachePath(bucketName, objectKey)

	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		return nil, errors.Wrap(err, "failed to create cache directory")
	}

	tmpPath := cachePath + ".tmp"
	f, err := os.Create(tmpPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create temp cache file")
	}

	if _, err := io.Copy(f, body); err != nil {
		f.Close()
		os.Remove(tmpPath)
		return nil, errors.Wrap(err, "failed to write cache file")
	}
	if err := f.Close(); err != nil {
		os.Remove(tmpPath)
		return nil, errors.Wrap(err, "failed to close cache file")
	}

	if err := os.Rename(tmpPath, cachePath); err != nil {
		os.Remove(tmpPath)
		return nil, errors.Wrap(err, "failed to rename cache file")
	}

	now := time.Now().UTC()
	expiresAt := now.Add(r.ttl)

	_, err = r.db.ExecContext(ctx,
		`INSERT INTO cache_entries (bucket_name, object_key, content_type, size, etag, cache_path, cached_at, expires_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT(bucket_name, object_key) DO UPDATE SET
		   content_type = excluded.content_type,
		   size = excluded.size,
		   etag = excluded.etag,
		   cache_path = excluded.cache_path,
		   cached_at = excluded.cached_at,
		   expires_at = excluded.expires_at`,
		bucketName, objectKey, contentType, size, etag, cachePath, now, expiresAt,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to upsert cache entry")
	}

	if _, err := r.Evict(ctx); err != nil {
		log.Printf("cache eviction error: %v", err)
	}

	return &domain.CacheEntry{
		BucketName:  bucketName,
		ObjectKey:   objectKey,
		ContentType: contentType,
		Size:        size,
		ETag:        etag,
		CachePath:   cachePath,
		CachedAt:    now,
		ExpiresAt:   expiresAt,
	}, nil
}

func (r *CacheRepository) OpenCacheFile(cachePath string) (io.ReadCloser, error) {
	f, err := os.Open(cachePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open cache file")
	}
	return f, nil
}

func (r *CacheRepository) InvalidateByETags(ctx context.Context, bucketName string, currentETags map[string]string) (int, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT object_key, etag, cache_path FROM cache_entries WHERE bucket_name = ?`,
		bucketName,
	)
	if err != nil {
		return 0, errors.Wrap(err, "failed to query cache entries for ETag check")
	}
	defer rows.Close()

	type staleEntry struct {
		objectKey string
		cachePath string
	}
	var stale []staleEntry

	for rows.Next() {
		var key, cachedETag, path string
		if err := rows.Scan(&key, &cachedETag, &path); err != nil {
			return 0, errors.Wrap(err, "failed to scan cache entry")
		}
		if s3ETag, exists := currentETags[key]; exists && s3ETag != cachedETag {
			stale = append(stale, staleEntry{objectKey: key, cachePath: path})
		}
	}
	if err := rows.Err(); err != nil {
		return 0, errors.Wrap(err, "failed to iterate cache entries")
	}

	for _, entry := range stale {
		_, err := r.db.ExecContext(ctx,
			`DELETE FROM cache_entries WHERE bucket_name = ? AND object_key = ?`,
			bucketName, entry.objectKey,
		)
		if err != nil {
			return 0, errors.Wrap(err, "failed to delete stale cache entry")
		}
		os.Remove(entry.cachePath)
	}

	return len(stale), nil
}

func (r *CacheRepository) Evict(ctx context.Context) (int, error) {
	if r.maxCacheSize <= 0 {
		return 0, nil
	}

	var totalSize int64
	if err := r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(size), 0) FROM cache_entries`).Scan(&totalSize); err != nil {
		return 0, errors.Wrap(err, "failed to get total cache size")
	}

	if totalSize <= r.maxCacheSize {
		return 0, nil
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT bucket_name, object_key, size, cache_path FROM cache_entries ORDER BY cached_at ASC`,
	)
	if err != nil {
		return 0, errors.Wrap(err, "failed to query cache entries for eviction")
	}
	defer rows.Close()

	type evictEntry struct {
		bucketName string
		objectKey  string
		size       int64
		cachePath  string
	}
	var toEvict []evictEntry

	for rows.Next() && totalSize > r.maxCacheSize {
		var e evictEntry
		if err := rows.Scan(&e.bucketName, &e.objectKey, &e.size, &e.cachePath); err != nil {
			return 0, errors.Wrap(err, "failed to scan cache entry for eviction")
		}
		toEvict = append(toEvict, e)
		totalSize -= e.size
	}
	if err := rows.Err(); err != nil {
		return 0, errors.Wrap(err, "failed to iterate cache entries for eviction")
	}

	for _, e := range toEvict {
		_, err := r.db.ExecContext(ctx,
			`DELETE FROM cache_entries WHERE bucket_name = ? AND object_key = ?`,
			e.bucketName, e.objectKey,
		)
		if err != nil {
			return 0, errors.Wrap(err, "failed to delete evicted cache entry")
		}
		os.Remove(e.cachePath)
	}

	return len(toEvict), nil
}

func (r *CacheRepository) CleanupExpired(ctx context.Context) (int, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT bucket_name, object_key, cache_path FROM cache_entries WHERE expires_at <= ?`,
		time.Now().UTC(),
	)
	if err != nil {
		return 0, errors.Wrap(err, "failed to query expired cache entries")
	}
	defer rows.Close()

	type expiredEntry struct {
		bucketName string
		objectKey  string
		cachePath  string
	}
	var expired []expiredEntry

	for rows.Next() {
		var e expiredEntry
		if err := rows.Scan(&e.bucketName, &e.objectKey, &e.cachePath); err != nil {
			return 0, errors.Wrap(err, "failed to scan expired cache entry")
		}
		expired = append(expired, e)
	}
	if err := rows.Err(); err != nil {
		return 0, errors.Wrap(err, "failed to iterate expired cache entries")
	}

	for _, e := range expired {
		_, err := r.db.ExecContext(ctx,
			`DELETE FROM cache_entries WHERE bucket_name = ? AND object_key = ?`,
			e.bucketName, e.objectKey,
		)
		if err != nil {
			return 0, errors.Wrap(err, "failed to delete expired cache entry")
		}
		os.Remove(e.cachePath)
	}

	return len(expired), nil
}

func (r *CacheRepository) Vacuum() error {
	_, err := r.db.Exec("VACUUM")
	if err != nil {
		return errors.Wrap(err, "failed to vacuum database")
	}
	return nil
}

func (r *CacheRepository) StartCleanupLoop(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				totalDeleted := 0

				deleted, err := r.CleanupExpired(ctx)
				if err != nil {
					log.Printf("cache cleanup error: %v", err)
				} else {
					totalDeleted += deleted
				}

				evicted, err := r.Evict(ctx)
				if err != nil {
					log.Printf("cache eviction error: %v", err)
				} else {
					totalDeleted += evicted
				}

				if totalDeleted > 0 {
					log.Printf("cache cleanup: deleted %d expired, evicted %d over-limit", deleted, evicted)
					if err := r.Vacuum(); err != nil {
						log.Printf("cache vacuum error: %v", err)
					}
				}
			}
		}
	}()
}

func (r *CacheRepository) cachePath(bucketName, objectKey string) string {
	hash := sha256.Sum256([]byte(objectKey))
	filename := fmt.Sprintf("%x", hash)
	return filepath.Join(r.cacheDir, bucketName, filename)
}
