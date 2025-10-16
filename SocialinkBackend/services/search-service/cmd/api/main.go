package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"socialink/search-service/internal/elasticsearch"
	"socialink/search-service/internal/handler"
	"socialink/search-service/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize Elasticsearch
	esAddresses := []string{getEnv("ELASTICSEARCH_URL", "http://localhost:9200")}
	esUsername := getEnv("ELASTICSEARCH_USERNAME", "")
	esPassword := getEnv("ELASTICSEARCH_PASSWORD", "")

	esClient, err := elasticsearch.NewClient(esAddresses, esUsername, esPassword)
	if err != nil {
		log.Fatalf("Failed to connect to Elasticsearch: %v", err)
	}
	log.Println("✅ Connected to Elasticsearch")

	// Initialize Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_URL", "localhost:6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("✅ Connected to Redis")

	// Initialize services
	searchService := service.NewSearchService(esClient, redisClient)
	autocompleteService := service.NewAutocompleteService(esClient, redisClient)
	indexingService := service.NewIndexingService(esClient, redisClient)
	hashtagService := service.NewHashtagService(esClient, redisClient)

	// Initialize handlers
	searchHandler := handler.NewSearchHandler(searchService)
	autocompleteHandler := handler.NewAutocompleteHandler(autocompleteService)
	indexingHandler := handler.NewIndexingHandler(indexingService)
	hashtagHandler := handler.NewHashtagHandler(hashtagService)

	// Setup Gin router
	if getEnv("GIN_MODE", "debug") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-User-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "Socialink Search Service",
			"version": "1.0.0",
		})
	})

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service":     "Socialink Search Service",
			"version":     "1.0.0",
			"description": "Multi-entity search with Elasticsearch",
			"features": []string{
				"Multi-entity search",
				"Real-time autocomplete",
				"Trending hashtags",
				"Related hashtags",
				"Advanced filters",
				"Search history",
				"Trending searches",
			},
			"docs": "/swagger",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Search routes
		search := v1.Group("/search")
		{
			search.GET("", searchHandler.Search)
			search.GET("/users", searchHandler.SearchUsers)
			search.GET("/posts", searchHandler.SearchPosts)
			search.GET("/takes", searchHandler.SearchTakes)
			search.GET("/history", searchHandler.GetSearchHistory)
			search.DELETE("/history", searchHandler.DeleteSearchHistory)
			search.GET("/trending", searchHandler.GetTrendingSearches)
		}

		// Autocomplete routes
		autocomplete := v1.Group("/autocomplete")
		{
			autocomplete.GET("", autocompleteHandler.Autocomplete)
			autocomplete.GET("/recent", autocompleteHandler.GetRecentSearches)
		}

		// Hashtag routes
		hashtags := v1.Group("/hashtags")
		{
			hashtags.GET("/trending", hashtagHandler.GetTrendingHashtags)
			hashtags.GET("/:tag/related", hashtagHandler.GetRelatedHashtags)
			hashtags.GET("/search", hashtagHandler.SearchHashtags)
		}

		// Indexing routes (should be protected in production)
		index := v1.Group("/index")
		{
			index.POST("/document", indexingHandler.IndexDocument)
			index.POST("/bulk", indexingHandler.BulkIndex)
			index.PUT("/document", indexingHandler.UpdateDocument)
			index.DELETE("/document", indexingHandler.DeleteDocument)
			index.POST("/reindex", indexingHandler.ReindexAll)
			index.GET("/stats", indexingHandler.GetIndexStats)
		}
	}

	// Server configuration
	port := getEnv("PORT", "8088")
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 Search service starting on port %s", port)
		log.Printf("📊 Elasticsearch: %s", esAddresses[0])
		log.Printf("💾 Redis: %s", getEnv("REDIS_URL", "localhost:6379"))
		log.Printf("🔍 Search endpoints: http://localhost:%s/api/v1/search", port)
		log.Printf("⚡ Autocomplete: http://localhost:%s/api/v1/autocomplete", port)
		log.Printf("🏷️  Hashtags: http://localhost:%s/api/v1/hashtags", port)
		log.Printf("📝 Indexing: http://localhost:%s/api/v1/index", port)
		log.Println("✅ Search service ready!")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down search service...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Close connections
	if err := redisClient.Close(); err != nil {
		log.Printf("Error closing Redis connection: %v", err)
	}

	log.Println("✅ Search service shut down successfully")
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
