package gameserver

import (
	"crypto/rand"
	"fmt"
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
	s = strings.ToLower(s)

	// Generate a random suffix
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%s_%x", s, b)
}
