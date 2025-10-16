package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"vignette/live-streaming-service/internal/model"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// ChatHub manages WebSocket connections for live stream chat
type ChatHub struct {
	// StreamID -> Map of clients
	streams map[uuid.UUID]map[*Client]bool
	
	// Register requests
	Register chan *Client
	
	// Unregister requests
	Unregister chan *Client
	
	// Broadcast messages to stream
	Broadcast chan *BroadcastMessage
	
	mu sync.RWMutex
}

type BroadcastMessage struct {
	StreamID uuid.UUID
	Type     string      `json:"type"` // comment, reaction, viewer_update, streamer_message
	Data     interface{} `json:"data"`
}

type Client struct {
	Hub      *ChatHub
	Conn     *websocket.Conn
	StreamID uuid.UUID
	UserID   uuid.UUID
	Send     chan []byte
}

func NewChatHub() *ChatHub {
	return &ChatHub{
		streams:    make(map[uuid.UUID]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *BroadcastMessage, 256),
	}
}

func (h *ChatHub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if h.streams[client.StreamID] == nil {
				h.streams[client.StreamID] = make(map[*Client]bool)
			}
			h.streams[client.StreamID][client] = true
			h.mu.Unlock()
			
			// Send viewer count update
			h.broadcastViewerCount(client.StreamID)
		
		case client := <-h.Unregister:
			h.mu.Lock()
			if clients, ok := h.streams[client.StreamID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.Send)
					
					if len(clients) == 0 {
						delete(h.streams, client.StreamID)
					}
				}
			}
			h.mu.Unlock()
			
			// Send viewer count update
			h.broadcastViewerCount(client.StreamID)
		
		case message := <-h.Broadcast:
			h.mu.RLock()
			clients := h.streams[message.StreamID]
			h.mu.RUnlock()
			
			// Marshal message
			jsonData, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshaling message: %v", err)
				continue
			}
			
			// Send to all clients in stream
			for client := range clients {
				select {
				case client.Send <- jsonData:
				default:
					// Client buffer full, disconnect
					close(client.Send)
					delete(clients, client)
				}
			}
		}
	}
}

func (h *ChatHub) BroadcastComment(streamID uuid.UUID, comment *model.StreamComment) {
	h.Broadcast <- &BroadcastMessage{
		StreamID: streamID,
		Type:     "comment",
		Data:     comment,
	}
}

func (h *ChatHub) BroadcastReaction(streamID uuid.UUID, reaction *model.StreamReaction) {
	h.Broadcast <- &BroadcastMessage{
		StreamID: streamID,
		Type:     "reaction",
		Data:     reaction,
	}
}

func (h *ChatHub) broadcastViewerCount(streamID uuid.UUID) {
	h.mu.RLock()
	count := len(h.streams[streamID])
	h.mu.RUnlock()
	
	h.Broadcast <- &BroadcastMessage{
		StreamID: streamID,
		Type:     "viewer_update",
		Data: map[string]interface{}{
			"viewer_count": count,
		},
	}
}

func (h *ChatHub) GetViewerCount(streamID uuid.UUID) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.streams[streamID])
}

// NewClient creates a new WebSocket client
func NewClient(hub *ChatHub, conn *websocket.Conn, streamID, userID uuid.UUID) *Client {
	return &Client{
		Hub:      hub,
		Conn:     conn,
		StreamID: streamID,
		UserID:   userID,
		Send:     make(chan []byte, 256),
	}
}

// ReadPump pumps messages from WebSocket to hub
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		
		// In production: Parse message and handle
		// (e.g., comment, reaction, etc.)
		log.Printf("Received message from user %s: %s", c.UserID, message)
	}
}

// WritePump pumps messages from hub to WebSocket
func (c *Client) WritePump() {
	defer c.Conn.Close()

	for {
		message, ok := <-c.Send
		if !ok {
			// Hub closed the channel
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return
		}
	}
}
