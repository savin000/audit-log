package config

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
)

const (
	serverPort         = uint32(9090)
	kafkaAddresses     = "broker1:9092,broker2:9092"
	kafkaGroupID       = "test-group-id"
	kafkaTopics        = "topic1,topic2"
	dlqEnabled         = "true"
	dlqTopic           = "dlq.topic"
	clickhouseHost     = "clickhouseHost"
	clickhousePort     = uint32(9001)
	clickhouseDatabase = "database"
	clickhouseUsername = "username"
	clickhousePassword = "password"
)

func TestGetConfig(t *testing.T) {
	t.Setenv("SERVER_PORT", strconv.Itoa(int(serverPort)))
	t.Setenv("KAFKA_ADDRESSES", kafkaAddresses)
	t.Setenv("KAFKA_GROUP_ID", kafkaGroupID)
	t.Setenv("KAFKA_TOPICS", kafkaTopics)
	t.Setenv("DLQ_ENABLED", dlqEnabled)
	t.Setenv("DLQ_TOPIC", dlqTopic)
	t.Setenv("CLICKHOUSE_HOST", clickhouseHost)
	t.Setenv("CLICKHOUSE_PORT", strconv.Itoa(int(clickhousePort)))
	t.Setenv("CLICKHOUSE_DATABASE", clickhouseDatabase)
	t.Setenv("CLICKHOUSE_USERNAME", clickhouseUsername)
	t.Setenv("CLICKHOUSE_PASSWORD", clickhousePassword)

	cfg, err := Get()

	assert.NoError(t, err)
	assert.Equal(t, serverPort, cfg.ServerPort)
	assert.Equal(t, strings.Split(kafkaAddresses, ","), cfg.KafkaAddresses)
	assert.Equal(t, kafkaGroupID, cfg.KafkaGroupID)
	assert.Equal(t, strings.Split(kafkaTopics, ","), cfg.KafkaTopics)
	assert.True(t, cfg.DLQEnabled)
	assert.Equal(t, dlqTopic, cfg.DLQTopic)
	assert.Equal(t, clickhouseHost, cfg.ClickhouseHost)
	assert.Equal(t, clickhousePort, cfg.ClickhousePort)
	assert.Equal(t, clickhouseDatabase, cfg.ClickhouseDatabase)
	assert.Equal(t, clickhouseUsername, cfg.ClickhouseUsername)
	assert.Equal(t, clickhousePassword, cfg.ClickhousePassword)
}
