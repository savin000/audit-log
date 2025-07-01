package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"log"
)

type Producer struct {
	SyncProducer sarama.SyncProducer
}

func NewProducer(brokers []string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{SyncProducer: producer}, err
}

func (p *Producer) SendMessage(msg *sarama.ProducerMessage) {
	partition, offset, err := p.SyncProducer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	} else {
		fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)
	}
}

func (p *Producer) Close() error {
	return p.SyncProducer.Close()
}
