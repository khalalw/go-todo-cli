# Makefile

# Variables
BINARY_NAME = todo-cli

# Build the binary
build:
	go build -o $(BINARY_NAME) main.go shell.go commands.go utils.go

# Run the binary
run: build
	./$(BINARY_NAME)

# Clean the binary
clean:
	rm -f $(BINARY_NAME)
