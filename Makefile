.PHONY: build run-cli run-cli-custom test clean deps test-coverage test-verbose help help-cli list-rules create-rule update-rule list-approvers create-approver update-approver process-invoice

# Default target
.DEFAULT_GOAL := help

# Build the CLI application
build:
	@echo "Building CLI application..."
	go build -o bin/backend-challenge-cli .
	@echo "✅ CLI build complete: bin/backend-challenge-cli"

# Run the CLI application
run-cli:
	@echo "Running CLI application..."
	@if [ ! -f "bin/backend-challenge-cli" ]; then \
		echo "❌ CLI not built. Run 'make build' first."; \
		exit 1; \
	fi
	./bin/backend-challenge-cli --company "Light"

# Run the CLI with custom company
run-cli-custom:
	@echo "Running CLI application with custom company..."
	@if [ -z "$(COMPANY)" ]; then \
		echo "❌ Error: COMPANY parameter is required"; \
		echo "Usage: make run-cli-custom COMPANY=\"MyCompany\""; \
		exit 1; \
	fi
	@if [ ! -f "bin/backend-challenge-cli" ]; then \
		echo "❌ CLI not built. Run 'make build' first."; \
		exit 1; \
	fi
	./bin/backend-challenge-cli --company "$(COMPANY)"

# Show CLI help
help-cli:
	@echo "Showing CLI help..."
	@if [ ! -f "bin/backend-challenge-cli" ]; then \
		echo "❌ CLI not built. Run 'make build' first."; \
		exit 1; \
	fi
	./bin/backend-challenge-cli --help

# CLI Commands for workflow rules
list-rules:
	@echo "Listing workflow rules..."
	@if [ ! -f "bin/backend-challenge-cli" ]; then \
		echo "❌ CLI not built. Run 'make build' first."; \
		exit 1; \
	fi
	./bin/backend-challenge-cli --company "Light" list-workflow-rules

create-rule:
	@echo "Creating workflow rule..."
	@if [ -z "$(RULE_ARGS)" ]; then \
		echo "❌ Error: RULE_ARGS parameter is required"; \
		echo "Usage: make create-rule RULE_ARGS=\"--min-amount 100 --max-amount 500 --approver-id 1 --approval-channel 0\""; \
		exit 1; \
	fi
	@if [ ! -f "bin/backend-challenge-cli" ]; then \
		echo "❌ CLI not built. Run 'make build' first."; \
		exit 1; \
	fi
	./bin/backend-challenge-cli --company "Light" create-workflow-rule $(RULE_ARGS)

update-rule:
	@echo "Updating workflow rule..."
	@if [ -z "$(RULE_ARGS)" ]; then \
		echo "❌ Error: RULE_ARGS parameter is required"; \
		echo "Usage: make update-rule RULE_ARGS=\"--id 1 --min-amount 200 --approver-id 2 --approval-channel 1\""; \
		exit 1; \
	fi
	@if [ ! -f "bin/backend-challenge-cli" ]; then \
		echo "❌ CLI not built. Run 'make build' first."; \
		exit 1; \
	fi
	./bin/backend-challenge-cli --company "Light" update-workflow-rule $(RULE_ARGS)

# CLI Commands for approvers
list-approvers:
	@echo "Listing approvers..."
	@if [ ! -f "bin/backend-challenge-cli" ]; then \
		echo "❌ CLI not built. Run 'make build' first."; \
		exit 1; \
	fi
	./bin/backend-challenge-cli --company "Light" list-approvers

create-approver:
	@echo "Creating approver..."
	@if [ -z "$(APPROVER_ARGS)" ]; then \
		echo "❌ Error: APPROVER_ARGS parameter is required"; \
		echo "Usage: make create-approver APPROVER_ARGS=\"--name 'John Doe' --role 'Manager' --email 'john@example.com' --slack-id 'U123456\""; \
		exit 1; \
	fi
	@if [ ! -f "bin/backend-challenge-cli" ]; then \
		echo "❌ CLI not built. Run 'make build' first."; \
		exit 1; \
	fi
	./bin/backend-challenge-cli --company "Light" create-approver $(APPROVER_ARGS)

update-approver:
	@echo "Updating approver..."
	@if [ -z "$(APPROVER_ARGS)" ]; then \
		echo "❌ Error: APPROVER_ARGS parameter is required"; \
		echo "Usage: make update-approver APPROVER_ARGS=\"--id 1 --name 'John Doe' --role 'Manager' --email 'john@example.com' --slack-id 'U123456\""; \
		exit 1; \
	fi
	@if [ ! -f "bin/backend-challenge-cli" ]; then \
		echo "❌ CLI not built. Run 'make build' first."; \
		exit 1; \
	fi
	./bin/backend-challenge-cli --company "Light" update-approver $(APPROVER_ARGS)

# Process invoice
process-invoice:
	@echo "Processing invoice..."
	@if [ ! -f "bin/backend-challenge-cli" ]; then \
		echo "❌ CLI not built. Run 'make build' first."; \
		exit 1; \
	fi
	./bin/backend-challenge-cli --company "Light" process-invoice

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
	rm -f backend-challenge-cli
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
	@echo "Invoice Approval Workflow CLI - Available Commands:"
	@echo ""
	@echo "📦 Build Commands:"
	@echo "  build         - Build the CLI application (creates bin/backend-challenge-cli)"
	@echo ""
	@echo "💻 CLI Commands:"
	@echo "  run-cli       - Run CLI with default company 'Light'"
	@echo "  run-cli-custom- Run CLI with custom company (requires COMPANY parameter)"
	@echo "  help-cli      - Show CLI help and available commands"
	@echo ""
	@echo "📋 Workflow Rules:"
	@echo "  list-rules    - List all workflow rules"
	@echo "  create-rule   - Create a new workflow rule (requires RULE_ARGS parameter)"
	@echo "  update-rule   - Update an existing workflow rule (requires RULE_ARGS parameter)"
	@echo ""
	@echo "👥 Approvers:"
	@echo "  list-approvers- List all approvers"
	@echo "  create-approver- Create a new approver (requires APPROVER_ARGS parameter)"
	@echo "  update-approver- Update an existing approver (requires APPROVER_ARGS parameter)"
	@echo ""
	@echo "📄 Invoice Processing:"
	@echo "  process-invoice- Process an invoice interactively"
	@echo ""
	@echo "🧪 Test Commands:"
	@echo "  test          - Run all tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  test-verbose  - Run tests with verbose output"
	@echo ""
	@echo "🔧 Utility Commands:"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Install/update dependencies"
	@echo "  help          - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make build"
	@echo "  make list-rules"
	@echo "  make create-rule RULE_ARGS=\"--min-amount 100 --max-amount 500 --approver-id 1 --approval-channel 0\""
	@echo "  make update-rule RULE_ARGS=\"--id 1 --min-amount 200 --approver-id 2 --approval-channel 1\""
	@echo "  make list-approvers"
	@echo "  make create-approver APPROVER_ARGS=\"--name 'John Doe' --role 'Manager' --email 'john@example.com' --slack-id 'U123456'\""
	@echo "  make update-approver APPROVER_ARGS=\"--id 1 --name 'Jane Smith' --role 'Senior Manager' --email 'jane@example.com' --slack-id 'U789012'\""
	@echo "  make process-invoice"
	@echo "  make test-coverage"
