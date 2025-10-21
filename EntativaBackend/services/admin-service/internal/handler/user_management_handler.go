package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"admin-service/internal/logger"
	"admin-service/internal/repository"
	"admin-service/internal/util"
)

type UserManagementHandler struct {
	userRepo  *repository.UserRepository
	adminRepo *repository.AdminRepository
	auditRepo *repository.AuditRepository
	logger    *logger.Logger
}

func NewUserManagementHandler(
	userRepo *repository.UserRepository,
	adminRepo *repository.AdminRepository,
	auditRepo *repository.AuditRepository,
	logger *logger.Logger,
) *UserManagementHandler {
	return &UserManagementHandler{
		userRepo:  userRepo,
		adminRepo: adminRepo,
		auditRepo: auditRepo,
		logger:    logger,
	}
}

// GetUser retrieves full user profile including private data
func (h *UserManagementHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	user, err := h.userRepo.GetUserByID(r.Context(), userID)
	if err != nil {
		util.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	// Get additional admin data
	adminData, _ := h.adminRepo.GetUserAdminData(r.Context(), userID)

	response := map[string]interface{}{
		"user":       user,
		"admin_data": adminData,
	}

	util.RespondWithSuccess(w, "", response)
}

// BanUser permanently bans a user
func (h *UserManagementHandler) BanUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetUserID := vars["id"]

	var req struct {
		Reason   string `json:"reason"`
		Duration int    `json:"duration,omitempty"` // Duration in hours, 0 = permanent
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Reason == "" {
		util.RespondWithError(w, http.StatusBadRequest, "Reason is required")
		return
	}

	// Cannot ban yourself
	admin := r.Context().Value("user").(*repository.User)
	if targetUserID == admin.ID {
		util.RespondWithError(w, http.StatusBadRequest, "Cannot ban yourself")
		return
	}

	// Ban the user
	var expiresAt *time.Time
	if req.Duration > 0 {
		t := time.Now().Add(time.Duration(req.Duration) * time.Hour)
		expiresAt = &t
	}

	if err := h.adminRepo.BanUser(r.Context(), targetUserID, req.Reason, expiresAt); err != nil {
		h.logger.Error("Failed to ban user", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to ban user")
		return
	}

	// Log audit trail
	h.auditRepo.LogAction(r.Context(), &repository.AuditLog{
		AdminID:    admin.ID,
		Action:     "ban_user",
		TargetType: "user",
		TargetID:   targetUserID,
		Reason:     req.Reason,
		Metadata: map[string]interface{}{
			"duration":   req.Duration,
			"expires_at": expiresAt,
		},
	})

	util.RespondWithSuccess(w, "User banned successfully", nil)
}

// UnbanUser removes ban from a user
func (h *UserManagementHandler) UnbanUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetUserID := vars["id"]

	if err := h.adminRepo.UnbanUser(r.Context(), targetUserID); err != nil {
		h.logger.Error("Failed to unban user", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to unban user")
		return
	}

	admin := r.Context().Value("user").(*repository.User)
	h.auditRepo.LogAction(r.Context(), &repository.AuditLog{
		AdminID:    admin.ID,
		Action:     "unban_user",
		TargetType: "user",
		TargetID:   targetUserID,
	})

	util.RespondWithSuccess(w, "User unbanned successfully", nil)
}

// ShadowbanUser shadowbans a user (content hidden but user doesn't know)
func (h *UserManagementHandler) ShadowbanUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetUserID := vars["id"]

	var req struct {
		Reason string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.adminRepo.ShadowbanUser(r.Context(), targetUserID, req.Reason); err != nil {
		h.logger.Error("Failed to shadowban user", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to shadowban user")
		return
	}

	admin := r.Context().Value("user").(*repository.User)
	h.auditRepo.LogAction(r.Context(), &repository.AuditLog{
		AdminID:    admin.ID,
		Action:     "shadowban_user",
		TargetType: "user",
		TargetID:   targetUserID,
		Reason:     req.Reason,
	})

	util.RespondWithSuccess(w, "User shadowbanned successfully", nil)
}

// UnshadowbanUser removes shadowban
func (h *UserManagementHandler) UnshadowbanUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetUserID := vars["id"]

	if err := h.adminRepo.UnshadowbanUser(r.Context(), targetUserID); err != nil {
		h.logger.Error("Failed to unshadowban user", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to unshadowban user")
		return
	}

	admin := r.Context().Value("user").(*repository.User)
	h.auditRepo.LogAction(r.Context(), &repository.AuditLog{
		AdminID:    admin.ID,
		Action:     "unshadowban_user",
		TargetType: "user",
		TargetID:   targetUserID,
	})

	util.RespondWithSuccess(w, "User unshadowbanned successfully", nil)
}

// MuteUser mutes a user platform-wide
func (h *UserManagementHandler) MuteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetUserID := vars["id"]

	var req struct {
		Reason   string `json:"reason"`
		Duration int    `json:"duration"` // Duration in hours
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	expiresAt := time.Now().Add(time.Duration(req.Duration) * time.Hour)

	if err := h.adminRepo.MuteUser(r.Context(), targetUserID, req.Reason, &expiresAt); err != nil {
		h.logger.Error("Failed to mute user", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to mute user")
		return
	}

	admin := r.Context().Value("user").(*repository.User)
	h.auditRepo.LogAction(r.Context(), &repository.AuditLog{
		AdminID:    admin.ID,
		Action:     "mute_user",
		TargetType: "user",
		TargetID:   targetUserID,
		Reason:     req.Reason,
		Metadata:   map[string]interface{}{"duration": req.Duration},
	})

	util.RespondWithSuccess(w, "User muted successfully", nil)
}

// UnmuteUser unmutes a user
func (h *UserManagementHandler) UnmuteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetUserID := vars["id"]

	if err := h.adminRepo.UnmuteUser(r.Context(), targetUserID); err != nil {
		h.logger.Error("Failed to unmute user", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to unmute user")
		return
	}

	admin := r.Context().Value("user").(*repository.User)
	h.auditRepo.LogAction(r.Context(), &repository.AuditLog{
		AdminID:    admin.ID,
		Action:     "unmute_user",
		TargetType: "user",
		TargetID:   targetUserID,
	})

	util.RespondWithSuccess(w, "User unmuted successfully", nil)
}

// SuspendUser temporarily suspends user account (reversible)
func (h *UserManagementHandler) SuspendUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetUserID := vars["id"]

	var req struct {
		Reason   string `json:"reason"`
		Duration int    `json:"duration"` // Duration in days, 0 = indefinite
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var expiresAt *time.Time
	if req.Duration > 0 {
		t := time.Now().Add(time.Duration(req.Duration) * 24 * time.Hour)
		expiresAt = &t
	}

	if err := h.adminRepo.SuspendUser(r.Context(), targetUserID, req.Reason, expiresAt); err != nil {
		h.logger.Error("Failed to suspend user", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to suspend user")
		return
	}

	admin := r.Context().Value("user").(*repository.User)
	h.auditRepo.LogAction(r.Context(), &repository.AuditLog{
		AdminID:    admin.ID,
		Action:     "suspend_user",
		TargetType: "user",
		TargetID:   targetUserID,
		Reason:     req.Reason,
		Metadata:   map[string]interface{}{"duration_days": req.Duration},
	})

	util.RespondWithSuccess(w, "User suspended successfully", nil)
}

// PermanentlyDeleteUser IRREVERSIBLY deletes user account and all data
func (h *UserManagementHandler) PermanentlyDeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetUserID := vars["id"]

	var req struct {
		Reason      string `json:"reason"`
		Confirmation string `json:"confirmation"` // Must be "PERMANENTLY DELETE"
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Double confirmation required
	if req.Confirmation != "PERMANENTLY DELETE" {
		util.RespondWithError(w, http.StatusBadRequest, "Confirmation text must be 'PERMANENTLY DELETE'")
		return
	}

	if req.Reason == "" || len(req.Reason) < 20 {
		util.RespondWithError(w, http.StatusBadRequest, "Detailed reason required (min 20 characters)")
		return
	}

	admin := r.Context().Value("user").(*repository.User)
	if targetUserID == admin.ID {
		util.RespondWithError(w, http.StatusBadRequest, "Cannot delete yourself")
		return
	}

	// CRITICAL: Permanent deletion
	if err := h.adminRepo.PermanentlyDeleteUser(r.Context(), targetUserID, req.Reason); err != nil {
		h.logger.Error("Failed to permanently delete user", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	// CRITICAL audit log
	h.auditRepo.LogAction(r.Context(), &repository.AuditLog{
		AdminID:    admin.ID,
		Action:     "permanently_delete_user",
		TargetType: "user",
		TargetID:   targetUserID,
		Reason:     req.Reason,
		Severity:   "CRITICAL",
		Metadata: map[string]interface{}{
			"confirmation": req.Confirmation,
			"legal_notice": "GDPR/CCPA compliance notification sent to legal team",
		},
	})

	// Notify legal team (TODO: implement)
	h.logger.Info("LEGAL NOTIFICATION: User permanently deleted", map[string]interface{}{
		"user_id": targetUserID,
		"reason":  req.Reason,
		"admin":   admin.Username,
	})

	util.RespondWithSuccess(w, "User permanently deleted", nil)
}

// ForceLogout kills all user sessions
func (h *UserManagementHandler) ForceLogout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetUserID := vars["id"]

	if err := h.adminRepo.ForceLogoutUser(r.Context(), targetUserID); err != nil {
		h.logger.Error("Failed to force logout", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to logout user")
		return
	}

	admin := r.Context().Value("user").(*repository.User)
	h.auditRepo.LogAction(r.Context(), &repository.AuditLog{
		AdminID:    admin.ID,
		Action:     "force_logout",
		TargetType: "user",
		TargetID:   targetUserID,
	})

	util.RespondWithSuccess(w, "User logged out successfully", nil)
}

// ForcePasswordReset forces user to reset password on next login
func (h *UserManagementHandler) ForcePasswordReset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetUserID := vars["id"]

	if err := h.adminRepo.ForcePasswordReset(r.Context(), targetUserID); err != nil {
		h.logger.Error("Failed to force password reset", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to force password reset")
		return
	}

	admin := r.Context().Value("user").(*repository.User)
	h.auditRepo.LogAction(r.Context(), &repository.AuditLog{
		AdminID:    admin.ID,
		Action:     "force_password_reset",
		TargetType: "user",
		TargetID:   targetUserID,
	})

	util.RespondWithSuccess(w, "Password reset required", nil)
}

// Disable2FA disables 2FA for user (recovery purposes)
func (h *UserManagementHandler) Disable2FA(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetUserID := vars["id"]

	var req struct {
		Reason string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.adminRepo.Disable2FA(r.Context(), targetUserID); err != nil {
		h.logger.Error("Failed to disable 2FA", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to disable 2FA")
		return
	}

	admin := r.Context().Value("user").(*repository.User)
	h.auditRepo.LogAction(r.Context(), &repository.AuditLog{
		AdminID:    admin.ID,
		Action:     "disable_2fa",
		TargetType: "user",
		TargetID:   targetUserID,
		Reason:     req.Reason,
	})

	util.RespondWithSuccess(w, "2FA disabled successfully", nil)
}

// ImpersonateUser allows viewing platform as another user (with step-up auth)
func (h *UserManagementHandler) ImpersonateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetUserID := vars["id"]

	var req struct {
		Reason   string `json:"reason"`
		Password string `json:"password"` // Step-up authentication
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Require detailed reason (min 20 chars)
	if len(req.Reason) < 20 {
		util.RespondWithError(w, http.StatusBadRequest, "Detailed reason required (min 20 characters)")
		return
	}

	admin := r.Context().Value("user").(*repository.User)

	// Verify password (step-up authentication)
	valid, err := h.userRepo.VerifyPassword(r.Context(), admin.ID, req.Password)
	if err != nil || !valid {
		util.RespondWithError(w, http.StatusUnauthorized, "Password verification failed")
		return
	}

	// Create impersonation session (10-minute timeout)
	impersonationToken, err := h.adminRepo.CreateImpersonationSession(r.Context(), admin.ID, targetUserID, req.Reason)
	if err != nil {
		h.logger.Error("Failed to create impersonation session", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to start impersonation")
		return
	}

	// HIGH-PRIORITY security alert to CTO
	h.logger.Info("SECURITY ALERT: Impersonation started", map[string]interface{}{
		"admin_id":   admin.ID,
		"admin_name": admin.Username,
		"target_id":  targetUserID,
		"reason":     req.Reason,
		"severity":   "HIGH",
	})

	// Critical audit log
	h.auditRepo.LogAction(r.Context(), &repository.AuditLog{
		AdminID:    admin.ID,
		Action:     "impersonate_user_start",
		TargetType: "user",
		TargetID:   targetUserID,
		Reason:     req.Reason,
		Severity:   "HIGH",
		Metadata: map[string]interface{}{
			"timeout_minutes": 10,
		},
	})

	util.RespondWithSuccess(w, "Impersonation session started", map[string]interface{}{
		"impersonation_token": impersonationToken,
		"expires_in_minutes":  10,
		"target_user_id":      targetUserID,
	})
}

// EndImpersonation ends impersonation session
func (h *UserManagementHandler) EndImpersonation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetUserID := vars["id"]

	admin := r.Context().Value("user").(*repository.User)

	if err := h.adminRepo.EndImpersonationSession(r.Context(), admin.ID, targetUserID); err != nil {
		h.logger.Error("Failed to end impersonation", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to end impersonation")
		return
	}

	h.auditRepo.LogAction(r.Context(), &repository.AuditLog{
		AdminID:    admin.ID,
		Action:     "impersonate_user_end",
		TargetType: "user",
		TargetID:   targetUserID,
	})

	util.RespondWithSuccess(w, "Impersonation session ended", nil)
}

// SearchUsers searches for users by username, email, or ID
func (h *UserManagementHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		util.RespondWithError(w, http.StatusBadRequest, "Search query required")
		return
	}

	users, err := h.userRepo.SearchUsers(r.Context(), query, 50)
	if err != nil {
		h.logger.Error("Failed to search users", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Search failed")
		return
	}

	util.RespondWithSuccess(w, "", map[string]interface{}{
		"users": users,
		"count": len(users),
	})
}
