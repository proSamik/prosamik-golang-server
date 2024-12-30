package models

import "time"

type Analytics struct {
	Date          string    `json:"date"`
	HomeViews     int       `json:"home_views"`
	AboutViews    int       `json:"about_views"`
	BlogsViews    int       `json:"blogs_views"`
	ProjectsViews int       `json:"projects_views"`
	FeedbackViews int       `json:"feedback_views"`
	UpdatedAt     time.Time `json:"updated_at"`
}
