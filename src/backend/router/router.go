package router

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"r2manager/handler"
)

func NewRouter(bucketsHandler *handler.BucketsHandler, objectsHandler *handler.ObjectsHandler, contentHandler *handler.ContentHandler, cacheHandler *handler.CacheHandler, settingsHandler *handler.SettingsHandler) *gin.Engine {
	r := gin.Default()

	trustedIPList := getTrustedIPList()
	if len(trustedIPList) > 0 {
		r.SetTrustedProxies(trustedIPList)
	} else {
		// CDN等の環境も設定できるようにしたいけどとりあえずこれで
	}

	api := r.Group("/api/v1")
	{
		api.GET("/buckets", bucketsHandler.GetBuckets)
		api.GET("/buckets/:bucketName/objects", objectsHandler.GetObjects)
		api.GET("/buckets/:bucketName/content/*key", contentHandler.GetContent)

		api.DELETE("/cache/content", cacheHandler.ClearContentCache)
		api.DELETE("/cache/api", cacheHandler.ClearAPICache)

		api.GET("/settings/buckets", settingsHandler.GetAllBucketSettings)
		api.GET("/settings/buckets/:bucketName", settingsHandler.GetBucketSettings)
		api.PUT("/settings/buckets/:bucketName", settingsHandler.UpdateBucketSettings)
	}

	return r
}

func getTrustedIPList() []string {
	env := os.Getenv("env")
	if env == "dev" {
		return []string{"192.168.0.0/24", "127.0.0.1"}
	}
	rawIpList := os.Getenv("IP_LIST")
	if rawIpList != "" {
		ipList := strings.Split(rawIpList, ",")
		results := []string{}
		for _, ip := range ipList {
			trimmed := strings.TrimSpace(ip)
			if trimmed == "" {
				continue
			}
			results = append(results, trimmed)
		}
		return results
	}

	return []string{}
}
