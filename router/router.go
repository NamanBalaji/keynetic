package router

import (
	"net/http"

	"github.com/NamanBalaji/keynetic/router/handler"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	kvApi := r.Group("/key-value-store")
	{
		kvApi.GET("/:key", handler.GetHandler)
		kvApi.DELETE("/:key", handler.DeleteHandler)
		kvApi.PUT("/:key", handler.PutHandler)
	}
	return r
}
