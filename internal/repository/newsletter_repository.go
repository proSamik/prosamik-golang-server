package repository

import (
	"database/sql"
	"fmt"
	"prosamik-backend/internal/database"
	"prosamik-backend/internal/models"
)

// NewsletterRepository handles database operations for newsletter subscriptions
type NewsletterRepository struct {
	db *sql.DB
}

// NewNewsletterRepository creates a new repository instance
func NewNewsletterRepository() *NewsletterRepository {
	return &NewsletterRepository{
		db: database.DB,
	}
}

// GetSubscriptionByEmail retrieves a subscription by email
func (r *NewsletterRepository) GetSubscriptionByEmail(email string) (*models.Newsletter, error) {
	query := `
        SELECT id, email, registration_timestamp, verified
        FROM newsletter_subscriptions
        WHERE email = $1
    `

	subscription := &models.Newsletter{}
	err := r.db.QueryRow(query, email).Scan(
		&subscription.ID,
		&subscription.Email,
		&subscription.RegistrationTime,
		&subscription.Verified,
	)

	// Return nil, nil if no subscription found
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("error getting subscription: %v", err)
	}

	return subscription, nil
}

// CreateSubscription adds a new newsletter subscription
func (r *NewsletterRepository) CreateSubscription(email string) (*models.Newsletter, error) {
	query := `
        INSERT INTO newsletter_subscriptions (email)
        VALUES ($1)
        RETURNING id, email, registration_timestamp, verified
    `

	subscription := &models.Newsletter{}
	err := r.db.QueryRow(query, email).Scan(
		&subscription.ID,
		&subscription.Email,
		&subscription.RegistrationTime,
		&subscription.Verified,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating subscription: %v", err)
	}

	return subscription, nil
}
