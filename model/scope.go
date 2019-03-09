package model

import "strings"

type Scope string

func (s Scope) String() string {
	return string(s)
}

const ScopeDefault = Scope("default")
const ScopeOpenID = Scope("scope")

// Scopes can be separated by white spaces. e.g. "openid profile".
func Scopes(scopes ...Scope) string {
	return strings.Join(scopes, " ")
}
