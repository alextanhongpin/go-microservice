package authsvc

import (
	"github.com/pkg/errors"

	"github.com/alextanhongpin/passwd"

	"github.com/alextanhongpin/go-microservice/model"
	"github.com/alextanhongpin/go-microservice/pkg/signer"
)

type (
	Service interface {
		Login(username, password string) (model.User, error)
		Register(username, password string) (model.User, error)
	}
	ServiceImpl struct {
		signer signer.Signer
		repo   Repository
	}
)

func (s *ServiceImpl) Login(username, password string) (*model.User, error) {
	user, err := s.repo.GetUser(username)
	if err != nil {
		return nil, errors.Wrap(err, "query user failed")
	}
	err = passwd.Verify(password, user.HashedPassword)
	return user, errors.Wrap(err, "verify password failed")

}
func (s *ServiceImpl) Register(username, password string) (*model.User, error) {
	hashedPassword, err := passwd.Hash(password)
	if err != nil {
		return nil, errors.Wrap(err, "hash password failed")
	}
	user, err := s.repo.CreateUser(username, hashedPassword)
	return user, errors.Wrap(err, "create user failed")
}
