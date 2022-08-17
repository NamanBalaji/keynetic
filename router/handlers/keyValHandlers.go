package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/NamanBalaji/keynetic/utils"
	"github.com/gin-gonic/gin"
)

type getSuccesResp struct {
	Exists  bool   `json:"doesExist"`
	Message string `json:"message,omitempty"`
	Value   string `json:"value,omitempty"`
}

type getFailResp struct {
	Exists  bool   `json:"doesExist"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type deleteSuccesResp struct {
	Exists  bool   `json:"doesExist"`
	Message string `json:"message,omitempty"`
}

type delteFailResp struct {
	Exists  bool   `json:"doesExist"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type putRequest struct {
	Value string `json:"value"`
}

type putSuccesResp struct {
	Replaced bool   `json:"replaced"`
	Message  string `json:"message,omitempty"`
}

type putFailResp struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// Handler for GET: /key-value-store/<key>
func GetKVHandler(c *gin.Context) {

	key := c.Param("key")
	val, err := utils.Store.Get(key)
	if err != nil {
		resp := getFailResp{
			Exists:  false,
			Error:   "Key does not exist",
			Message: "Error in GET",
		}
		c.JSON(http.StatusNotFound, resp)
		return
	}

	resp := getSuccesResp{
		Exists:  true,
		Message: "Retrieved successfully",
		Value:   val,
	}
	c.JSON(http.StatusOK, resp)
}

// Handler for DELETE: /key-value-store/<key>
func DeleteKVHandler(c *gin.Context) {

	key := c.Param("key")
	err := utils.Store.Delete(key)
	if err != nil {
		resp := delteFailResp{
			Exists:  false,
			Message: "Error in DELETE",
			Error:   "Key does not exist",
		}
		c.JSON(http.StatusNotFound, resp)
		return
	}

	resp := deleteSuccesResp{
		Message: "Deleted successfully",
		Exists:  true,
	}
	c.JSON(http.StatusOK, resp)
}

// Handler for PUT: /key-value-store/<key>
func PutKVHandler(c *gin.Context) {
	key := c.Param("key")

	if len(key) > 50 {
		resp := putFailResp{
			Error:   "Key is too long",
			Message: "Error in PUT",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	jsonData, err := io.ReadAll(c.Request.Body)
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
