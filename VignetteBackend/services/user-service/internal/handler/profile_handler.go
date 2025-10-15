package handler

import (
	"net/http"

	"vignette/user-service/internal/model"
	"vignette/user-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProfileHandler struct {
	profileService *service.ProfileService
}

func NewProfileHandler(profileService *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
	}
}

// GetProfile retrieves user profile
// @Summary Get user profile
// @Description Get complete profile information for a user
// @Tags profile
// @Security BearerAuth
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} model.ProfileResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /profile/{username} [get]
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"message": "The provided user ID is not valid",
		})
		return
	}

	// Increment profile views
	viewerID, exists := c.Get("user_id")
	if exists {
		viewerUUID, _ := uuid.Parse(viewerID.(string))
		// Only increment if viewer is not the profile owner
		if viewerUUID != userID {
			h.profileService.IncrementProfileViews(c.Request.Context(), userID)
		}
	}

	profile, err := h.profileService.GetProfileWithUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get profile",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    profile,
	})
}

// GetMyProfile retrieves the authenticated user's profile
// @Summary Get my profile
// @Description Get complete profile information for authenticated user
// @Tags profile
// @Security BearerAuth
// @Produce json
// @Success 200 {object} model.ProfileResponse
// @Failure 401 {object} map[string]interface{}
// @Router /profile/me [get]
func (h *ProfileHandler) GetMyProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User not authenticated",
		})
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"message": "Invalid user ID format",
		})
		return
	}

	profile, err := h.profileService.GetProfileWithUser(c.Request.Context(), userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get profile",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    profile,
	})
}

// UpdateProfileExtended updates extended profile information
// @Summary Update extended profile
// @Description Update category, gender, pronouns, etc.
// @Tags profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.UpdateProfileExtendedRequest true "Profile update request"
// @Success 200 {object} model.ProfileResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /profile/extended [put]
func (h *ProfileHandler) UpdateProfileExtended(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User not authenticated",
		})
		return
	}

	var req model.UpdateProfileExtendedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	userUUID, _ := uuid.Parse(userID.(string))
	profile, err := h.profileService.UpdateProfileExtended(c.Request.Context(), userUUID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update profile",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile updated successfully",
		"data":    profile.ToProfileResponse(),
	})
}

// AddLinkInBio adds a link to bio
// @Summary Add link in bio
// @Description Add a clickable link to profile bio
// @Tags profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.AddLinkInBioRequest true "Link request"
// @Success 200 {object} model.ProfileResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /profile/links [post]
func (h *ProfileHandler) AddLinkInBio(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User not authenticated",
		})
		return
	}

	var req model.AddLinkInBioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	userUUID, _ := uuid.Parse(userID.(string))
	profile, err := h.profileService.AddLinkInBio(c.Request.Context(), userUUID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to add link",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Link added successfully",
		"data":    profile.ToProfileResponse(),
	})
}

// RemoveLinkInBio removes a link from bio
// @Summary Remove link from bio
// @Description Remove a clickable link from profile bio
// @Tags profile
// @Security BearerAuth
// @Param link_id path string true "Link ID"
// @Success 200 {object} model.ProfileResponse
// @Failure 401 {object} map[string]interface{}
// @Router /profile/links/{link_id} [delete]
func (h *ProfileHandler) RemoveLinkInBio(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User not authenticated",
		})
		return
	}

	linkID := c.Param("link_id")
	userUUID, _ := uuid.Parse(userID.(string))

	profile, err := h.profileService.RemoveLinkInBio(c.Request.Context(), userUUID, linkID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to remove link",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Link removed successfully",
		"data":    profile.ToProfileResponse(),
	})
}

// AddHighlight adds a story highlight
// @Summary Add story highlight
// @Description Add a story highlight to profile
// @Tags profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.AddHighlightRequest true "Highlight request"
// @Success 200 {object} model.ProfileResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /profile/highlights [post]
func (h *ProfileHandler) AddHighlight(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User not authenticated",
		})
		return
	}

	var req model.AddHighlightRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	userUUID, _ := uuid.Parse(userID.(string))
	profile, err := h.profileService.AddHighlight(c.Request.Context(), userUUID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to add highlight",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Highlight added successfully",
		"data":    profile.ToProfileResponse(),
	})
}

// RemoveHighlight removes a story highlight
// @Summary Remove story highlight
// @Description Remove a story highlight from profile
// @Tags profile
// @Security BearerAuth
// @Param highlight_id path string true "Highlight ID"
// @Success 200 {object} model.ProfileResponse
// @Failure 401 {object} map[string]interface{}
// @Router /profile/highlights/{highlight_id} [delete]
func (h *ProfileHandler) RemoveHighlight(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User not authenticated",
		})
		return
	}

	highlightID := c.Param("highlight_id")
	userUUID, _ := uuid.Parse(userID.(string))

	profile, err := h.profileService.RemoveHighlight(c.Request.Context(), userUUID, highlightID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to remove highlight",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Highlight removed successfully",
		"data":    profile.ToProfileResponse(),
	})
}

// UpdateContactOptions updates contact options for business account
// @Summary Update contact options
// @Description Update business/creator contact options
// @Tags profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.UpdateContactOptionsRequest true "Contact options"
// @Success 200 {object} model.ProfileResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /profile/contact-options [put]
func (h *ProfileHandler) UpdateContactOptions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User not authenticated",
		})
		return
	}

	var req model.UpdateContactOptionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	userUUID, _ := uuid.Parse(userID.(string))
	profile, err := h.profileService.UpdateContactOptions(c.Request.Context(), userUUID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update contact options",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Contact options updated successfully",
		"data":    profile.ToProfileResponse(),
	})
}

// EnableCreatorAccount switches to creator account
// @Summary Switch to creator account
// @Description Enable creator account features
// @Tags profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.EnableCreatorAccountRequest true "Creator account settings"
// @Success 200 {object} model.ProfileResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /profile/creator/enable [post]
func (h *ProfileHandler) EnableCreatorAccount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User not authenticated",
		})
		return
	}

	var req model.EnableCreatorAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	userUUID, _ := uuid.Parse(userID.(string))
	profile, err := h.profileService.EnableCreatorAccount(c.Request.Context(), userUUID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to enable creator account",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Creator account enabled successfully! You now have access to insights and analytics",
		"data":    profile.ToProfileResponse(),
	})
}

// EnableBusinessAccount switches to business account
// @Summary Switch to business account
// @Description Enable business account features
// @Tags profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.EnableBusinessAccountRequest true "Business account settings"
// @Success 200 {object} model.ProfileResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /profile/business/enable [post]
func (h *ProfileHandler) EnableBusinessAccount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User not authenticated",
		})
		return
	}

	var req model.EnableBusinessAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	userUUID, _ := uuid.Parse(userID.(string))
	profile, err := h.profileService.EnableBusinessAccount(c.Request.Context(), userUUID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to enable business account",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Business account enabled successfully! You can now add contact buttons and business hours",
		"data":    profile.ToProfileResponse(),
	})
}

// UpdateAvailability updates availability status
// @Summary Update availability
// @Description Update availability status for professional account
// @Tags profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.UpdateAvailabilityRequest true "Availability request"
// @Success 200 {object} model.ProfileResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /profile/availability [put]
func (h *ProfileHandler) UpdateAvailability(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User not authenticated",
		})
		return
	}

	var req model.UpdateAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	userUUID, _ := uuid.Parse(userID.(string))
	profile, err := h.profileService.UpdateAvailability(c.Request.Context(), userUUID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update availability",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Availability updated successfully",
		"data":    profile.ToProfileResponse(),
	})
}
