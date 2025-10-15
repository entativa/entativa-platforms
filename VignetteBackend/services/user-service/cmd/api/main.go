package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"vignette/user-service/internal/config"
	"vignette/user-service/internal/handler"
	"vignette/user-service/internal/middleware"
	"vignette/user-service/internal/repository"
	"vignette/user-service/internal/service"
	"vignette/user-service/pkg/cache"
	"vignette/user-service/pkg/database"
	grpcServer "vignette/user-service/internal/grpc"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Set Gin mode
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("✓ Connected to PostgreSQL database")

	// Wait for database to be ready
	if err := database.WaitForDB(db, 10, 2*time.Second); err != nil {
		log.Fatalf("Database not ready: %v", err)
	}

	// Initialize Redis cache
	var redisCache *cache.RedisCache
	redisCache, err = cache.NewRedisCache(&cfg.Redis)
	if err != nil {
		log.Printf("⚠️  Redis not available: %v", err)
		log.Println("   Running without caching and rate limiting")
	} else {
		log.Println("✓ Connected to Redis cache")
	}

	// Initialize Kafka producer
	var kafkaProducer *service.KafkaProducer
	if cfg.Kafka.Enabled {
		kafkaProducer = service.NewKafkaProducer(cfg.Kafka.Brokers)
		if kafkaProducer != nil {
			log.Println("✓ Kafka producer initialized")
			defer kafkaProducer.Close()
		}
	}

	// Initialize S3/MinIO media service
	var mediaService *service.MediaService
	if cfg.S3.AccessKeyID != "" {
		s3Config := &service.S3Config{
			Endpoint:        cfg.S3.Endpoint,
			AccessKeyID:     cfg.S3.AccessKeyID,
			SecretAccessKey: cfg.S3.SecretAccessKey,
			BucketName:      cfg.S3.BucketName,
			Region:          cfg.S3.Region,
			CDNURL:          cfg.S3.CDNURL,
			UsePathStyle:    cfg.S3.UsePathStyle,
		}
		mediaService, err = service.NewMediaService(s3Config)
		if err != nil {
			log.Printf("⚠️  Media service not available: %v", err)
		} else {
			log.Println("✓ Media upload service initialized (S3/MinIO)")
		}
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	twoFactorRepo := repository.NewTwoFactorRepository(db)
	passwordResetRepo := repository.NewPasswordResetRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, sessionRepo, cfg)
	twoFactorService := service.NewTwoFactorService(twoFactorRepo, userRepo)
	passwordResetService := service.NewPasswordResetService(userRepo, passwordResetRepo, kafkaProducer)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	twoFactorHandler := handler.NewTwoFactorHandler(twoFactorService)
	passwordResetHandler := handler.NewPasswordResetHandler(passwordResetService)
	mediaHandler := handler.NewMediaHandler(mediaService, userRepo)

	// Setup HTTP router
	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggerMiddleware())
	
	// Apply rate limiting if Redis is available
	if redisCache != nil {
		router.Use(middleware.APIRateLimitMiddleware(redisCache))
	}

	// Health check endpoint
	router.GET("/health", authHandler.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public auth routes
		auth := v1.Group("/auth")
		{
			// Apply specific rate limits for auth endpoints
			if redisCache != nil {
				auth.POST("/signup", middleware.SignupRateLimitMiddleware(redisCache), authHandler.Signup)
				auth.POST("/login", middleware.LoginRateLimitMiddleware(redisCache), authHandler.Login)
			} else {
				auth.POST("/signup", authHandler.Signup)
				auth.POST("/login", authHandler.Login)
			}

			// Password reset (public)
			auth.POST("/password-reset/request", passwordResetHandler.RequestPasswordReset)
			auth.POST("/password-reset/confirm", passwordResetHandler.ResetPassword)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			// User profile
			protected.GET("/auth/me", authHandler.Me)
			protected.POST("/auth/logout", authHandler.Logout)

			// 2FA management
			protected.POST("/auth/2fa/setup", twoFactorHandler.Setup2FA)
			protected.POST("/auth/2fa/enable", twoFactorHandler.Enable2FA)
			protected.POST("/auth/2fa/verify", twoFactorHandler.Verify2FA)
			protected.POST("/auth/2fa/disable", twoFactorHandler.Disable2FA)

			// Media upload
			if mediaService != nil {
				protected.POST("/media/profile-picture", mediaHandler.UploadProfilePicture)
			}
		}
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start gRPC server if enabled
	var grpcSrv *grpcServer.Server
	if cfg.GRPC.Enabled {
		grpcSrv = grpcServer.NewGRPCServer(
			cfg.GRPC.Port,
			userRepo,
			authService,
			twoFactorService,
		)
		go func() {
			if err := grpcSrv.Start(); err != nil {
				log.Printf("gRPC server error: %v", err)
			}
		}()
	}

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("🚀 Vignette User Service starting on port %s", cfg.Server.Port)
		log.Printf("📝 Environment: %s", cfg.Server.Environment)
		log.Println("🔐 Meta-level authentication enabled (instant access, no verification)")
		log.Println("✨ Enhanced Features:")
		log.Println("   • Two-Factor Authentication (TOTP)")
		log.Println("   • Password Reset Flows")
		log.Println("   • Account Recovery")
		if mediaService != nil {
			log.Println("   • Profile Picture Upload (S3/MinIO)")
		}
		if redisCache != nil {
			log.Println("   • Rate Limiting (Redis)")
			log.Println("   • Response Caching (Redis)")
		}
		if kafkaProducer != nil {
			log.Println("   • Event Publishing (Kafka)")
		}
		if cfg.GRPC.Enabled {
			log.Printf("   • gRPC Server (port %s)", cfg.GRPC.Port)
		}
		log.Println("✨ Ready to accept connections!")
		
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Stop gRPC server
	if grpcSrv != nil {
		grpcSrv.Stop()
	}

	// Close Redis
	if redisCache != nil {
		redisCache.Close()
	}

	log.Println("✓ Server shutdown complete")
}
