package handler

import (
	"encoding/json"
	"net/http"
	"prosamik-backend/internal/data"
	"prosamik-backend/pkg/models"
)

func HandleReposList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	repos := make([]models.RepoListItem, 0, len(data.ReposList))

	for title, repoInfo := range data.ReposList {
		repos = append(repos, models.RepoListItem{
			Title:       title,
			RepoPath:    repoInfo.Path,
			Description: repoInfo.Description,
		})
	}

	response := models.RepoListResponse{
		Repos: repos,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
