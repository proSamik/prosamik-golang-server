package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"prosamik-backend/internal/database"
	"prosamik-backend/pkg/models"
	"strings"
	"time"
)

// RepoViewRepository handles database operations for repository views
type RepoViewRepository struct {
	db *sql.DB
}

// NewRepoViewRepository creates a new RepoViewRepository instance
func NewRepoViewRepository() *RepoViewRepository {
	return &RepoViewRepository{
		db: database.DB,
	}
}

// normalizeRepoPath ensures consistent repo path format
func normalizeRepoPath(repoPath string) string {
	return strings.TrimSpace(repoPath)
}

// RecordView increments the view count for a repository path on the current date
// If no record exists for today, it creates a new one
func (r *RepoViewRepository) RecordView(repoPath string) error {
	normalizedPath := normalizeRepoPath(repoPath)
	today := time.Now().UTC().Truncate(24 * time.Hour)

	query := `
		INSERT INTO repo_views (repo_path, view_date, view_count, created_at, updated_at)
		VALUES ($1, $2, 1, NOW(), NOW())
		ON CONFLICT (repo_path, view_date)
		DO UPDATE SET
			view_count = repo_views.view_count + 1,
			updated_at = NOW()
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

	_, err = stmt.Exec(normalizedPath, today)
	if err != nil {
		return fmt.Errorf("record view error: %w", err)
	}

	return nil
}

// GetTotalViewsForRepo returns total view count for a specific repository path
func (r *RepoViewRepository) GetTotalViewsForRepo(repoPath string) (int, error) {
	normalizedPath := normalizeRepoPath(repoPath)

	query := `
		SELECT COALESCE(SUM(view_count), 0) as total_views
		FROM repo_views
		WHERE repo_path = $1
	`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("prepare statement error: %w", err)
	}

	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = fmt.Errorf("statement close error: %v: %w", closeErr, err)
		}
	}()

	var totalViews int
	err = stmt.QueryRow(normalizedPath).Scan(&totalViews)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("scan error: %w", err)
	}

	return totalViews, nil
}

// GetViewsByDateRange returns view data for a repo within a date range
func (r *RepoViewRepository) GetViewsByDateRange(repoPath string, startDate, endDate time.Time) ([]*models.RepoView, error) {
	normalizedPath := normalizeRepoPath(repoPath)

	query := `
		SELECT id, repo_path, view_date, view_count, created_at, updated_at
		FROM repo_views
		WHERE repo_path = $1
		AND view_date BETWEEN $2 AND $3
		ORDER BY view_date ASC
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

	rows, err := stmt.Query(normalizedPath, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = fmt.Errorf("rows close error: %v: %w", err, err)
		}
	}(rows)

	var views []*models.RepoView
	for rows.Next() {
		view := &models.RepoView{}
		err := rows.Scan(
			&view.ID,
			&view.RepoPath,
			&view.ViewDate,
			&view.ViewCount,
			&view.CreatedAt,
			&view.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		views = append(views, view)
	}

	return views, nil
}

// GetTopRepositories returns the most viewed repositories
func (r *RepoViewRepository) GetTopRepositories(limit int) ([]*models.RepoView, error) {
	query := `
		SELECT repo_path, SUM(view_count) as total_views
		FROM repo_views
		GROUP BY repo_path
		ORDER BY total_views DESC
		LIMIT $1
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

	rows, err := stmt.Query(limit)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = fmt.Errorf("rows close error: %v: %w", err, err)
		}
	}(rows)

	var results []*models.RepoView
	for rows.Next() {
		view := &models.RepoView{}
		var totalViews int
		err := rows.Scan(&view.RepoPath, &totalViews)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		view.ViewCount = totalViews
		results = append(results, view)
	}

	return results, nil
}

// GetTopRepositoriesInDateRange returns the most viewed repositories within a date range
func (r *RepoViewRepository) GetTopRepositoriesInDateRange(startDate, endDate time.Time, limit int) ([]*models.RepoView, error) {
	query := `
		SELECT repo_path, SUM(view_count) as total_views
		FROM repo_views
		WHERE view_date BETWEEN $1 AND $2
		GROUP BY repo_path
		ORDER BY total_views DESC
		LIMIT $3
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

	rows, err := stmt.Query(startDate, endDate, limit)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = fmt.Errorf("rows close error: %v: %w", err, err)
		}
	}(rows)

	var results []*models.RepoView
	for rows.Next() {
		view := &models.RepoView{}
		var totalViews int
		err := rows.Scan(&view.RepoPath, &totalViews)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		view.ViewCount = totalViews
		results = append(results, view)
	}

	return results, nil
}
