package kv_handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/NamanBalaji/keynetic/utils"
	"github.com/gin-gonic/gin"
)

type putSuccesResp struct {
	Replaced bool   `json:"replaced"`
	Message  string `json:"message,omitempty"`
}

type putFailResp struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type putRequest struct {
	Value string `json:"value"`
}

func PutHandler(c *gin.Context) {
	key := c.Param("key")

	if len(key) > 50 {
		resp := putFailResp{
			Error:   "Key is too long",
			Message: "Error in PUT",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("invalid request body [ERROR]: %s", err)
		return
	}

	var body putRequest
	err = json.Unmarshal(jsonData, &body)
	if err != nil {
		log.Printf("invalid body format [ERROR]: %s", err)
		return
	}

	if body.Value == "" {
		resp := putFailResp{
			Error:   "Value is missing",
			Message: "Error in PUT",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	replaced, _ := utils.Store.Put(key, body.Value)
	if replaced {
		resp := putSuccesResp{
			Message:  "Updated successfully",
			Replaced: replaced,
		}
		c.JSON(http.StatusOK, resp)
		return
	} else {
		resp := putSuccesResp{
			Message:  "Added successfully",
			Replaced: replaced,
		}
		c.JSON(http.StatusCreated, resp)
		return
	}
}
