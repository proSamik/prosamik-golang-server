package handler

import (
	"encoding/json"
	"net/http"
	"prosamik-backend/internal/data"
	"prosamik-backend/pkg/models"
)

func HandleBlogsList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	repos := make([]models.RepoListItem, 0, len(data.OrderedReposList))

	for i := len(data.OrderedReposList) - 1; i >= 0; i-- {
		item := data.OrderedReposList[i]
		repos = append(repos, models.RepoListItem{
			Title:       item.Title,
			RepoPath:    item.Info.Path,
			Description: item.Info.Description,
			Tags:        item.Info.Tags,
			ViewsCount:  item.Info.ViewsCount,
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
