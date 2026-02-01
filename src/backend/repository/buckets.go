package repository

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"

	"r2manager/domain"
)

type BucketRepository struct {
	client *s3.Client
}

func NewBucketRepository(client *s3.Client) *BucketRepository {
	return &BucketRepository{client: client}
}

func (r *BucketRepository) GetBuckets(ctx context.Context) ([]domain.Bucket, error) {
	output, err := r.client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to ListBuckets")
	}

	buckets := make([]domain.Bucket, 0, len(output.Buckets))
	for _, b := range output.Buckets {
		bucket := domain.Bucket{}
		if b.Name != nil {
			bucket.Name = *b.Name
		}
		if b.CreationDate != nil {
			bucket.CreationDate = *b.CreationDate
		}
		buckets = append(buckets, bucket)
	}

	return buckets, nil
}
