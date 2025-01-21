package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.POST("/register", Register)
	server.POST("/login", Login)
	server.GET("/verify", Verify)
	server.GET("/perm/:id", VerifyPerms)
	server.POST("/perm", AddPerm)
	server.DELETE("/perm", RemovePerm)
}