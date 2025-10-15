package handler

import (
	"net/http"

	"vignette/meta-user-service/internal/model"
	"vignette/meta-user-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MetaUserHandler struct {
	userService *service.MetaUserService
	syncService *service.CrossPlatformSyncService
}

func NewMetaUserHandler(userService *service.MetaUserService, syncService *service.CrossPlatformSyncService) *MetaUserHandler {
	return &MetaUserHandler{
		userService: userService,
		syncService: syncService,
	}
}

// CreateUser handles meta user creation
func (h *MetaUserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract request metadata
	createReq := &service.CreateMetaUserRequest{
		Email:       req.Email,
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
		IPAddress:   c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
		DeviceID:    c.GetHeader("X-Device-ID"),
		DeviceType:  c.GetHeader("X-Device-Type"),
		OS:          c.GetHeader("X-OS"),
		Region:      c.GetHeader("X-Region"),
	}

	user, err := h.userService.CreateMetaUser(c.Request.Context(), createReq)
	if err != nil {
		if err == service.ErrHighRiskDetected {
			c.JSON(http.StatusForbidden, gin.H{"error": "Account creation blocked due to security concerns"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user_id":  user.ID,
		"meta_id":  user.MetaID,
		"email":    user.Email,
		"status":   user.Status,
	})
}

// Authenticate handles user authentication
func (h *MetaUserHandler) Authenticate(c *gin.Context) {
	var req AuthenticationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authReq := &service.AuthenticationRequest{
		Email:      req.Email,
		Password:   req.Password,
		IPAddress:  c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		DeviceID:   c.GetHeader("X-Device-ID"),
		DeviceType: c.GetHeader("X-Device-Type"),
	}

	resp, err := h.userService.Authenticate(c.Request.Context(), authReq)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		if err == service.ErrUserSuspended {
			c.JSON(http.StatusForbidden, gin.H{"error": "Account suspended"})
			return
		}
		if err == service.ErrUserBanned {
			c.JSON(http.StatusForbidden, gin.H{"error": "Account banned"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication failed"})
		return
	}

	if resp.RequiresTwoFactor {
		c.JSON(http.StatusAccepted, gin.H{
			"requires_2fa": true,
			"user_id":      resp.UserID,
		})
		return
	}

	if resp.AnomalyDetected {
		c.JSON(http.StatusAccepted, gin.H{
			"requires_additional_auth": true,
			"anomaly_detected":         true,
			"anomaly_score":            resp.AnomalyScore,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"session_token": resp.SessionToken,
		"user": gin.H{
			"id":      resp.User.ID,
			"meta_id": resp.User.MetaID,
			"email":   resp.User.Email,
		},
	})
}

// GetUser retrieves user information
func (h *MetaUserHandler) GetUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// LinkPlatform links a platform account
func (h *MetaUserHandler) LinkPlatform(c *gin.Context) {
	var req LinkPlatformRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userService.LinkPlatformAccount(c.Request.Context(), req.MetaUserID, req.Platform, req.PlatformUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to link platform account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// UpdateTrustScore manually triggers trust score update
func (h *MetaUserHandler) UpdateTrustScore(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.userService.UpdateTrustScore(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update trust score"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// EnableSync enables cross-platform synchronization
func (h *MetaUserHandler) EnableSync(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.syncService.EnableCrossPlatformSync(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enable sync"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DisableSync disables cross-platform synchronization
func (h *MetaUserHandler) DisableSync(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.syncService.DisableCrossPlatformSync(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disable sync"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// UpdatePrivacySettings updates user privacy settings
func (h *MetaUserHandler) UpdatePrivacySettings(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var settings model.PrivacySettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.syncService.SyncPrivacySettings(c.Request.Context(), userID, settings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update privacy settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// Request/Response types

type CreateUserRequest struct {
	Email       string  `json:"email" binding:"required,email"`
	Password    string  `json:"password" binding:"required,min=8"`
	PhoneNumber *string `json:"phone_number,omitempty"`
}

type AuthenticationRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LinkPlatformRequest struct {
	MetaUserID     uuid.UUID `json:"meta_user_id" binding:"required"`
	Platform       string    `json:"platform" binding:"required"`
	PlatformUserID uuid.UUID `json:"platform_user_id" binding:"required"`
}
