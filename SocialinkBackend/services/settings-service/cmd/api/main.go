package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/entativa/socialink/settings-service/internal/crypto"
	"github.com/entativa/socialink/settings-service/internal/handler"
	"github.com/entativa/socialink/settings-service/internal/repository"
	"github.com/entativa/socialink/settings-service/internal/service"
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
		dbURL = "postgresql://postgres:postgres@localhost:5432/socialink_settings?sslmode=disable"
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
	settingsRepo := repository.NewSettingsRepository(db)
	keyBackupRepo := repository.NewKeyBackupRepository(db)

	// Initialize services
	encryptionSvc := crypto.NewEncryptionService()
	settingsService := service.NewSettingsService(settingsRepo, keyBackupRepo, encryptionSvc)

	// Initialize handlers
	settingsHandler := handler.NewSettingsHandler(settingsService)

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
		c.JSON(200, gin.H{"status": "healthy", "service": "socialink-settings-service"})
	})

	// Root
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "socialink-settings-service",
			"version": "1.0.0",
			"features": []string{
				"Comprehensive app settings (appearance, privacy, notifications, etc)",
				"Encrypted chat key backup with PIN/Passphrase",
				"Multiple storage options (Entativa servers, local, iCloud, Google Drive)",
				"Double-encryption (Signal + PIN/Passphrase)",
				"PBKDF2 key derivation (100,000 iterations)",
				"bcrypt password hashing",
				"AES-256-GCM encryption",
				"Security audit logging",
				"Settings change history",
				"Authorities only get metadata (keys are encrypted!)",
			},
		})
	})

	// API v1
	v1 := router.Group("/api/v1", authMiddleware)
	{
		// Settings
		v1.GET("/settings", settingsHandler.GetSettings)
		v1.PUT("/settings", settingsHandler.UpdateSettings)

		// Encrypted key backup
		v1.POST("/keys/backup", settingsHandler.CreateKeyBackup)
		v1.POST("/keys/restore", settingsHandler.RestoreKeyBackup)
		v1.GET("/keys/backup", settingsHandler.GetKeyBackupInfo)
		v1.DELETE("/keys/backup", settingsHandler.DeleteKeyBackup)

		// Storage location info
		v1.GET("/storage-locations", settingsHandler.GetStorageLocations)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8101"
	}

	log.Printf("🚀 Socialink Settings Service starting on port %s", port)
	log.Printf("🔐 Encrypted key backup enabled with PIN/Passphrase protection")
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
