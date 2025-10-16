package handler

import (
	"fmt"
	"net/http"

	"vignette/post-service/internal/model"
	"vignette/post-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TakesHandler struct {
	takesService *service.TakesService
}

func NewTakesHandler(takesService *service.TakesService) *TakesHandler {
	return &TakesHandler{
		takesService: takesService,
	}
}

// CreateTake creates a new Take
// @Summary Create Take
// @Description Create a new Take (short video)
// @Tags takes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param take body model.CreateTakeRequest true "Take data"
// @Success 201 {object} model.TakeResponse
// @Router /takes [post]
func (h *TakesHandler) CreateTake(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userUUID, _ := uuid.Parse(userID.(string))

	var req model.CreateTakeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	take, err := h.takesService.CreateTake(c.Request.Context(), userUUID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": take})
}

// GetTake retrieves a Take by ID
// @Summary Get Take
// @Description Get a Take by ID
// @Tags takes
// @Produce json
// @Param take_id path string true "Take ID"
// @Success 200 {object} model.TakeResponse
// @Router /takes/{take_id} [get]
func (h *TakesHandler) GetTake(c *gin.Context) {
	takeIDStr := c.Param("take_id")
	takeID, err := uuid.Parse(takeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Take ID"})
		return
	}

	take, err := h.takesService.GetTake(c.Request.Context(), takeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": take})
}

// GetTrendingTakes retrieves trending Takes
// @Summary Get trending Takes
// @Description Get trending Takes
// @Tags takes
// @Produce json
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} []model.TakeResponse
// @Router /takes/trending [get]
func (h *TakesHandler) GetTrendingTakes(c *gin.Context) {
	limit := 20
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	takes, err := h.takesService.GetTrendingTakes(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": takes, "count": len(takes)})
}

// CreateBTT creates Behind-the-Takes content
// @Summary Create Behind-the-Takes
// @Description Add BTT content to a Take
// @Tags takes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param take_id path string true "Take ID"
// @Param btt body model.CreateBTTRequest true "BTT data"
// @Success 201 {object} model.BTTResponse
// @Router /takes/{take_id}/btt [post]
func (h *TakesHandler) CreateBTT(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userUUID, _ := uuid.Parse(userID.(string))
	takeIDStr := c.Param("take_id")
	takeID, err := uuid.Parse(takeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Take ID"})
		return
	}

	var req model.CreateBTTRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	btt, err := h.takesService.CreateBTT(c.Request.Context(), takeID, userUUID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": btt})
}

// GetBTT retrieves Behind-the-Takes content
// @Summary Get Behind-the-Takes
// @Description Get BTT content for a Take
// @Tags takes
// @Produce json
// @Param take_id path string true "Take ID"
// @Success 200 {object} model.BTTResponse
// @Router /takes/{take_id}/btt [get]
func (h *TakesHandler) GetBTT(c *gin.Context) {
	takeIDStr := c.Param("take_id")
	takeID, err := uuid.Parse(takeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Take ID"})
		return
	}

	btt, err := h.takesService.GetBTT(c.Request.Context(), takeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": btt})
}

// GetTrendingBTT retrieves trending BTT content
// @Summary Get trending BTT
// @Description Get trending Behind-the-Takes content
// @Tags takes
// @Produce json
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} []model.BTTResponse
// @Router /takes/btt/trending [get]
func (h *TakesHandler) GetTrendingBTT(c *gin.Context) {
	limit := 20
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	bttList, err := h.takesService.GetTrendingBTT(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": bttList, "count": len(bttList)})
}

// CreateTemplate creates a template from a Take
// @Summary Create Template
// @Description Create a reusable template from a Take
// @Tags templates
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param take_id path string true "Take ID"
// @Param template body model.CreateTemplateRequest true "Template data"
// @Success 201 {object} model.TemplateResponse
// @Router /takes/{take_id}/template [post]
func (h *TakesHandler) CreateTemplate(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userUUID, _ := uuid.Parse(userID.(string))
	takeIDStr := c.Param("take_id")
	takeID, err := uuid.Parse(takeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Take ID"})
		return
	}

	var req model.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template, err := h.takesService.CreateTemplate(c.Request.Context(), takeID, userUUID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": template})
}

// GetTemplates retrieves available templates
// @Summary Get templates
// @Description Get Takes templates
// @Tags templates
// @Produce json
// @Param category query string false "Category filter"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} []model.TemplateResponse
// @Router /templates [get]
func (h *TakesHandler) GetTemplates(c *gin.Context) {
	var category *string
	if cat := c.Query("category"); cat != "" {
		category = &cat
	}

	limit := 20
	offset := 0
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if o := c.Query("offset"); o != "" {
		fmt.Sscanf(o, "%d", &offset)
	}

	templates, err := h.takesService.GetTemplates(c.Request.Context(), category, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": templates, "count": len(templates)})
}

// GetActiveTrends retrieves active trends
// @Summary Get active trends
// @Description Get currently active Takes trends
// @Tags trends
// @Produce json
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} []model.TrendResponse
// @Router /trends [get]
func (h *TakesHandler) GetActiveTrends(c *gin.Context) {
	limit := 20
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	trends, err := h.takesService.GetActiveTrends(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": trends, "count": len(trends)})
}

// GetTrendTakes retrieves Takes for a trend (with originator deep-link)
// @Summary Get trend Takes
// @Description Get all Takes for a specific trend (deep-linked to originator)
// @Tags trends
// @Produce json
// @Param trend_id path string true "Trend ID"
// @Param limit query int false "Limit" default(30)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} map[string]interface{}
// @Router /trends/{trend_id}/takes [get]
func (h *TakesHandler) GetTrendTakes(c *gin.Context) {
	trendIDStr := c.Param("trend_id")
	trendID, err := uuid.Parse(trendIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Trend ID"})
		return
	}

	limit := 30
	offset := 0
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if o := c.Query("offset"); o != "" {
		fmt.Sscanf(o, "%d", &offset)
	}

	takes, trend, err := h.takesService.GetTrendTakes(c.Request.Context(), trendID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"trend": gin.H{
			"id":                trend.ID,
			"keyword":           trend.Keyword,
			"display_name":      trend.DisplayName,
			"originator_id":     trend.OriginatorID, // Deep-link to originator
			"origin_take_id":    trend.OriginTakeID, // Deep-link to original Take
			"participant_count": trend.ParticipantCount,
		},
		"takes": takes,
		"count": len(takes),
	})
}

// JoinTrend allows user to join a trend with their Take
// @Summary Join trend
// @Description Join an existing trend or create new one
// @Tags trends
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param take_id path string true "Take ID"
// @Param request body model.JoinTrendRequest true "Trend keyword"
// @Success 200 {object} model.TrendResponse
// @Router /takes/{take_id}/trend [post]
func (h *TakesHandler) JoinTrend(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userUUID, _ := uuid.Parse(userID.(string))
	takeIDStr := c.Param("take_id")
	takeID, err := uuid.Parse(takeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Take ID"})
		return
	}

	var req model.JoinTrendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trend, err := h.takesService.JoinOrCreateTrend(c.Request.Context(), userUUID, req.TrendKeyword, takeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    trend,
		"message": "Successfully joined trend!",
	})
}
