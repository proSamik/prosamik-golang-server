package models

import "time"

// MarkdownDocument represents the response model for the README content
type MarkdownDocument struct {
	Content    string           `json:"content"`    // HTML content converted from Markdown
	RawContent string           `json:"rawContent"` // Original raw Markdown content
	Metadata   DocumentMetadata `json:"metadata"`   // Metadata about the document
}

// DocumentMetadata holds the metadata for the document (e.g., title, repository)
type DocumentMetadata struct {
	Title       string    `json:"title"`       // Title of the document (e.g., "README - repo")
	Repository  string    `json:"repository"`  // Repository name
	LastUpdated time.Time `json:"lastUpdated"` // Timestamp of the last update
	Author      string    `json:"author"`      // Author of the repository (owner)
	Description string    `json:"description"` // Description or summary of the document
}

type RepoListItem struct {
	Title       string `json:"title"`
	RepoPath    string `json:"repoPath"`
	Description string `json:"description"`
}

type RepoListResponse struct {
	Repos []RepoListItem `json:"repos"`
}
