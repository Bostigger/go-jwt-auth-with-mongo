package main

import (
	"github.com/bostigger/go-jwt-mongo/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(cors.Default())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Go JWT API blazing"})
	})
	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	err := router.Run(":" + port)
	if err != nil {
		return
	}
}
