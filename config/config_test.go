package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/timrourke/maschera/m/v2/config"
)

func stubValidEnvVars(t *testing.T) {
	t.Setenv("APP_ENV", "prod")
	t.Setenv("KAFKA_BROKERS", "broker1,broker2,broker3")
	t.Setenv("KAFKA_TOPIC_PII", "pii_data")
	t.Setenv("PII_MASKER_SECRET", "some-secret-key")
}

func TestConfig(t *testing.T) {
	t.Run("AppEnv", func(t *testing.T) {
		t.Run("Returns app env from environment variable", func(t *testing.T) {
			stubValidEnvVars(t)

			c := config.NewConfig()

			assert.Equal(t, config.AppEnvProduction, c.AppEnv())
		})

		t.Run("Panics if environment variable is missing", func(t *testing.T) {
			stubValidEnvVars(t)
			t.Setenv("APP_ENV", "")

			assert.Panicsf(t, func() { config.NewConfig() }, "Missing required environment variable APP_ENV")
		})

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

		t.Run("Parses dev", func(t *testing.T) {
			stubValidEnvVars(t)
			t.Setenv("APP_ENV", "test")

			c := config.NewConfig()

			assert.Equal(t, config.AppEnvTest, c.AppEnv())
		})
	})

	t.Run("KafkaBrokers", func(t *testing.T) {
		t.Run("Returns brokers from environment variable", func(t *testing.T) {
			stubValidEnvVars(t)

			c := config.NewConfig()

			assert.Equal(t, []string{"broker1", "broker2", "broker3"}, c.KafkaBrokers())
		})

		t.Run("Trims brokers", func(t *testing.T) {
			stubValidEnvVars(t)
			t.Setenv("KAFKA_BROKERS", " broker1 , broker2 , broker3 ")

			c := config.NewConfig()

			assert.Equal(t, []string{"broker1", "broker2", "broker3"}, c.KafkaBrokers())
		})

		t.Run("Panics if environment variable is missing", func(t *testing.T) {
			stubValidEnvVars(t)
			t.Setenv("KAFKA_BROKERS", "")

			assert.Panicsf(t, func() { config.NewConfig() }, "Missing required environment variable KAFKA_BROKERS")
		})

		t.Run("Panics if environment variable is empty", func(t *testing.T) {
			stubValidEnvVars(t)
			t.Setenv("KAFKA_BROKERS", " ")

			assert.Panicsf(t, func() { config.NewConfig() }, "Required environment variable KAFKA_BROKERS cannot contain empty brokers")
		})

		t.Run("Panics if environment contains empty broker", func(t *testing.T) {
			stubValidEnvVars(t)
			t.Setenv("KAFKA_BROKERS", "broker1, ,broker3")

			assert.Panicsf(t, func() { config.NewConfig() }, "Required environment variable KAFKA_BROKERS cannot contain empty brokers")
		})
	})

	t.Run("KafkaTopicPII", func(t *testing.T) {
		t.Run("Returns topic from environment variable", func(t *testing.T) {
			stubValidEnvVars(t)

			c := config.NewConfig()

			assert.Equal(t, "pii_data", c.KafkaTopicPII())
		})

		t.Run("Trims topic", func(t *testing.T) {
			stubValidEnvVars(t)
			t.Setenv("KAFKA_TOPIC_PII", " pii_data ")

			c := config.NewConfig()

			assert.Equal(t, "pii_data", c.KafkaTopicPII())
		})

		t.Run("Panics if environment variable is missing", func(t *testing.T) {
			stubValidEnvVars(t)
			t.Setenv("KAFKA_TOPIC_PII", "")

			assert.Panicsf(t, func() { config.NewConfig() }, "Missing required environment variable KAFKA_TOPIC_PII")
		})

		t.Run("Panics if environment variable is empty", func(t *testing.T) {
			stubValidEnvVars(t)
			t.Setenv("KAFKA_TOPIC_PII", "   ")

			assert.Panicsf(t, func() { config.NewConfig() }, "Required environment variable KAFKA_TOPIC_PII cannot be empty")
		})
	})

	t.Run("PIIMaskerSecret", func(t *testing.T) {
		t.Run("Returns secret from environment variable", func(t *testing.T) {
			stubValidEnvVars(t)

			c := config.NewConfig()

			assert.Equal(t, "some-secret-key", c.PIIMaskerSecret())
		})

		t.Run("Trims secret", func(t *testing.T) {
			stubValidEnvVars(t)
			t.Setenv("PII_MASKER_SECRET", " some-secret-key    ")

			c := config.NewConfig()

			assert.Equal(t, "some-secret-key", c.PIIMaskerSecret())
		})

		t.Run("Panics if environment variable is missing", func(t *testing.T) {
			stubValidEnvVars(t)
			t.Setenv("PII_MASKER_SECRET", "")

			assert.Panicsf(t, func() { config.NewConfig() }, "Missing required environment variable PII_MASKER_SECRET")
		})

		t.Run("Panics if environment variable is empty", func(t *testing.T) {
			stubValidEnvVars(t)
			t.Setenv("PII_MASKER_SECRET", "   ")

			assert.Panicsf(t, func() { config.NewConfig() }, "Required environment variable PII_MASKER_SECRET cannot be empty")
		})
	})
}
