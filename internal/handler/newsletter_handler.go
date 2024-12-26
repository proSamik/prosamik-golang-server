package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Newsletter represents the structure of the newsletter subscription data
type Newsletter struct {
	Email string `json:"email"`
}

// HandleNewsletter processes newsletter subscriptions and logs them
func HandleNewsletter(w http.ResponseWriter, r *http.Request) {
	// Check for correct HTTP method - This is important because we only want to handle POST requests
	// just like in your frontend useSubscribeNewsletter hook which uses POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the JSON body to newsletter struct - This matches the JSON being sent
	// from your frontend where you're sending { email }
	var newsletter Newsletter
	if err := json.NewDecoder(r.Body).Decode(&newsletter); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Check if email field is empty - Basic validation similar to your frontend
	if newsletter.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Log the subscription with timestamp - Since we don't have a DB, we'll just print to console
	log.Printf("[%s] New newsletter subscription: %s",
		time.Now().Format(time.RFC3339),
		newsletter.Email,
	)

	// Send a success response - This matches the SubscriptionResponse interface
	// in your frontend code
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"success": true,
		"message": "Successfully subscribed to newsletter",
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
