package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/entativa/vignette/live-streaming-service/internal/handler"
	mediagrpc "github.com/entativa/vignette/live-streaming-service/internal/grpc"
	"github.com/entativa/vignette/live-streaming-service/internal/repository"
	"github.com/entativa/vignette/live-streaming-service/internal/service"
	"github.com/entativa/vignette/live-streaming-service/internal/websocket"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // In production: validate origin
	},
}

func main() {
	// Load environment
	godotenv.Load()

	// Database connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgresql://postgres:postgres@localhost:5432/vignette_streaming?sslmode=disable"
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

	// Initialize gRPC client for media service
	mediaServiceAddr := os.Getenv("MEDIA_SERVICE_GRPC_URL")
	if mediaServiceAddr == "" {
		mediaServiceAddr = "localhost:50051"
	}

	mediaGRPC, err := mediagrpc.NewMediaServiceClient(mediaServiceAddr)
	if err != nil {
		log.Printf("⚠️  Media service not available: %v", err)
		// Continue without media service (some features won't work)
	}
	defer func() {
		if mediaGRPC != nil {
			mediaGRPC.Close()
		}
	}()

	// Initialize WebSocket hub
	chatHub := websocket.NewChatHub()
	go chatHub.Run()

	// Initialize repositories
	streamRepo := repository.NewStreamRepository(db)
	viewerRepo := repository.NewViewerRepository(db)
	commentRepo := repository.NewCommentRepository(db)

	// Initialize services
	// Mock Redis and Kafka for now
	var mockRedis *service.RedisClient
	var mockKafka *service.KafkaProducer
	
	streamingService := service.NewStreamingService(
		streamRepo, viewerRepo, commentRepo,
		mediaGRPC, mockKafka, mockRedis,
	)

	// Initialize handlers
	streamHandler := handler.NewStreamHandler(streamingService)

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

	// WebSocket endpoint for live chat
	router.GET("/ws/stream/:stream_id", func(c *gin.Context) {
		streamID, err := uuid.Parse(c.Param("stream_id"))
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid stream ID"})
			return
		}

		userID, _ := uuid.Parse(c.GetHeader("X-User-ID"))

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("WebSocket upgrade failed: %v", err)
			return
		}

		client := &websocket.Client{
			conn:     conn,
			streamID: streamID,
			userID:   userID,
			send:     make(chan []byte, 256),
		}

		chatHub.register <- client

		go client.WritePump()
		go client.ReadPump()
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "vignette-live-streaming"})
	})

	// Root
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "vignette-live-streaming-service",
			"version": "1.0.0",
			"features": []string{
				"YouTube-quality streaming (up to 4K)",
				"Follower threshold (100 followers)",
				"Real-time comments",
				"Live reactions",
				"Viewer analytics",
				"Stream recording (VOD)",
				"gRPC integration with media service",
				"WebSocket live chat",
			},
		})
	})

	// API v1
	v1 := router.Group("/api/v1", authMiddleware)
	{
		// Stream management
		v1.POST("/streams", streamHandler.CreateStream)
		v1.GET("/streams/:id", streamHandler.GetStream)
		v1.POST("/streams/:id/start", streamHandler.StartStream)
		v1.POST("/streams/:id/end", streamHandler.EndStream)
		v1.GET("/streams/live", streamHandler.GetLiveStreams)
		v1.GET("/streams/eligibility", streamHandler.CheckEligibility)

		// Interactions
		v1.POST("/streams/:id/comments", streamHandler.PostComment)
		v1.POST("/streams/:id/reactions", streamHandler.AddReaction)

		// Analytics
		v1.GET("/streams/:id/analytics", streamHandler.GetAnalytics)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8097"
	}

	log.Printf("🚀 Vignette Live Streaming Service starting on port %s", port)
	log.Printf("📺 WebSocket available at ws://localhost:%s/ws/stream/:stream_id", port)
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
