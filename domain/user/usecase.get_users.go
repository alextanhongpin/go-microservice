package user

type (
	getUsersRepository interface {
		BelongingToPage() ([]User, error)
	}
	getUsersUseCase interface {
		GetUsers() ([]User, error)
	}
	GetUsersUseCase struct {
		users getUsersRepository
	}
)

func NewGetUsersUseCase(users getUsersRepository) *GetUsersUseCase {
	return &GetUsersUseCase{users}
}

func (g *GetUsersUseCase) GetUsers() ([]User, error) {
	return g.users.BelongingToPage()
}
