package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"prosamik-backend/internal/cache"
	"prosamik-backend/internal/fetcher"
	"prosamik-backend/internal/parser"
	"prosamik-backend/pkg/models"
	"regexp"
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

// constructGitHubAPIURL parses a GitHub URL and returns API URL and repository details
func constructGitHubAPIURL(githubURL string) (string, string, string, string, string, error) {
	repoPrefix := "https://github.com/"
	if !strings.HasPrefix(githubURL, repoPrefix) {
		return "", "", "", "", "", fmt.Errorf("URL must start with %s", repoPrefix)
	}

	repoPath := strings.TrimPrefix(githubURL, repoPrefix)
	parts := strings.Split(repoPath, "/")

	if len(parts) < 2 {
		return "", "", "", "", "", fmt.Errorf("invalid GitHub URL format: %s", githubURL)
	}

	owner := parts[0]
	repo := parts[1]
	var filePath string
	branchName := "main" // default branch

	if strings.Contains(githubURL, "/blob/") {
		branchName = parts[3]
		filePath = strings.Join(parts[4:], "/")
		apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s",
			owner, repo, filePath, branchName)
		return apiURL, owner, repo, filePath, branchName, nil
	} else if strings.Contains(githubURL, "/tree/") {
		branchName = parts[3]
		filePath = strings.Join(parts[4:], "/") + "/README.md"
		apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s",
			owner, repo, filePath, branchName)
		return apiURL, owner, repo, filePath, branchName, nil
	} else {
		filePath = "README.md"
		apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/README.md",
			owner, repo)
		return apiURL, owner, repo, filePath, branchName, nil
	}
}

// processImageURLs converts relative image URLs to raw.githubusercontent.com URLs
func processImageURLs(content, owner, repo, branch, markdownPath string) string {
	markdownDir := filepath.Dir(markdownPath)

	// Handle Markdown image syntax ![alt](path)
	// Using simpler pattern that matches any path not starting with http:// or https://
	mdPattern := regexp.MustCompile(`!\[(.*?)\]\(((?:\./|[^)h]|h[^t]|ht[^t]|htt[^p]|http[^:/]|https[^:/])[^)]*)\)`)
	content = mdPattern.ReplaceAllStringFunc(content, func(match string) string {
		parts := mdPattern.FindStringSubmatch(match)
		if len(parts) < 3 {
			return match
		}

		altText := parts[1]
		imagePath := parts[2]

		// Skip URLs that somehow matched our pattern
		if strings.HasPrefix(imagePath, "http://") || strings.HasPrefix(imagePath, "https://") {
			return match
		}

		// If path starts with ./, remove it and join with markdownDir
		// Otherwise, treat it as relative to markdown directory
		var fullPath string
		if strings.HasPrefix(imagePath, "./") {
			relPath := strings.TrimPrefix(imagePath, "./")
			fullPath = filepath.Join(markdownDir, relPath)
		} else {
			// For paths not starting with ./, treat them relative to markdown directory
			fullPath = filepath.Join(markdownDir, imagePath)
		}

		fullPath = filepath.ToSlash(fullPath)

		// If alt text is empty, use the last part of the path
		if altText == "" {
			pathParts := strings.Split(fullPath, "/")
			if len(pathParts) > 0 {
				fileName := pathParts[len(pathParts)-1]
				altText = strings.TrimSuffix(fileName, filepath.Ext(fileName))
			}
		}

		rawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s",
			owner, repo, branch, fullPath)

		return fmt.Sprintf("![%s](%s)", altText, rawURL)
	})

	// Handle HTML image syntax <img src="path" />
	htmlPattern := regexp.MustCompile(`<img[^>]+src=["']((?:\./|[^"'h]|h[^t]|ht[^t]|htt[^p]|http[^:/]|https[^:/])[^"']*)["']`)
	content = htmlPattern.ReplaceAllStringFunc(content, func(match string) string {
		parts := htmlPattern.FindStringSubmatch(match)
		if len(parts) < 2 {
			return match
		}

		imagePath := parts[1]

		// Skip URLs that somehow matched our pattern
		if strings.HasPrefix(imagePath, "http://") || strings.HasPrefix(imagePath, "https://") {
			return match
		}

		// If path starts with ./, remove it and join with markdownDir
		// Otherwise, treat it as relative to markdown directory
		var fullPath string
		if strings.HasPrefix(imagePath, "./") {
			relPath := strings.TrimPrefix(imagePath, "./")
			fullPath = filepath.Join(markdownDir, relPath)
		} else {
			// For paths not starting with ./, treat them relative to markdown directory
			fullPath = filepath.Join(markdownDir, imagePath)
		}

		fullPath = filepath.ToSlash(fullPath)

		rawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s",
			owner, repo, branch, fullPath)

		return strings.Replace(match, parts[1], rawURL, 1)
	})

	return content
}

// getFileName gets filename of the Markdown file
func getFileName(filePath string) string {
	parts := strings.Split(filePath, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

// Fix relative path resolution in image URLs
func processImageURL(url string) string {
	// Remove any ../ from the URL path
	return strings.ReplaceAll(url, "/../", "/")
}

// MarkdownHandler processes GitHub markdown content and returns rendered HTML
func MarkdownHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}

	// Try to get from cache first
	cached, err := cache.GetCachedContent(r.Context(), url)
	if err == nil && cached != nil {
		// Unmarshal the cached response
		var response models.MarkdownDocument
		if err := json.Unmarshal([]byte(cached.Content), &response); err != nil {
			fmt.Printf("Warning: failed to unmarshal cached response: %v\n", err)
			// Continue with normal processing since cache read failed
		} else {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, "Failed to encode cached response", http.StatusInternalServerError)
				return
			}
			return
		}
	}

	// If not in cache or error, proceed with normal processing
	apiURL, owner, repo, filePath, branch, err := constructGitHubAPIURL(url)
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

	// Process image URLs before converting to HTML
	processedContent := processImageURLs(markdownContent, owner, repo, branch, filePath)

	// Construct commits API URL and fetch last updated time
	commitsURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits?path=%s&sha=%s&page=1&per_page=1",
		owner, repo, filePath, branch)
	lastUpdated, err := fetcher.FetchLastCommitData(r.Context(), commitsURL)
	if err != nil && branch == "main" {
		// If the main branch fails, try with "master" branch
		masterCommitsURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits?path=%s&sha=%s&page=1&per_page=1",
			owner, repo, filePath, "master")
		lastUpdated, err = fetcher.FetchLastCommitData(r.Context(), masterCommitsURL)
		if err != nil {
			http.Error(w, "Failed to fetch document metadata from both main and master branches", http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		http.Error(w, "Failed to fetch document metadata", http.StatusInternalServerError)
		return
	}

	renderedHTML, err := parser.ConvertMarkdownToHTML(processedContent)
	if err != nil {
		http.Error(w, "Failed to convert Markdown to HTML", http.StatusInternalServerError)
		return
	}

	// Get the title based on URL type
	title := repo // default title
	if strings.Contains(url, "/blob/") || strings.Contains(url, "/tree/") {
		title = getFileName(filePath)
	}

	// Get description from content if available
	description := "This is the README for the repository." // default description
	if len(markdownContent) > 0 {
		// Take the first 200 characters, trim to last complete word
		if len(markdownContent) > 100 {
			description = markdownContent[:100]
			lastSpace := strings.LastIndex(description, " ")
			if lastSpace > 0 {
				description = description[:lastSpace] + "..."
			}
		} else {
			description = markdownContent
		}
	}

	response := models.MarkdownDocument{
		Content: renderedHTML,
		//RawContent: markdownContent,
		Metadata: models.DocumentMetadata{
			Title:       title,
			Repository:  repo,
			LastUpdated: lastUpdated,
			Author:      owner,
			Description: description,
		},
	}

	// Cache the response before sending
	responseBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("Warning: failed to marshal response for caching: %v\n", err)
	} else {
		// Store in cache
		if err := cache.SetCachedContent(r.Context(), url, &cache.CachedContent{
			Content:     string(responseBytes),
			LastUpdated: lastUpdated,
		}); err != nil {
			fmt.Printf("Warning: failed to cache response: %v\n", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
