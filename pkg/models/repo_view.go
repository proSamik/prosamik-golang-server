package models

import "time"

// RepoView represents a GitHub repository view tracking record
type RepoView struct {
	ID        int64     `json:"id"`
	RepoPath  string    `json:"repo_path"`  // Full GitHub URL path
	ViewDate  time.Time `json:"view_date"`  // Date of the view
	ViewCount int       `json:"view_count"` // Number of views on this date
	CreatedAt time.Time `json:"created_at"` // Record creation time
	UpdatedAt time.Time `json:"updated_at"` // Record last update time
}
