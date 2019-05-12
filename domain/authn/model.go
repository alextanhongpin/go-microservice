package authn

// User represents the user entity.
type User struct {
	ID             string `json:"id"`
	Email          string `json:"email,omitempty" validate:"omitempty,email"`
	HashedPassword string `json:"-"`
}
