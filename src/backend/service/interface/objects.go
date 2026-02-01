package serviceif

import (
	"context"

	"r2manager/domain"
)

type ObjectRepository interface {
	GetObjects(ctx context.Context, bucketName string) ([]domain.Object, error)
}

type ObjectService interface {
	GetObjects(ctx context.Context, bucketName string) ([]domain.Object, error)
}
