package authenticator

// User represents the user entity.
type User struct {
	ID             string `json:"-"`
	HashedPassword string `json:"-"`
}
