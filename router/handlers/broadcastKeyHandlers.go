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

func BroadcastKeyPut(c *gin.Context) {

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

	for key, val := range body.CausalMetadata {
		utils.Vc[key] = val
	}
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
	if replaced {
		resp := types.PutSuccesResp{
			Message:        "Updated successfully",
			Replaced:       replaced,
			CausalMetadata: body.CausalMetadata,
		}
		c.JSON(http.StatusOK, resp)
		return
	} else {
		resp := types.PutSuccesResp{
			Message:        "Added successfully",
			Replaced:       replaced,
			CausalMetadata: body.CausalMetadata,
		}
		c.JSON(http.StatusCreated, resp)
		return
	}
}

func BroadcastKeyDelete(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("invalid request body [ERROR]: %s", err)
		return
	}

	var body types.DeleteReq
	err = json.Unmarshal(jsonData, &body)
	if err != nil {
		log.Printf("invalid body format [ERROR]: %s", err)
		return
	}

	for key, val := range body.CausalMetadata {
		utils.Vc[key] = val
	}
	key := c.Param("key")
	err = utils.Store.Delete(key)
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

	resp := types.DeleteSuccesResp{
		Message:        "Deleted successfully",
		Exists:         true,
		CausalMetadata: body.CausalMetadata,
	}
	c.JSON(http.StatusOK, resp)
}
