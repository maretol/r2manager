package repository

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"

	serviceif "r2manager/service/interface"
)

type UploadRepository struct {
	client *s3.Client
}

func NewUploadRepository(client *s3.Client) *UploadRepository {
	return &UploadRepository{client: client}
}

func (r *UploadRepository) PutObject(ctx context.Context, bucketName, key, contentType string, body io.ReadSeeker) (string, error) {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
		Body:        body,
	}

	output, err := r.client.PutObject(ctx, input)
	if err != nil {
		return "", errors.Wrap(err, "failed to PutObject")
	}

	etag := ""
	if output.ETag != nil {
		etag = *output.ETag
	}

	return etag, nil
}

func (r *UploadRepository) PutObjectIfNotExists(ctx context.Context, bucketName, key, contentType string, body io.ReadSeeker) (string, error) {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
		Body:        body,
		IfNoneMatch: aws.String("*"),
	}

	output, err := r.client.PutObject(ctx, input)
	if err != nil {
		// R2: 412 PreconditionFailed はオブジェクトが既に存在することを示す
		var respErr interface{ HTTPStatusCode() int }
		if errors.As(err, &respErr) && respErr.HTTPStatusCode() == 412 {
			return "", serviceif.ErrObjectAlreadyExists
		}
		return "", errors.Wrap(err, "failed to PutObject")
	}

	etag := ""
	if output.ETag != nil {
		etag = *output.ETag
	}

	return etag, nil
}

