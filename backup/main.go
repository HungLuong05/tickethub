package main

import (
	"log"
	"os"

	"tickethub.com/backup/pg"
	"tickethub.com/backup/kafka"

	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
)

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

	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)

	kafka.Consumer = &kafka.KafkaConsumer{}
	err = kafka.Consumer.ConnectKafka()
	if err != nil {
		log.Fatal("Could not connect to Kafka: ", err)
	}

	defer func() {
		if err := kafka.Consumer.Conn.Close(); err != nil {
			log.Fatalf("Failed to close Kafka connection: %v", err)
		}
	} ()

	kafka.Consumer.ConsumeMessage("ticket")
}