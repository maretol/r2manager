package repository

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"

	"r2manager/domain"
)

type ContentRepository struct {
	client *s3.Client
}

func NewContentRepository(client *s3.Client) *ContentRepository {
	return &ContentRepository{client: client}
}

func (r *ContentRepository) GetContent(ctx context.Context, bucketName, objectKey string) (*domain.ObjectContent, error) {
	output, err := r.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to GetObject")
	}

	contentType := "application/octet-stream"
	if output.ContentType != nil {
		contentType = *output.ContentType
	}

	var size int64
	if output.ContentLength != nil {
		size = *output.ContentLength
	}

	etag := ""
	if output.ETag != nil {
		etag = *output.ETag
	}

	return &domain.ObjectContent{
		Body:        output.Body,
		ContentType: contentType,
		Size:        size,
		ETag:        etag,
	}, nil
}
