package service

import (
	"context"

	"r2manager/domain"
	serviceif "r2manager/service/interface"
)

type BucketService struct {
	repo serviceif.BucketRepository
}

func NewBucketService(repo serviceif.BucketRepository) *BucketService {
	return &BucketService{repo: repo}
}

func (s *BucketService) GetBuckets(ctx context.Context) ([]domain.Bucket, error) {
	return s.repo.GetBuckets(ctx)
}
