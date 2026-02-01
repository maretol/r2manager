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
