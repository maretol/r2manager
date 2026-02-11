package service

import (
	"bytes"
	"context"
	"io"
	"log"
	"strings"
	"time"

	"github.com/pkg/errors"

	"r2manager/progress"
	serviceif "r2manager/service/interface"
)

type UploadService struct {
	repo      serviceif.UploadRepository
	listCache serviceif.ListCacheRepository
}

func NewUploadService(repo serviceif.UploadRepository, listCache serviceif.ListCacheRepository) *UploadService {
	return &UploadService{repo: repo, listCache: listCache}
}

func (s *UploadService) UploadObject(ctx context.Context, bucketName, key, contentType string, body io.Reader, size int64, overwrite bool, onProgress serviceif.ProgressCallback) (*serviceif.UploadResult, error) {
	key = sanitizeObjectPath(key)
	if key == "" {
		return nil, errors.New("invalid key")
	}

	// リクエストボディを一度バッファに読み込み、io.ReadSeeker として渡すことで
	// SDK がリトライ時にボディを巻き戻せるようにする
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, body); err != nil {
		return nil, errors.Wrap(err, "failed to read request body")
	}

	var reader io.ReadSeeker
	baseReader := bytes.NewReader(buf.Bytes())
	if onProgress != nil {
		progressReader, err := progress.NewProgressReadSeeker(baseReader, onProgress, 100*time.Millisecond)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create progress reader")
		}
		reader = progressReader
	} else {
		reader = baseReader
	}

	var etag string
	var putErr error
	if overwrite {
		etag, putErr = s.repo.PutObject(ctx, bucketName, key, contentType, reader)
	} else {
		etag, putErr = s.repo.PutObjectIfNotExists(ctx, bucketName, key, contentType, reader)
	}
	if putErr != nil {
		return nil, errors.Wrap(putErr, "failed to upload object")
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
	path = sanitizeObjectPath(path)
	if path == "" {
		return nil, errors.New("invalid path")
	}
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

// sanitizeObjectPath はオブジェクトキーのパスを正規化・検証する。
// 先頭スラッシュの除去、連続スラッシュの正規化、パストラバーサルの排除を行い、
// 不正なパスの場合は空文字を返す。
func sanitizeObjectPath(path string) string {
	path = strings.TrimSpace(path)
	// 先頭のスラッシュを除去
	path = strings.TrimLeft(path, "/")
	if path == "" {
		return ""
	}

	// 連続スラッシュを正規化
	for strings.Contains(path, "//") {
		path = strings.ReplaceAll(path, "//", "/")
	}

	// パストラバーサルを含むパスを拒否
	for segment := range strings.SplitSeq(strings.TrimSuffix(path, "/"), "/") {
		if segment == ".." || segment == "." {
			return ""
		}
	}

	return path
}
