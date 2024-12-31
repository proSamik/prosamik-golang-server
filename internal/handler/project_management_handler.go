package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"prosamik-backend/internal/repository"
	"prosamik-backend/pkg/models"
	"strconv"
	"strings"
)

// Helper function to extract ID from URL path
func getProjectIDFromPath(path string) (int64, error) {
	segments := strings.Split(path, "/")
	if len(segments) < 4 {
		return 0, fmt.Errorf("invalid URL")
	}
	return strconv.ParseInt(segments[len(segments)-1], 10, 64)
}

// HandleProjectManagement renders the project management page
func HandleProjectManagement(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	repo := repository.NewProjectRepository()
	projects, err := repo.GetAllProjects()
	if err != nil {
		log.Printf("Error fetching projects: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Page: "project-management",
		Data: struct{ Projects []*models.Project }{
			Projects: projects,
		},
	}

	err = templates.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

// HandleProjectSearch handles searching for projects
func HandleProjectSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("search")

	repo := repository.NewProjectRepository()
	projects, err := repo.SearchProjects(query)
	if err != nil {
		log.Printf("Error searching projects: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "project-list", struct{ Projects []*models.Project }{Projects: projects})
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

// HandleProjectAdd handles adding a new project
func HandleProjectAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		renderProjectFormError(w, "Invalid form data")
		return
	}

	// Create project from form data
	project := &models.Project{
		Title:       strings.TrimSpace(r.FormValue("title")),
		Path:        strings.TrimSpace(r.FormValue("path")),
		Description: strings.TrimSpace(r.FormValue("description")),
		Tags:        strings.TrimSpace(r.FormValue("tags")),
	}

	// Validate required fields
	if project.Title == "" || project.Path == "" {
		renderProjectFormError(w, "Title and path are required")
		return
	}

	// Validate description length
	if len(project.Description) > 5000 {
		renderProjectFormError(w, "Description cannot exceed 5000 characters")
		return
	}

	// Validate path format
	if err := validateProjectPath(project.Path); err != nil {
		renderProjectFormError(w, err.Error())
		return
	}

	// Validate and format tags
	validTags, err := validateProjectTags(project.Tags)
	if err != nil {
		renderProjectFormError(w, err.Error())
		return
	}
	project.Tags = validTags

	repo := repository.NewProjectRepository()

	if err := validateProjectUniqueness(project, repo); err != nil {
		renderProjectFormError(w, err.Error())
		return
	}

	// Create the project
	err = repo.CreateProject(project)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			if strings.Contains(err.Error(), "projects_path_key") {
				renderProjectFormError(w, "A project with this path already exists")
			} else {
				renderProjectFormError(w, "A project with this title already exists")
			}
			return
		}
		log.Printf("Error adding project: %v", err)
		renderProjectFormError(w, "Failed to add project")
		return
	}

	// Render a success message
	w.Header().Set("Content-Type", "text/html")
	err = templates.ExecuteTemplate(w, "project-form-message", struct{ Error string }{Error: ""})
	if err != nil {
		log.Printf("Template error: %v", err)
	}
}

// HandleProjectEdit handles rendering the edit form for a project post
func HandleProjectEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract project ID from URL
	segments := strings.Split(r.URL.Path, "/")
	if len(segments) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(segments[len(segments)-1], 10, 64)
	if err != nil {
		log.Printf("Invalid project ID: %v", err)
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	repo := repository.NewProjectRepository()
	project, err := repo.GetProject(id)
	if err != nil {
		log.Printf("Error fetching project: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if project == nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	// Render edit form template
	err = templates.ExecuteTemplate(w, "edit-form", project)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

// HandleProjectUpdate handles updating a project post
func HandleProjectUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract project ID from URL
	segments := strings.Split(r.URL.Path, "/")
	if len(segments) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(segments[len(segments)-1], 10, 64)
	if err != nil {
		log.Printf("Invalid project ID: %v", err)
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	// Parse form data
	err = r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Update project in repository
	project := &models.Project{
		ID:          id,
		Title:       r.FormValue("title"),
		Path:        r.FormValue("path"),
		Description: r.FormValue("description"),
		Tags:        r.FormValue("tags"),
	}

	repo := repository.NewProjectRepository()
	err = repo.UpdateProject(project)
	if err != nil {
		log.Printf("Error updating project: %v", err)
		http.Error(w, "Failed to update project", http.StatusInternalServerError)
		return
	}

	// Return updated project content
	err = templates.ExecuteTemplate(w, "project-content", project)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

// HandleProjectCancelEdit handles canceling project edit
func HandleProjectCancelEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract project ID from URL
	segments := strings.Split(r.URL.Path, "/")
	if len(segments) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(segments[len(segments)-1], 10, 64)
	if err != nil {
		log.Printf("Invalid project ID: %v", err)
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	repo := repository.NewProjectRepository()
	project, err := repo.GetProject(id)
	if err != nil {
		log.Printf("Error fetching project: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if project == nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	// Return original project content
	err = templates.ExecuteTemplate(w, "project-content", project)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

// HandleProjectDelete handles deleting a project
func HandleProjectDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract project ID from URL
	id, err := getProjectIDFromPath(r.URL.Path)
	if err != nil {
		log.Printf("Invalid project ID: %v", err)
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	repo := repository.NewProjectRepository()
	err = repo.DeleteProject(id)
	if err != nil {
		log.Printf("Error deleting project: %v", err)
		http.Error(w, "Failed to delete project", http.StatusInternalServerError)
		return
	}

	// Return successful status
	w.WriteHeader(http.StatusOK)
}

// Helper function to render form errors
func renderProjectFormError(w http.ResponseWriter, message string) {
	err := templates.ExecuteTemplate(w, "project-form-message", struct{ Error string }{
		Error: message,
	})
	if err != nil {
		log.Printf("Template error: %v", err)
	}
}

// validateProjectTags checks if tags are properly formatted
func validateProjectTags(tags string) (string, error) {
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

// validateProjectPath checks if the path is a valid URL starting with http
func validateProjectPath(path string) error {
	if !strings.HasPrefix(path, "http") {
		return fmt.Errorf("path must start with or https://")
	}

	_, err := url.Parse(path)
	if err != nil {
		return fmt.Errorf("invalid URL format: %v", err)
	}

	return nil
}

// validateProjectUniqueness performs concurrent validation checks for project uniqueness and URL validity
func validateProjectUniqueness(project *models.Project, repo *repository.ProjectRepository) error {
	pathCheckChan := make(chan error, 1)
	titleCheckChan := make(chan error, 1)
	urlCheckChan := make(chan error, 1)

	// Check path existence in DB
	go func() {
		existingPath, err := repo.GetProjectByPath(project.Path)
		if err != nil {
			pathCheckChan <- fmt.Errorf("database error: %v", err)
			return
		}
		if existingPath != nil {
			pathCheckChan <- fmt.Errorf("a project with this path already exists")
			return
		}
		pathCheckChan <- nil
	}()

	// Check title existence in DB
	go func() {
		existingTitle, err := repo.GetProjectByTitle(project.Title)
		if err != nil {
			titleCheckChan <- fmt.Errorf("database error: %v", err)
			return
		}
		if existingTitle != nil {
			titleCheckChan <- fmt.Errorf("a project with this title already exists")
			return
		}
		titleCheckChan <- nil
	}()

	// Check URL validity using markdown handler
	go func() {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", fmt.Sprintf("/md?url=%s", project.Path), nil)
		if err != nil {
			urlCheckChan <- fmt.Errorf("failed to create request: %v", err)
			return
		}
		MarkdownHandler(w, req)
		if w.Code != http.StatusOK {
			urlCheckChan <- fmt.Errorf("content not found at specified URL")
			return
		}
		urlCheckChan <- nil
	}()

	// Wait for all checks and return first error
	if err := <-pathCheckChan; err != nil {
		return err
	}
	if err := <-titleCheckChan; err != nil {
		return err
	}
	if err := <-urlCheckChan; err != nil {
		return err
	}

	return nil
}
