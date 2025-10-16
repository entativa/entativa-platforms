package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"vignette/creator-service/internal/handler"
	"vignette/creator-service/internal/repository"
	"vignette/creator-service/internal/service"
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
		dbURL = "postgresql://postgres:postgres@localhost:5432/vignette_creator?sslmode=disable"
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
	creatorRepo := repository.NewCreatorRepository(db)

	// Initialize services
	var mockKafka *service.KafkaProducer
	var mockRedis *service.RedisClient

	creatorService := service.NewCreatorService(creatorRepo, mockKafka, mockRedis)

	// Initialize handlers
	creatorHandler := handler.NewCreatorHandler(creatorService)

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
		c.JSON(200, gin.H{"status": "healthy", "service": "vignette-creator-service"})
	})

	// Root
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "vignette-creator-service",
			"version": "1.0.0",
			"features": []string{
				"Instagram-style creator tools",
				"Professional accounts (Personal, Business, Creator)",
				"Creator analytics & insights",
				"Audience demographics",
				"Content performance tracking",
				"Monetization program (10K followers, 100 posts)",
				"Creator badges (Verified, Partner, Top Creator, etc)",
				"Earnings tracking",
				"Top content discovery",
				"Follower growth analytics",
			},
		})
	})

	// API v1
	v1 := router.Group("/api/v1", authMiddleware)
	{
		// Profile management
		v1.POST("/profile", creatorHandler.CreateProfile)
		v1.GET("/profile", creatorHandler.GetProfile)
		v1.PUT("/profile", creatorHandler.UpdateProfile)

		// Analytics
		v1.GET("/analytics/overview", creatorHandler.GetAnalyticsOverview)
		v1.GET("/analytics/audience", creatorHandler.GetAudienceInsights)
		v1.GET("/analytics/content/top", creatorHandler.GetTopContent)

		// Monetization
		v1.POST("/monetization/apply", creatorHandler.ApplyForMonetization)

		// Admin endpoints (badge awarding)
		v1.POST("/admin/users/:user_id/badges/:badge", creatorHandler.AwardBadge)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8100"
	}

	log.Printf("🚀 Vignette Creator Service starting on port %s", port)
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
