# ProSamik Golang Server

It fetches markdown content from the GitHub repository
and converts it to required HTML and sends the response to the /GET request. 

It also has the dashboard in HTMX which supports adding the list of blogs,
projects and see website analytics for the prosamik nextjs app.

The instructions on how to write the Markdown content that it can understand are [here](https://github.com/proSamik/README-template)

To understand the project in-depth [click here](https://prosamik.com/projects/prosamik-golang-server)

---

## Prerequisites

- Go 1.22.0 (with toolchain 1.23.1)
- Docker (for containerized deployment)
- PostgreSQL
- Git

## Installation

1. Clone the repository:
   ```bash
   git clone git@github.com:proSamik/prosamik-golang-server.git
   cd prosamik-golang-server
   ```

2. Install Go dependencies:
   ```bash
   go mod tidy
   ```

## Environment Setup

1. Create a `.env` file in the root directory with the following configuration:
   ```env
   # SMTP Configuration
   # I have used gmail, you can all aso use your own gmail id and app password of it
   SMTP_HOST=<your-smtp-host>
   SMTP_PORT=<your-smtp-port>
   SMTP_USER=<your-smtp-email>
   SMTP_PASSWORD=<your-smtp-password>
   FEEDBACK_RECIPIENT_EMAIL=<recipient-email>

   # Authentication
   JWT_SECRET_KEY=<your-secret-token>
   ADMIN_PASSWORD=<admin-password>

   # GitHub Integration
   GITHUB_TOKEN=<your-github-token>

   # Database Configuration
   DB_HOST=<your-database-host>
   DB_PORT=<your-database-port>
   DB_USER=<your-database-username>
   DB_PASSWORD=<your-database-pwd>
   DB_NAME=<your-database-name>

   # Application Port
   PORT=10000
   ```

## Running the Application

### Local Development

1. Ensure PostgreSQL is running and accessible with the configured credentials

2. Run the application:
   ```bash
   go run cmd/server/main.go
   ```

### Using Docker

1. Build the Docker image:
   ```bash
   docker build -t <your-image-name> -f Dockerfile .
   ```

2. Run the container:
   ```bash
   docker run -p 10000:10000 --env-file .env <your-image-name>
   ```
---

## Additional Notes

- The application runs on port 10000 by default
- Make sure to set up proper SMTP credentials for email functionality
- Database migrations are automatically handled by the application
- Ensure all environment variables are properly set before running the application

## Suggested Environment Values

- For `JWT_SECRET_KEY`: Use a strong, random string (minimum 32 characters)
- For `ADMIN_PASSWORD`: Use a strong password with mixed case, numbers, and special characters
- For `GITHUB_TOKEN`: Create a personal access token with appropriate permissions
- For `DB_NAME`: Choose a meaningful name for your database

## Troubleshooting

If you encounter any issues:

1. Verify all environment variables are correctly set
2. Ensure PostgreSQL is running and accessible
3. Check if port 10000 is available
4. Verify Go version compatibility

## Contributing

Instructions for contributing:

1. First open the issue
2. Solve the issue and open pull request
3. Assign me as the reviewer

----

