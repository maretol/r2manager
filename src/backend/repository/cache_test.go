package repository

import (
	"context"
	"database/sql"
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

func newTestCacheRepo(db *sql.DB, cacheDir string, ttl time.Duration) *CacheRepository {
	// Reset singleton for each test
	once = sync.Once{}
	repo = nil
	return NewCacheRepository(db, cacheDir, ttl)
}

func insertTestEntry(t *testing.T, db *sql.DB, bucketName, objectKey, cachePath string, expiresAt time.Time) {
	t.Helper()
	_, err := db.Exec(
		`INSERT INTO cache_entries (bucket_name, object_key, content_type, size, etag, cache_path, cached_at, expires_at)
		 VALUES (?, ?, 'text/plain', 100, 'etag1', ?, ?, ?)`,
		bucketName, objectKey, cachePath, time.Now().UTC(), expiresAt,
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
