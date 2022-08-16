package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/NamanBalaji/keynetic/requests"
	"github.com/NamanBalaji/keynetic/router"
	"github.com/NamanBalaji/keynetic/utils"
)

const port = "8085"

func main() {

	views := strings.Split(os.Getenv("VIEWS"), ",")
	socketAddr := os.Getenv("SOCKET_ADDRESS")

	utils.InitViews(views, socketAddr)
	utils.InitStore()

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

	// ask for a replicas key value store
	for _, replica := range utils.View.Views {
		if replica != utils.View.SocketAddr {
			res, err := requests.GetKeyValueStore(replica)
			if err != nil {
				json.NewDecoder(res.Body).Decode(&utils.Store.Database)
				break
			}
		}
	}

	log.Printf("HTTP server started on port %s", endpoint)
	_ = server.ListenAndServe()
}
