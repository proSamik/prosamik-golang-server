# Build stage
FROM golang:1.22 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code, templates, and migrations
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/main.go

# Final stage
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the compiled binary, templates, and migrations from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/internal/templates ./internal/templates
COPY --from=builder /app/internal/database/migrations ./internal/database/migrations

# Expose the application port
EXPOSE 10000

# Command to run the application
CMD ["./main"]
