package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/maksimUlitin/internal/app/controllers"
)

func AuthRouter(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("users/singUp", controllers.SignUp())
	incomingRoutes.POST("users/login", controllers.Login())
}
