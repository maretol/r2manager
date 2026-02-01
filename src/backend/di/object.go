package di

import (
	"database/sql"
	appconfig "r2manager/config"
	"r2manager/handler"
	"r2manager/repository"
	service "r2manager/service/model"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func CreateObjectsHandler(s3Client *s3.Client, db *sql.DB, cacheCfg *appconfig.CacheConfig) *handler.ObjectsHandler {
	cacheRepo := repository.NewCacheRepository(db, cacheCfg.CacheDir, cacheCfg.TTL)

	objectRepo := repository.NewObjectRepository(s3Client)
	objectService := service.NewObjectService(objectRepo, cacheRepo)
	objectsHandler := handler.NewObjectsHandler(objectService)

	return objectsHandler
}
