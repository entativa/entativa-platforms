package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"socialink/post-service/internal/handler"
	"socialink/post-service/internal/repository"
	"socialink/post-service/internal/service"
	"socialink/post-service/pkg/database"
	"socialink/post-service/pkg/kafka"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"net/http"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize database
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnvAsInt("DB_PORT", 5432),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "socialink_posts"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	db, err := database.NewPostgresConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("✓ Database connection established")

	// Initialize Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       getEnvAsInt("REDIS_DB", 0),
		PoolSize: 100,
	})

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Redis connection failed: %v", err)
	} else {
		log.Println("✓ Redis connection established")
	}

	// Initialize Kafka producer
	kafkaBrokers := []string{getEnv("KAFKA_BROKERS", "localhost:9092")}
	kafkaProducer := kafka.NewProducer(kafkaBrokers)
	defer kafkaProducer.Close()

	log.Println("✓ Kafka producer initialized")

	// Initialize repositories
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	saveRepo := repository.NewSaveRepository(db)

	// Initialize services
	postService := service.NewPostService(postRepo, likeRepo, commentRepo, saveRepo, redisClient, kafkaProducer)
	commentService := service.NewCommentService(commentRepo, postRepo, redisClient, kafkaProducer)
	likeService := service.NewLikeService(likeRepo, postRepo, commentRepo, redisClient, kafkaProducer)

	// Initialize handlers
	postHandler := handler.NewPostHandler(postService)
	commentHandler := handler.NewCommentHandler(commentService)
	likeHandler := handler.NewLikeHandler(likeService)
	saveHandler := handler.NewSaveHandler(postService)

	// Setup Gin router
	if getEnv("GIN_MODE", "debug") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Middleware
	router.Use(corsMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "socialink-post-service",
			"time":    time.Now(),
		})
	})

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Post routes
		posts := v1.Group("/posts")
		{
			posts.POST("", authMiddleware(), postHandler.CreatePost)
			posts.GET("/feed", authMiddleware(), postHandler.GetFeed)
			posts.GET("/explore", postHandler.GetExplorePosts)
			posts.GET("/reels", postHandler.GetReels)
			posts.GET("/hashtag/:hashtag", postHandler.GetPostsByHashtag)
			posts.GET("/:post_id", postHandler.GetPost)
			posts.PUT("/:post_id", authMiddleware(), postHandler.UpdatePost)
			posts.DELETE("/:post_id", authMiddleware(), postHandler.DeletePost)
			posts.GET("/user/:user_id", postHandler.GetUserPosts)

			// Comment routes (nested)
			posts.POST("/:post_id/comments", authMiddleware(), commentHandler.CreateComment)
			posts.GET("/:post_id/comments", commentHandler.GetComments)

			// Like routes (nested)
			posts.POST("/:post_id/like", authMiddleware(), likeHandler.LikePost)
			posts.DELETE("/:post_id/like", authMiddleware(), likeHandler.UnlikePost)
			posts.GET("/:post_id/likes", likeHandler.GetPostLikers)

			// Save routes (nested)
			posts.POST("/:post_id/save", authMiddleware(), saveHandler.SavePost)
			posts.DELETE("/:post_id/save", authMiddleware(), saveHandler.UnsavePost)
		}

		// Comment routes (standalone)
		comments := v1.Group("/comments")
		{
			comments.GET("/:comment_id/replies", commentHandler.GetReplies)
			comments.PUT("/:comment_id", authMiddleware(), commentHandler.UpdateComment)
			comments.DELETE("/:comment_id", authMiddleware(), commentHandler.DeleteComment)
			comments.POST("/:comment_id/like", authMiddleware(), likeHandler.LikeComment)
			comments.DELETE("/:comment_id/like", authMiddleware(), likeHandler.UnlikeComment)
		}

		// Saved posts
		v1.GET("/saved", authMiddleware(), saveHandler.GetSavedPosts)
	}

	// Start server
	port := getEnv("PORT", "8084")
	addr := fmt.Sprintf(":%s", port)

	log.Printf("🚀 Socialink Post Service starting on %s", addr)

	// Graceful shutdown
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

// Middleware functions

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Authorization header required",
			})
			c.Abort()
			return
		}

		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "X-User-ID header required",
			})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	var value int
	if valueStr := os.Getenv(key); valueStr != "" {
		fmt.Sscanf(valueStr, "%d", &value)
		return value
	}
	return defaultValue
}
