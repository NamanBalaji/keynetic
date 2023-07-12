package handlers

import (
	"net/http"

	"github.com/NamanBalaji/keynetic/types"
	"github.com/NamanBalaji/keynetic/utils"
	"github.com/gin-gonic/gin"
)

// Handler for GET: /store
func GetStoreHandler(c *gin.Context) {
	resp := types.GetStoreResponse{
		Store: utils.Store.Database,
	}
	c.JSON(http.StatusOK, resp)
}

// Handler for GET: /vector-clock
func GetVectorClock(c *gin.Context) {
	resp := types.GetVectorClockResponse{
		VectorClock: utils.Vc,
	}
	c.JSON(http.StatusOK, resp)
}

// Handler for GET: /shard
func GetShard(c *gin.Context) {
	resp := types.GetShardResponse{
		Shard: utils.Shard.Shards,
	}
	c.JSON(http.StatusOK, resp)
}
