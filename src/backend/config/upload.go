package config

import (
	"os"
	"strconv"
)

const defaultMaxUploadSizeMB = 100

type UploadConfig struct {
	MaxUploadSize int64 // bytes
}

func LoadUploadConfigFromEnv() *UploadConfig {
	maxSizeMB := defaultMaxUploadSizeMB
	if v := os.Getenv("UPLOAD_MAX_SIZE_MB"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			maxSizeMB = parsed
		}
	}

	return &UploadConfig{
		MaxUploadSize: int64(maxSizeMB) * 1024 * 1024,
	}
}
