package main

import (
	"fmt"
	"log"
	"os"

	grpcclient "vignette/api-gateway/internal/grpc"
	"vignette/api-gateway/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment
	godotenv.Load()

	// Configure service URLs
	serviceConfig := &grpcclient.ServiceConfig{
		UserServiceURL:          getEnv("USER_SERVICE_GRPC", "localhost:50001"),
		PostServiceURL:          getEnv("POST_SERVICE_GRPC", "localhost:50002"),
		MessagingServiceURL:     getEnv("MESSAGING_SERVICE_GRPC", "localhost:50003"),
		SettingsServiceURL:      getEnv("SETTINGS_SERVICE_GRPC", "localhost:50004"),
		MediaServiceURL:         getEnv("MEDIA_SERVICE_GRPC", "localhost:50051"),
		StoryServiceURL:         getEnv("STORY_SERVICE_GRPC", "localhost:50005"),
		SearchServiceURL:        getEnv("SEARCH_SERVICE_GRPC", "localhost:50006"),
		NotificationServiceURL:  getEnv("NOTIFICATION_SERVICE_GRPC", "localhost:50007"),
		FeedServiceURL:          getEnv("FEED_SERVICE_GRPC", "localhost:50008"),
		CommunityServiceURL:     getEnv("COMMUNITY_SERVICE_GRPC", "localhost:50009"),
		RecommendationServiceURL: getEnv("RECOMMENDATION_SERVICE_GRPC", "localhost:50010"),
		StreamingServiceURL:     getEnv("STREAMING_SERVICE_GRPC", "localhost:50011"),
		CreatorServiceURL:       getEnv("CREATOR_SERVICE_GRPC", "localhost:50012"),
	}

	// Initialize gRPC clients
	grpcClients, err := grpcclient.NewGRPCClients(serviceConfig)
	if err != nil {
		log.Printf("Warning: Some gRPC services not available: %v", err)
	}
	defer grpcClients.Close()

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

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "vignette-api-gateway",
			"grpc_clients": map[string]bool{
				"user":           grpcClients.UserConn != nil,
				"post":           grpcClients.PostConn != nil,
				"messaging":      grpcClients.MessagingConn != nil,
				"settings":       grpcClients.SettingsConn != nil,
				"media":          grpcClients.MediaConn != nil,
				"story":          grpcClients.StoryConn != nil,
				"search":         grpcClients.SearchConn != nil,
				"notification":   grpcClients.NotificationConn != nil,
				"feed":           grpcClients.FeedConn != nil,
				"community":      grpcClients.CommunityConn != nil,
				"recommendation": grpcClients.RecommendationConn != nil,
				"streaming":      grpcClients.StreamingConn != nil,
				"creator":        grpcClients.CreatorConn != nil,
			},
		})
	})

	// Root
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "vignette-api-gateway",
			"version": "1.0.0",
			"description": "API Gateway for all Vignette microservices",
			"features": []string{
				"gRPC-based inter-service communication",
				"Unified REST API for native clients",
				"JWT authentication",
				"Request routing to microservices",
				"Service health monitoring",
			},
			"services": []string{
				"User Service",
				"Post Service",
				"Messaging Service (E2EE)",
				"Settings Service",
				"Media Service",
				"Story Service",
				"Search Service",
				"Notification Service",
				"Feed Service",
				"Community Service",
				"Recommendation Service",
				"Live Streaming Service",
				"Creator Service",
			},
		})
	})

	// API v1 routes (authenticated)
	v1 := router.Group("/api/v1")
	v1.Use(middleware.AuthMiddleware())
	{
		// User routes
		v1.GET("/users/:id", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "User endpoint - gRPC call to User Service"})
		})
		v1.PUT("/users/:id", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Update user - gRPC call to User Service"})
		})
		
		// Post routes
		v1.POST("/posts", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Create post - gRPC call to Post Service"})
		})
		v1.GET("/posts/:id", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Get post - gRPC call to Post Service"})
		})
		
		// Messaging routes
		v1.POST("/messages", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Send message - gRPC call to Messaging Service"})
		})
		v1.GET("/messages", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Get messages - gRPC call to Messaging Service"})
		})
		
		// Settings routes
		v1.GET("/settings", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Get settings - gRPC call to Settings Service"})
		})
		v1.PUT("/settings", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Update settings - gRPC call to Settings Service"})
		})
		
		// Story routes
		v1.POST("/stories", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Create story - gRPC call to Story Service"})
		})
		v1.GET("/stories", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Get stories - gRPC call to Story Service"})
		})
		
		// Community routes
		v1.POST("/communities", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Create community - gRPC call to Community Service"})
		})
		v1.GET("/communities/:id", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Get community - gRPC call to Community Service"})
		})
		
		// Streaming routes
		v1.POST("/streams", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Create stream - gRPC call to Streaming Service"})
		})
		v1.GET("/streams/live", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Get live streams - gRPC call to Streaming Service"})
		})
		
		// Creator routes
		v1.GET("/creator/profile", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Get creator profile - gRPC call to Creator Service"})
		})
		v1.GET("/creator/analytics", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Get analytics - gRPC call to Creator Service"})
		})
	}

	// Public routes (no auth)
	public := router.Group("/api/v1/public")
	{
		public.GET("/feed", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Public feed - gRPC call to Feed Service"})
		})
		public.GET("/search", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Search - gRPC call to Search Service"})
		})
	}

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("🚀 Vignette API Gateway starting on port %s", port)
	log.Printf("📡 Routing requests to microservices via gRPC")
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
