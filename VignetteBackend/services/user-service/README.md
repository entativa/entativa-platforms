# Vignette User Authentication Service

Meta-level authentication service for Vignette (Instagram-like social platform). Provides instant access authentication without email/phone verification.

## 🚀 Features

- **Meta-Level Authentication**: Instant account creation and access
- **No Verification Required**: Users can start using the platform immediately after signup
- **JWT-based Authentication**: Secure token-based authentication
- **Session Management**: Track user sessions across devices
- **Password Security**: Bcrypt password hashing
- **Instagram-Style Usernames**: Validate usernames with Instagram-like rules
- **RESTful API**: Clean and intuitive API design

## 📋 Prerequisites

- Go 1.21 or higher
- PostgreSQL 12+
- Redis (optional, for future caching)

## 🛠️ Installation

1. Clone the repository
2. Copy `.env.example` to `.env` and configure your environment variables
3. Install dependencies:
```bash
go mod download
```

4. Run the service:
```bash
go run cmd/api/main.go
```

## 🔧 Environment Variables

See `.env.example` for all available configuration options.

## 📚 API Documentation

### Base URL
```
http://localhost:8002/api/v1
```

### Endpoints

#### 1. Health Check
```http
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "service": "vignette-user-service",
  "version": "1.0.0"
}
```

#### 2. Sign Up
```http
POST /auth/signup
```

**Request Body:**
```json
{
  "username": "johndoe",
  "email": "john.doe@example.com",
  "full_name": "John Doe",
  "password": "SecurePassword123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Account created successfully! Welcome to Vignette!",
  "data": {
    "user": {
      "id": "uuid",
      "username": "johndoe",
      "email": "john.doe@example.com",
      "full_name": "John Doe",
      "is_private": false,
      "is_verified": false,
      "is_active": true,
      "followers_count": 0,
      "following_count": 0,
      "posts_count": 0,
      "created_at": "2025-10-15T..."
    },
    "access_token": "eyJhbGci...",
    "token_type": "Bearer",
    "expires_in": 86400
  }
}
```

#### 3. Login
```http
POST /auth/login
```

**Request Body:**
```json
{
  "username_or_email": "johndoe",
  "password": "SecurePassword123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful! Welcome back!",
  "data": {
    "user": { ... },
    "access_token": "eyJhbGci...",
    "token_type": "Bearer",
    "expires_in": 86400
  }
}
```

#### 4. Get Current User
```http
GET /auth/me
Authorization: Bearer <access_token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "username": "johndoe",
    "email": "john.doe@example.com",
    "full_name": "John Doe",
    ...
  }
}
```

#### 5. Logout
```http
POST /auth/logout
Authorization: Bearer <access_token>
```

**Response:**
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

## 🗄️ Database Schema

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    username VARCHAR(30) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    full_name VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20),
    bio VARCHAR(150),
    website VARCHAR(255),
    profile_picture_url TEXT,
    is_private BOOLEAN DEFAULT false,
    is_verified BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    is_deleted BOOLEAN DEFAULT false,
    followers_count INTEGER DEFAULT 0,
    following_count INTEGER DEFAULT 0,
    posts_count INTEGER DEFAULT 0,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### Sessions Table
```sql
CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    device_info VARCHAR(255),
    ip_address VARCHAR(45),
    user_agent TEXT,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_active_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

## 🏗️ Architecture

```
user-service/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── handler/
│   │   └── auth_handler.go      # HTTP handlers
│   ├── middleware/
│   │   ├── auth_middleware.go   # JWT authentication middleware
│   │   ├── cors_middleware.go   # CORS middleware
│   │   └── logger_middleware.go # Logging middleware
│   ├── model/
│   │   ├── user.go              # User models
│   │   └── session.go           # Session models
│   ├── repository/
│   │   ├── user_repository.go   # User data access
│   │   └── session_repository.go # Session data access
│   ├── service/
│   │   └── auth_service.go      # Business logic
│   └── util/
│       ├── jwt.go               # JWT utilities
│       ├── password.go          # Password hashing
│       └── validation.go        # Input validation (Instagram-style)
├── pkg/
│   └── database/
│       ├── postgres.go          # Database connection
│       └── migrations.go        # Database migrations
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

## 🐳 Docker

Build and run with Docker:

```bash
docker build -t vignette-user-service .
docker run -p 8002:8002 --env-file .env vignette-user-service
```

## 🔒 Security Features

- **Password Hashing**: Bcrypt with default cost factor
- **JWT Tokens**: HS256 algorithm with configurable expiration
- **HTTPS Ready**: Support for TLS/SSL in production
- **CORS Protection**: Configurable CORS policies
- **SQL Injection Prevention**: Parameterized queries
- **Rate Limiting Ready**: Middleware structure supports rate limiting

## 📝 Username Rules (Instagram-Style)

- **Length**: 3-30 characters
- **Allowed characters**: Letters, numbers, periods (.), underscores (_)
- **Restrictions**:
  - Cannot start or end with a period
  - Cannot have consecutive periods
  - Must be unique across the platform

## 🎯 Meta-Level Authentication Philosophy

This service implements Meta's approach to user onboarding:
- **Frictionless Signup**: No email verification barriers
- **Instant Access**: Users can immediately start using the platform
- **Trust First**: Build trust through user experience, not verification gates
- **Progressive Verification**: Future enhancements can add optional verification for features

## 📝 License

Proprietary - Vignette Platform

## 👥 Contributing

This is a proprietary service for the Vignette platform.
