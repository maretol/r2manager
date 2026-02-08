package infrastructure

import (
	"database/sql"

	"github.com/pkg/errors"

	_ "modernc.org/sqlite"
)

const schema = `
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
CREATE TABLE IF NOT EXISTS bucket_settings (
    bucket_name TEXT NOT NULL PRIMARY KEY,
    public_url  TEXT NOT NULL DEFAULT ''
);
`

func NewSQLiteDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open SQLite database")
	}

	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, errors.Wrap(err, "failed to set WAL mode")
	}

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, errors.Wrap(err, "failed to run schema migration")
	}

	return db, nil
}
