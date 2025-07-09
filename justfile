# Show available recipes
default:
    @just --list

# Run tests
test:
    go test -v ./...

# Run tests with race detector
test-race:
    go test -v -race ./...

# Run tests with coverage
test-coverage:
    go test -v -race -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html

# Run linter
lint:
    golangci-lint run --config=.golangci.yml

# Format code
fmt:
    go fmt ./...
    goimports -w .

# Build all packages
build:
    go build -v ./...

# Clean build artifacts
clean:
    go clean ./...
    rm -f coverage.out coverage.html

# Download dependencies
deps:
    go mod download

# Tidy dependencies
tidy:
    go mod tidy

# Run all checks (used in CI)
check: fmt lint test-race
