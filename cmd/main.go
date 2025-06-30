package main

import (
	"context"
	"github.com/savin000/audit-log/internal/server"
	"github.com/savin000/audit-log/internal/server/handlers"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/savin000/audit-log/config"
	"github.com/savin000/audit-log/internal/clickhouse"
	"github.com/savin000/audit-log/internal/kafka"
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

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for {
			err := consumerGroup.Consume(ctx, cfg.KafkaTopics, handler)
			if err != nil {
				log.Fatalf("Error from consumerGroup: %v", err)
			}
		}
	}()

	go func() {
		h := &handlers.Handler{Ch: client}
		httpServer := server.New(cfg.ServerPort, h)
		err := httpServer.ListenAndServe()
		if err != nil {
			log.Fatalf("Http server error: %v", err)
		}
	}()

	<-stop
	wg.Wait()
}
