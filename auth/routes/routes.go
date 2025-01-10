package routes

import (
	// "bluebid.com/auth/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	// authenticated := server.Group("/")
	// authenticated.Use(middlewares.Authenticate)
	// // authenticated.GET("/verify", verify)

	server.POST("/register", Register)
	server.POST("/login", Login)
}