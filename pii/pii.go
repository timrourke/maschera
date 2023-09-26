package pii

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/hashicorp/go-multierror"
	"github.com/segmentio/kafka-go"
	"github.com/timrourke/maschera/m/v2/hasher"
	"github.com/timrourke/maschera/m/v2/log"
	"go.uber.org/zap"
)

type KafkaPIIReader interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	Close() error
}

type KafkaMaskedWriter interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type Masker interface {
	Mask(ctx context.Context) error
	Shutdown() error
}

type masker struct {
	hasher            hasher.HmacSha256
	jsonFieldsWithPII []string
	kafkaMaskedWriter KafkaMaskedWriter
	kafkaPIIReader    KafkaPIIReader
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
			if errors.Is(err, io.EOF) {
				m.logger.Debug("Kafka PII reader closed")
				return nil
			}

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

		maskedJSON, err := m.maskPIIFields(msg)
		if err != nil {
			m.logger.Error("Error masking PII fields", zap.Error(err))
			continue
		}

		err = m.writeMaskedJSON(ctx, maskedJSON)
		switch err.(type) {
		case nil:
			m.logger.Debug("Successfully wrote message")
			continue
		case kafka.WriteErrors:
			m.logger.Error("Error writing message", zap.Error(err))
			return err
		default:
			panic(err)
		}
	}
}

func (m *masker) maskPIIFields(msg kafka.Message) ([]byte, error) {
	var jsonValue map[string]interface{}
	err := json.Unmarshal(msg.Value, &jsonValue)
	if err != nil {
		m.logger.Error("Error deserializing message to JSON", zap.Error(err))
		return nil, err
	}

	for _, field := range m.jsonFieldsWithPII {
		if _, ok := jsonValue[field]; ok {
			hash, err := m.hasher.Sign(msg.Value)
			if err != nil {
				m.logger.Error("Error signing message", zap.Error(err))
				return nil, err
			}

			jsonValue[field] = string(hash)
		}
	}

	return json.Marshal(jsonValue)
}

func (m *masker) writeMaskedJSON(ctx context.Context, maskedJSON []byte) error {
	maskedMessage := kafka.Message{Value: maskedJSON}

	return m.kafkaMaskedWriter.WriteMessages(ctx, maskedMessage)
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
	jsonFieldsWithPII []string,
	kafkaMaskedWriter KafkaMaskedWriter,
	kafkaPIIReader KafkaPIIReader,
	logger log.Logger,
) Masker {
	return &masker{
		hasher:            hasher,
		jsonFieldsWithPII: jsonFieldsWithPII,
		kafkaMaskedWriter: kafkaMaskedWriter,
		kafkaPIIReader:    kafkaPIIReader,
		logger:            logger,
	}
}
