package di

import (
	"database/sql"
	appconfig "r2manager/config"
	"r2manager/handler"
	"r2manager/repository"
	service "r2manager/service/model"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func CreateObjectsHandler(s3Client *s3.Client, db *sql.DB, cacheCfg *appconfig.CacheConfig, listCache *repository.ListCacheRepository) *handler.ObjectsHandler {
	var opts []repository.CacheOption
	if cacheCfg.MaxCacheSize > 0 {
		opts = append(opts, repository.WithMaxCacheSize(cacheCfg.MaxCacheSize))
	}
	cacheRepo := repository.NewCacheRepository(db, cacheCfg.CacheDir, cacheCfg.TTL, opts...)

	objectRepo := repository.NewObjectRepository(s3Client)
	objectService := service.NewObjectService(objectRepo, cacheRepo, listCache)
	objectsHandler := handler.NewObjectsHandler(objectService)

	return objectsHandler
}
