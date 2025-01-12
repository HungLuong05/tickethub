package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.POST("/register", Register)
	server.POST("/login", Login)
	server.POST("/verify", Verify)
}