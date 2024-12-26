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
	// Check for correct HTTP method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the JSON body to feedback struct
	var feedback Feedback
	if err := json.NewDecoder(r.Body).Decode(&feedback); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Check if required fields are empty
	if feedback.Name == "" || feedback.Email == "" || feedback.Message == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Send the feedback email
	if err := sendFeedbackEmail(feedback); err != nil {
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
func sendFeedbackEmail(feedback Feedback) error {
	// Retrieve SMTP configuration from environment variables
	smtpHost := os.Getenv("SMTP_HOST")                 // e.g., "smtp.gmail.com"
	smtpPort := os.Getenv("SMTP_PORT")                 // e.g., "587"
	smtpUser := os.Getenv("SMTP_USER")                 // Your SMTP username/email
	smtpPassword := os.Getenv("SMTP_PASSWORD")         // Your SMTP password/app password
	recipient := os.Getenv("FEEDBACK_RECIPIENT_EMAIL") // Feedback recipient email

	// Check if any of the environment variables are empty
	if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPassword == "" || recipient == "" {
		return fmt.Errorf("one or more required environment variables are missing")
	}

	// Construct the email headers and body
	subject := "Feedback on prosamik.com"
	headers := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n",
		smtpUser, recipient, subject,
	)
	body := fmt.Sprintf(
		"Name: %s\nEmail: %s\nMessage:\n%s",
		feedback.Name, feedback.Email, feedback.Message,
	)

	// Combine headers and body
	emailContent := headers + body

	// Set up authentication
	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)

	// Send the email
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
