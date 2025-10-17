package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware logs HTTP requests
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(startTime)

		// Get status code
		statusCode := c.Writer.Status()

		// Log format
		log.Printf(
			"[%s] %s %s | Status: %d | Latency: %v | IP: %s | User-Agent: %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Request.Proto,
			statusCode,
			latency,
			c.ClientIP(),
			c.Request.UserAgent(),
		)
	}
}
