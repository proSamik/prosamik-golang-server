package fetcher

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"prosamik-backend/internal/auth"
	"strings"
	"time"
)

type GitHubFile struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Content string `json:"content"`
}

// FetchContentFromGitHubURL fetches content from a GitHub URL using the constructed API URL
func FetchContentFromGitHubURL(ctx context.Context, apiURL string) (string, error) {

	req, err := createRequest(ctx, apiURL)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("making request: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			fmt.Printf("warning: failed to close response body: %v\n", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned non-OK status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response body: %w", err)
	}

	// Parse the JSON response into GitHubFile struct
	var fileContent GitHubFile
	if err := json.Unmarshal(body, &fileContent); err != nil {
		return "", fmt.Errorf("error unmarshalling GitHub response: %v", err)
	}

	// If the content is base64 encoded, decode it
	decodedContent, err := decodeBase64Content(fileContent.Content)
	if err != nil {
		return "", fmt.Errorf("error decoding base64 content: %v", err)
	}

	return decodedContent, nil
}

func createRequest(ctx context.Context, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	ghToken := auth.GetGitHubToken()
	if ghToken == "" {
		return nil, fmt.Errorf("GitHub token is empty")
	}

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", "Bearer "+ghToken)
	return req, nil
}

// decodeBase64Content decodes the base64 content from GitHub API response
func decodeBase64Content(content string) (string, error) {
	cleanContent := strings.ReplaceAll(content, "\n", "")
	decodedBytes, err := base64.StdEncoding.DecodeString(cleanContent)
	if err != nil {
		return "", fmt.Errorf("base64 decoding failed: %w", err)
	}

	if len(decodedBytes) == 0 {
		return "", fmt.Errorf("decoded content is empty")
	}

	return string(decodedBytes), nil
}
