package config

import (
	"log"

	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	Conn sarama.SyncProducer
}

var (
	Producer *KafkaProducer
)

func ConfigProducer() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Version = sarama.V3_9_0_0
	return config
}

func (producerProvider *KafkaProducer) ConnectKafka() error {
	brokers := []string{"ticket-kafka-kafka-bootstrap:9092"}
	config := ConfigProducer()
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Printf("Could not connect to Kafka: %v", err)
		return err
	}
	log.Println("Connected to Kafka")
	producerProvider.Conn = producer
	return nil
}

func (producerProvider *KafkaProducer) SendMessage(topic string, message string) error {
	var partition int32 = 0
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Partition: partition,
		Value: sarama.StringEncoder(message),
	}
	_, _, err := producerProvider.Conn.SendMessage(msg)
	if err != nil {
		log.Printf("Could not send message to Kafka: %v", err)
		return err
	}
	log.Println("Message sent to Kafka")
	return nil
}

func (producerProvider *KafkaProducer) Close() {
	producerProvider.Conn.Close()
}