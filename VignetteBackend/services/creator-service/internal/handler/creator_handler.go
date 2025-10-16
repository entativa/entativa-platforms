package handler

import (
	"net/http"
	"strconv"

	"github.com/entativa/vignette/creator-service/internal/model"
	"github.com/entativa/vignette/creator-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreatorHandler struct {
	service *service.CreatorService
}

func NewCreatorHandler(service *service.CreatorService) *CreatorHandler {
	return &CreatorHandler{service: service}
}

// CreateProfile creates a creator profile
func (h *CreatorHandler) CreateProfile(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))

	var req model.CreateCreatorProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := h.service.CreateCreatorProfile(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, profile)
}

// GetProfile gets creator profile
func (h *CreatorHandler) GetProfile(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))

	profile, err := h.service.GetCreatorProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateProfile updates creator profile
func (h *CreatorHandler) UpdateProfile(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))

	var req model.UpdateCreatorProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.UpdateCreatorProfile(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated"})
}

// GetAnalyticsOverview gets analytics overview
func (h *CreatorHandler) GetAnalyticsOverview(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	period := c.DefaultQuery("period", "30d")

	overview, err := h.service.GetAnalyticsOverview(c.Request.Context(), userID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, overview)
}

// GetAudienceInsights gets audience insights
func (h *CreatorHandler) GetAudienceInsights(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	insights, err := h.service.GetAudienceInsights(c.Request.Context(), userID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, insights)
}

// GetTopContent gets top performing content
func (h *CreatorHandler) GetTopContent(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	contentType := c.DefaultQuery("type", "post")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	topContent, err := h.service.GetTopContent(c.Request.Context(), userID, contentType, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"content": topContent, "count": len(topContent)})
}

// ApplyForMonetization applies for monetization
func (h *CreatorHandler) ApplyForMonetization(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))

	var req model.MonetizationApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.ApplyForMonetization(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "application submitted", "status": "pending"})
}

// AwardBadge awards a badge (admin endpoint)
func (h *CreatorHandler) AwardBadge(c *gin.Context) {
	targetUserID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	badge := model.CreatorBadge(c.Param("badge"))

	err = h.service.AwardBadge(c.Request.Context(), targetUserID, badge)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "badge awarded", "badge": badge})
}
