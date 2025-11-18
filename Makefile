.PHONY: build clean test test-verbose test-coverage test-race install run help

# Build variables
BINARY_NAME=witness
MAIN_PATH=./cmd/witness
BUILD_DIR=.
GO=go

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(BUILD_DIR)/$(BINARY_NAME)
	rm -f *.gif *.mp4
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	$(GO) test ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests (verbose)..."
	$(GO) test -v ./...

# Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test -cover ./...
	@echo ""
	@echo "For detailed coverage report, run:"
	@echo "  go test -coverprofile=coverage.out ./..."
	@echo "  go tool cover -html=coverage.out"

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	$(GO) test -race ./...

# Install to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	$(GO) install $(MAIN_PATH)
	@echo "Install complete"

# Run with default settings
run: build
	./$(BINARY_NAME)

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run ./...

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  clean         - Remove build artifacts and output files"
	@echo "  test          - Run tests"
	@echo "  test-verbose  - Run tests with verbose output"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  test-race     - Run tests with race detection"
	@echo "  install       - Install to GOPATH/bin"
	@echo "  run           - Build and run with default settings"
	@echo "  fmt           - Format code"
	@echo "  lint          - Run linter"
	@echo "  help          - Show this help message"
