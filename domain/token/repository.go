package token

import "time"

type Reader interface {
	WithValue(token string) (Entity, error)
}
type Writer interface {
	Create(userID, token string) (bool, error)
	Delete(token string) (bool, error)
	DeleteExpired(ttl time.Duration) (int64, error)
}

type Repository interface {
	Reader
	Writer
}
