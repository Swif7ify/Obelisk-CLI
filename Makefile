# Obelisk CLI Makefile

# Variables
BINARY_NAME=obelisk
VERSION?=0.1.0
BUILD_DIR=bin
GO=go
GOFLAGS=-ldflags="-s -w -X main.Version=$(VERSION)"

# Colors for output
CYAN=\033[0;36m
GREEN=\033[0;32m
YELLOW=\033[0;33m
RED=\033[0;31m
NC=\033[0m # No Color

.PHONY: all build test clean install uninstall lint fmt vet coverage help

# Default target
all: clean test build

## help: Display this help message
help:
	@echo "$(CYAN)Obelisk CLI - Makefile Commands$(NC)"
	@echo ""
	@echo "$(GREEN)Building:$(NC)"
	@echo "  make build          - Build the binary for current platform"
	@echo "  make build-all      - Build for all platforms"
	@echo "  make install        - Install binary to system"
	@echo "  make uninstall      - Remove installed binary"
	@echo ""
	@echo "$(GREEN)Testing:$(NC)"
	@echo "  make test           - Run all tests"
	@echo "  make test-verbose   - Run tests with verbose output"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make test-race      - Run tests with race detector"
	@echo "  make bench          - Run benchmarks"
	@echo ""
	@echo "$(GREEN)Code Quality:$(NC)"
	@echo "  make lint           - Run linter"
	@echo "  make fmt            - Format code"
	@echo "  make vet            - Run go vet"
	@echo "  make check          - Run all checks (fmt, vet, lint, test)"
	@echo ""
	@echo "$(GREEN)Maintenance:$(NC)"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make deps           - Download dependencies"
	@echo "  make tidy           - Tidy go.mod"
	@echo ""
	@echo "$(GREEN)Distribution:$(NC)"
	@echo "  make installer      - Build Windows MSI installer"
	@echo "  make release        - Create release artifacts"

## build: Build the binary for current platform
build:
	@echo "$(CYAN)Building $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "$(GREEN)✓ Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

## build-all: Build for all platforms
build-all:
	@echo "$(CYAN)Building for all platforms...$(NC)"
	@mkdir -p $(BUILD_DIR)
	
	@echo "$(YELLOW)Building for Windows (amd64)...$(NC)"
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	
	@echo "$(YELLOW)Building for Linux (amd64)...$(NC)"
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	
	@echo "$(YELLOW)Building for macOS (amd64)...$(NC)"
	GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	
	@echo "$(YELLOW)Building for macOS (arm64)...$(NC)"
	GOOS=darwin GOARCH=arm64 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	
	@echo "$(GREEN)✓ All builds complete$(NC)"

## test: Run all tests
test:
	@echo "$(CYAN)Running tests...$(NC)"
	$(GO) test -v ./...
	@echo "$(GREEN)✓ Tests passed$(NC)"

## test-verbose: Run tests with verbose output
test-verbose:
	@echo "$(CYAN)Running tests (verbose)...$(NC)"
	$(GO) test -v -count=1 ./...

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "$(CYAN)Running tests with coverage...$(NC)"
	@mkdir -p coverage
	$(GO) test -v -coverprofile=coverage/coverage.out -covermode=atomic ./...
	$(GO) tool cover -html=coverage/coverage.out -o coverage/coverage.html
	@echo "$(GREEN)✓ Coverage report: coverage/coverage.html$(NC)"

## test-race: Run tests with race detector
test-race:
	@echo "$(CYAN)Running tests with race detector...$(NC)"
	$(GO) test -v -race ./...
	@echo "$(GREEN)✓ Race tests passed$(NC)"

## bench: Run benchmarks
bench:
	@echo "$(CYAN)Running benchmarks...$(NC)"
	$(GO) test -bench=. -benchmem ./...

## lint: Run linter
lint:
	@echo "$(CYAN)Running linter...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout=5m; \
		echo "$(GREEN)✓ Linting passed$(NC)"; \
	else \
		echo "$(YELLOW)⚠ golangci-lint not installed. Install: https://golangci-lint.run/usage/install/$(NC)"; \
	fi

## fmt: Format code
fmt:
	@echo "$(CYAN)Formatting code...$(NC)"
	$(GO) fmt ./...
	@echo "$(GREEN)✓ Code formatted$(NC)"

## vet: Run go vet
vet:
	@echo "$(CYAN)Running go vet...$(NC)"
	$(GO) vet ./...
	@echo "$(GREEN)✓ Vet passed$(NC)"

## check: Run all checks
check: fmt vet lint test
	@echo "$(GREEN)✓ All checks passed$(NC)"

## clean: Remove build artifacts
clean:
	@echo "$(CYAN)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -rf coverage
	@rm -f coverage.out coverage.html
	@echo "$(GREEN)✓ Clean complete$(NC)"

## deps: Download dependencies
deps:
	@echo "$(CYAN)Downloading dependencies...$(NC)"
	$(GO) mod download
	@echo "$(GREEN)✓ Dependencies downloaded$(NC)"

## tidy: Tidy go.mod
tidy:
	@echo "$(CYAN)Tidying go.mod...$(NC)"
	$(GO) mod tidy
	@echo "$(GREEN)✓ go.mod tidied$(NC)"

## install: Install binary to system
install: build
	@echo "$(CYAN)Installing $(BINARY_NAME)...$(NC)"
	@if [ "$(shell uname)" = "Windows_NT" ]; then \
		powershell -ExecutionPolicy Bypass -File install.ps1; \
	else \
		sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/; \
		sudo chmod +x /usr/local/bin/$(BINARY_NAME); \
	fi
	@echo "$(GREEN)✓ Installed to system$(NC)"

## uninstall: Remove installed binary
uninstall:
	@echo "$(CYAN)Uninstalling $(BINARY_NAME)...$(NC)"
	@if [ "$(shell uname)" = "Windows_NT" ]; then \
		powershell -ExecutionPolicy Bypass -File uninstall.ps1; \
	else \
		sudo rm -f /usr/local/bin/$(BINARY_NAME); \
	fi
	@echo "$(GREEN)✓ Uninstalled from system$(NC)"

## installer: Build Windows MSI installer
installer: build
	@echo "$(CYAN)Building Windows MSI installer...$(NC)"
	@if [ "$(shell uname)" = "Windows_NT" ]; then \
		powershell -ExecutionPolicy Bypass -File installer/build-installer.ps1; \
	else \
		echo "$(RED)✗ MSI installer can only be built on Windows$(NC)"; \
		exit 1; \
	fi

## release: Create release artifacts
release: clean build-all
	@echo "$(CYAN)Creating release artifacts...$(NC)"
	@mkdir -p release
	
	@echo "$(YELLOW)Calculating checksums...$(NC)"
	@cd $(BUILD_DIR) && sha256sum * > ../release/checksums.txt
	
	@echo "$(YELLOW)Copying binaries...$(NC)"
	@cp $(BUILD_DIR)/* release/
	
	@echo "$(GREEN)✓ Release artifacts ready in release/$(NC)"

# Development helpers
.PHONY: dev run watch

## dev: Run in development mode
dev: build
	@./$(BUILD_DIR)/$(BINARY_NAME)

## run: Run without building
run:
	@$(GO) run .

## watch: Watch for changes and rebuild (requires entr)
watch:
	@if command -v entr >/dev/null 2>&1; then \
		find . -name '*.go' | entr -r make run; \
	else \
		echo "$(RED)✗ entr not installed. Install: http://eradman.com/entrproject/$(NC)"; \
	fi

# Made with Bob
