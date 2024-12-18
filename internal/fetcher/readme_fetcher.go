package fetcher

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"prosamik-backend/internal/auth"
)

type ReadmeResponse struct {
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

// FetchReadmeContent retrieves README content from GitHub
func FetchReadmeContent(ctx context.Context, owner, repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/readme", owner, repo)
	log.Printf("Fetching README from: %s", url)

	req, err := createRequest(ctx, url)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned non-OK status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response body: %w", err)
	}

	var readmeResp ReadmeResponse
	if err := json.Unmarshal(body, &readmeResp); err != nil {
		return "", fmt.Errorf("parsing JSON response: %w", err)
	}

	if readmeResp.Content == "" {
		return "", fmt.Errorf("empty content received from GitHub API")
	}

	decodedContent, err := decodeReadmeContent(readmeResp.Content)
	if err != nil {
		return "", fmt.Errorf("decoding content: %w", err)
	}

	log.Printf("Successfully fetched and decoded README (length: %d characters)", len(decodedContent))
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

func decodeReadmeContent(content string) (string, error) {
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
