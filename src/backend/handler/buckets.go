package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	serviceif "r2manager/service/interface"
)

type BucketsHandler struct {
	service serviceif.BucketService
}

func NewBucketsHandler(service serviceif.BucketService) *BucketsHandler {
	return &BucketsHandler{service: service}
}

func (bh *BucketsHandler) GetBuckets(ctx *gin.Context) {
	buckets, err := bh.service.GetBuckets(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"buckets": buckets})
}
