package authnimpl

import (
	"github.com/alextanhongpin/go-microservice/domain/user"
	"github.com/alextanhongpin/passwd"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) HashPassword(password string) (string, error) {
	pwd, err := passwd.Hash(password)
	return pwd, err
}

func (s *Service) ComparePassword(user user.Entity, password string) bool {
	err := passwd.Verify(password, user.HashedPassword)
	return err == nil
}
