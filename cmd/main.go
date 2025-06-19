package main

import (
	"context"
	"github.com/savin000/audit-log/config"
	"github.com/savin000/audit-log/internal/clickhouse"
	"github.com/savin000/audit-log/internal/kafka"
	"log"
)

func main() {
	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("Failed to get environment variables: %v", err)
	}

	clickhouseCfg := clickhouse.Config{
		Host:     cfg.ClickhouseHost,
		Port:     cfg.ClickhousePort,
		Database: cfg.ClickhouseDatabase,
		Username: cfg.ClickhouseUsername,
		Password: cfg.ClickhousePassword,
	}

	client, err := clickhouse.New(clickhouseCfg)
	if err != nil {
		log.Fatalf("ClickHouse connection error: %v", err)
	}
	defer func() { _ = client.Close() }()

	err = client.CreateAuditLogTable()
	if err != nil {
		log.Fatalf("Failed to create AuditLog table: %v", err)
	}

	consumerGroup, err := kafka.NewConsumerGroup(cfg.KafkaAddresses, cfg.KafkaGroupID)
	if err != nil {
		log.Fatalf("Kafka error: %v", err)
	}
	defer func() { _ = consumerGroup.Close() }()

	ctx := context.Background()
	handler := &kafka.ConsumerGroupHandler{Ch: client}
	for {
		err := consumerGroup.Consume(ctx, cfg.KafkaTopics, handler)
		if err != nil {
			log.Printf("Error from consumerGroup: %v", err)
			panic(err)
		}
	}
}
