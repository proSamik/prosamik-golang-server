package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"prosamik-backend/internal/fetcher"
	"prosamik-backend/internal/parser"
	"prosamik-backend/pkg/models"
	"strings"
)

// GitHubCommit represents a single commit in GitHub's API response
type GitHubCommit struct {
	Sha    string `json:"sha"`
	Commit struct {
		Author struct {
			Name  string `json:"name"`
			Email string `json:"email"`
			Date  string `json:"date"`
		} `json:"author"`
		Committer struct {
			Name  string `json:"name"`
			Email string `json:"email"`
			Date  string `json:"date"`
		} `json:"committer"`
		Message string `json:"message"`
	} `json:"commit"`
}

func constructGitHubAPIURL(githubURL string) (string, string, string, string, error) {
	// Previous URL parsing logic remains the same...
	// But now we also return the file path for getting the last updated time

	repoPrefix := "https://github.com/"
	if !strings.HasPrefix(githubURL, repoPrefix) {
		return "", "", "", "", fmt.Errorf("URL must start with %s", repoPrefix)
	}

	repoPath := strings.TrimPrefix(githubURL, repoPrefix)
	parts := strings.Split(repoPath, "/")

	if len(parts) < 2 {
		return "", "", "", "", fmt.Errorf("invalid GitHub URL format: %s", githubURL)
	}

	owner := parts[0]
	repo := parts[1]
	var filePath string

	if strings.Contains(githubURL, "/blob/") {
		filePath = strings.Join(parts[4:], "/")
		branchName := parts[3]
		apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s",
			owner, repo, filePath, branchName)
		return apiURL, owner, repo, filePath, nil
	} else if strings.Contains(githubURL, "/tree/") {
		branchName := parts[3]
		filePath = strings.Join(parts[4:], "/") + "/README.md"
		apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s",
			owner, repo, filePath, branchName)
		return apiURL, owner, repo, filePath, nil
	} else {
		filePath = "README.md"
		apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/README.md", owner, repo)
		return apiURL, owner, repo, filePath, nil
	}
}

func MarkdownHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}

	// Now getting filePath as well from the URL constructor
	apiURL, owner, repo, filePath, err := constructGitHubAPIURL(url)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error constructing GitHub API URL: %v", err),
			http.StatusBadRequest)
		return
	}

	markdownContent, err := fetcher.FetchContentFromGitHubURL(r.Context(), apiURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching content: %v", err),
			http.StatusInternalServerError)
		return
	}

	// Construct commits API URL and fetch last updated time
	commitsURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits?path=%s&page=1&per_page=1",
		owner, repo, filePath)
	lastUpdated, err := fetcher.FetchLastCommitData(r.Context(), commitsURL)
	if err != nil {
		http.Error(w, "Failed to fetch document metadata", http.StatusInternalServerError)
		return
	}

	renderedHTML, err := parser.ConvertMarkdownToHTML(markdownContent)
	if err != nil {
		http.Error(w, "Failed to convert Markdown to HTML", http.StatusInternalServerError)
		return
	}

	response := models.MarkdownDocument{
		Content:    renderedHTML,
		RawContent: markdownContent,
		Metadata: models.DocumentMetadata{
			Title:       repo,
			Repository:  repo,
			LastUpdated: lastUpdated, // Now using the actual last updated time
			Author:      owner,
			Description: "This is the README for the repository.",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
