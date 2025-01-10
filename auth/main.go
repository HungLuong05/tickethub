package main

import (
	"log"

	"bluebid.com/auth/config"
	"bluebid.com/auth/routes"
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

// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"

// 	"github.com/jackc/pgx/v5"
// )

// func main() {
// 	// Connection string
// 	connStr := "postgres://youruser:yourpassword@localhost:5432/yourdb"

// 	// Connect to PostgreSQL
// 	conn, err := pgx.Connect(context.Background(), connStr)
// 	if err != nil {
// 		log.Fatalf("Unable to connect to database: %v\n", err)
// 	}
// 	defer conn.Close(context.Background())

// 	// Test the connection
// 	var greeting string
// 	err = conn.QueryRow(context.Background(), "SELECT 'Hello, PostgreSQL!'").Scan(&greeting)
// 	if err != nil {
// 		log.Fatalf("Query failed: %v\n", err)
// 	}
// 	fmt.Println(greeting)
// }