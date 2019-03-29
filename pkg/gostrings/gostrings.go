package gostrings

import "strings"

func IsEmpty(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}
