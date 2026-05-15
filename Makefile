BINARY_NAME=obelisk
BUILD_DIR=bin
VERSION=0.1.0
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
BUILD_DATE=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS=-ldflags "-X github.com/Swif7ify/Obelisk-CLI/cmd.Version=$(VERSION) -X github.com/Swif7ify/Obelisk-CLI/cmd.Commit=$(COMMIT) -X github.com/Swif7ify/Obelisk-CLI/cmd.BuildDate=$(BUILD_DATE)"

.PHONY: all build run test clean install lint tidy

all: clean build

build:
	@echo "🏗️  Building $(BINARY_NAME) v$(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "✅ Built $(BUILD_DIR)/$(BINARY_NAME)"

run: build
	@./$(BUILD_DIR)/$(BINARY_NAME) check

test:
	@echo "🧪 Running tests..."
	go test -v -race ./...

clean:
	@echo "🧹 Cleaning..."
	@rm -rf $(BUILD_DIR)
	@go clean

install: build
	@echo "📦 Installing $(BINARY_NAME)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME) 2>/dev/null || \
		cp $(BUILD_DIR)/$(BINARY_NAME) $(HOME)/go/bin/$(BINARY_NAME)
	@echo "✅ Installed"

lint:
	@echo "🔍 Linting..."
	@golangci-lint run ./... 2>/dev/null || echo "Install golangci-lint for linting"

tidy:
	@echo "📦 Tidying modules..."
	go mod tidy
