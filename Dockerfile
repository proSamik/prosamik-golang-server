# Use the official Go image for building the app (Go 1.22)
FROM golang:1.22 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main ./cmd/server/main.go

# Use a minimal image for running the application
FROM debian:bullseye-slim

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Expose the port the application listens on
EXPOSE 10000

# Run the binary
CMD ["./main"]
