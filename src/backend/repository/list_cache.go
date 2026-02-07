package repository

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"

	"r2manager/domain"
)

const (
	BucketsCacheTTL = 60 * time.Minute
	ObjectsCacheTTL = 10 * time.Minute
	CleanupInterval = 30 * time.Minute
)

type ListCacheRepository struct {
	buckets *cache.Cache
	objects *cache.Cache
}

func NewListCacheRepository() *ListCacheRepository {
	return &ListCacheRepository{
		buckets: cache.New(BucketsCacheTTL, CleanupInterval),
		objects: cache.New(ObjectsCacheTTL, CleanupInterval),
	}
}

const bucketsCacheKey = "buckets"

func (r *ListCacheRepository) GetBuckets() ([]domain.Bucket, bool) {
	if val, found := r.buckets.Get(bucketsCacheKey); found {
		if buckets, ok := val.([]domain.Bucket); ok {
			return buckets, true
		}
	}
	return nil, false
}

func (r *ListCacheRepository) SetBuckets(buckets []domain.Bucket) {
	r.buckets.Set(bucketsCacheKey, buckets, cache.DefaultExpiration)
}

func (r *ListCacheRepository) InvalidateBuckets() {
	r.buckets.Delete(bucketsCacheKey)
}

func objectsCacheKey(bucketName, prefix string) string {
	return fmt.Sprintf("%s:%s", bucketName, prefix)
}

func (r *ListCacheRepository) GetObjects(bucketName, prefix string) (*domain.ListObjectsResult, bool) {
	key := objectsCacheKey(bucketName, prefix)
	if val, found := r.objects.Get(key); found {
		if result, ok := val.(*domain.ListObjectsResult); ok {
			return result, true
		}
	}
	return nil, false
}

func (r *ListCacheRepository) SetObjects(bucketName, prefix string, result *domain.ListObjectsResult) {
	key := objectsCacheKey(bucketName, prefix)
	r.objects.Set(key, result, cache.DefaultExpiration)
}

func (r *ListCacheRepository) InvalidateObjects(bucketName string) {
	// Invalidate all cached objects for a specific bucket
	for key := range r.objects.Items() {
		if len(key) > len(bucketName) && key[:len(bucketName)+1] == bucketName+":" {
			r.objects.Delete(key)
		}
	}
}

func (r *ListCacheRepository) InvalidateAllObjects() {
	r.objects.Flush()
}

func (r *ListCacheRepository) InvalidateAll() {
	r.buckets.Flush()
	r.objects.Flush()
}
