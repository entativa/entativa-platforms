package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"messaging-service/internal/logger"
	"messaging-service/internal/repository"
	"messaging-service/internal/util"
)

type ConversationHandler struct {
	conversationRepo *repository.ConversationRepository
	messageRepo      *repository.MessageRepository
	logger           *logger.Logger
}

func NewConversationHandler(
	conversationRepo *repository.ConversationRepository,
	messageRepo *repository.MessageRepository,
	logger *logger.Logger,
) *ConversationHandler {
	return &ConversationHandler{
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
		logger:           logger,
	}
}

// GetConversations returns all conversations for a user
func (h *ConversationHandler) GetConversations(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*repository.User)
	if !ok {
		util.RespondWithUnauthorized(w, "")
		return
	}

	conversations, err := h.conversationRepo.GetUserConversations(r.Context(), user.ID)
	if err != nil {
		h.logger.Error("Failed to fetch conversations", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch conversations")
		return
	}

	util.RespondWithSuccess(w, "", map[string]interface{}{
		"conversations": conversations,
	})
}

// GetConversation returns a specific conversation
func (h *ConversationHandler) GetConversation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	conversationID := vars["id"]

	user, ok := r.Context().Value("user").(*repository.User)
	if !ok {
		util.RespondWithUnauthorized(w, "")
		return
	}

	conversation, err := h.conversationRepo.GetByID(r.Context(), conversationID)
	if err != nil {
		util.RespondWithNotFound(w, "Conversation not found")
		return
	}

	// Verify user is participant
	participants, err := h.conversationRepo.GetParticipants(r.Context(), conversationID)
	if err != nil || !contains(participants, user.ID) {
		util.RespondWithForbidden(w, "Access denied")
		return
	}

	util.RespondWithSuccess(w, "", conversation)
}

// CreateConversation creates a new conversation
func (h *ConversationHandler) CreateConversation(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*repository.User)
	if !ok {
		util.RespondWithUnauthorized(w, "")
		return
	}

	var req struct {
		ParticipantIDs []string `json:"participant_ids"`
		Type           string   `json:"type"` // 'direct' or 'group'
		Name           string   `json:"name,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate
	if len(req.ParticipantIDs) == 0 {
		util.RespondWithError(w, http.StatusBadRequest, "At least one participant required")
		return
	}

	if req.Type == "" {
		req.Type = "direct"
	}

	if req.Type == "direct" && len(req.ParticipantIDs) != 1 {
		util.RespondWithError(w, http.StatusBadRequest, "Direct conversations require exactly 1 other participant")
		return
	}

	// Check if direct conversation already exists
	if req.Type == "direct" {
		existing, err := h.conversationRepo.FindDirectConversation(r.Context(), user.ID, req.ParticipantIDs[0])
		if err == nil && existing != nil {
			util.RespondWithSuccess(w, "Conversation exists", existing)
			return
		}
	}

	// Add current user to participants
	allParticipants := append(req.ParticipantIDs, user.ID)

	// Create conversation
	conversation := &repository.Conversation{
		ID:   util.GenerateUUID(),
		Type: req.Type,
		Name: req.Name,
	}

	if err := h.conversationRepo.CreateConversation(r.Context(), conversation, allParticipants); err != nil {
		h.logger.Error("Failed to create conversation", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to create conversation")
		return
	}

	util.RespondWithCreated(w, "Conversation created", conversation)
}

// MarkAsRead marks a conversation as read
func (h *ConversationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	conversationID := vars["id"]

	user, ok := r.Context().Value("user").(*repository.User)
	if !ok {
		util.RespondWithUnauthorized(w, "")
		return
	}

	if err := h.conversationRepo.MarkAsRead(r.Context(), conversationID, user.ID); err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to mark as read")
		return
	}

	util.RespondWithSuccess(w, "Marked as read", nil)
}
