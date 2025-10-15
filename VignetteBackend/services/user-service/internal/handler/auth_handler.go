package handler

import (
	"errors"
	"net/http"

	"vignette/user-service/internal/model"
	"vignette/user-service/internal/repository"
	"vignette/user-service/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Signup handles user registration
// @Summary Register a new user
// @Description Create a new Vignette account with instant access (Meta-level authentication)
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.SignupRequest true "Signup request"
// @Success 201 {object} model.AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/signup [post]
func (h *AuthHandler) Signup(c *gin.Context) {
	var req model.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// Get client IP and User-Agent
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// Create user
	authResp, err := h.authService.Signup(&req, ipAddress, userAgent)
	if err != nil {
		if errors.Is(err, repository.ErrEmailExists) {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Email already exists",
				"message": "An account with this email already exists",
			})
			return
		}
		if errors.Is(err, repository.ErrUsernameExists) {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Username already taken",
				"message": "This username is already taken. Please choose another one",
			})
			return
		}
		if errors.Is(err, service.ErrInvalidUsername) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid username",
				"message": "Username must be 3-30 characters and can only contain letters, numbers, periods, and underscores",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create account",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Account created successfully! Welcome to Vignette!",
		"data":    authResp,
	})
}

// Login handles user authentication
// @Summary Login to Vignette
// @Description Authenticate user with username/email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Login request"
// @Success 200 {object} model.AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// Get client IP and User-Agent
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// Authenticate user
	authResp, err := h.authService.Login(&req, ipAddress, userAgent)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid credentials",
				"message": "The username/email or password you entered is incorrect",
			})
			return
		}
		if errors.Is(err, service.ErrUserNotActive) {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Account not active",
				"message": "Your account has been deactivated",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Login failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful! Welcome back!",
		"data":    authResp,
	})
}

// Logout handles user logout
// @Summary Logout from Vignette
// @Description Invalidate user session
// @Tags auth
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User not authenticated",
		})
		return
	}

	// Logout user
	err := h.authService.Logout(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Logout failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logged out successfully",
	})
}

// Me returns the current authenticated user
// @Summary Get current user
// @Description Get authenticated user information
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} model.UserResponse
// @Failure 401 {object} map[string]interface{}
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	// User data is set by auth middleware
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User not authenticated",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

// HealthCheck returns service health status
func (h *AuthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "vignette-user-service",
		"version": "1.0.0",
	})
}
