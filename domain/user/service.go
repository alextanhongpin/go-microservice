package user

type Service interface {
	HashPassword(password string) (string, error)
	ComparePassword(user Entity, password string) bool
}
