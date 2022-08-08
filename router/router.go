package router

import (
	"net/http"

	"github.com/NamanBalaji/keynetic/router/broadcast_handler"
	"github.com/NamanBalaji/keynetic/router/kv_handler"
	"github.com/NamanBalaji/keynetic/router/views_handler"
	"github.com/gin-gonic/gin"
)

func InitMainRouter() *gin.Engine {

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "I'm alive"})
	})

	kvApi := r.Group("/key-value-store")
	{
		kvApi.GET("/:key", kv_handler.GetHandler)
		kvApi.DELETE("/:key", kv_handler.DeleteHandler)
		kvApi.PUT("/:key", kv_handler.PutHandler)
	}
	r.GET("/store", kv_handler.GetStoreHandler)

	r.PUT("/broadcast-put/:ip", broadcast_handler.BroadcastPut)
	r.DELETE("/broadcast-delete/:ip", broadcast_handler.BroadcastDelete)

	r.GET("/key-value-store-view", views_handler.GetViewHandler)
	r.PUT("/key-value-store-view", views_handler.PutViewsHandler)
	r.DELETE("/key-value-store-view", views_handler.DeleteViewHandler)

	return r
}
