package pii

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/segmentio/kafka-go"
	"github.com/timrourke/maschera/m/v2/hasher"
	"github.com/timrourke/maschera/m/v2/log"
	"go.uber.org/zap"
)

type Masker interface {
	Mask(ctx context.Context) error
	Shutdown() error
}

type masker struct {
	hasher            hasher.HmacSha256
	kafkaMaskedWriter *kafka.Writer
	kafkaPIIReader    *kafka.Reader
	logger            log.Logger
}

func (m *masker) Mask(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		msg, err := m.kafkaPIIReader.ReadMessage(ctx)
		if err != nil {
			m.logger.Error("Error reading message from Kafka", zap.Error(err))
			continue
		}

		m.logger.Debug(
			"Received message",
			zap.String("topic", msg.Topic),
			zap.Int("partition", msg.Partition),
			zap.Int64("offset", msg.Offset),
			zap.String("key", string(msg.Key)),
		)

		hash, err := m.hasher.Sign(msg.Value)
		if err != nil {
			m.logger.Error("Error signing message", zap.Error(err))
			continue
		}

		hashedMessage := kafka.Message{Value: hash}

		switch err := m.kafkaMaskedWriter.WriteMessages(ctx, hashedMessage).(type) {
		case nil:
			m.logger.Debug("Successfully wrote message")
			break
		case kafka.WriteErrors:
			m.logger.Error("Error writing message", zap.Error(err))
			break
		default:
			panic(err)
		}
	}
}

func (m *masker) Shutdown() error {
	var result *multierror.Error

	err := m.kafkaPIIReader.Close()
	if err != nil {
		result = multierror.Append(result, err)
	}

	err = m.kafkaMaskedWriter.Close()
	if err != nil {
		result = multierror.Append(result, err)
	}

	return result.ErrorOrNil()
}

func NewMasker(
	hasher hasher.HmacSha256,
	kafkaMaskedWriter *kafka.Writer,
	kafkaPIIReader *kafka.Reader,
	logger log.Logger,
) Masker {
	return &masker{
		hasher:            hasher,
		kafkaMaskedWriter: kafkaMaskedWriter,
		kafkaPIIReader:    kafkaPIIReader,
		logger:            logger,
	}
}
