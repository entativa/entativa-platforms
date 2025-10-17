#!/bin/bash

set -e

echo "🔄 Running database migrations for all services..."

# Wait for PostgreSQL to be ready
echo "⏳ Waiting for PostgreSQL..."
until docker exec socialink-postgres pg_isready -U postgres > /dev/null 2>&1; do
  sleep 1
done
echo "✅ PostgreSQL is ready!"

# Run migrations for each service
echo ""
echo "📦 User Service migrations..."
docker exec socialink-postgres psql -U postgres -d socialink_users -f /migrations/user-service/*.sql || true

echo "📦 Post Service migrations..."
docker exec socialink-postgres psql -U postgres -d socialink_posts -f /migrations/posting-service/*.sql || true

echo "📦 Messaging Service migrations..."
docker exec socialink-postgres psql -U postgres -d socialink_messaging -f /migrations/messaging-service/*.sql || true

echo "📦 Settings Service migrations..."
docker exec socialink-postgres psql -U postgres -d socialink_settings -f /migrations/settings-service/*.sql || true

echo "📦 Notification Service migrations..."
docker exec socialink-postgres psql -U postgres -d socialink_notifications -f /migrations/notification-service/*.sql || true

echo "📦 Community Service migrations..."
docker exec socialink-postgres psql -U postgres -d socialink_communities -f /migrations/community-service/*.sql || true

echo "📦 Streaming Service migrations..."
docker exec socialink-postgres psql -U postgres -d socialink_streaming -f /migrations/live-streaming-service/*.sql || true

echo "📦 Creator Service migrations..."
docker exec socialink-postgres psql -U postgres -d socialink_creator -f /migrations/creator-service/*.sql || true

echo ""
echo "✅ All migrations completed successfully!"
