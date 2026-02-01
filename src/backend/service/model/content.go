package service

import (
	"context"

	"github.com/pkg/errors"

	"r2manager/domain"
	serviceif "r2manager/service/interface"
)

type ContentService struct {
	contentRepo serviceif.ContentRepository
	cacheRepo   serviceif.CacheRepository
}

func NewContentService(contentRepo serviceif.ContentRepository, cacheRepo serviceif.CacheRepository) *ContentService {
	return &ContentService{
		contentRepo: contentRepo,
		cacheRepo:   cacheRepo,
	}
}

func (s *ContentService) GetContent(ctx context.Context, bucketName, objectKey string) (*domain.ObjectContent, error) {
	entry, err := s.cacheRepo.Lookup(ctx, bucketName, objectKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup cache")
	}

	if entry != nil {
		body, err := s.cacheRepo.OpenCacheFile(entry.CachePath)
		if err == nil {
			return &domain.ObjectContent{
				Body:        body,
				ContentType: entry.ContentType,
				Size:        entry.Size,
				ETag:        entry.ETag,
				CacheHit:    true,
			}, nil
		}
		// Cache file missing on disk; fall through to fetch from S3
	}

	content, err := s.contentRepo.GetContent(ctx, bucketName, objectKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get content from S3")
	}

	entry, err = s.cacheRepo.Store(ctx, bucketName, objectKey, content.Body, content.ContentType, content.Size, content.ETag)
	content.Body.Close()
	if err != nil {
		return nil, errors.Wrap(err, "failed to store cache")
	}

	body, err := s.cacheRepo.OpenCacheFile(entry.CachePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open cached file after store")
	}

	return &domain.ObjectContent{
		Body:        body,
		ContentType: entry.ContentType,
		Size:        entry.Size,
		ETag:        entry.ETag,
		CacheHit:    false,
	}, nil
}
