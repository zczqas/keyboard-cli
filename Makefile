# Makefile for Keyboard CLI Project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=keyboard-cli
BINARY_UNIX=$(BINARY_NAME)_unix

# Directories
BIN_DIR=./bin/$(BINARY_NAME)
SRC_DIR=.
CMD_DIR=./cmd/keyboard-cli

# Build flags
LDFLAGS=-s -w

# Default target
all: test build


# Build the application
build:
	@mkdir -p bin
	$(GOBUILD) -ldflags="$(LDFLAGS)" -o $(BIN_DIR) $(CMD_DIR)

# Build for different platforms
build-linux:
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="$(LDFLAGS)" -o $(BIN_DIR) $(CMD_DIR)

# Run tests
test:
	$(GOTEST) -v ./...

# Clean up build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Install dependencies
deps:
	$(GOGET) -v ./...
	$(GOGET) -u ./...

# Run the application
run:
	$(GOBUILD) -ldflags="$(LDFLAGS)" -o $(BIN_DIR) $(CMD_DIR)
	./$(BIN_DIR)

# Lint the project
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Vet the code
vet:
	go vet ./...

# Full check (test, lint, vet)
check: test lint vet

# Cross-compilation
cross: build-linux

# Help target
help:
	@echo "Available targets:"
	@echo "  all      - Run tests and build the application"
	@echo "  build    - Build the application"
	@echo "  test     - Run tests"
	@echo "  clean    - Remove build artifacts"
	@echo "  deps     - Install dependencies"
	@echo "  run      - Build and run the application"
	@echo "  lint     - Run golangci-lint"
	@echo "  fmt      - Format the code"
	@echo "  vet      - Vet the code"
	@echo "  cross    - Cross-compile for different platforms"
	@echo "  help     - Show this help message"

.PHONY: all build build-linux test clean deps run lint fmt vet check cross help