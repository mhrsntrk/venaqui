.PHONY: build install test clean version help

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || echo "unknown")

LDFLAGS = -X 'main.version=$(VERSION)' -X 'main.commit=$(COMMIT)' -X 'main.date=$(DATE)'

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the binary
	@echo "Building venaqui..."
	@go build -ldflags "$(LDFLAGS)" -o venaqui ./cmd/venaqui

install: ## Install the binary to $GOPATH/bin
	@echo "Installing venaqui..."
	@go install -ldflags "$(LDFLAGS)" ./cmd/venaqui

test: ## Run tests
	@go test ./...

clean: ## Remove build artifacts
	@rm -f venaqui
	@rm -rf dist/

version: ## Show version information
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Date: $(DATE)"

release: ## Build release binaries for multiple platforms
	@echo "Building release binaries..."
	@GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/venaqui-linux-amd64 ./cmd/venaqui
	@GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/venaqui-darwin-amd64 ./cmd/venaqui
	@GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/venaqui-darwin-arm64 ./cmd/venaqui
	@GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/venaqui-windows-amd64.exe ./cmd/venaqui
	@echo "Release binaries built in dist/"
