package user

type (
	// UseCase groups the usecases together.
	UseCase struct {
		getUsersUseCase
		userInfoUseCase
	}
	repository interface {
		WithID(id string) (User, error)
		BelongingToPage() ([]User, error)
	}
)

func NewUseCase(repo repository) *UseCase {
	return &UseCase{
		getUsersUseCase: NewGetUsersUseCase(repo),
		userInfoUseCase: NewUserInfoUseCase(repo),
	}
}
