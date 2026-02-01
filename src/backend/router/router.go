package router

import (
	"github.com/gin-gonic/gin"

	"r2manager/handler"
)

func NewRouter(bucketsHandler *handler.BucketsHandler, objectsHandler *handler.ObjectsHandler, contentHandler *handler.ContentHandler) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.GET("/buckets", bucketsHandler.GetBuckets)
		api.GET("/buckets/:bucketName/objects", objectsHandler.GetObjects)
		api.GET("/buckets/:bucketName/content/*key", contentHandler.GetContent)
	}

	return r
}
