package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/NamanBalaji/keynetic/requests"
	"github.com/NamanBalaji/keynetic/types"
	"github.com/NamanBalaji/keynetic/utils"
	"github.com/gin-gonic/gin"
)

// Handler for GET: /key-value-store/<key>
func GetKVHandler(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("invalid request body [ERROR]: %s", err)
		return
	}

	var body types.PutRequest
	err = json.Unmarshal(jsonData, &body)
	if err != nil {
		log.Printf("invalid body format [ERROR]: %s", err)
		return
	}

	syncStoreAndVc(utils.StringToMap(body.CausalMetadata))
	key := c.Param("key")
	val, err := utils.Store.Get(key)
	if err != nil {
		resp := types.GetFailResp{
			Exists:         false,
			Error:          "Key does not exist",
			Message:        "Error in GET",
			CausalMetadata: body.CausalMetadata,
		}
		c.JSON(http.StatusNotFound, resp)
		return
	}

	resp := types.GetSuccesResp{
		Exists:         true,
		Message:        "Retrieved successfully",
		Value:          val,
		CausalMetadata: utils.MapToString(utils.Vc),
	}
	c.JSON(http.StatusOK, resp)

}

// Handler for DELETE: /key-value-store/<key>
func DeleteKVHandler(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("invalid request body [ERROR]: %s", err)
		return
	}

	var body types.PutRequest
	err = json.Unmarshal(jsonData, &body)
	if err != nil {
		log.Printf("invalid body format [ERROR]: %s", err)
		return
	}

	syncStoreAndVc(utils.StringToMap(body.CausalMetadata))

	key := c.Param("key")

	_, err = utils.Store.Get(key)
	if err != nil {
		resp := types.DeleteFailResp{
			Exists:         false,
			Message:        "Error in DELETE",
			Error:          "Key does not exist",
			CausalMetadata: body.CausalMetadata,
		}
		c.JSON(http.StatusNotFound, resp)
		return
	}
	utils.Vc[utils.View.SocketAddr] = utils.Vc[utils.View.SocketAddr] + 1
	incrementVCDeleteSteps(key, utils.MapToString(utils.Vc))
	utils.Store.Delete(key)

	resp := types.DeleteSuccesResp{
		Message:        "Deleted successfully",
		Exists:         true,
		CausalMetadata: utils.MapToString(utils.Vc),
	}
	c.JSON(http.StatusOK, resp)
}

// Handler for PUT: /key-value-store/<key>
func PutKVHandler(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("invalid request body [ERROR]: %s", err)
		return
	}

	var body types.PutRequest
	err = json.Unmarshal(jsonData, &body)
	if err != nil {
		log.Printf("invalid body format [ERROR]: %s", err)
		return
	}

	syncStoreAndVc(utils.StringToMap(body.CausalMetadata))

	key := c.Param("key")

	if len(key) > 50 {
		resp := types.PutFailResp{
			Error:          "Key is too long",
			Message:        "Error in PUT",
			CausalMetadata: body.CausalMetadata,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if body.Value == "" {
		resp := types.PutFailResp{
			Error:          "Value is missing",
			Message:        "Error in PUT",
			CausalMetadata: body.CausalMetadata,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	replaced, _ := utils.Store.Put(key, body.Value)

	utils.Vc[utils.View.SocketAddr] = utils.Vc[utils.View.SocketAddr] + 1
	incrementVCPutSteps(key, body.Value, utils.MapToString(utils.Vc))
	if replaced {
		resp := types.PutSuccesResp{
			Message:        "Updated successfully",
			Replaced:       replaced,
			CausalMetadata: utils.MapToString(utils.Vc),
		}
		c.JSON(http.StatusOK, resp)
		return
	} else {
		resp := types.PutSuccesResp{
			Message:        "Added successfully",
			Replaced:       replaced,
			CausalMetadata: utils.MapToString(utils.Vc),
		}
		c.JSON(http.StatusCreated, resp)
		return
	}
}

func syncStoreAndVc(causalMetadata map[string]int) {
	for key, val := range causalMetadata {
		if utils.Vc[key] < val {
			for _, replica := range utils.View.Views {
				if replica != utils.View.SocketAddr {
					var storeRes types.GetStoreResponse
					var vectorClockRes types.GetVectorClockResponse

					res, err1 := requests.GetKeyValueStore(replica)
					if err1 == nil {
						jsonData, _ := io.ReadAll(res.Body)
						err1 = json.Unmarshal(jsonData, &storeRes)
					}

					res, err2 := requests.GetVectorClock(replica)
					if err2 == nil {
						jsonData, _ := io.ReadAll(res.Body)
						err2 = json.Unmarshal(jsonData, &vectorClockRes)
					}

					err3 := requests.BroadcastPutView(replica, utils.View.SocketAddr)
					if err3 == nil && err2 == nil && err1 == nil {
						utils.SetVectorClock(vectorClockRes.VectorClock)
						utils.SetStore(storeRes.Store)
						break
					}
				}
			}
		}
	}
}

func incrementVCPutSteps(key, val string, causalMetadata string) {
	var down []string
	for _, replica := range utils.View.Views {
		if replica != utils.View.SocketAddr {
			err := requests.BroadcastPutKey(key, val, replica, causalMetadata)
			if err != nil {
				down = append(down, replica)
			}
		}
	}
	for _, d := range down {
		utils.View.RemoveFromView(d)
	}
	for _, replica := range utils.View.Views {
		requests.BroadcastDeleteView(replica, down...)
	}
}

func incrementVCDeleteSteps(key string, causalMetadata string) {
	var down []string
	for _, replica := range utils.View.Views {
		if replica != utils.View.SocketAddr {
			err := requests.BroadcastDeleteKey(key, replica, causalMetadata)
			if err != nil {
				down = append(down, replica)
			}
		}
	}
	for _, d := range down {
		utils.View.RemoveFromView(d)
	}
	for _, replica := range utils.View.Views {
		requests.BroadcastDeleteView(replica, down...)
	}
}
