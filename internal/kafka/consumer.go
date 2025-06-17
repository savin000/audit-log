package kafka

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/savin000/audit-log/internal/clickhouse"
	"log"
)

type ConsumerGroupHandler struct {
	Ch *clickhouse.Client
}

func (handler *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error { return nil }

func (handler *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (handler *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(msg.Value), msg.Timestamp, msg.Topic)

		var logEntry clickhouse.AuditLog
		err := json.Unmarshal(msg.Value, &logEntry)

		if err != nil {
			log.Fatalf("Error reading message: %v", err)
		}

		err = handler.Ch.AddAuditLog(logEntry)
		if err != nil {
			log.Fatalf("Failed to insert into ClickHouse: %v", err)
		} else {
			session.MarkMessage(msg, "")
		}
	}
	return nil
}

func NewConsumerGroup(addrs []string, groupID string) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.AutoCommit.Enable = true

	consumerGroup, err := sarama.NewConsumerGroup(addrs, groupID, config)
	if err != nil {
		return nil, err
	}

	return consumerGroup, nil
}
