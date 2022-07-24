package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/NamanBalaji/keynetic/router"
	"github.com/NamanBalaji/keynetic/router/handler/utils"
)

const port = "8085"

func main() {
	routesInit := router.InitRouter()
	utils.InitStore()

	endpoint := fmt.Sprintf(":%s", port)
	server := &http.Server{
		Addr:    endpoint,
		Handler: routesInit,
	}

	log.Printf("HTTP server started on port %s", endpoint)
	_ = server.ListenAndServe()
}
