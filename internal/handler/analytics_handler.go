package handler

import (
	"fmt"
	"net/http"
	"prosamik-backend/internal/repository"
	"strconv"
	"strings"
)

type AnalyticsHandlerInterface struct {
	projectRepo   *repository.ProjectRepository
	blogRepo      *repository.BlogRepository
	analyticsRepo *repository.AnalyticsRepository
}

func NewAnalyticsHandler() *AnalyticsHandlerInterface {
	return &AnalyticsHandlerInterface{
		projectRepo:   repository.NewProjectRepository(),
		blogRepo:      repository.NewBlogRepository(),
		analyticsRepo: repository.NewAnalyticsRepository(),
	}
}

// HandleAnalytics is the main handler function that routes to appropriate sub-handlers
func HandleAnalytics(w http.ResponseWriter, r *http.Request) {
	handler := NewAnalyticsHandler()

	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check which type of analytics to handle
	if page := r.URL.Query().Get("page"); page != "" {
		handler.handlePageAnalytics(w, page) // Removed the request parameter
		return
	}

	handler.handleContentAnalytics(w, r)
}

// handlePageAnalytics handles the page view analytics
func (h *AnalyticsHandlerInterface) handlePageAnalytics(w http.ResponseWriter, page string) {
	if !h.isValidPage(page) {
		http.Error(w, "Invalid page parameter", http.StatusBadRequest)
		return
	}

	if err := h.analyticsRepo.IncrementPageViewCount(page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// handleContentAnalytics handles the project and blog analytics
func (h *AnalyticsHandlerInterface) handleContentAnalytics(w http.ResponseWriter, r *http.Request) {
	contentType := r.URL.Query().Get("type")
	idStr := r.URL.Query().Get("id")

	id, err := h.parseID(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	if err := h.updateContentViewCount(contentType, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Helper functions
func (h *AnalyticsHandlerInterface) isValidPage(page string) bool {
	validPages := map[string]bool{
		"home":     true,
		"about":    true,
		"blogs":    true,
		"projects": true,
		"feedback": true,
	}
	return validPages[strings.ToLower(page)]
}

func (h *AnalyticsHandlerInterface) parseID(idStr string) (int64, error) {
	return strconv.ParseInt(idStr, 10, 64)
}

func (h *AnalyticsHandlerInterface) updateContentViewCount(contentType string, id int64) error {
	switch contentType {
	case "project":
		return h.projectRepo.IncrementProjectViewCount(id)
	case "blog":
		return h.blogRepo.IncrementBlogViewCount(id)
	default:
		return fmt.Errorf("invalid content type")
	}
}
