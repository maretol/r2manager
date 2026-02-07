package di

import (
	"database/sql"

	appconfig "r2manager/config"
	"r2manager/handler"
	"r2manager/repository"
)

func CreateCacheHandler(db *sql.DB, cacheCfg *appconfig.CacheConfig, listCache *repository.ListCacheRepository) *handler.CacheHandler {
	cacheRepo := repository.NewCacheRepository(db, cacheCfg.CacheDir, cacheCfg.TTL)
	return handler.NewCacheHandler(cacheRepo, listCache)
}
