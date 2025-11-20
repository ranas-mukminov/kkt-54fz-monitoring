.PHONY: help build test test-integration lint security-check perf-check check clean deps run docker-build docker-up docker-down

# Variables
BINARY_NAME=kkt-monitor
MAIN_PATH=./cmd/kkt-monitor
BUILD_DIR=./build
GO=go
GOTEST=$(GO) test
GOLINT=golangci-lint
DOCKER_COMPOSE=docker-compose

# Build info
VERSION?=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse HEAD)
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deps: ## Install dependencies
	$(GO) mod download
	$(GO) mod tidy

build: deps ## Build the application
	@echo "Building $(BINARY_NAME)..."
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

run: build ## Build and run the application
	$(BUILD_DIR)/$(BINARY_NAME) --config configs/config.yaml

test: ## Run unit tests
	@echo "Running unit tests..."
	$(GOTEST) -v -race -coverprofile=coverage.txt -covermode=atomic ./...

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	$(GOTEST) -v -tags=integration ./test/integration/...

test-coverage: test ## Run tests and generate coverage report
	$(GO) tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run linter
	@echo "Running linter..."
	@if ! command -v $(GOLINT) > /dev/null; then \
		echo "golangci-lint not found. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	$(GOLINT) run --timeout=5m ./...

security-check: ## Run security checks
	@echo "Running security checks..."
	@if ! command -v gosec > /dev/null; then \
		echo "gosec not found. Installing..."; \
		go install github.com/securego/gosec/v2/cmd/gosec@latest; \
	fi
	gosec -fmt=json -out=security-report.json ./...
	@echo "Security report generated: security-report.json"

perf-check: ## Run performance benchmarks
	@echo "Running performance benchmarks..."
	$(GOTEST) -bench=. -benchmem -cpuprofile=cpu.prof -memprofile=mem.prof ./...

check: lint test security-check ## Run all checks (lint, test, security)

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.txt coverage.html
	rm -f security-report.json
	rm -f cpu.prof mem.prof
	$(GO) clean

docker-build: ## Build Docker image
	docker build -t $(BINARY_NAME):$(VERSION) -f deployments/docker/Dockerfile .

docker-up: ## Start services with docker-compose
	$(DOCKER_COMPOSE) -f deployments/docker/docker-compose.yaml up -d

docker-down: ## Stop services with docker-compose
	$(DOCKER_COMPOSE) -f deployments/docker/docker-compose.yaml down

fmt: ## Format code
	$(GO) fmt ./...
	goimports -w .

vet: ## Run go vet
	$(GO) vet ./...

mod-update: ## Update dependencies
	$(GO) get -u ./...
	$(GO) mod tidy

install: build ## Install the binary
	$(GO) install $(LDFLAGS) $(MAIN_PATH)

.DEFAULT_GOAL := help
