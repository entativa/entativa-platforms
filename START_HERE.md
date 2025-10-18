# 🚀 START HERE - Auth System Quick Start

**Real talk:** This is a complete, working authentication system. No stubs, no TODOs. Here's how to get it running in 10 minutes.

---

## ⚡ Super Quick Start (10 Minutes)

### 1. Prerequisites Check
```bash
# Check if you have everything (30 seconds)
which psql     # PostgreSQL
which go       # Go 1.21+
which jq       # JSON processor (for testing)

# If missing:
brew install postgresql go jq  # macOS
# or
apt install postgresql golang-go jq  # Linux
```

### 2. Start PostgreSQL
```bash
# macOS
brew services start postgresql@14

# Linux
sudo systemctl start postgresql

# Docker (easiest)
docker run --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:14
```

### 3. Setup Entativa Backend (3 minutes)
```bash
cd /workspace/EntativaBackend/services/user-service

# Run automated setup
./scripts/setup-dev.sh

# This will:
# - Create .env file
# - Install dependencies
# - Create database
# - Run migrations
# - Build the app
```

### 4. Start Entativa Backend (instant)
```bash
make run

# You should see:
# INFO: Starting entativa-user-service on port 8001
# INFO: Connected to database: entativa_users
# INFO: Server listening on port 8001
```

### 5. Setup Vignette Backend (3 minutes)
```bash
# Open new terminal
cd /workspace/VignetteBackend/services/user-service

./scripts/setup-dev.sh
```

### 6. Start Vignette Backend (instant)
```bash
make run

# You should see:
# INFO: Starting vignette-user-service on port 8002
# INFO: Server listening on port 8002
```

### 7. Test Everything (2 minutes)
```bash
# Open new terminal
cd /workspace

# Run automated tests
./test-auth-complete.sh

# You should see all green checkmarks:
# ✅ Health checks passed
# ✅ Entativa sign up successful
# ✅ Vignette sign up successful
# ✅ Cross-platform SSO successful
# ... and more
# 🎉 All tests passed!
```

### 8. Test Mobile Apps (optional)
```bash
# iOS - Entativa
cd /workspace/EntativaiOS
open Entativa.xcodeproj
# Press Cmd+R in Xcode

# iOS - Vignette
cd /workspace/VignetteiOS
open Vignette.xcodeproj
# Press Cmd+R in Xcode

# Android - Entativa
cd /workspace/EntativaAndroid
./gradlew installDebug

# Android - Vignette
cd /workspace/VignetteAndroid
./gradlew installDebug
```

---

## 🎯 What You Can Do Right Now

### cURL Examples

**Sign up on Entativa:**
```bash
curl -X POST http://localhost:8001/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Neo",
    "last_name": "Qiss",
    "email": "neo@entativa.com",
    "password": "Secure123",
    "birthday": "1990-01-01",
    "gender": "prefer_not_to_say"
  }'

# Response: Your access token + user data
```

**Login:**
```bash
curl -X POST http://localhost:8001/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email_or_username": "neo@entativa.com",
    "password": "Secure123"
  }'
```

**Sign in with Vignette (cross-platform):**
```bash
# First create Vignette account
curl -X POST http://localhost:8002/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "username": "neoqiss",
    "email": "neo@vignette.app",
    "full_name": "Neo Qiss",
    "password": "Secure123"
  }'

# Get the token from response, then:
curl -X POST http://localhost:8001/api/v1/auth/cross-platform/signin \
  -H "Content-Type: application/json" \
  -d '{
    "platform": "vignette",
    "access_token": "YOUR_VIGNETTE_TOKEN"
  }'

# Magic! Now you have Entativa account too!
```

---

## 📱 Mobile App Features

### Entativa App
- **Multi-step sign-up** (name → email/pass → birthday/gender)
- **Password strength meter** (live updates)
- **Age validation** (must be 13+)
- **Gender selection** with chips
- **Sign in with Vignette** button
- **Biometric login** (Face ID/Touch ID)
- **Forgot password** flow

### Vignette App
- **Single-page sign-up** (all fields at once)
- **Username validation** (Instagram rules)
- **Auto-lowercase** username
- **Real-time hints** ("Can contain letters, numbers...")
- **Sign in with Entativa** button
- **Minimal design** (Instagram-style)
- **Password requirements** with checkmarks

---

## 🔧 Troubleshooting

### "Connection refused"
```bash
# Check if services are running
lsof -i :8001  # Entativa
lsof -i :8002  # Vignette

# If not, start them:
cd EntativaBackend/services/user-service && make run
```

### "Database does not exist"
```bash
# Run setup script again
cd EntativaBackend/services/user-service
./scripts/setup-dev.sh

# Or manually create:
createdb entativa_users
createdb vignette_users
```

### "Package not found" (Go)
```bash
# In the service directory:
go mod download
go mod tidy
```

### iOS won't build
```bash
# Clean Xcode cache
rm -rf ~/Library/Developer/Xcode/DerivedData/*

# In Xcode: Product → Clean Build Folder
```

### Android won't build
```bash
./gradlew clean
./gradlew build --refresh-dependencies
```

---

## 📚 Documentation

Read these in order:

1. **START_HERE.md** (this file) - Quick start
2. **COMPLETE_SETUP_GUIDE.md** - Detailed setup instructions
3. **REAL_FINAL_STATUS.md** - What's actually done
4. **AUTH_SYSTEM_COMPLETE.md** - Technical deep dive
5. **HONEST_IMPLEMENTATION_STATUS.md** - No-BS status

---

## ✅ Verification Checklist

Run through this to verify everything works:

```bash
# 1. Health checks
curl http://localhost:8001/health
curl http://localhost:8002/health

# 2. Sign up
curl -X POST http://localhost:8001/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"first_name":"Test","last_name":"User","email":"test@example.com","password":"Test1234","birthday":"1995-01-01","gender":"male"}'

# 3. Login  
curl -X POST http://localhost:8001/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email_or_username":"test@example.com","password":"Test1234"}'

# 4. Forgot password
curl -X POST http://localhost:8001/api/v1/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com"}'
```

All should return `{"success": true, ...}`

---

## 🎉 You're Ready!

Everything is implemented and working. Just:

1. **Run the setup scripts** (already created)
2. **Start the servers** (one command each)
3. **Test with the script** (automated)
4. **Build mobile apps** (Xcode/Android Studio)

**No configuration needed beyond basic .env setup!**

---

**Questions?** Everything is documented in the code itself. Check the files! 📖

**Want to see it work?** Run `./test-auth-complete.sh` right now! 🚀

**Ready to ship?** Read `COMPLETE_SETUP_GUIDE.md` for production deployment! 🎯
