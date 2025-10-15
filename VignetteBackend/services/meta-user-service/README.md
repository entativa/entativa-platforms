# Meta User Service - Vignette

## Overview

The Meta User Service is a PhD-level engineered microservice that provides advanced cross-platform user management for the Meta ecosystem (Socialink and Vignette platforms). It implements sophisticated features including:

- **Cross-Platform Identity Management**: Unified user identity across Socialink and Vignette
- **ML-Based Fraud Detection**: Real-time fraud detection using machine learning algorithms
- **Behavioral Anomaly Detection**: Advanced anomaly detection for suspicious user activities
- **Dynamic Trust Scoring**: ML-powered trust score calculation with multi-factor analysis
- **Advanced Security**: Multi-factor authentication, biometrics, device fingerprinting, and more
- **Event Sourcing**: Complete audit trail with event sourcing pattern
- **GDPR/CCPA Compliance**: Advanced data rights management and privacy controls
- **Cross-Platform Synchronization**: Real-time data sync between platforms
- **Circuit Breaker Pattern**: Resilient service communication
- **Multi-Level Caching**: Redis-based caching with intelligent invalidation

## Architecture

### Domain-Driven Design (DDD)

The service follows DDD principles with clear separation of concerns:

- **Domain Layer**: Core business models and logic (`internal/model`)
- **Application Layer**: Service orchestration (`internal/service`)
- **Infrastructure Layer**: Database, cache, messaging (`pkg/database`, `pkg/cache`, `pkg/kafka`)
- **Interface Layer**: HTTP handlers and API (`internal/handler`)

### Machine Learning Components

#### Fraud Detector
- Multi-factor fraud analysis during signup
- Email pattern analysis
- IP reputation scoring
- Device fingerprint analysis
- Velocity attack detection
- Ensemble scoring with weighted features

#### Behavior Analyzer
- Location-based anomaly detection
- Impossible travel detection
- Temporal pattern analysis
- Device anomaly detection
- Behavioral baseline deviation analysis

#### Trust Score Engine
- Multi-factor trust scoring:
  - Account age (15%)
  - Verification status (20%)
  - Activity score (15%)
  - Violation history (15%)
  - Positive interactions (15%)
  - Device trust (10%)
  - Security compliance (5%)
  - Anomaly penalty (5%)
- Non-linear transformation for better distribution
- Explainable AI with score breakdown

### Event Sourcing

All user state changes are captured as events in the `meta_user_events` table:
- Complete audit trail
- Event replay capability
- Temporal queries
- Compliance support

## API Endpoints

### User Management

```
POST   /api/v1/users                    # Create user
POST   /api/v1/users/authenticate       # Authenticate user
GET    /api/v1/users/:id                # Get user details
POST   /api/v1/users/:id/trust-score    # Update trust score
```

### Platform Linking

```
POST   /api/v1/platforms/link           # Link platform account
```

### Cross-Platform Sync

```
POST   /api/v1/sync/:id/enable          # Enable sync
POST   /api/v1/sync/:id/disable         # Disable sync
```

### Privacy Settings

```
PUT    /api/v1/privacy/:id/settings     # Update privacy settings
```

## Configuration

Environment variables:

```bash
# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
ENVIRONMENT=production

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-password
DB_NAME=meta_users
DB_SSL_MODE=require
DB_MAX_CONNS=100
DB_MIN_CONNS=10

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_POOL_SIZE=100

# Kafka
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=meta-user-events
KAFKA_GROUP_ID=meta-user-service

# Security
JWT_SECRET=your-secret-key
JWT_EXPIRATION=3600
BCRYPT_COST=10

# ML Configuration
ML_FRAUD_THRESHOLD=0.7
ML_ANOMALY_THRESHOLD=0.75
ML_TRUST_SCORE_VERSION=v1.0
```

## Database Schema

### meta_users Table

Primary table storing unified user data:
- Core identity (ID, meta_id, email, phone)
- Security profile (2FA, biometrics, security keys)
- Privacy settings (visibility, data sharing)
- Platform links (Socialink, Vignette)
- Device fingerprints
- Anomaly detection data
- Trust score and risk level
- Compliance data (GDPR, CCPA)

### meta_user_events Table

Event sourcing table:
- Event ID and type
- Event data (JSON)
- Timestamp
- User reference

## Running the Service

### Local Development

```bash
# Install dependencies
go mod download

# Run migrations
psql -U postgres -d meta_users -f migrations/001_create_meta_users_table.up.sql

# Run the service
go run cmd/api/main.go
```

### Docker

```bash
# Build image
docker build -t vignette-meta-user-service .

# Run container
docker run -p 8080:8080 \
  -e DB_HOST=postgres \
  -e REDIS_HOST=redis \
  -e KAFKA_BROKERS=kafka:9092 \
  vignette-meta-user-service
```

## Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run integration tests
go test -tags=integration ./test/...
```

## Monitoring

The service exposes Prometheus metrics at `/metrics`:
- Request rates
- Error rates
- Response times
- Trust score distributions
- Fraud detection metrics
- Anomaly detection metrics

## Security Features

### Multi-Factor Authentication
- TOTP (Time-based One-Time Password)
- SMS verification
- Email verification
- Biometric authentication
- Hardware security keys (WebAuthn/FIDO2)

### Device Fingerprinting
- Browser fingerprinting
- Device ID tracking
- IP address monitoring
- Geolocation tracking
- Behavioral patterns

### Anomaly Detection
- Real-time behavioral analysis
- Impossible travel detection
- Unusual login patterns
- Mass action detection
- Scraping detection

## Cross-Platform Features

### Account Linking
- Link Socialink and Vignette accounts
- Unified authentication
- Single sign-on (SSO)
- Cross-platform profile sync

### Data Synchronization
- Real-time data sync
- Conflict resolution
- Privacy settings sync
- Activity tracking

## Compliance

### GDPR
- Right to access
- Right to rectification
- Right to erasure
- Right to portability
- Data export functionality

### CCPA
- Opt-out mechanisms
- Data transparency
- Privacy controls

## Performance

### Caching Strategy
- L1: Redis cache (15-minute TTL)
- Cache invalidation on updates
- Pattern-based cache clearing

### Database Optimization
- Connection pooling
- Indexed queries
- JSON field optimization
- Event sourcing for audit

### Scalability
- Horizontal scaling supported
- Stateless design
- Distributed caching
- Event-driven architecture

## Future Enhancements

- [ ] GraphQL API
- [ ] Advanced ML models (TensorFlow integration)
- [ ] Real-time WebSocket support
- [ ] Advanced biometric authentication
- [ ] Blockchain-based identity verification
- [ ] Zero-knowledge proof authentication
- [ ] Federated learning for privacy-preserving ML
- [ ] Advanced encryption (homomorphic encryption)
- [ ] Quantum-resistant cryptography

## License

Proprietary - Meta Platforms Inc.
