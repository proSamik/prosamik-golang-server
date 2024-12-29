package handler

import (
	"encoding/json"
	"net/http"
	"prosamik-backend/internal/repository"
	"prosamik-backend/pkg/models"
)

func HandleProjectsList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Initialize repository
	projectRepo := repository.NewProjectRepository()

	// Fetch projects using repository
	projects, err := projectRepo.GetAllProjects()
	if err != nil {
		http.Error(w, "Failed to fetch projects", http.StatusInternalServerError)
		return
	}

	// Convert projects to RepoListItems format
	repos := make([]models.RepoListItem, 0, len(projects))
	for _, project := range projects {
		repos = append(repos, models.RepoListItem{
			Title:       project.Title,
			RepoPath:    project.Path, // Path maps to RepoPath
			Description: project.Description,
			Tags:        project.Tags,
			ViewsCount:  project.ViewsCount,
		})
	}

	// Create the response in required format
	response := models.RepoListResponse{
		Repos: repos,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
