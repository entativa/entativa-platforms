# Quick Integration Guide - Profile Management Services

## 🚀 5-Minute Integration

### Step 1: Run Database Migrations

**Socialink**:
```bash
psql -U postgres -d socialink_db -f SocialinkBackend/services/user-service/migrations/003_create_profiles_table.up.sql
```

**Vignette**:
```bash
psql -U postgres -d vignette_db -f VignetteBackend/services/user-service/migrations/004_create_profiles_table.up.sql
```

### Step 2: Add Routes to Main.go

**For Both Services** - Add to `cmd/api/main.go`:

```go
// After initializing userRepo and redisCache:

// Initialize profile repository
profileRepo := repository.NewProfileRepository(db, redisCache)

// Initialize profile service  
profileService := service.NewProfileService(profileRepo, userRepo)

// Initialize profile handler
profileHandler := handler.NewProfileHandler(profileService)

// Add to your router (after auth routes):
profile := r.Group("/profile")
profile.Use(authMiddleware.RequireAuth()) // Your existing auth middleware
{
    // Get profiles
    profile.GET("/me", profileHandler.GetMyProfile)
    profile.GET("/:user_id", profileHandler.GetProfile)
    
    // Socialink-specific
    profile.PUT("/info", profileHandler.UpdateProfileInfo)          // Basic info
    profile.POST("/work", profileHandler.AddWorkExperience)         // Add work
    profile.DELETE("/work/:work_id", profileHandler.RemoveWorkExperience)
    profile.POST("/education", profileHandler.AddEducation)         // Add education
    profile.DELETE("/education/:education_id", profileHandler.RemoveEducation)
    profile.PUT("/contact", profileHandler.UpdateContactInfo)       // Contact info
    profile.PUT("/social-links", profileHandler.UpdateSocialLinks)  // Social media
    profile.PUT("/visibility", profileHandler.UpdateVisibility)     // Privacy
    
    // Vignette-specific (replace above with these)
    profile.PUT("/extended", profileHandler.UpdateProfileExtended)  // Category, gender, etc.
    profile.POST("/links", profileHandler.AddLinkInBio)            // Link in bio
    profile.DELETE("/links/:link_id", profileHandler.RemoveLinkInBio)
    profile.POST("/highlights", profileHandler.AddHighlight)        // Story highlights
    profile.DELETE("/highlights/:highlight_id", profileHandler.RemoveHighlight)
    profile.PUT("/contact-options", profileHandler.UpdateContactOptions)
    profile.POST("/creator/enable", profileHandler.EnableCreatorAccount)
    profile.POST("/business/enable", profileHandler.EnableBusinessAccount)
    profile.PUT("/availability", profileHandler.UpdateAvailability)
}
```

### Step 3: Test It!

```bash
# Get your profile (works for both platforms)
curl http://localhost:8080/profile/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Should auto-create profile with defaults on first access
```

## 📋 Complete Route List

### Socialink Routes
```
GET    /profile/me
GET    /profile/:user_id
PUT    /profile/info
POST   /profile/work
DELETE /profile/work/:work_id
POST   /profile/education  
DELETE /profile/education/:education_id
PUT    /profile/contact
PUT    /profile/social-links
PUT    /profile/visibility
```

### Vignette Routes
```
GET    /profile/me
GET    /profile/:user_id
PUT    /profile/extended
POST   /profile/links
DELETE /profile/links/:link_id
POST   /profile/highlights
DELETE /profile/highlights/:highlight_id
PUT    /profile/contact-options
POST   /profile/creator/enable
POST   /profile/business/enable
PUT    /profile/availability
```

## 🔧 Troubleshooting

### Profile not found?
Profiles are auto-created on first access. Just call `/profile/me` and it creates one with defaults.

### Cache not working?
Ensure Redis is running: `redis-cli ping` should return `PONG`

### Migration errors?
Make sure the users table exists first (from your existing migrations)

## 📚 Full Documentation

See `PROFILE_MANAGEMENT_IMPLEMENTATION.md` for complete details, examples, and features.

---

**That's it! Profile management is now live on both platforms!** 🎉
