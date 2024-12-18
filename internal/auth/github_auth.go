package auth

import (
	"errors"
	"os"
)

// ValidateGitHubToken checks if GitHub token is set
func ValidateGitHubToken() error {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return errors.New("GITHUB_TOKEN environment variable is not set")
	}
	return nil
}

// GetGitHubToken retrieves the GitHub token
func GetGitHubToken() string {
	return os.Getenv("GITHUB_TOKEN")
}
