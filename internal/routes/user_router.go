package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/maksimUlitin/internal/controllers"
	"github.com/maksimUlitin/internal/lib"
	"github.com/maksimUlitin/internal/middleware"
)

func UserRouter(incomingRoutes *gin.Engine) {
	logger.Info("Initializing user routes")

	incomingRoutes.Use(middleware.Authenticate())
	logger.Info("Authentication middleware applied to user routes")

	incomingRoutes.GET("/users", func(c *gin.Context) {
		logger.Info("Received request to get all users",
			"method", "GET",
			"path", "/users",
			"ip", c.ClientIP())
		controllers.GetUsers()(c)
	})

	incomingRoutes.GET("/users/:users_id", func(c *gin.Context) {
		logger.Info("Received request to get specific user",
			"method", "GET",
			"path", "/users/:users_id",
			"user_id", c.Param("users_id"),
			"ip", c.ClientIP())
		controllers.GetUser()(c)
	})

	logger.Info("User routes initialized successfully")
}
