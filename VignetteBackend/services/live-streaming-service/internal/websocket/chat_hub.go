package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/entativa/vignette/live-streaming-service/internal/model"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// ChatHub manages WebSocket connections for live stream chat
type ChatHub struct {
	// StreamID -> Map of clients
	streams map[uuid.UUID]map[*Client]bool
	
	// Register requests
	register chan *Client
	
	// Unregister requests
	unregister chan *Client
	
	// Broadcast messages to stream
	broadcast chan *BroadcastMessage
	
	mu sync.RWMutex
}

type BroadcastMessage struct {
	StreamID uuid.UUID
	Type     string      `json:"type"` // comment, reaction, viewer_update, streamer_message
	Data     interface{} `json:"data"`
}

type Client struct {
	hub      *ChatHub
	conn     *websocket.Conn
	streamID uuid.UUID
	userID   uuid.UUID
	send     chan []byte
}

func NewChatHub() *ChatHub {
	return &ChatHub{
		streams:    make(map[uuid.UUID]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *BroadcastMessage, 256),
	}
}

func (h *ChatHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.streams[client.streamID] == nil {
				h.streams[client.streamID] = make(map[*Client]bool)
			}
			h.streams[client.streamID][client] = true
			h.mu.Unlock()
			
			// Send viewer count update
			h.broadcastViewerCount(client.streamID)
		
		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.streams[client.streamID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.send)
					
					if len(clients) == 0 {
						delete(h.streams, client.streamID)
					}
				}
			}
			h.mu.Unlock()
			
			// Send viewer count update
			h.broadcastViewerCount(client.streamID)
		
		case message := <-h.broadcast:
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
				case client.send <- jsonData:
				default:
					// Client buffer full, disconnect
					close(client.send)
					delete(clients, client)
				}
			}
		}
	}
}

func (h *ChatHub) BroadcastComment(streamID uuid.UUID, comment *model.StreamComment) {
	h.broadcast <- &BroadcastMessage{
		StreamID: streamID,
		Type:     "comment",
		Data:     comment,
	}
}

func (h *ChatHub) BroadcastReaction(streamID uuid.UUID, reaction *model.StreamReaction) {
	h.broadcast <- &BroadcastMessage{
		StreamID: streamID,
		Type:     "reaction",
		Data:     reaction,
	}
}

func (h *ChatHub) broadcastViewerCount(streamID uuid.UUID) {
	h.mu.RLock()
	count := len(h.streams[streamID])
	h.mu.RUnlock()
	
	h.broadcast <- &BroadcastMessage{
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

// ReadPump pumps messages from WebSocket to hub
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		
		// In production: Parse message and handle
		// (e.g., comment, reaction, etc.)
		log.Printf("Received message from user %s: %s", c.userID, message)
	}
}

// WritePump pumps messages from hub to WebSocket
func (c *Client) WritePump() {
	defer c.conn.Close()

	for {
		message, ok := <-c.send
		if !ok {
			// Hub closed the channel
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return
		}
	}
}
