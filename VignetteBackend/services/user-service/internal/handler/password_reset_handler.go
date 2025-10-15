package handler

import (
	"net/http"

	"vignette/user-service/internal/model"
	"vignette/user-service/internal/service"

	"github.com/gin-gonic/gin"
)

type PasswordResetHandler struct {
	passwordResetService *service.PasswordResetService
}

func NewPasswordResetHandler(passwordResetService *service.PasswordResetService) *PasswordResetHandler {
	return &PasswordResetHandler{passwordResetService: passwordResetService}
}

// RequestPasswordReset initiates password reset process
// @Summary Request Password Reset
// @Description Send password reset link to user's email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.RequestPasswordResetRequest true "Email"
// @Success 200 {object} map[string]interface{}
// @Router /auth/password-reset/request [post]
func (h *PasswordResetHandler) RequestPasswordReset(c *gin.Context) {
	var req model.RequestPasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	err := h.passwordResetService.RequestPasswordReset(req.Email)
	if err != nil {
		// Don't reveal if email exists
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "If an account exists with this email, you will receive password reset instructions",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "If an account exists with this email, you will receive password reset instructions",
	})
}

// ResetPassword resets user password with token
// @Summary Reset Password
// @Description Reset password using reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.ResetPasswordRequest true "Reset token and new password"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/password-reset/confirm [post]
func (h *PasswordResetHandler) ResetPassword(c *gin.Context) {
	var req model.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	err := h.passwordResetService.ResetPassword(req.Token, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to reset password",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password reset successfully. You can now log in with your new password",
	})
}
