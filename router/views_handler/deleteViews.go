package views_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/NamanBalaji/keynetic/utils"
	"github.com/gin-gonic/gin"
)

type deleteViewSucces struct {
	Message string `json:"message,omitempty"`
}

type deleteViewFail struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func DeleteViewHandler(c *gin.Context) {
	ctx, can := context.WithTimeout(context.Background(), 1*time.Second)
	defer can()

	var addr string
	body, _ := ioutil.ReadAll(c.Request.Body)

	strBody := string(body[:])
	json.NewDecoder(strings.NewReader(strBody)).Decode(&addr)
	defer c.Request.Body.Close()

	_, ok := utils.Contains(addr)
	// check if view present
	if !ok {
		resp := deleteViewFail{
			Message: "Socket address does not exist in the view",
			Error:   "Error in DELETE",
		}
		c.JSON(http.StatusNotFound, resp)
		return
	}

	utils.RemoveFromView(addr)

	for _, v := range utils.Views {
		if v != utils.SocketAddr {
			url := fmt.Sprintf("http://%s/broadcast-delete/%s", v, addr)
			req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
			_, _ = http.DefaultClient.Do(req)
		}
	}

	resp := deleteViewSucces{
		Message: "Replica deleted successfully from the view",
	}
	c.JSON(http.StatusOK, resp)

}
