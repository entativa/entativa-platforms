package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"user-service/internal/config"
	"user-service/internal/handler"
	"user-service/internal/logger"
	"user-service/internal/middleware"
	"user-service/internal/repository"
	"user-service/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	
	// Initialize logger
	appLogger := logger.NewLogger()
	appLogger.Info("Starting %s on port %s", cfg.Server.ServiceName, cfg.Server.Port)
	
	// Connect to database
	db, err := connectDatabase(cfg)
	if err != nil {
		appLogger.Fatal("Failed to connect to database", err)
	}
	defer db.Close()
	
	appLogger.Info("Connected to database: %s", cfg.Database.Name)
	
	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	
	// Initialize services
	emailService := service.NewEmailService()
	auditLog := service.NewAuditLog(db)
	
	// Initialize handlers
	authHandler := handler.NewAuthHandler(
		userRepo,
		sessionRepo,
		tokenRepo,
		emailService,
		auditLog,
		appLogger,
		cfg,
	)
	
	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(userRepo, appLogger)
	
	// Setup routes
	router := SetupRoutes(authHandler, authMiddleware)
	
	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	// Start cleanup goroutine for expired sessions and tokens
	go cleanupExpiredData(sessionRepo, tokenRepo, appLogger)
	
	// Start server in a goroutine
	go func() {
		appLogger.Info("Server listening on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Server error", err)
		}
	}()
	
	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	appLogger.Info("Shutting down server...")
	
	// Graceful shutdown with 30 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := server.Shutdown(ctx); err != nil {
		appLogger.Error("Server forced to shutdown", err)
	}
	
	appLogger.Info("Server stopped")
}

// connectDatabase establishes database connection
func connectDatabase(cfg *config.Config) (*sql.DB, error) {
	dsn := cfg.GetDatabaseDSN()
	
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	
	// Configure connection pool
	db.SetMaxOpenConns(cfg.Database.MaxConnections)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConnections)
	db.SetConnMaxLifetime(time.Hour)
	
	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	return db, nil
}

// cleanupExpiredData periodically cleans up expired sessions and tokens
func cleanupExpiredData(
	sessionRepo *repository.SessionRepository,
	tokenRepo *repository.TokenRepository,
	logger *logger.Logger,
) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	
	for range ticker.C {
		ctx := context.Background()
		
		// Clean up expired sessions
		if err := sessionRepo.DeleteExpiredSessions(ctx); err != nil {
			logger.Error("Failed to delete expired sessions", err)
		}
		
		// Clean up expired password reset tokens
		if err := tokenRepo.DeleteExpiredTokens(ctx); err != nil {
			logger.Error("Failed to delete expired tokens", err)
		}
	}
}
