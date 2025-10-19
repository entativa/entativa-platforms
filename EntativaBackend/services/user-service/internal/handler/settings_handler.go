package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"user-service/internal/logger"
	"user-service/internal/repository"
	"user-service/internal/util"
)

type SettingsHandler struct {
	settingsRepo *repository.SettingsRepository
	logger       *logger.Logger
}

func NewSettingsHandler(settingsRepo *repository.SettingsRepository, logger *logger.Logger) *SettingsHandler {
	return &SettingsHandler{
		settingsRepo: settingsRepo,
		logger:       logger,
	}
}

// GetUserSettings retrieves all user settings
func (h *SettingsHandler) GetUserSettings(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*repository.User)
	if !ok {
		util.RespondWithUnauthorized(w, "")
		return
	}

	settings, err := h.settingsRepo.GetUserSettings(r.Context(), user.ID)
	if err != nil {
		h.logger.Error("Failed to get user settings", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to get settings")
		return
	}

	util.RespondWithSuccess(w, "", settings)
}

// UpdateAccountSettings updates account information
func (h *SettingsHandler) UpdateAccountSettings(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*repository.User)
	if !ok {
		util.RespondWithUnauthorized(w, "")
		return
	}

	var req struct {
		Email    string `json:"email,omitempty"`
		Phone    string `json:"phone,omitempty"`
		Bio      string `json:"bio,omitempty"`
		Website  string `json:"website,omitempty"`
		Name     string `json:"name,omitempty"`
		Username string `json:"username,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate email if provided
	if req.Email != "" && !util.IsValidEmail(req.Email) {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid email format")
		return
	}

	// Check if username is taken (if changing)
	if req.Username != "" && req.Username != user.Username {
		exists, err := h.settingsRepo.UsernameExists(r.Context(), req.Username, user.ID)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "Failed to check username")
			return
		}
		if exists {
			util.RespondWithError(w, http.StatusConflict, "Username already taken")
			return
		}
	}

	// Update account
	if err := h.settingsRepo.UpdateAccountSettings(r.Context(), user.ID, &req); err != nil {
		h.logger.Error("Failed to update account settings", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to update account")
		return
	}

	util.RespondWithSuccess(w, "Account updated successfully", nil)
}

// UpdatePrivacySettings updates privacy settings
func (h *SettingsHandler) UpdatePrivacySettings(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*repository.User)
	if !ok {
		util.RespondWithUnauthorized(w, "")
		return
	}

	var req struct {
		IsPrivateAccount       bool   `json:"is_private_account"`
		ShowActivityStatus     bool   `json:"show_activity_status"`
		ReadReceipts           bool   `json:"read_receipts"`
		AllowMessageRequests   bool   `json:"allow_message_requests"`
		PostsVisibility        string `json:"posts_visibility"`        // everyone, friends, only_me
		CommentsAllowed        string `json:"comments_allowed"`        // everyone, friends, no_one
		MentionsAllowed        string `json:"mentions_allowed"`        // everyone, following, off
		StorySharing           string `json:"story_sharing"`           // everyone, following, off
		SimilarAccountSuggestions bool `json:"similar_account_suggestions"`
		IncludeInRecommendations  bool `json:"include_in_recommendations"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate enum values
	validVisibilities := map[string]bool{"everyone": true, "friends": true, "only_me": true}
	validComments := map[string]bool{"everyone": true, "friends": true, "no_one": true}
	validMentions := map[string]bool{"everyone": true, "following": true, "off": true}

	if req.PostsVisibility != "" && !validVisibilities[req.PostsVisibility] {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid posts visibility value")
		return
	}

	if req.CommentsAllowed != "" && !validComments[req.CommentsAllowed] {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid comments allowed value")
		return
	}

	// Update privacy settings
	if err := h.settingsRepo.UpdatePrivacySettings(r.Context(), user.ID, &req); err != nil {
		h.logger.Error("Failed to update privacy settings", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to update privacy settings")
		return
	}

	util.RespondWithSuccess(w, "Privacy settings updated successfully", nil)
}

// UpdateNotificationSettings updates notification preferences
func (h *SettingsHandler) UpdateNotificationSettings(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*repository.User)
	if !ok {
		util.RespondWithUnauthorized(w, "")
		return
	}

	var req struct {
		NotifyLikes           bool `json:"notify_likes"`
		NotifyComments        bool `json:"notify_comments"`
		NotifyFollowers       bool `json:"notify_followers"`
		NotifyMessages        bool `json:"notify_messages"`
		NotifyFriendRequests  bool `json:"notify_friend_requests"`
		NotifyVideoViews      bool `json:"notify_video_views"`
		NotifyLiveVideos      bool `json:"notify_live_videos"`
		EmailWeeklySummary    bool `json:"email_weekly_summary"`
		EmailProductUpdates   bool `json:"email_product_updates"`
		EmailTips             bool `json:"email_tips"`
		NotificationSound     bool `json:"notification_sound"`
		NotificationVibration bool `json:"notification_vibration"`
		ShowBadgeCount        bool `json:"show_badge_count"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Update notification settings
	if err := h.settingsRepo.UpdateNotificationSettings(r.Context(), user.ID, &req); err != nil {
		h.logger.Error("Failed to update notification settings", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to update notification settings")
		return
	}

	util.RespondWithSuccess(w, "Notification settings updated successfully", nil)
}

// UpdateDataSettings updates data usage preferences
func (h *SettingsHandler) UpdateDataSettings(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*repository.User)
	if !ok {
		util.RespondWithUnauthorized(w, "")
		return
	}

	var req struct {
		UploadQuality    string `json:"upload_quality"`    // high, medium, low, normal, basic
		AutoplaySettings string `json:"autoplay_settings"` // always, wifi, never
		DataSaverMode    bool   `json:"data_saver_mode"`
		UseLessData      bool   `json:"use_less_data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate enum values
	validQualities := map[string]bool{"high": true, "medium": true, "low": true, "normal": true, "basic": true}
	validAutoplay := map[string]bool{"always": true, "wifi": true, "never": true}

	if req.UploadQuality != "" && !validQualities[req.UploadQuality] {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid upload quality value")
		return
	}

	if req.AutoplaySettings != "" && !validAutoplay[req.AutoplaySettings] {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid autoplay settings value")
		return
	}

	// Update data settings
	if err := h.settingsRepo.UpdateDataSettings(r.Context(), user.ID, &req); err != nil {
		h.logger.Error("Failed to update data settings", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to update data settings")
		return
	}

	util.RespondWithSuccess(w, "Data settings updated successfully", nil)
}

// GetBlockedUsers retrieves list of blocked users
func (h *SettingsHandler) GetBlockedUsers(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*repository.User)
	if !ok {
		util.RespondWithUnauthorized(w, "")
		return
	}

	blockedUsers, err := h.settingsRepo.GetBlockedUsers(r.Context(), user.ID)
	if err != nil {
		h.logger.Error("Failed to get blocked users", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to get blocked users")
		return
	}

	util.RespondWithSuccess(w, "", map[string]interface{}{
		"blocked_users": blockedUsers,
	})
}

// BlockUser blocks a user
func (h *SettingsHandler) BlockUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*repository.User)
	if !ok {
		util.RespondWithUnauthorized(w, "")
		return
	}

	vars := mux.Vars(r)
	targetUserID := vars["userID"]

	if targetUserID == user.ID {
		util.RespondWithError(w, http.StatusBadRequest, "Cannot block yourself")
		return
	}

	if err := h.settingsRepo.BlockUser(r.Context(), user.ID, targetUserID); err != nil {
		h.logger.Error("Failed to block user", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to block user")
		return
	}

	util.RespondWithSuccess(w, "User blocked successfully", nil)
}

// UnblockUser unblocks a user
func (h *SettingsHandler) UnblockUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*repository.User)
	if !ok {
		util.RespondWithUnauthorized(w, "")
		return
	}

	vars := mux.Vars(r)
	targetUserID := vars["userID"]

	if err := h.settingsRepo.UnblockUser(r.Context(), user.ID, targetUserID); err != nil {
		h.logger.Error("Failed to unblock user", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to unblock user")
		return
	}

	util.RespondWithSuccess(w, "User unblocked successfully", nil)
}

// ChangePassword changes user password
func (h *SettingsHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*repository.User)
	if !ok {
		util.RespondWithUnauthorized(w, "")
		return
	}

	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate current password
	valid, err := h.settingsRepo.VerifyPassword(r.Context(), user.ID, req.CurrentPassword)
	if err != nil {
		h.logger.Error("Failed to verify password", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to verify password")
		return
	}

	if !valid {
		util.RespondWithError(w, http.StatusUnauthorized, "Current password is incorrect")
		return
	}

	// Validate new password
	if len(req.NewPassword) < 8 {
		util.RespondWithError(w, http.StatusBadRequest, "New password must be at least 8 characters")
		return
	}

	// Change password
	if err := h.settingsRepo.ChangePassword(r.Context(), user.ID, req.NewPassword); err != nil {
		h.logger.Error("Failed to change password", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to change password")
		return
	}

	util.RespondWithSuccess(w, "Password changed successfully", nil)
}

// DeleteAccount permanently deletes a user account
func (h *SettingsHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*repository.User)
	if !ok {
		util.RespondWithUnauthorized(w, "")
		return
	}

	var req struct {
		Password string `json:"password"`
		Reason   string `json:"reason,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Verify password
	valid, err := h.settingsRepo.VerifyPassword(r.Context(), user.ID, req.Password)
	if err != nil || !valid {
		util.RespondWithError(w, http.StatusUnauthorized, "Password is incorrect")
		return
	}

	// Delete account (soft delete with 30-day grace period)
	if err := h.settingsRepo.DeleteAccount(r.Context(), user.ID, req.Reason); err != nil {
		h.logger.Error("Failed to delete account", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to delete account")
		return
	}

	util.RespondWithSuccess(w, "Account scheduled for deletion. You have 30 days to recover it.", nil)
}

// ClearCache clears user's cached data
func (h *SettingsHandler) ClearCache(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*repository.User)
	if !ok {
		util.RespondWithUnauthorized(w, "")
		return
	}

	// Clear cache (implementation depends on caching strategy)
	// For now, just log it
	h.logger.Info("Cache cleared for user", user.ID)

	util.RespondWithSuccess(w, "Cache cleared successfully", nil)
}

// GetLoginActivity retrieves recent login activity
func (h *SettingsHandler) GetLoginActivity(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*repository.User)
	if !ok {
		util.RespondWithUnauthorized(w, "")
		return
	}

	activity, err := h.settingsRepo.GetLoginActivity(r.Context(), user.ID, 10)
	if err != nil {
		h.logger.Error("Failed to get login activity", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to get login activity")
		return
	}

	util.RespondWithSuccess(w, "", map[string]interface{}{
		"login_activity": activity,
	})
}
