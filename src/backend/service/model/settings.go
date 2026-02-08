package service

import (
	"context"

	"r2manager/domain"
	serviceif "r2manager/service/interface"
)

type SettingsService struct {
	repo serviceif.SettingsRepository
}

func NewSettingsService(repo serviceif.SettingsRepository) *SettingsService {
	return &SettingsService{repo: repo}
}

func (s *SettingsService) GetAllBucketSettings(ctx context.Context) ([]domain.BucketSettings, error) {
	return s.repo.GetAllBucketSettings(ctx)
}

func (s *SettingsService) GetBucketSettings(ctx context.Context, bucketName string) (*domain.BucketSettings, error) {
	return s.repo.GetBucketSettings(ctx, bucketName)
}

func (s *SettingsService) UpdateBucketPublicUrl(ctx context.Context, bucketName, publicUrl string) error {
	return s.repo.UpsertBucketSettings(ctx, bucketName, publicUrl)
}
