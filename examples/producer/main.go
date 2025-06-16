package main

import (
	"encoding/json"
	"fmt"
	"github.com/savin000/audit-log/pkg/kafka"
	"log"
	"time"

	"github.com/IBM/sarama"
)

func main() {
	brokers := []string{"localhost:9092"}
	topic := "audit-logs"

	producer, err := kafka.NewProducer(brokers)
	if err != nil {
		log.Fatalf("Kafka error: %v", err)
	}
	defer func() { _ = producer.Close() }()

	payload := map[string]string{
		"metadata": "test",
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder("my-key"),
		Value: sarama.ByteEncoder(jsonBytes),
	}

	for i := 0; i < 5; i++ {
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.Printf("Failed to send message: %v", err)
		} else {
			fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)
		}

		time.Sleep(5 * time.Second)
	}
}
