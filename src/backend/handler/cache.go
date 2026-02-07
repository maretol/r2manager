package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"r2manager/repository"
)

type CacheHandler struct {
	cacheRepo     *repository.CacheRepository
	listCacheRepo *repository.ListCacheRepository
}

func NewCacheHandler(cacheRepo *repository.CacheRepository, listCacheRepo *repository.ListCacheRepository) *CacheHandler {
	return &CacheHandler{
		cacheRepo:     cacheRepo,
		listCacheRepo: listCacheRepo,
	}
}

func (h *CacheHandler) ClearContentCache(ctx *gin.Context) {
	bucketName := ctx.Query("bucket")
	objectKey := ctx.Query("key")

	var affected int64
	var err error

	switch {
	case bucketName != "" && objectKey != "":
		affected, err = h.cacheRepo.ClearByKey(ctx.Request.Context(), bucketName, objectKey)
	case bucketName != "":
		affected, err = h.cacheRepo.ClearByBucket(ctx.Request.Context(), bucketName)
	default:
		affected, err = h.cacheRepo.ClearAll(ctx.Request.Context())
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Content cache cleared",
		"deleted": affected,
	})
}

func (h *CacheHandler) ClearAPICache(ctx *gin.Context) {
	cacheType := ctx.Query("type")
	bucketName := ctx.Query("bucket")

	var message string

	switch cacheType {
	case "buckets":
		h.listCacheRepo.InvalidateBuckets()
		message = "Buckets cache cleared"
	case "objects":
		if bucketName != "" {
			h.listCacheRepo.InvalidateObjects(bucketName)
			message = "Objects cache cleared for bucket: " + bucketName
		} else {
			h.listCacheRepo.InvalidateAllObjects()
			message = "All objects cache cleared"
		}
	default:
		h.listCacheRepo.InvalidateAll()
		message = "API cache cleared"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}
