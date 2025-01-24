package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/maksimUlitin/internal/lib"
)

func APIRouter(incomingRoutes *gin.Engine) {
	logger.Info("Initializing API routes")

	incomingRoutes.GET("/api-1", func(c *gin.Context) {
		logger.Info("Access to api-1",
			"method", "GET",
			"path", "/api-1",
			"ip", c.ClientIP())
		c.JSON(200, gin.H{"success": "Access granted for api-1"})
	})

	incomingRoutes.GET("/api-2", func(c *gin.Context) {
		logger.Info("Access to api-2",
			"method", "GET",
			"path", "/api-2",
			"ip", c.ClientIP())
		c.JSON(200, gin.H{"success": "Access granted for api-2"})
	})

	logger.Info("API routes initialized successfully")
}
