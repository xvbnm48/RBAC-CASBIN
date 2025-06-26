.PHONY: build run seed test clean help

# Variables
APP_NAME := library-api
PORT := 8080

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@go build -o bin/$(APP_NAME) cmd/server/main.go
	@echo "Build completed!"

# Run the application
run:
	@echo "Starting $(APP_NAME) on port $(PORT)..."
	@go run cmd/server/main.go

# Seed the database
seed:
	@echo "Seeding database..."
	@go run cmd/seed/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf bin/
	@go clean

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Create database
createdb:
	@echo "Creating database..."
	@createdb library_db || echo "Database might already exist"

# Setup project (install deps, create db, run migrations, seed data)
setup: deps createdb seed
	@echo "Project setup completed!"

# Development workflow
dev: setup run

# Show help
help:
	@echo "Available commands:"
	@echo "  build    - Build the application"
	@echo "  run      - Run the application"
	@echo "  seed     - Seed the database with initial data"
	@echo "  test     - Run tests"
	@echo "  clean    - Clean build artifacts"
	@echo "  deps     - Download and tidy dependencies"
	@echo "  createdb - Create PostgreSQL database"
	@echo "  setup    - Complete project setup"
	@echo "  dev      - Setup and run for development"
	@echo "  help     - Show this help message"
