package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"user-service/internal/repository"
	"user-service/internal/util"
)

// TakesHandler handles takes/reels requests
type TakesHandler struct {
	takesRepo *repository.TakesRepository
	logger    *logger.Logger
}

// NewTakesHandler creates a new takes handler
func NewTakesHandler(takesRepo *repository.TakesRepository, logger *logger.Logger) *TakesHandler {
	return &TakesHandler{
		takesRepo: takesRepo,
		logger:    logger,
	}
}

// TakeResponse represents a take in API responses
type TakeResponse struct {
	ID            string   `json:"id"`
	UserID        string   `json:"user_id"`
	Username      string   `json:"username"`
	UserAvatar    string   `json:"user_avatar,omitempty"`
	VideoURL      string   `json:"video_url"`
	ThumbnailURL  string   `json:"thumbnail_url,omitempty"`
	Caption       string   `json:"caption"`
	AudioName     string   `json:"audio_name"`
	AudioURL      string   `json:"audio_url,omitempty"`
	Duration      int      `json:"duration"`
	LikesCount    int      `json:"likes_count"`
	CommentsCount int      `json:"comments_count"`
	SharesCount   int      `json:"shares_count"`
	ViewsCount    int      `json:"views_count"`
	IsLiked       bool     `json:"is_liked"`
	IsSaved       bool     `json:"is_saved"`
	Hashtags      []string `json:"hashtags,omitempty"`
	CreatedAt     string   `json:"created_at"`
}

// GetFeed returns takes feed for user
func (h *TakesHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
	// Get pagination params
	page := 1
	limit := 10
	
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}
	
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l <= 50 {
			limit = l
		}
	}
	
	// Get current user ID from context (if authenticated)
	currentUserID := ""
	if user, ok := r.Context().Value("user").(*repository.User); ok {
		currentUserID = user.ID
	}
	
	// Fetch takes feed
	takes, err := h.takesRepo.GetFeed(r.Context(), currentUserID, page, limit)
	if err != nil {
		h.logger.Error("Failed to fetch takes feed", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch takes")
		return
	}
	
	// Map to response
	takesResponse := make([]TakeResponse, len(takes))
	for i, take := range takes {
		takesResponse[i] = mapTakeToResponse(take)
	}
	
	util.RespondWithSuccess(w, "", map[string]interface{}{
		"takes":    takesResponse,
		"page":     page,
		"limit":    limit,
		"has_more": len(takes) == limit,
	})
}

// GetTakeByID returns a specific take
func (h *TakesHandler) GetTakeByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	takeID := vars["id"]
	
	// Get current user ID from context (if authenticated)
	currentUserID := ""
	if user, ok := r.Context().Value("user").(*repository.User); ok {
		currentUserID = user.ID
	}
	
	// Fetch take
	take, err := h.takesRepo.GetByID(r.Context(), takeID, currentUserID)
	if err != nil {
		util.RespondWithError(w, http.StatusNotFound, "Take not found")
		return
	}
	
	// Increment view count (async)
	go h.takesRepo.IncrementViews(r.Context(), takeID)
	
	util.RespondWithSuccess(w, "", mapTakeToResponse(take))
}

// LikeTake likes a take
func (h *TakesHandler) LikeTake(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	takeID := vars["id"]
	
	// Get user from context
	user, ok := r.Context().Value("user").(*repository.User)
	if !ok {
		util.RespondWithUnauthorized(w, "")
		return
	}
	
	// Toggle like
	err := h.takesRepo.LikeTake(r.Context(), takeID, user.ID)
	if err != nil {
		h.logger.Error("Failed to like take", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to like take")
		return
	}
	
	// Get updated take
	take, err := h.takesRepo.GetByID(r.Context(), takeID, user.ID)
	if err != nil {
		util.RespondWithError(w, http.StatusNotFound, "Take not found")
		return
	}
	
	util.RespondWithSuccess(w, "Take liked successfully", mapTakeToResponse(take))
}

// UnlikeTake unlikes a take
func (h *TakesHandler) UnlikeTake(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	takeID := vars["id"]
	
	// Get user from context
	user, ok := r.Context().Value("user").(*repository.User)
	if !ok {
		util.RespondWithUnauthorized(w, "")
		return
	}
	
	// Unlike
	err := h.takesRepo.UnlikeTake(r.Context(), takeID, user.ID)
	if err != nil {
		h.logger.Error("Failed to unlike take", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to unlike take")
		return
	}
	
	// Get updated take
	take, err := h.takesRepo.GetByID(r.Context(), takeID, user.ID)
	if err != nil {
		util.RespondWithError(w, http.StatusNotFound, "Take not found")
		return
	}
	
	util.RespondWithSuccess(w, "Take unliked successfully", mapTakeToResponse(take))
}

// GetComments returns comments for a take
func (h *TakesHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	takeID := vars["id"]
	
	// Get pagination params
	page := 1
	limit := 20
	
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}
	
	// Fetch comments
	comments, err := h.takesRepo.GetComments(r.Context(), takeID, page, limit)
	if err != nil {
		h.logger.Error("Failed to fetch comments", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch comments")
		return
	}
	
	util.RespondWithSuccess(w, "", map[string]interface{}{
		"comments": comments,
		"page":     page,
		"limit":    limit,
		"has_more": len(comments) == limit,
	})
}

// AddComment adds a comment to a take
func (h *TakesHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	takeID := vars["id"]
	
	// Get user from context
	user, ok := r.Context().Value("user").(*repository.User)
	if !ok {
		util.RespondWithUnauthorized(w, "")
		return
	}
	
	// Parse request
	var req struct {
		Text string `json:"text"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	if req.Text == "" {
		util.RespondWithError(w, http.StatusBadRequest, "Comment text is required")
		return
	}
	
	// Create comment
	comment := &repository.TakeComment{
		ID:        util.GenerateUUID(),
		TakeID:    takeID,
		UserID:    user.ID,
		Username:  user.Username,
		Text:      util.SanitizeInput(req.Text),
		CreatedAt: time.Now(),
	}
	
	err := h.takesRepo.AddComment(r.Context(), comment)
	if err != nil {
		h.logger.Error("Failed to add comment", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to add comment")
		return
	}
	
	util.RespondWithCreated(w, "Comment added successfully", comment)
}

// Helper function to map take to response
func mapTakeToResponse(take *repository.Take) TakeResponse {
	return TakeResponse{
		ID:            take.ID,
		UserID:        take.UserID,
		Username:      take.Username,
		UserAvatar:    take.UserAvatar,
		VideoURL:      take.VideoURL,
		ThumbnailURL:  take.ThumbnailURL,
		Caption:       take.Caption,
		AudioName:     take.AudioName,
		AudioURL:      take.AudioURL,
		Duration:      take.Duration,
		LikesCount:    take.LikesCount,
		CommentsCount: take.CommentsCount,
		SharesCount:   take.SharesCount,
		ViewsCount:    take.ViewsCount,
		IsLiked:       take.IsLiked,
		IsSaved:       take.IsSaved,
		Hashtags:      take.Hashtags,
		CreatedAt:     take.CreatedAt.Format(time.RFC3339),
	}
}
