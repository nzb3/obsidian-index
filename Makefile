# Makefile for obsidian-index

# Variables
BINARY_NAME=obsidian-index
MAIN_PATH=./cmd
BUILD_DIR=build
VERSION?=1.0.0
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X github.com/nzb3/obsidian-index/internal/version.Version=$(VERSION) -X github.com/nzb3/obsidian-index/internal/version.GitCommit=$(GIT_COMMIT) -X github.com/nzb3/obsidian-index/internal/version.BuildDate=$(BUILD_DATE)"

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Build with custom version
.PHONY: build-version
build-version:
	@echo "Building $(BINARY_NAME) with custom version..."
	@mkdir -p $(BUILD_DIR)
	@if [ -z "$(CUSTOM_VERSION)" ]; then \
		echo "Usage: make build-version CUSTOM_VERSION=1.2.3"; \
		exit 1; \
	fi
	@CUSTOM_GIT_COMMIT=$$(git rev-parse --short HEAD 2>/dev/null || echo "unknown"); \
	CUSTOM_BUILD_DATE=$$(date -u +"%Y-%m-%dT%H:%M:%SZ"); \
	go build -ldflags "-X github.com/nzb3/obsidian-index/internal/version.Version=$(CUSTOM_VERSION) -X github.com/nzb3/obsidian-index/internal/version.GitCommit=$$CUSTOM_GIT_COMMIT -X github.com/nzb3/obsidian-index/internal/version.BuildDate=$$CUSTOM_BUILD_DATE" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Build for multiple platforms
.PHONY: build-all
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	
	# macOS
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	
	# Linux
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	
	# Windows
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test ./...

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linter
.PHONY: lint
lint:
	@echo "Running linter..."
	golangci-lint run

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Run the application locally
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Create release archives
.PHONY: release
release: build-all
	@echo "Creating release archives..."
	@cd $(BUILD_DIR) && \
	for binary in *; do \
		if [ -f "$$binary" ]; then \
			tar -czf "$$binary.tar.gz" "$$binary"; \
		fi; \
	done

# Install locally (for testing)
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

# Uninstall
.PHONY: uninstall
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	sudo rm -f /usr/local/bin/$(BINARY_NAME)

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  build-version - Build with custom version (use CUSTOM_VERSION=1.2.3)"
	@echo "  build-all     - Build for all platforms"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run tests"
	@echo "  fmt           - Format code"
	@echo "  lint          - Run linter"
	@echo "  deps          - Install dependencies"
	@echo "  run           - Build and run locally"
	@echo "  release       - Create release archives"
	@echo "  install       - Install locally to /usr/local/bin"
	@echo "  uninstall     - Remove from /usr/local/bin"
	@echo "  help          - Show this help"

