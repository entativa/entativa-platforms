# Vignette Story Service 🎬

**Production-grade ephemeral story service** with interactive stickers, highlights, and comprehensive analytics.

---

## 🔥 Features

### Core Stories
- ⏱️ **24-hour ephemeral stories** (Instagram/Snapchat style)
- 📸 **Multiple formats**: Images, videos, text
- 🎨 **Custom backgrounds**: Colors, gradients
- 🔒 **Privacy controls**: Public, followers, close friends, custom
- 👀 **View tracking**: "Seen by" lists with timestamps
- 💬 **Story replies**: Private messages to story owner
- 🎵 **Music integration**: Add soundtracks

### Interactive Stickers 🎯
- **Polls** - Multiple choice voting with live percentages
- **Quizzes** - Test followers (correct/wrong tracking)
- **Questions** - Open-ended responses
- **Countdowns** - Event timers with follower notifications
- **Sliders** - Emoji sliders for ratings
- **Mentions** - Tag other users
- **Location** - Share where you are
- **Music** - Soundtrack stickers

### Story Highlights 📌
- Permanent story collections
- Custom covers (images, emojis, colors)
- Organize and reorder
- Pin important highlights
- Archive management

### Analytics & Insights 📊
- View counts and unique viewers
- Reach percentage
- Interaction rates
- Peak viewing times
- Sticker performance metrics
- Top viewers list
- Engagement scoring

### Close Friends 💚
- Private close friends lists
- Stories visible only to close friends
- Fast Redis-backed lookups

---

## 🏗️ Architecture

```
Story Service (Python/FastAPI)
├── MongoDB (Story data)
├── Redis (Caching + Real-time)
├── Background Tasks (Expiration)
└── gRPC Client (Media service)
```

### Tech Stack
- **FastAPI** - High-performance async API
- **Motor** - Async MongoDB driver
- **Redis** - Caching and real-time features
- **MongoDB** - Document storage (perfect for stories)
- **APScheduler** - Background job scheduling
- **Pydantic** - Data validation

---

## 📊 Database Schema

### Stories Collection
```json
{
  "_id": "ObjectId",
  "user_id": "string",
  "story_type": "image|video|text",
  "media_id": "string",
  "media_url": "string",
  "thumbnail_url": "string",
  "text_content": "string",
  "background_color": "string",
  "background_gradient": ["color1", "color2"],
  "duration": 5,
  "stickers": [
    {
      "sticker_type": "poll|quiz|question|countdown|slider",
      "x": 0.5, "y": 0.5,
      "width": 0.5, "height": 0.3,
      "rotation": 0,
      "poll": {
        "question": "string",
        "options": [
          {"text": "string", "votes": 0, "percentage": 0}
        ],
        "total_votes": 0
      }
    }
  ],
  "privacy": "public|followers|close_friends|custom",
  "close_friends_only": false,
  "viewers": [
    {
      "user_id": "string",
      "username": "string",
      "viewed_at": "datetime"
    }
  ],
  "view_count": 0,
  "replies": [],
  "reply_count": 0,
  "is_expired": false,
  "is_deleted": false,
  "created_at": "datetime",
  "expires_at": "datetime",
  "updated_at": "datetime"
}
```

### Highlights Collection
```json
{
  "_id": "ObjectId",
  "user_id": "string",
  "title": "string",
  "description": "string",
  "cover": {
    "story_id": "string",
    "media_url": "string",
    "emoji": "string",
    "color": "#FF6B6B"
  },
  "story_ids": ["story_id1", "story_id2"],
  "story_count": 0,
  "is_archived": false,
  "is_pinned": false,
  "order_index": 0,
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### Close Friends Collection
```json
{
  "user_id": "string",
  "friend_user_id": "string",
  "added_at": "datetime",
  "is_mutual": false
}
```

---

## 🚀 Quick Start

### Prerequisites
- Python 3.11+
- MongoDB 5.0+
- Redis 6.0+

### Installation

```bash
# Navigate to service
cd VignetteBackend/services/story-service

# Install dependencies
pip install -r requirements.txt

# Copy environment file
cp .env.example .env

# Edit .env with your settings
nano .env
```

### Configuration

```env
# MongoDB
MONGODB_URL=mongodb://localhost:27017
MONGODB_DB_NAME=vignette_stories

# Redis
REDIS_URL=redis://localhost:6379

# Story Settings
STORY_EXPIRY_HOURS=24
MAX_STORIES_PER_USER=100
```

### Run

```bash
# Development
python -m uvicorn app.main:app --reload --port 8085

# Production
uvicorn app.main:app --host 0.0.0.0 --port 8085 --workers 4
```

### Docker

```bash
# Build
docker build -t vignette-story-service .

# Run
docker run -d \
  -p 8085:8085 \
  -e MONGODB_URL=mongodb://mongo:27017 \
  -e REDIS_URL=redis://redis:6379 \
  vignette-story-service
```

---

## 📡 API Endpoints

### Stories

```http
POST   /api/v1/stories              - Create story
GET    /api/v1/stories/feed         - Get story feed
GET    /api/v1/stories/{story_id}   - Get story
GET    /api/v1/stories/user/{user_id} - Get user stories
POST   /api/v1/stories/{story_id}/view - Mark as viewed
POST   /api/v1/stories/{story_id}/reply - Reply to story
POST   /api/v1/stories/{story_id}/interact - Interact with sticker
DELETE /api/v1/stories/{story_id}   - Delete story
```

### Highlights

```http
POST   /api/v1/highlights              - Create highlight
GET    /api/v1/highlights/user/{user_id} - Get user highlights
GET    /api/v1/highlights/{highlight_id} - Get highlight
PATCH  /api/v1/highlights/{highlight_id} - Update highlight
POST   /api/v1/highlights/{highlight_id}/stories - Add stories
DELETE /api/v1/highlights/{highlight_id}/stories - Remove stories
DELETE /api/v1/highlights/{highlight_id} - Delete highlight
POST   /api/v1/highlights/reorder      - Reorder highlights
```

### Analytics

```http
GET /api/v1/viewers/story/{story_id} - Get story viewers
GET /api/v1/viewers/story/{story_id}/insights - Get story insights
GET /api/v1/viewers/user/insights - Get user insights
```

### Close Friends

```http
POST   /api/v1/close-friends/{friend_user_id} - Add close friend
DELETE /api/v1/close-friends/{friend_user_id} - Remove close friend
GET    /api/v1/close-friends/                 - Get close friends list
```

---

## 💡 Usage Examples

### Create a Story with Poll

```bash
curl -X POST http://localhost:8085/api/v1/stories \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user123" \
  -d '{
    "story_type": "image",
    "media_id": "media-uuid",
    "duration": 10,
    "privacy": "followers",
    "stickers": [
      {
        "sticker_type": "poll",
        "x": 0.5,
        "y": 0.7,
        "width": 0.8,
        "height": 0.3,
        "poll": {
          "question": "Best programming language?",
          "options": [
            {"text": "Python", "votes": 0},
            {"text": "Go", "votes": 0},
            {"text": "Rust", "votes": 0}
          ]
        }
      }
    ]
  }'
```

### Vote on Poll

```bash
curl -X POST http://localhost:8085/api/v1/stories/story-id/interact \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user456" \
  -d '{
    "story_id": "story-id",
    "sticker_index": 0,
    "option_index": 0
  }'
```

### Create Highlight

```bash
curl -X POST http://localhost:8085/api/v1/highlights \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user123" \
  -d '{
    "title": "Summer 2025",
    "description": "Best vacation ever!",
    "story_ids": ["story1", "story2", "story3"],
    "cover": {
      "emoji": "☀️",
      "color": "#FFD700"
    },
    "is_pinned": true
  }'
```

### Get Story Analytics

```bash
curl -X GET http://localhost:8085/api/v1/viewers/story/story-id/insights \
  -H "X-User-ID: user123"
```

Response:
```json
{
  "story_id": "story-id",
  "total_views": 1523,
  "unique_viewers": 847,
  "reach_percentage": 42.3,
  "total_interactions": 234,
  "interaction_rate": 27.6,
  "replies": 12,
  "views_by_hour": {
    "14": 234,
    "15": 456,
    "16": 389
  },
  "peak_viewing_hour": 15,
  "sticker_interactions": {
    "sticker_0_poll": 156,
    "sticker_1_quiz": 78
  }
}
```

---

## 🎯 Interactive Stickers Guide

### Poll Sticker
```json
{
  "sticker_type": "poll",
  "poll": {
    "question": "Which one?",
    "options": [
      {"text": "Option A", "votes": 0},
      {"text": "Option B", "votes": 0}
    ]
  }
}
```
**Features**: Live vote percentages, one vote per user

### Quiz Sticker
```json
{
  "sticker_type": "quiz",
  "quiz": {
    "question": "Capital of France?",
    "options": [
      {"text": "London", "is_correct": false},
      {"text": "Paris", "is_correct": true},
      {"text": "Berlin", "is_correct": false}
    ],
    "correct_answer_index": 1
  }
}
```
**Features**: Correct/wrong feedback, attempt tracking

### Question Sticker
```json
{
  "sticker_type": "question",
  "question": {
    "text": "Ask me anything!"
  }
}
```
**Features**: Open-ended responses, view all answers

### Slider Sticker
```json
{
  "sticker_type": "slider",
  "slider": {
    "question": "How much do you like this?",
    "emoji": "😍",
    "min_value": 0,
    "max_value": 100
  }
}
```
**Features**: Emoji slider (0-100), average calculation

### Countdown Sticker
```json
{
  "sticker_type": "countdown",
  "countdown": {
    "title": "Product Launch",
    "end_time": "2025-10-20T00:00:00Z"
  }
}
```
**Features**: Live countdown, follower notifications

---

## 🔧 Background Services

### Story Expiration
- Runs every 5 minutes
- Marks stories older than 24 hours as expired
- Clears related caches
- Optional: Hard delete after 30 days

### Implementation
```python
# In app/services/expiration_service.py
async def expire_old_stories():
    """Mark expired stories"""
    now = datetime.utcnow()
    
    await db.stories.update_many(
        {
            "is_expired": False,
            "expires_at": {"$lte": now}
        },
        {"$set": {"is_expired": True}}
    )
```

---

## 📈 Performance

### Caching Strategy
- **Story Cache**: 1 hour TTL
- **User Stories**: Invalidated on new story
- **Highlights**: 1 hour TTL
- **Close Friends**: 1 hour TTL
- **View Status**: 24 hour TTL

### Database Indexes
```javascript
// Stories
db.stories.createIndex({"user_id": 1})
db.stories.createIndex({"created_at": -1})
db.stories.createIndex({"expires_at": 1})
db.stories.createIndex({"user_id": 1, "created_at": -1})
db.stories.createIndex({"is_expired": 1, "expires_at": 1})

// Highlights
db.highlights.createIndex({"user_id": 1})
db.highlights.createIndex({"user_id": 1, "order_index": 1})

// Close Friends
db.close_friends.createIndex({"user_id": 1, "friend_user_id": 1}, {unique: true})
```

---

## 🧪 Testing

```bash
# Run tests
pytest

# With coverage
pytest --cov=app --cov-report=html

# Specific test
pytest tests/test_stories.py
```

---

## 🚀 Deployment

### Environment Variables
```env
# Production settings
DEBUG=False
MONGODB_URL=mongodb://prod-mongo:27017
REDIS_URL=redis://prod-redis:6379
JWT_SECRET=strong-secret-key-here
```

### Docker Compose
```yaml
version: '3.8'
services:
  story-service:
    build: .
    ports:
      - "8085:8085"
    environment:
      - MONGODB_URL=mongodb://mongo:27017
      - REDIS_URL=redis://redis:6379
    depends_on:
      - mongo
      - redis
  
  mongo:
    image: mongo:7
    volumes:
      - mongo-data:/data/db
  
  redis:
    image: redis:7-alpine
    
volumes:
  mongo-data:
```

---

## 📊 Monitoring

### Health Check
```bash
curl http://localhost:8085/health
```

Response:
```json
{
  "status": "healthy",
  "service": "Vignette Story Service",
  "version": "1.0.0"
}
```

### Metrics (TODO)
- Story creation rate
- View rate
- Interaction rate
- Cache hit rate
- API response times

---

## 🎯 Roadmap

- [ ] Story mentions with notifications
- [ ] Story music integration
- [ ] Link stickers
- [ ] AR filters (via media service)
- [ ] Story ads
- [ ] Story gifting
- [ ] Cross-posting to feed
- [ ] Story templates
- [ ] Collaborative stories

---

## 📝 License

Proprietary - Entativa

---

## 🤝 Contributing

Internal project - contact team lead for contribution guidelines.

---

**Vignette Story Service** - Built with ❤️ by the Entativa team
