package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/savin000/audit-log/internal/clickhouse"
	"log"
)

type ConsumerGroupHandler struct {
	Ch          *clickhouse.Client
	DLQEnabled  bool
	DLQTopic    string
	DLQProducer *Producer
}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error { return nil }

func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(msg.Value), msg.Timestamp, msg.Topic)

		var auditLog clickhouse.AuditLog
		err := json.Unmarshal(msg.Value, &auditLog)
		if err != nil {
			if h.DLQEnabled {
				msg := &sarama.ProducerMessage{
					Topic: h.DLQTopic,
					Key:   sarama.StringEncoder(msg.Key),
					Value: sarama.StringEncoder(msg.Value),
				}

				h.DLQProducer.SendMessage(msg)
			} else {
				return fmt.Errorf("error reading message: %w", err)
			}
		} else {
			session.MarkMessage(msg, "")
		}

		err = h.Ch.AddAuditLog(auditLog)
		if err != nil {
			return fmt.Errorf("failed to insert into ClickHouse: %w", err)
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
