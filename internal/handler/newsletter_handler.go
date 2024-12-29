package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"prosamik-backend/internal/repository"
	"prosamik-backend/pkg/models"
	"strconv"
	"strings"
)

// NewsletterRequest represents the incoming request structure
type NewsletterRequest struct {
	Email string `json:"email"`
}

type NewsletterManagementData struct {
	Subscriptions []*models.Newsletter
}

// HandleNewsletter processes newsletter subscriptions
func HandleNewsletter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req NewsletterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	repo := repository.NewNewsletterRepository()

	// Check if email exists
	existing, err := repo.GetSubscriptionByEmail(req.Email)
	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := models.NewsletterResponse{
		Success: true,
	}

	if existing != nil {
		response.Message = "Already subscribed to newsletter"
	} else {
		// Create new subscription
		_, err := repo.CreateSubscription(req.Email)
		if err != nil {
			log.Printf("Error creating subscription: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		response.Message = "Successfully subscribed to newsletter"
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
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
