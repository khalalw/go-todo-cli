# Makefile for go-todo-cli

# Variables
BINARY_NAME=todo-cli
BUILD_DIR=.
MAIN_PACKAGE=./cmd/todo

# The arguments passed to make run
ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
# ...and turn them into do-nothing targets
$(eval $(ARGS):;@:)

# Phony targets
.PHONY: all build clean test run help fmt lint vet check deps

# Default target
all: build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)

# Clean up
clean:
	@echo "Cleaning..."
	@rm -f $(BUILD_DIR)/$(BINARY_NAME)
	@go clean

# Run all tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BINARY_NAME) $(ARGS)

# Format the code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint the code
lint:
	@echo "Linting code..."
	@golangci-lint run

# Vet the code
vet:
	@echo "Vetting code..."
	@go vet ./...

# Check formatting, linting, and tests
check: fmt vet lint test

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy

# Help
help:
	@echo "Available targets:"
	@echo "  make build      - Build the application"
	@echo "  make clean      - Clean up built files"
	@echo "  make test       - Run all tests"
	@echo "  make run [args] - Build and run the application with optional arguments"
	@echo "  make fmt        - Format the code"
	@echo "  make lint       - Lint the code"
	@echo "  make vet        - Vet the code"
	@echo "  make check      - Run fmt, vet, lint, and test"
	@echo "  make deps       - Install dependencies"
	@echo "  make help       - Show this help message"
