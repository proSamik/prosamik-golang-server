package handler

import (
	"fmt"
	"log"
	"net/http"
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

	// Create the blog
	err = repo.CreateBlog(blog)
	if err != nil {
		log.Printf("Error adding blog: %v", err)
		renderFormError(w, "Failed to add blog")
		return
	}

	// Just render the success message
	w.Header().Set("Content-Type", "text/html")
	err = templates.ExecuteTemplate(w, "blog-form-message", struct{ Error string }{Error: ""})
	if err != nil {
		log.Printf("Template error: %v", err)
	}
}

// HandleBlogEdit handles rendering edit inputs for a specific blog field
func HandleBlogEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract blog ID and field from URL
	segments := strings.Split(r.URL.Path, "/")
	if len(segments) < 6 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(segments[len(segments)-2], 10, 64)
	if err != nil {
		log.Printf("Invalid blog ID: %v", err)
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	field := segments[len(segments)-1]

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

	// Prepare data for edit input template
	editData := struct {
		ID     int64
		Field  string
		Value  string
		Target string
	}{
		ID:     id,
		Field:  field,
		Target: fmt.Sprintf("#%s-cell-%d", field, id),
	}

	// Get the correct value based on the field
	switch field {
	case "title":
		editData.Value = blog.Title
	case "path":
		editData.Value = blog.Path
	case "description":
		editData.Value = blog.Description
	case "tags":
		editData.Value = blog.Tags
	default:
		http.Error(w, "Invalid field", http.StatusBadRequest)
		return
	}

	// Render edit input template
	err = templates.ExecuteTemplate(w, "edit-input", editData)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

// HandleBlogUpdate handles updating a specific blog field
func HandleBlogUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
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

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
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

	// Get the first form key
	var fieldName string
	for key := range r.Form {
		fieldName = key
		break
	}

	// Update the specific field
	fieldValue := strings.TrimSpace(r.FormValue(fieldName))
	switch fieldName {
	case "title":
		if fieldValue == "" {
			http.Error(w, "Title cannot be empty", http.StatusBadRequest)
			return
		}
		blog.Title = fieldValue
	case "path":
		if fieldValue == "" {
			http.Error(w, "Path cannot be empty", http.StatusBadRequest)
			return
		}
		blog.Path = fieldValue
	case "description":
		blog.Description = fieldValue
	case "tags":
		blog.Tags = fieldValue
	default:
		http.Error(w, "Invalid field", http.StatusBadRequest)
		return
	}

	// Update the blog
	err = repo.UpdateBlog(blog)
	if err != nil {
		log.Printf("Error updating blog: %v", err)
		http.Error(w, "Failed to update blog", http.StatusInternalServerError)
		return
	}

	// Return the updated value
	_, err = w.Write([]byte(fieldValue))
	if err != nil {
		log.Printf("Error writing response: %v", err)
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

// HandleBlogCancelEdit handles canceling edit for a specific blog field
func HandleBlogCancelEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract blog ID and field from URL
	segments := strings.Split(r.URL.Path, "/")
	if len(segments) < 5 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(segments[4], 10, 64)
	if err != nil {
		log.Printf("Invalid blog ID: %v", err)
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	field := segments[5]

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

	// Return the original value based on the field
	var value string
	switch field {
	case "title":
		value = blog.Title
	case "path":
		value = blog.Path
	case "description":
		value = blog.Description
	case "tags":
		value = blog.Tags
	default:
		http.Error(w, "Invalid field", http.StatusBadRequest)
		return
	}

	_, err = w.Write([]byte(value))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
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
