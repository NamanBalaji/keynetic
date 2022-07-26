package router

import (
	"net/http"

	"github.com/NamanBalaji/keynetic/router/handler"
	"github.com/gin-gonic/gin"
)

func InitMainRouter() *gin.Engine {

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "I'm alive"})
	})

	kvApi := r.Group("/key-value-store")
	{
		kvApi.GET("/:key", handler.GetHandler)
		kvApi.DELETE("/:key", handler.DeleteHandler)
		kvApi.PUT("/:key", handler.PutHandler)
	}
	return r
}

func InitForwardRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "I'm forwarding server"})
	})

	kvApi := r.Group("/key-value-store")
	{
		kvApi.GET("/:key", handler.ForwardHandler)
		kvApi.DELETE("/:key", handler.ForwardHandler)
		kvApi.PUT("/:key", handler.ForwardHandler)
	}
	return r
}
