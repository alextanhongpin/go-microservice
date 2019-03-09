package api

import "github.com/alextanhongpin/pkg/set"

type Roles map[Role]set.Set

func (r Roles) Can(target string) []Role {
	var result []Role
	for role, scopes := range r {
		if scopes.Has(target) {
			result = append(result, role)
		}
	}
	return result
}
