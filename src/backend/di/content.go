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
	var opts []repository.CacheOption
	if cacheCfg.MaxCacheSize > 0 {
		opts = append(opts, repository.WithMaxCacheSize(cacheCfg.MaxCacheSize))
	}
	cacheRepo := repository.NewCacheRepository(db, cacheCfg.CacheDir, cacheCfg.TTL, opts...)

	contentRepo := repository.NewContentRepository(s3Client)
	contentService := service.NewContentService(contentRepo, cacheRepo)
	contentHandler := handler.NewContentHandler(contentService)

	return contentHandler
}
