package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	appconfig "r2manager/config"
	"r2manager/di"
	"r2manager/infrastructure"
	"r2manager/progress"
	"r2manager/repository"
	"r2manager/router"
)

func main() {
	r2cfg, err := appconfig.LoadR2ConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	s3Client, err := appconfig.NewS3Client(context.Background(), r2cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Cache config and DB
	cacheCfg := appconfig.LoadCacheConfigFromEnv()

	if err := os.MkdirAll(filepath.Dir(cacheCfg.DBPath), 0755); err != nil {
		log.Fatalf("failed to create DB directory: %v", err)
	}
	if err := os.MkdirAll(cacheCfg.CacheDir, 0755); err != nil {
		log.Fatalf("failed to create cache directory: %v", err)
	}

	db, err := infrastructure.NewSQLiteDB(cacheCfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// List cache (shared between buckets and objects)
	listCache := repository.NewListCacheRepository()

	// Upload config
	uploadCfg := appconfig.LoadUploadConfigFromEnv()

	// Progress store
	progressStore := progress.NewUploadProgressStore()

	// DI wiring
	bh := di.CreateBucketsHandler(s3Client, listCache)
	oh := di.CreateObjectsHandler(s3Client, db, cacheCfg, listCache)
	ch := di.CreateContentHandler(s3Client, db, cacheCfg)
	cah := di.CreateCacheHandler(db, cacheCfg, listCache)
	sh := di.CreateSettingsHandler(db)
	uh := di.CreateUploadHandler(s3Client, listCache, uploadCfg, progressStore)
	uph := di.CreateUploadProgressHandler(progressStore)

	// Start background cache cleanup
	var opts []repository.CacheOption
	if cacheCfg.MaxCacheSize > 0 {
		opts = append(opts, repository.WithMaxCacheSize(cacheCfg.MaxCacheSize))
	}
	cacheRepo := repository.NewCacheRepository(db, cacheCfg.CacheDir, cacheCfg.TTL, opts...)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cacheRepo.StartCleanupLoop(ctx, cacheCfg.CleanupInterval)

	// Start progress store cleanup
	progressStore.StartCleanupLoop(ctx)

	// Start server
	r := router.NewRouter(bh, oh, ch, cah, sh, uh, uph)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
