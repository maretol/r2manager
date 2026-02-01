package repository

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"

	"r2manager/domain"
)

type ObjectRepository struct {
	client *s3.Client
}

func NewObjectRepository(client *s3.Client) *ObjectRepository {
	return &ObjectRepository{client: client}
}

func (r *ObjectRepository) GetObjects(ctx context.Context, bucketName string) ([]domain.Object, error) {
	output, err := r.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to ListObjectsV2")
	}

	objects := make([]domain.Object, 0, len(output.Contents))
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

	return objects, nil
}
