# 🚀 Meta-Level Authentication Services - Quick Start Guide

This guide will help you get both Socialink and Vignette authentication services up and running in minutes.

## 📋 Prerequisites

Before you begin, ensure you have the following installed:

- **Go**: Version 1.21 or higher ([Install Go](https://golang.org/dl/))
- **PostgreSQL**: Version 12 or higher ([Install PostgreSQL](https://www.postgresql.org/download/))
- **curl**: For testing APIs
- **jq** (optional): For pretty-printing JSON responses

## 🗄️ Database Setup

### 1. Create Databases

```bash
# Connect to PostgreSQL
psql -U postgres

# Create databases
CREATE DATABASE socialink_users;
CREATE DATABASE vignette_users;

# Create users (optional, for production)
CREATE USER socialink WITH PASSWORD 'socialink_password';
CREATE USER vignette WITH PASSWORD 'vignette_password';

# Grant privileges
GRANT ALL PRIVILEGES ON DATABASE socialink_users TO socialink;
GRANT ALL PRIVILEGES ON DATABASE vignette_users TO vignette;

# Exit psql
\q
```

### 2. Database Migrations

The services automatically run migrations on startup. Tables will be created automatically:
- ✅ `users` table
- ✅ `sessions` table
- ✅ Indexes for performance

## 🔵 Starting Socialink Service

### Method 1: Using the Run Script (Recommended)

```bash
cd /workspace/SocialinkBackend/services/user-service
./run.sh
```

### Method 2: Manual Start

```bash
cd /workspace/SocialinkBackend/services/user-service

# Copy environment file
cp .env.example .env

# Edit .env with your database credentials
# nano .env

# Download dependencies
go mod download

# Run the service
go run cmd/api/main.go
```

**Service will start on**: `http://localhost:8001`

**Logs should show**:
```
✓ Connected to PostgreSQL database
🚀 Socialink User Service starting on port 8001
📝 Environment: development
🔐 Meta-level authentication enabled (instant access, no verification)
✨ Ready to accept connections!
```

## 🟣 Starting Vignette Service

### Method 1: Using the Run Script (Recommended)

```bash
cd /workspace/VignetteBackend/services/user-service
./run.sh
```

### Method 2: Manual Start

```bash
cd /workspace/VignetteBackend/services/user-service

# Copy environment file
cp .env.example .env

# Edit .env with your database credentials
# nano .env

# Download dependencies
go mod download

# Run the service
go run cmd/api/main.go
```

**Service will start on**: `http://localhost:8002`

## 🧪 Testing the Services

### Quick Health Check

**Socialink**:
```bash
curl http://localhost:8001/health
```

**Vignette**:
```bash
curl http://localhost:8002/health
```

### Comprehensive API Tests

**Test Socialink** (Signup, Login, Get User, Logout):
```bash
/workspace/test-socialink-auth.sh
```

**Test Vignette** (Signup, Login, Get User, Logout):
```bash
/workspace/test-vignette-auth.sh
```

## 📝 Manual API Testing

### Socialink - Create Account

```bash
curl -X POST http://localhost:8001/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Jane",
    "last_name": "Smith",
    "email": "jane.smith@example.com",
    "password": "MySecurePass123",
    "birthday": "1998-03-20",
    "gender": "female"
  }'
```

### Socialink - Login

```bash
curl -X POST http://localhost:8001/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email_or_username": "jane.smith@example.com",
    "password": "MySecurePass123"
  }'
```

### Vignette - Create Account

```bash
curl -X POST http://localhost:8002/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "username": "janesmith",
    "email": "jane.smith@example.com",
    "full_name": "Jane Smith",
    "password": "MySecurePass123"
  }'
```

### Vignette - Login

```bash
curl -X POST http://localhost:8002/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username_or_email": "janesmith",
    "password": "MySecurePass123"
  }'
```

### Get Current User (Protected Endpoint)

Replace `<ACCESS_TOKEN>` with the token from signup/login response:

**Socialink**:
```bash
curl -X GET http://localhost:8001/api/v1/auth/me \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

**Vignette**:
```bash
curl -X GET http://localhost:8002/api/v1/auth/me \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

## 🐳 Docker Deployment (Alternative)

### Build Docker Images

**Socialink**:
```bash
cd /workspace/SocialinkBackend/services/user-service
docker build -t socialink-user-service .
```

**Vignette**:
```bash
cd /workspace/VignetteBackend/services/user-service
docker build -t vignette-user-service .
```

### Run with Docker

**Socialink**:
```bash
docker run -p 8001:8001 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=5432 \
  -e DB_USER=socialink \
  -e DB_PASSWORD=socialink_password \
  -e DB_NAME=socialink_users \
  socialink-user-service
```

**Vignette**:
```bash
docker run -p 8002:8002 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=5432 \
  -e DB_USER=vignette \
  -e DB_PASSWORD=vignette_password \
  -e DB_NAME=vignette_users \
  vignette-user-service
```

## 🔧 Configuration

### Environment Variables

Both services use the same configuration structure. See `.env.example` for all options:

**Key Variables**:
- `PORT` - Service port (8001 for Socialink, 8002 for Vignette)
- `DB_HOST` - PostgreSQL host
- `DB_PORT` - PostgreSQL port
- `DB_USER` - Database user
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name
- `JWT_SECRET` - Secret key for JWT signing (CHANGE IN PRODUCTION!)
- `JWT_ACCESS_TOKEN_TTL` - Token lifetime in seconds (default: 86400 = 24h)

### Production Considerations

For production deployment:

1. **Change JWT Secret**: Use a strong, random secret key
2. **Enable SSL**: Set `DB_SSL_MODE=require`
3. **Set Environment**: `ENVIRONMENT=production`
4. **Configure CORS**: Set specific `ALLOWED_ORIGINS`
5. **Use HTTPS**: Deploy behind a reverse proxy (nginx, Caddy)
6. **Database**: Use connection pooling and read replicas
7. **Monitoring**: Add Prometheus metrics
8. **Logging**: Configure structured logging

## 📊 API Endpoints Summary

### Socialink (Port 8001)

| Method | Endpoint | Auth Required | Description |
|--------|----------|---------------|-------------|
| GET | `/health` | No | Health check |
| POST | `/api/v1/auth/signup` | No | Create account |
| POST | `/api/v1/auth/login` | No | Login |
| GET | `/api/v1/auth/me` | Yes | Get current user |
| POST | `/api/v1/auth/logout` | Yes | Logout |

### Vignette (Port 8002)

| Method | Endpoint | Auth Required | Description |
|--------|----------|---------------|-------------|
| GET | `/health` | No | Health check |
| POST | `/api/v1/auth/signup` | No | Create account |
| POST | `/api/v1/auth/login` | No | Login |
| GET | `/api/v1/auth/me` | Yes | Get current user |
| POST | `/api/v1/auth/logout` | Yes | Logout |

## 🔍 Troubleshooting

### Service won't start

**Check PostgreSQL is running**:
```bash
pg_isready
```

**Check if port is already in use**:
```bash
lsof -i :8001  # for Socialink
lsof -i :8002  # for Vignette
```

### Database connection errors

**Verify database exists**:
```bash
psql -U postgres -l | grep socialink
psql -U postgres -l | grep vignette
```

**Check database credentials in `.env`**

### Cannot create user (email/username exists)

This is expected if you're testing multiple times. Either:
- Use a different email/username
- Clear the database: `TRUNCATE users, sessions CASCADE;`

## 📚 Additional Resources

- **Socialink README**: `/workspace/SocialinkBackend/services/user-service/README.md`
- **Vignette README**: `/workspace/VignetteBackend/services/user-service/README.md`
- **Implementation Summary**: `/workspace/META_AUTH_IMPLEMENTATION_SUMMARY.md`

## 🎉 Success Indicators

You'll know everything is working when:

✅ Services start without errors  
✅ Health check returns `{"status": "healthy"}`  
✅ You can create a new account (signup)  
✅ You receive an access token  
✅ You can access protected endpoints with the token  
✅ You can login with existing credentials  
✅ Database shows user and session records  

## 🚀 Next Steps

Once both services are running:

1. **Integrate with Frontend**: Use the access tokens for authenticated requests
2. **Add More Features**: Profile updates, password reset, etc.
3. **Connect Other Microservices**: User service provides authentication for all services
4. **Set Up API Gateway**: Route requests through a unified gateway
5. **Add Monitoring**: Prometheus, Grafana for observability
6. **Deploy to Production**: Kubernetes, AWS, GCP, or your preferred platform

## 💡 Pro Tips

- Keep your JWT secret secure and rotate it periodically
- Monitor session table size and clean up expired sessions
- Use Redis for session caching in production
- Implement rate limiting for auth endpoints
- Add logging and monitoring from day one
- Test authentication flows thoroughly before going live

---

**Need Help?** Check the README files or review the implementation summary document!

**Happy Building!** 🚀✨
