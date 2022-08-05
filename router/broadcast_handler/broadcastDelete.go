package broadcast_handler

import (
	"net/http"

	"github.com/NamanBalaji/keynetic/utils"
	"github.com/gin-gonic/gin"
)

func BroadcastDelete(c *gin.Context) {
	ip := c.Param("ip")
	err := utils.RemoveFromView(ip)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "view deleted"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "view not present"})
		return
	}
}
