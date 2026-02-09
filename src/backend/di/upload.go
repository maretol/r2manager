package di

import (
	appconfig "r2manager/config"
	"r2manager/handler"
	"r2manager/repository"
	service "r2manager/service/model"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func CreateUploadHandler(s3Client *s3.Client, listCache *repository.ListCacheRepository, uploadCfg *appconfig.UploadConfig) *handler.UploadHandler {
	uploadRepo := repository.NewUploadRepository(s3Client)
	uploadService := service.NewUploadService(uploadRepo, listCache)
	uploadHandler := handler.NewUploadHandler(uploadService, uploadCfg.MaxUploadSize)

	return uploadHandler
}
