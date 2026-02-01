package domain

import "time"

type CacheEntry struct {
	BucketName  string
	ObjectKey   string
	ContentType string
	Size        int64
	ETag        string
	CachePath   string
	CachedAt    time.Time
	ExpiresAt   time.Time
}
