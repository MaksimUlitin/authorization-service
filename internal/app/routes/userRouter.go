package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/maksimUlitin/internal/app/controllers"
	"github.com/maksimUlitin/internal/app/middleware"
)

func UserRouter(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticante())
	incomingRoutes.GET("/users", controllers.GetUsers())
	incomingRoutes.GET("/users/:users_id", controllers.GetUser())
}
