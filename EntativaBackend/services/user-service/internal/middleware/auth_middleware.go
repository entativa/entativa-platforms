package middleware

import (
	"context"
	"net/http"

	"user-service/internal/logger"
	"user-service/internal/repository"
	"user-service/internal/util"
)

// AuthMiddleware handles JWT authentication
type AuthMiddleware struct {
	userRepo *repository.UserRepository
	logger   *logger.Logger
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(userRepo *repository.UserRepository, logger *logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		userRepo: userRepo,
		logger:   logger,
	}
}

// RequireAuth middleware ensures the request has a valid JWT token
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			util.RespondWithUnauthorized(w, "Authorization header required")
			return
		}
		
		// Extract token
		token, err := util.ExtractTokenFromHeader(authHeader)
		if err != nil {
			util.RespondWithUnauthorized(w, err.Error())
			return
		}
		
		// Parse and validate token
		claims, err := util.ParseAccessToken(token)
		if err != nil {
			util.RespondWithUnauthorized(w, "Invalid or expired token")
			return
		}
		
		// Get user from database
		user, err := m.userRepo.FindByID(r.Context(), claims.UserID)
		if err != nil {
			util.RespondWithUnauthorized(w, "User not found")
			return
		}
		
		// Check if user is active
		if !user.IsActive {
			util.RespondWithForbidden(w, "Account is deactivated")
			return
		}
		
		// Add user to request context
		ctx := context.WithValue(r.Context(), "user", user)
		ctx = context.WithValue(ctx, "user_id", user.ID)
		
		// Continue to next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth middleware attempts to authenticate but doesn't require it
func (m *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		
		if authHeader != "" {
			token, err := util.ExtractTokenFromHeader(authHeader)
			if err == nil {
				claims, err := util.ParseAccessToken(token)
				if err == nil {
					user, err := m.userRepo.FindByID(r.Context(), claims.UserID)
					if err == nil && user.IsActive {
						ctx := context.WithValue(r.Context(), "user", user)
						ctx = context.WithValue(ctx, "user_id", user.ID)
						r = r.WithContext(ctx)
					}
				}
			}
		}
		
		next.ServeHTTP(w, r)
	})
}

// RequireFounder middleware ensures the user is the founder (@neoqiss)
func (m *AuthMiddleware) RequireFounder(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value("user").(*repository.User)
		if !ok {
			util.RespondWithUnauthorized(w, "")
			return
		}
		
		// Check if user is the founder
		if user.Username != "neoqiss" {
			util.RespondWithForbidden(w, "Founder access required")
			m.logger.Warn("Unauthorized admin access attempt", nil)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}
