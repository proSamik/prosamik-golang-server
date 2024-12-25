package data

type RepoInfo struct {
	Path        string
	Description string
}

var ReposList = map[string]RepoInfo{
	"Demo Template": {
		Path:        "https://github.com/proSamik/demo-template",
		Description: "A demo template to show how content will be rendered in the UI",
	},
	"Direct link": {
		Path:        "https://github.com/proSamik/airbnb-analytics/blob/main/mock_data/README.md",
		Description: "Direct Link of a markdown",
	},
	"Smart Parking System": {
		Path:        "https://github.com/proSamik/Smart-Parking-System-using-8051-MCU",
		Description: "An embedded system project implementing smart parking solution using 8051 microcontroller",
	},
	"AI Receipt": {
		Path:        "https://github.com/proSamik/AiReceipt",
		Description: "AI-powered receipt scanner and expense tracker application",
	},
	"About me": {
		Path:        "https://github.com/proSamik/proSamik",
		Description: "My personal portfolio and profile repository",
	},
	"ProSamik Frontend App": {
		Path:        "https://github.com/proSamik/prosamik-frontend-app",
		Description: "Frontend application built with React and TypeScript",
	},
	"ProSamik Server": {
		Path:        "https://github.com/proSamik/prosamik-server",
		Description: "Backend server implementation for ProSamik applications",
	},
	"Airbnb Analytics": {
		Path:        "https://github.com/proSamik/airbnb-analytics",
		Description: "Data analysis project for Airbnb listings and pricing trends",
	},
	"Grocery App backend": {
		Path:        "https://github.com/proSamik/grocery-backend",
		Description: "Backend API for a grocery shopping application",
	},
	"To Do List API with caching(using SpringBoot)": {
		Path:        "https://github.com/proSamik/Spring-Boot-Todo-List-API-with-Caching",
		Description: "Spring Boot-based Todo List API implementation with caching mechanisms",
	},
	"Task Management API(using Go)": {
		Path:        "https://github.com/proSamik/go-task-management-api",
		Description: "Golang-based task management REST API",
	},
}
