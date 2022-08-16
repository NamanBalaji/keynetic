package handlers

import (
	"net/http"

	"github.com/NamanBalaji/keynetic/utils"
	"github.com/gin-gonic/gin"
)

// Handler for DELETE: /broadcast-view/<addr>
func BroadcastViewDelete(c *gin.Context) {
	ip := c.Param("ip")
	err := utils.View.RemoveFromView(ip)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "view deleted"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "view not present"})
		return
	}
}

// Handler for PUT: /broadcast-view/<addr>
func BroadcastViewPut(c *gin.Context) {
	ip := c.Param("ip")
	_, ok := utils.View.Contains(ip)
	if ok {
		c.JSON(http.StatusOK, gin.H{"message": "view already added"})
		return
	} else {
		utils.View.AddToView(ip)
		c.JSON(http.StatusOK, gin.H{"message": "view added successfully"})
		return
	}
}
