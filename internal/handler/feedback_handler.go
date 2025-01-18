package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

// Feedback represents the structure of the feedback form data
type Feedback struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

// HandleFeedback processes feedback submissions and sends them via SMTP
func HandleFeedback(w http.ResponseWriter, r *http.Request) {
	// Check for the correct HTTP method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the site parameter from the URL
	site := r.URL.Query().Get("site")

	// Decode the JSON body to feedback struct
	var feedback Feedback
	if err := json.NewDecoder(r.Body).Decode(&feedback); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Different validation logic based on site
	if site == "githubme" {
		// For githubme site, only email and message are required
		if feedback.Email == "" || feedback.Message == "" {
			http.Error(w, "Email and message are required", http.StatusBadRequest)
			return
		}
	} else {
		// For default site, all fields are required
		if feedback.Name == "" || feedback.Email == "" || feedback.Message == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}
	}

	// Send the feedback email with site information
	if err := sendFeedbackEmail(feedback, site); err != nil {
		http.Error(w, "Failed to send feedback", http.StatusInternalServerError)
		return
	}

	// Send a success response if everything is fine
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"success": "true",
		"message": "Feedback submitted successfully",
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// sendFeedbackEmail sends the feedback email using SMTP
func sendFeedbackEmail(feedback Feedback, site string) error {
	// Retrieve SMTP configuration from environment variables
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	recipient := os.Getenv("FEEDBACK_RECIPIENT_EMAIL")

	// Check if any of the environment variables are empty
	if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPassword == "" || recipient == "" {
		return fmt.Errorf("one or more required environment variables are missing")
	}

	// Set subject based on site
	subject := "Feedback on prosamik.com"
	if site == "githubme" {
		subject = "Feedback on githubme.com"
	}

	// Construct the email headers and body
	headers := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n",
		smtpUser, recipient, subject,
	)

	// Handle optional name field
	nameField := "N/A"
	if feedback.Name != "" {
		nameField = feedback.Name
	}

	body := fmt.Sprintf(
		"Name: %s\nEmail: %s\nMessage:\n%s",
		nameField, feedback.Email, feedback.Message,
	)

	// The Rest of the email sending logic remains the same
	emailContent := headers + body
	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", smtpHost, smtpPort),
		auth,
		smtpUser,
		[]string{recipient},
		[]byte(emailContent),
	)
	if err != nil {
		log.Printf("Failed to send email: %v\n", err)
		return err
	}

	log.Println("Feedback email sent successfully")
	return nil
}
