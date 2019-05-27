package token

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// TTL is the time-to-live for the generated token.
const TTL = 10 * time.Minute

// Entity represents the temporary token that is send when resetting password.
type Entity struct {
	Token     string    `json:"-"`
	CreatedAt time.Time `json:"-"`
	UserID    string    `json:"-"`
}

// HasExpired checks if the token has expired in the given duration.
func (e *Entity) HasExpired(ttl time.Duration) bool {
	return time.Since(e.CreatedAt) > ttl
}

// New returns a new Token entity with the given user id.
func New(userID string) Entity {
	return Entity{
		Token:     uuid.Must(uuid.NewV4()).String(),
		CreatedAt: time.Now(),
		UserID:    userID,
	}
}

// Hashed returns the encrypted token.
func (e *Entity) Hashed() string {
	return HashToken(e.Token)
}
