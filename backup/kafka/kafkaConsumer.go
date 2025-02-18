package kafka

import (
	"log"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
	"tickethub.com/backup/models"
)

type KafkaConsumer struct {
	Conn sarama.Consumer
}

var (
	Consumer *KafkaConsumer
)

func ConfigConsumer() *sarama.Config {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	return config
}

func (consumerProvider *KafkaConsumer) ConnectKafka() error {
	brokers := []string{"ticket-kafka-kafka-bootstrap:9092"}
	config := ConfigConsumer()
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Printf("Could not connect to Kafka: %v", err)
		return err
	}

	log.Println("Connected to Kafka")
	consumerProvider.Conn = consumer
	return nil
}

func (consumerProvider *KafkaConsumer) ConsumeMessage(topic string) {
	messages, err := consumerProvider.Conn.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	msgCount := 0

	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-messages.Errors():
				log.Printf("Error: %v\n", err.Err)
			case msg := <-messages.Messages():
				msgCount++
				log.Println("Received messages", string(msg.Key), string(msg.Value))

				ticket := models.Ticket{}
				ticket.ParseTicketString(string(msg.Value))
				err := ticket.UpdateTicket()
				if err != nil {
					log.Println("Error updating ticket: ", err)
				}
			case <-signals:
				log.Println("Interrupt is detected")
				doneCh <- struct{}{}	
			}
		}
	} ()

	<-doneCh
	log.Printf("Processed %d messages", msgCount)
}

func (consumerProvider *KafkaConsumer) Consume(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Consumed message: %s\n", message.Value)
		session.MarkMessage(message, "")
	}
	return nil
}

func (consumerProvider *KafkaConsumer) Close() {
	if err := consumerProvider.Conn.Close(); err != nil {
		log.Printf("Could not close Kafka connection: %v", err)
	}
	log.Println("Closed Kafka connection")
}
