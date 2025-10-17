package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"socialink/event-service/internal/handler"
	"socialink/event-service/internal/repository"
	"socialink/event-service/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment
	godotenv.Load()

	// Database connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgresql://postgres:postgres@localhost:5432/socialink_events?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("✅ Connected to database")

	// Initialize repositories
	eventRepo := repository.NewEventRepository(db)
	rsvpRepo := repository.NewRSVPRepository(db)

	// Initialize services
	var mockKafka *service.KafkaProducer
	var mockRedis *service.RedisClient

	eventService := service.NewEventService(eventRepo, rsvpRepo, mockKafka, mockRedis)

	// Initialize handlers
	eventHandler := handler.NewEventHandler(eventService)

	// Setup router
	router := gin.Default()

	// CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Auth middleware (simplified)
	authMiddleware := func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		if userID != "" {
			c.Set("user_id", userID)
		}
		c.Next()
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "socialink-event-service"})
	})

	// Root
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "socialink-event-service",
			"version": "1.0.0",
			"features": []string{
				"Facebook-style events",
				"In-person & virtual events",
				"RSVP (Going, Interested, Not Going)",
				"Event invitations",
				"Event check-in",
				"Recurring events",
				"Co-hosts",
				"Event discussions",
				"Location-based discovery",
				"Full-text search",
				"Privacy controls",
			},
		})
	})

	// API v1
	v1 := router.Group("/api/v1", authMiddleware)
	{
		// Event management
		v1.POST("/events", eventHandler.CreateEvent)
		v1.GET("/events/:id", eventHandler.GetEvent)
		v1.PUT("/events/:id", eventHandler.UpdateEvent)
		v1.DELETE("/events/:id", eventHandler.CancelEvent)

		// Discovery
		v1.GET("/events", eventHandler.GetUpcomingEvents)
		v1.GET("/events/search", eventHandler.SearchEvents)
		v1.GET("/events/nearby", eventHandler.GetNearbyEvents)

		// RSVP
		v1.POST("/events/:id/rsvp", eventHandler.RSVP)
		v1.DELETE("/events/:id/rsvp", eventHandler.RemoveRSVP)
		v1.POST("/events/:id/checkin", eventHandler.CheckIn)

		// Attendees
		v1.GET("/events/:id/attendees", eventHandler.GetEventAttendees)
		v1.GET("/events/:id/stats", eventHandler.GetEventStats)

		// User events
		v1.GET("/users/events", eventHandler.GetUserEvents)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8099"
	}

	log.Printf("🚀 Socialink Event Service starting on port %s", port)
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
