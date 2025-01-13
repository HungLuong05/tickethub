package main

import (
	"log"
	"sync"

	"tickethub.com/auth/config"
	"tickethub.com/auth/routes"
	"tickethub.com/auth/grpc"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

	pg.DB = &pg.PGPool{}

	err = pg.DB.ConnectDB()
	if err != nil {
		log.Fatal("Could not connect to the database: ", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		server := gin.Default()
		err = server.SetTrustedProxies([]string{"127.0.0.1"})
		if err != nil {
			log.Fatal("Could not set trusted proxies: ", err)
		}

		routes.RegisterRoutes(server)
		server.Run(":8000")
	} ()

	go grpc.StartGrpcServer(&wg)

	wg.Wait()
}