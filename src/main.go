package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/NamanBalaji/keynetic/router"
	"github.com/NamanBalaji/keynetic/utils"
)

const port = "8085"

func main() {
	utils.InitEnvVars()

	routesInit := router.InitMainRouter()

	utils.InitStore()

	endpoint := fmt.Sprintf(":%s", port)
	server := &http.Server{
		Addr:    endpoint,
		Handler: routesInit,
	}

	// send put request to all replicas
	ctx, can := context.WithTimeout(context.Background(), 1*time.Second)

	for _, replica := range utils.Views {
		if replica != utils.SocketAddr {
			url := fmt.Sprintf("http://%s/broadcast-put/%s", replica, utils.SocketAddr)
			req, _ := http.NewRequestWithContext(ctx, http.MethodPut, url, nil)
			_, _ = http.DefaultClient.Do(req)
		}
	}
	can()

	// ask for a replicas key value store
	ctx, can = context.WithTimeout(context.Background(), 1*time.Second)

	for _, replica := range utils.Views {
		if replica != utils.SocketAddr {
			url := fmt.Sprintf("http://%s/store", replica)
			req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				json.NewDecoder(res.Body).Decode(&utils.Store.Database)
				break
			}
		}
	}
	can()

	log.Printf("HTTP server started on port %s", endpoint)
	_ = server.ListenAndServe()
}
