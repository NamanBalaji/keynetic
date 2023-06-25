package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/NamanBalaji/keynetic/types"
	"github.com/NamanBalaji/keynetic/utils"
	"github.com/gin-gonic/gin"
)

// Handler for PUT: /broadcast-shard/:shardId
func BroadcastShardPut(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("invalid request body [ERROR]: %s", err)
		return
	}

	var body types.ShardAddMemberRequest
	err = json.Unmarshal(jsonData, &body)
	if err != nil {
		log.Printf("invalid body format [ERROR]: %s", err)
		return
	}

	shardId := c.Param("shardId")

	shardIdInt, err := strconv.Atoi(shardId)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)

		return
	}

	addrs := utils.Shard.Shards[shardIdInt]
	exist := false
	for _, addr := range addrs {
		if addr == body.SocketAddress {
			exist = true
		}
	}

	if !exist {
		utils.Shard.Shards[shardIdInt] = append(utils.Shard.Shards[shardIdInt], body.SocketAddress)
	}

	c.JSON(http.StatusOK, nil)
}
