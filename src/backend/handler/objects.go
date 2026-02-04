package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	serviceif "r2manager/service/interface"
)

type ObjectsHandler struct {
	service serviceif.ObjectService
}

func NewObjectsHandler(service serviceif.ObjectService) *ObjectsHandler {
	return &ObjectsHandler{service: service}
}

func (oh *ObjectsHandler) GetObjects(ctx *gin.Context) {
	bucketName := ctx.Param("bucketName")
	if bucketName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bucketName is required"})
		return
	}

	params := serviceif.ListObjectsParams{
		Prefix:    ctx.Query("prefix"),
		Delimiter: ctx.Query("delimiter"),
	}

	result, err := oh.service.GetObjects(ctx.Request.Context(), bucketName, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}
