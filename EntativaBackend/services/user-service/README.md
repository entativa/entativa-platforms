# Socialink User Authentication Service

Meta-level authentication service for Socialink (Facebook-like social platform). Provides instant access authentication without email/phone verification.

## 🚀 Features

- **Meta-Level Authentication**: Instant account creation and access
- **No Verification Required**: Users can start using the platform immediately after signup
- **JWT-based Authentication**: Secure token-based authentication
- **Session Management**: Track user sessions across devices
- **Password Security**: Bcrypt password hashing
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
http://localhost:8001/api/v1
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
  "service": "socialink-user-service",
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
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "password": "SecurePassword123",
  "birthday": "1995-05-15",
  "gender": "male"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Account created successfully! Welcome to Socialink!",
  "data": {
    "user": {
      "id": "uuid",
      "first_name": "John",
      "last_name": "Doe",
      "email": "john.doe@example.com",
      "username": "john.doe1234",
      "birthday": "1995-05-15T00:00:00Z",
      "gender": "male",
      "is_active": true,
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
  "email_or_username": "john.doe@example.com",
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
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "username": "john.doe1234",
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
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(30) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    birthday DATE NOT NULL,
    gender VARCHAR(20) NOT NULL,
    phone_number VARCHAR(20),
    bio TEXT,
    profile_picture_url TEXT,
    cover_photo_url TEXT,
    is_active BOOLEAN DEFAULT true,
    is_deleted BOOLEAN DEFAULT false,
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
│       └── validation.go        # Input validation
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
docker build -t socialink-user-service .
docker run -p 8001:8001 --env-file .env socialink-user-service
```

## 🔒 Security Features

- **Password Hashing**: Bcrypt with default cost factor
- **JWT Tokens**: HS256 algorithm with configurable expiration
- **HTTPS Ready**: Support for TLS/SSL in production
- **CORS Protection**: Configurable CORS policies
- **SQL Injection Prevention**: Parameterized queries
- **Rate Limiting Ready**: Middleware structure supports rate limiting

## 🎯 Meta-Level Authentication Philosophy

This service implements Meta's approach to user onboarding:
- **Frictionless Signup**: No email verification barriers
- **Instant Access**: Users can immediately start using the platform
- **Trust First**: Build trust through user experience, not verification gates
- **Progressive Verification**: Future enhancements can add optional verification for features

## 📝 License

Proprietary - Socialink Platform

## 👥 Contributing

This is a proprietary service for the Socialink platform.
