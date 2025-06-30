package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	ServerPort         uint32   `env:"SERVER_PORT" envDefault:"8080"`
	KafkaAddresses     []string `env:"KAFKA_ADDRESSES" envDefault:"localhost:9092" envSeparator:","`
	KafkaGroupID       string   `env:"KAFKA_GROUP_ID" envDefault:"audit-log-default-consumer"`
	KafkaTopics        []string `env:"KAFKA_TOPICS" envDefault:"default.audit.log" envSeparator:","`
	ClickhouseHost     string   `env:"CLICKHOUSE_HOST" envDefault:"localhost"`
	ClickhousePort     uint32   `env:"CLICKHOUSE_PORT" envDefault:"9000"`
	ClickhouseDatabase string   `env:"CLICKHOUSE_DATABASE" envDefault:"default"`
	ClickhouseUsername string   `env:"CLICKHOUSE_USERNAME" envDefault:"clickuser"`
	ClickhousePassword string   `env:"CLICKHOUSE_PASSWORD" envDefault:"clickpassword"`
}

func Get() (*Config, error) {
	cfg := Config{}

	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
