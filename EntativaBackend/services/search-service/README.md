# Socialink Search Service 🔍

**Production-grade search service** powered by Elasticsearch with real-time autocomplete, trending hashtags, and advanced filters.

---

## 🔥 Features

### Core Search
- ⚡ **Multi-entity search** - Search across users, posts, Takes, hashtags, and locations
- 🎯 **Field boosting** - Username^3, display_name^2, caption^2 for relevance
- 🔤 **Fuzzy matching** - Typo tolerance with AUTO fuzziness
- 📊 **Advanced filters** - Date range, engagement metrics, media type, location
- 💾 **Redis caching** - 5-15 minute TTLs for performance

### Autocomplete
- ⚡ **Real-time suggestions** - Sub-100ms completion suggester
- 🎯 **Popular searches** - Include trending queries
- 📝 **Recent searches** - User's search history
- 🚀 **Fuzzy matching** - Handle typos

### Hashtags
- 📈 **Trending hashtags** - With growth rates and rankings
- 🔗 **Related hashtags** - More Like This algorithm
- 🔍 **Hashtag search** - Fuzzy matching
- 📊 **Usage statistics** - Post count, Take count, growth rate

### Analytics
- 📊 **Search history** - Track user searches
- 📈 **Trending searches** - Real-time trending queries
- 💡 **Smart suggestions** - Suggest alternatives for low results

---

## 🏗️ Architecture

```
Search Service (Go + Elasticsearch + Redis)
├── Gin (REST API)
├── Elasticsearch
│   ├── users (3 shards)
│   ├── posts (5 shards)
│   ├── takes (5 shards)
│   ├── hashtags (3 shards)
│   └── locations (3 shards)
├── Redis (Caching + Analytics)
│   ├── Search results (5 min)
│   ├── Autocomplete (15 min)
│   ├── Trending searches
│   └── Search history
└── gRPC (Optional - for indexing events)
```

---

## 📡 API Endpoints

### Search
```
GET  /api/v1/search              - Multi-entity search
GET  /api/v1/search/users        - Search users
GET  /api/v1/search/posts        - Search posts
GET  /api/v1/search/takes        - Search Takes
GET  /api/v1/search/history      - Get search history
DEL  /api/v1/search/history      - Clear search history
GET  /api/v1/search/trending     - Trending searches
```

### Autocomplete
```
GET /api/v1/autocomplete         - Get suggestions
GET /api/v1/autocomplete/recent  - Recent searches
```

### Hashtags
```
GET /api/v1/hashtags/trending         - Trending hashtags
GET /api/v1/hashtags/:tag/related     - Related hashtags
GET /api/v1/hashtags/search           - Search hashtags
```

### Indexing (Internal/Admin)
```
POST   /api/v1/index/document    - Index document
POST   /api/v1/index/bulk        - Bulk index
PUT    /api/v1/index/document    - Update document
DELETE /api/v1/index/document    - Delete document
POST   /api/v1/index/reindex     - Reindex all
GET    /api/v1/index/stats       - Index statistics
```

---

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Elasticsearch 8.x
- Redis 6.0+

### Installation

```bash
cd SocialinkBackend/services/search-service

# Install dependencies
go mod download

# Copy environment file
cp .env.example .env

# Edit configuration
nano .env
```

### Configuration

```env
PORT=8087
ELASTICSEARCH_URL=http://localhost:9200
REDIS_URL=localhost:6379
```

### Run

```bash
# Development
go run cmd/api/main.go

# Production
go build -o search-service cmd/api/main.go
./search-service
```

### Docker

```bash
# Build
docker build -t socialink-search-service .

# Run
docker run -d \
  -p 8087:8087 \
  -e ELASTICSEARCH_URL=http://elasticsearch:9200 \
  -e REDIS_URL=redis:6379 \
  socialink-search-service
```

---

## 💡 Usage Examples

### Multi-Entity Search

```bash
curl -X GET "http://localhost:8087/api/v1/search?query=john&type=all&limit=20"
```

Response:
```json
{
  "query": "john",
  "type": "all",
  "total_hits": 247,
  "results": [
    {
      "id": "user-uuid",
      "type": "user",
      "score": 8.5,
      "data": {
        "username": "john_doe",
        "display_name": "John Doe",
        "verified": true
      },
      "snippet": "John Doe - Creator & Photographer"
    },
    {
      "id": "post-uuid",
      "type": "post",
      "score": 7.2,
      "data": {
        "caption": "Amazing sunset with @john!",
        "likes_count": 1234
      }
    }
  ],
  "took_ms": 45
}
```

### User Search with Filters

```bash
curl -X GET "http://localhost:8087/api/v1/search/users?query=photographer&verified=true&min_followers=1000"
```

### Autocomplete

```bash
curl -X GET "http://localhost:8087/api/v1/autocomplete?query=joh&type=user&limit=10"
```

Response:
```json
{
  "query": "joh",
  "suggestions": [
    {
      "text": "john_doe",
      "type": "user",
      "score": 9.5,
      "metadata": {
        "verified": true,
        "follower_count": 12500
      }
    },
    {
      "text": "johnny_test",
      "type": "user",
      "score": 8.2
    }
  ],
  "took_ms": 12
}
```

### Trending Hashtags

```bash
curl -X GET "http://localhost:8087/api/v1/hashtags/trending?limit=20"
```

Response:
```json
{
  "hashtags": [
    {
      "tag": "dancechallenge",
      "display_tag": "#DanceChallenge",
      "usage_count": 15234,
      "post_count": 8932,
      "take_count": 6302,
      "growth_rate": 85.3,
      "rank": 1,
      "is_trending": true
    }
  ],
  "updated_at": "2025-10-15T10:30:00Z"
}
```

### Related Hashtags

```bash
curl -X GET "http://localhost:8087/api/v1/hashtags/dance/related?limit=10"
```

### Indexing a Document

```bash
curl -X POST "http://localhost:8087/api/v1/index/document" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "index",
    "document_type": "user",
    "document_id": "user-uuid",
    "data": {
      "username": "john_doe",
      "display_name": "John Doe",
      "bio": "Creator & Photographer",
      "verified": true,
      "follower_count": 12500
    }
  }'
```

---

## 📊 Elasticsearch Indices

### Users Index
```
Shards: 3
Mappings:
  - username (text + keyword + completion)
  - display_name (text + keyword)
  - bio (text)
  - verified (boolean)
  - follower_count (integer)
```

### Posts Index
```
Shards: 5 (high volume)
Mappings:
  - caption (text, analyzed)
  - content (text, analyzed)
  - hashtags (keyword[])
  - media_type (keyword)
  - likes_count (integer)
  - created_at (date)
```

### Takes Index
```
Shards: 5 (high volume)
Mappings:
  - caption (text)
  - hashtags (keyword[])
  - views_count (integer)
  - likes_count (integer)
  - remix_count (integer)
```

### Hashtags Index
```
Shards: 3
Mappings:
  - tag (text + keyword + completion)
  - usage_count (long)
  - growth_rate (float)
  - is_trending (boolean)
```

### Locations Index
```
Shards: 3
Mappings:
  - name (text + keyword + completion)
  - coordinates (geo_point)
  - city (text)
  - country (keyword)
```

---

## 📈 Performance

### Caching Strategy
- **Search results**: 5 minutes
- **Autocomplete**: 15 minutes
- **Trending hashtags**: 5 minutes
- **Related hashtags**: 10 minutes

### Query Performance
- **Search**: Sub-100ms typical
- **Autocomplete**: Sub-50ms typical
- **Trending**: Sub-20ms (cached)

### Elasticsearch Optimization
- **Field boosting**: username^3, display_name^2
- **Minimum should match**: 75%
- **Fuzziness**: AUTO (typo tolerance)
- **Completion suggesters**: Optimized for prefix matching

---

## 🔧 Development

### Project Structure
```
search-service/
├── cmd/api/main.go          (Entry point)
├── internal/
│   ├── elasticsearch/
│   │   ├── client.go        (ES client wrapper)
│   │   ├── indices.go       (Index mappings)
│   │   └── queries.go       (Query builders)
│   ├── service/
│   │   ├── search_service.go
│   │   ├── autocomplete_service.go
│   │   ├── indexing_service.go
│   │   └── hashtag_service.go
│   ├── handler/
│   │   ├── search_handler.go
│   │   ├── autocomplete_handler.go
│   │   ├── hashtag_handler.go
│   │   └── indexing_handler.go
│   └── model/
│       ├── search.go
│       ├── document.go
│       └── hashtag.go
├── go.mod
├── Dockerfile
└── README.md
```

### Code Statistics
```
Total Lines: 4,750+
Files: 20
Handlers: 30+ endpoints
Services: 4
Models: 25+
```

---

## 🎯 Key Features Explained

### Multi-Match Query
```go
fields: [
  "username^3",        // Highest priority
  "display_name^2",    // Medium priority
  "caption^2",
  "bio"                // Standard priority
]
fuzziness: "AUTO"      // Typo tolerance
minimum_should_match: "75%"
```

### Trending Algorithm
```
Growth Rate = (current_usage - previous_usage) / previous_usage * 100
Rank by: growth_rate DESC, usage_count DESC
```

### Related Hashtags
Uses Elasticsearch "More Like This" query:
- Finds similar hashtags based on co-occurrence
- Filters by usage count
- Returns top N results

---

## 🚀 Deployment

### Environment Variables
```env
PORT=8087
ELASTICSEARCH_URL=http://elasticsearch:9200
ELASTICSEARCH_USERNAME=elastic
ELASTICSEARCH_PASSWORD=changeme
REDIS_URL=redis:6379
REDIS_PASSWORD=
GIN_MODE=release
```

### Docker Compose
```yaml
version: '3.8'
services:
  search-service:
    build: .
    ports:
      - "8087:8087"
    environment:
      - ELASTICSEARCH_URL=http://elasticsearch:9200
      - REDIS_URL=redis:6379
    depends_on:
      - elasticsearch
      - redis
  
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.11.0
    environment:
      - discovery.type=single-node
    volumes:
      - es-data:/usr/share/elasticsearch/data
  
  redis:
    image: redis:7-alpine

volumes:
  es-data:
```

---

## 📊 Monitoring

### Health Check
```bash
curl http://localhost:8087/health
```

### Index Statistics
```bash
curl http://localhost:8087/api/v1/index/stats
```

---

**Socialink Search Service** - Powered by Elasticsearch 🔍⚡
