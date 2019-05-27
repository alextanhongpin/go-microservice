package user

type Reader interface {
	WithEmail(email string) (Entity, error)
	WithID(userID string) (Entity, error)
}

type Writer interface {
	Create(username, password string) (Entity, error)
	ChangePassword(userID, password string) (bool, error)
}

type Repository interface {
	Reader
	Writer
}
