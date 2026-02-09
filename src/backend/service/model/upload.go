package service

import (
	"context"
	"io"
	"log"
	"strings"

	"github.com/pkg/errors"

	serviceif "r2manager/service/interface"
)

type UploadService struct {
	repo      serviceif.UploadRepository
	listCache serviceif.ListCacheRepository
}

func NewUploadService(repo serviceif.UploadRepository, listCache serviceif.ListCacheRepository) *UploadService {
	return &UploadService{repo: repo, listCache: listCache}
}

func (s *UploadService) UploadObject(ctx context.Context, bucketName, key, contentType string, body io.Reader, size int64, overwrite bool) (*serviceif.UploadResult, error) {
	if !overwrite {
		exists, err := s.repo.HeadObject(ctx, bucketName, key)
		if err != nil {
			return nil, errors.Wrap(err, "failed to check object existence")
		}
		if exists {
			return nil, serviceif.ErrObjectAlreadyExists
		}
	}

	etag, err := s.repo.PutObject(ctx, bucketName, key, contentType, body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to upload object")
	}

	// Invalidate list cache for this bucket
	s.listCache.InvalidateObjects(bucketName)
	log.Printf("uploaded object: bucket=%s key=%s size=%d", bucketName, key, size)

	return &serviceif.UploadResult{
		Key:  key,
		Size: size,
		ETag: etag,
	}, nil
}

func (s *UploadService) CreateDirectory(ctx context.Context, bucketName, path string) (*serviceif.UploadResult, error) {
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	etag, err := s.repo.PutObject(ctx, bucketName, path, "application/x-directory", strings.NewReader(""))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create directory")
	}

	s.listCache.InvalidateObjects(bucketName)
	log.Printf("created directory: bucket=%s path=%s", bucketName, path)

	return &serviceif.UploadResult{
		Key:  path,
		Size: 0,
		ETag: etag,
	}, nil
}

func (s *UploadService) ObjectExists(ctx context.Context, bucketName, key string) (bool, error) {
	return s.repo.HeadObject(ctx, bucketName, key)
}
