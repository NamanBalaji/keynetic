package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/NamanBalaji/keynetic/types"
	"github.com/NamanBalaji/keynetic/utils"
	"github.com/gin-gonic/gin"
)

func ReshardStorePutHandler(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("invalid request body [ERROR]: %s", err)
		return
	}

	var body types.ReshardStoreRequest
	err = json.Unmarshal(jsonData, &body)
	if err != nil {
		log.Printf("invalid body format [ERROR]: %s", err)
		return
	}

	utils.SetStore(body.Store)

	c.JSON(http.StatusOK, nil)
}

func ReshardShardPutHandler(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("invalid request body [ERROR]: %s", err)
		return
	}

	var body types.ReshardShardRequest
	err = json.Unmarshal(jsonData, &body)
	if err != nil {
		log.Printf("invalid body format [ERROR]: %s", err)
		return
	}

	utils.Shard.Clear()

	utils.Shard.Shards = body.Shards
	utils.Shard.ShardCount = len(body.Shards)

	for shardId := range utils.Shard.Shards {
		if utils.IsReplicaInShard(utils.View.SocketAddr, shardId, utils.Shard.Shards) {
			utils.Shard.ShardID = shardId
		}
	}

	for k := range utils.Vc {
		utils.Vc[k] = 0
	}

	c.JSON(http.StatusOK, nil)
}
