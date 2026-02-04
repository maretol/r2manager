package di

import (
	"r2manager/handler"
	"r2manager/repository"
	service "r2manager/service/model"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func CreateBucketsHandler(s3Client *s3.Client, listCache *repository.ListCacheRepository) *handler.BucketsHandler {
	bucketRepo := repository.NewBucketRepository(s3Client)
	bucketService := service.NewBucketService(bucketRepo, listCache)
	bucketsHandler := handler.NewBucketsHandler(bucketService)

	return bucketsHandler
}
