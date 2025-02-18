package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.GET("/api/ticket/:id", GetTicket)
	server.POST("/api/ticket", CreateTicket)
	server.POST("/api/ticket/:id", UpdateTicket)
}