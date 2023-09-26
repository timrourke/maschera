package config

import (
	"fmt"
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
	KafkaConsumerGroupID() string
	KafkaTopicMasked() string
	KafkaTopicPII() string
	PIIMaskerSecret() string
}

type config struct {
	appEnv               AppEnv
	kafkaBrokers         []string
	kafkaConsumerGroupID string
	kafkaTopicMasked     string
	kafkaTopicPII        string
	piiMaskerSecret      string
}

func (c *config) AppEnv() AppEnv {
	return c.appEnv
}

func (c *config) KafkaBrokers() []string {
	return c.kafkaBrokers
}

func (c *config) KafkaConsumerGroupID() string {
	return c.kafkaConsumerGroupID
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

func StringFromEnvOrPanic(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("Missing required environment variable %s", key))
	}

	vTrimmed := strings.TrimSpace(v)
	if vTrimmed == "" {
		panic(fmt.Sprintf("Required environment variable %s cannot be empty", key))
	}

	return vTrimmed
}

func appEnvFromEnv() AppEnv {
	v := StringFromEnvOrPanic("APP_ENV")

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
	v := StringFromEnvOrPanic("KAFKA_BROKERS")

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

func kafkaConsumerGroupIDFromEnv() string {
	return StringFromEnvOrPanic("KAFKA_CONSUMER_GROUP_ID")
}

func kafkaTopicMaskedFromEnv() string {
	return StringFromEnvOrPanic("KAFKA_TOPIC_MASKED")
}

func kafkaTopicPIIFromEnv() string {
	return StringFromEnvOrPanic("KAFKA_TOPIC_PII")
}

func piiMaskerSecretFromEnv() string {
	return StringFromEnvOrPanic("PII_MASKER_SECRET")
}

func NewConfig() Config {
	return &config{
		appEnv:               appEnvFromEnv(),
		kafkaBrokers:         kafkaBrokersFromEnv(),
		kafkaConsumerGroupID: kafkaConsumerGroupIDFromEnv(),
		kafkaTopicMasked:     kafkaTopicMaskedFromEnv(),
		kafkaTopicPII:        kafkaTopicPIIFromEnv(),
		piiMaskerSecret:      piiMaskerSecretFromEnv(),
	}
}
