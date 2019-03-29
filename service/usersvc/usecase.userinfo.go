package usersvc

import "errors"

type (
	userInfoRepository interface {
		WithID(id string) (User, error)
	}
	UserInfoUseCase struct {
		users userInfoRepository
	}
)

func NewUserInfoUseCase(users userInfoRepository) *UserInfoUseCase {
	return &UserInfoUseCase{users}
}

func (u *UserInfoUseCase) UserInfo(id string) (*User, error) {
	if len(id) == 0 {
		return nil, errors.New("id id required")
	}
	user, err := u.users.WithID(id)
	return &user, err
}
