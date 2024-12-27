package handler

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"prosamik-backend/internal/auth"
)

var templates = template.Must(template.ParseGlob("internal/templates/*.html"))

type PageData struct {
	Page  string
	Data  interface{}
	Error string
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/samik", http.StatusSeeOther)
		return
	}
	// If it's any other unhandled path, return 404
	http.NotFound(w, r)
}

func HandleAdminLoginUsingJWT(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		data := PageData{
			Page: "login",
		}
		err := templates.ExecuteTemplate(w, "base", data)
		if err != nil {
			log.Printf("Template error: %v", err)
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}

	case "POST":
		username := r.FormValue("username")
		password := r.FormValue("password")
		envPassword := os.Getenv("ADMIN_PASSWORD")

		if username == "" || password == "" {
			data := PageData{
				Page:  "login",
				Error: "Username and password are required",
			}
			if err := templates.ExecuteTemplate(w, "base", data); err != nil {
				log.Printf("Template error: %v", err)
				http.Error(w, "Failed to render template", http.StatusInternalServerError)
			}
			return
		}

		if username != "admin" || password != envPassword {
			data := PageData{
				Page:  "login",
				Error: "Invalid username or password",
			}
			if err := templates.ExecuteTemplate(w, "base", data); err != nil {
				log.Printf("Template error: %v", err)
				http.Error(w, "Failed to render template", http.StatusInternalServerError)
			}
			return
		}

		// If we get here, credentials are valid
		token, err := auth.GenerateToken(username)
		if err != nil {
			log.Printf("Token generation error: %v", err)
			data := PageData{
				Page:  "login",
				Error: "Authentication error occurred",
			}
			if err := templates.ExecuteTemplate(w, "base", data); err != nil {
				log.Printf("Template error: %v", err)
				http.Error(w, "Failed to render template", http.StatusInternalServerError)
			}
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			MaxAge:   24 * 60 * 60,
		})

		// Successful login
		http.Redirect(w, r, "/samik", http.StatusSeeOther)
		return

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleDashboard(w http.ResponseWriter, r *http.Request) {
	// Get username from JWT token for personalized welcome
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		http.Redirect(w, r, "/samik/login", http.StatusSeeOther)
		return
	}

	claims, err := auth.ValidateToken(cookie.Value)
	if err != nil {
		log.Printf("Token validation error: %v", err)
		http.Redirect(w, r, "/samik/login", http.StatusSeeOther)
		return
	}

	data := PageData{
		Page: "dashboard",
		Data: claims.Username,
	}
	err = templates.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

func HandleAdminLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	log.Printf("User logged out successfully")
	http.Redirect(w, r, "/samik/login", http.StatusSeeOther)
}
