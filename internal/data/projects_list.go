package data

type ProjectInfo struct {
	Path        string
	Description string
	Tags        string
	ViewsCount  int
	Order       int
}

// OrderedProjectsList maintains the order of projects (reversed order - newest first)
var OrderedProjectsList = []struct {
	Title string
	Info  ProjectInfo
}{
	{
		Title: "Task Management API(using Go)",
		Info: ProjectInfo{
			Path:        "https://github.com/proSamik/go-task-management-api",
			Description: "Golang-based task management REST API",
			Tags:        "golang,api,rest,backend",
			ViewsCount:  1,
			Order:       0,
		},
	},
	{
		Title: "To Do List API with caching(using SpringBoot)",
		Info: ProjectInfo{
			Path:        "https://github.com/proSamik/Spring-Boot-Todo-List-API-with-Caching",
			Description: "Spring Boot-based Todo List API implementation with caching mechanisms",
			Tags:        "java, api,rest,backend",
			ViewsCount:  1,
			Order:       1,
		},
	},
	{
		Title: "Grocery App backend",
		Info: ProjectInfo{
			Path:        "https://github.com/proSamik/grocery-backend",
			Description: "Backend API for a grocery shopping application",
			Tags:        "golang,api,rest,backend",
			ViewsCount:  1,
			Order:       2,
		},
	},
	{
		Title: "Airbnb Analytics",
		Info: ProjectInfo{
			Path:        "https://github.com/proSamik/airbnb-analytics",
			Description: "Data analysis project for Airbnb listings and pricing trends",
			Tags:        "golang,api,rest,backend",
			ViewsCount:  1,
			Order:       3,
		},
	},
	{
		Title: "ProSamik Server",
		Info: ProjectInfo{
			Path:        "https://github.com/proSamik/prosamik-server",
			Description: "Backend server implementation for ProSamik applications",
			Tags:        "golang,api,rest,backend",
			ViewsCount:  1,
			Order:       4,
		},
	},
	{
		Title: "ProSamik Frontend App",
		Info: ProjectInfo{
			Path:        "https://github.com/proSamik/prosamik-frontend-app",
			Description: "Frontend application built with React and TypeScript",
			Tags:        "golang,api,rest,backend",
			ViewsCount:  1,
			Order:       5,
		},
	},
	{
		Title: "About me",
		Info: ProjectInfo{
			Path:        "https://github.com/proSamik/proSamik",
			Description: "My personal portfolio and profile repository",
			Tags:        "golang,api,rest,backend",
			ViewsCount:  1,
			Order:       6,
		},
	},
	{
		Title: "AI Receipt",
		Info: ProjectInfo{
			Path:        "https://github.com/proSamik/AiReceipt",
			Description: "AI-powered receipt scanner and expense tracker application",
			Tags:        "golang,api,rest,backend",
			ViewsCount:  1,
			Order:       7,
		},
	},
	{
		Title: "Smart Parking System",
		Info: ProjectInfo{
			Path:        "https://github.com/proSamik/Smart-Parking-System-using-8051-MCU",
			Description: "An embedded system project implementing smart parking solution using 8051 microcontroller",
			Tags:        "golang,api,rest,backend",
			ViewsCount:  1,
			Order:       8,
		},
	},
	{
		Title: "Direct link",
		Info: ProjectInfo{
			Path:        "https://github.com/proSamik/airbnb-analytics/blob/main/mock_data/README.md",
			Description: "Direct Link of a markdown",
			Tags:        "golang,api,rest,backend",
			ViewsCount:  1,
			Order:       9,
		},
	},
}

// ProjectsList maintains backward compatibility
var ProjectsList = make(map[string]ProjectInfo)

func init() {
	// Initialize the map from the ordered slice
	for _, item := range OrderedProjectsList {
		ProjectsList[item.Title] = ProjectInfo{
			Path:        item.Info.Path,
			Description: item.Info.Description,
		}
	}
}
