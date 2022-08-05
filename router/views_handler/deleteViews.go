package views_handler

import (
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

}
