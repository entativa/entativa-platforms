package handler

import (
	"net/http"
	"strconv"

	"vignette/search-service/internal/model"
	"vignette/search-service/internal/service"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	searchService *service.SearchService
}

func NewSearchHandler(searchService *service.SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

// Search performs a multi-entity search
// @Summary Multi-entity search
// @Description Search across users, posts, Takes, hashtags, and locations
// @Tags search
// @Accept json
// @Produce json
// @Param query query string true "Search query"
// @Param type query string false "Search type (user, post, take, hashtag, location, all)" default(all)
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Param verified query bool false "Filter verified users"
// @Param has_media query bool false "Filter posts with media"
// @Param media_type query string false "Filter by media type (image, video)"
// @Param min_likes query int false "Minimum likes"
// @Param min_views query int false "Minimum views"
// @Success 200 {object} model.SearchResponse
// @Router /search [get]
func (h *SearchHandler) Search(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
		return
	}

	searchType := model.SearchType(c.DefaultQuery("type", "all"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Build request
	req := &model.SearchRequest{
		Query:  query,
		Type:   searchType,
		Limit:  limit,
		Offset: offset,
	}

	// Parse filters
	if verified := c.Query("verified"); verified != "" {
		val := verified == "true"
		req.Filters.Verified = &val
	}

	if hasMedia := c.Query("has_media"); hasMedia != "" {
		val := hasMedia == "true"
		req.Filters.HasMedia = &val
	}

	if mediaType := c.Query("media_type"); mediaType != "" {
		req.Filters.MediaType = mediaType
	}

	if minLikes := c.Query("min_likes"); minLikes != "" {
		val, _ := strconv.Atoi(minLikes)
		req.Filters.MinLikes = val
	}

	if minViews := c.Query("min_views"); minViews != "" {
		val, _ := strconv.Atoi(minViews)
		req.Filters.MinViews = val
	}

	// Get user ID from header (optional)
	if userID := c.GetHeader("X-User-ID"); userID != "" {
		req.UserID = userID
	}

	// Perform search
	result, err := h.searchService.Search(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// SearchUsers searches for users
// @Summary Search users
// @Description Search for users by username, display name, or bio
// @Tags search
// @Produce json
// @Param query query string true "Search query"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Param verified query bool false "Filter verified users"
// @Param min_followers query int false "Minimum followers"
// @Success 200 {object} model.SearchResponse
// @Router /search/users [get]
func (h *SearchHandler) SearchUsers(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	req := &model.SearchRequest{
		Query:  query,
		Type:   model.SearchTypeUser,
		Limit:  limit,
		Offset: offset,
	}

	// Parse filters
	if verified := c.Query("verified"); verified != "" {
		val := verified == "true"
		req.Filters.Verified = &val
	}

	if minFollowers := c.Query("min_followers"); minFollowers != "" {
		val, _ := strconv.Atoi(minFollowers)
		req.Filters.MinFollowers = val
	}

	if userID := c.GetHeader("X-User-ID"); userID != "" {
		req.UserID = userID
	}

	result, err := h.searchService.Search(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// SearchPosts searches for posts
// @Summary Search posts
// @Description Search for posts by caption or content
// @Tags search
// @Produce json
// @Param query query string true "Search query"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} model.SearchResponse
// @Router /search/posts [get]
func (h *SearchHandler) SearchPosts(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	req := &model.SearchRequest{
		Query:  query,
		Type:   model.SearchTypePost,
		Limit:  limit,
		Offset: offset,
	}

	if userID := c.GetHeader("X-User-ID"); userID != "" {
		req.UserID = userID
	}

	result, err := h.searchService.Search(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// SearchTakes searches for Takes
// @Summary Search Takes
// @Description Search for Takes by caption
// @Tags search
// @Produce json
// @Param query query string true "Search query"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} model.SearchResponse
// @Router /search/takes [get]
func (h *SearchHandler) SearchTakes(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	req := &model.SearchRequest{
		Query:  query,
		Type:   model.SearchTypeTake,
		Limit:  limit,
		Offset: offset,
	}

	if userID := c.GetHeader("X-User-ID"); userID != "" {
		req.UserID = userID
	}

	result, err := h.searchService.Search(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetSearchHistory gets user's search history
// @Summary Get search history
// @Description Get user's recent searches
// @Tags search
// @Produce json
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} []model.SearchHistory
// @Router /search/history [get]
func (h *SearchHandler) GetSearchHistory(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID required"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	history, err := h.searchService.GetSearchHistory(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"history": history,
		"count":   len(history),
	})
}

// DeleteSearchHistory deletes user's search history
// @Summary Delete search history
// @Description Clear user's search history
// @Tags search
// @Success 200 {object} map[string]interface{}
// @Router /search/history [delete]
func (h *SearchHandler) DeleteSearchHistory(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID required"})
		return
	}

	err := h.searchService.DeleteSearchHistory(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Search history cleared",
	})
}

// GetTrendingSearches gets trending searches
// @Summary Get trending searches
// @Description Get currently trending search queries
// @Tags search
// @Produce json
// @Param type query string false "Search type" default(all)
// @Param limit query int false "Limit" default(10)
// @Success 200 {object} []model.TrendingSearch
// @Router /search/trending [get]
func (h *SearchHandler) GetTrendingSearches(c *gin.Context) {
	searchType := model.SearchType(c.DefaultQuery("type", "all"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	trending, err := h.searchService.GetTrendingSearches(c.Request.Context(), searchType, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"trending": trending,
		"count":    len(trending),
	})
}
