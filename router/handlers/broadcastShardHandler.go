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
					jsonData, err := io.ReadAll(res.Body)
					if err != nil {
						log.Printf("invalid request body [ERROR]: %s", err)
						c.JSON(http.StatusInternalServerError, err)
						return
					}
					err = json.Unmarshal(jsonData, &shardRes)
					if err != nil {
						log.Printf("invalid body format [ERROR]: %s", err)
						c.JSON(http.StatusInternalServerError, err)
						return
					}

					for k, v := range shardRes.Shard {
						utils.Shard.Shards[k] = v
					}

					break
				}
			}
		}

		var storeRes types.GetStoreResponse
		var vectorClockRes types.GetVectorClockResponse

		for _, replica := range utils.View.Views {
			if replica != utils.View.SocketAddr && utils.IsReplicaInShard(replica, shardIdInt, utils.Shard.Shards) && replica != body.ClientAddress {

				res, err := requests.GetKeyValueStore(replica)
				if err == nil {
					jsonData, err := io.ReadAll(res.Body)
					if err != nil {
						log.Printf("invalid request body [ERROR]: %s", err)
						c.JSON(http.StatusInternalServerError, err)
						return
					}
					err = json.Unmarshal(jsonData, &storeRes)
					if err != nil {
						log.Printf("invalid body format [ERROR]: %s", err)
						c.JSON(http.StatusInternalServerError, err)
						return
					}
					break
				}

			}
		}

		for _, replica := range utils.View.Views {
			utils.Vc[replica] = 0
			if replica != utils.View.SocketAddr && utils.IsReplicaInShard(replica, shardIdInt, utils.Shard.Shards) && replica != body.ClientAddress {

				res, err := requests.GetVectorClock(replica)
				if err == nil {
					jsonData, err := io.ReadAll(res.Body)
					if err != nil {
						log.Printf("invalid request body [ERROR]: %s", err)
						c.JSON(http.StatusInternalServerError, err)
						return
					}
					err = json.Unmarshal(jsonData, &vectorClockRes)
					if err != nil {
						log.Printf("invalid body format [ERROR]: %s", err)
						c.JSON(http.StatusInternalServerError, err)
						return
					}
					break
				}

			}
		}

		utils.SetStore(storeRes.Store)
		utils.SetVectorClock(vectorClockRes.VectorClock)

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
