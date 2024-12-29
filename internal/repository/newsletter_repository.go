package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"prosamik-backend/internal/database"
	"prosamik-backend/pkg/models"
	"strings"
)

type NewsletterRepository struct {
	db *sql.DB
}

func NewNewsletterRepository() *NewsletterRepository {
	return &NewsletterRepository{
		db: database.DB,
	}
}

// normalizeEmail ensures consistent email format
func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func (r *NewsletterRepository) GetSubscriptionByEmail(email string) (*models.Newsletter, error) {
	query := `
        SELECT id, email, registration_timestamp, verified
        FROM newsletter_subscriptions
        WHERE LOWER(email) = LOWER($1)
    `

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("prepare statement error: %w", err)
	}

	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = fmt.Errorf("statement close error: %v: %w", closeErr, err)
		}
	}()

	subscription := &models.Newsletter{}
	err = stmt.QueryRow(normalizeEmail(email)).Scan(
		&subscription.ID,
		&subscription.Email,
		&subscription.RegistrationTime,
		&subscription.Verified,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}

	return subscription, nil
}

func (r *NewsletterRepository) CreateSubscription(email string) (*models.Newsletter, error) {
	query := `
        INSERT INTO newsletter_subscriptions (email)
        VALUES (LOWER($1))
        RETURNING id, email, registration_timestamp, verified
    `

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("prepare statement error: %w", err)
	}

	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = fmt.Errorf("statement close error: %v: %w", closeErr, err)
		}
	}()

	subscription := &models.Newsletter{}
	err = stmt.QueryRow(normalizeEmail(email)).Scan(
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

func (r *NewsletterRepository) AddSubscription(newsletter *models.Newsletter) error {
	query := `
        INSERT INTO newsletter_subscriptions (email, registration_timestamp, verified)
        VALUES (LOWER($1), $2, $3)
        RETURNING id
    `

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("prepare statement error: %w", err)
	}

	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = fmt.Errorf("statement close error: %v: %w", closeErr, err)
		}
	}()

	var id int64
	err = stmt.QueryRow(
		normalizeEmail(newsletter.Email),
		newsletter.RegistrationTime,
		newsletter.Verified,
	).Scan(&id)

	if err != nil {
		return fmt.Errorf("create subscription error: %w", err)
	}

	newsletter.ID = id
	return nil
}

func (r *NewsletterRepository) SearchSubscriptions(query string) ([]*models.Newsletter, error) {
	searchQuery := `
        SELECT id, email, registration_timestamp, verified
        FROM newsletter_subscriptions
        WHERE LOWER(email) LIKE LOWER($1)
        ORDER BY registration_timestamp ASC
    `

	stmt, err := r.db.Prepare(searchQuery)
	if err != nil {
		return nil, fmt.Errorf("prepare statement error: %w", err)
	}

	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = fmt.Errorf("statement close error: %v: %w", closeErr, err)
		}
	}()

	rows, err := stmt.Query("%" + strings.ToLower(query) + "%")
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = fmt.Errorf("rows close error: %v: %w", err, err)
		}
	}(rows)

	var subscriptions []*models.Newsletter
	for rows.Next() {
		subscription := &models.Newsletter{}
		err := rows.Scan(
			&subscription.ID,
			&subscription.Email,
			&subscription.RegistrationTime,
			&subscription.Verified,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, nil
}

func (r *NewsletterRepository) UpdateSubscription(id int64, email string) error {
	query := `
        UPDATE newsletter_subscriptions
        SET email = LOWER($1)
        WHERE id = $2
    `

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("prepare statement error: %w", err)
	}

	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = fmt.Errorf("statement close error: %v: %w", closeErr, err)
		}
	}()

	result, err := stmt.Exec(normalizeEmail(email), id)
	if err != nil {
		return fmt.Errorf("update error: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected error: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no subscription found with id: %d", id)
	}

	return nil
}

func (r *NewsletterRepository) GetAllSubscriptions() ([]*models.Newsletter, error) {
	query := `
        SELECT id, email, registration_timestamp, verified
        FROM newsletter_subscriptions
        ORDER BY registration_timestamp ASC
    `

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("prepare statement error: %w", err)
	}

	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = fmt.Errorf("statement close error: %v: %w", closeErr, err)
		}
	}()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = fmt.Errorf("rows close error: %v: %w", err, err)
		}
	}(rows)

	var subscriptions []*models.Newsletter
	for rows.Next() {
		subscription := &models.Newsletter{}
		err := rows.Scan(
			&subscription.ID,
			&subscription.Email,
			&subscription.RegistrationTime,
			&subscription.Verified,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, nil
}

func (r *NewsletterRepository) DeleteSubscription(id int64) error {
	query := `
        DELETE FROM newsletter_subscriptions
        WHERE id = $1
    `

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("prepare statement error: %w", err)
	}

	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = fmt.Errorf("statement close error: %v: %w", closeErr, err)
		}
	}()

	result, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("delete error: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected error: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no subscription found with id: %d", id)
	}

	return nil
}

// GetSubscription retrieves a single subscription by ID
func (r *NewsletterRepository) GetSubscription(id int64) (*models.Newsletter, error) {
	query := `
        SELECT id, email, registration_timestamp, verified
        FROM newsletter_subscriptions
        WHERE id = $1
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
	err = stmt.QueryRow(id).Scan(
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
