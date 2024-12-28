package handler

import (
	"encoding/json"
	"net/http"
	"prosamik-backend/internal/data"
	"prosamik-backend/pkg/models"
)

func HandleProjectsList(w http.ResponseWriter, r *http.Request) {
	// 1. Check if the method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Initialize a slice to store project items
	projects := make([]models.RepoListItem, 0, len(data.OrderedReposList))

	// 3. Iterate through the data in reverse order (newest first)
	for i := len(data.OrderedReposList) - 1; i >= 0; i-- {
		item := data.OrderedReposList[i]
		projects = append(projects, models.RepoListItem{
			Title:       item.Title,
			RepoPath:    item.Info.Path,
			Description: item.Info.Description,
			Tags:        item.Info.Tags,
			ViewsCount:  item.Info.ViewsCount,
		})
	}

	// 4. Create the response structure
	response := models.RepoListResponse{
		Repos: projects,
	}

	// 5. Set header and encode response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
