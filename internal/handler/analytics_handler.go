package handler

import (
	"net/http"
	"prosamik-backend/internal/repository"
	"strconv"
)

type AnalyticsHandlerInterface struct {
	projectRepo *repository.ProjectRepository
	blogRepo    *repository.BlogRepository
}

func HandleAnalytics(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get query parameters
	contentType := r.URL.Query().Get("type")
	idStr := r.URL.Query().Get("id")

	// Convert id to int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var updateErr error
	switch contentType {
	case "project":
		projectRepo := repository.NewProjectRepository()
		updateErr = projectRepo.IncrementProjectViewCount(id)
	case "blog":
		blogRepo := repository.NewBlogRepository()
		updateErr = blogRepo.IncrementBlogViewCount(id)
	default:
		http.Error(w, "Invalid content type", http.StatusBadRequest)
		return
	}

	if updateErr != nil {
		http.Error(w, updateErr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
