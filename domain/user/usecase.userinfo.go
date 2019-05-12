package user

// ? Change to root-level /me endpoint.
import (
	"errors"

	"github.com/alextanhongpin/go-microservice/pkg/gostrings"
)

type (
	userInfoRepository interface {
		WithID(id string) (User, error)
	}
	userInfoUseCase interface {
		UserInfo(id string) (*User, error)
	}
	UserInfoUseCase struct {
		users userInfoRepository
	}
)

func NewUserInfoUseCase(users userInfoRepository) *UserInfoUseCase {
	return &UserInfoUseCase{users}
}

func (u *UserInfoUseCase) UserInfo(id string) (*User, error) {
	if gostrings.IsEmpty(id) {
		return nil, errors.New("id is required")
	}
	user, err := u.users.WithID(id)
	return &user, err
}
