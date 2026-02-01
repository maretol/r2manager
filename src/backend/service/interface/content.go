package serviceif

import (
	"context"
	"io"

	"r2manager/domain"
)

type ContentRepository interface {
	GetContent(ctx context.Context, bucketName, objectKey string) (*domain.ObjectContent, error)
}

type CacheRepository interface {
	Lookup(ctx context.Context, bucketName, objectKey string) (*domain.CacheEntry, error)
	Store(ctx context.Context, bucketName, objectKey string, body io.Reader, contentType string, size int64, etag string) (*domain.CacheEntry, error)
	OpenCacheFile(cachePath string) (io.ReadCloser, error)
	InvalidateByETags(ctx context.Context, bucketName string, currentETags map[string]string) (int, error)
}

type ContentService interface {
	GetContent(ctx context.Context, bucketName, objectKey string) (*domain.ObjectContent, error)
}
