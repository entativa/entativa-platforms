package middleware

import (
	"net/http"
	"strings"

	"vignette/user-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthMiddleware creates a middleware that validates JWT tokens
func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Missing authorization header",
			})
			c.Abort()
			return
		}

		// Check if it's a Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid authorization header format. Use: Bearer <token>",
			})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		claims, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user data in context
		c.Set("user_id", claims.UserID.String())
		c.Set("email", claims.Email)
		c.Set("username", claims.Username)
		c.Set("full_name", claims.FullName)
		c.Set("user", map[string]interface{}{
			"id":        claims.UserID,
			"email":     claims.Email,
			"username":  claims.Username,
			"full_name": claims.FullName,
		})

		c.Next()
	}
}

// OptionalAuthMiddleware is a middleware that sets user context if token is present
func OptionalAuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		token := parts[1]
		claims, err := authService.ValidateToken(token)
		if err != nil {
			c.Next()
			return
		}

		// Set user data in context
		c.Set("user_id", claims.UserID.String())
		c.Set("email", claims.Email)
		c.Set("username", claims.Username)

		c.Next()
	}
}

// GetUserID extracts user ID from context
func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return uuid.Nil, false
	}

	return userID, true
}
