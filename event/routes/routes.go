package routes

import (
	"tickethub.com/event/proto"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine, grpcClient proto.EventPermClient) {
	// eventRoutes := server.Group("/api/event")
	// {
	// 	eventRoutes.POST("/", func(c *gin.Context) { CreateEvent(c, grpcClient) })
	// 	eventRoutes.GET("/", GetEvents)
	// 	eventRoutes.GET("/:id", GetEvent)
	// 	eventRoutes.PUT("/:id", UpdateEvent)
	// 	eventRoutes.DELETE("/:id", func(c *gin.Context) { DeleteEvent(c, grpcClient) })
	// }
	server.POST("/api/event", func(c *gin.Context) { CreateEvent(c, grpcClient) })
	server.GET("/api/event", GetEvents)
	server.GET("/api/event/:id", GetEvent)
	server.PUT("/api/event/:id", UpdateEvent)
	server.DELETE("/api/event/:id", func(c *gin.Context) { DeleteEvent(c, grpcClient) })
}