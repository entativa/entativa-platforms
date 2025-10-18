# 🚀 Advanced Features Guide - Meta-Level Authentication Services

## 📋 Table of Contents

1. [Overview](#overview)
2. [Two-Factor Authentication (2FA)](#two-factor-authentication-2fa)
3. [Password Reset](#password-reset)
4. [Account Recovery](#account-recovery)
5. [Profile Picture Upload](#profile-picture-upload)
6. [Rate Limiting](#rate-limiting)
7. [Redis Caching](#redis-caching)
8. [gRPC Communication](#grpc-communication)
9. [Kafka Event Streaming](#kafka-event-streaming)
10. [Configuration](#configuration)

---

## Overview

Both Entativa and Vignette services now include enterprise-grade advanced features while maintaining the Meta-level instant-access authentication philosophy.

### New Features Summary

✅ **Two-Factor Authentication** - TOTP-based 2FA with backup codes  
✅ **Password Reset** - Secure token-based password reset flows  
✅ **Account Recovery** - Multiple recovery methods  
✅ **Media Upload** - S3/MinIO integration for profile pictures  
✅ **Rate Limiting** - Redis-powered intelligent rate limiting  
✅ **Caching** - Redis caching for improved performance  
✅ **gRPC** - Microservice communication protocol  
✅ **Kafka** - Event-driven architecture support  

---

## Two-Factor Authentication (2FA)

### Overview

TOTP-based Two-Factor Authentication using industry-standard algorithms. Compatible with Google Authenticator, Authy, and other authenticator apps.

### Features

- **TOTP Generation**: Time-based One-Time Passwords
- **QR Code**: Easy setup via QR code scanning
- **Backup Codes**: 10 backup codes for recovery
- **Optional**: Users can enable/disable 2FA anytime
- **Security**: Codes expire after 30 seconds

### API Endpoints

#### 1. Setup 2FA
```http
POST /api/v1/auth/2fa/setup
Authorization: Bearer <access_token>
```

**Response:**
```json
{
  "success": true,
  "message": "2FA setup initiated. Scan the QR code with your authenticator app",
  "data": {
    "secret": "JBSWY3DPEHPK3PXP",
    "qr_code_url": "otpauth://totp/Entativa:user@example.com?secret=...",
    "backup_codes": [
      "ABCD1234EFGH5678",
      "WXYZ9876QRST5432",
      ...
    ]
  }
}
```

#### 2. Enable 2FA
```http
POST /api/v1/auth/2fa/enable
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "code": "123456"
}
```

#### 3. Verify 2FA Code
```http
POST /api/v1/auth/2fa/verify
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "code": "123456"
}
```

#### 4. Disable 2FA
```http
POST /api/v1/auth/2fa/disable
Authorization: Bearer <access_token>
```

### Usage Example

```bash
# 1. Setup 2FA
curl -X POST http://localhost:8001/api/v1/auth/2fa/setup \
  -H "Authorization: Bearer YOUR_TOKEN" | jq '.'

# 2. Scan QR code with authenticator app

# 3. Enable 2FA with code from app
curl -X POST http://localhost:8001/api/v1/auth/2fa/enable \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"code": "123456"}'

# 4. Verify 2FA (during login)
curl -X POST http://localhost:8001/api/v1/auth/2fa/verify \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"code": "789012"}'
```

### Database Schema

```sql
CREATE TABLE two_factor_auth (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE REFERENCES users(id),
    secret VARCHAR(255) NOT NULL,
    is_enabled BOOLEAN DEFAULT false,
    backup_codes JSONB,
    created_at TIMESTAMP NOT NULL,
    enabled_at TIMESTAMP,
    last_used_at TIMESTAMP
);
```

---

## Password Reset

### Overview

Secure password reset flow using time-limited tokens. Tokens expire after 1 hour.

### Features

- **Secure Tokens**: Cryptographically secure random tokens
- **Time-Limited**: Tokens expire after 1 hour
- **Email Notification**: Via Kafka events (optional)
- **One-Time Use**: Tokens can only be used once
- **Privacy**: Doesn't reveal if email exists

### API Endpoints

#### 1. Request Password Reset
```http
POST /api/v1/auth/password-reset/request
Content-Type: application/json

{
  "email": "user@example.com"
}
```

**Response:**
```json
{
  "success": true,
  "message": "If an account exists with this email, you will receive password reset instructions"
}
```

#### 2. Reset Password
```http
POST /api/v1/auth/password-reset/confirm
Content-Type: application/json

{
  "token": "secure-reset-token-from-email",
  "new_password": "NewSecurePassword123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Password reset successfully. You can now log in with your new password"
}
```

### Usage Example

```bash
# 1. Request password reset
curl -X POST http://localhost:8001/api/v1/auth/password-reset/request \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com"}'

# 2. Check email for reset token (or Kafka event)

# 3. Reset password with token
curl -X POST http://localhost:8001/api/v1/auth/password-reset/confirm \
  -H "Content-Type: application/json" \
  -d '{
    "token": "abc123def456...",
    "new_password": "NewPassword123"
  }'
```

### Database Schema

```sql
CREATE TABLE password_reset_tokens (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    token VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL
);
```

---

## Account Recovery

### Overview

Multiple recovery methods including email, phone, and security questions (extensible).

### Features

- **Email Recovery**: Primary recovery method
- **Phone Recovery**: SMS-based recovery (future)
- **Security Questions**: Additional recovery option (future)
- **Extensible**: Easy to add new recovery methods

### Database Schema

```sql
CREATE TABLE account_recovery_methods (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    type VARCHAR(50) NOT NULL, -- email, phone, security_question
    value TEXT NOT NULL,
    is_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP NOT NULL
);
```

---

## Profile Picture Upload

### Overview

S3/MinIO integration for media uploads with automatic optimization and CDN support.

### Features

- **S3 Compatible**: Works with AWS S3 or MinIO
- **File Validation**: Type and size validation
- **CDN Support**: Optional CDN URL configuration
- **Auto Update**: Automatically updates user profile
- **Secure**: Public-read ACL for easy access

### API Endpoints

#### 1. Upload Profile Picture (Entativa)
```http
POST /api/v1/media/profile-picture
Authorization: Bearer <access_token>
Content-Type: multipart/form-data

file=@profile.jpg
```

**Response:**
```json
{
  "success": true,
  "message": "Profile picture uploaded successfully",
  "data": {
    "url": "https://cdn.example.com/profile-pictures/user-id/file.jpg"
  }
}
```

#### 2. Upload Cover Photo (Entativa)
```http
POST /api/v1/media/cover-photo
Authorization: Bearer <access_token>
Content-Type: multipart/form-data

file=@cover.jpg
```

### Supported Formats

- **Image Types**: JPG, JPEG, PNG, GIF, WEBP
- **Max Size (Profile)**: 5MB
- **Max Size (Cover)**: 10MB

### Configuration

```env
S3_ENDPOINT=http://localhost:9000           # MinIO or S3 endpoint
S3_ACCESS_KEY_ID=your-access-key           # S3 access key
S3_SECRET_ACCESS_KEY=your-secret-key       # S3 secret key
S3_BUCKET_NAME=entativa-media             # Bucket name
S3_REGION=us-east-1                        # AWS region
S3_CDN_URL=https://cdn.example.com         # Optional CDN
S3_USE_PATH_STYLE=true                     # true for MinIO
```

### Usage Example

```bash
# Upload profile picture
curl -X POST http://localhost:8001/api/v1/media/profile-picture \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@/path/to/picture.jpg"

# Upload cover photo
curl -X POST http://localhost:8001/api/v1/media/cover-photo \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@/path/to/cover.jpg"
```

---

## Rate Limiting

### Overview

Redis-powered intelligent rate limiting to prevent abuse and ensure fair usage.

### Features

- **Per-IP Limiting**: Track requests by IP address
- **Per-User Limiting**: Track authenticated user requests
- **Endpoint-Specific**: Different limits for different endpoints
- **Headers**: Rate limit info in response headers
- **Graceful Degradation**: Falls back to no limiting if Redis unavailable

### Rate Limits

| Endpoint | Limit | Window |
|----------|-------|--------|
| `/auth/signup` | 3 requests | 1 hour per IP |
| `/auth/login` | 5 requests | 15 minutes per IP |
| `/auth/*` | 100 requests | 1 minute per user |
| General API | 100 requests | 1 minute per user |

### Response Headers

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1634567890
```

### Rate Limit Exceeded Response

```json
{
  "error": "Rate limit exceeded",
  "message": "Too many requests. Please try again in 15m0s"
}
```

### Configuration

```env
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

---

## Redis Caching

### Overview

Redis-based caching layer for improved performance and reduced database load.

### Features

- **Session Caching**: Cache user sessions
- **User Profile Caching**: Cache frequently accessed profiles
- **Token Validation**: Cache JWT validation results
- **TTL Support**: Configurable expiration times
- **Automatic Invalidation**: Cache invalidation on updates

### Cache Keys

```
user:profile:{user_id}           # User profile cache
session:{session_id}             # Session cache
rate_limit:{path}:{identifier}   # Rate limit counters
token:validation:{token_hash}    # Token validation cache
```

### Cache Strategies

1. **Cache-Aside**: Application manages cache
2. **TTL-Based Expiration**: Automatic cache expiration
3. **Write-Through**: Update cache on data changes
4. **Cache Invalidation**: Manual invalidation when needed

---

## gRPC Communication

### Overview

gRPC server for high-performance microservice communication.

### Features

- **Protocol Buffers**: Efficient binary serialization
- **Type Safety**: Strongly typed contracts
- **Streaming**: Bi-directional streaming support
- **Performance**: Lower latency than REST
- **Service Discovery**: Ready for service mesh integration

### Protobuf Definition

```protobuf
syntax = "proto3";

package entativa.user;

service UserService {
  rpc GetUser(GetUserRequest) returns (UserResponse);
  rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
  rpc GetUsersByIDs(GetUsersByIDsRequest) returns (GetUsersByIDsResponse);
  rpc VerifyTwoFactor(VerifyTwoFactorRequest) returns (VerifyTwoFactorResponse);
}
```

### Configuration

```env
GRPC_PORT=9001                 # Entativa gRPC port
GRPC_ENABLED=true              # Enable/disable gRPC server
```

### Ports

- **Entativa gRPC**: Port 9001
- **Vignette gRPC**: Port 9002

### Usage (Go Client Example)

```go
conn, err := grpc.Dial("localhost:9001", grpc.WithInsecure())
client := pb.NewUserServiceClient(conn)

resp, err := client.GetUser(ctx, &pb.GetUserRequest{
    UserId: "user-uuid",
})
```

---

## Kafka Event Streaming

### Overview

Apache Kafka integration for event-driven architecture and async processing.

### Features

- **Event Publishing**: Publish user events to Kafka
- **Async Processing**: Non-blocking event publishing
- **Topic-Based**: Organized by event topics
- **Scalable**: Kafka's distributed architecture
- **Optional**: Can be disabled if not needed

### Events Published

#### 1. Password Reset Requested
```json
{
  "event_type": "password_reset_requested",
  "user_id": "uuid",
  "email": "user@example.com",
  "token": "reset-token",
  "expires_at": "2025-10-15T12:00:00Z",
  "timestamp": "2025-10-15T11:00:00Z"
}
```

#### 2. Password Changed
```json
{
  "event_type": "password_changed",
  "user_id": "uuid",
  "email": "user@example.com",
  "timestamp": "2025-10-15T11:00:00Z"
}
```

#### 3. User Signup (Future)
```json
{
  "event_type": "user_signup",
  "user_id": "uuid",
  "email": "user@example.com",
  "timestamp": "2025-10-15T11:00:00Z"
}
```

### Topics

- `user-events`: All user-related events
- `auth-events`: Authentication events (future)
- `media-events`: Media upload events (future)

### Configuration

```env
KAFKA_BROKERS=localhost:9092
KAFKA_ENABLED=false            # Enable/disable Kafka
```

### Consumer Example (Python)

```python
from kafka import KafkaConsumer
import json

consumer = KafkaConsumer(
    'user-events',
    bootstrap_servers=['localhost:9092'],
    value_deserializer=lambda m: json.loads(m.decode('utf-8'))
)

for message in consumer:
    event = message.value
    if event['event_type'] == 'password_reset_requested':
        send_password_reset_email(event)
```

---

## Configuration

### Complete Environment Variables

```env
# Server
PORT=8001
ENVIRONMENT=development
ALLOWED_ORIGINS=*

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=entativa
DB_PASSWORD=entativa_password
DB_NAME=entativa_users
DB_SSL_MODE=disable
DB_MAX_CONNECTIONS=100
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=3600

# Redis (Caching & Rate Limiting)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=change-this-in-production
JWT_ACCESS_TOKEN_TTL=86400
JWT_REFRESH_TOKEN_TTL=604800

# S3/MinIO (Media Uploads)
S3_ENDPOINT=http://localhost:9000
S3_ACCESS_KEY_ID=
S3_SECRET_ACCESS_KEY=
S3_BUCKET_NAME=entativa-media
S3_REGION=us-east-1
S3_CDN_URL=
S3_USE_PATH_STYLE=true

# Kafka (Event Streaming)
KAFKA_BROKERS=localhost:9092
KAFKA_ENABLED=false

# gRPC (Microservices)
GRPC_PORT=9001
GRPC_ENABLED=true
```

### Feature Toggles

All advanced features can be individually enabled/disabled:

- **Redis**: Set empty `REDIS_HOST` to disable
- **S3/MinIO**: Set empty `S3_ACCESS_KEY_ID` to disable
- **Kafka**: Set `KAFKA_ENABLED=false` to disable
- **gRPC**: Set `GRPC_ENABLED=false` to disable

### Production Recommendations

1. **Redis**: Use Redis Cluster for high availability
2. **S3**: Use CloudFront or similar CDN
3. **Kafka**: Use managed Kafka (Confluent, AWS MSK)
4. **gRPC**: Use service mesh (Istio, Linkerd)
5. **Secrets**: Use secret management (Vault, AWS Secrets Manager)

---

## Testing

### Test 2FA Flow

```bash
# 1. Setup 2FA
TOKEN=$(curl -s -X POST http://localhost:8001/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email_or_username":"user@example.com","password":"password"}' \
  | jq -r '.data.access_token')

# 2. Get QR code
curl -X POST http://localhost:8001/api/v1/auth/2fa/setup \
  -H "Authorization: Bearer $TOKEN" | jq '.data.qr_code_url'

# 3. Enable 2FA with code from authenticator app
curl -X POST http://localhost:8001/api/v1/auth/2fa/enable \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"code":"123456"}'
```

### Test Password Reset

```bash
# 1. Request reset
curl -X POST http://localhost:8001/api/v1/auth/password-reset/request \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com"}'

# 2. Reset password (get token from Kafka or database)
curl -X POST http://localhost:8001/api/v1/auth/password-reset/confirm \
  -H "Content-Type: application/json" \
  -d '{"token":"reset-token","new_password":"NewPassword123"}'
```

### Test Rate Limiting

```bash
# Trigger rate limit
for i in {1..10}; do
  curl -X POST http://localhost:8001/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email_or_username":"test","password":"test"}'
  echo ""
done
```

---

## Troubleshooting

### Redis Connection Failed
- Check Redis is running: `redis-cli ping`
- Verify Redis host/port in `.env`
- Services will run without Redis but without caching/rate limiting

### S3 Upload Failed
- Verify S3 credentials
- Check bucket exists and has correct permissions
- For MinIO, ensure `S3_USE_PATH_STYLE=true`

### Kafka Not Connected
- Check Kafka is running: `kafka-topics.sh --list --bootstrap-server localhost:9092`
- Verify `KAFKA_BROKERS` in `.env`
- Services will run without Kafka but events won't be published

### gRPC Port Conflict
- Change `GRPC_PORT` if 9001/9002 are in use
- Or disable gRPC: `GRPC_ENABLED=false`

---

## Summary

Both Entativa and Vignette now include:

✅ **Enterprise Security**: 2FA, password reset, account recovery  
✅ **Performance**: Redis caching and rate limiting  
✅ **Scalability**: gRPC and Kafka for microservices  
✅ **Media**: S3/MinIO integration  
✅ **Flexibility**: All features are optional and configurable  

**Total New API Endpoints**: 8 per service  
**Total New Database Tables**: 2 per service  
**Total New Features**: 8 major features  

All features maintain the Meta-level philosophy of instant access while adding optional security enhancements!

---

**For more details, see the main documentation in the service README files.**
