package views_handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/NamanBalaji/keynetic/utils"
	"github.com/gin-gonic/gin"
)

type getViewSucces struct {
	Message string `json:"message,omitempty"`
	View    string `json:"view,omitempty"`
}

// type getViewFail struct {
// 	Message string `json:"message,omitempty"`
// 	Error   string `json:"error,omitempty"`
// }

func GetViewHandler(c *gin.Context) {
	// send heartbeat to each view
	ctx, can := context.WithTimeout(context.Background(), 1*time.Second)
	defer can()
	var downInstances []string
	for _, v := range utils.Views {
		if v != utils.SocketAddr {
			url := fmt.Sprintf("http://%s/ping", v)
			req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			_, err := http.DefaultClient.Do(req)
			if err != nil {
				downInstances = append(downInstances, v)
				utils.RemoveFromView(v)
			}
		}

	}

	// broadcast down instances for deletion
	for _, v := range utils.Views {
		if v != utils.SocketAddr {
			for _, d := range downInstances {
				url := fmt.Sprintf("http://%s/broadcast-delete/%s", v, d)
				req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
				_, err := http.DefaultClient.Do(req)
				if err != nil {
					log.Printf("instance %s down", v)
				}
			}
		}
	}

	resp := getViewSucces{
		Message: "View retrieved successfully",
		View:    fmt.Sprint(strings.Join(utils.Views[:], ",")),
	}
	c.JSON(http.StatusOK, resp)
}
