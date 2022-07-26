package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/NamanBalaji/keynetic/router"
	"github.com/NamanBalaji/keynetic/utils"
	"github.com/gin-gonic/gin"
)

const port = "8085"

func main() {
	utils.InitEnvVars()
	var routesInit *gin.Engine
	if utils.Env.FwdAddr != "" {
		routesInit = router.InitForwardRouter()
	} else {
		routesInit = router.InitMainRouter()
	}

	utils.InitStore()

	endpoint := fmt.Sprintf(":%s", port)
	server := &http.Server{
		Addr:    endpoint,
		Handler: routesInit,
	}

	log.Printf("HTTP server started on port %s", endpoint)
	_ = server.ListenAndServe()
}
