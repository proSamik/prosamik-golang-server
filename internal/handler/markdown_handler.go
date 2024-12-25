package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"prosamik-backend/internal/fetcher"
	"prosamik-backend/internal/parser"
	"prosamik-backend/pkg/models"
	"strings"
	"time"
)

// constructGitHubAPIURL constructs the GitHub API URL based on the GitHub URL
func constructGitHubAPIURL(githubURL string) (string, string, string, error) {

	// Step 1: Ensure the URL starts with the correct GitHub prefix
	repoPrefix := "https://github.com/"
	if !strings.HasPrefix(githubURL, repoPrefix) {
		return "", "", "", fmt.Errorf("URL must start with %s", repoPrefix)
	}

	// Strip the prefix and split to get the owner and repo
	repoPath := strings.TrimPrefix(githubURL, repoPrefix)
	parts := strings.Split(repoPath, "/")

	if len(parts) < 2 {
		return "", "", "", fmt.Errorf("invalid GitHub URL format: %s", githubURL)
	}

	owner := parts[0]
	repo := parts[1]

	// Step 2: Determine URL type and construct API URL
	if strings.Contains(githubURL, "/blob/") {
		// If it's a file (blob), construct the API URL with the file path
		filePath := strings.Join(parts[4:], "/") // This will give the path after "blob/branch_name/"
		branchName := parts[3]                   // This is the branch name

		apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s", owner, repo, filePath, branchName)
		return apiURL, owner, repo, nil
	} else if strings.Contains(githubURL, "/tree/") {
		// If it's a directory or branch reference, construct the API URL for the directory
		branchName := parts[3]                        // This is the branch name
		directoryPath := strings.Join(parts[4:], "/") // Directory path after "tree/branch_name/"

		// Default file name is "README.md" if directory is specified
		apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s/README.md?ref=%s", owner, repo, directoryPath, branchName)
		return apiURL, owner, repo, nil
	} else {
		// If it's just the repository URL, fetch the default README
		apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/README.md", owner, repo)
		return apiURL, owner, repo, nil
	}
}

func MarkdownHandler(w http.ResponseWriter, r *http.Request) {
	// Get GitHub URL from query parameter
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}

	// Construct the corresponding GitHub API URL
	apiURL, owner, repo, err := constructGitHubAPIURL(url)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error constructing GitHub API URL: %v", err), http.StatusBadRequest)
		return
	}

	// Fetch content from GitHub using the constructed API URL
	markdownContent, err := fetcher.FetchContentFromGitHubURL(r.Context(), apiURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching content: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert Markdown to HTML
	renderedHTML, err := parser.ConvertMarkdownToHTML(markdownContent)
	if err != nil {
		http.Error(w, "Failed to convert Markdown to HTML", http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := models.MarkdownDocument{
		Content:    renderedHTML,    // The HTML content processed from Markdown
		RawContent: markdownContent, // Original raw Markdown content
		Metadata: models.DocumentMetadata{
			Title:       repo,
			Repository:  repo,
			LastUpdated: time.Now(),                               // Set the current time for last updated
			Author:      owner,                                    // The repository owner (Author)
			Description: "This is the README for the repository.", // Description field
		},
	}

	// Send response with Content-Type as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
