package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/maksimUlitin/internal/controllers"
	"github.com/maksimUlitin/internal/middleware"
)

func UserRouter(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controllers.GetUsers())
	incomingRoutes.GET("/users/:users_id", controllers.GetUser())
}
