package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"prosamik-backend/internal/database"
	"prosamik-backend/pkg/models"
	"strings"
)

type BlogRepository struct {
	db *sql.DB
}

func NewBlogRepository() *BlogRepository {
	return &BlogRepository{
		db: database.DB,
	}
}

var closeErr error

// Helper function to normalize strings
func normalizeBlogString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// GetBlogByTitle retrieves a blog by its title
func (r *BlogRepository) GetBlogByTitle(title string) (*models.Blog, error) {
	query := `
        SELECT id, title, path, description, tags, views_count
        FROM blogs
        WHERE LOWER(title) = LOWER($1)
    `

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("prepare statement error: %w", err)
	}

	var closeErr error
	defer func() {
		if cerr := closeStmt(stmt); cerr != nil {
			// If there's no error from the function, use the close error
			if closeErr == nil {
				closeErr = cerr
			}
		}
	}()

	blog := &models.Blog{}
	err = stmt.QueryRow(normalizeBlogString(title)).Scan(
		&blog.ID,
		&blog.Title,
		&blog.Path,
		&blog.Description,
		&blog.Tags,
		&blog.ViewsCount,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}

	return blog, nil
}

// GetAllBlogs retrieves all blog posts
func (r *BlogRepository) GetAllBlogs() ([]*models.Blog, error) {
	query := `
        SELECT id, title, path, description, tags, views_count
        FROM blogs
        ORDER BY id DESC
    `

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("prepare statement error: %w", err)
	}

	defer func() {
		if cerr := closeStmt(stmt); cerr != nil {
			// If there's no error from the function, use the close error
			if closeErr == nil {
				closeErr = cerr
			}
		}
	}()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer func() {
		if cerr := closeStmt(stmt); cerr != nil {
			// If there's no error from the function, use the close error
			if closeErr == nil {
				closeErr = cerr
			}
		}
	}()

	var blogs []*models.Blog
	for rows.Next() {
		blog := &models.Blog{}
		err := rows.Scan(
			&blog.ID,
			&blog.Title,
			&blog.Path,
			&blog.Description,
			&blog.Tags,
			&blog.ViewsCount,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		blogs = append(blogs, blog)
	}

	return blogs, nil
}

// CreateBlog adds a new blog post
func (r *BlogRepository) CreateBlog(blog *models.Blog) error {
	query := `
        INSERT INTO blogs (title, path, description, tags)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("prepare statement error: %w", err)
	}

	defer func() {
		if cerr := closeStmt(stmt); cerr != nil {
			// If there's no error from the function, use the close error
			if closeErr == nil {
				closeErr = cerr
			}
		}
	}()

	err = stmt.QueryRow(
		strings.TrimSpace(blog.Title),
		strings.TrimSpace(blog.Path),
		strings.TrimSpace(blog.Description),
		strings.TrimSpace(blog.Tags),
	).Scan(&blog.ID)

	if err != nil {
		return fmt.Errorf("create blog error: %w", err)
	}

	return nil
}

// UpdateBlog updates an existing blog post
func (r *BlogRepository) UpdateBlog(blog *models.Blog) error {
	query := `
        UPDATE blogs
        SET title = $1, path = $2, description = $3, tags = $4
        WHERE id = $5
    `

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("prepare statement error: %w", err)
	}

	defer func() {
		if cerr := closeStmt(stmt); cerr != nil {
			// If there's no error from the function, use the close error
			if closeErr == nil {
				closeErr = cerr
			}
		}
	}()

	result, err := stmt.Exec(
		strings.TrimSpace(blog.Title),
		strings.TrimSpace(blog.Path),
		strings.TrimSpace(blog.Description),
		strings.TrimSpace(blog.Tags),
		blog.ID,
	)
	if err != nil {
		return fmt.Errorf("update error: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected error: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no blog found with id: %d", blog.ID)
	}

	return nil
}

// DeleteBlog removes a blog post
func (r *BlogRepository) DeleteBlog(id int64) error {
	query := `
        DELETE FROM blogs
        WHERE id = $1
    `

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("prepare statement error: %w", err)
	}

	defer func() {
		if cerr := closeStmt(stmt); cerr != nil {
			// If there's no error from the function, use the close error
			if closeErr == nil {
				closeErr = cerr
			}
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
		return fmt.Errorf("no blog found with id: %d", id)
	}

	return nil
}

// SearchBlogs searches for blogs by title, path, tags, or description
func (r *BlogRepository) SearchBlogs(query string) ([]*models.Blog, error) {
	searchQuery := `
        SELECT id, title, path, description, tags, views_count
        FROM blogs
        WHERE LOWER(title) LIKE LOWER($1) 
           OR LOWER(path) LIKE LOWER($1)
           OR LOWER(tags) LIKE LOWER($1)
           OR LOWER(description) LIKE LOWER($1)
        ORDER BY id DESC
    `

	stmt, err := r.db.Prepare(searchQuery)
	if err != nil {
		return nil, fmt.Errorf("prepare statement error: %w", err)
	}

	defer func() {
		if cerr := closeStmt(stmt); cerr != nil {
			// If there's no error from the function, use the close error
			if closeErr == nil {
				closeErr = cerr
			}
		}
	}()

	normalizedQuery := "%" + normalizeBlogString(query) + "%"
	rows, err := stmt.Query(normalizedQuery)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = fmt.Errorf("rows close error: %v: %w", err, err)
		}
	}(rows)

	var blogs []*models.Blog
	for rows.Next() {
		blog := &models.Blog{}
		err := rows.Scan(
			&blog.ID,
			&blog.Title,
			&blog.Path,
			&blog.Description,
			&blog.Tags,
			&blog.ViewsCount,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		blogs = append(blogs, blog)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return blogs, nil
}

// GetBlog retrieves a single blog by ID
func (r *BlogRepository) GetBlog(id int64) (*models.Blog, error) {
	query := `
        SELECT id, title, path, description, tags, views_count
        FROM blogs
        WHERE id = $1
    `

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("prepare statement error: %w", err)
	}

	defer func() {
		if cerr := closeStmt(stmt); cerr != nil {
			// If there's no error from the function, use the close error
			if closeErr == nil {
				closeErr = cerr
			}
		}
	}()

	blog := &models.Blog{}
	err = stmt.QueryRow(id).Scan(
		&blog.ID,
		&blog.Title,
		&blog.Path,
		&blog.Description,
		&blog.Tags,
		&blog.ViewsCount,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}

	return blog, nil
}

// GetBlogByPath retrieves a blog by its path
func (r *BlogRepository) GetBlogByPath(path string) (*models.Blog, error) {
	query := `
        SELECT id, title, path, description, tags, views_count
        FROM blogs
        WHERE path = $1
    `

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("prepare statement error: %w", err)
	}

	defer func() {
		if cerr := closeStmt(stmt); cerr != nil {
			// If there's no error from the function, use the close error
			if closeErr == nil {
				closeErr = cerr
			}
		}
	}()

	blog := &models.Blog{}
	err = stmt.QueryRow(path).Scan(
		&blog.ID,
		&blog.Title,
		&blog.Path,
		&blog.Description,
		&blog.Tags,
		&blog.ViewsCount,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	return blog, nil
}

// Helper function to handle statement closing
func closeStmt(stmt *sql.Stmt) error {
	if err := stmt.Close(); err != nil {
		return fmt.Errorf("error closing statement: %w", err)
	}
	return nil
}
