package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"prosamik-backend/internal/fetcher"
	"prosamik-backend/internal/parser"
	"prosamik-backend/pkg/models"
	"time"
)

// HandleReadmeRequest processes the README fetching and Markdown-to-HTML conversion
func HandleReadmeRequest(w http.ResponseWriter, r *http.Request) {
	// Extract query parameters
	owner := r.URL.Query().Get("owner")
	repo := r.URL.Query().Get("repo")

	if owner == "" || repo == "" {
		http.Error(w, "Owner and repository are required", http.StatusBadRequest)
		return
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// Fetch README content
	readmeContent, err := fetcher.FetchReadmeContent(ctx, owner, repo)
	if err != nil {
		http.Error(w, "Failed to fetch README", http.StatusInternalServerError)
		return
	}

	// Convert Markdown to HTML
	renderedHTML, err := parser.ConvertMarkdownToHTML(readmeContent)
	if err != nil {
		http.Error(w, "Failed to convert Markdown to HTML", http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := models.MarkdownDocument{
		Content:    renderedHTML,  // The HTML content processed from Markdown
		RawContent: readmeContent, // Original raw Markdown content
		Metadata: models.DocumentMetadata{
			Title:       "README - " + repo, // Dynamically set the title
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
