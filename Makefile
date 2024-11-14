# Colors for help message
BLUE = \033[34m
NC = \033[0m
SHELL = /bin/bash

# Output directory
OUTPUT := ./output
OUTPUT_BIN := yanji
LDFLAGS := $(shell source ./scripts/make/version.sh && generate_ldflags)

# Go build flags
GO_BUILD_FLAGS := -ldflags "$(LDFLAGS)"

# Targets
.PHONY: help build run test clean lint fmt vet

help:  ## Show help message.
	@printf "Usage:\n"
	@printf "  make $(BLUE)<target>$(NC)\n\n"
	@printf "Targets:\n"
	@perl -nle'print $& if m{^[a-zA-Z0-9_-]+:.*?## .*$$}' $(MAKEFILE_LIST) | \
		sort | \
		awk 'BEGIN {FS = ":.*?## "}; \
		{printf "$(BLUE)  %-18s$(NC) %s\n", $$1, $$2}'

build:  ## Build the Go server.
	@echo "Building the Go server..."
	go generate ./...
	go build $(GO_BUILD_FLAGS) -o $(OUTPUT)/$(OUTPUT_BIN) ./cmd/app

run:  ## Run the Go server.
	@echo "Running the Go server..."
	go run $(GO_BUILD_FLAGS) ./cmd/app

test:  ## Run tests.
	@echo "Running tests..."
	go test -v ./...

clean:  ## Clean the project.
	@echo "Cleaning the project..."
	rm -rf $(OUTPUT)

lint:  ## Run golangci-lint.
	@echo "Running golangci-lint..."
	golangci-lint run ./...

fmt:  ## Format Go code.
	@echo "Formatting Go code..."
	go fmt ./...

vet:  ## Run go vet.
	@echo "Running go vet..."
	go vet ./...