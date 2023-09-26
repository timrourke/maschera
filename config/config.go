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
	JSONFieldsWithPII() []string
	KafkaBrokers() []string
	KafkaConsumerGroupID() string
	KafkaTopicMasked() string
	KafkaTopicPII() string
	PIIMaskerSecret() string
}

type config struct {
	appEnv               AppEnv
	jsonFieldsWithPII    []string
	kafkaBrokers         []string
	kafkaConsumerGroupID string
	kafkaTopicMasked     string
	kafkaTopicPII        string
	piiMaskerSecret      string
}

func (c *config) AppEnv() AppEnv {
	return c.appEnv
}

func (c *config) JSONFieldsWithPII() []string {
	return c.jsonFieldsWithPII
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

func StringSliceFromCommaSeparatedEnvOrPanic(key string) []string {
	v := StringFromEnvOrPanic(key)

	s := strings.Split(v, ",")
	if len(s) == 0 {
		panic(fmt.Sprintf("Required environment variable %s cannot be empty", key))
	}

	sTrimmed := make([]string, 0, len(s))
	for _, e := range s {
		if strings.TrimSpace(e) == "" {
			panic(fmt.Sprintf("Required environment variable %s cannot contain empty elements", key))
		}

		sTrimmed = append(sTrimmed, strings.TrimSpace(e))
	}

	return sTrimmed
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

func jsonFieldsWithPIIFromEnv() []string {
	return StringSliceFromCommaSeparatedEnvOrPanic("JSON_FIELDS_WITH_PII")
}

func kafkaBrokersFromEnv() []string {
	return StringSliceFromCommaSeparatedEnvOrPanic("KAFKA_BROKERS")
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
		jsonFieldsWithPII:    jsonFieldsWithPIIFromEnv(),
		kafkaBrokers:         kafkaBrokersFromEnv(),
		kafkaConsumerGroupID: kafkaConsumerGroupIDFromEnv(),
		kafkaTopicMasked:     kafkaTopicMaskedFromEnv(),
		kafkaTopicPII:        kafkaTopicPIIFromEnv(),
		piiMaskerSecret:      piiMaskerSecretFromEnv(),
	}
}
