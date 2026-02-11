package serviceif

import (
	"context"
	"io"

	"github.com/pkg/errors"
)

var ErrObjectAlreadyExists = errors.New("object already exists")

// ProgressCallback はアップロード進捗のコールバック関数型。
// nilの場合、進捗追跡は行われない。
type ProgressCallback func(bytesProcessed int64)

type UploadRepository interface {
	PutObject(ctx context.Context, bucketName, key, contentType string, body io.ReadSeeker) (string, error)
	PutObjectIfNotExists(ctx context.Context, bucketName, key, contentType string, body io.ReadSeeker) (string, error)
}

type UploadResult struct {
	Key  string `json:"key"`
	Size int64  `json:"size"`
	ETag string `json:"etag"`
}

type UploadService interface {
	UploadObject(ctx context.Context, bucketName, key, contentType string, body io.Reader, size int64, overwrite bool, onProgress ProgressCallback) (*UploadResult, error)
	CreateDirectory(ctx context.Context, bucketName, path string) (*UploadResult, error)
}
