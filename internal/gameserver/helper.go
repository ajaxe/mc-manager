package gameserver

import (
	"regexp"
	"strings"
)

func ToContainerName(s string) string {
	p := regexp.MustCompile("(?i)[^a-z0-9_]")
	// Remove all non-alphanumeric characters and replace spaces with underscores
	v:= p.ReplaceAllString(s, "_")

	// Convert to lowercase
	return strings.ToLower(v)
}
