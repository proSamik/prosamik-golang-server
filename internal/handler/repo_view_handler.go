package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"prosamik-backend/internal/repository"
	"strconv"
	"time"
)

// RepoViewsResponse represents the structure for the repository views API response
type RepoViewsResponse struct {
	TotalViews int `json:"total_views"`
}

// TopRepositoriesResponse represents the structure for the top repositories API response
type TopRepositoriesResponse struct {
	Repositories []RepoViewItem `json:"repositories"`
}

// RepoViewItem represents a single repository with its view count
type RepoViewItem struct {
	RepoPath  string `json:"repo_path"`
	ViewCount int    `json:"view_count"`
}

// GetRepoViewsHandler returns the total views for a repository
func GetRepoViewsHandler(w http.ResponseWriter, r *http.Request) {
	repoPath := r.URL.Query().Get("repo_path")
	if repoPath == "" {
		http.Error(w, "repo_path parameter is missing", http.StatusBadRequest)
		return
	}

	repo := repository.NewRepoViewRepository()
	count, err := repo.GetTotalViewsForRepo(repoPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching view count: %v", err), http.StatusInternalServerError)
		return
	}

	response := RepoViewsResponse{
		TotalViews: count,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GetTopRepositoriesHandler returns the most viewed repositories
func GetTopRepositoriesHandler(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10 // default limit

	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
		limit = parsedLimit
	}

	repo := repository.NewRepoViewRepository()
	topRepos, err := repo.GetTopRepositories(limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching top repositories: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert to response format
	repoItems := make([]RepoViewItem, 0, len(topRepos))
	for _, item := range topRepos {
		repoItems = append(repoItems, RepoViewItem{
			RepoPath:  item.RepoPath,
			ViewCount: item.ViewCount,
		})
	}

	response := TopRepositoriesResponse{
		Repositories: repoItems,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GetRepoViewsByDateRangeHandler returns views for a repository within a date range
func GetRepoViewsByDateRangeHandler(w http.ResponseWriter, r *http.Request) {
	repoPath := r.URL.Query().Get("repo_path")
	if repoPath == "" {
		http.Error(w, "repo_path parameter is missing", http.StatusBadRequest)
		return
	}

	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var startDate, endDate time.Time
	var err error

	// Parse start date or use 30 days ago as default
	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			http.Error(w, "Invalid start_date format (use YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
	} else {
		startDate = time.Now().UTC().AddDate(0, 0, -30) // Default: 30 days ago
	}

	// Parse end date or use today as default
	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			http.Error(w, "Invalid end_date format (use YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
	} else {
		endDate = time.Now().UTC() // Default: today
	}

	repo := repository.NewRepoViewRepository()
	views, err := repo.GetViewsByDateRange(repoPath, startDate, endDate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching views by date range: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(views); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
