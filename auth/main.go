package main

import (
	"log"

	"tickethub.com/auth/config"
	"tickethub.com/auth/routes"
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

	server := gin.Default()
	err = server.SetTrustedProxies([]string{"35.20.176.45", "127.0.0.1"})
	if err != nil {
		log.Fatal("Could not set trusted proxies: ", err)
	}

	routes.RegisterRoutes(server)
	server.Run(":8000")
}