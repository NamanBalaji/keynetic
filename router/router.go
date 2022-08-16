package router

import (
	"net/http"

	"github.com/NamanBalaji/keynetic/router/handlers"
	"github.com/gin-gonic/gin"
)

func InitMainRouter() *gin.Engine {

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "I'm alive"})
	})

	kvApi := r.Group("/key-value-store")
	{
		kvApi.GET("/:key", handlers.GetKVHandler)
		kvApi.DELETE("/:key", handlers.DeleteKVHandler)
		kvApi.PUT("/:key", handlers.PutKVHandler)
	}
	r.GET("/store", handlers.GetStoreHandler)

	r.PUT("/broadcast-view/:ip", handlers.BroadcastViewPut)
	r.DELETE("/broadcast-view/:ip", handlers.BroadcastViewDelete)

	r.GET("/key-value-store-view", handlers.GetViewHandler)
	r.PUT("/key-value-store-view", handlers.PutViewHandler)
	r.DELETE("/key-value-store-view", handlers.DeleteViewHandler)

	return r
}
