package pii_test

import (
	"context"
	"io"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/timrourke/maschera/m/v2/hasher"
	"github.com/timrourke/maschera/m/v2/log"
	"github.com/timrourke/maschera/m/v2/pii"
	"go.uber.org/zap"
)

var (
	shaHasher         = hasher.NewSha256("some-super-secret-value")
	jsonFieldsWithPII = []string{"email", "firstName", "lastName"}
	nilLogger         = log.NewLogger(zap.NewNop())
)

type stubKafkaMaskedWriter struct {
	MessagesWritten []string
}

func (s *stubKafkaMaskedWriter) WriteMessages(_ context.Context, msgs ...kafka.Message) error {
	for _, msg := range msgs {
		s.MessagesWritten = append(s.MessagesWritten, string(msg.Value))
	}

	return nil
}

func (s *stubKafkaMaskedWriter) Close() error {
	return nil
}

type stubKafkaPIIReader struct {
	MessagesToRead []kafka.Message
}

func (s *stubKafkaPIIReader) ReadMessage(_ context.Context) (kafka.Message, error) {
	if len(s.MessagesToRead) == 0 {
		return kafka.Message{}, io.EOF
	}

	msg, rest := s.MessagesToRead[0], s.MessagesToRead[1:]
	s.MessagesToRead = rest

	return msg, nil
}

func (s *stubKafkaPIIReader) Close() error {
	return nil
}

func buildMasker(msgsToRead []string) (pii.Masker, *stubKafkaPIIReader, *stubKafkaMaskedWriter) {
	var kafkaMsgsToRead []kafka.Message
	for _, msgStr := range msgsToRead {
		kafkaMsgsToRead = append(kafkaMsgsToRead, kafka.Message{Value: []byte(msgStr)})
	}

	kafkaPIIReader := &stubKafkaPIIReader{MessagesToRead: kafkaMsgsToRead}
	kafkaMaskedWriter := &stubKafkaMaskedWriter{}

	masker := pii.NewMasker(
		shaHasher,
		jsonFieldsWithPII,
		kafkaMaskedWriter,
		kafkaPIIReader,
		nilLogger,
	)

	return masker, kafkaPIIReader, kafkaMaskedWriter
}

func TestMasker(t *testing.T) {
	t.Run("Mask", func(t *testing.T) {
		expectedMessages := []string{
			"{\"email\":\"someone@example.com\",\"firstName\":\"Kelly\",\"lastName\":\"Baskerson\"}",
		}

		masker, _, stubWriter := buildMasker(expectedMessages)

		err := masker.Mask(context.Background())

		assert.Nil(t, err)
		assert.True(t, len(stubWriter.MessagesWritten) == 1)
		assert.Equal(
			t,
			"{\"email\":\"upwFtySM0mBjNxUdZAAy6D7LXMC22idHTwXp_HWuVHs=\",\"firstName\":\"upwFtySM0mBjNxUdZAAy6D7LXMC22idHTwXp_HWuVHs=\",\"lastName\":\"upwFtySM0mBjNxUdZAAy6D7LXMC22idHTwXp_HWuVHs=\"}",
			stubWriter.MessagesWritten[0],
		)
	})
}
