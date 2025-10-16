package handler

import (
	"net/http"

	"github.com/entativa/vignette/settings-service/internal/model"
	"github.com/entativa/vignette/settings-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SettingsHandler struct {
	service *service.SettingsService
}

func NewSettingsHandler(service *service.SettingsService) *SettingsHandler {
	return &SettingsHandler{service: service}
}

// GetSettings gets user settings
func (h *SettingsHandler) GetSettings(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))

	settings, err := h.service.GetOrCreateSettings(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// UpdateSettings updates user settings
func (h *SettingsHandler) UpdateSettings(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))

	var req model.UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.UpdateSettings(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "settings updated"})
}

// CreateKeyBackup creates encrypted key backup
func (h *SettingsHandler) CreateKeyBackup(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))

	var req model.CreateKeyBackupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.CreateKeyBackup(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "key backup created successfully"})
}

// RestoreKeyBackup restores encrypted key backup
func (h *SettingsHandler) RestoreKeyBackup(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))

	var req model.RestoreKeyBackupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.RestoreKeyBackup(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetKeyBackupInfo gets key backup info
func (h *SettingsHandler) GetKeyBackupInfo(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))

	info, err := h.service.GetKeyBackupInfo(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, info)
}

// DeleteKeyBackup deletes key backup
func (h *SettingsHandler) DeleteKeyBackup(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))

	err := h.service.DeleteKeyBackup(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "key backup deleted"})
}

// GetStorageLocations gets info about storage location options
func (h *SettingsHandler) GetStorageLocations(c *gin.Context) {
	locations := h.service.GetStorageLocationInfo()
	c.JSON(http.StatusOK, gin.H{"locations": locations})
}
