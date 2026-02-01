package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	appconfig "r2manager/config"
	"r2manager/infrastructure"
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

	// DI wiring
	h := wireHandlers(s3Client, db, cacheCfg)

	// Start server
	r := router.NewRouter(h.buckets, h.objects, h.content)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
