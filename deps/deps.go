package deps

import (
	"github.com/segmentio/kafka-go"
	"github.com/timrourke/maschera/m/v2/app"
	"github.com/timrourke/maschera/m/v2/config"
	"github.com/timrourke/maschera/m/v2/hasher"
	"github.com/timrourke/maschera/m/v2/log"
	"github.com/timrourke/maschera/m/v2/pii"
	"go.uber.org/zap"
)

const tenKB = 10e6

type deps struct {
	cfg               config.Config
	kafkaPIIReader    *kafka.Reader
	kafkaMaskedWriter *kafka.Writer
	logger            log.Logger
	piiMasker         pii.Masker
	hasher            hasher.HmacSha256
}

type Deps interface {
	Config() config.Config
	Hasher() hasher.HmacSha256
	KafkaPIIReader() *kafka.Reader
	KafkaMaskedWriter() *kafka.Writer
	Logger() log.Logger
}

func (d *deps) Config() config.Config {
	if d.cfg != nil {
		return d.cfg
	}

	d.cfg = config.NewConfig()

	return d.cfg
}

func (d *deps) Hasher() hasher.HmacSha256 {
	if d.hasher != nil {
		return d.hasher
	}

	d.hasher = hasher.NewSha256(d.Config().PIIMaskerSecret())

	return d.hasher
}

func (d *deps) Logger() log.Logger {
	if d.logger != nil {
		return d.logger
	}

	var logFunc func(options ...zap.Option) (*zap.Logger, error)
	switch d.Config().AppEnv() {
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

	d.logger = log.NewLogger(zapLogger)

	return d.logger
}

func (d *deps) KafkaPIIReader() *kafka.Reader {
	if d.kafkaPIIReader != nil {
		return d.kafkaPIIReader
	}

	d.kafkaPIIReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  d.Config().KafkaBrokers(),
		Topic:    d.Config().KafkaTopicPII(),
		MaxBytes: tenKB,
	})

	return d.kafkaPIIReader
}

func (d *deps) KafkaMaskedWriter() *kafka.Writer {
	if d.kafkaMaskedWriter != nil {
		return d.kafkaMaskedWriter
	}

	d.kafkaMaskedWriter = &kafka.Writer{
		Addr:         kafka.TCP(d.Config().KafkaBrokers()...),
		Topic:        d.Config().KafkaTopicMasked(),
		RequiredAcks: kafka.RequireAll,
	}

	return d.kafkaMaskedWriter
}

func (d *deps) PIIMasker() pii.Masker {
	if d.piiMasker != nil {
		return d.piiMasker
	}

	d.piiMasker = pii.NewMasker(
		d.Hasher(),
		d.KafkaMaskedWriter(),
		d.KafkaPIIReader(),
		d.Logger(),
	)

	return d.piiMasker
}

func BuildApp() app.App {
	d := &deps{}

	return app.NewApp(d.Logger(), d.PIIMasker())
}
