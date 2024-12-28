package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"prosamik-backend/internal/database"
	"prosamik-backend/pkg/models"
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

	// Prepare the statement
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("prepare statement error: %w", err)
	}

	// Use closure to handle defer error
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = fmt.Errorf("statement close error: %v: %w", closeErr, err)
		}
	}()

	subscription := &models.Newsletter{}
	err = stmt.QueryRow(email).Scan(
		&subscription.ID,
		&subscription.Email,
		&subscription.RegistrationTime,
		&subscription.Verified,
	)

	// Return nil, nil if no subscription found
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
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

	// Prepare the statement
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("prepare statement error: %w", err)
	}

	// Use closure to handle defer error
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = fmt.Errorf("statement close error: %v: %w", closeErr, err)
		}
	}()

	subscription := &models.Newsletter{}
	err = stmt.QueryRow(email).Scan(
		&subscription.ID,
		&subscription.Email,
		&subscription.RegistrationTime,
		&subscription.Verified,
	)

	if err != nil {
		return nil, fmt.Errorf("create subscription error: %w", err)
	}

	return subscription, nil
}
