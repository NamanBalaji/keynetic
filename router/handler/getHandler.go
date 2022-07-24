package handler

import (
	"net/http"

	"github.com/NamanBalaji/keynetic/router/handler/utils"
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

func GetHandler(c *gin.Context) {

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
