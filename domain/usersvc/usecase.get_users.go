package usersvc

type (
	usersGetter interface {
		BelongingToPage() ([]User, error)
	}
	GetUsersUseCase struct {
		users usersGetter
	}
)

func NewGetUsersUseCase(users usersGetter) *GetUsersUseCase {
	return &GetUsersUseCase{users}
}

func (g *GetUsersUseCase) GetUsers() ([]User, error) {
	return g.users.BelongingToPage()
}
