package middleware

import (
	"fmt"
	"net/http"
	"time"

	"vignette/user-service/pkg/cache"

	"github.com/gin-gonic/gin"
)

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware(redis *cache.RedisCache, maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if redis == nil {
			c.Next()
			return
		}

		// Get client identifier (IP address or user ID if authenticated)
		identifier := c.ClientIP()
		if userID, exists := c.Get("user_id"); exists {
			identifier = fmt.Sprintf("user:%s", userID)
		}

		// Create rate limit key
		key := fmt.Sprintf("rate_limit:%s:%s", c.Request.URL.Path, identifier)

		// Increment counter
		count, err := redis.IncrementWithExpiry(key, window)
		if err != nil {
			// If Redis fails, allow request (fail open)
			c.Next()
			return
		}

		// Check if rate limit exceeded
		if count > int64(maxRequests) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": fmt.Sprintf("Too many requests. Please try again in %v", window),
			})
			c.Abort()
			return
		}

		// Add rate limit headers
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", maxRequests))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", maxRequests-int(count)))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(window).Unix()))

		c.Next()
	}
}

// LoginRateLimitMiddleware specific rate limiter for login endpoints
func LoginRateLimitMiddleware(redis *cache.RedisCache) gin.HandlerFunc {
	return RateLimitMiddleware(redis, 5, 15*time.Minute) // 5 attempts per 15 minutes
}

// SignupRateLimitMiddleware specific rate limiter for signup endpoints
func SignupRateLimitMiddleware(redis *cache.RedisCache) gin.HandlerFunc {
	return RateLimitMiddleware(redis, 3, 1*time.Hour) // 3 signups per hour per IP
}

// APIRateLimitMiddleware general rate limiter for API endpoints
func APIRateLimitMiddleware(redis *cache.RedisCache) gin.HandlerFunc {
	return RateLimitMiddleware(redis, 100, 1*time.Minute) // 100 requests per minute
}
