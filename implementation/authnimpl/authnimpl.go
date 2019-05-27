package authnimpl

import (
	"github.com/pkg/errors"
)

var (
	ErrTokenExpired              = errors.New("token expired")
	ErrDuplicateUser             = errors.New("user already exists")
	ErrInvalidRequest            = errors.New("invalid request")
	ErrInvalidPassword           = errors.New("invalid password")
	ErrInvalidUsernameOrPassword = errors.New("invalid username or password")
)
