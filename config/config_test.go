package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/timrourke/maschera/m/v2/config"
)

func stubValidEnvVars(t *testing.T) {
	t.Setenv("APP_ENV", "prod")
	t.Setenv("KAFKA_BROKERS", "broker1,broker2,broker3")
	t.Setenv("KAFKA_CONSUMER_GROUP_ID", "some-consumer-group-id")
	t.Setenv("KAFKA_TOPIC_MASKED", "masked_data")
	t.Setenv("KAFKA_TOPIC_PII", "pii_data")
	t.Setenv("PII_MASKER_SECRET", "some-secret-key")
}

func TestStringFromEnvOrPanic(t *testing.T) {
	t.Run("Returns value from environment variable", func(t *testing.T) {
		t.Setenv("SOME_ENV_VAR", "some-value")

		assert.Equal(t, "some-value", config.StringFromEnvOrPanic("SOME_ENV_VAR"))
	})

	t.Run("Panics if environment variable is missing", func(t *testing.T) {
		t.Setenv("SOME_ENV_VAR", "")

		assert.Panicsf(t, func() { config.StringFromEnvOrPanic("SOME_ENV_VAR") }, "Missing required environment variable SOME_ENV_VAR")
	})

	t.Run("Panics if environment variable is empty", func(t *testing.T) {
		t.Setenv("SOME_ENV_VAR", " ")

		assert.Panicsf(t, func() { config.StringFromEnvOrPanic("SOME_ENV_VAR") }, "Required environment variable SOME_ENV_VAR cannot be empty")
	})
}

func TestConfig(t *testing.T) {
	t.Run("Exposes config values", func(t *testing.T) {
		stubValidEnvVars(t)

		c := config.NewConfig()

		assert.Equal(t, config.AppEnvProduction, c.AppEnv())
		assert.Equal(t, []string{"broker1", "broker2", "broker3"}, c.KafkaBrokers())
		assert.Equal(t, "some-consumer-group-id", c.KafkaConsumerGroupID())
		assert.Equal(t, "masked_data", c.KafkaTopicMasked())
		assert.Equal(t, "pii_data", c.KafkaTopicPII())
		assert.Equal(t, "some-secret-key", c.PIIMaskerSecret())
	})

	t.Run("AppEnv", func(t *testing.T) {
		t.Run("Parses prod", func(t *testing.T) {
			stubValidEnvVars(t)
			t.Setenv("APP_ENV", "prod")

			c := config.NewConfig()

			assert.Equal(t, config.AppEnvProduction, c.AppEnv())
		})

		t.Run("Parses dev", func(t *testing.T) {
			stubValidEnvVars(t)
			t.Setenv("APP_ENV", "dev")

			c := config.NewConfig()

			assert.Equal(t, config.AppEnvDevelopment, c.AppEnv())
		})

		t.Run("Parses test", func(t *testing.T) {
			stubValidEnvVars(t)
			t.Setenv("APP_ENV", "test")

			c := config.NewConfig()

			assert.Equal(t, config.AppEnvTest, c.AppEnv())
		})
	})

	t.Run("KafkaBrokers", func(t *testing.T) {
		t.Run("Trims brokers", func(t *testing.T) {
			stubValidEnvVars(t)
			t.Setenv("KAFKA_BROKERS", " broker1 , broker2 , broker3 ")

			c := config.NewConfig()

			assert.Equal(t, []string{"broker1", "broker2", "broker3"}, c.KafkaBrokers())
		})

		t.Run("Panics if environment contains empty broker", func(t *testing.T) {
			stubValidEnvVars(t)
			t.Setenv("KAFKA_BROKERS", "broker1, ,broker3")

			assert.Panicsf(t, func() { config.NewConfig() }, "Required environment variable KAFKA_BROKERS cannot contain empty brokers")
		})
	})
}
