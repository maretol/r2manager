package handler

import (
	"errors"
	"mime"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"

	"r2manager/domain"
	"r2manager/progress"
	serviceif "r2manager/service/interface"
)

type UploadHandler struct {
	service       serviceif.UploadService
	maxUploadSize int64
	progressStore *progress.UploadProgressStore
}

func NewUploadHandler(service serviceif.UploadService, maxUploadSize int64, progressStore *progress.UploadProgressStore) *UploadHandler {
	return &UploadHandler{service: service, maxUploadSize: maxUploadSize, progressStore: progressStore}
}

func (h *UploadHandler) UploadObject(ctx *gin.Context) {
	bucketName := ctx.Param("bucketName")
	if bucketName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bucketName is required"})
		return
	}

	key := ctx.Param("key")
	if key == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "key is required"})
		return
	}

	overwrite := ctx.Query("overwrite") == "true"

	// Upload ID の取得（ヘッダー優先、クエリパラメータにフォールバック）
	uploadID := ctx.GetHeader("X-Upload-ID")
	if uploadID == "" {
		uploadID = ctx.Query("upload_id")
	}

	// マルチパートのオーバーヘッド分を加算してボディサイズを制限する
	// FormFile によるパース前に制限をかけることで、巨大リクエストによるリソース消費を防ぐ
	const multipartOverhead = 4096
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, h.maxUploadSize+multipartOverhead)

	// Phase 1: リクエストボディの受信進捗を追跡
	if uploadID != "" {
		h.progressStore.Register(uploadID)
		totalBytes := max(ctx.Request.ContentLength, 0)
		progressReader, err := progress.NewProgressReadCloser(
			ctx.Request.Body,
			func(bytesProcessed int64) {
				h.progressStore.Publish(uploadID, domain.UploadEvent{
					EventType: domain.EventProgress,
					Data: domain.UploadProgress{
						UploadID:       uploadID,
						Phase:          domain.PhaseReceiving,
						BytesProcessed: bytesProcessed,
						TotalBytes:     totalBytes,
					},
				})
			},
			100*time.Millisecond,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.Request.Body = progressReader
	}

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			if uploadID != "" {
				h.progressStore.Publish(uploadID, domain.UploadEvent{
					EventType: domain.EventError,
					Data:      domain.UploadError{UploadID: uploadID, Error: "file too large"},
				})
			}
			ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error":    "file too large",
				"max_size": h.maxUploadSize,
			})
			return
		}
		if uploadID != "" {
			h.progressStore.Publish(uploadID, domain.UploadEvent{
				EventType: domain.EventError,
				Data:      domain.UploadError{UploadID: uploadID, Error: "file is required"},
			})
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	if header.Size > h.maxUploadSize {
		if uploadID != "" {
			h.progressStore.Publish(uploadID, domain.UploadEvent{
				EventType: domain.EventError,
				Data:      domain.UploadError{UploadID: uploadID, Error: "file too large"},
			})
		}
		ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"error":    "file too large",
			"max_size": h.maxUploadSize,
		})
		return
	}

	contentType := detectContentType(header.Filename, header.Header.Get("Content-Type"))

	// Phase 2 のコールバックを準備
	var uploadingCallback serviceif.ProgressCallback
	if uploadID != "" {
		uploadingCallback = func(bytesProcessed int64) {
			h.progressStore.Publish(uploadID, domain.UploadEvent{
				EventType: domain.EventProgress,
				Data: domain.UploadProgress{
					UploadID:       uploadID,
					Phase:          domain.PhaseUploading,
					BytesProcessed: bytesProcessed,
					TotalBytes:     header.Size,
				},
			})
		}
	}

	result, err := h.service.UploadObject(ctx.Request.Context(), bucketName, key, contentType, file, header.Size, overwrite, uploadingCallback)
	if err != nil {
		if uploadID != "" {
			h.progressStore.Publish(uploadID, domain.UploadEvent{
				EventType: domain.EventError,
				Data:      domain.UploadError{UploadID: uploadID, Error: err.Error()},
			})
		}
		if errors.Is(err, serviceif.ErrObjectAlreadyExists) {
			ctx.JSON(http.StatusConflict, gin.H{"error": "object already exists", "code": "CONFLICT"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if uploadID != "" {
		h.progressStore.Publish(uploadID, domain.UploadEvent{
			EventType: domain.EventComplete,
			Data:      domain.UploadComplete{UploadID: uploadID, Result: result},
		})
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
