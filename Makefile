.PHONY: build run test clean

# Build variables
BINARY_NAME=shreadbox
MAIN_FILE=cmd/api/main.go

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_FILE)

# Run the application
run:
	$(GORUN) $(MAIN_FILE)

# Run tests
test:
	$(GOTEST) -v ./...

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf storage/*

# Create necessary directories
setup:
	mkdir -p storage
	mkdir -p web/static
	cp -n .env.example .env || true

# Install dependencies
deps:
	$(GOCMD) mod download
	$(GOCMD) mod tidy

# Run all setup tasks
init: deps setup

# Development tasks
dev: init run 