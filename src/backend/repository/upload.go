package repository

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
)

type UploadRepository struct {
	client *s3.Client
}

func NewUploadRepository(client *s3.Client) *UploadRepository {
	return &UploadRepository{client: client}
}

func (r *UploadRepository) PutObject(ctx context.Context, bucketName, key, contentType string, body io.Reader) (string, error) {
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

func (r *UploadRepository) HeadObject(ctx context.Context, bucketName, key string) (bool, error) {
	_, err := r.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		// R2/S3-compatible: check HTTP 404 status code for non-existent objects
		var respErr interface{ HTTPStatusCode() int }
		if errors.As(err, &respErr) && respErr.HTTPStatusCode() == 404 {
			return false, nil
		}
		return false, errors.Wrap(err, "failed to HeadObject")
	}
	return true, nil
}
