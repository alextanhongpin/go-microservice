package authn

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"time"

	"github.com/pkg/errors"
)

// Errors.
var (
	ErrInvalidRequest            = errors.New("invalid request")
	ErrInvalidPassword           = errors.New("invalid pssword")
	ErrInvalidUsernameOrPassword = errors.New("invalid username or password")
	ErrTokenExpired              = errors.New("token expired")
)

// User represents the user entity.
type User struct {
	ID             string `json:"id"`
	Email          string `json:"email,omitempty" validate:"omitempty,email"`
	HashedPassword string `json:"-"`
}

type Token struct {
	Token     string    `json:"-"`
	CreatedAt time.Time `json:"-"`
	UserID    string    `json:"-"`
}

func (t *Token) HasExpired(ttl time.Duration) bool {
	return time.Since(t.CreatedAt) > ttl
}

// Static methods.
func constantTimeCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

func hashToken(plaintext string) string {
	h := sha256.New()
	h.Write([]byte(plaintext))
	b := h.Sum(nil)
	return hex.EncodeToString(b)
}
