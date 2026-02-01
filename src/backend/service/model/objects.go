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
}

func NewObjectService(repo serviceif.ObjectRepository, cacheRepo serviceif.CacheRepository) *ObjectService {
	return &ObjectService{repo: repo, cacheRepo: cacheRepo}
}

func (s *ObjectService) GetObjects(ctx context.Context, bucketName string) ([]domain.Object, error) {
	objects, err := s.repo.GetObjects(ctx, bucketName)
	if err != nil {
		return nil, err
	}

	// Build ETag map and invalidate stale cache entries
	currentETags := make(map[string]string, len(objects))
	for _, obj := range objects {
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

	return objects, nil
}
