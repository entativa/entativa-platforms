package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"messaging-service/internal/logger"
	"messaging-service/internal/repository"
	"messaging-service/internal/util"
)

type BackupHandler struct {
	backupRepo *repository.BackupRepository
	logger     *logger.Logger
}

func NewBackupHandler(backupRepo *repository.BackupRepository, logger *logger.Logger) *BackupHandler {
	return &BackupHandler{
		backupRepo: backupRepo,
		logger:     logger,
	}
}

// GetBackupSettings retrieves user's backup settings
func (h *BackupHandler) GetBackupSettings(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	settings, err := h.backupRepo.GetBackupSettings(r.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get backup settings", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to get settings")
		return
	}

	util.RespondWithSuccess(w, "", settings)
}

// UpdateBackupSettings updates user's backup preferences
func (h *BackupHandler) UpdateBackupSettings(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req struct {
		BackupEnabled        bool   `json:"backup_enabled"`
		BackupLocation       string `json:"backup_location"`        // our_servers, google_drive, icloud
		AutoBackupEnabled    bool   `json:"auto_backup_enabled"`
		AutoBackupFrequency  string `json:"auto_backup_frequency"`  // daily, weekly, monthly
		AutoBackupWifiOnly   bool   `json:"auto_backup_wifi_only"`
		KeepBackupsCount     int    `json:"keep_backups_count"`
		ThirdPartyAccountID  string `json:"third_party_account_id,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate backup location
	validLocations := map[string]bool{
		"our_servers":  true,
		"google_drive": true,
		"icloud":       true,
	}
	if !validLocations[req.BackupLocation] {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid backup location")
		return
	}

	// Validate frequency
	validFrequencies := map[string]bool{
		"daily":   true,
		"weekly":  true,
		"monthly": true,
	}
	if req.AutoBackupEnabled && !validFrequencies[req.AutoBackupFrequency] {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid backup frequency")
		return
	}

	if err := h.backupRepo.UpdateBackupSettings(r.Context(), userID, &req); err != nil {
		h.logger.Error("Failed to update backup settings", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to update settings")
		return
	}

	util.RespondWithSuccess(w, "Backup settings updated", nil)
}

// AcknowledgeThirdPartyWarning records user acknowledgment of third-party backup risks
func (h *BackupHandler) AcknowledgeThirdPartyWarning(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	if err := h.backupRepo.AcknowledgeThirdPartyWarning(r.Context(), userID); err != nil {
		h.logger.Error("Failed to acknowledge warning", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to acknowledge")
		return
	}

	util.RespondWithSuccess(w, "Warning acknowledged", nil)
}

// SetupBackupKey creates backup encryption key from user's PIN/passphrase
func (h *BackupHandler) SetupBackupKey(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req struct {
		PIN        string `json:"pin,omitempty"`         // 6-8 digit PIN
		Passphrase string `json:"passphrase,omitempty"`  // or passphrase (min 12 chars)
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate PIN or passphrase
	if req.PIN == "" && req.Passphrase == "" {
		util.RespondWithError(w, http.StatusBadRequest, "PIN or passphrase required")
		return
	}

	if req.PIN != "" && (len(req.PIN) < 6 || len(req.PIN) > 8) {
		util.RespondWithError(w, http.StatusBadRequest, "PIN must be 6-8 digits")
		return
	}

	if req.Passphrase != "" && len(req.Passphrase) < 12 {
		util.RespondWithError(w, http.StatusBadRequest, "Passphrase must be at least 12 characters")
		return
	}

	// Use passphrase if provided, otherwise PIN
	secret := req.Passphrase
	if secret == "" {
		secret = req.PIN
	}

	// Create backup key
	if err := h.backupRepo.CreateBackupKey(r.Context(), userID, secret); err != nil {
		h.logger.Error("Failed to create backup key", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to create backup key")
		return
	}

	util.RespondWithSuccess(w, "Backup key created successfully", nil)
}

// CreateBackup creates an encrypted backup of user's messages
func (h *BackupHandler) CreateBackup(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req struct {
		BackupType string `json:"backup_type"` // 'full' or 'incremental'
		PIN        string `json:"pin,omitempty"`
		Passphrase string `json:"passphrase,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate backup type
	if req.BackupType != "full" && req.BackupType != "incremental" {
		req.BackupType = "full"
	}

	// Use passphrase if provided, otherwise PIN
	secret := req.Passphrase
	if secret == "" {
		secret = req.PIN
	}

	if secret == "" {
		util.RespondWithError(w, http.StatusBadRequest, "PIN or passphrase required")
		return
	}

	// Create backup
	backupID, err := h.backupRepo.CreateBackup(r.Context(), userID, secret, req.BackupType)
	if err != nil {
		h.logger.Error("Failed to create backup", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to create backup")
		return
	}

	util.RespondWithSuccess(w, "Backup created successfully", map[string]interface{}{
		"backup_id": backupID,
	})
}

// GetBackupHistory retrieves user's backup history
func (h *BackupHandler) GetBackupHistory(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	backups, err := h.backupRepo.GetBackupHistory(r.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get backup history", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to get history")
		return
	}

	util.RespondWithSuccess(w, "", map[string]interface{}{
		"backups": backups,
	})
}

// RestoreBackup restores messages from a backup
func (h *BackupHandler) RestoreBackup(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	vars := mux.Vars(r)
	backupID := vars["id"]

	var req struct {
		PIN        string `json:"pin,omitempty"`
		Passphrase string `json:"passphrase,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Use passphrase if provided, otherwise PIN
	secret := req.Passphrase
	if secret == "" {
		secret = req.PIN
	}

	if secret == "" {
		util.RespondWithError(w, http.StatusBadRequest, "PIN or passphrase required")
		return
	}

	// Restore backup
	restoredCount, err := h.backupRepo.RestoreBackup(r.Context(), userID, backupID, secret)
	if err != nil {
		h.logger.Error("Failed to restore backup", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to restore backup")
		return
	}

	util.RespondWithSuccess(w, "Backup restored successfully", map[string]interface{}{
		"messages_restored": restoredCount,
	})
}

// DeleteBackup deletes a specific backup
func (h *BackupHandler) DeleteBackup(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	vars := mux.Vars(r)
	backupID := vars["id"]

	if err := h.backupRepo.DeleteBackup(r.Context(), userID, backupID); err != nil {
		h.logger.Error("Failed to delete backup", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to delete backup")
		return
	}

	util.RespondWithSuccess(w, "Backup deleted successfully", nil)
}

// GetBackupActivityLog retrieves backup activity history
func (h *BackupHandler) GetBackupActivityLog(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	logs, err := h.backupRepo.GetBackupActivityLog(r.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get activity log", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to get activity log")
		return
	}

	util.RespondWithSuccess(w, "", map[string]interface{}{
		"activity_log": logs,
	})
}
