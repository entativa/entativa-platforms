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

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	
	"messaging-service/internal/config"
	"messaging-service/internal/encryption"
	"messaging-service/internal/handler"
	"messaging-service/internal/logger"
	"messaging-service/internal/middleware"
	"messaging-service/internal/repository"
	"messaging-service/internal/websocket"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	appLogger := logger.NewLogger()
	appLogger.Info("Starting Entativa Messaging Service...")

	// Initialize database
	db, err := sql.Open("postgres", cfg.DatabaseDSN())
	if err != nil {
		appLogger.Fatal("Failed to connect to database", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		appLogger.Fatal("Failed to ping database", err)
	}
	appLogger.Info("Connected to database successfully")

	// Initialize repositories
	messageRepo := repository.NewMessageRepository(db)
	conversationRepo := repository.NewConversationRepository(db)
	prekeyRepo := repository.NewPrekeyRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	// Initialize encryption service
	encryptionService := encryption.NewSignalProtocolService(prekeyRepo, sessionRepo)

	// Initialize WebSocket hub
	wsHub := ws.NewHub(messageRepo, conversationRepo, encryptionService, appLogger)
	go wsHub.Run()

	// Initialize handlers
	messageHandler := handler.NewMessageHandler(messageRepo, conversationRepo, encryptionService, wsHub, appLogger)
	conversationHandler := handler.NewConversationHandler(conversationRepo, messageRepo, appLogger)
	keyHandler := handler.NewKeyHandler(prekeyRepo, appLogger)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret, appLogger)
	corsMiddleware := middleware.NewCORSMiddleware(cfg.AllowedOrigins)
	loggingMiddleware := middleware.NewLoggingMiddleware(appLogger)

	// Setup router
	router := setupRouter(
		messageHandler,
		conversationHandler,
		keyHandler,
		authMiddleware,
		corsMiddleware,
		loggingMiddleware,
		wsHub,
	)

	// Create server
	srv := &http.Server{
		Addr:         cfg.ServerAddress,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	go func() {
		appLogger.Info(fmt.Sprintf("Server starting on %s", cfg.ServerAddress))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Failed to start server", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Fatal("Server forced to shutdown", err)
	}

	appLogger.Info("Server stopped")
}

func setupRouter(
	messageHandler *handler.MessageHandler,
	conversationHandler *handler.ConversationHandler,
	keyHandler *handler.KeyHandler,
	authMiddleware *middleware.AuthMiddleware,
	corsMiddleware *middleware.CORSMiddleware,
	loggingMiddleware *middleware.LoggingMiddleware,
	wsHub *ws.Hub,
) *mux.Router {
	router := mux.NewRouter()

	// Apply global middleware
	router.Use(corsMiddleware.Handler)
	router.Use(loggingMiddleware.Handler)

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	}).Methods("GET", "OPTIONS")

	// WebSocket endpoint
	router.HandleFunc("/api/v1/ws", authMiddleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wsHub.ServeWS(w, r)
	}))).Methods("GET")

	// Key management (for E2EE)
	router.HandleFunc("/api/v1/keys/prekeys", authMiddleware.RequireAuth(http.HandlerFunc(keyHandler.UploadPrekeys))).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/keys/prekeys/{userID}", authMiddleware.RequireAuth(http.HandlerFunc(keyHandler.GetPrekeyBundle))).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/keys/identity", authMiddleware.RequireAuth(http.HandlerFunc(keyHandler.UploadIdentityKey))).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/keys/identity/{userID}", authMiddleware.RequireAuth(http.HandlerFunc(keyHandler.GetIdentityKey))).Methods("GET", "OPTIONS")

	// Conversations
	router.HandleFunc("/api/v1/conversations", authMiddleware.RequireAuth(http.HandlerFunc(conversationHandler.GetConversations))).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/conversations", authMiddleware.RequireAuth(http.HandlerFunc(conversationHandler.CreateConversation))).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/conversations/{id}", authMiddleware.RequireAuth(http.HandlerFunc(conversationHandler.GetConversation))).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/conversations/{id}/mark-read", authMiddleware.RequireAuth(http.HandlerFunc(conversationHandler.MarkAsRead))).Methods("POST", "OPTIONS")

	// Messages
	router.HandleFunc("/api/v1/conversations/{conversationID}/messages", authMiddleware.RequireAuth(http.HandlerFunc(messageHandler.GetMessages))).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/conversations/{conversationID}/messages", authMiddleware.RequireAuth(http.HandlerFunc(messageHandler.SendMessage))).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/messages/{id}", authMiddleware.RequireAuth(http.HandlerFunc(messageHandler.DeleteMessage))).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/v1/messages/{id}/delivered", authMiddleware.RequireAuth(http.HandlerFunc(messageHandler.MarkDelivered))).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/messages/{id}/read", authMiddleware.RequireAuth(http.HandlerFunc(messageHandler.MarkRead))).Methods("POST", "OPTIONS")

	// Typing indicators
	router.HandleFunc("/api/v1/conversations/{conversationID}/typing", authMiddleware.RequireAuth(http.HandlerFunc(messageHandler.SendTypingIndicator))).Methods("POST", "OPTIONS")

	// Media upload (encrypted)
	router.HandleFunc("/api/v1/media/upload", authMiddleware.RequireAuth(http.HandlerFunc(messageHandler.UploadMedia))).Methods("POST", "OPTIONS")

	return router
}
