package serviceif

import (
	"context"

	"r2manager/domain"
)

type SettingsRepository interface {
	GetAllBucketSettings(ctx context.Context) ([]domain.BucketSettings, error)
	GetBucketSettings(ctx context.Context, bucketName string) (*domain.BucketSettings, error)
	UpsertBucketSettings(ctx context.Context, bucketName, publicUrl string) error
}

type SettingsService interface {
	GetAllBucketSettings(ctx context.Context) ([]domain.BucketSettings, error)
	GetBucketSettings(ctx context.Context, bucketName string) (*domain.BucketSettings, error)
	UpdateBucketPublicUrl(ctx context.Context, bucketName, publicUrl string) error
}
