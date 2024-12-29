package repository

import (
	"database/sql"
	"fmt"
	"prosamik-backend/internal/database"
	"time"
)

type AnalyticsRepository struct {
	db *sql.DB
}

func NewAnalyticsRepository() *AnalyticsRepository {
	return &AnalyticsRepository{
		db: database.DB,
	}
}

func (r *AnalyticsRepository) IncrementPageViewCount(page string) error {
	today := time.Now().Format("2006-01-02")

	// Use separate queries for each page type to avoid SQL injection
	var query string
	switch page {
	case "home":
		query = `
            INSERT INTO analytics (date, home_views, updated_at)
            VALUES ($1, 1, CURRENT_TIMESTAMP)
            ON CONFLICT (date) 
            DO UPDATE SET 
                home_views = analytics.home_views + 1,
                updated_at = CURRENT_TIMESTAMP`
	case "about":
		query = `
            INSERT INTO analytics (date, about_views, updated_at)
            VALUES ($1, 1, CURRENT_TIMESTAMP)
            ON CONFLICT (date) 
            DO UPDATE SET 
                about_views = analytics.about_views + 1,
                updated_at = CURRENT_TIMESTAMP`
	case "blogs":
		query = `
            INSERT INTO analytics (date, blogs_views, updated_at)
            VALUES ($1, 1, CURRENT_TIMESTAMP)
            ON CONFLICT (date) 
            DO UPDATE SET 
                blogs_views = analytics.blogs_views + 1,
                updated_at = CURRENT_TIMESTAMP`
	case "projects":
		query = `
            INSERT INTO analytics (date, projects_views, updated_at)
            VALUES ($1, 1, CURRENT_TIMESTAMP)
            ON CONFLICT (date) 
            DO UPDATE SET 
                projects_views = analytics.projects_views + 1,
                updated_at = CURRENT_TIMESTAMP`
	case "feedback":
		query = `
            INSERT INTO analytics (date, feedback_views, updated_at)
            VALUES ($1, 1, CURRENT_TIMESTAMP)
            ON CONFLICT (date) 
            DO UPDATE SET 
                feedback_views = analytics.feedback_views + 1,
                updated_at = CURRENT_TIMESTAMP`
	default:
		return fmt.Errorf("invalid page type: %s", page)
	}

	// Prepare the statement
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("prepare statement error: %w", err)
	}

	// Use closure to handle defer error
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = fmt.Errorf("statement close error: %v: %w", closeErr, err)
		}
	}()

	_, err = stmt.Exec(today)
	if err != nil {
		return fmt.Errorf("increment view count error: %w", err)
	}

	return nil
}

func (r *AnalyticsRepository) GetAnalytics(startDate, endDate string) (map[string]map[string]int, error) {
	query := `
        SELECT date, home_views, about_views, blogs_views, projects_views, feedback_views 
        FROM analytics 
        WHERE date BETWEEN $1 AND $2
        ORDER BY date DESC
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

	rows, err := stmt.Query(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = fmt.Errorf("rows close error: %v: %w", err, err)
		}
	}(rows)

	stats := make(map[string]map[string]int)

	for rows.Next() {
		var date time.Time
		var home, about, blogs, projects, feedback int

		err := rows.Scan(&date, &home, &about, &blogs, &projects, &feedback)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		dateStr := date.Format("2006-01-02")
		stats[dateStr] = map[string]int{
			"home":     home,
			"about":    about,
			"blogs":    blogs,
			"projects": projects,
			"feedback": feedback,
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return stats, nil
}
