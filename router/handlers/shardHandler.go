package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/NamanBalaji/keynetic/requests"
	"github.com/NamanBalaji/keynetic/types"
	"github.com/NamanBalaji/keynetic/utils"
	"github.com/gin-gonic/gin"
)

// Handler for GET: /key-value-store-shard/shard-ids
func GetShardIds(c *gin.Context) {
	var shardIds []int

	for id := range utils.Shard.Shards {
		shardIds = append(shardIds, id)
	}

	resp := types.GetShardIdsResponse{
		ShardIds: shardIds,
		Message:  "Shard IDs retrieved successfully",
	}

	c.JSON(http.StatusOK, resp)
}

// Handler for GET: /key-value-store-shard/node-shard-id
func GetNodeShardId(c *gin.Context) {
	resp := types.GetNodeShardIdResponse{
		Message: "Shard ID of the node retrieved successfully",
		ShardID: utils.Shard.ShardID,
	}

	c.JSON(http.StatusOK, resp)
}

// Handler for GET: /key-value-store-shard/shard-id-members/
func GetShardMembers(c *gin.Context) {
	shardId := c.Param("shardId")

	shardIdInt, err := strconv.Atoi(shardId)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)

		return
	}

	resp := types.GetShardIdMembersResponse{
		Message:        "Members of shard ID retrieved successfully",
		ShardIdMembers: utils.Shard.Shards[shardIdInt],
	}

	c.JSON(http.StatusOK, resp)
}

// Handler for GET: /key-value-store-shard/shard-id-key-count/
func GetShardKeyCount(c *gin.Context) {
	resp := types.GetShardIDKeyCountResponse{
		Message:  "Key count of shard ID retrieved successfully",
		KeyCount: len(utils.Store.Database),
	}

	c.JSON(http.StatusOK, resp)
}

// Handler for PUT: /key-value-store-shard/add-member/:shardId
func ShardAddMember(c *gin.Context) {
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

	var down []string
	for _, replica := range utils.View.Views {
		if replica != utils.View.SocketAddr {
			err := requests.BroadcastPutShard(replica, body.SocketAddress, shardIdInt)
			if err != nil {
				down = append(down, replica)
			}
		}
	}

	for _, d := range down {
		utils.View.RemoveFromView(d)
	}
	for _, replica := range utils.View.Views {
		requests.BroadcastDeleteView(replica, down...)
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
