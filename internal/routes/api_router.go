package routes

import (
	"github.com/gin-gonic/gin"
)

func APIRouter(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/api-1", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-1"})
	})
	incomingRoutes.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-2"})
	})
}
