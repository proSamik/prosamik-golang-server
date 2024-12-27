package auth

import (
	"os"
	"strings"
)

// GetGitHubToken retrieves the GitHub token
func GetGitHubToken() string {
	token := os.Getenv("GITHUB_TOKEN")
	return strings.TrimSpace(token) // Removes extra spaces and newlines
}
