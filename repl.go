package main

import (
	"strings"
)

// cleanInput trims whitespace, lowercases the input, and splits it into words.
func cleanInput(text string) []string {
	trimmed := strings.TrimSpace(text)
	lowered := strings.ToLower(trimmed)
	// strings.Fields splits the string on any whitespace.
	words := strings.Fields(lowered)
	return words
}
