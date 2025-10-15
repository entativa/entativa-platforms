# Socialink Post Service

Enterprise-grade posting service for **Entativa's Socialink** (Facebook-like) platform.

## Features

### Core Posting
- ✅ Create posts with text and media attachments
- ✅ Privacy controls (public, friends, only me, custom)
- ✅ Edit and delete posts
- ✅ Location tagging
- ✅ User tagging
- ✅ Feelings and activities
- ✅ Media attachments via gRPC

### Engagement
- ✅ Facebook-style reactions (like, love, haha, wow, sad, angry, care)
- ✅ Nested comments (replies to comments)
- ✅ Comment editing and deletion
- ✅ Share posts with custom caption
- ✅ Real-time engagement counts

### Performance
- ✅ Redis caching (posts, comments, feed)
- ✅ PostgreSQL with optimized indexes
- ✅ Kafka event publishing
- ✅ Cursor-based pagination
- ✅ Connection pooling

---

## Architecture

### Tech Stack
- **Language**: Go 1.21
- **Framework**: Gin
- **Database**: PostgreSQL
- **Cache**: Redis
- **Events**: Kafka
- **Service Comm**: gRPC (to Media Service)

### Design Patterns
- Repository pattern
- Service layer
- Dependency injection
- Event-driven architecture

---

## API Endpoints

### Posts
```
POST   /api/v1/posts                  - Create post
GET    /api/v1/posts/:post_id         - Get post
PUT    /api/v1/posts/:post_id         - Update post
DELETE /api/v1/posts/:post_id         - Delete post
GET    /api/v1/posts/feed             - Get personalized feed
GET    /api/v1/posts/trending         - Get trending posts
GET    /api/v1/posts/user/:user_id    - Get user's posts
```

### Comments
```
POST   /api/v1/posts/:post_id/comments      - Add comment
GET    /api/v1/posts/:post_id/comments      - Get comments
GET    /api/v1/comments/:comment_id/replies - Get replies
PUT    /api/v1/comments/:comment_id         - Update comment
DELETE /api/v1/comments/:comment_id         - Delete comment
```

### Likes/Reactions
```
POST   /api/v1/posts/:post_id/like          - Like post
DELETE /api/v1/posts/:post_id/like          - Unlike post
GET    /api/v1/posts/:post_id/likes         - Get likers
POST   /api/v1/comments/:comment_id/like    - Like comment
DELETE /api/v1/comments/:comment_id/like    - Unlike comment
```

### Shares
```
POST   /api/v1/posts/:post_id/share         - Share post
GET    /api/v1/posts/:post_id/shares        - Get shares
DELETE /api/v1/shares/:share_id             - Delete share
```

---

## Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL 14+
- Redis 6+
- Kafka (optional)

### Setup

```bash
# Install dependencies
go mod download

# Copy environment file
cp .env.example .env

# Edit .env with your configuration
nano .env

# Run migrations
psql -U postgres -c "CREATE DATABASE socialink_posts;"
psql -U postgres -d socialink_posts -f migrations/001_create_posts_table.up.sql
psql -U postgres -d socialink_posts -f migrations/002_create_comments_table.up.sql
psql -U postgres -d socialink_posts -f migrations/003_create_likes_table.up.sql
psql -U postgres -d socialink_posts -f migrations/004_create_shares_table.up.sql

# Run service
go run cmd/api/main.go
```

---

## Usage Examples

### Create Post
```bash
curl -X POST http://localhost:8084/api/v1/posts \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "X-User-ID: $USER_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Hello, Socialink! 🎉",
    "privacy": "public",
    "media_ids": ["media-uuid-1", "media-uuid-2"],
    "location": "San Francisco, CA"
  }'
```

### Get Feed
```bash
curl http://localhost:8084/api/v1/posts/feed \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "X-User-ID: $USER_ID"
```

### Add Comment
```bash
curl -X POST http://localhost:8084/api/v1/posts/$POST_ID/comments \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "X-User-ID: $USER_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Great post! 👍"
  }'
```

### Like Post (with Reaction)
```bash
curl -X POST http://localhost:8084/api/v1/posts/$POST_ID/like \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "X-User-ID: $USER_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "reaction_type": "love"
  }'
```

### Share Post
```bash
curl -X POST http://localhost:8084/api/v1/posts/$POST_ID/share \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "X-User-ID: $USER_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "caption": "Check this out!",
    "privacy": "friends"
  }'
```

---

## Database Schema

### Posts
- ID, User ID, Content
- Media IDs (JSONB array)
- Privacy settings
- Location, Tagged users
- Feelings, Activities
- Engagement counts
- Soft deletion

### Comments
- ID, Post ID, User ID
- Parent ID (for nested replies)
- Content, Optional media
- Likes count
- Soft deletion

### Likes
- ID, User ID
- Post ID or Comment ID
- Reaction type
- Unique constraint per user

### Shares
- ID, User ID, Original Post ID
- Optional caption
- Privacy settings
- Unique constraint per user

---

## Configuration

### Environment Variables
- `PORT` - HTTP server port (default: 8084)
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` - PostgreSQL config
- `REDIS_ADDR` - Redis connection
- `KAFKA_BROKERS` - Kafka brokers
- `MEDIA_SERVICE_GRPC` - Media service gRPC address

---

## Performance

### Caching Strategy
- Posts: 1 hour TTL
- Feed: 10 minutes TTL
- Comments: 30 minutes TTL
- Trending: 5 minutes TTL

### Database Optimization
- 9+ indexes on posts table
- GIN indexes for JSONB arrays
- Full-text search indexes
- Partial indexes for filtered queries
- Foreign key constraints with CASCADE

---

## Events Published

### Kafka Topics
- `post-events` - Post created, updated, deleted
- `post-events` - Comment created, updated, deleted
- `post-events` - Post liked, unliked
- `post-events` - Post shared, unshared

### Event Schema
```json
{
  "event_type": "post.created",
  "post_id": "uuid",
  "user_id": "uuid",
  "privacy": "public",
  "created_at": "2025-10-15T..."
}
```

---

## Integration

### Media Service
- gRPC communication on port 50051
- Automatic thumbnail generation
- Blurhash for progressive loading
- Deduplication

### User Service
- User authentication via JWT
- Profile information enrichment
- Friend relationship checking (for privacy)

---

## Production Deployment

### Docker
```bash
docker build -t socialink-post-service:latest .
docker run -p 8084:8084 \
  -e DATABASE_URL="postgresql://..." \
  socialink-post-service:latest
```

### Kubernetes
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: socialink-post-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: socialink-post-service
  template:
    metadata:
      labels:
        app: socialink-post-service
    spec:
      containers:
      - name: socialink-post-service
        image: socialink-post-service:latest
        ports:
        - containerPort: 8084
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: url
```

---

## Monitoring

### Metrics
- Post creation rate
- Comment creation rate
- Like/reaction distribution
- API latency
- Cache hit ratio
- Database query performance

### Health Check
```bash
curl http://localhost:8084/health
```

---

## Development

### Run Tests
```bash
go test ./...
```

### Format Code
```bash
go fmt ./...
```

### Lint
```bash
golangci-lint run
```

---

Built with ❤️ by **Entativa Engineering**  
**Socialink** - Connecting people, sharing moments
