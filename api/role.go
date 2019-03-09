package api

// Role represents the actor of the system. One user can have only one Role.
type Role string

func (r Role) String() string {
	return string(r)
}

const (
	RoleUser  = Role("user")
	RoleGuest = Role("guest")
	RoleAdmin = Role("admin")
)
