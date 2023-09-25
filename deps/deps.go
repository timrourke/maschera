package deps

import (
	"github.com/segmentio/kafka-go"
	"github.com/timrourke/maschera/m/v2/app"
	"github.com/timrourke/maschera/m/v2/config"
	"github.com/timrourke/maschera/m/v2/log"
	"github.com/timrourke/maschera/m/v2/pii"
	"go.uber.org/zap"
)

const tenKB = 10e6

type deps struct {
	cfg            config.Config
	kafkaPIIReader *kafka.Reader
	logger         log.Logger
	piiMasker      pii.Masker
}

type Deps interface {
	Config() config.Config
	Logger() log.Logger
	KafkaPIIReader() *kafka.Reader
}

func (c *deps) Config() config.Config {
	if c.cfg != nil {
		return c.cfg
	}

	c.cfg = config.NewConfig()

	return c.cfg
}

func (c *deps) Logger() log.Logger {
	if c.logger != nil {
		return c.logger
	}

	var logFunc func(options ...zap.Option) (*zap.Logger, error)
	switch c.Config().AppEnv() {
	case config.AppEnvDevelopment:
		logFunc = zap.NewDevelopment
		break
	case config.AppEnvProduction:
		logFunc = zap.NewProduction
		break
	case config.AppEnvTest:
		logFunc = func(_ ...zap.Option) (*zap.Logger, error) { return zap.NewNop(), nil }
		break
	}

	zapLogger, err := logFunc()
	if err != nil {
		panic(err)
	}

	c.logger = log.NewLogger(zapLogger)

	return c.logger
}

func (c *deps) KafkaPIIReader() *kafka.Reader {
	if c.kafkaPIIReader != nil {
		return c.kafkaPIIReader
	}

	c.kafkaPIIReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  c.Config().KafkaBrokers(),
		Topic:    c.Config().KafkaTopicPII(),
		MaxBytes: tenKB,
	})

	return c.kafkaPIIReader
}

func (d *deps) PIIMasker() pii.Masker {
	if d.piiMasker != nil {
		return d.piiMasker
	}

	d.piiMasker = pii.NewMasker(d.Logger(), d.KafkaPIIReader())

	return d.piiMasker
}

func BuildApp() app.App {
	d := &deps{}

	return app.NewApp(d.Logger(), d.PIIMasker())
}
