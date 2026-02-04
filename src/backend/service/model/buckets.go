package service

import (
	"context"
	"log"

	"r2manager/domain"
	serviceif "r2manager/service/interface"
)

type BucketService struct {
	repo      serviceif.BucketRepository
	listCache serviceif.ListCacheRepository
}

func NewBucketService(repo serviceif.BucketRepository, listCache serviceif.ListCacheRepository) *BucketService {
	return &BucketService{repo: repo, listCache: listCache}
}

func (s *BucketService) GetBuckets(ctx context.Context) ([]domain.Bucket, error) {
	// Check cache first
	if buckets, found := s.listCache.GetBuckets(); found {
		log.Printf("cache hit: buckets (%d items)", len(buckets))
		return buckets, nil
	}

	log.Printf("cache miss: buckets")

	// Fetch from R2
	buckets, err := s.repo.GetBuckets(ctx)
	if err != nil {
		return nil, err
	}

	// Store in cache
	s.listCache.SetBuckets(buckets)
	log.Printf("cache stored: buckets (%d items)", len(buckets))

	return buckets, nil
}
