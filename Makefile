# Makefile for Vine Pod

# Application name
APP_NAME = vine-pod

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS = -ldflags "-X github.com/wepala/vine-pod/pkg/version.Version=$(VERSION) \
                   -X github.com/wepala/vine-pod/pkg/version.Commit=$(COMMIT) \
                   -X github.com/wepala/vine-pod/pkg/version.BuildTime=$(BUILD_TIME)"

# Directories
BIN_DIR = bin
CMD_DIR = cmd/$(APP_NAME)

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BIN_DIR)
	go build $(LDFLAGS) -o $(BIN_DIR)/$(APP_NAME) ./$(CMD_DIR)

# Build for multiple platforms
.PHONY: build-all
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BIN_DIR)/$(APP_NAME)-linux-amd64 ./$(CMD_DIR)
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BIN_DIR)/$(APP_NAME)-linux-arm64 ./$(CMD_DIR)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BIN_DIR)/$(APP_NAME)-darwin-amd64 ./$(CMD_DIR)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BIN_DIR)/$(APP_NAME)-darwin-arm64 ./$(CMD_DIR)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BIN_DIR)/$(APP_NAME)-windows-amd64.exe ./$(CMD_DIR)

# Run the application
.PHONY: run
run: build
	./$(BIN_DIR)/$(APP_NAME)

# Run all tests
.PHONY: test
test:
	go test -v ./...

# Run unit tests only
.PHONY: test-unit
test-unit:
	go test -v ./internal/... ./pkg/...

# Run BDD tests
.PHONY: test-bdd
test-bdd:
	@if [ -d "test/features" ]; then \
		go test -v ./test/; \
	else \
		echo "No BDD features found in test/features/"; \
	fi

# Run integration tests
.PHONY: test-integration
test-integration:
	go test -v -tags=integration ./test/integration/...

# Run tests with coverage
.PHONY: test-cover
test-cover:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Generate mocks using moq
.PHONY: generate-mocks
generate-mocks:
	go generate ./...

# Run linter
.PHONY: lint
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, running go vet and go fmt instead"; \
		go vet ./...; \
		go fmt ./...; \
	fi

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Tidy dependencies
.PHONY: tidy
tidy:
	go mod tidy

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf $(BIN_DIR)
	rm -f coverage.out coverage.html

# Build Docker image
.PHONY: docker-build
docker-build:
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		-t $(APP_NAME):$(VERSION) \
		-t $(APP_NAME):latest \
		.

# Run Docker container
.PHONY: docker-run
docker-run: docker-build
	docker run --rm -p 8080:8080 $(APP_NAME):latest

# Development setup
.PHONY: dev-setup
dev-setup:
	@echo "Setting up development environment..."
	go mod download
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@if ! command -v godog >/dev/null 2>&1; then \
		echo "Installing godog..."; \
		go install github.com/cucumber/godog/cmd/godog@latest; \
	fi

# TDD workflow (format, lint, and test)
.PHONY: dev-test
dev-test: fmt lint test

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build            - Build the application"
	@echo "  build-all        - Build for multiple platforms"
	@echo "  run              - Build and run the application"
	@echo "  test             - Run all tests"
	@echo "  test-unit        - Run unit tests only"
	@echo "  test-bdd         - Run BDD tests"
	@echo "  test-integration - Run integration tests"
	@echo "  test-cover       - Run tests with coverage"
	@echo "  generate-mocks   - Generate mocks using moq"
	@echo "  lint             - Run linter"
	@echo "  fmt              - Format code"
	@echo "  tidy             - Tidy dependencies"
	@echo "  clean            - Clean build artifacts"
	@echo "  docker-build     - Build Docker image"
	@echo "  docker-run       - Build and run Docker container"
	@echo "  dev-setup        - Setup development environment"
	@echo "  dev-test         - TDD workflow (format, lint, test)"
	@echo "  help             - Show this help message"