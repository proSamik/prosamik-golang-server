package handler

import (
	"log"
	"net/http"
	"prosamik-backend/internal/repository"
	"prosamik-backend/pkg/models"
	"strconv"
	"strings"
	"time"
)

type NewsletterManagementData struct {
	Subscriptions []*models.Newsletter
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
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

func HandleNewsletterSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("search")

	repo := repository.NewNewsletterRepository()
	subscriptions, err := repo.SearchSubscriptions(query)
	if err != nil {
		log.Printf("Error searching subscriptions: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Pass subscriptions directly to the template
	err = templates.ExecuteTemplate(w, "newsletter-table", struct {
		Subscriptions []*models.Newsletter
	}{
		Subscriptions: subscriptions,
	})
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

func HandleNewsletterAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := normalizeEmail(r.FormValue("email"))
	if email == "" {
		err := templates.ExecuteTemplate(w, "form-message", struct{ Error string }{
			Error: "Email is required",
		})
		if err != nil {
			log.Printf("Template error: %v", err)
		}
		return
	}

	repo := repository.NewNewsletterRepository()

	// First check if email already exists (case insensitive)
	existing, err := repo.GetSubscriptionByEmail(email)
	if err != nil {
		log.Printf("Error checking existing email: %v", err)
		err = templates.ExecuteTemplate(w, "form-message", struct{ Error string }{
			Error: "Internal server error",
		})
		if err != nil {
			log.Printf("Template error: %v", err)
		}
		return
	}

	if existing != nil {
		err = templates.ExecuteTemplate(w, "form-message", struct{ Error string }{
			Error: "This email is already subscribed",
		})
		if err != nil {
			log.Printf("Template error: %v", err)
		}
		return
	}

	newsletter := &models.Newsletter{
		Email:            email, // Store normalized email
		RegistrationTime: time.Now(),
		Verified:         false,
	}

	err = repo.AddSubscription(newsletter)
	if err != nil {
		log.Printf("Error adding subscription: %v", err)
		err = templates.ExecuteTemplate(w, "form-message", struct{ Error string }{
			Error: "Failed to add subscription",
		})
		if err != nil {
			log.Printf("Template error: %v", err)
		}
		return
	}

	err = templates.ExecuteTemplate(w, "form-message", struct{ Error string }{})
	if err != nil {
		log.Printf("Template error: %v", err)
	}
}

func HandleNewsletterUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	email := normalizeEmail(r.FormValue("email"))
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	repo := repository.NewNewsletterRepository()
	err = repo.UpdateSubscription(id, email)
	if err != nil {
		log.Printf("Error updating subscription: %v", err)
		http.Error(w, "Failed to update subscription", http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(email))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func HandleNewsletterEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL
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

	// Get the subscription
	repo := repository.NewNewsletterRepository()
	subscription, err := repo.GetSubscription(id)
	if err != nil {
		log.Printf("Error fetching subscription: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "email-edit", subscription)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

func HandleNewsletterCancelEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
	subscription, err := repo.GetSubscription(id)
	if err != nil {
		log.Printf("Error fetching subscription: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(subscription.Email))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
