package main

import (
	"context"
	"github.com/savin000/audit-log/internal/clickhouse"
	"github.com/savin000/audit-log/internal/kafka"
	"log"
)

func main() {
	brokers := []string{"localhost:9092"}
	groupID := "default-group"
	topic := []string{"audit-logs"}
	cfg := clickhouse.Config{
		Host:     "localhost",
		Port:     9000,
		Database: "default",
		Username: "clickuser",
		Password: "clickpassword",
	}

	client, err := clickhouse.New(cfg)
	if err != nil {
		log.Fatalf("ClickHouse connection error: %v", err)
	}
	defer func() { _ = client.Close() }()

	err = client.CreateAuditLogTable()
	if err != nil {
		log.Fatalf("Failed to create AuditLog table: %v", err)
	}

	consumerGroup, err := kafka.NewConsumerGroup(brokers, groupID)

	if err != nil {
		log.Fatalf("Kafka error: %v", err)
	}
	defer func() { _ = consumerGroup.Close() }()

	ctx := context.Background()
	handler := &kafka.ConsumerGroupHandler{Ch: client}

	for {
		err := consumerGroup.Consume(ctx, topic, handler)
		if err != nil {
			log.Printf("Error from consumerGroup: %v", err)
			panic(err)
		}
	}
}
