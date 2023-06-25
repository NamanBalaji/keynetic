package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/NamanBalaji/keynetic/requests"
	"github.com/NamanBalaji/keynetic/router"
	"github.com/NamanBalaji/keynetic/types"
	"github.com/NamanBalaji/keynetic/utils"
)

const port = "8085"

func main() {

	views := strings.Split(os.Getenv("VIEW"), ",")
	socketAddr := os.Getenv("SOCKET_ADDRESS")
	shardCount, _ := strconv.Atoi(os.Getenv("SHARD_COUNT"))

	utils.InitViews(views, socketAddr)
	utils.InitShard(shardCount, socketAddr)
	utils.InitStore()
	utils.InitVectorClock(views)

	routesInit := router.InitMainRouter()

	endpoint := fmt.Sprintf(":%s", port)
	server := &http.Server{
		Addr:    endpoint,
		Handler: routesInit,
	}

	// send put request to all replicas
	for _, replica := range utils.View.Views {
		if replica != utils.View.SocketAddr {
			requests.BroadcastPutView(replica, utils.View.SocketAddr)
		}
	}

	var storeRes types.GetStoreResponse
	var vectorClockRes types.GetVectorClockResponse

	// ask for a replicas key value store
	if utils.Shard.ShardID != -1 {
		for _, replica := range utils.View.Views {
			if replica != utils.View.SocketAddr && isReplicaInShard(replica, utils.Shard.ShardID, utils.Shard.Shards) {
				res, err := requests.GetKeyValueStore(replica)
				if err == nil {
					jsonData, err := io.ReadAll(res.Body)
					if err != nil {
						log.Printf("invalid request body [ERROR]: %s", err)
						return
					}
					err = json.Unmarshal(jsonData, &storeRes)
					if err != nil {
						log.Printf("invalid body format [ERROR]: %s", err)
						return
					}
					break
				}
			}
		}

		// ask for a replicas vector clock
		for _, replica := range utils.View.Views {
			if replica != utils.View.SocketAddr && isReplicaInShard(replica, utils.Shard.ShardID, utils.Shard.Shards) {
				res, err := requests.GetVectorClock(replica)
				if err == nil {
					jsonData, err := io.ReadAll(res.Body)
					if err != nil {
						log.Printf("invalid request body [ERROR]: %s", err)
						return
					}
					err = json.Unmarshal(jsonData, &vectorClockRes)
					if err != nil {
						log.Printf("invalid body format [ERROR]: %s", err)
						return
					}
					break
				}
			}
		}
	} else {
		log.Printf("shardID uninitialized, will update kvStore and vector clocks later")
	}

	utils.SetStore(storeRes.Store)
	utils.SetVectorClock(vectorClockRes.VectorClock)

	log.Printf("HTTP server started on port %s", endpoint)
	_ = server.ListenAndServe()
}

func isReplicaInShard(replica string, shardId int, shardMap map[int][]string) bool {
	shards, ok := shardMap[shardId]
	if !ok {
		return false
	}

	for _, r := range shards {
		if r == replica {
			return true
		}
	}

	return false
}
