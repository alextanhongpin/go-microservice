package usersvc

import "errors"

type UseCase struct {
	UserInfo UserInfoUseCase
}

type (
	UserInfoUseCase           func(id string) (User, error)
	UserInfoUseCaseRepository interface {
		WithID(id string) (User, error)
	}
)

func NewUserInfoUseCase(users UserInfoUseCaseRepository) UserInfoUseCase {
	return func(id string) (u User, err error) {
		if len(id) == 0 {
			err = errors.New("id is required")
			return
		}
		u, err = users.WithID(id)
		return
	}
}

type (
	GetUsersUseCase           func() ([]User, error)
	GetUsersUseCaseRepository interface {
		BelongingToPage() ([]User, error)
	}
)

func NewGetUsersUseCase(
	users GetUsersUseCaseRepository,
) GetUsersUseCase {
	return func() ([]User, error) {
		return users.BelongingToPage()
	}
}
