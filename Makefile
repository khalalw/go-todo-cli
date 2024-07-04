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
.PHONY: all build clean test run help

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

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy

# Help
help:
	@echo "Available targets:"
	@echo "  make build     - Build the application"
	@echo "  make clean     - Clean up built files"
	@echo "  make test      - Run all tests"
	@echo "  make run [args]- Build and run the application with optional arguments"
	@echo "  make deps      - Install dependencies"
	@echo "  make help      - Show this help message"