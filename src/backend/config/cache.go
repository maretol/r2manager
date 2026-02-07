package config

import (
	"os"
	"strconv"
	"time"
)

type CacheConfig struct {
	DBPath          string
	CacheDir        string
	TTL             time.Duration
	CleanupInterval time.Duration
	MaxCacheSize    int64
}

func LoadCacheConfigFromEnv() *CacheConfig {
	dbPath := os.Getenv("CACHE_DB_PATH")
	if dbPath == "" {
		dbPath = "./data/cache.db"
	}

	cacheDir := os.Getenv("CACHE_DIR")
	if cacheDir == "" {
		cacheDir = "./data/cache"
	}

	ttlMinutes := 120
	if v := os.Getenv("CACHE_TTL_MINUTES"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			ttlMinutes = parsed
		}
	}

	cleanupIntervalMinutes := 60
	if v := os.Getenv("CACHE_CLEANUP_INTERVAL_MINUTES"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			cleanupIntervalMinutes = parsed
		}
	}

	var maxCacheSize int64
	if v := os.Getenv("CACHE_MAX_SIZE_MB"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil && parsed > 0 {
			maxCacheSize = parsed * 1024 * 1024
		}
	}

	return &CacheConfig{
		DBPath:          dbPath,
		CacheDir:        cacheDir,
		TTL:             time.Duration(ttlMinutes) * time.Minute,
		CleanupInterval: time.Duration(cleanupIntervalMinutes) * time.Minute,
		MaxCacheSize:    maxCacheSize,
	}
}
