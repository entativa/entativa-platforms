package middleware

import (
	"context"
	"net/http"
	"strings"

	"admin-service/internal/logger"
	"admin-service/internal/repository"
)

type FounderAuthMiddleware struct {
	userRepo *repository.UserRepository
	logger   *logger.Logger
}

func NewFounderAuthMiddleware(userRepo *repository.UserRepository, logger *logger.Logger) *FounderAuthMiddleware {
	return &FounderAuthMiddleware{
		userRepo: userRepo,
		logger:   logger,
	}
}

// RequireFounder ensures only the founder (@neoqiss) can access admin endpoints
func (m *FounderAuthMiddleware) RequireFounder(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract JWT token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.logger.Warn("Missing Authorization header in admin request")
			http.Error(w, `{"error":"Unauthorized: Missing authentication"}`, http.StatusUnauthorized)
			return
		}

		// Parse Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			m.logger.Warn("Invalid Authorization header format")
			http.Error(w, `{"error":"Unauthorized: Invalid token format"}`, http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// Verify JWT and get user
		user, err := m.userRepo.GetUserFromToken(r.Context(), token)
		if err != nil {
			m.logger.Error("Failed to verify token", err)
			http.Error(w, `{"error":"Unauthorized: Invalid token"}`, http.StatusUnauthorized)
			return
		}

		// CRITICAL: Verify founder account
		if user.Username != "neoqiss" {
			m.logger.Warn("Non-founder attempted admin access", map[string]interface{}{
				"username": user.Username,
				"user_id":  user.ID,
				"ip":       getIPAddress(r),
			})
			http.Error(w, `{"error":"Forbidden: Founder access required"}`, http.StatusForbidden)
			m.logSecurityEvent(r, user, "unauthorized_admin_access_attempt")
			return
		}

		if !user.IsFounder {
			m.logger.Warn("User @neoqiss missing is_founder flag", map[string]interface{}{
				"user_id": user.ID,
			})
			http.Error(w, `{"error":"Forbidden: Founder flag not set"}`, http.StatusForbidden)
			return
		}

		// Verify device whitelist (if device_id is provided)
		deviceID := r.Header.Get("X-Device-ID")
		if deviceID != "" {
			isWhitelisted, err := m.userRepo.IsDeviceWhitelisted(r.Context(), user.ID, deviceID)
			if err != nil {
				m.logger.Error("Failed to check device whitelist", err)
			} else if !isWhitelisted {
				m.logger.Warn("Unrecognized device attempting admin access", map[string]interface{}{
					"user_id":   user.ID,
					"device_id": deviceID,
					"ip":        getIPAddress(r),
				})
				http.Error(w, `{"error":"Forbidden: Device not recognized"}`, http.StatusForbidden)
				m.logSecurityEvent(r, user, "unrecognized_device_admin_access")
				return
			}
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), "user", user)
		ctx = context.WithValue(ctx, "is_admin", true)
		ctx = context.WithValue(ctx, "admin_level", 10) // Supreme Founder level

		m.logger.Info("Founder admin access granted", map[string]interface{}{
			"username": user.Username,
			"endpoint": r.URL.Path,
			"method":   r.Method,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *FounderAuthMiddleware) logSecurityEvent(r *http.Request, user *repository.User, eventType string) {
	// Log security event for monitoring
	m.logger.Error("SECURITY EVENT: "+eventType, nil, map[string]interface{}{
		"user_id":   user.ID,
		"username":  user.Username,
		"ip":        getIPAddress(r),
		"user_agent": r.UserAgent(),
		"endpoint":  r.URL.Path,
		"method":    r.Method,
	})
}

func getIPAddress(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxied requests)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Take the first IP in the list
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	// Remove port if present
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}
