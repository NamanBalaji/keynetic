package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/NamanBalaji/keynetic/types"
	"github.com/NamanBalaji/keynetic/utils"
	"github.com/gin-gonic/gin"
)

// Handler for GET: /key-value-store/<key>
func GetKVHandler(c *gin.Context) {

	key := c.Param("key")
	val, err := utils.Store.Get(key)
	if err != nil {
		resp := types.GetFailResp{
			Exists:  false,
			Error:   "Key does not exist",
			Message: "Error in GET",
		}
		c.JSON(http.StatusNotFound, resp)
		return
	}

	resp := types.GetSuccesResp{
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
		resp := types.DeleteFailResp{
			Exists:  false,
			Message: "Error in DELETE",
			Error:   "Key does not exist",
		}
		c.JSON(http.StatusNotFound, resp)
		return
	}

	resp := types.DeleteSuccesResp{
		Message: "Deleted successfully",
		Exists:  true,
	}
	c.JSON(http.StatusOK, resp)
}

// Handler for PUT: /key-value-store/<key>
func PutKVHandler(c *gin.Context) {
	key := c.Param("key")

	if len(key) > 50 {
		resp := types.PutFailResp{
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

	var body types.PutRequest
	err = json.Unmarshal(jsonData, &body)
	if err != nil {
		log.Printf("invalid body format [ERROR]: %s", err)
		return
	}

	if body.Value == "" {
		resp := types.PutFailResp{
			Error:   "Value is missing",
			Message: "Error in PUT",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	replaced, _ := utils.Store.Put(key, body.Value)
	if replaced {
		resp := types.PutSuccesResp{
			Message:  "Updated successfully",
			Replaced: replaced,
		}
		c.JSON(http.StatusOK, resp)
		return
	} else {
		resp := types.PutSuccesResp{
			Message:  "Added successfully",
			Replaced: replaced,
		}
		c.JSON(http.StatusCreated, resp)
		return
	}
}
