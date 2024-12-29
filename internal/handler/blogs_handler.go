package handler

import (
	"encoding/json"
	"net/http"
	"prosamik-backend/internal/repository"
	"prosamik-backend/pkg/models"
)

func HandleBlogsList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Initialize repository
	blogRepo := repository.NewBlogRepository()

	// Fetch blogs using repository
	blogs, err := blogRepo.GetAllBlogs()
	if err != nil {
		http.Error(w, "Failed to fetch blogs", http.StatusInternalServerError)
		return
	}

	// Convert blogs to RepoListItems format
	repos := make([]models.RepoListItem, 0, len(blogs))
	for _, blog := range blogs {
		repos = append(repos, models.RepoListItem{
			Title:       blog.Title,
			RepoPath:    blog.Path,
			Description: blog.Description,
			Tags:        blog.Tags,
			ViewsCount:  blog.ViewsCount,
			ID:          int(blog.ID),
			Type:        "blog",
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
