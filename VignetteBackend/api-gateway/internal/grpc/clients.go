package grpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GRPCClients holds all gRPC service clients
type GRPCClients struct {
	UserConn       *grpc.ClientConn
	PostConn       *grpc.ClientConn
	MessagingConn  *grpc.ClientConn
	SettingsConn   *grpc.ClientConn
	MediaConn      *grpc.ClientConn
	StoryConn      *grpc.ClientConn
	SearchConn     *grpc.ClientConn
	NotificationConn *grpc.ClientConn
	FeedConn       *grpc.ClientConn
	CommunityConn  *grpc.ClientConn
	RecommendationConn *grpc.ClientConn
	StreamingConn  *grpc.ClientConn
	CreatorConn    *grpc.ClientConn
}

// ServiceConfig holds service endpoint configuration
type ServiceConfig struct {
	UserServiceURL          string
	PostServiceURL          string
	MessagingServiceURL     string
	SettingsServiceURL      string
	MediaServiceURL         string
	StoryServiceURL         string
	SearchServiceURL        string
	NotificationServiceURL  string
	FeedServiceURL          string
	CommunityServiceURL     string
	RecommendationServiceURL string
	StreamingServiceURL     string
	CreatorServiceURL       string
}

// NewGRPCClients creates gRPC connections to all services
func NewGRPCClients(config *ServiceConfig) (*GRPCClients, error) {
	clients := &GRPCClients{}
	
	// Connect to User Service
	userConn, err := createConnection(config.UserServiceURL, "User Service")
	if err != nil {
		log.Printf("Warning: %v", err)
	} else {
		clients.UserConn = userConn
	}
	
	// Connect to Post Service
	postConn, err := createConnection(config.PostServiceURL, "Post Service")
	if err != nil {
		log.Printf("Warning: %v", err)
	} else {
		clients.PostConn = postConn
	}
	
	// Connect to Messaging Service
	messagingConn, err := createConnection(config.MessagingServiceURL, "Messaging Service")
	if err != nil {
		log.Printf("Warning: %v", err)
	} else {
		clients.MessagingConn = messagingConn
	}
	
	// Connect to Settings Service
	settingsConn, err := createConnection(config.SettingsServiceURL, "Settings Service")
	if err != nil {
		log.Printf("Warning: %v", err)
	} else {
		clients.SettingsConn = settingsConn
	}
	
	// Connect to Media Service
	mediaConn, err := createConnection(config.MediaServiceURL, "Media Service")
	if err != nil {
		log.Printf("Warning: %v", err)
	} else {
		clients.MediaConn = mediaConn
	}
	
	// Connect to Story Service
	storyConn, err := createConnection(config.StoryServiceURL, "Story Service")
	if err != nil {
		log.Printf("Warning: %v", err)
	} else {
		clients.StoryConn = storyConn
	}
	
	// Connect to Search Service
	searchConn, err := createConnection(config.SearchServiceURL, "Search Service")
	if err != nil {
		log.Printf("Warning: %v", err)
	} else {
		clients.SearchConn = searchConn
	}
	
	// Connect to Notification Service
	notificationConn, err := createConnection(config.NotificationServiceURL, "Notification Service")
	if err != nil {
		log.Printf("Warning: %v", err)
	} else {
		clients.NotificationConn = notificationConn
	}
	
	// Connect to Feed Service
	feedConn, err := createConnection(config.FeedServiceURL, "Feed Service")
	if err != nil {
		log.Printf("Warning: %v", err)
	} else {
		clients.FeedConn = feedConn
	}
	
	// Connect to Community Service
	communityConn, err := createConnection(config.CommunityServiceURL, "Community Service")
	if err != nil {
		log.Printf("Warning: %v", err)
	} else {
		clients.CommunityConn = communityConn
	}
	
	// Connect to Recommendation Service
	recommendationConn, err := createConnection(config.RecommendationServiceURL, "Recommendation Service")
	if err != nil {
		log.Printf("Warning: %v", err)
	} else {
		clients.RecommendationConn = recommendationConn
	}
	
	// Connect to Streaming Service
	streamingConn, err := createConnection(config.StreamingServiceURL, "Streaming Service")
	if err != nil {
		log.Printf("Warning: %v", err)
	} else {
		clients.StreamingConn = streamingConn
	}
	
	// Connect to Creator Service
	creatorConn, err := createConnection(config.CreatorServiceURL, "Creator Service")
	if err != nil {
		log.Printf("Warning: %v", err)
	} else {
		clients.CreatorConn = creatorConn
	}
	
	log.Println("✅ gRPC clients initialized")
	return clients, nil
}

func createConnection(url, serviceName string) (*grpc.ClientConn, error) {
	if url == "" {
		return nil, fmt.Errorf("%s URL not configured", serviceName)
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	conn, err := grpc.DialContext(
		ctx,
		url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s at %s: %w", serviceName, url, err)
	}
	
	log.Printf("✅ Connected to %s at %s", serviceName, url)
	return conn, nil
}

// Close closes all gRPC connections
func (c *GRPCClients) Close() {
	if c.UserConn != nil {
		c.UserConn.Close()
	}
	if c.PostConn != nil {
		c.PostConn.Close()
	}
	if c.MessagingConn != nil {
		c.MessagingConn.Close()
	}
	if c.SettingsConn != nil {
		c.SettingsConn.Close()
	}
	if c.MediaConn != nil {
		c.MediaConn.Close()
	}
	if c.StoryConn != nil {
		c.StoryConn.Close()
	}
	if c.SearchConn != nil {
		c.SearchConn.Close()
	}
	if c.NotificationConn != nil {
		c.NotificationConn.Close()
	}
	if c.FeedConn != nil {
		c.FeedConn.Close()
	}
	if c.CommunityConn != nil {
		c.CommunityConn.Close()
	}
	if c.RecommendationConn != nil {
		c.RecommendationConn.Close()
	}
	if c.StreamingConn != nil {
		c.StreamingConn.Close()
	}
	if c.CreatorConn != nil {
		c.CreatorConn.Close()
	}
	
	log.Println("✅ All gRPC connections closed")
}
