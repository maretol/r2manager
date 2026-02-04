package service

import (
	"context"
	"log"

	"r2manager/domain"
	serviceif "r2manager/service/interface"
)

type ObjectService struct {
	repo      serviceif.ObjectRepository
	cacheRepo serviceif.CacheRepository
	listCache serviceif.ListCacheRepository
}

func NewObjectService(repo serviceif.ObjectRepository, cacheRepo serviceif.CacheRepository, listCache serviceif.ListCacheRepository) *ObjectService {
	return &ObjectService{repo: repo, cacheRepo: cacheRepo, listCache: listCache}
}

func (s *ObjectService) GetObjects(ctx context.Context, bucketName string, params serviceif.ListObjectsParams) (*domain.ListObjectsResult, error) {
	cacheKey := bucketName + ":" + params.Prefix

	// Check list cache first
	if result, found := s.listCache.GetObjects(bucketName, params.Prefix); found {
		log.Printf("cache hit: objects [%s] (%d items)", cacheKey, len(result.Objects))
		return result, nil
	}

	log.Printf("cache miss: objects [%s]", cacheKey)

	// Fetch from R2
	result, err := s.repo.GetObjects(ctx, bucketName, params)
	if err != nil {
		return nil, err
	}

	// Store in list cache
	s.listCache.SetObjects(bucketName, params.Prefix, result)
	log.Printf("cache stored: objects [%s] (%d items)", cacheKey, len(result.Objects))

	// Build ETag map and invalidate stale content cache entries
	currentETags := make(map[string]string, len(result.Objects))
	for _, obj := range result.Objects {
		if obj.ETag != "" {
			currentETags[obj.Key] = obj.ETag
		}
	}

	if len(currentETags) > 0 {
		invalidated, err := s.cacheRepo.InvalidateByETags(ctx, bucketName, currentETags)
		if err != nil {
			log.Printf("warning: failed to invalidate cache by ETags: %v", err)
		} else if invalidated > 0 {
			log.Printf("invalidated %d stale cache entries for bucket %s", invalidated, bucketName)
		}
	}

	return result, nil
}
