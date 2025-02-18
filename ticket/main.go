package main

import (
	"log"
	"os"

	"tickethub.com/ticket/config"
	"tickethub.com/ticket/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/IBM/sarama"
)

func main() {
	if os.Getenv("KUBERNETES_SERVICE_HOST") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)

	config.DB = &config.RedisConn{}
	err := config.DB.ConnectDB()
	if err != nil {
		log.Fatal("Could not connect to Redis: ", err)
	}

	config.Producer = &config.KafkaProducer{}
	err = config.Producer.ConnectKafka()
	if err != nil {
		log.Fatal("Could not connect to Kafka: ", err)
	}

	err = config.Producer.SendMessage("test-topic", "Hello from Ticket Service")
	if err != nil {
		log.Fatalf("Failed to send message to Kafka: %v", err)
	}

	defer config.Producer.Close()

	server := gin.Default()
	routes.RegisterRoutes(server)
	server.Run(":8002")
}