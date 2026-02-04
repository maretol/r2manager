package serviceif

import (
	"context"

	"r2manager/domain"
)

type BucketRepository interface {
	GetBuckets(ctx context.Context) ([]domain.Bucket, error)
}

type BucketService interface {
	GetBuckets(ctx context.Context) ([]domain.Bucket, error)
}

type ListCacheRepository interface {
	GetBuckets() ([]domain.Bucket, bool)
	SetBuckets(buckets []domain.Bucket)
	InvalidateBuckets()
	GetObjects(bucketName, prefix string) (*domain.ListObjectsResult, bool)
	SetObjects(bucketName, prefix string, result *domain.ListObjectsResult)
	InvalidateObjects(bucketName string)
	InvalidateAll()
}
