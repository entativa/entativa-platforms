package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/entativa/socialink/community-service/internal/handler"
	"github.com/entativa/socialink/community-service/internal/repository"
	"github.com/entativa/socialink/community-service/internal/service"
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
		dbURL = "postgresql://postgres:postgres@localhost:5432/socialink_community?sslmode=disable"
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
	communityRepo := repository.NewCommunityRepository(db)
	memberRepo := repository.NewMemberRepository(db)

	// Initialize services
	communityService := service.NewCommunityService(communityRepo, memberRepo)

	// Initialize handlers
	communityHandler := handler.NewCommunityHandler(communityService)

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

	// Auth middleware (simplified - in production use proper JWT validation)
	authMiddleware := func(c *gin.Context) {
		// Extract user_id from header or token
		userID := c.GetHeader("X-User-ID")
		if userID != "" {
			c.Set("user_id", userID)
		}
		c.Next()
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "socialink-community-service"})
	})

	// Root
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "socialink-community-service",
			"version": "1.0.0",
			"features": []string{
				"Community management",
				"Granular permissions",
				"Role-based access control",
				"Moderation tools",
				"Join requests & invites",
				"Ban & mute system",
				"Rules & guidelines",
				"Analytics & insights",
			},
		})
	})

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Public routes
		v1.GET("/communities", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "list communities"})
		})
		v1.GET("/communities/:id", communityHandler.GetCommunity)

		// Protected routes
		protected := v1.Group("/", authMiddleware)
		{
			// Community CRUD
			protected.POST("/communities", communityHandler.CreateCommunity)
			protected.PUT("/communities/:id", communityHandler.UpdateCommunity)

			// Membership
			protected.POST("/communities/:id/join", communityHandler.JoinCommunity)
			protected.POST("/communities/:id/leave", communityHandler.LeaveCommunity)

			// Moderation
			protected.POST("/communities/:id/members/:user_id/ban", communityHandler.BanMember)
			protected.PUT("/communities/:id/members/:user_id/role", communityHandler.UpdateMemberRole)
			protected.PUT("/communities/:id/members/:user_id/permissions", communityHandler.UpdateMemberPermissions)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8093"
	}

	log.Printf("🚀 Socialink Community Service starting on port %s", port)
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
