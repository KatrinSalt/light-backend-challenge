.PHONY: build run test clean deps test-coverage test-verbose help

# Default target
.DEFAULT_GOAL := help

# Build the application
build:
	@echo "Building invoice approval workflow application..."
	go build -o bin/app .
	@echo "✅ Build complete: bin/app"

# Run the application with default settings
run:
	@echo "Running invoice approval workflow with default company 'Light'..."
	go run . -company "Light"

# Run the application with custom settings
# Usage: make run-custom COMPANY="MyCompany" DEPARTMENTS="Finance,Marketing,HR"
run-custom:
	@echo "Running invoice approval workflow with custom settings..."
	@if [ -z "$(COMPANY)" ]; then \
		echo "❌ Error: COMPANY parameter is required"; \
		echo "Usage: make run-custom COMPANY=\"MyCompany\" DEPARTMENTS=\"Finance,Marketing\""; \
		exit 1; \
	fi
	@if [ -n "$(DEPARTMENTS)" ]; then \
		go run . -company "$(COMPANY)" -departments "$(DEPARTMENTS)"; \
	else \
		go run . -company "$(COMPANY)"; \
	fi

# Run the application with any flags
# Usage: make run-with FLAGS="-company MyCompany -departments Finance,Marketing"
run-with:
	@echo "Running invoice approval workflow with custom flags..."
	@if [ -z "$(FLAGS)" ]; then \
		echo "❌ Error: FLAGS parameter is required"; \
		echo "Usage: make run-with FLAGS=\"-company MyCompany -departments Finance,Marketing\""; \
		exit 1; \
	fi
	go run . $(FLAGS)

# Run the application with help
help-app:
	@echo "Showing application help..."
	go run . --help

# Run all tests
test:
	@echo "Running all tests..."
	go test ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -cover ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean
	@echo "✅ Clean complete"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download
	@echo "✅ Dependencies installed"

# Show this help message
help:
	@echo "Invoice Approval Workflow - Available Commands:"
	@echo ""
	@echo "  build         - Build the application (creates bin/app)"
	@echo "  run           - Run with default company 'Light'"
	@echo "  run-custom    - Run with custom company and departments (requires COMPANY parameter)"
	@echo "  run-with      - Run with any custom flags (requires FLAGS parameter)"
	@echo "  help-app      - Show application help and available flags"
	@echo "  test          - Run all tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  test-verbose  - Run tests with verbose output"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Install/update dependencies"
	@echo "  help          - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make build"
	@echo "  make run"
	@echo "  make run-custom COMPANY=\"MyCompany\" DEPARTMENTS=\"Finance,Marketing,HR\""
	@echo "  make run-with FLAGS=\"-company MyCompany -departments Finance,Marketing\""
	@echo "  make run-with FLAGS=\"-company MyCompany -slack-connection-string my-slack -email-connection-string my-email\""
	@echo "  make test-coverage"
