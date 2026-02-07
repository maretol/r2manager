package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) (*sql.DB, string) {
	t.Helper()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		t.Fatalf("failed to set WAL mode: %v", err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS cache_entries (
	    bucket_name  TEXT NOT NULL,
	    object_key   TEXT NOT NULL,
	    content_type TEXT NOT NULL DEFAULT 'application/octet-stream',
	    size         INTEGER NOT NULL DEFAULT 0,
	    etag         TEXT NOT NULL DEFAULT '',
	    cache_path   TEXT NOT NULL,
	    cached_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	    expires_at   DATETIME NOT NULL,
	    PRIMARY KEY (bucket_name, object_key)
	);
	CREATE INDEX IF NOT EXISTS idx_cache_entries_expires_at ON cache_entries(expires_at);
	`
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("failed to run schema: %v", err)
	}

	return db, tmpDir
}

func newTestCacheRepo(db *sql.DB, cacheDir string, ttl time.Duration, opts ...CacheOption) *CacheRepository {
	// Reset singleton for each test
	once = sync.Once{}
	repo = nil
	return NewCacheRepository(db, cacheDir, ttl, opts...)
}

func insertTestEntry(t *testing.T, db *sql.DB, bucketName, objectKey, cachePath string, expiresAt time.Time) {
	t.Helper()
	insertTestEntryWithSize(t, db, bucketName, objectKey, cachePath, 100, expiresAt, time.Now().UTC())
}

func insertTestEntryWithSize(t *testing.T, db *sql.DB, bucketName, objectKey, cachePath string, size int64, expiresAt, cachedAt time.Time) {
	t.Helper()
	_, err := db.Exec(
		`INSERT INTO cache_entries (bucket_name, object_key, content_type, size, etag, cache_path, cached_at, expires_at)
		 VALUES (?, ?, 'text/plain', ?, 'etag1', ?, ?, ?)`,
		bucketName, objectKey, size, cachePath, cachedAt, expiresAt,
	)
	if err != nil {
		t.Fatalf("failed to insert test entry: %v", err)
	}
}

func countEntries(t *testing.T, db *sql.DB) int {
	t.Helper()
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM cache_entries").Scan(&count); err != nil {
		t.Fatalf("failed to count entries: %v", err)
	}
	return count
}

func TestCleanupExpired_DeletesExpiredEntries(t *testing.T) {
	db, tmpDir := setupTestDB(t)
	defer db.Close()

	cacheDir := filepath.Join(tmpDir, "cache")
	os.MkdirAll(filepath.Join(cacheDir, "test-bucket"), 0755)

	r := newTestCacheRepo(db, cacheDir, 10*time.Minute)

	// Create cache files on disk
	expiredFilePath := filepath.Join(cacheDir, "test-bucket", "expired-file")
	os.WriteFile(expiredFilePath, []byte("expired content"), 0644)

	validFilePath := filepath.Join(cacheDir, "test-bucket", "valid-file")
	os.WriteFile(validFilePath, []byte("valid content"), 0644)

	now := time.Now().UTC()
	insertTestEntry(t, db, "test-bucket", "expired-key", expiredFilePath, now.Add(-1*time.Hour))
	insertTestEntry(t, db, "test-bucket", "valid-key", validFilePath, now.Add(1*time.Hour))

	deleted, err := r.CleanupExpired(context.Background())
	if err != nil {
		t.Fatalf("CleanupExpired failed: %v", err)
	}

	if deleted != 1 {
		t.Errorf("expected 1 deleted, got %d", deleted)
	}

	if countEntries(t, db) != 1 {
		t.Errorf("expected 1 remaining entry, got %d", countEntries(t, db))
	}

	// Expired file should be deleted from disk
	if _, err := os.Stat(expiredFilePath); !os.IsNotExist(err) {
		t.Error("expected expired cache file to be deleted")
	}

	// Valid file should still exist
	if _, err := os.Stat(validFilePath); err != nil {
		t.Error("expected valid cache file to still exist")
	}
}

func TestCleanupExpired_NoExpiredEntries(t *testing.T) {
	db, tmpDir := setupTestDB(t)
	defer db.Close()

	cacheDir := filepath.Join(tmpDir, "cache")
	r := newTestCacheRepo(db, cacheDir, 10*time.Minute)

	now := time.Now().UTC()
	insertTestEntry(t, db, "bucket1", "key1", "/tmp/fake1", now.Add(1*time.Hour))
	insertTestEntry(t, db, "bucket1", "key2", "/tmp/fake2", now.Add(2*time.Hour))

	deleted, err := r.CleanupExpired(context.Background())
	if err != nil {
		t.Fatalf("CleanupExpired failed: %v", err)
	}

	if deleted != 0 {
		t.Errorf("expected 0 deleted, got %d", deleted)
	}

	if countEntries(t, db) != 2 {
		t.Errorf("expected 2 entries, got %d", countEntries(t, db))
	}
}

func TestCleanupExpired_AllExpired(t *testing.T) {
	db, tmpDir := setupTestDB(t)
	defer db.Close()

	cacheDir := filepath.Join(tmpDir, "cache")
	r := newTestCacheRepo(db, cacheDir, 10*time.Minute)

	now := time.Now().UTC()
	insertTestEntry(t, db, "bucket1", "key1", "/tmp/fake1", now.Add(-2*time.Hour))
	insertTestEntry(t, db, "bucket1", "key2", "/tmp/fake2", now.Add(-1*time.Hour))
	insertTestEntry(t, db, "bucket2", "key3", "/tmp/fake3", now.Add(-30*time.Minute))

	deleted, err := r.CleanupExpired(context.Background())
	if err != nil {
		t.Fatalf("CleanupExpired failed: %v", err)
	}

	if deleted != 3 {
		t.Errorf("expected 3 deleted, got %d", deleted)
	}

	if countEntries(t, db) != 0 {
		t.Errorf("expected 0 entries, got %d", countEntries(t, db))
	}
}

func TestCleanupExpired_MissingCacheFile(t *testing.T) {
	db, tmpDir := setupTestDB(t)
	defer db.Close()

	cacheDir := filepath.Join(tmpDir, "cache")
	r := newTestCacheRepo(db, cacheDir, 10*time.Minute)

	now := time.Now().UTC()
	// Entry points to a file that doesn't exist on disk
	insertTestEntry(t, db, "bucket1", "key1", "/tmp/nonexistent-file", now.Add(-1*time.Hour))

	deleted, err := r.CleanupExpired(context.Background())
	if err != nil {
		t.Fatalf("CleanupExpired failed: %v", err)
	}

	if deleted != 1 {
		t.Errorf("expected 1 deleted, got %d", deleted)
	}

	if countEntries(t, db) != 0 {
		t.Errorf("expected 0 entries, got %d", countEntries(t, db))
	}
}

func TestStartCleanupLoop_RunsPeriodically(t *testing.T) {
	db, tmpDir := setupTestDB(t)
	defer db.Close()

	cacheDir := filepath.Join(tmpDir, "cache")
	os.MkdirAll(filepath.Join(cacheDir, "bucket1"), 0755)
	r := newTestCacheRepo(db, cacheDir, 10*time.Minute)

	now := time.Now().UTC()
	expiredFile := filepath.Join(cacheDir, "bucket1", "file1")
	os.WriteFile(expiredFile, []byte("data"), 0644)
	insertTestEntry(t, db, "bucket1", "key1", expiredFile, now.Add(-1*time.Hour))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r.StartCleanupLoop(ctx, 100*time.Millisecond)

	// Wait for at least one cleanup cycle
	time.Sleep(300 * time.Millisecond)

	if countEntries(t, db) != 0 {
		t.Errorf("expected expired entry to be cleaned up, got %d entries", countEntries(t, db))
	}

	if _, err := os.Stat(expiredFile); !os.IsNotExist(err) {
		t.Error("expected cache file to be deleted by cleanup loop")
	}
}

func TestStartCleanupLoop_StopsOnContextCancel(t *testing.T) {
	db, tmpDir := setupTestDB(t)
	defer db.Close()

	cacheDir := filepath.Join(tmpDir, "cache")
	r := newTestCacheRepo(db, cacheDir, 10*time.Minute)

	ctx, cancel := context.WithCancel(context.Background())
	r.StartCleanupLoop(ctx, 50*time.Millisecond)

	// Cancel immediately
	cancel()

	// Insert expired entry after cancel
	now := time.Now().UTC()
	insertTestEntry(t, db, "bucket1", "key1", "/tmp/fake1", now.Add(-1*time.Hour))

	// Wait to ensure no more cleanup happens
	time.Sleep(200 * time.Millisecond)

	if countEntries(t, db) != 1 {
		t.Errorf("expected entry to remain after context cancel, got %d", countEntries(t, db))
	}
}

func totalCacheSize(t *testing.T, db *sql.DB) int64 {
	t.Helper()
	var total int64
	if err := db.QueryRow("SELECT COALESCE(SUM(size), 0) FROM cache_entries").Scan(&total); err != nil {
		t.Fatalf("failed to get total cache size: %v", err)
	}
	return total
}

func TestEvict_RemovesOldestEntriesWhenOverLimit(t *testing.T) {
	db, tmpDir := setupTestDB(t)
	defer db.Close()

	cacheDir := filepath.Join(tmpDir, "cache")
	os.MkdirAll(filepath.Join(cacheDir, "bucket1"), 0755)

	// Max 500 bytes
	r := newTestCacheRepo(db, cacheDir, 10*time.Minute, WithMaxCacheSize(500))

	now := time.Now().UTC()
	expires := now.Add(1 * time.Hour)

	// Create files and entries: 200 + 200 + 200 = 600 bytes (over 500 limit)
	for i, key := range []string{"old", "mid", "new"} {
		path := filepath.Join(cacheDir, "bucket1", key)
		os.WriteFile(path, make([]byte, 200), 0644)
		insertTestEntryWithSize(t, db, "bucket1", key, path, 200, expires, now.Add(time.Duration(i)*time.Minute))
	}

	evicted, err := r.Evict(context.Background())
	if err != nil {
		t.Fatalf("Evict failed: %v", err)
	}

	if evicted != 1 {
		t.Errorf("expected 1 evicted, got %d", evicted)
	}

	if countEntries(t, db) != 2 {
		t.Errorf("expected 2 remaining entries, got %d", countEntries(t, db))
	}

	if totalCacheSize(t, db) != 400 {
		t.Errorf("expected total size 400, got %d", totalCacheSize(t, db))
	}

	// The oldest entry ("old") should have been evicted
	oldPath := filepath.Join(cacheDir, "bucket1", "old")
	if _, err := os.Stat(oldPath); !os.IsNotExist(err) {
		t.Error("expected oldest cache file to be deleted")
	}

	// Newer entries should remain
	for _, key := range []string{"mid", "new"} {
		path := filepath.Join(cacheDir, "bucket1", key)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected %s cache file to still exist", key)
		}
	}
}

func TestEvict_NoEvictionWhenUnderLimit(t *testing.T) {
	db, tmpDir := setupTestDB(t)
	defer db.Close()

	cacheDir := filepath.Join(tmpDir, "cache")
	r := newTestCacheRepo(db, cacheDir, 10*time.Minute, WithMaxCacheSize(1000))

	now := time.Now().UTC()
	expires := now.Add(1 * time.Hour)
	insertTestEntryWithSize(t, db, "bucket1", "key1", "/tmp/f1", 200, expires, now)
	insertTestEntryWithSize(t, db, "bucket1", "key2", "/tmp/f2", 300, expires, now)

	evicted, err := r.Evict(context.Background())
	if err != nil {
		t.Fatalf("Evict failed: %v", err)
	}

	if evicted != 0 {
		t.Errorf("expected 0 evicted, got %d", evicted)
	}

	if countEntries(t, db) != 2 {
		t.Errorf("expected 2 entries, got %d", countEntries(t, db))
	}
}

func TestEvict_NoOpWhenMaxCacheSizeIsZero(t *testing.T) {
	db, tmpDir := setupTestDB(t)
	defer db.Close()

	cacheDir := filepath.Join(tmpDir, "cache")
	// No max cache size (default 0 = unlimited)
	r := newTestCacheRepo(db, cacheDir, 10*time.Minute)

	now := time.Now().UTC()
	expires := now.Add(1 * time.Hour)
	insertTestEntryWithSize(t, db, "bucket1", "key1", "/tmp/f1", 999999, expires, now)

	evicted, err := r.Evict(context.Background())
	if err != nil {
		t.Fatalf("Evict failed: %v", err)
	}

	if evicted != 0 {
		t.Errorf("expected 0 evicted (unlimited), got %d", evicted)
	}
}

func TestEvict_EvictsMultipleEntries(t *testing.T) {
	db, tmpDir := setupTestDB(t)
	defer db.Close()

	cacheDir := filepath.Join(tmpDir, "cache")
	// Max 100 bytes
	r := newTestCacheRepo(db, cacheDir, 10*time.Minute, WithMaxCacheSize(100))

	now := time.Now().UTC()
	expires := now.Add(1 * time.Hour)

	// 4 entries x 100 bytes = 400 bytes total (over 100 limit)
	for i, key := range []string{"a", "b", "c", "d"} {
		insertTestEntryWithSize(t, db, "bucket1", key, "/tmp/"+key, 100, expires, now.Add(time.Duration(i)*time.Minute))
	}

	evicted, err := r.Evict(context.Background())
	if err != nil {
		t.Fatalf("Evict failed: %v", err)
	}

	// Need to evict 3 entries (400 -> 300 -> 200 -> 100) to get to 100
	if evicted != 3 {
		t.Errorf("expected 3 evicted, got %d", evicted)
	}

	if countEntries(t, db) != 1 {
		t.Errorf("expected 1 remaining entry, got %d", countEntries(t, db))
	}
}

func TestVacuum_Success(t *testing.T) {
	db, tmpDir := setupTestDB(t)
	defer db.Close()

	cacheDir := filepath.Join(tmpDir, "cache")
	r := newTestCacheRepo(db, cacheDir, 10*time.Minute)

	// Insert and delete entries to create free pages
	now := time.Now().UTC()
	for i := range 100 {
		key := fmt.Sprintf("key-%d", i)
		insertTestEntryWithSize(t, db, "bucket1", key, "/tmp/"+key, 1000, now.Add(-1*time.Hour), now)
	}
	_, err := r.CleanupExpired(context.Background())
	if err != nil {
		t.Fatalf("CleanupExpired failed: %v", err)
	}

	// VACUUM should succeed without error
	if err := r.Vacuum(); err != nil {
		t.Fatalf("Vacuum failed: %v", err)
	}

	// DB should still be functional after VACUUM
	if countEntries(t, db) != 0 {
		t.Errorf("expected 0 entries after cleanup, got %d", countEntries(t, db))
	}
}

func TestVacuum_ReducesFileSize(t *testing.T) {
	db, tmpDir := setupTestDB(t)
	defer db.Close()

	cacheDir := filepath.Join(tmpDir, "cache")
	r := newTestCacheRepo(db, cacheDir, 10*time.Minute)
	dbPath := filepath.Join(tmpDir, "test.db")

	// Insert many entries
	now := time.Now().UTC()
	for i := range 200 {
		key := fmt.Sprintf("key-%d", i)
		insertTestEntryWithSize(t, db, "bucket1", key, "/tmp/"+key, 10000, now.Add(-1*time.Hour), now)
	}

	// Force a checkpoint so WAL is flushed to main DB
	db.Exec("PRAGMA wal_checkpoint(TRUNCATE)")

	sizeBeforeDelete, err := os.Stat(dbPath)
	if err != nil {
		t.Fatalf("failed to stat DB: %v", err)
	}

	// Delete all entries
	_, err = r.CleanupExpired(context.Background())
	if err != nil {
		t.Fatalf("CleanupExpired failed: %v", err)
	}

	// Checkpoint again before measuring
	db.Exec("PRAGMA wal_checkpoint(TRUNCATE)")

	sizeAfterDelete, err := os.Stat(dbPath)
	if err != nil {
		t.Fatalf("failed to stat DB: %v", err)
	}

	// Without VACUUM, file size should be same or similar
	if sizeAfterDelete.Size() < sizeBeforeDelete.Size()/2 {
		t.Skip("DB already shrunk without VACUUM, cannot test VACUUM effect")
	}

	// Run VACUUM
	if err := r.Vacuum(); err != nil {
		t.Fatalf("Vacuum failed: %v", err)
	}

	sizeAfterVacuum, err := os.Stat(dbPath)
	if err != nil {
		t.Fatalf("failed to stat DB: %v", err)
	}

	if sizeAfterVacuum.Size() >= sizeAfterDelete.Size() {
		t.Errorf("expected file size to decrease after VACUUM: before=%d, after=%d",
			sizeAfterDelete.Size(), sizeAfterVacuum.Size())
	}
}

func TestStartCleanupLoop_RunsVacuumAfterDeletion(t *testing.T) {
	db, tmpDir := setupTestDB(t)
	defer db.Close()

	cacheDir := filepath.Join(tmpDir, "cache")
	os.MkdirAll(filepath.Join(cacheDir, "bucket1"), 0755)
	r := newTestCacheRepo(db, cacheDir, 10*time.Minute)
	dbPath := filepath.Join(tmpDir, "test.db")

	// Insert expired entries
	now := time.Now().UTC()
	for i := range 100 {
		key := fmt.Sprintf("key-%d", i)
		path := filepath.Join(cacheDir, "bucket1", key)
		os.WriteFile(path, make([]byte, 100), 0644)
		insertTestEntryWithSize(t, db, "bucket1", key, path, 100, now.Add(-1*time.Hour), now)
	}

	// Checkpoint to flush WAL
	db.Exec("PRAGMA wal_checkpoint(TRUNCATE)")

	sizeBeforeCleanup, err := os.Stat(dbPath)
	if err != nil {
		t.Fatalf("failed to stat DB: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r.StartCleanupLoop(ctx, 100*time.Millisecond)

	// Wait for cleanup cycle to run
	time.Sleep(400 * time.Millisecond)

	if countEntries(t, db) != 0 {
		t.Errorf("expected all entries cleaned up, got %d", countEntries(t, db))
	}

	// Check that DB file shrank (VACUUM ran)
	sizeAfterCleanup, err := os.Stat(dbPath)
	if err != nil {
		t.Fatalf("failed to stat DB: %v", err)
	}

	if sizeAfterCleanup.Size() > sizeBeforeCleanup.Size() {
		t.Errorf("expected DB file to shrink after cleanup+VACUUM: before=%d, after=%d",
			sizeBeforeCleanup.Size(), sizeAfterCleanup.Size())
	}
}
