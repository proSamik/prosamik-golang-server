package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"prosamik-backend/internal/repository"
	"time"
)

// RepoViewsAnalyticsData holds the data for the repo view analytics template
type RepoViewsAnalyticsData struct {
	TopRepositories   []RepoViewItem            `json:"top_repositories"`
	ViewsByDateRange  map[string]map[string]int `json:"views_by_date_range"`
	StartDate         string                    `json:"start_date"`
	EndDate           string                    `json:"end_date"`
	ChartHTML         template.HTML             `json:"chart_html"`
	TotalViews        int                       `json:"total_views"`
	ChartDataRepoPath []string                  `json:"chart_data_repo_path"`
	ChartDataViews    []int                     `json:"chart_data_views"`
}

// HandleRepoViewAnalytics handles the repo view analytics page
func HandleRepoViewAnalytics(w http.ResponseWriter, r *http.Request) {
	// Get top repositories (default 10)
	repoViewRepo := repository.NewRepoViewRepository()
	topRepos, err := repoViewRepo.GetTopRepositories(10)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching top repositories: %v", err), http.StatusInternalServerError)
		return
	}

	// Format for template
	topRepoItems := make([]RepoViewItem, 0, len(topRepos))
	chartLabels := make([]string, 0, len(topRepos))
	chartData := make([]int, 0, len(topRepos))
	totalViews := 0

	for _, item := range topRepos {
		// Extract the repository name from the full path for chart labels
		repoName := extractRepoNameFromPath(item.RepoPath)

		topRepoItems = append(topRepoItems, RepoViewItem{
			RepoPath:  item.RepoPath,
			ViewCount: item.ViewCount,
		})

		chartLabels = append(chartLabels, repoName)
		chartData = append(chartData, item.ViewCount)
		totalViews += item.ViewCount
	}

	// Default to last 7 days, starting from 2 days ago
	endDate := time.Now().UTC().AddDate(0, 0, -2)
	startDate := endDate.AddDate(0, 0, -7)

	// Create data for the template
	data := RepoViewsAnalyticsData{
		TopRepositories:   topRepoItems,
		StartDate:         startDate.Format("2006-01-02"),
		EndDate:           endDate.Format("2006-01-02"),
		TotalViews:        totalViews,
		ChartDataRepoPath: chartLabels,
		ChartDataViews:    chartData,
	}

	// Generate chart HTML
	data.ChartHTML = generateRepoViewChartHTML(chartLabels, chartData)

	// Use the global templates variable and proper page name
	pageData := PageData{
		Page: "repo-view-analytics", // Use the exact page name that matches in base.html
		Data: data,
	}

	// Execute the template
	err = templates.ExecuteTemplate(w, "base", pageData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
	}
}

// HandleRepoViewFilter handles filtering repo views by date range
func HandleRepoViewFilter(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	startDateStr := r.Form.Get("startDate")
	endDateStr := r.Form.Get("endDate")
	repoPath := r.Form.Get("repoPath")

	// Parse dates
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		// Default to last 7 days, starting from 2 days ago
		endDate := time.Now().UTC().AddDate(0, 0, -2)
		startDate = endDate.AddDate(0, 0, -7)
		startDateStr = startDate.Format("2006-01-02")
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		endDate = time.Now().UTC().AddDate(0, 0, -2) // 2 days ago
		endDateStr = endDate.Format("2006-01-02")
	}

	repoViewRepo := repository.NewRepoViewRepository()

	var data RepoViewsAnalyticsData

	// If a specific repo is selected, get its views by date range
	if repoPath != "" {
		views, err := repoViewRepo.GetViewsByDateRange(repoPath, startDate, endDate)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching views by date range: %v", err), http.StatusInternalServerError)
			return
		}

		// Format data for the chart
		viewsByDate := make(map[string]map[string]int)
		dateLabels := make([]string, 0)
		viewCounts := make([]int, 0)
		total := 0

		for _, view := range views {
			dateStr := view.ViewDate.Format("2006-01-02")

			if viewsByDate[dateStr] == nil {
				viewsByDate[dateStr] = make(map[string]int)
				dateLabels = append(dateLabels, dateStr)
			}

			viewsByDate[dateStr][repoPath] = view.ViewCount
			viewCounts = append(viewCounts, view.ViewCount)
			total += view.ViewCount
		}

		data = RepoViewsAnalyticsData{
			ViewsByDateRange: viewsByDate,
			StartDate:        startDateStr,
			EndDate:          endDateStr,
			TotalViews:       total,
			ChartHTML:        generateDateRangeChartHTML(dateLabels, viewCounts, extractRepoNameFromPath(repoPath)),
		}
	} else {
		// Otherwise get top repos for the date range
		topRepos, err := repoViewRepo.GetTopRepositoriesInDateRange(startDate, endDate, 10)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching top repositories: %v", err), http.StatusInternalServerError)
			return
		}

		// Format for template
		topRepoItems := make([]RepoViewItem, 0, len(topRepos))
		chartLabels := make([]string, 0, len(topRepos))
		chartData := make([]int, 0, len(topRepos))
		totalViews := 0

		for _, item := range topRepos {
			repoName := extractRepoNameFromPath(item.RepoPath)

			topRepoItems = append(topRepoItems, RepoViewItem{
				RepoPath:  item.RepoPath,
				ViewCount: item.ViewCount,
			})

			chartLabels = append(chartLabels, repoName)
			chartData = append(chartData, item.ViewCount)
			totalViews += item.ViewCount
		}

		data = RepoViewsAnalyticsData{
			TopRepositories:   topRepoItems,
			StartDate:         startDateStr,
			EndDate:           endDateStr,
			TotalViews:        totalViews,
			ChartDataRepoPath: chartLabels,
			ChartDataViews:    chartData,
			ChartHTML:         generateRepoViewChartHTML(chartLabels, chartData),
		}
	}

	// Create a PageData struct for the template
	pageData := PageData{
		Page: "repo-view-analytics",
		Data: data,
	}

	// Use the global templates for rendering
	err = templates.ExecuteTemplate(w, "repo-view-data-section", pageData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
	}
}

// extractRepoNameFromPath extracts a readable repo name from a GitHub URL
func extractRepoNameFromPath(path string) string {
	// For GitHub URLs, we'll extract the owner/repo part
	const prefix = "https://github.com/"
	if len(path) > len(prefix) && path[:len(prefix)] == prefix {
		parts := path[len(prefix):]

		// Find the first slash or end of string
		slashIndex := -1
		for i, char := range parts {
			if char == '/' {
				slashIndex = i
				break
			}
		}

		if slashIndex != -1 {
			// If there's a slash, extract owner/repo
			secondSlashIndex := -1
			for i, char := range parts[slashIndex+1:] {
				if char == '/' {
					secondSlashIndex = slashIndex + 1 + i
					break
				}
			}

			if secondSlashIndex != -1 {
				return parts[:secondSlashIndex]
			}
			return parts
		}
		return parts
	}
	return path
}

// generateRepoViewChartHTML generates the HTML for the repository views chart
func generateRepoViewChartHTML(labels []string, data []int) template.HTML {
	// Convert labels and data to JavaScript arrays
	labelsJSON := "["
	for i, label := range labels {
		if i > 0 {
			labelsJSON += ", "
		}
		labelsJSON += fmt.Sprintf("'%s'", label)
	}
	labelsJSON += "]"

	dataJSON := "["
	for i, value := range data {
		if i > 0 {
			dataJSON += ", "
		}
		dataJSON += fmt.Sprintf("%d", value)
	}
	dataJSON += "]"

	// Generate the chart HTML using Chart.js
	chartHTML := fmt.Sprintf(`
        <canvas id="repoViewsChart"></canvas>
        <script>
            // Use immediately invoked function to initialize chart as soon as it's added to DOM
            (function() {
                // Destroy any existing chart with same ID to prevent duplicates
                const existingChart = Chart.getChart("repoViewsChart");
                if (existingChart) {
                    existingChart.destroy();
                }
                
                const ctx = document.getElementById('repoViewsChart').getContext('2d');
                const repoViewsChart = new Chart(ctx, {
                    type: 'bar',
                    data: {
                        labels: %s,
                        datasets: [{
                            label: 'Repository Views',
                            data: %s,
                            backgroundColor: 'rgba(54, 162, 235, 0.6)',
                            borderColor: 'rgba(54, 162, 235, 1)',
                            borderWidth: 1
                        }]
                    },
                    options: {
                        responsive: true,
                        scales: {
                            y: {
                                beginAtZero: true,
                                title: {
                                    display: true,
                                    text: 'Views'
                                }
                            },
                            x: {
                                title: {
                                    display: true,
                                    text: 'Repository'
                                }
                            }
                        },
                        plugins: {
                            legend: {
                                position: 'top',
                            },
                            title: {
                                display: true,
                                text: 'Top Repositories by Views'
                            }
                        }
                    }
                });
            })();
        </script>
    `, labelsJSON, dataJSON)

	return template.HTML(chartHTML)
}

// generateDateRangeChartHTML generates the HTML for the date range chart
func generateDateRangeChartHTML(dates []string, viewCounts []int, repoName string) template.HTML {
	// Convert dates and data to JavaScript arrays
	datesJSON := "["
	for i, date := range dates {
		if i > 0 {
			datesJSON += ", "
		}
		datesJSON += fmt.Sprintf("'%s'", date)
	}
	datesJSON += "]"

	dataJSON := "["
	for i, value := range viewCounts {
		if i > 0 {
			dataJSON += ", "
		}
		dataJSON += fmt.Sprintf("%d", value)
	}
	dataJSON += "]"

	// Generate the chart HTML using Chart.js
	chartHTML := fmt.Sprintf(`
        <canvas id="dateRangeChart"></canvas>
        <script>
            // Use immediately invoked function to initialize chart as soon as it's added to DOM
            (function() {
                // Destroy any existing chart with same ID to prevent duplicates
                const existingChart = Chart.getChart("dateRangeChart");
                if (existingChart) {
                    existingChart.destroy();
                }
                
                const ctx = document.getElementById('dateRangeChart').getContext('2d');
                const dateRangeChart = new Chart(ctx, {
                    type: 'line',
                    data: {
                        labels: %s,
                        datasets: [{
                            label: 'Daily Views for %s',
                            data: %s,
                            backgroundColor: 'rgba(75, 192, 192, 0.6)',
                            borderColor: 'rgba(75, 192, 192, 1)',
                            borderWidth: 2,
                            tension: 0.3,
                            fill: true
                        }]
                    },
                    options: {
                        responsive: true,
                        scales: {
                            y: {
                                beginAtZero: true,
                                title: {
                                    display: true,
                                    text: 'Views'
                                }
                            },
                            x: {
                                title: {
                                    display: true,
                                    text: 'Date'
                                }
                            }
                        },
                        plugins: {
                            legend: {
                                position: 'top',
                            },
                            title: {
                                display: true,
                                text: 'Repository Views Over Time'
                            }
                        }
                    }
                });
            })();
        </script>
    `, datesJSON, repoName, dataJSON)

	return template.HTML(chartHTML)
}
