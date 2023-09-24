package config

import (
	"os"
	"strings"
)

type Config interface {
	KafkaBrokers() []string
	KafkaTopicPII() string
}

type config struct {
	kafkaBrokers  []string
	kafkaTopicPII string
}

func (c config) KafkaBrokers() []string {
	return c.kafkaBrokers
}

func (c config) KafkaTopicPII() string {
	return c.kafkaTopicPII
}

func kafkaBrokersFromEnv() []string {
	v, ok := os.LookupEnv("KAFKA_BROKERS")
	if !ok {
		panic("Missing required environment variable KAFKA_BROKERS")
	}

	brokers := strings.Split(v, ",")
	if len(brokers) == 0 {
		panic("Required environment variable KAFKA_BROKERS cannot be empty")
	}

	brokersTrimmed := make([]string, 0, len(brokers))
	for _, broker := range brokers {
		if strings.TrimSpace(broker) == "" {
			panic("Required environment variable KAFKA_BROKERS cannot contain empty brokers")
		}

		brokersTrimmed = append(brokersTrimmed, strings.TrimSpace(broker))
	}

	return brokersTrimmed
}

func kafkaTopicPIIFromEnv() string {
	v, ok := os.LookupEnv("KAFKA_TOPIC_PII")
	if !ok {
		panic("Missing required environment variable KAFKA_TOPIC_PII")
	}

	topicTrimmed := strings.TrimSpace(v)
	if topicTrimmed == "" {
		panic("Required environment variable KAFKA_TOPIC_PII cannot be empty")
	}

	return topicTrimmed
}

func NewConfig() Config {
	return &config{
		kafkaBrokers:  kafkaBrokersFromEnv(),
		kafkaTopicPII: kafkaTopicPIIFromEnv(),
	}
}
