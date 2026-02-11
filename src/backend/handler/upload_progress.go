package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"r2manager/progress"
)

type UploadProgressHandler struct {
	store *progress.UploadProgressStore
}

func NewUploadProgressHandler(store *progress.UploadProgressStore) *UploadProgressHandler {
	return &UploadProgressHandler{store: store}
}

// GetUploadProgress は SSE ストリームでアップロード進捗を配信する。
// GET /api/v1/uploads/:uploadId/progress
func (h *UploadProgressHandler) GetUploadProgress(ctx *gin.Context) {
	uploadID := ctx.Param("uploadId")
	if uploadID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "uploadId is required"})
		return
	}

	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("X-Accel-Buffering", "no")

	eventCh, unsubscribe := h.store.Subscribe(uploadID)
	defer unsubscribe()

	clientGone := ctx.Request.Context().Done()

	ctx.Stream(func(w io.Writer) bool {
		select {
		case <-clientGone:
			return false
		case event, ok := <-eventCh:
			if !ok {
				return false
			}
			data, _ := json.Marshal(event.Data)
			fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event.EventType, string(data))
			ctx.Writer.Flush()

			if event.EventType == "complete" || event.EventType == "error" {
				return false
			}
			return true
		}
	})
}
