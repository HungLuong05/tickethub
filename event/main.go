package main

import (
	"log"
	"os"

	"tickethub.com/event/config"
	"tickethub.com/event/proto"
	"tickethub.com/event/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var grpcClient proto.EventPermClient

func main() {
	if os.Getenv("KUBERNETES_SERVICE_HOST") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	pg.DB = &pg.PGPool{}

	err := pg.DB.ConnectDB()
	if err != nil {
		log.Fatal("Could not connect to the database: ", err)
	}

	conn, err := grpc.NewClient("auth:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	} else {
		log.Println("Connected to gRPC server")
	}
	defer conn.Close()
	grpcClient = proto.NewEventPermClient(conn)

	server := gin.Default()
	// err = server.SetTrustedProxies([]string{"127.0.0.1"})
	// if err != nil {
	// 	log.Fatal("Could not set trusted proxies: ", err)
	// }

	routes.RegisterRoutes(server, grpcClient)
	server.Run(":8001")
}