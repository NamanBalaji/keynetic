package handlers

import (
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/NamanBalaji/keynetic/requests"
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

	var body types.BroadcastShardPutRequest
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

	if body.SocketAddress == utils.View.SocketAddr {
		utils.Shard.ShardID = shardIdInt

		for _, replica := range utils.View.Views {
			if replica != utils.View.SocketAddr && replica != body.ClientAddress {
				var shardRes types.GetShardResponse
				res, err := requests.GetShard(replica)
				if err == nil {
					jsonData, _ := io.ReadAll(res.Body)
					json.Unmarshal(jsonData, &shardRes)
					for k, v := range shardRes.Shard {
						utils.Shard.Shards[k] = v
					}
					break
				}
			}
		}

		for _, replica := range utils.View.Views {
			if replica != utils.View.SocketAddr && utils.IsReplicaInShard(replica, shardIdInt, utils.Shard.Shards) && replica != body.ClientAddress {
				var storeRes types.GetStoreResponse
				res, err := requests.GetKeyValueStore(replica)
				if err == nil {
					jsonData, _ := io.ReadAll(res.Body)
					json.Unmarshal(jsonData, &storeRes)
					for k, v := range storeRes.Store {
						utils.Store.Put(k, v)
					}
					break
				}

			}
		}

		for _, replica := range utils.View.Views {
			utils.Vc[replica] = 0
			if replica != utils.View.SocketAddr && utils.IsReplicaInShard(replica, shardIdInt, utils.Shard.Shards) && replica != body.ClientAddress {
				var vectorClockRes types.GetVectorClockResponse
				res, err := requests.GetVectorClock(replica)
				if err == nil {
					jsonData, _ := io.ReadAll(res.Body)
					json.Unmarshal(jsonData, &vectorClockRes)
					for k, v := range vectorClockRes.VectorClock {
						utils.Vc[replica] = int(math.Max(float64(utils.Vc[k]), float64(v)))
					}
					break
				}

			}
		}
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
