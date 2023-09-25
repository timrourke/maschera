package pii

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/timrourke/maschera/m/v2/log"
	"go.uber.org/zap"
)

type Masker interface {
	Mask(ctx context.Context) error
	Shutdown() error
}

type masker struct {
	logger         log.Logger
	kafkaPIIReader *kafka.Reader
}

func (m *masker) Mask(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	for {
		msg, err := m.kafkaPIIReader.ReadMessage(ctx)
		if err != nil {
			m.logger.Error("Error reading message from Kafka", zap.Error(err))
			continue
		}

		m.logger.Info("Received message: " + string(msg.Value))

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

func (m *masker) Shutdown() error {
	return m.kafkaPIIReader.Close()
}

func NewMasker(logger log.Logger, kafkaPIIReader *kafka.Reader) Masker {
	return &masker{
		logger:         logger,
		kafkaPIIReader: kafkaPIIReader,
	}
}
