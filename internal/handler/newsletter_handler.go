package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"prosamik-backend/internal/repository"
	"prosamik-backend/pkg/models"
	"strings"
)

// NewsletterRequest represents the incoming request structure
type NewsletterRequest struct {
	Email string `json:"email"`
}

func normalizeEmailSignup(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// HandleNewsletterSignup processes newsletter subscriptions
func HandleNewsletterSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req NewsletterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Normalize email before processing
	req.Email = normalizeEmailSignup(req.Email)
	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	repo := repository.NewNewsletterRepository()

	// Check if email exists (will be case insensitive since we normalized the input)
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
		// Create a new subscription with normalized email
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
