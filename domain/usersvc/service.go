package usersvc

type (
	getUsersUseCase interface {
		GetUsers() ([]User, error)
	}
	userInfoUseCase interface {
		UserInfo(id string) (*User, error)
	}

	// Service groups the usecases together.
	Service struct {
		getUsersUseCase
		userInfoUseCase
	}

	repository interface {
		WithID(id string) (User, error)
		BelongingToPage() ([]User, error)
	}
)

func NewService(repo repository) *Service {
	return &Service{
		getUsersUseCase: NewGetUsersUseCase(repo),
		userInfoUseCase: NewUserInfoUseCase(repo),
	}
}
