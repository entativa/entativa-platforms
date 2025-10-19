package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"messaging-service/internal/encryption"
	"messaging-service/internal/logger"
	"messaging-service/internal/repository"
	"messaging-service/internal/util"
	ws "messaging-service/internal/websocket"
)

type MessageHandler struct {
	messageRepo       *repository.MessageRepository
	conversationRepo  *repository.ConversationRepository
	encryptionService *encryption.SignalProtocolService
	wsHub             *ws.Hub
	logger            *logger.Logger
}

func NewMessageHandler(
	messageRepo *repository.MessageRepository,
	conversationRepo *repository.ConversationRepository,
	encryptionService *encryption.SignalProtocolService,
	wsHub *ws.Hub,
	logger *logger.Logger,
) *MessageHandler {
	return &MessageHandler{
		messageRepo:       messageRepo,
		conversationRepo:  conversationRepo,
		encryptionService: encryptionService,
		wsHub:             wsHub,
		logger:            logger,
	}
}

// SendMessage handles sending a new message
func (h *MessageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	conversationID := vars["conversationID"]

	user, ok := r.Context().Value("user").(*repository.User)
	if !ok {
		util.RespondWithUnauthorized(w, "")
		return
	}

	var req struct {
		Content     string   `json:"content"`
		ContentType string   `json:"content_type"`
		MediaURL    string   `json:"media_url,omitempty"`
		ReplyTo     string   `json:"reply_to,omitempty"`
		ExpiresIn   int      `json:"expires_in,omitempty"` // Seconds (disappearing messages)
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate
	if req.Content == "" && req.MediaURL == "" {
		util.RespondWithError(w, http.StatusBadRequest, "Message content or media required")
		return
	}

	if req.ContentType == "" {
		req.ContentType = "text"
	}

	// Get conversation participants
	participants, err := h.conversationRepo.GetParticipants(r.Context(), conversationID)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to get participants")
		return
	}

	// Verify user is participant
	isParticipant := false
	for _, p := range participants {
		if p == user.ID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		util.RespondWithForbidden(w, "You are not a participant in this conversation")
		return
	}

	// Encrypt message for each recipient (except sender)
	encryptedMessages := make(map[string]string)
	for _, participantID := range participants {
		if participantID == user.ID {
			continue
		}

		// Encrypt using Signal Protocol
		encrypted, err := h.encryptionService.EncryptMessage(user.ID, participantID, req.Content)
		if err != nil {
			h.logger.Error(fmt.Sprintf("Failed to encrypt message for %s", participantID), err)
			// Continue with other recipients
			continue
		}

		encryptedMessages[participantID] = encrypted
	}

	// Store message
	var expiresAt *time.Time
	if req.ExpiresIn > 0 {
		t := time.Now().Add(time.Duration(req.ExpiresIn) * time.Second)
		expiresAt = &t
	}

	message := &repository.Message{
		ID:               util.GenerateUUID(),
		ConversationID:   conversationID,
		SenderID:         user.ID,
		EncryptedContent: encryptedMessages, // Map of recipientID -> encrypted content
		ContentType:      req.ContentType,
		MediaURL:         req.MediaURL,
		ReplyTo:          req.ReplyTo,
		SentAt:           time.Now(),
		ExpiresAt:        expiresAt,
	}

	if err := h.messageRepo.CreateMessage(r.Context(), message); err != nil {
		h.logger.Error("Failed to create message", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to send message")
		return
	}

	// Update conversation last message time
	if err := h.conversationRepo.UpdateLastMessage(r.Context(), conversationID, time.Now()); err != nil {
		h.logger.Error("Failed to update conversation", err)
	}

	// Broadcast to WebSocket clients
	wsMessage := &ws.Message{
		Type: "message",
		Payload: map[string]interface{}{
			"message_id":      message.ID,
			"conversation_id": conversationID,
			"sender_id":       user.ID,
			"content_type":    req.ContentType,
			"sent_at":         message.SentAt.Unix(),
		},
	}

	// Send to each participant
	for _, participantID := range participants {
		if participantID != user.ID {
			// Add encrypted content for this recipient
			wsMessage.Payload["encrypted_content"] = encryptedMessages[participantID]
			h.wsHub.SendToUser(participantID, wsMessage)
		}
	}

	util.RespondWithCreated(w, "Message sent", map[string]interface{}{
		"message_id": message.ID,
		"sent_at":    message.SentAt.Unix(),
	})
}

// GetMessages retrieves messages for a conversation
func (h *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	conversationID := vars["conversationID"]

	user, ok := r.Context().Value("user").(*repository.User)
	if !ok {
		util.RespondWithUnauthorized(w, "")
		return
	}

	// Verify user is participant
	participants, err := h.conversationRepo.GetParticipants(r.Context(), conversationID)
	if err != nil || !contains(participants, user.ID) {
		util.RespondWithForbidden(w, "Access denied")
		return
	}

	// Get pagination params
	page := 1
	limit := 50
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		fmt.Sscanf(pageStr, "%d", &page)
	}

	// Fetch messages
	messages, err := h.messageRepo.GetMessages(r.Context(), conversationID, user.ID, page, limit)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch messages")
		return
	}

	// Decrypt messages for this user
	decryptedMessages := make([]map[string]interface{}, 0)
	for _, msg := range messages {
		// Decrypt content
		var plaintext string
		if msg.SenderID == user.ID {
			// Sender can read their own message (not encrypted for themselves)
			plaintext = msg.EncryptedContent[user.ID] // Or store plaintext separately
		} else {
			// Decrypt using Signal Protocol
			encrypted := msg.EncryptedContent[user.ID]
			plaintext, err = h.encryptionService.DecryptMessage(user.ID, msg.SenderID, encrypted)
			if err != nil {
				h.logger.Error("Failed to decrypt message", err)
				plaintext = "[Failed to decrypt]"
			}
		}

		decryptedMessages = append(decryptedMessages, map[string]interface{}{
			"id":              msg.ID,
			"sender_id":       msg.SenderID,
			"content":         plaintext,
			"content_type":    msg.ContentType,
			"media_url":       msg.MediaURL,
			"reply_to":        msg.ReplyTo,
			"sent_at":         msg.SentAt.Unix(),
			"delivered_at":    msg.DeliveredAt,
			"read_at":         msg.ReadAt,
			"is_sender":       msg.SenderID == user.ID,
		})
	}

	util.RespondWithSuccess(w, "", map[string]interface{}{
		"messages": decryptedMessages,
		"page":     page,
		"limit":    limit,
	})
}

// MarkDelivered marks a message as delivered
func (h *MessageHandler) MarkDelivered(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageID := vars["id"]

	user, ok := r.Context().Value("user").(*repository.User)
	if !ok {
		util.RespondWithUnauthorized(w, "")
		return
	}

	if err := h.messageRepo.MarkDelivered(r.Context(), messageID, user.ID); err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to mark as delivered")
		return
	}

	// Notify sender via WebSocket
	message, _ := h.messageRepo.GetByID(r.Context(), messageID)
	if message != nil {
		wsMsg := &ws.Message{
			Type: "delivered",
			Payload: map[string]interface{}{
				"message_id": messageID,
				"user_id":    user.ID,
				"timestamp":  time.Now().Unix(),
			},
		}
		h.wsHub.SendToUser(message.SenderID, wsMsg)
	}

	util.RespondWithSuccess(w, "Marked as delivered", nil)
}

// MarkRead marks a message as read
func (h *MessageHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageID := vars["id"]

	user, ok := r.Context().Value("user").(*repository.User)
	if !ok {
		util.RespondWithUnauthorized(w, "")
		return
	}

	if err := h.messageRepo.MarkRead(r.Context(), messageID, user.ID); err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to mark as read")
		return
	}

	// Notify sender via WebSocket
	message, _ := h.messageRepo.GetByID(r.Context(), messageID)
	if message != nil {
		wsMsg := &ws.Message{
			Type: "read",
			Payload: map[string]interface{}{
				"message_id": messageID,
				"user_id":    user.ID,
				"timestamp":  time.Now().Unix(),
			},
		}
		h.wsHub.SendToUser(message.SenderID, wsMsg)
	}

	util.RespondWithSuccess(w, "Marked as read", nil)
}

// DeleteMessage soft-deletes a message
func (h *MessageHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageID := vars["id"]

	user, ok := r.Context().Value("user").(*repository.User)
	if !ok {
		util.RespondWithUnauthorized(w, "")
		return
	}

	// Verify user is sender
	message, err := h.messageRepo.GetByID(r.Context(), messageID)
	if err != nil || message.SenderID != user.ID {
		util.RespondWithForbidden(w, "Only sender can delete messages")
		return
	}

	if err := h.messageRepo.DeleteMessage(r.Context(), messageID); err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to delete message")
		return
	}

	// Notify participants via WebSocket
	participants, _ := h.conversationRepo.GetParticipants(r.Context(), message.ConversationID)
	wsMsg := &ws.Message{
		Type: "message_deleted",
		Payload: map[string]interface{}{
			"message_id":      messageID,
			"conversation_id": message.ConversationID,
			"timestamp":       time.Now().Unix(),
		},
	}

	for _, participantID := range participants {
		if participantID != user.ID {
			h.wsHub.SendToUser(participantID, wsMsg)
		}
	}

	util.RespondWithSuccess(w, "Message deleted", nil)
}

// SendTypingIndicator sends a typing indicator
func (h *MessageHandler) SendTypingIndicator(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	conversationID := vars["conversationID"]

	user, ok := r.Context().Value("user").(*repository.User)
	if !ok {
		util.RespondWithUnauthorized(w, "")
		return
	}

	var req struct {
		IsTyping bool `json:"is_typing"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// Broadcast typing indicator via WebSocket
	wsMsg := &ws.Message{
		Type: "typing",
		Payload: map[string]interface{}{
			"conversation_id": conversationID,
			"user_id":         user.ID,
			"is_typing":       req.IsTyping,
			"timestamp":       time.Now().Unix(),
		},
	}

	// Send to all participants
	participants, _ := h.conversationRepo.GetParticipants(r.Context(), conversationID)
	for _, participantID := range participants {
		if participantID != user.ID {
			h.wsHub.SendToUser(participantID, wsMsg)
		}
	}

	util.RespondWithSuccess(w, "", nil)
}

// UploadMedia handles encrypted media upload
func (h *MessageHandler) UploadMedia(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max
		util.RespondWithError(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	file, header, err := r.FormFile("media")
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Media file required")
		return
	}
	defer file.Close()

	// TODO: Upload encrypted media to S3/CDN
	// For now, return placeholder URL
	mediaURL := fmt.Sprintf("https://cdn.entativa.com/media/%s", util.GenerateUUID())

	util.RespondWithSuccess(w, "Media uploaded", map[string]interface{}{
		"media_url": mediaURL,
		"filename":  header.Filename,
		"size":      header.Size,
	})
}

// Helper function
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
