package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"r2manager/domain"
	serviceif "r2manager/service/interface"
)

type SettingsHandler struct {
	service serviceif.SettingsService
}

func NewSettingsHandler(service serviceif.SettingsService) *SettingsHandler {
	return &SettingsHandler{service: service}
}

func (h *SettingsHandler) GetAllBucketSettings(ctx *gin.Context) {
	settings, err := h.service.GetAllBucketSettings(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if settings == nil {
		settings = []domain.BucketSettings{}
	}
	ctx.JSON(http.StatusOK, gin.H{"settings": settings})
}

func (h *SettingsHandler) GetBucketSettings(ctx *gin.Context) {
	bucketName := ctx.Param("bucketName")
	if bucketName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bucketName is required"})
		return
	}

	settings, err := h.service.GetBucketSettings(ctx.Request.Context(), bucketName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if settings == nil {
		ctx.JSON(http.StatusOK, domain.BucketSettings{
			BucketName: bucketName,
			PublicUrl:  "",
		})
		return
	}

	ctx.JSON(http.StatusOK, settings)
}

type updateBucketSettingsRequest struct {
	PublicUrl string `json:"public_url"`
}

func (h *SettingsHandler) UpdateBucketSettings(ctx *gin.Context) {
	bucketName := ctx.Param("bucketName")
	if bucketName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bucketName is required"})
		return
	}

	var req updateBucketSettingsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.service.UpdateBucketPublicUrl(ctx.Request.Context(), bucketName, req.PublicUrl); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, domain.BucketSettings{
		BucketName: bucketName,
		PublicUrl:  req.PublicUrl,
	})
}
