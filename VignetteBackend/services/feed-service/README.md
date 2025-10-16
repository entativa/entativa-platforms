# Vignette Feed Service 🔥

**World-class feed algorithms competing with TikTok's For You!**

---

## 🎯 Overview

The Vignette Feed Service provides **3 sophisticated feed algorithms** for maximum user engagement and discovery:

1. **Home Feed** - TikTok-style For You feed (viral discovery)
2. **Circle Feed** - Friends + nearby content (social graph)
3. **Surprise & Delight** - Balanced exploration (60% interests, 30% explore, 10% surprise)

---

## 🚀 Features

### Feed Algorithms

#### 1. Home Feed (For You)
**The TikTok Competitor!**

Optimized for:
- 🔥 Viral discovery (35% weight)
- 🎯 Personalization (30% weight)
- ⭐ Quality content (15% weight)
- 🆕 Freshness (15% weight)
- 🎨 Diversity (5% weight)

**Use case**: Discovering trending content, new creators, viral videos

#### 2. Circle Feed
**The social graph prioritizer!**

Optimized for:
- 👥 Social connections (50% weight)
- 📍 Location proximity (25% weight)
- 💬 Engagement (15% weight)
- 🆕 Recency (10% weight)

**Use case**: Staying connected with friends and nearby community

#### 3. Surprise & Delight Feed
**The balanced explorer!**

Mix:
- 60% Known interests (familiar content)
- 30% Exploration (adjacent topics)
- 10% Surprises (random viral)

**Use case**: Curated discovery, less overwhelming than Home feed

### Core Features
- ✅ Real-time signal tracking (views, likes, shares, skips)
- ✅ User profiling & interest extraction
- ✅ Collaborative filtering
- ✅ Location-based recommendations
- ✅ Smart diversity constraints
- ✅ Redis caching for performance
- ✅ Feed metrics & analytics

---

## 📡 API Endpoints

### Get Feed
```
GET /api/v1/feed/{feed_type}?limit=20&offset=0
```

**Feed types**:
- `home` - For You feed
- `circle` - Friends + nearby
- `surprise_delight` - Balanced exploration

**Query params**:
- `limit` (1-100, default: 20)
- `offset` (default: 0)
- `seen_ids` (comma-separated)
- `refresh` (boolean)
- `latitude` (for Circle feed)
- `longitude` (for Circle feed)

**Response**:
```json
{
  "feed_type": "home",
  "items": [
    {
      "content_id": "...",
      "content_type": "take",
      "user_id": "...",
      "username": "...",
      "caption": "...",
      "media_urls": [...],
      "likes_count": 1234,
      "score": 0.87,
      "source": "viral",
      "rank": 0,
      ...
    }
  ],
  "next_offset": 20,
  "has_more": true
}
```

### Track Signal
```
POST /api/v1/feed/signal
```

**Body**:
```json
{
  "user_id": "...",
  "content_id": "...",
  "content_type": "take",
  "signal_type": "like",
  "time_spent_seconds": 15.5,
  "completion_rate": 0.95
}
```

**Signal types**:
- `view` - Content viewed
- `like` - Content liked
- `comment` - User commented
- `share` - Content shared
- `save` - Content saved
- `skip` - Content skipped
- `hide` - Content hidden
- `report` - Content reported

### Get Feed Metrics
```
GET /api/v1/feed/metrics/{feed_type}?hours=24
```

**Response**:
```json
{
  "feed_type": "home",
  "user_id": "...",
  "items_shown": 150,
  "items_engaged": 45,
  "avg_time_spent_seconds": 18.5,
  "skip_rate": 0.12,
  "engagement_rate": 0.30
}
```

---

## 🏗️ Architecture

```
Feed Service
├── API Layer (FastAPI)
├── Feed Service (Orchestration)
├── ML Layer
│   └── Recommendation Engine (3 algorithms)
├── Services
│   ├── Ranking Service (Candidate retrieval)
│   └── Personalization Service (User profiling)
├── Storage
│   ├── MongoDB (Content data)
│   ├── Redis (Caching, signals)
│   └── Elasticsearch (Content search - optional)
└── Integration
    ├── Post Service (Content data)
    ├── User Service (Social graph)
    └── Story Service (Story content)
```

---

## 🧠 How It Works

### 1. User Profile Building
- Extracts interests from engagement history
- Identifies favorite creators, topics, hashtags
- Builds social graph (friends, following, followers)
- Tracks engagement patterns & preferences

### 2. Candidate Retrieval
- **Home**: Viral content + trending topics + similar content
- **Circle**: Friends' content + following + nearby
- **Surprise**: Mix of known interests + exploration + viral

### 3. Ranking
- Scores each candidate using feed-specific algorithm
- Applies diversity constraints (prevent creator fatigue)
- Interleaves categories for natural flow

### 4. Learning
- Tracks all user signals (views, likes, skips)
- Updates user profile in real-time
- Adapts future recommendations

---

## ⚙️ Configuration

### Environment Variables
```env
# Service
PORT=8085
DEBUG=false

# Databases
MONGODB_URL=mongodb://localhost:27017
MONGODB_DB=vignette_feed
POSTGRES_URL=postgresql://...
REDIS_URL=redis://localhost:6379
ELASTICSEARCH_URL=http://localhost:9200

# Feed Settings
DEFAULT_FEED_SIZE=20
MAX_FEED_SIZE=100
CANDIDATE_POOL_SIZE=500

# Cache TTL
USER_PROFILE_TTL=3600
VIRAL_CONTENT_TTL=300
```

### Ranking Weights
Adjust in `app/config.py`:

```python
# Home Feed
HOME_VIRAL_WEIGHT = 0.35
HOME_INTEREST_WEIGHT = 0.30
HOME_QUALITY_WEIGHT = 0.15
HOME_RECENCY_WEIGHT = 0.15
HOME_DIVERSITY_WEIGHT = 0.05

# Circle Feed
CIRCLE_SOCIAL_WEIGHT = 0.50
CIRCLE_PROXIMITY_WEIGHT = 0.25
CIRCLE_ENGAGEMENT_WEIGHT = 0.15
CIRCLE_RECENCY_WEIGHT = 0.10

# Surprise & Delight
SURPRISE_KNOWN_RATIO = 0.60
SURPRISE_EXPLORE_RATIO = 0.30
SURPRISE_RANDOM_RATIO = 0.10
```

---

## 🚀 Quick Start

### Installation
```bash
cd VignetteBackend/services/feed-service

# Install dependencies
pip install -r requirements.txt
```

### Run
```bash
# Development
python -m app.main

# Production
uvicorn app.main:app --host 0.0.0.0 --port 8085 --workers 4
```

### Docker
```bash
docker build -t vignette-feed-service .
docker run -p 8085:8085 vignette-feed-service
```

---

## 📊 Performance

### Targets
- **Feed generation**: <200ms
- **Candidate retrieval**: <100ms
- **Ranking**: <50ms
- **Signal tracking**: <10ms

### Optimization
- Redis caching (user profiles, viral content)
- Batch candidate retrieval
- Async/await everywhere
- Connection pooling
- Query optimization

---

## 🎯 Feed Selection Guide

**Use Home Feed when**:
- User wants to discover new content
- Exploring trending topics
- Looking for viral entertainment

**Use Circle Feed when**:
- User wants to see friends' updates
- Checking local/nearby content
- Staying connected with close circle

**Use Surprise & Delight when**:
- User wants curated exploration
- Less overwhelming than Home feed
- Balanced mix of familiar + new

---

## 🔥 Why This Beats TikTok

| Feature | Vignette | TikTok |
|---------|----------|--------|
| Feed Algorithms | **3** | 1 (For You) |
| Social Priority | ✅ Circle | ❌ |
| Balanced Exploration | ✅ S&D | ❌ |
| User Control | ✅ Toggle | ❌ |
| Location-based | ✅ | Limited |
| Signal Tracking | **8 types** | Basic |

---

## 📈 Metrics & Monitoring

Track these KPIs:
- **Engagement rate** (likes + comments + shares / views)
- **Skip rate** (how often users skip content)
- **Time spent** (average per item)
- **Feed diversity** (unique creators, topics)
- **Completion rate** (for Takes)

---

## 🎊 Summary

**Vignette Feed Service** provides **3 world-class feed algorithms** that give users:
- 🔥 **Discovery** (Home feed - TikTok competitor)
- 👥 **Connection** (Circle feed - social graph)
- 🎨 **Exploration** (Surprise & Delight - balanced)

**Tech**: Python + FastAPI + MongoDB + Redis + Elasticsearch  
**Performance**: Sub-200ms feed generation  
**Status**: Production-ready  

**LET'S GOOOOO! 🚀🔥**
