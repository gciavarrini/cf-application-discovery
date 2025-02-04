.DEFAULT_GOAL := build

GINKGO := ginkgo
COVERPROFILE := coverage.out
COVERPKG := ./...

.PHONY: fmt vet build lint test coverage-html clean help

# Format Go code
fmt:
	@echo "Formatting Go code..."
	go fmt ./...

# Run static analysis and vetting
vet: fmt
	@echo "Running go vet..."
	go vet ./...

build: vet
	@echo "Building the application..."
	go build

lint:
	@echo "Running staticcheck..."
	staticcheck ./...
	go vet ./...

# Run tests with coverage
test:
	@echo "Running tests with coverage..."
	$(GINKGO) -v --cover --coverprofile=$(COVERPROFILE) --coverpkg=$(COVERPKG) ./...

# Generate HTML coverage report from tests
coverage-html: test
	@echo "Generating HTML coverage report..."
	go tool cover -html=$(COVERPROFILE) -o coverage.html

# Clean up generated files
clean:
	@echo "Cleaning up coverage files..."
	rm -f $(COVERPROFILE) coverage.html

# Help target to show available commands
help:
	@echo "Available commands:"
	@echo "  make test            - Run tests with coverage"
	@echo "  make coverage-html   - Generate HTML coverage report"
	@echo "  make clean           - Remove coverage files"
	@echo "  make help            - Show this help message"
