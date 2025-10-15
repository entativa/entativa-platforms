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

	"socialink/meta-user-service/internal/config"
	"socialink/meta-user-service/internal/handler"
	"socialink/meta-user-service/internal/repository"
	"socialink/meta-user-service/internal/service"
	"socialink/meta-user-service/pkg/cache"
	"socialink/meta-user-service/pkg/database"
	"socialink/meta-user-service/pkg/kafka"
	"socialink/meta-user-service/pkg/ml"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.NewPostgresConnection(database.PostgresConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
		MaxConns: cfg.Database.MaxConns,
		MinConns: cfg.Database.MinConns,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Successfully connected to PostgreSQL")

	// Initialize Redis cache
	redisCache, err := cache.NewRedisCache(
		cfg.GetRedisAddr(),
		cfg.Redis.Password,
		cfg.Redis.DB,
		cfg.Redis.PoolSize,
	)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Successfully connected to Redis")

	// Initialize Kafka event publisher
	eventPublisher := kafka.NewEventPublisher(cfg.Kafka.Brokers, cfg.Kafka.Topic)

	log.Println("Successfully initialized Kafka producer")

	// Initialize ML components
	fraudDetector := ml.NewFraudDetector()
	behaviorAnalyzer := ml.NewBehaviorAnalyzer()
	trustScoreEngine := ml.NewTrustScoreEngine()

	log.Println("Successfully initialized ML components")

	// Initialize repository
	metaUserRepo := repository.NewMetaUserRepository(db, redisCache)

	// Initialize cross-platform clients (would use actual gRPC clients in production)
	vignetteClient := &service.VignetteClient{}
	socialinkClient := &service.SocialinkClient{}

	// Initialize services
	crossPlatformSync := service.NewCrossPlatformSyncService(
		metaUserRepo,
		eventPublisher,
		vignetteClient,
		socialinkClient,
	)

	metaUserService := service.NewMetaUserService(
		metaUserRepo,
		fraudDetector,
		behaviorAnalyzer,
		trustScoreEngine,
		eventPublisher,
		crossPlatformSync,
	)

	log.Println("Successfully initialized services")

	// Initialize handlers
	metaUserHandler := handler.NewMetaUserHandler(metaUserService, crossPlatformSync)

	// Setup router
	router := setupRouter(metaUserHandler, cfg.Server.Environment)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting Meta User Service on %s:%d", cfg.Server.Host, cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}

func setupRouter(metaUserHandler *handler.MetaUserHandler, environment string) *gin.Engine {
	if environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// User management
		users := v1.Group("/users")
		{
			users.POST("", metaUserHandler.CreateUser)
			users.POST("/authenticate", metaUserHandler.Authenticate)
			users.GET("/:id", metaUserHandler.GetUser)
			users.POST("/:id/trust-score", metaUserHandler.UpdateTrustScore)
		}

		// Platform linking
		platforms := v1.Group("/platforms")
		{
			platforms.POST("/link", metaUserHandler.LinkPlatform)
		}

		// Cross-platform sync
		sync := v1.Group("/sync")
		{
			sync.POST("/:id/enable", metaUserHandler.EnableSync)
			sync.POST("/:id/disable", metaUserHandler.DisableSync)
		}

		// Privacy settings
		privacy := v1.Group("/privacy")
		{
			privacy.PUT("/:id/settings", metaUserHandler.UpdatePrivacySettings)
		}
	}

	return router
}
