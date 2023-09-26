package config

import (
	"os"
	"strings"
)

type AppEnv int64

const (
	AppEnvDevelopment AppEnv = iota
	AppEnvProduction
	AppEnvTest
)

type Config interface {
	AppEnv() AppEnv
	KafkaBrokers() []string
	KafkaTopicMasked() string
	KafkaTopicPII() string
	PIIMaskerSecret() string
}

type config struct {
	appEnv           AppEnv
	kafkaBrokers     []string
	kafkaTopicMasked string
	kafkaTopicPII    string
	piiMaskerSecret  string
}

func (c *config) AppEnv() AppEnv {
	return c.appEnv
}

func (c *config) KafkaBrokers() []string {
	return c.kafkaBrokers
}

func (c *config) KafkaTopicMasked() string {
	return c.kafkaTopicMasked
}

func (c *config) KafkaTopicPII() string {
	return c.kafkaTopicPII
}

func (c *config) PIIMaskerSecret() string {
	return c.piiMaskerSecret
}

func appEnvFromEnv() AppEnv {
	v, ok := os.LookupEnv("APP_ENV")
	if !ok {
		panic("Missing required environment variable APP_ENV")
	}

	switch strings.ToLower(v) {
	case "dev":
		return AppEnvDevelopment
	case "prod":
		return AppEnvProduction
	case "test":
		return AppEnvTest
	default:
		panic("Invalid value for environment variable APP_ENV")
	}
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

func kafkaTopicMaskedFromEnv() string {
	v, ok := os.LookupEnv("KAFKA_TOPIC_MASKED")
	if !ok {
		panic("Missing required environment variable KAFKA_TOPIC_MASKED")
	}

	topicTrimmed := strings.TrimSpace(v)
	if topicTrimmed == "" {
		panic("Required environment variable KAFKA_TOPIC_MASKED cannot be empty")
	}

	return topicTrimmed
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

func piiMaskerSecretFromEnv() string {
	v, ok := os.LookupEnv("PII_MASKER_SECRET")
	if !ok {
		panic("Missing required environment variable PII_MASKER_SECRET")
	}

	secretTrimmed := strings.TrimSpace(v)
	if secretTrimmed == "" {
		panic("Required environment variable PII_MASKER_SECRET cannot be empty")
	}

	return secretTrimmed
}

func NewConfig() Config {
	return &config{
		appEnv:           appEnvFromEnv(),
		kafkaBrokers:     kafkaBrokersFromEnv(),
		kafkaTopicMasked: kafkaTopicMaskedFromEnv(),
		kafkaTopicPII:    kafkaTopicPIIFromEnv(),
		piiMaskerSecret:  piiMaskerSecretFromEnv(),
	}
}
