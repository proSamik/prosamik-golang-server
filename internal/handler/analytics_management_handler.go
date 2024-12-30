package handler

import (
	"log"
	"net/http"
	"prosamik-backend/internal/repository"
	"sort"
	"time"
)

// AnalyticsManagementData holds the structured analytics data
type AnalyticsManagementData struct {
	Stats     map[string]map[string]int
	StartDate string
	EndDate   string
	MaxValue  int               // Added from GraphData
	Dates     []string          // Added from GraphData
	Pages     []string          // Added from GraphData
	Colors    map[string]string // Added from GraphData
	LineData  map[string][]int  // Added from GraphData
}

// Since we're not using GraphData struct anymore, rename function to be more specific
func prepareAnalyticsData(stats map[string]map[string]int) AnalyticsManagementData {
	// Initialize colors for each page
	colors := map[string]string{
		"home":     "#3b82f6", // blue
		"about":    "#10b981", // green
		"blogs":    "#f59e0b", // amber
		"projects": "#ef4444", // red
		"feedback": "#8b5cf6", // purple
	}

	// Get sorted dates
	dates := make([]string, 0, len(stats))
	for date := range stats {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	// Define pages in specific order
	pages := []string{"home", "about", "blogs", "projects", "feedback"}

	// Prepare line data and find max value
	lineData := make(map[string][]int)
	maxValue := 0

	for _, page := range pages {
		values := make([]int, len(dates))
		for i, date := range dates {
			value := stats[date][page]
			values[i] = value
			if value > maxValue {
				maxValue = value
			}
		}
		lineData[page] = values
	}

	// Return AnalyticsManagementData directly
	return AnalyticsManagementData{
		Stats:    stats,
		MaxValue: maxValue,
		Dates:    dates,
		Pages:    pages,
		Colors:   colors,
		LineData: lineData,
	}
}

func HandleAnalyticsManagement(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get date range from query parameters
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	// Default to last 7 days if no dates provided
	if startDate == "" || endDate == "" {
		now := time.Now()
		endDate = now.Format("2006-01-02")
		startDate = now.AddDate(0, 0, -7).Format("2006-01-02")
	}

	repo := repository.NewAnalyticsRepository()
	stats, err := repo.GetAnalytics(startDate, endDate)
	if err != nil {
		log.Printf("Error fetching analytics: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get all data at once
	analyticsData := prepareAnalyticsData(stats)
	// Add the dates to the data
	analyticsData.StartDate = startDate
	analyticsData.EndDate = endDate

	data := PageData{
		Page: "analytics-management",
		Data: analyticsData,
	}

	if r.Header.Get("HX-Request") == "true" {
		err = templates.ExecuteTemplate(w, "analytics-data-section", data.Data)
	} else {
		err = templates.ExecuteTemplate(w, "base", data)
	}

	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// HandleAnalyticsFilter handles HTMX requests for filtering analytics
func HandleAnalyticsFilter(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	repo := repository.NewAnalyticsRepository()
	stats, err := repo.GetAnalytics(startDate, endDate)
	if err != nil {
		log.Printf("Error fetching analytics: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get all data at once
	analyticsData := prepareAnalyticsData(stats)
	// Add the dates to the data
	analyticsData.StartDate = startDate
	analyticsData.EndDate = endDate

	// Wrap the analytics data in PageData structure
	data := PageData{
		Page: "analytics-management",
		Data: analyticsData,
	}

	// Now pass the properly structured data to the template
	err = templates.ExecuteTemplate(w, "analytics-data-section", data)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
