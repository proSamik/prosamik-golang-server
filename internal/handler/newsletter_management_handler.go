package handler

import (
	"log"
	"net/http"
	"prosamik-backend/internal/repository"
	"prosamik-backend/pkg/models"
	"strconv"
	"strings"
)

type NewsletterManagementData struct {
	Subscriptions []*models.Newsletter
}

func HandleNewsletterManagement(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	repo := repository.NewNewsletterRepository()
	subscriptions, err := repo.GetAllSubscriptions()
	if err != nil {
		log.Printf("Error fetching subscriptions: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Changed this part - now we're passing data to the base template
	data := PageData{
		Page: "newsletter", // This tells base.html which template to use
		Data: NewsletterManagementData{
			Subscriptions: subscriptions,
		},
	}

	// Changed to execute the base template instead
	err = templates.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

func HandleNewsletterDelete(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path
	path := r.URL.Path
	segments := strings.Split(path, "/")
	if len(segments) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	idStr := segments[len(segments)-1]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	repo := repository.NewNewsletterRepository()
	err = repo.DeleteSubscription(id)
	if err != nil {
		log.Printf("Error deleting subscription: %v", err)
		http.Error(w, "Failed to delete subscription", http.StatusInternalServerError)
		return
	}

	// Successfully deleted
	w.WriteHeader(http.StatusOK)
}
