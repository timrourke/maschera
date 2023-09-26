package hasher

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

type HmacSha256 interface {
	Sign(message []byte) ([]byte, error)
}

type hmacSha256 struct {
	secret []byte
}

func (s *hmacSha256) Sign(message []byte) ([]byte, error) {
	h := hmac.New(sha256.New, s.secret)
	_, err := h.Write(message)
	if err != nil {
		return nil, err
	}

	signature := h.Sum(nil)
	dst := make([]byte, base64.URLEncoding.EncodedLen(len(signature)))

	base64.URLEncoding.Encode(dst, signature)
	return dst, nil
}

func NewSha256(secret string) HmacSha256 {
	return &hmacSha256{
		secret: []byte(secret),
	}
}
