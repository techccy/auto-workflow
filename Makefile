.PHONY: build clean run test install help

# Binary name
BINARY_NAME=auto-workflow

# Build directory
BUILD_DIR=build

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) cmd/main.go
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -rf output
	@rm -rf screenshots
	@rm -rf advanced_output
	@echo "Clean complete"

run-simple: build ## Run simple workflow example
	@echo "Running simple workflow..."
	$(BUILD_DIR)/$(BINARY_NAME) -workflow examples/simple_workflow.json

run-adb: build ## Run ADB workflow example
	@echo "Running ADB workflow..."
	$(BUILD_DIR)/$(BINARY_NAME) -workflow examples/adb_workflow.json

run-web: build ## Run web workflow example
	@echo "Running web workflow..."
	$(BUILD_DIR)/$(BINARY_NAME) -workflow examples/web_workflow.json

run-advanced: build ## Run advanced workflow example
	@echo "Running advanced workflow..."
	$(BUILD_DIR)/$(BINARY_NAME) -workflow examples/advanced_workflow.json

list-devices: build ## List connected ADB devices
	@echo "Listing connected ADB devices..."
	$(BUILD_DIR)/$(BINARY_NAME) -list-devices

test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v ./...

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

fmt: ## Format code
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	$(GOCMD) vet ./...

lint: fmt vet ## Run linting (format and vet)

all: clean deps build ## Clean, download dependencies, and build

install: build ## Install the binary to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	$(GOCMD) install cmd/main.go
	@echo "Installation complete"
