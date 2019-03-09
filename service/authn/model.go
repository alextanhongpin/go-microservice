package authn

// User represents the user entity.
type User struct {
	ID             string `json:"-"`
	Email          string `json:"-"`
	HashedPassword string `json:"-"`
}
