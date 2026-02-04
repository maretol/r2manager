package repository

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"

	"r2manager/domain"
	serviceif "r2manager/service/interface"
)

type ObjectRepository struct {
	client *s3.Client
}

func NewObjectRepository(client *s3.Client) *ObjectRepository {
	return &ObjectRepository{client: client}
}

func (r *ObjectRepository) GetObjects(ctx context.Context, bucketName string, params serviceif.ListObjectsParams) (*domain.ListObjectsResult, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	}
	if params.Prefix != "" {
		input.Prefix = aws.String(params.Prefix)
	}
	if params.Delimiter != "" {
		input.Delimiter = aws.String(params.Delimiter)
	}

	output, err := r.client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ListObjectsV2")
	}

	// Combine Contents and CommonPrefixes into objects
	objects := make([]domain.Object, 0, len(output.Contents)+len(output.CommonPrefixes))

	// Add regular objects from Contents
	for _, o := range output.Contents {
		obj := domain.Object{}
		if o.Key != nil {
			obj.Key = *o.Key
		}
		if o.Size != nil {
			obj.Size = *o.Size
		}
		if o.LastModified != nil {
			obj.LastModified = *o.LastModified
		}
		if o.ETag != nil {
			obj.ETag = *o.ETag
		}
		objects = append(objects, obj)
	}

	// Add folder markers from CommonPrefixes
	for _, cp := range output.CommonPrefixes {
		if cp.Prefix != nil {
			objects = append(objects, domain.Object{
				Key: *cp.Prefix,
			})
		}
	}

	result := &domain.ListObjectsResult{
		Objects:   objects,
		Prefix:    params.Prefix,
		Delimiter: params.Delimiter,
	}

	if output.IsTruncated != nil {
		result.IsTruncated = *output.IsTruncated
	}
	if output.NextContinuationToken != nil {
		result.NextContinuationToken = *output.NextContinuationToken
	}

	return result, nil
}
