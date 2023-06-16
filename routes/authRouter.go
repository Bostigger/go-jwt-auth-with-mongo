package routes

import (
	"github.com/bostigger/go-jwt-mongo/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("api/user/create", controllers.CreateUser)
	incomingRoutes.POST("api/user/login", controllers.LoginUser)
}
