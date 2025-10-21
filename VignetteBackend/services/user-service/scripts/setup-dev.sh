#!/bin/bash

set -e

echo "🚀 Setting up Vignette User Service for development..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if .env exists
if [ ! -f ".env" ]; then
    echo -e "${YELLOW}Creating .env from .env.example...${NC}"
    cp .env.example .env
    echo -e "${GREEN}✓ Created .env file${NC}"
    echo -e "${YELLOW}⚠️  Please update .env with your actual configuration${NC}"
fi

# Install Go dependencies
echo -e "${YELLOW}Installing Go dependencies...${NC}"
go mod download
go mod tidy
echo -e "${GREEN}✓ Dependencies installed${NC}"

# Check PostgreSQL connection
echo -e "${YELLOW}Checking PostgreSQL connection...${NC}"
source .env
export PGPASSWORD=$DB_PASSWORD

if psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c '\l' > /dev/null 2>&1; then
    echo -e "${GREEN}✓ PostgreSQL connection successful${NC}"
else
    echo -e "${RED}✗ PostgreSQL connection failed${NC}"
    echo -e "${YELLOW}Please ensure PostgreSQL is running and credentials are correct${NC}"
    exit 1
fi

# Create database if it doesn't exist
echo -e "${YELLOW}Creating database if not exists...${NC}"
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c "CREATE DATABASE $DB_NAME;" 2>/dev/null || echo -e "${GREEN}✓ Database already exists${NC}"

# Run migrations
echo -e "${YELLOW}Running database migrations...${NC}"
for migration in migrations/*.sql; do
    echo "  Running $(basename $migration)..."
    psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f $migration
done
echo -e "${GREEN}✓ Migrations complete${NC}"

# Build the application
echo -e "${YELLOW}Building application...${NC}"
go build -o bin/user-service cmd/api/main.go cmd/api/routes.go
echo -e "${GREEN}✓ Build complete${NC}"

echo ""
echo -e "${GREEN}🎉 Setup complete!${NC}"
echo ""
echo "To start the server:"
echo "  make run"
echo ""
echo "To run migrations manually:"
echo "  make migrate-up"
echo ""
echo "To test the API:"
echo "  curl http://localhost:8002/health"
echo ""
