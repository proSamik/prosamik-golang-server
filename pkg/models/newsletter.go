package models

import "time"

// Newsletter represents the newsletter subscription data structure
type Newsletter struct {
	ID               int64     `json:"id"`
	Email            string    `json:"email"`
	RegistrationTime time.Time `json:"registration_time"`
	Verified         bool      `json:"verified"`
}

// NewsletterResponse represents the API response structure
type NewsletterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
