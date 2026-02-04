package serviceif

import (
	"context"

	"r2manager/domain"
)

type ListObjectsParams struct {
	Prefix    string
	Delimiter string
}

type ObjectRepository interface {
	GetObjects(ctx context.Context, bucketName string, params ListObjectsParams) (*domain.ListObjectsResult, error)
}

type ObjectService interface {
	GetObjects(ctx context.Context, bucketName string, params ListObjectsParams) (*domain.ListObjectsResult, error)
}
