package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"log"
	"net/http"
	"prosamik-backend/internal/cache"
	"prosamik-backend/internal/repository"
	"sort"
	"time"
)

// AnalyticsManagementData holds the structured analytics data
type AnalyticsManagementData struct {
	Stats     map[string]map[string]int
	StartDate string
	EndDate   string
	MaxValue  int
	Dates     []string
	Pages     []string
	Colors    map[string]string
	LineData  map[string][]int
	ChartHTML string // New field for the interactive chart
}

func prepareChartData(stats map[string]map[string]int, dates []string) string {
	line := charts.NewLine()

	// Format dates for display
	formattedDates := make([]string, len(dates))
	for i, date := range dates {
		t, err := time.Parse("2006-01-02", date)
		if err == nil {
			formattedDates[i] = t.Format("02 Jan") // Format as "dd MMM"
		} else {
			formattedDates[i] = date
		}
	}

	// Basic chart configuration
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "100%",
			Height: "400px",
		}),
		charts.WithTitleOpts(opts.Title{
			Title: "Page Views Over Time",
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Trigger: "axis",
		}),
		charts.WithLegendOpts(opts.Legend{
			Right:  "10%",
			Orient: "vertical",
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "slider",
			XAxisIndex: []int{0},
			Start:      0,
			End:        100,
		}),
	)

	line.SetXAxis(formattedDates)

	// Define the pages in groups
	pageGroups := map[string][]struct {
		key   string // key in stats map
		name  string // display name
		color string // line color
	}{
		"Main": {
			{key: "home", name: "Home", color: "#3b82f6"},         // blue
			{key: "about", name: "About", color: "#10b981"},       // green
			{key: "blogs", name: "Blogs", color: "#f59e0b"},       // yellow
			{key: "projects", name: "Projects", color: "#ef4444"}, // red
			{key: "feedback", name: "Feedback", color: "#8b5cf6"}, // purple
		},
		"Githubme": {
			{key: "githubme_home", name: "Githubme Home", color: "#6366f1"},         // indigo
			{key: "githubme_about", name: "Githubme About", color: "#ec4899"},       // pink
			{key: "githubme_markdown", name: "Githubme Markdown", color: "#14b8a6"}, // teal
		},
	}

	// Add series for each group
	for groupName, pages := range pageGroups {
		for _, page := range pages {
			var values []opts.LineData
			for _, date := range dates {
				values = append(values, opts.LineData{
					Value: stats[date][page.key],
				})
			}

			// Configure line style based on the group
			lineStyle := opts.LineStyle{
				Width: 2,
			}
			if groupName == "Githubme" {
				lineStyle.Type = "dashed"
			}

			// Add series with options
			series := line.AddSeries(page.name, values)
			series.SetSeriesOptions(
				charts.WithLineStyleOpts(lineStyle),
				charts.WithItemStyleOpts(opts.ItemStyle{
					Color: page.color,
				}),
			)
		}
	}

	// Render to HTML
	buf := new(bytes.Buffer)
	err := line.Render(buf)
	if err != nil {
		log.Printf("Error rendering chart: %v", err)
		return ""
	}
	return buf.String()
}

func prepareAnalyticsData(stats map[string]map[string]int) AnalyticsManagementData {
	// Define page groups
	pageGroups := map[string][]string{
		"Main": {
			"home", "about", "blogs", "projects", "feedback",
		},
		"Githubme": {
			"githubme_home", "githubme_about", "githubme_markdown",
		},
	}

	// Define colors by group
	colors := map[string]string{
		// Main pages
		"home":     "#3b82f6", // blue
		"about":    "#10b981", // green
		"blogs":    "#f59e0b", // yellow
		"projects": "#ef4444", // red
		"feedback": "#8b5cf6", // purple

		// Githubme pages (using a different color palette)
		"githubme_home":     "#6366f1", // indigo
		"githubme_about":    "#ec4899", // pink
		"githubme_markdown": "#14b8a6", // teal
	}

	// Get and sort dates
	dates := make([]string, 0, len(stats))
	for date := range stats {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	// Combine all pages while maintaining group order
	var pages []string
	for _, groupPages := range pageGroups {
		pages = append(pages, groupPages...)
	}

	// Calculate line data and max value
	lineData := make(map[string][]int)
	maxValue := 0

	// Calculate by group for better organization
	for _, groupPages := range pageGroups {
		for _, page := range groupPages {
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
	}

	// Prepare data for the chart
	chartHTML := prepareChartData(stats, dates)

	return AnalyticsManagementData{
		Stats:     stats,
		MaxValue:  maxValue,
		Dates:     dates,
		Pages:     pages,
		Colors:    colors,
		LineData:  lineData,
		ChartHTML: chartHTML,
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

func HandleCacheMonitoring(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cacheStats, err := cache.GetCacheStats(r.Context())
	if err != nil {
		fmt.Printf("Failed to get cache stats: %v\n", err)
		http.Error(w, "Failed to get cache statistics", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Page: "cache-monitoring",
		Data: cacheStats,
	}

	// Execute the base template which includes cache-monitoring
	err = templates.ExecuteTemplate(w, "base", data) // Changed from "cache-monitoring" to "base"
	if err != nil {
		fmt.Printf("Template error: %v\n", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

func HandleCacheStats(w http.ResponseWriter, r *http.Request) {
	stats, err := cache.GetCacheStats(r.Context())
	if err != nil {
		fmt.Printf("Error getting cache stats: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Page: "cache-monitoring",
		Data: stats,
	}

	if r.Header.Get("HX-Request") == "true" {
		// Only render the cache-stats section for HTMX requests
		err = templates.ExecuteTemplate(w, "cache-stats", data.Data) // Changed to pass just the stats data
	} else {
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(stats)
	}

	if err != nil {
		fmt.Printf("Error rendering response: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
