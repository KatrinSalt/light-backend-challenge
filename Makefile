.PHONY: build run test clean db-test demo

# Build the application
build:
	go build -o bin/app cmd/cli/main.go

# Run the application
run:
	go run cmd/cli/main.go

# Run all tests
test:
	go test ./...

# Run database tests specifically
db-test:
	go test ./db/sqlite/...

# Run the demo
demo:
	go run cmd/cli/main.go

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run tests with coverage
test-coverage:
	go test -cover ./...

# Run tests with verbose output
test-verbose:
	go test -v ./...




