package token

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
)

type Service interface {
	CreateAccessToken(userID string) (string, error)
}

// Static methods.

// HashToken takes a plaintext token and returns an encrypted token.
func HashToken(plaintext string) string {
	h := sha256.New()
	h.Write([]byte(plaintext))
	b := h.Sum(nil)
	return hex.EncodeToString(b)
}

// ConstantTimeCompare takes two strings and compare them.
func ConstantTimeCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
