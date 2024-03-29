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

	r.PUT("/broadcast-kv/:key", handlers.BroadcastKeyPut)
	r.DELETE("/broadcast-kv/:key", handlers.BroadcastKeyDelete)

	r.GET("/store", handlers.GetStoreHandler)
	r.GET("/vector-clock", handlers.GetVectorClock)
	r.GET("/shard", handlers.GetShard)

	r.PUT("/broadcast-view/:ip", handlers.BroadcastViewPut)
	r.DELETE("/broadcast-view/:ip", handlers.BroadcastViewDelete)

	r.GET("/key-value-store-view", handlers.GetViewHandler)
	r.PUT("/key-value-store-view", handlers.PutViewHandler)
	r.DELETE("/key-value-store-view", handlers.DeleteViewHandler)

	shardApi := r.Group("/key-value-store-shard")
	{
		shardApi.GET("/shard-ids", handlers.GetShardIds)
		shardApi.GET("/node-shard-id", handlers.GetNodeShardId)
		shardApi.GET("/shard-id-members/:shardId", handlers.GetShardMembers)
		shardApi.GET("/shard-id-key-count/:shardId", handlers.GetShardKeyCount)
		shardApi.PUT("/add-member/:shardId", handlers.ShardAddMember)
		shardApi.PUT("/reshard", handlers.ReshardHandler)
	}

	r.PUT("/broadcast-shard/:shardId", handlers.BroadcastShardPut)

	broadcastReshardApi := r.Group("/broadcast-reshard")
	{
		broadcastReshardApi.PUT("/shard", handlers.ReshardShardPutHandler)
		broadcastReshardApi.PUT("/store", handlers.ReshardStorePutHandler)
	}

	return r
}
