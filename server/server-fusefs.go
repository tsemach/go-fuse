package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func serverFuseFS(_ context.Context, r *gin.Engine) {
	r.POST("/api/v1/filesystems", func(c *gin.Context) {
		var jsonData FuseFSMountPoint
			
		if err := c.ShouldBindJSON(&jsonData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"mountpoint":  jsonData.Mountpoint,
			"targetpath": jsonData.Targetpath,
		})

		fmt.Println(jsonData)
		c.JSON(http.StatusOK, gin.H{"message": "[routeFuseFS]: POST:/api/v1/filesystems", "data": jsonData})
	})

}
