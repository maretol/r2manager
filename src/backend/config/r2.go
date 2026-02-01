package config

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
)

type R2Config struct {
	AccountID       string
	AccessKeyID     string
	SecretAccessKey string
}

func LoadR2ConfigFromEnv() (*R2Config, error) {
	accountID := os.Getenv("R2_ACCOUNT_ID")
	accessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")

	if accountID == "" || accessKeyID == "" || secretAccessKey == "" {
		return nil, errors.New("R2_ACCOUNT_ID, R2_ACCESS_KEY_ID, R2_SECRET_ACCESS_KEY must be set")
	}

	return &R2Config{
		AccountID:       accountID,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	}, nil
}

func NewS3Client(ctx context.Context, r2cfg *R2Config) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("auto"),
		config.WithBaseEndpoint(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", r2cfg.AccountID)),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(r2cfg.AccessKeyID, r2cfg.SecretAccessKey, "")),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load AWS config")
	}

	return s3.NewFromConfig(cfg), nil
}
