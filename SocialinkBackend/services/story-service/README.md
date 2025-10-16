# Socialink Story Service 🎬

**Production-grade ephemeral story service** with interactive stickers, highlights, and comprehensive analytics for the Socialink platform (Facebook-like).

---

## 🔥 Features

### Core Stories
- ⏱️ **24-hour ephemeral stories** (Facebook Stories style)
- 📸 **Multiple formats**: Images, videos, text
- 🎨 **Custom backgrounds**: Colors, gradients
- 🔒 **Privacy controls**: Public, friends, close friends, custom
- 👀 **View tracking**: "Seen by" lists with timestamps
- 💬 **Story replies**: Private messages to story owner
- 🎵 **Music integration**: Add soundtracks

### Interactive Stickers 🎯
- **Polls** - Multiple choice voting with live percentages
- **Quizzes** - Test friends (correct/wrong tracking)
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

## 🚀 Quick Start

### Prerequisites
- Python 3.11+
- MongoDB 5.0+
- Redis 6.0+

### Installation

```bash
# Navigate to service
cd SocialinkBackend/services/story-service

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
MONGODB_DB_NAME=socialink_stories

# Redis
REDIS_URL=redis://localhost:6379

# Story Settings
STORY_EXPIRY_HOURS=24
MAX_STORIES_PER_USER=100
```

### Run

```bash
# Development
python -m uvicorn app.main:app --reload --port 8086

# Production
uvicorn app.main:app --host 0.0.0.0 --port 8086 --workers 4
```

### Docker

```bash
# Build
docker build -t socialink-story-service .

# Run
docker run -d \
  -p 8086:8086 \
  -e MONGODB_URL=mongodb://mongo:27017 \
  -e REDIS_URL=redis://redis:6379 \
  socialink-story-service
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
curl -X POST http://localhost:8086/api/v1/stories \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user123" \
  -d '{
    "story_type": "image",
    "media_id": "media-uuid",
    "duration": 10,
    "privacy": "friends",
    "stickers": [
      {
        "sticker_type": "poll",
        "x": 0.5,
        "y": 0.7,
        "width": 0.8,
        "height": 0.3,
        "poll": {
          "question": "Coming to the party tonight?",
          "options": [
            {"text": "Yes!", "votes": 0},
            {"text": "Maybe", "votes": 0},
            {"text": "Can'\''t make it", "votes": 0}
          ]
        }
      }
    ]
  }'
```

### Vote on Poll

```bash
curl -X POST http://localhost:8086/api/v1/stories/story-id/interact \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user456" \
  -d '{
    "story_id": "story-id",
    "sticker_index": 0,
    "option_index": 0
  }'
```

---

## 🎯 Facebook-Style Features

### Story Buckets
- Stories grouped by user
- Close friends appear first
- Unseen stories highlighted
- Gradient rings for new content

### Privacy Levels
- **Public** - Anyone can see
- **Friends** - Only friends
- **Close Friends** - Inner circle only
- **Custom** - Specific people

### Engagement
- Story replies (private DMs)
- Interactive stickers
- View tracking
- Analytics for creators

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

## 🚀 Deployment

### Environment Variables
```env
# Production settings
DEBUG=False
MONGODB_URL=mongodb://prod-mongo:27017
REDIS_URL=redis://prod-redis:6379
JWT_SECRET=strong-secret-key-here
PORT=8086
```

---

## 📊 Monitoring

### Health Check
```bash
curl http://localhost:8086/health
```

Response:
```json
{
  "status": "healthy",
  "service": "Socialink Story Service",
  "version": "1.0.0"
}
```

---

**Socialink Story Service** - Built with ❤️ by the Entativa team
