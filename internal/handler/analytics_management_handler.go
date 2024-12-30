package handler

import (
	"bytes"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
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
	MaxValue  int
	Dates     []string
	Pages     []string
	Colors    map[string]string
	LineData  map[string][]int
	ChartHTML string // New field for the interactive chart
}

func prepareChartData(stats map[string]map[string]int, dates []string, pages []string, colors map[string]string) string {
	line := charts.NewLine()

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

	line.SetXAxis(dates)

	// Add data series
	for _, page := range pages {
		var values []opts.LineData
		for _, date := range dates {
			values = append(values, opts.LineData{
				Value: stats[date][page],
			})
		}

		line.AddSeries(page, values).
			SetSeriesOptions(
				charts.WithLineChartOpts(opts.LineChart{}),
				charts.WithItemStyleOpts(opts.ItemStyle{
					Color: colors[page],
				}),
			)
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
	colors := map[string]string{
		"home":     "#3b82f6",
		"about":    "#10b981",
		"blogs":    "#f59e0b",
		"projects": "#ef4444",
		"feedback": "#8b5cf6",
	}

	dates := make([]string, 0, len(stats))
	for date := range stats {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	pages := []string{"home", "about", "blogs", "projects", "feedback"}

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

	chartHTML := prepareChartData(stats, dates, pages, colors)

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
