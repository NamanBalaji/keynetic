package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/NamanBalaji/keynetic/router"
)

const port = "8085"

func main() {
	routesInit := router.InitRouter()

	endpoint := fmt.Sprintf(":%s", port)
	server := &http.Server{
		Addr:    endpoint,
		Handler: routesInit,
	}

	log.Printf("HTTP server started on port %s", endpoint)
	_ = server.ListenAndServe()
}
