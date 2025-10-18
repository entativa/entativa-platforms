# gRPC Integration - Media Service ✅

## Complete Service-to-Service Communication

**Status**: ✅ **PRODUCTION-READY gRPC Integration**

---

## 📊 What Was Implemented

### 1. Protocol Buffers Definition (`media_service.proto`)
**Location**: `services/media-service/proto/media_service.proto`

#### Service Methods (8 RPCs)
```protobuf
service MediaService {
  rpc UploadMedia(UploadMediaRequest) returns (UploadMediaResponse);
  rpc UploadMediaStream(stream UploadChunk) returns (UploadMediaResponse);
  rpc GetMedia(GetMediaRequest) returns (MediaResponse);
  rpc DeleteMedia(DeleteMediaRequest) returns (DeleteMediaResponse);
  rpc ProcessMedia(ProcessMediaRequest) returns (ProcessMediaResponse);
  rpc GetSignedUrl(GetSignedUrlRequest) returns (GetSignedUrlResponse);
  rpc BatchGetMedia(BatchGetMediaRequest) returns (BatchGetMediaResponse);
  rpc BatchDeleteMedia(BatchDeleteMediaRequest) returns (BatchDeleteMediaResponse);
}
```

#### Key Enums
- `MediaType`: IMAGE, VIDEO, AUDIO, DOCUMENT
- `MediaPurpose`: PROFILE_PICTURE, COVER_PHOTO, POST_ATTACHMENT, MESSAGE_ATTACHMENT, STORY, AVATAR
- `ProcessingStatus`: PENDING, PROCESSING, COMPLETED, FAILED, CANCELLED

---

### 2. Rust gRPC Server (`grpc_server.rs`)
**Location**: `services/media-service/src/grpc_server.rs`
**Lines**: 450+ lines of production code

#### Features
- ✅ Full implementation of all 8 RPC methods
- ✅ Database integration (PostgreSQL)
- ✅ Redis caching
- ✅ Storage backend abstraction (S3/MinIO/Local)
- ✅ Automatic thumbnail generation
- ✅ Blurhash generation
- ✅ File validation
- ✅ Deduplication
- ✅ Ownership verification
- ✅ Batch operations

#### Server Configuration
```rust
// Dual server setup: HTTP + gRPC
- HTTP Server: Port 8083
- gRPC Server: Port 50051
```

---

### 3. Go gRPC Client (`media/client.go`)
**Location**: `services/user-service/pkg/media/client.go`
**Lines**: 250+ lines

#### Client Methods
```go
- NewClient(address string) (*Client, error)
- UploadMedia(ctx, req) (*UploadMediaResponse, error)
- UploadMediaStream(ctx, chunks) (*UploadMediaResponse, error)
- GetMedia(ctx, mediaID) (*MediaResponse, error)
- DeleteMedia(ctx, mediaID, userID) error
- GetSignedUrl(ctx, mediaID, expirySeconds) (string, error)
- BatchGetMedia(ctx, mediaIDs) ([]*MediaResponse, error)
- BatchDeleteMedia(ctx, mediaIDs, userID) (int32, []string, error)
```

#### Convenience Methods
```go
- UploadProfilePicture(ctx, data, filename, userID) (*UploadMediaResponse, error)
- UploadCoverPhoto(ctx, data, filename, userID) (*UploadMediaResponse, error)
```

---

### 4. Media Handler (`media_handler.go`)
**Location**: `services/user-service/internal/handler/media_handler.go`
**Lines**: 250+ lines

#### HTTP Endpoints
```
POST   /api/v1/media/profile-picture  - Upload profile picture
POST   /api/v1/media/cover-photo      - Upload cover photo
GET    /api/v1/media/:media_id        - Get media info
DELETE /api/v1/media/:media_id        - Delete media
```

#### Features
- ✅ JWT authentication
- ✅ File size validation (10MB profile, 20MB cover)
- ✅ Multipart form handling
- ✅ Automatic thumbnail generation
- ✅ Error handling

---

## 🏗️ Architecture

### Service Communication Flow
```
┌─────────────────┐     gRPC      ┌──────────────────┐
│  User Service   │◄─────────────►│  Media Service   │
│  (Go - Gin)     │   Port 50051   │  (Rust - Actix)  │
│                 │                │                  │
│  Profile Pics   │                │  Image Process   │
│  Cover Photos   │                │  Video Process   │
│  Posts          │                │  Storage (S3)    │
└─────────────────┘                └──────────────────┘
         │                                  │
         │                                  │
         ▼                                  ▼
   ┌─────────┐                        ┌─────────┐
   │   JWT   │                        │  Redis  │
   │  Auth   │                        │  Cache  │
   └─────────┘                        └─────────┘
         │                                  │
         └──────────► PostgreSQL ◄──────────┘
```

### Profile Picture Upload Flow
```
1. User uploads image via HTTP (multipart/form-data)
   ↓
2. User Service receives file
   ↓
3. User Service validates file (size, type)
   ↓
4. User Service calls Media Service via gRPC
   ↓
5. Media Service:
   - Validates file
   - Detects MIME type
   - Generates unique filename
   - Uploads to S3/storage
   - Generates thumbnails (150px, 300px, 600px)
   - Generates blurhash
   - Saves to PostgreSQL
   - Caches in Redis
   ↓
6. Returns media_id, URLs, thumbnails
   ↓
7. User Service updates profile with media URLs
```

---

## 🔧 Configuration

### Rust Media Service
```toml
# Cargo.toml additions
tonic = "0.10"
tonic-build = "0.10"
prost = "0.12"
prost-types = "0.12"
```

```rust
// config.rs
pub struct ServerConfig {
    pub host: String,
    pub port: u16,          // HTTP: 8083
    pub grpc_port: u16,     // gRPC: 50051
    pub workers: usize,
    pub max_connections: usize,
}
```

### Go User Service
```go
// go.mod already has:
google.golang.org/grpc v1.60.0
google.golang.org/protobuf v1.31.0
```

---

## 📝 Proto Generation

### For Rust (Automatic via build.rs)
```rust
// build.rs
fn main() -> Result<(), Box<dyn std::error::Error>> {
    tonic_build::configure()
        .build_server(true)
        .build_client(false)
        .compile(&["proto/media_service.proto"], &["proto"])?;
    Ok(())
}
```

### For Go (Manual Script)
```bash
# scripts/generate_proto.sh
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/media/media_service.proto
```

---

## 🚀 Usage Examples

### Profile Picture Upload (Go)
```go
// In handler
mediaClient, _ := media.NewClient("localhost:50051")
defer mediaClient.Close()

// Upload profile picture
resp, err := mediaClient.UploadProfilePicture(
    ctx,
    fileData,
    "avatar.jpg",
    userID,
)

// Response contains:
// - media_id
// - url (full size)
// - thumbnail_url (150px)
// - width, height
// - blurhash
```

### Cover Photo Upload (Go)
```go
// Upload cover photo
resp, err := mediaClient.UploadCoverPhoto(
    ctx,
    fileData,
    "cover.jpg",
    userID,
)

// Automatically generates thumbnails:
// - thumb: 400x150
// - medium: 800x300
// - large: 1600x600
```

### Direct gRPC Call (Go)
```go
req := &pb.UploadMediaRequest{
    Data:        fileData,
    Filename:    "photo.jpg",
    ContentType: "image/jpeg",
    UserId:      userID,
    Purpose:     pb.MediaPurpose_PROFILE_PICTURE,
    ProcessingOptions: &pb.ProcessingOptions{
        GenerateThumbnails: true,
        GenerateBlurhash:   true,
        Quality:            92,
        MaxWidth:           2048,
        MaxHeight:          2048,
    },
}

resp, err := mediaClient.UploadMedia(ctx, req)
```

---

## 🔐 Security

### Authentication
- User Service: JWT authentication on HTTP endpoints
- gRPC: User ID passed in request, verified by media service
- Ownership validation on delete/update operations

### Authorization
```rust
// Media service verifies ownership
if media.user_id != user_id {
    return Err(Status::permission_denied("Not authorized"));
}
```

### Validation
- File size limits (configurable)
- MIME type validation
- Magic byte verification
- Extension checking
- Content scanning ready

---

## 📊 Performance

### gRPC Benefits
- **Binary Protocol**: Faster than JSON/HTTP
- **HTTP/2**: Multiplexing, server push
- **Streaming**: Large file support via chunking
- **Code Generation**: Type-safe, no runtime overhead

### Optimizations
- Connection pooling
- Request timeouts (60s upload, 10s get)
- Message size limits (100MB max)
- Redis caching (1-hour TTL)
- Thumbnail generation async

---

## 🧪 Testing

### Test Profile Picture Upload
```bash
# Via HTTP API
curl -X POST http://localhost:8080/api/v1/media/profile-picture \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -F "file=@profile.jpg"

# Response:
{
  "success": true,
  "data": {
    "media_id": "uuid",
    "url": "https://storage.../profile.jpg",
    "thumbnail_url": "https://storage.../thumb.jpg",
    "width": 2048,
    "height": 2048,
    "blurhash": "L6PZfSi_.AyE_3t7t7R**0o#DgR4"
  }
}
```

### Test gRPC Directly
```bash
# Using grpcurl
grpcurl -plaintext \
  -d '{
    "filename": "test.jpg",
    "content_type": "image/jpeg",
    "user_id": "uuid",
    "purpose": 1,
    "data": "..."
  }' \
  localhost:50051 media.MediaService/UploadMedia
```

---

## 🔄 Integration Points

### Profile Service Integration
```go
// After uploading profile picture
resp, _ := mediaClient.UploadProfilePicture(...)

// Update user profile
profileService.UpdateProfileInfo(ctx, userID, &model.UpdateProfileInfoRequest{
    ProfilePictureURL: &resp.Url,
    ProfilePictureID:  &resp.MediaId,
})
```

### Cover Photo Integration
```go
// After uploading cover photo
resp, _ := mediaClient.UploadCoverPhoto(...)

// Update user profile
profileService.UpdateProfileInfo(ctx, userID, &model.UpdateProfileInfoRequest{
    CoverPhotoURL: &resp.Url,
    CoverPhotoID:  &resp.MediaId,
})
```

---

## 📈 Scalability

### Horizontal Scaling
- Multiple media service instances
- gRPC load balancing
- Shared PostgreSQL + Redis
- S3 for stateless storage

### Vertical Scaling
- Connection pooling
- Worker threads (CPU-based)
- Async I/O throughout
- Zero-copy operations

---

## 🎯 Next Steps

### For Posting Service
```go
// Post attachments will use same media service
mediaClient.UploadMedia(ctx, &pb.UploadMediaRequest{
    Purpose: pb.MediaPurpose_POST_ATTACHMENT,
    // ... file data
})
```

### For Messaging Service
```go
// Message attachments
mediaClient.UploadMedia(ctx, &pb.UploadMediaRequest{
    Purpose: pb.MediaPurpose_MESSAGE_ATTACHMENT,
    // ... file data
})
```

### For Stories Service
```go
// Story media
mediaClient.UploadMedia(ctx, &pb.UploadMediaRequest{
    Purpose: pb.MediaPurpose_STORY,
    // ... file data
})
```

---

## 📚 Files Created/Modified

### Rust Media Service
```
services/media-service/
├── proto/
│   └── media_service.proto          (NEW - 200 lines)
├── build.rs                         (NEW - 5 lines)
├── src/
│   ├── grpc_server.rs              (NEW - 450 lines)
│   ├── main.rs                     (MODIFIED - gRPC server startup)
│   └── config.rs                   (MODIFIED - grpc_port added)
└── Cargo.toml                      (MODIFIED - tonic deps)
```

### Go User Service
```
services/user-service/
├── proto/
│   └── media/
│       └── media_service.proto     (COPIED from media-service)
├── pkg/
│   └── media/
│       └── client.go               (NEW - 250 lines)
├── internal/
│   └── handler/
│       └── media_handler.go        (NEW - 250 lines)
└── scripts/
    └── generate_proto.sh           (NEW - script)
```

---

## ✅ Production Checklist

### Media Service
- [x] gRPC server implementation
- [x] All 8 RPC methods
- [x] Database integration
- [x] Redis caching
- [x] Storage backends (S3/MinIO/Local)
- [x] Thumbnail generation
- [x] Blurhash generation
- [x] File validation
- [x] Error handling
- [x] Ownership verification

### User Service
- [x] gRPC client
- [x] Profile picture upload
- [x] Cover photo upload
- [x] HTTP endpoints
- [x] JWT authentication
- [x] File validation
- [x] Error handling
- [x] Helper methods

---

## 🏆 Benefits Achieved

### For Developers
✅ Type-safe service communication  
✅ Auto-generated code  
✅ Clear API contracts  
✅ Easy testing  
✅ Versioning support  

### For Operations
✅ High performance (binary protocol)  
✅ Efficient (HTTP/2 multiplexing)  
✅ Observable (built-in metrics)  
✅ Reliable (connection management)  
✅ Scalable (stateless, load-balanced)  

### For Users
✅ Fast uploads  
✅ Auto thumbnails  
✅ Progressive loading (blurhash)  
✅ Responsive images  
✅ Reliable delivery  

---

## 🚀 Ready for Posting Service

**With gRPC integration complete, we're ready to build the posting service!**

The posting service will use the same media client to handle:
- Post images
- Post videos
- Multiple attachments per post
- Thumbnail generation
- CDN delivery

**All infrastructure is in place. Let's build the posting service next!** 🎉
