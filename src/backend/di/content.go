package di

import (
	"database/sql"
	appconfig "r2manager/config"
	"r2manager/handler"
	"r2manager/repository"
	service "r2manager/service/model"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func CreateContentHandler(s3Client *s3.Client, db *sql.DB, cacheCfg *appconfig.CacheConfig) *handler.ContentHandler {
	cacheRepo := repository.NewCacheRepository(db, cacheCfg.CacheDir, cacheCfg.TTL)

	contentRepo := repository.NewContentRepository(s3Client)
	contentService := service.NewContentService(contentRepo, cacheRepo)
	contentHandler := handler.NewContentHandler(contentService)

	return contentHandler
}
