# Makefile

# Variables
BINARY_NAME = todo-cli
CMD_DIR = ./cmd/todo-cli
INTERNAL_DIRS = ./internal/shell ./internal/commands
PKG_DIR = ./pkg/todo

# Build the binary
build:
	@echo "Building the binary..."
	GO111MODULE=on go build -o $(BINARY_NAME) $(CMD_DIR)/main.go

# Run the binary
run: build
	@echo "Running the binary..."
	./$(BINARY_NAME)

# Run tests
test:
	@echo "Running tests..."
	GO111MODULE=on go test $(INTERNAL_DIRS) $(PKG_DIR) -v

# Clean the binary
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
