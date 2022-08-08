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

type putViewSucces struct {
	Message string `json:"message,omitempty"`
}

type putViewFail struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func PutViewsHandler(c *gin.Context) {
	// get body
	var addr string
	body, _ := ioutil.ReadAll(c.Request.Body)
	strBody := string(body[:])
	json.NewDecoder(strings.NewReader(strBody)).Decode(&addr)
	defer c.Request.Body.Close()

	ctx, can := context.WithTimeout(context.Background(), 1*time.Second)
	defer can()

	// check if already in view
	_, ok := utils.Contains(addr)
	if !ok {
		for _, v := range utils.Views {
			if v != utils.SocketAddr {
				url := fmt.Sprintf("http://%s/broadcast-put/%s", v, addr)
				req, _ := http.NewRequestWithContext(ctx, http.MethodPut, url, nil)
				_, _ = http.DefaultClient.Do(req)
			}
		}
		utils.Views = append(utils.Views, addr)
		resp := putViewSucces{
			Message: "Replica added successfully to the view",
		}
		c.JSON(http.StatusOK, resp)
		return
	}

	resp := putViewFail{
		Message: "",
		Error:   "",
	}
	c.JSON(http.StatusBadRequest, resp)
}
