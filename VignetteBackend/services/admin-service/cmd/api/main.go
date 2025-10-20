package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"admin-service/internal/config"
	"admin-service/internal/handler"
	"admin-service/internal/logger"
	"admin-service/internal/middleware"
	"admin-service/internal/repository"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	appLogger := logger.NewLogger()
	appLogger.Info("Starting Entativa Admin Service...")

	// Initialize database
	db, err := sql.Open("postgres", cfg.DatabaseDSN())
	if err != nil {
		appLogger.Fatal("Failed to connect to database", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		appLogger.Fatal("Failed to ping database", err)
	}
	appLogger.Info("Connected to database successfully")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	adminRepo := repository.NewAdminRepository(db)
	contentRepo := repository.NewContentRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	// Initialize handlers
	userMgmtHandler := handler.NewUserManagementHandler(userRepo, adminRepo, auditRepo, appLogger)
	contentModerationHandler := handler.NewContentModerationHandler(contentRepo, auditRepo, appLogger)
	platformControlHandler := handler.NewPlatformControlHandler(adminRepo, auditRepo, appLogger)
	analyticsHandler := handler.NewAnalyticsHandler(adminRepo, appLogger)
	auditHandler := handler.NewAuditHandler(auditRepo, appLogger)
	securityHandler := handler.NewSecurityHandler(sessionRepo, adminRepo, auditRepo, appLogger)

	// Initialize middleware
	founderAuthMiddleware := middleware.NewFounderAuthMiddleware(userRepo, appLogger)
	auditMiddleware := middleware.NewAuditMiddleware(auditRepo, appLogger)
	corsMiddleware := middleware.NewCORSMiddleware(cfg.AllowedOrigins)
	loggingMiddleware := middleware.NewLoggingMiddleware(appLogger)

	// Setup routes
	router := setupRoutes(
		userMgmtHandler,
		contentModerationHandler,
		platformControlHandler,
		analyticsHandler,
		auditHandler,
		securityHandler,
		founderAuthMiddleware,
		auditMiddleware,
		corsMiddleware,
		loggingMiddleware,
	)

	// Create server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		appLogger.Info("Admin Service listening on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Failed to start server", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down admin service...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		appLogger.Fatal("Admin service forced to shutdown", err)
	}

	appLogger.Info("Admin service stopped")
}

func setupRoutes(
	userMgmt *handler.UserManagementHandler,
	contentMod *handler.ContentModerationHandler,
	platform *handler.PlatformControlHandler,
	analytics *handler.AnalyticsHandler,
	audit *handler.AuditHandler,
	security *handler.SecurityHandler,
	founderAuth *middleware.FounderAuthMiddleware,
	auditMiddleware *middleware.AuditMiddleware,
	corsMiddleware *middleware.CORSMiddleware,
	loggingMiddleware *middleware.LoggingMiddleware,
) *mux.Router {
	r := mux.NewRouter()

	// Global middleware
	r.Use(corsMiddleware.Handler)
	r.Use(loggingMiddleware.Handler)

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"admin"}`))
	}).Methods("GET")

	// Admin API (all require founder authentication)
	api := r.PathPrefix("/api/admin").Subrouter()
	api.Use(founderAuth.RequireFounder)
	api.Use(auditMiddleware.LogAction)

	// User Management
	api.HandleFunc("/users/{id}", userMgmt.GetUser).Methods("GET")
	api.HandleFunc("/users/{id}/ban", userMgmt.BanUser).Methods("POST")
	api.HandleFunc("/users/{id}/unban", userMgmt.UnbanUser).Methods("POST")
	api.HandleFunc("/users/{id}/shadowban", userMgmt.ShadowbanUser).Methods("POST")
	api.HandleFunc("/users/{id}/unshadowban", userMgmt.UnshadowbanUser).Methods("POST")
	api.HandleFunc("/users/{id}/mute", userMgmt.MuteUser).Methods("POST")
	api.HandleFunc("/users/{id}/unmute", userMgmt.UnmuteUser).Methods("POST")
	api.HandleFunc("/users/{id}/suspend", userMgmt.SuspendUser).Methods("POST")
	api.HandleFunc("/users/{id}/permanently-delete", userMgmt.PermanentlyDeleteUser).Methods("DELETE")
	api.HandleFunc("/users/{id}/force-logout", userMgmt.ForceLogout).Methods("POST")
	api.HandleFunc("/users/{id}/reset-password", userMgmt.ForcePasswordReset).Methods("POST")
	api.HandleFunc("/users/{id}/disable-2fa", userMgmt.Disable2FA).Methods("POST")
	api.HandleFunc("/users/{id}/impersonate", userMgmt.ImpersonateUser).Methods("POST")
	api.HandleFunc("/users/{id}/end-impersonation", userMgmt.EndImpersonation).Methods("POST")
	api.HandleFunc("/users/search", userMgmt.SearchUsers).Methods("GET")

	// Content Moderation
	api.HandleFunc("/content/{id}", contentMod.GetContent).Methods("GET")
	api.HandleFunc("/content/{id}", contentMod.DeleteContent).Methods("DELETE")
	api.HandleFunc("/content/{id}/edit", contentMod.EditContent).Methods("PUT")
	api.HandleFunc("/content/{id}/pin", contentMod.PinContent).Methods("POST")
	api.HandleFunc("/content/{id}/unpin", contentMod.UnpinContent).Methods("POST")
	api.HandleFunc("/content/{id}/feature", contentMod.FeatureContent).Methods("POST")
	api.HandleFunc("/content/{id}/unfeature", contentMod.UnfeatureContent).Methods("POST")
	api.HandleFunc("/content/{id}/restore", contentMod.RestoreContent).Methods("POST")
	api.HandleFunc("/content/moderation-queue", contentMod.GetModerationQueue).Methods("GET")
	api.HandleFunc("/content/bulk-delete", contentMod.BulkDeleteContent).Methods("POST")

	// Platform Control
	api.HandleFunc("/features/{feature}/toggle", platform.ToggleFeature).Methods("POST")
	api.HandleFunc("/features", platform.GetAllFeatures).Methods("GET")
	api.HandleFunc("/notifications/broadcast", platform.BroadcastNotification).Methods("POST")
	api.HandleFunc("/maintenance/enable", platform.EnableMaintenanceMode).Methods("POST")
	api.HandleFunc("/maintenance/disable", platform.DisableMaintenanceMode).Methods("POST")
	api.HandleFunc("/platform/health", platform.GetPlatformHealth).Methods("GET")
	api.HandleFunc("/killswitch/{switch}", platform.ActivateKillSwitch).Methods("POST")
	api.HandleFunc("/killswitch/{switch}/deactivate", platform.DeactivateKillSwitch).Methods("POST")

	// Analytics
	api.HandleFunc("/analytics/live", analytics.GetLiveMetrics).Methods("GET")
	api.HandleFunc("/analytics/trending", analytics.GetTrendingContent).Methods("GET")
	api.HandleFunc("/analytics/users", analytics.GetUserMetrics).Methods("GET")
	api.HandleFunc("/analytics/engagement", analytics.GetEngagementMetrics).Methods("GET")
	api.HandleFunc("/analytics/revenue", analytics.GetRevenueMetrics).Methods("GET")

	// Security & Audit
	api.HandleFunc("/audit/logs", audit.GetAuditLogs).Methods("GET")
	api.HandleFunc("/audit/logs/{id}", audit.GetAuditLogDetail).Methods("GET")
	api.HandleFunc("/security/sessions", security.GetActiveSessions).Methods("GET")
	api.HandleFunc("/security/sessions/{id}/kill", security.KillSession).Methods("POST")
	api.HandleFunc("/security/devices", security.GetWhitelistedDevices).Methods("GET")
	api.HandleFunc("/security/devices/{id}/whitelist", security.WhitelistDevice).Methods("POST")
	api.HandleFunc("/security/devices/{id}/revoke", security.RevokeDevice).Methods("DELETE")
	api.HandleFunc("/security/ip-blocks", security.GetIPBlocks).Methods("GET")
	api.HandleFunc("/security/ip-blocks", security.BlockIPRange).Methods("POST")

	return r
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}
