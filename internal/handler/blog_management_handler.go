package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"prosamik-backend/internal/repository"
	"prosamik-backend/pkg/models"
	"strconv"
	"strings"
)

// Helper function to extract ID from URL path
func getBlogIDFromPath(path string) (int64, error) {
	segments := strings.Split(path, "/")
	if len(segments) < 4 {
		return 0, fmt.Errorf("invalid URL")
	}
	return strconv.ParseInt(segments[len(segments)-1], 10, 64)
}

// HandleBlogManagement renders the blog management page
func HandleBlogManagement(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	repo := repository.NewBlogRepository()
	blogs, err := repo.GetAllBlogs()
	if err != nil {
		log.Printf("Error fetching blogs: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Page: "blog-management",
		Data: struct{ Blogs []*models.Blog }{
			Blogs: blogs,
		},
	}

	err = templates.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

// HandleBlogSearch handles searching for blogs
func HandleBlogSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("search")

	repo := repository.NewBlogRepository()
	blogs, err := repo.SearchBlogs(query)
	if err != nil {
		log.Printf("Error searching blogs: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "blog-list", struct{ Blogs []*models.Blog }{Blogs: blogs})
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

// HandleBlogAdd handles adding a new blog
func HandleBlogAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		renderFormError(w, "Invalid form data")
		return
	}

	// Create blog from form data
	blog := &models.Blog{
		Title:       strings.TrimSpace(r.FormValue("title")),
		Path:        strings.TrimSpace(r.FormValue("path")),
		Description: strings.TrimSpace(r.FormValue("description")),
		Tags:        strings.TrimSpace(r.FormValue("tags")),
	}

	// Validate required fields
	if blog.Title == "" || blog.Path == "" {
		renderFormError(w, "Title and path are required")
		return
	}

	// Validate description length
	if len(blog.Description) > 5000 {
		renderFormError(w, "Description cannot exceed 5000 characters")
		return
	}

	// Validate path format
	if err := validatePath(blog.Path); err != nil {
		renderFormError(w, err.Error())
		return
	}

	// Validate and format tags
	validTags, err := validateTags(blog.Tags)
	if err != nil {
		renderFormError(w, err.Error())
		return
	}
	blog.Tags = validTags

	repo := repository.NewBlogRepository()

	// Check for existing title
	existing, err := repo.GetBlogByTitle(blog.Title)
	if err != nil {
		log.Printf("Error checking existing title: %v", err)
		renderFormError(w, "Internal server error")
		return
	}

	if existing != nil {
		renderFormError(w, "A blog with this title already exists")
		return
	}

	// Check for existing path
	existingPath, err := repo.GetBlogByPath(blog.Path)
	if err != nil {
		log.Printf("Error checking existing path: %v", err)
		renderFormError(w, "Internal server error")
		return
	}

	if existingPath != nil {
		renderFormError(w, "A blog with this path already exists")
		return
	}

	// Create the blog
	err = repo.CreateBlog(blog)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			if strings.Contains(err.Error(), "blogs_path_key") {
				renderFormError(w, "A blog with this path already exists")
			} else {
				renderFormError(w, "A blog with this title already exists")
			}
			return
		}
		log.Printf("Error adding blog: %v", err)
		renderFormError(w, "Failed to add blog")
		return
	}

	// Render a success message
	w.Header().Set("Content-Type", "text/html")
	err = templates.ExecuteTemplate(w, "blog-form-message", struct{ Error string }{Error: ""})
	if err != nil {
		log.Printf("Template error: %v", err)
	}
}

// HandleBlogEdit handles rendering the edit form for a blog post
func HandleBlogEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract blog ID from URL
	segments := strings.Split(r.URL.Path, "/")
	if len(segments) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(segments[len(segments)-1], 10, 64)
	if err != nil {
		log.Printf("Invalid blog ID: %v", err)
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	repo := repository.NewBlogRepository()
	blog, err := repo.GetBlog(id)
	if err != nil {
		log.Printf("Error fetching blog: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if blog == nil {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}

	// Render edit form template
	err = templates.ExecuteTemplate(w, "edit-form", blog)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

// HandleBlogUpdate handles updating a blog post
func HandleBlogUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract blog ID from URL
	segments := strings.Split(r.URL.Path, "/")
	if len(segments) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(segments[len(segments)-1], 10, 64)
	if err != nil {
		log.Printf("Invalid blog ID: %v", err)
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	// Parse form data
	err = r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Update blog in repository
	blog := &models.Blog{
		ID:          id,
		Title:       r.FormValue("title"),
		Path:        r.FormValue("path"),
		Description: r.FormValue("description"),
		Tags:        r.FormValue("tags"),
	}

	repo := repository.NewBlogRepository()
	err = repo.UpdateBlog(blog)
	if err != nil {
		log.Printf("Error updating blog: %v", err)
		http.Error(w, "Failed to update blog", http.StatusInternalServerError)
		return
	}

	// Return updated blog content
	err = templates.ExecuteTemplate(w, "blog-content", blog)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

// HandleBlogCancelEdit handles canceling blog edit
func HandleBlogCancelEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract blog ID from URL
	segments := strings.Split(r.URL.Path, "/")
	if len(segments) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(segments[len(segments)-1], 10, 64)
	if err != nil {
		log.Printf("Invalid blog ID: %v", err)
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	repo := repository.NewBlogRepository()
	blog, err := repo.GetBlog(id)
	if err != nil {
		log.Printf("Error fetching blog: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if blog == nil {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}

	// Return original blog content
	err = templates.ExecuteTemplate(w, "blog-content", blog)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

// HandleBlogDelete handles deleting a blog
func HandleBlogDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract blog ID from URL
	id, err := getBlogIDFromPath(r.URL.Path)
	if err != nil {
		log.Printf("Invalid blog ID: %v", err)
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	repo := repository.NewBlogRepository()
	err = repo.DeleteBlog(id)
	if err != nil {
		log.Printf("Error deleting blog: %v", err)
		http.Error(w, "Failed to delete blog", http.StatusInternalServerError)
		return
	}

	// Return successful status
	w.WriteHeader(http.StatusOK)
}

// Helper function to render form errors
func renderFormError(w http.ResponseWriter, message string) {
	err := templates.ExecuteTemplate(w, "blog-form-message", struct{ Error string }{
		Error: message,
	})
	if err != nil {
		log.Printf("Template error: %v", err)
	}
}

// validateTags checks if tags are properly formatted
func validateTags(tags string) (string, error) {
	if tags == "" {
		return "", nil
	}

	// Split tags and process each one
	tagList := strings.Split(tags, ",")
	var validTags []string

	for _, tag := range tagList {
		// Trim spaces
		tag = strings.TrimSpace(tag)

		// Check for periods
		if strings.Contains(tag, ".") {
			return "", fmt.Errorf("tags cannot contain periods: %s", tag)
		}

		if tag != "" {
			validTags = append(validTags, tag)
		}
	}

	return strings.Join(validTags, ","), nil
}

// validatePath checks if the path is a valid URL starting with http
func validatePath(path string) error {
	if !strings.HasPrefix(path, "http") {
		return fmt.Errorf("path must start with http:// or https://")
	}

	_, err := url.Parse(path)
	if err != nil {
		return fmt.Errorf("invalid URL format: %v", err)
	}

	return nil
}
