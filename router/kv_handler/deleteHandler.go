package kv_handler

import (
	"net/http"

	"github.com/NamanBalaji/keynetic/utils"
	"github.com/gin-gonic/gin"
)

type deleteSuccesResp struct {
	Exists  bool   `json:"doesExist"`
	Message string `json:"message,omitempty"`
}

type delteFailResp struct {
	Exists  bool   `json:"doesExist"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func DeleteHandler(c *gin.Context) {

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
