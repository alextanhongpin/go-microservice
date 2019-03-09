package model

type Role string

func (r Role) String() string {
	return string(r)
}

const RoleUser = Role("user")
const RoleGuest = Role("guest")
const RoleAdmin = Role("admin")
