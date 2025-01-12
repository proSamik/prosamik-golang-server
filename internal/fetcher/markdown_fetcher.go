package fetcher

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"prosamik-backend/internal/auth"
	"prosamik-backend/internal/cache"
	"strings"
	"time"
)

// GitHubFile represents the structure for file content response
type GitHubFile struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Content string `json:"content"`
}

// GitHubCommit represents a single commit in the commit API response
type GitHubCommit struct {
	Commit struct {
		Committer struct {
			Date string `json:"date"`
		} `json:"committer"`
	} `json:"commit"`
}

// FetchContentFromGitHubURL fetches file content from GitHub API with caching
func FetchContentFromGitHubURL(ctx context.Context, apiURL string) (string, error) {
	// Try to get from cache first
	cached, err := cache.GetCachedContent(ctx, apiURL)
	if err != nil {
		fmt.Printf("Cache error: %v, falling back to GitHub API\n", err)
	} else if cached != nil {
		// If found in the cache, verify if content is still fresh by checking the last commit
		commitsURL := getCommitsURL(apiURL)
		lastUpdated, err := FetchLastCommitData(ctx, commitsURL)
		if err != nil {
			fmt.Printf("Error checking last commit, using cached content: %v\n", err)
			return cached.Content, nil
		}

		// If content hasn't changed, return a cached version
		if !lastUpdated.After(cached.LastUpdated) {
			return cached.Content, nil
		}
		fmt.Printf("Cache outdated, fetching fresh content\n")
	}

	// Fetch fresh content from GitHub
	content, err := fetchFreshContent(ctx, apiURL)
	if err != nil {
		return "", err
	}

	// Get last commit time for caching
	lastUpdated, err := FetchLastCommitData(ctx, getCommitsURL(apiURL))
	if err != nil {
		fmt.Printf("Warning: couldn't get last commit time: %v\n", err)
		return content, nil
	}

	// Cache the new content
	if err := cache.SetCachedContent(ctx, apiURL, &cache.CachedContent{
		Content:     content,
		LastUpdated: lastUpdated,
	}); err != nil {
		fmt.Printf("Warning: failed to cache content: %v\n", err)
	}

	return content, nil
}

// fetchFreshContent contains the original content fetching logic
func fetchFreshContent(ctx context.Context, apiURL string) (string, error) {
	body, err := makeGitHubRequest(ctx, apiURL)
	if err != nil {
		return "", err
	}

	var fileContent GitHubFile
	if err := json.Unmarshal(body, &fileContent); err != nil {
		return "", fmt.Errorf("error unmarshalling GitHub file response: %v", err)
	}

	decodedContent, err := decodeBase64Content(fileContent.Content)
	if err != nil {
		return "", fmt.Errorf("error decoding base64 content: %v", err)
	}

	return decodedContent, nil
}

func getCommitsURL(contentURL string) string {
	parts := strings.Split(contentURL, "/contents/")
	if len(parts) != 2 {
		return contentURL
	}
	return parts[0] + "/commits?path=" + parts[1]
}

// FetchLastCommitData fetches the last commit information for a file or repository
func FetchLastCommitData(ctx context.Context, apiURL string) (time.Time, error) {
	body, err := makeGitHubRequest(ctx, apiURL)
	if err != nil {
		return time.Time{}, err
	}

	var commits []GitHubCommit
	if err := json.Unmarshal(body, &commits); err != nil {
		fmt.Printf("Raw GitHub commits response: %s\n", string(body))
		return time.Time{}, fmt.Errorf("error unmarshalling GitHub commits response: %v", err)
	}

	if len(commits) == 0 {
		return time.Time{}, fmt.Errorf("no commits found")
	}

	lastUpdated, err := time.Parse(time.RFC3339, commits[0].Commit.Committer.Date)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing commit date: %v", err)
	}

	return lastUpdated, nil
}

// makeGitHubRequest makes a generic HTTP request to GitHub API
func makeGitHubRequest(ctx context.Context, url string) ([]byte, error) {
	req, err := createRequest(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			fmt.Printf("warning: failed to close response body: %v\n", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned non-OK status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	return body, nil
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
