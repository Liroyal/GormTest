#!/bin/bash

# Setup script for GormTest development environment
# This script sets up the development environment for the project

set -e

echo "Setting up GormTest development environment..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Please install Go first."
    exit 1
fi

# Check if PostgreSQL is installed
if ! command -v psql &> /dev/null; then
    echo "PostgreSQL is not installed. Please install PostgreSQL first."
    exit 1
fi

# Install dependencies
echo "Installing Go dependencies..."
go mod tidy
go mod download

# Set up environment variables
if [ ! -f .env ]; then
    echo "Creating .env file..."
    cat > .env << EOF
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=
DB_NAME=gormtest
DB_SSLMODE=disable
SERVER_HOST=localhost
SERVER_PORT=8080
EOF
    echo ".env file created. Please update it with your database credentials."
fi

# Create database if it doesn't exist
echo "Setting up database..."
createdb gormtest 2>/dev/null || echo "Database 'gormtest' may already exist"

echo "Setup complete! You can now run 'make run' to start the application."
