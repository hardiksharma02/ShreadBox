#!/bin/bash

# Make script exit on first error
set -e

# Print commands before executing them
set -x

# Create necessary directories
mkdir -p storage
mkdir -p web/static
mkdir -p web/templates

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    cat > .env << EOL
# Server Configuration
PORT=8080
GIN_MODE=debug  # Set to 'release' in production

# File Settings
MAX_FILE_SIZE=10  # Maximum file size in MB
STORAGE_PATH=./storage
CLEANUP_INTERVAL=5m  # Format: 1h, 5m, 30s, etc.

# Rate Limiting
RATE_LIMIT=100  # Requests per minute
RATE_BURST=5    # Maximum burst size
EOL
fi

# Initialize Go module and get dependencies
go mod tidy
go mod download

# Build the application
go build -o shreadbox cmd/api/main.go

# Make the binary executable
chmod +x shreadbox

echo "Setup completed successfully!"
echo "You can now run the application with: ./shreadbox"
echo "Or use 'make run' to run it through the Makefile" 