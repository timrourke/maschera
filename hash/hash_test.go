package hash_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/timrourke/maschera/m/v2/hash"
)

const (
	expectedHmacSha256 = "hlEbYWLslrxz6P6xRWfrgzRLnFaUYq3Tb4R-5aZ5Rn0="
	message            = "some-message"
	secret             = "some-secure-secret"
)

func TestHmacSha256_Sign(t *testing.T) {
	result, err := hash.NewSha256(secret).Sign([]byte(message))

	assert.Nil(t, err)
	assert.Equal(t, expectedHmacSha256, string(result))
}
