package mystrings

import "strings"

func ReplaceAll(s, old, new string) string {
	return strings.Replace(s, old, new, -1)
}
