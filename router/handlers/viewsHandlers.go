package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/NamanBalaji/keynetic/requests"
	"github.com/NamanBalaji/keynetic/utils"
	"github.com/gin-gonic/gin"
)

type getViewSucces struct {
	Message string `json:"message,omitempty"`
	View    string `json:"view,omitempty"`
}

type deleteViewSucces struct {
	Message string `json:"message,omitempty"`
}

type deleteViewFail struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type putViewSucces struct {
	Message string `json:"message,omitempty"`
}

type putViewFail struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// Handler function for GET: /key-value-store-view
func GetViewHandler(c *gin.Context) {
	// send heartbeats to check for down replicas
	downInstances := requests.SendHeartbeats()

	// remove down replicas from view
	for _, instance := range downInstances {
		utils.View.RemoveFromView(instance)
	}

	// broadcast delete of down replicas to alive ones
	for _, v := range utils.View.Views {
		if v != utils.View.SocketAddr {
			requests.BroadcastDeleteView(v, downInstances...)
		}
	}

	resp := getViewSucces{
		Message: "View retrieved successfully",
		View:    fmt.Sprint(strings.Join(utils.View.Views[:], ",")),
	}
	c.JSON(http.StatusOK, resp)
}

// Handler function for DELETE: /key-value-store-view
func DeleteViewHandler(c *gin.Context) {

	addr := getAddrFromBody(c)

	// replicas that do not respond back to delete requests
	var downInstances []string

	_, ok := utils.View.Contains(addr)

	//if addr not present in views return error response
	if !ok {
		resp := deleteViewFail{
			Message: "Socket address does not exist in the view",
			Error:   "Error in DELETE",
		}
		c.JSON(http.StatusNotFound, resp)
		return
	}

	// remove from personal view
	utils.View.RemoveFromView(addr)
	// broadcast request to delete address from views of other replicas
	for _, v := range utils.View.Views {
		if v != utils.View.SocketAddr {
			down, err := requests.BroadcastDeleteView(v, addr)
			if err != nil {
				downInstances = append(downInstances, down)
			}
		}
	}

	removeDeadReplicasAndBroadcastDelete(downInstances)

	resp := deleteViewSucces{
		Message: "Replica deleted successfully from the view",
	}
	c.JSON(http.StatusOK, resp)
}

// Handler function for PUT: /key-value-store-view
func PutViewHandler(c *gin.Context) {
	addr := getAddrFromBody(c)

	var downInstances []string

	// if address already present in view return error response
	_, ok := utils.View.Contains(addr)
	if ok {
		resp := putViewFail{
			Message: "Socket address already exists in the view",
			Error:   "Error in PUT",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	for _, replica := range utils.View.Views {
		if replica != utils.View.SocketAddr {
			err := requests.BroadcastPutView(replica, addr)
			if err != nil {
				downInstances = append(downInstances, replica)
			}
		}
	}

	removeDeadReplicasAndBroadcastDelete(downInstances)

	utils.View.AddToView(addr)
	resp := putViewSucces{
		Message: "Replica added successfully to the view",
	}
	c.JSON(http.StatusOK, resp)
}

// helper function to parse body
func getAddrFromBody(c *gin.Context) string {
	var addr string
	body, _ := io.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()
	strBody := string(body[:])
	json.NewDecoder(strings.NewReader(strBody)).Decode(&addr)
	return addr
}

func removeDeadReplicasAndBroadcastDelete(dead []string) {
	for _, d := range dead {
		utils.View.RemoveFromView(d)
	}
	for _, v := range utils.View.Views {
		if v != utils.View.SocketAddr {
			requests.BroadcastDeleteView(v, dead...)
		}
	}
}
