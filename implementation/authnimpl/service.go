package authnimpl

import (
	"strings"

	"github.com/alextanhongpin/go-microservice/domain/user"
	"github.com/alextanhongpin/passwd"
	"github.com/pkg/errors"
)

var ErrPasswordRequired = errors.New("password is required")

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) HashPassword(password string) (string, error) {
	if len(strings.TrimSpace(password)) == 0 {
		return "", ErrPasswordRequired
	}
	pwd, err := passwd.Hash(password)
	return pwd, err
}

func (s *Service) ComparePassword(user user.Entity, password string) bool {
	err := passwd.Verify(password, user.HashedPassword)
	return err == nil
}
