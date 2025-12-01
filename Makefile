.PHONY: build install clean test run build-all release help

# Binary name
BINARY_NAME=dusty

# Build flags
LDFLAGS=-ldflags="-s -w"

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build binary for current platform
	go build $(LDFLAGS) -o $(BINARY_NAME)

install: ## Install binary to $GOPATH/bin
	go install $(LDFLAGS)

run: ## Run the application on current directory
	go run main.go

test: ## Run tests
	go test -v ./...

clean: ## Remove built binaries and dist directory
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*
	rm -rf dist/

build-all: ## Build for all platforms (Linux, macOS, Windows)
	@echo "Building for all platforms..."
	@mkdir -p dist
	@echo "Building for Linux amd64..."
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64
	@echo "Building for Linux arm64..."
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64
	@echo "Building for macOS amd64..."
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64
	@echo "Building for macOS arm64..."
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64
	@echo "Building for Windows amd64..."
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe
	@echo "Building for Windows arm64..."
	@GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-arm64.exe
	@echo "Done! Binaries are in dist/"
	@ls -lh dist/

release: clean test build-all ## Run tests, clean, and build for all platforms

fmt: ## Format Go code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

lint: fmt vet ## Run formatting and vetting

deps: ## Download dependencies
	go mod download
	go mod tidy

upgrade-deps: ## Upgrade all dependencies
	go get -u ./...
	go mod tidy
