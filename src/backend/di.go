package main

import (
	"database/sql"

	"r2manager/handler"
	"r2manager/repository"
	service "r2manager/service/model"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	appconfig "r2manager/config"
)

type handlers struct {
	buckets *handler.BucketsHandler
	objects *handler.ObjectsHandler
	content *handler.ContentHandler
}

func wireHandlers(s3Client *s3.Client, db *sql.DB, cacheCfg *appconfig.CacheConfig) *handlers {
	bucketRepo := repository.NewBucketRepository(s3Client)
	bucketService := service.NewBucketService(bucketRepo)
	bucketsHandler := handler.NewBucketsHandler(bucketService)

	cacheRepo := repository.NewCacheRepository(db, cacheCfg.CacheDir, cacheCfg.TTL)

	objectRepo := repository.NewObjectRepository(s3Client)
	objectService := service.NewObjectService(objectRepo, cacheRepo)
	objectsHandler := handler.NewObjectsHandler(objectService)

	contentRepo := repository.NewContentRepository(s3Client)
	contentService := service.NewContentService(contentRepo, cacheRepo)
	contentHandler := handler.NewContentHandler(contentService)

	return &handlers{
		buckets: bucketsHandler,
		objects: objectsHandler,
		content: contentHandler,
	}
}
