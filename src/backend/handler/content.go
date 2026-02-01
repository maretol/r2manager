package handler

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	serviceif "r2manager/service/interface"
)

type ContentHandler struct {
	service serviceif.ContentService
}

func NewContentHandler(service serviceif.ContentService) *ContentHandler {
	return &ContentHandler{service: service}
}

func (ch *ContentHandler) GetContent(ctx *gin.Context) {
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

	content, err := ch.service.GetContent(ctx.Request.Context(), bucketName, key)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer content.Body.Close()

	if content.CacheHit {
		ctx.Header("X-Cache", "HIT")
	} else {
		ctx.Header("X-Cache", "MISS")
	}

	if content.ETag != "" {
		ctx.Header("ETag", content.ETag)
	}

	ctx.Header("Content-Type", content.ContentType)
	if content.Size > 0 {
		ctx.Header("Content-Length", strconv.FormatInt(content.Size, 10))
	}

	ctx.Status(http.StatusOK)
	io.Copy(ctx.Writer, content.Body)
}
