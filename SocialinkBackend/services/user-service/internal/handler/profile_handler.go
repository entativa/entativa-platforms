package handler

import (
	"net/http"

	"socialink/user-service/internal/model"
	"socialink/user-service/internal/service"

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

// GetProfile retrieves user profile with all information
// @Summary Get user profile
// @Description Get complete profile information for a user
// @Tags profile
// @Security BearerAuth
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} model.ProfileResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /profile/{user_id} [get]
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

// UpdateProfileInfo updates basic profile information
// @Summary Update profile information
// @Description Update basic profile details
// @Tags profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.UpdateProfileInfoRequest true "Profile update request"
// @Success 200 {object} model.ProfileResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /profile/info [put]
func (h *ProfileHandler) UpdateProfileInfo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User not authenticated",
		})
		return
	}

	var req model.UpdateProfileInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	userUUID, _ := uuid.Parse(userID.(string))
	profile, err := h.profileService.UpdateProfileInfo(c.Request.Context(), userUUID, &req)
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

// AddWorkExperience adds work experience to profile
// @Summary Add work experience
// @Description Add work experience entry to profile
// @Tags profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.AddWorkExperienceRequest true "Work experience request"
// @Success 200 {object} model.ProfileResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /profile/work [post]
func (h *ProfileHandler) AddWorkExperience(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User not authenticated",
		})
		return
	}

	var req model.AddWorkExperienceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	userUUID, _ := uuid.Parse(userID.(string))
	profile, err := h.profileService.AddWorkExperience(c.Request.Context(), userUUID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to add work experience",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Work experience added successfully",
		"data":    profile.ToProfileResponse(),
	})
}

// RemoveWorkExperience removes work experience from profile
// @Summary Remove work experience
// @Description Remove work experience entry from profile
// @Tags profile
// @Security BearerAuth
// @Param work_id path string true "Work ID"
// @Success 200 {object} model.ProfileResponse
// @Failure 401 {object} map[string]interface{}
// @Router /profile/work/{work_id} [delete]
func (h *ProfileHandler) RemoveWorkExperience(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User not authenticated",
		})
		return
	}

	workID := c.Param("work_id")
	userUUID, _ := uuid.Parse(userID.(string))

	profile, err := h.profileService.RemoveWorkExperience(c.Request.Context(), userUUID, workID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to remove work experience",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Work experience removed successfully",
		"data":    profile.ToProfileResponse(),
	})
}

// AddEducation adds education to profile
// @Summary Add education
// @Description Add education entry to profile
// @Tags profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.AddEducationRequest true "Education request"
// @Success 200 {object} model.ProfileResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /profile/education [post]
func (h *ProfileHandler) AddEducation(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User not authenticated",
		})
		return
	}

	var req model.AddEducationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	userUUID, _ := uuid.Parse(userID.(string))
	profile, err := h.profileService.AddEducation(c.Request.Context(), userUUID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to add education",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Education added successfully",
		"data":    profile.ToProfileResponse(),
	})
}

// RemoveEducation removes education from profile
// @Summary Remove education
// @Description Remove education entry from profile
// @Tags profile
// @Security BearerAuth
// @Param education_id path string true "Education ID"
// @Success 200 {object} model.ProfileResponse
// @Failure 401 {object} map[string]interface{}
// @Router /profile/education/{education_id} [delete]
func (h *ProfileHandler) RemoveEducation(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User not authenticated",
		})
		return
	}

	educationID := c.Param("education_id")
	userUUID, _ := uuid.Parse(userID.(string))

	profile, err := h.profileService.RemoveEducation(c.Request.Context(), userUUID, educationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to remove education",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Education removed successfully",
		"data":    profile.ToProfileResponse(),
	})
}

// UpdateContactInfo updates contact information
// @Summary Update contact information
// @Description Update user contact details
// @Tags profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.UpdateContactInfoRequest true "Contact info request"
// @Success 200 {object} model.ProfileResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /profile/contact [put]
func (h *ProfileHandler) UpdateContactInfo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User not authenticated",
		})
		return
	}

	var req model.UpdateContactInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	userUUID, _ := uuid.Parse(userID.(string))
	profile, err := h.profileService.UpdateContactInfo(c.Request.Context(), userUUID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update contact info",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Contact information updated successfully",
		"data":    profile.ToProfileResponse(),
	})
}

// UpdateSocialLinks updates social media links
// @Summary Update social links
// @Description Update social media profile links
// @Tags profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.UpdateSocialLinksRequest true "Social links request"
// @Success 200 {object} model.ProfileResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /profile/social-links [put]
func (h *ProfileHandler) UpdateSocialLinks(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User not authenticated",
		})
		return
	}

	var req model.UpdateSocialLinksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	userUUID, _ := uuid.Parse(userID.(string))
	profile, err := h.profileService.UpdateSocialLinks(c.Request.Context(), userUUID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update social links",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Social links updated successfully",
		"data":    profile.ToProfileResponse(),
	})
}

// UpdateVisibility updates profile visibility settings
// @Summary Update profile visibility
// @Description Update who can see various parts of your profile
// @Tags profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.UpdateVisibilityRequest true "Visibility settings"
// @Success 200 {object} model.ProfileResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /profile/visibility [put]
func (h *ProfileHandler) UpdateVisibility(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User not authenticated",
		})
		return
	}

	var req model.UpdateVisibilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	userUUID, _ := uuid.Parse(userID.(string))
	profile, err := h.profileService.UpdateVisibility(c.Request.Context(), userUUID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update visibility",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Visibility settings updated successfully",
		"data":    profile.ToProfileResponse(),
	})
}
