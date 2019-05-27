package api

import "strings"

// Scope represents the Resource the Actor can access. A Resource can have
// multiples scopes, e.g. Resource for Events can have "read:events",
// "delete:events" etc.
type Scope string

func (s Scope) String() string {
	return string(s)
}

func (s Scope) Equal(scope string) bool {
	return string(s) == scope
}

const (
	ScopeDefault = Scope("default")
	ScopeOpenID  = Scope("scope")
	ScopeProfile = Scope("profile")
)

// Scopes can be separated by white spaces. e.g. "openid profile".
func Scopes(scopes ...Scope) string {
	result := make([]string, len(scopes))
	for i, scope := range scopes {
		result[i] = scope.String()
	}
	return strings.Join(result, " ")
}
