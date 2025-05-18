package gameserver

import (
	"strings"
)

func ToContainerName(s string) string {
	// Remove all non-alphanumeric characters and replace spaces with underscores
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, ".", "_")
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.ReplaceAll(s, "\\", "_")

	// Convert to lowercase
	return strings.ToLower(s)
}
