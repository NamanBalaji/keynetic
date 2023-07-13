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
	shardId := c.Param("shardId")
	shardIdInt, err := strconv.Atoi(shardId)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)

		return
	}

	if utils.Shard.ShardID != shardIdInt {
		inserted := false
		index := 0

		var down []string
		var keyCountResponse types.GetShardKeyCountResponse

		for !inserted && index < len(utils.Shard.Shards[shardIdInt]) {
			node := utils.Shard.Shards[shardIdInt][index]
			if node != utils.View.SocketAddr {
				res, err := requests.GetShardKeyCount(node, shardIdInt)
				if err == nil {
					jsonData, _ := io.ReadAll(res.Body)
					json.Unmarshal(jsonData, &keyCountResponse)
					inserted = true
				} else {
					down = append(down, node)
				}

				index++
			}
		}

		for _, d := range down {
			utils.View.RemoveFromView(d)
		}
		for _, replica := range utils.View.Views {
			requests.BroadcastDeleteView(replica, down...)
		}

		if inserted {
			c.JSON(http.StatusOK, keyCountResponse)
			return
		}

		c.JSON(http.StatusInternalServerError, nil)
		return
	}

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
			err := requests.BroadcastPutShard(replica, body.SocketAddress, utils.View.SocketAddr, shardIdInt)
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

func ReshardHandler(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("invalid request body [ERROR]: %s", err)
		return
	}

	var body types.ReshardRequest
	err = json.Unmarshal(jsonData, &body)
	if err != nil {
		log.Printf("invalid body format [ERROR]: %s", err)
		return
	}

	var storeRes types.GetStoreResponse

	if len(utils.View.Views)/body.ShardCount >= 2 {
		// get all key value
		for shardId, shardList := range utils.Shard.Shards {
			if !utils.IsReplicaInShard(utils.View.SocketAddr, shardId, utils.Shard.Shards) {
				for _, shard := range shardList {
					res, err := requests.GetKeyValueStore(shard)
					if err == nil {
						jsonData, _ := io.ReadAll(res.Body)
						json.Unmarshal(jsonData, &storeRes)
					}
					for k, v := range storeRes.Store {
						utils.Store.Put(k, v)
					}
				}
			}
		}

		// empty shardlist
		for k := range utils.Shard.Shards {
			delete(utils.Shard.Shards, k)
		}

		// empty vector clock
		for k := range utils.Vc {
			utils.Vc[k] = 0
		}

		utils.Shard.ShardCount = body.ShardCount

		nodesInShard := len(utils.View.Views) / utils.Shard.ShardCount
		nodesSoFar := 0
		shardIdx := 1

		for _, view := range utils.View.Views {
			if shardIdx <= utils.Shard.ShardCount {
				if view == utils.View.SocketAddr {
					utils.Shard.ShardID = shardIdx
				}

				if nodesSoFar < nodesInShard {
					utils.Shard.Shards[shardIdx] = append(utils.Shard.Shards[shardIdx], view)
					nodesSoFar++
				} else {
					shardIdx++
					if shardIdx <= utils.Shard.ShardCount {
						nodesSoFar = 0
						utils.Shard.Shards[shardIdx] = append(utils.Shard.Shards[shardIdx], view)
						nodesSoFar++
						if view == utils.View.SocketAddr {
							utils.Shard.ShardID = shardIdx
						}
					}
				}
			}
		}

		if (len(utils.View.Views) % utils.Shard.ShardCount) == 1 {
			utils.Shard.Shards[shardIdx-1] = append(utils.Shard.Shards[shardIdx-1], utils.View.Views[len(utils.View.Views)-1])
		}
	} else {
		c.JSON(http.StatusBadRequest, "noy enough nodes to have redundancy in shards")
		return
	}

	for shard := range utils.Shard.Shards {
		tempKvStore := make(map[string]string)
		for k, v := range utils.Store.Database {
			if utils.Shard.HashShardIndex(k) == shard {
				tempKvStore[k] = v
			}
		}

		for _, updatedShard := range utils.Shard.Shards[shard] {
			if updatedShard != utils.View.SocketAddr {
				requests.BroadcstReshardShardPut(updatedShard, utils.Shard.Shards)
				requests.BroadcstReshardStorePut(updatedShard, tempKvStore)
			}
		}
	}

	tempKvStore := make(map[string]string)
	for k, v := range utils.Store.Database {
		if utils.Shard.HashShardIndex(k) == utils.Shard.ShardID {
			tempKvStore[k] = v
		}
	}

	utils.SetStore(tempKvStore)

	resp := types.ReshardResponse{
		Message: "Resharding done successfully",
	}

	c.JSON(http.StatusOK, resp)
}
