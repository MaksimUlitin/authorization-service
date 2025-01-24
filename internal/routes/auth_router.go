package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/maksimUlitin/internal/controllers"
	"github.com/maksimUlitin/internal/lib"
)

func AuthRouter(incomingRoutes *gin.Engine) {
	logger.Info("Initializing authentication routes")

	incomingRoutes.POST("users/signup", func(c *gin.Context) {
		logger.Info("Received signup request",
			"method", "POST",
			"path", "/users/signup",
			"ip", c.ClientIP())
		controllers.SignUp()(c)
	})

	incomingRoutes.POST("users/login", func(c *gin.Context) {
		logger.Info("Received login request",
			"method", "POST",
			"path", "/users/login",
			"ip", c.ClientIP())
		controllers.Login()(c)
	})

	logger.Info("Authentication routes initialized successfully")
}
