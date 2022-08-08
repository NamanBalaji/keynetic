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

func BroadcastPut(c *gin.Context) {
	ip := c.Param("ip")
	_, ok := utils.Contains(ip)
	if ok {
		c.JSON(http.StatusOK, gin.H{"message": "view already added"})
		return
	} else {
		utils.Views = append(utils.Views, ip)
		c.JSON(http.StatusOK, gin.H{"message": "view added successfully"})
		return
	}
}
