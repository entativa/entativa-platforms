package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"messaging-service/internal/encryption"
	"messaging-service/internal/logger"
	"messaging-service/internal/repository"
)

// Hub maintains active WebSocket connections and handles message broadcasting
type Hub struct {
	// Registered clients (userID -> *Client)
	clients map[string]*Client
	
	// Register requests from clients
	register chan *Client
	
	// Unregister requests from clients
	unregister chan *Client
	
	// Inbound messages from clients
	broadcast chan *Message
	
	// Mutex for thread-safe access
	mu sync.RWMutex
	
	// Dependencies
	messageRepo       *repository.MessageRepository
	conversationRepo  *repository.ConversationRepository
	encryptionService *encryption.SignalProtocolService
	logger            *logger.Logger
}

// Client represents a connected WebSocket client
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	userID   string
	username string
}

// Message types
type Message struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, validate origin
	},
}

func NewHub(
	messageRepo *repository.MessageRepository,
	conversationRepo *repository.ConversationRepository,
	encryptionService *encryption.SignalProtocolService,
	logger *logger.Logger,
) *Hub {
	return &Hub{
		clients:           make(map[string]*Client),
		register:          make(chan *Client),
		unregister:        make(chan *Client),
		broadcast:         make(chan *Message, 256),
		messageRepo:       messageRepo,
		conversationRepo:  conversationRepo,
		encryptionService: encryptionService,
		logger:            logger,
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.userID] = client
			h.mu.Unlock()
			h.logger.Info(fmt.Sprintf("Client registered: %s", client.userID))
			
			// Send online status to contacts
			h.broadcastPresence(client.userID, "online")
			
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)
				close(client.send)
			}
			h.mu.Unlock()
			h.logger.Info(fmt.Sprintf("Client unregistered: %s", client.userID))
			
			// Send offline status to contacts
			h.broadcastPresence(client.userID, "offline")
			
		case message := <-h.broadcast:
			h.handleMessage(message)
		}
	}
}

// ServeWS handles WebSocket requests from clients
func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	user, ok := r.Context().Value("user").(*repository.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Upgrade connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("Failed to upgrade connection", err)
		return
	}

	// Create client
	client := &Client{
		hub:      h,
		conn:     conn,
		send:     make(chan []byte, 256),
		userID:   user.ID,
		username: user.Username,
	}

	h.register <- client

	// Start goroutines
	go client.writePump()
	go client.readPump()
}

// handleMessage processes incoming messages
func (h *Hub) handleMessage(msg *Message) {
	switch msg.Type {
	case "message":
		h.handleChatMessage(msg)
	case "typing":
		h.handleTypingIndicator(msg)
	case "read":
		h.handleReadReceipt(msg)
	case "delivered":
		h.handleDeliveryReceipt(msg)
	default:
		h.logger.Warn(fmt.Sprintf("Unknown message type: %s", msg.Type))
	}
}

// handleChatMessage handles incoming chat messages
func (h *Hub) handleChatMessage(msg *Message) {
	recipientID, ok := msg.Payload["recipient_id"].(string)
	if !ok {
		return
	}

	// Get recipient client
	h.mu.RLock()
	recipient, isOnline := h.clients[recipientID]
	h.mu.RUnlock()

	// Marshal message
	data, err := json.Marshal(msg)
	if err != nil {
		h.logger.Error("Failed to marshal message", err)
		return
	}

	// If recipient is online, send immediately
	if isOnline {
		select {
		case recipient.send <- data:
		default:
			// Send buffer full, close connection
			h.unregister <- recipient
		}
	}

	// Message is stored in database regardless (for offline delivery)
}

// handleTypingIndicator handles typing indicators
func (h *Hub) handleTypingIndicator(msg *Message) {
	conversationID, ok := msg.Payload["conversation_id"].(string)
	if !ok {
		return
	}

	senderID, ok := msg.Payload["sender_id"].(string)
	if !ok {
		return
	}

	// Get all participants in conversation
	participants, err := h.conversationRepo.GetParticipants(conversationID)
	if err != nil {
		h.logger.Error("Failed to get participants", err)
		return
	}

	// Broadcast to all participants except sender
	data, _ := json.Marshal(msg)
	h.mu.RLock()
	for _, participant := range participants {
		if participant != senderID {
			if client, ok := h.clients[participant]; ok {
				select {
				case client.send <- data:
				default:
				}
			}
		}
	}
	h.mu.RUnlock()
}

// handleReadReceipt handles read receipts
func (h *Hub) handleReadReceipt(msg *Message) {
	// Similar to typing indicator
	h.handleTypingIndicator(msg)
}

// handleDeliveryReceipt handles delivery receipts
func (h *Hub) handleDeliveryReceipt(msg *Message) {
	// Similar to typing indicator
	h.handleTypingIndicator(msg)
}

// broadcastPresence sends presence updates to contacts
func (h *Hub) broadcastPresence(userID, status string) {
	presenceMsg := Message{
		Type: "presence",
		Payload: map[string]interface{}{
			"user_id": userID,
			"status":  status,
			"timestamp": time.Now().Unix(),
		},
	}

	data, _ := json.Marshal(presenceMsg)

	// Get user's contacts (TODO: fetch from contacts service)
	// For now, broadcast to all connected clients
	h.mu.RLock()
	for _, client := range h.clients {
		if client.userID != userID {
			select {
			case client.send <- data:
			default:
			}
		}
	}
	h.mu.RUnlock()
}

// SendToUser sends a message to a specific user
func (h *Hub) SendToUser(userID string, message *Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	h.mu.RLock()
	client, ok := h.clients[userID]
	h.mu.RUnlock()

	if !ok {
		return fmt.Errorf("user not connected")
	}

	select {
	case client.send <- data:
		return nil
	default:
		return fmt.Errorf("send buffer full")
	}
}

// Client read pump
const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512 * 1024 // 512 KB
)

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, messageData, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.hub.logger.Error("WebSocket error", err)
			}
			break
		}

		// Parse message
		var msg Message
		if err := json.Unmarshal(messageData, &msg); err != nil {
			c.hub.logger.Error("Failed to unmarshal message", err)
			continue
		}

		// Add sender ID to payload
		msg.Payload["sender_id"] = c.userID

		// Broadcast message
		c.hub.broadcast <- &msg
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
