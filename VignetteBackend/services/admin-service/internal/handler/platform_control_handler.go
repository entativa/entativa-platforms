package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"admin-service/internal/logger"
	"admin-service/internal/repository"
	"admin-service/internal/util"
)

type PlatformControlHandler struct {
	adminRepo *repository.AdminRepository
	auditRepo *repository.AuditRepository
	logger    *logger.Logger
}

func NewPlatformControlHandler(
	adminRepo *repository.AdminRepository,
	auditRepo *repository.AuditRepository,
	logger *logger.Logger,
) *PlatformControlHandler {
	return &PlatformControlHandler{
		adminRepo: adminRepo,
		auditRepo: auditRepo,
		logger:    logger,
	}
}

// ToggleFeature enables/disables a platform feature
func (h *PlatformControlHandler) ToggleFeature(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	featureName := vars["feature"]

	var req struct {
		Enabled            bool   `json:"enabled"`
		RolloutPercentage  int    `json:"rollout_percentage,omitempty"` // 0-100
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate rollout percentage
	if req.RolloutPercentage < 0 || req.RolloutPercentage > 100 {
		req.RolloutPercentage = 100
	}

	if err := h.adminRepo.ToggleFeature(r.Context(), featureName, req.Enabled, req.RolloutPercentage); err != nil {
		h.logger.Error("Failed to toggle feature", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to toggle feature")
		return
	}

	admin := r.Context().Value("user").(*repository.User)
	h.auditRepo.LogAction(r.Context(), &repository.AuditLog{
		AdminID:    admin.ID,
		Action:     "toggle_feature",
		TargetType: "platform",
		TargetID:   featureName,
		Metadata: map[string]interface{}{
			"enabled":    req.Enabled,
			"percentage": req.RolloutPercentage,
		},
	})

	util.RespondWithSuccess(w, "Feature toggled successfully", map[string]interface{}{
		"feature": featureName,
		"enabled": req.Enabled,
	})
}

// GetAllFeatures retrieves all feature flags
func (h *PlatformControlHandler) GetAllFeatures(w http.ResponseWriter, r *http.Request) {
	features, err := h.adminRepo.GetAllFeatures(r.Context())
	if err != nil {
		h.logger.Error("Failed to get features", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to get features")
		return
	}

	util.RespondWithSuccess(w, "", map[string]interface{}{
		"features": features,
	})
}

// BroadcastNotification sends a platform-wide notification
func (h *PlatformControlHandler) BroadcastNotification(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title          string `json:"title"`
		Message        string `json:"message"`
		Type           string `json:"type"`           // info, warning, critical, announcement
		TargetAudience string `json:"target_audience"` // all, premium, creators
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" || req.Message == "" {
		util.RespondWithError(w, http.StatusBadRequest, "Title and message required")
		return
	}

	if req.Type == "" {
		req.Type = "info"
	}
	if req.TargetAudience == "" {
		req.TargetAudience = "all"
	}

	admin := r.Context().Value("user").(*repository.User)

	// Create broadcast notification
	notificationID, err := h.adminRepo.CreateBroadcast(r.Context(), &repository.Broadcast{
		Title:          req.Title,
		Message:        req.Message,
		Type:           req.Type,
		TargetAudience: req.TargetAudience,
		CreatedBy:      admin.ID,
	})

	if err != nil {
		h.logger.Error("Failed to create broadcast", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to create broadcast")
		return
	}

	h.auditRepo.LogAction(r.Context(), &repository.AuditLog{
		AdminID:    admin.ID,
		Action:     "broadcast_notification",
		TargetType: "platform",
		TargetID:   notificationID,
		Metadata: map[string]interface{}{
			"title":           req.Title,
			"type":            req.Type,
			"target_audience": req.TargetAudience,
		},
	})

	util.RespondWithSuccess(w, "Notification broadcast scheduled", map[string]interface{}{
		"notification_id": notificationID,
	})
}

// EnableMaintenanceMode enables platform maintenance mode
func (h *PlatformControlHandler) EnableMaintenanceMode(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Message      string `json:"message"`
		EstimatedEnd string `json:"estimated_end,omitempty"` // ISO 8601 timestamp
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	admin := r.Context().Value("user").(*repository.User)

	if err := h.adminRepo.EnableMaintenanceMode(r.Context(), req.Message, req.EstimatedEnd, admin.ID); err != nil {
		h.logger.Error("Failed to enable maintenance mode", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to enable maintenance mode")
		return
	}

	h.auditRepo.LogAction(r.Context(), &repository.AuditLog{
		AdminID:    admin.ID,
		Action:     "enable_maintenance_mode",
		TargetType: "platform",
		Severity:   "HIGH",
		Metadata: map[string]interface{}{
			"message":       req.Message,
			"estimated_end": req.EstimatedEnd,
		},
	})

	util.RespondWithSuccess(w, "Maintenance mode enabled", nil)
}

// DisableMaintenanceMode disables maintenance mode
func (h *PlatformControlHandler) DisableMaintenanceMode(w http.ResponseWriter, r *http.Request) {
	if err := h.adminRepo.DisableMaintenanceMode(r.Context()); err != nil {
		h.logger.Error("Failed to disable maintenance mode", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to disable maintenance mode")
		return
	}

	admin := r.Context().Value("user").(*repository.User)
	h.auditRepo.LogAction(r.Context(), &repository.AuditLog{
		AdminID:    admin.ID,
		Action:     "disable_maintenance_mode",
		TargetType: "platform",
	})

	util.RespondWithSuccess(w, "Maintenance mode disabled", nil)
}

// GetPlatformHealth retrieves current platform health metrics
func (h *PlatformControlHandler) GetPlatformHealth(w http.ResponseWriter, r *http.Request) {
	health, err := h.adminRepo.GetPlatformHealth(r.Context())
	if err != nil {
		h.logger.Error("Failed to get platform health", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to get health metrics")
		return
	}

	util.RespondWithSuccess(w, "", health)
}

// ActivateKillSwitch activates an emergency kill switch
func (h *PlatformControlHandler) ActivateKillSwitch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	switchName := vars["switch"]

	var req struct {
		Reason string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Reason == "" {
		util.RespondWithError(w, http.StatusBadRequest, "Reason required for kill switch")
		return
	}

	admin := r.Context().Value("user").(*repository.User)

	if err := h.adminRepo.ActivateKillSwitch(r.Context(), switchName, req.Reason, admin.ID); err != nil {
		h.logger.Error("Failed to activate kill switch", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to activate kill switch")
		return
	}

	// CRITICAL audit log
	h.auditRepo.LogAction(r.Context(), &repository.AuditLog{
		AdminID:    admin.ID,
		Action:     "activate_kill_switch",
		TargetType: "platform",
		TargetID:   switchName,
		Reason:     req.Reason,
		Severity:   "CRITICAL",
	})

	h.logger.Info("CRITICAL: Kill switch activated", map[string]interface{}{
		"switch": switchName,
		"reason": req.Reason,
		"admin":  admin.Username,
	})

	util.RespondWithSuccess(w, "Kill switch activated", map[string]interface{}{
		"switch": switchName,
		"active": true,
	})
}

// DeactivateKillSwitch deactivates a kill switch
func (h *PlatformControlHandler) DeactivateKillSwitch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	switchName := vars["switch"]

	if err := h.adminRepo.DeactivateKillSwitch(r.Context(), switchName); err != nil {
		h.logger.Error("Failed to deactivate kill switch", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to deactivate kill switch")
		return
	}

	admin := r.Context().Value("user").(*repository.User)
	h.auditRepo.LogAction(r.Context(), &repository.AuditLog{
		AdminID:    admin.ID,
		Action:     "deactivate_kill_switch",
		TargetType: "platform",
		TargetID:   switchName,
		Severity:   "HIGH",
	})

	util.RespondWithSuccess(w, "Kill switch deactivated", nil)
}
