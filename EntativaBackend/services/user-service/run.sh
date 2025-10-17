#!/bin/bash

# Socialink User Service Start Script

set -e

echo "🚀 Starting Socialink User Authentication Service..."
echo ""

# Check if .env exists, if not copy from example
if [ ! -f .env ]; then
    echo "📋 Creating .env from .env.example..."
    cp .env.example .env
    echo "⚠️  Please configure your .env file with proper credentials"
    echo ""
fi

# Download dependencies
echo "📦 Downloading Go dependencies..."
go mod download
echo ""

# Run the service
echo "✨ Starting service on port 8001..."
echo "🔗 Health check: http://localhost:8001/health"
echo "📚 API Base URL: http://localhost:8001/api/v1"
echo ""
go run cmd/api/main.go
