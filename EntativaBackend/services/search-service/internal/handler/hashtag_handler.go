package handler

import (
	"net/http"
	"strconv"

	"socialink/search-service/internal/model"
	"socialink/search-service/internal/service"

	"github.com/gin-gonic/gin"
)

type HashtagHandler struct {
	hashtagService *service.HashtagService
}

func NewHashtagHandler(hashtagService *service.HashtagService) *HashtagHandler {
	return &HashtagHandler{
		hashtagService: hashtagService,
	}
}

// GetTrendingHashtags gets trending hashtags
// @Summary Get trending hashtags
// @Description Get currently trending hashtags with growth rates
// @Tags hashtags
// @Produce json
// @Param limit query int false "Limit" default(20)
// @Param category query string false "Category filter"
// @Success 200 {object} model.TrendingHashtagsResponse
// @Router /hashtags/trending [get]
func (h *HashtagHandler) GetTrendingHashtags(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	category := c.Query("category")

	req := &model.TrendingHashtagsRequest{
		Limit:    limit,
		Category: category,
	}

	result, err := h.hashtagService.GetTrendingHashtags(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetRelatedHashtags gets hashtags related to a given hashtag
// @Summary Get related hashtags
// @Description Get hashtags similar to the given hashtag
// @Tags hashtags
// @Produce json
// @Param tag path string true "Hashtag (without #)"
// @Param limit query int false "Limit" default(10)
// @Success 200 {object} model.RelatedHashtagsResponse
// @Router /hashtags/{tag}/related [get]
func (h *HashtagHandler) GetRelatedHashtags(c *gin.Context) {
	tag := c.Param("tag")
	if tag == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tag parameter is required"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	req := &model.RelatedHashtagsRequest{
		Tag:   tag,
		Limit: limit,
	}

	result, err := h.hashtagService.GetRelatedHashtags(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// SearchHashtags searches for hashtags
// @Summary Search hashtags
// @Description Search for hashtags by name
// @Tags hashtags
// @Produce json
// @Param query query string true "Search query"
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} []model.HashtagStats
// @Router /hashtags/search [get]
func (h *HashtagHandler) SearchHashtags(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	req := &model.HashtagSearchRequest{
		Query: query,
		Limit: limit,
	}

	hashtags, err := h.hashtagService.SearchHashtags(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"hashtags": hashtags,
		"count":    len(hashtags),
	})
}
