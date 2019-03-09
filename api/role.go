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

type Roles map[Role][]string

func (r Roles) Can(target string) (result []Role) {
	for role, scopes := range r {
		for _, scope := range scopes {
			if scope == target {
				result = append(result, role)
				break
			}
		}
	}
	return
}
