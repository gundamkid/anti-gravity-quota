.PHONY: build clean test install help release

# Build variables
APP_NAME := ag-quota
VERSION := 0.1.3
BUILD_DIR := dist
INSTALL_DIR := /usr/local/bin

# Build the application for current platform
build:
	@echo "Building $(APP_NAME) v$(VERSION)..."
	@go build -ldflags "-X main.version=$(VERSION)" -o $(APP_NAME) ./cmd/ag-quota
	@echo "✓ Build complete: ./$(APP_NAME)"

# Build for all platforms
release:
	@./build.sh

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(APP_NAME)
	@rm -rf $(BUILD_DIR)
	@echo "✓ Clean complete"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Install to system
install: build
	@echo "Installing $(APP_NAME) to $(INSTALL_DIR)..."
	@sudo cp $(APP_NAME) $(INSTALL_DIR)/
	@echo "✓ Installed to $(INSTALL_DIR)/$(APP_NAME)"

# Uninstall from system
uninstall:
	@echo "Uninstalling $(APP_NAME)..."
	@sudo rm -f $(INSTALL_DIR)/$(APP_NAME)
	@echo "✓ Uninstalled"

# Run the application
run: build
	@./$(APP_NAME)

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✓ Format complete"

# Lint code
lint:
	@echo "Linting code..."
	@$(shell go env GOPATH)/bin/golangci-lint run
	@echo "✓ Lint complete"

# Show help
help:
	@echo "Available targets:"
	@echo "  build     - Build for current platform"
	@echo "  release   - Build for all platforms"
	@echo "  clean     - Remove build artifacts"
	@echo "  test      - Run tests"
	@echo "  install   - Install to system (requires sudo)"
	@echo "  uninstall - Remove from system (requires sudo)"
	@echo "  run       - Build and run"
	@echo "  fmt       - Format code"
	@echo "  lint      - Lint code"
	@echo "  help      - Show this help message"
