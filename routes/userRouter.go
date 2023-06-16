package routes

import (
	"github.com/bostigger/go-jwt-mongo/controllers"
	"github.com/bostigger/go-jwt-mongo/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middlewares.Authenticate)
	incomingRoutes.GET("api/users/get-user/:user_id", controllers.GetUserByID)
	incomingRoutes.GET("api/users/get-user", controllers.GetUsers)
}
