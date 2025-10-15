package handler

import (
	"net/http"

	"vignette/user-service/internal/middleware"
	"vignette/user-service/internal/model"
	"vignette/user-service/internal/service"

	"github.com/gin-gonic/gin"
)

type TwoFactorHandler struct {
	twoFactorService *service.TwoFactorService
}

func NewTwoFactorHandler(twoFactorService *service.TwoFactorService) *TwoFactorHandler {
	return &TwoFactorHandler{twoFactorService: twoFactorService}
}

// Setup2FA initializes 2FA for the user
// @Summary Setup Two-Factor Authentication
// @Description Generate QR code and backup codes for 2FA
// @Tags 2fa
// @Security BearerAuth
// @Produce json
// @Success 200 {object} model.TwoFactorSetupResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/2fa/setup [post]
func (h *TwoFactorHandler) Setup2FA(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User not authenticated",
		})
		return
	}

	response, err := h.twoFactorService.SetupTwoFactor(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to setup 2FA",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "2FA setup initiated. Scan the QR code with your authenticator app",
		"data":    response,
	})
}

// Enable2FA enables 2FA after verification
// @Summary Enable Two-Factor Authentication
// @Description Enable 2FA by verifying TOTP code
// @Tags 2fa
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.Enable2FARequest true "Verification code"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/2fa/enable [post]
func (h *TwoFactorHandler) Enable2FA(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var req model.Enable2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	err := h.twoFactorService.EnableTwoFactor(userID, req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to enable 2FA",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Two-factor authentication enabled successfully",
	})
}

// Verify2FA verifies a 2FA code
// @Summary Verify 2FA Code
// @Description Verify TOTP code or backup code
// @Tags 2fa
// @Accept json
// @Produce json
// @Param request body model.Verify2FARequest true "2FA code"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/2fa/verify [post]
func (h *TwoFactorHandler) Verify2FA(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var req model.Verify2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	valid, err := h.twoFactorService.VerifyTwoFactorCode(userID, req.Code)
	if err != nil || !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid 2FA code",
			"message": "The code you entered is incorrect",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "2FA verification successful",
	})
}

// Disable2FA disables 2FA for the user
// @Summary Disable Two-Factor Authentication
// @Description Disable 2FA for the authenticated user
// @Tags 2fa
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/2fa/disable [post]
func (h *TwoFactorHandler) Disable2FA(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	err := h.twoFactorService.DisableTwoFactor(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to disable 2FA",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Two-factor authentication disabled",
	})
}
