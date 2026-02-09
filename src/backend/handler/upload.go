package handler

import (
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	serviceif "r2manager/service/interface"
)

type UploadHandler struct {
	service       serviceif.UploadService
	maxUploadSize int64
}

func NewUploadHandler(service serviceif.UploadService, maxUploadSize int64) *UploadHandler {
	return &UploadHandler{service: service, maxUploadSize: maxUploadSize}
}

func (h *UploadHandler) UploadObject(ctx *gin.Context) {
	bucketName := ctx.Param("bucketName")
	if bucketName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bucketName is required"})
		return
	}

	key := ctx.Param("key")
	key = strings.TrimPrefix(key, "/")
	if key == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "key is required"})
		return
	}

	overwrite := ctx.Query("overwrite") == "true"

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	if header.Size > h.maxUploadSize {
		ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"error":    "file too large",
			"max_size": h.maxUploadSize,
		})
		return
	}

	contentType := detectContentType(header.Filename, header.Header.Get("Content-Type"))

	result, err := h.service.UploadObject(ctx.Request.Context(), bucketName, key, contentType, file, header.Size, overwrite)
	if err != nil {
		if errors.Is(err, serviceif.ErrObjectAlreadyExists) {
			ctx.JSON(http.StatusConflict, gin.H{"error": "object already exists", "code": "CONFLICT"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (h *UploadHandler) CreateDirectory(ctx *gin.Context) {
	bucketName := ctx.Param("bucketName")
	if bucketName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bucketName is required"})
		return
	}

	var req struct {
		Path string `json:"path" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
		return
	}

	result, err := h.service.CreateDirectory(ctx.Request.Context(), bucketName, req.Path)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func detectContentType(filename, provided string) string {
	if provided != "" && provided != "application/octet-stream" {
		return provided
	}

	ext := filepath.Ext(filename)
	if ext != "" {
		if ct := mime.TypeByExtension(ext); ct != "" {
			return ct
		}
	}

	return "application/octet-stream"
}
