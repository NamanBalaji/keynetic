package handlers

import (
	"net/http"

	"github.com/NamanBalaji/keynetic/utils"
	"github.com/gin-gonic/gin"
)

type GetStoreResponse struct {
	Store map[string]string `json:"store"`
}

type GetVectorClockResponse struct {
	VectorClock map[string]int `json:"vectorClock"`
}

// Handler for GET: /store
func GetStoreHandler(c *gin.Context) {
	resp := GetStoreResponse{
		Store: utils.Store.Database,
	}
	c.JSON(http.StatusOK, resp)
}

// Handler for GET: /vector-clock
func GetVectorClock(c *gin.Context) {
	resp := GetVectorClockResponse{
		VectorClock: utils.Vc,
	}
	c.JSON(http.StatusOK, resp)
}
