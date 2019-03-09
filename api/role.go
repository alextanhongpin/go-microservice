package api

// Role represents the actor of the system. One user can have only one Role.
type Role string

func (r Role) String() string {
	return string(r)
}

const (
	RoleUser  = Role("user")
	RoleOwner = Role("owner")
	RoleGuest = Role("guest")
	RoleAdmin = Role("admin")
	// RestrictedUser: Possibly for external client.
	// ReadOnly
)
