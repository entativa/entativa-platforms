package handler

import (
	"net/http"
	"strconv"

	"vignette/search-service/internal/model"
	"vignette/search-service/internal/service"

	"github.com/gin-gonic/gin"
)

type AutocompleteHandler struct {
	autocompleteService *service.AutocompleteService
}

func NewAutocompleteHandler(autocompleteService *service.AutocompleteService) *AutocompleteHandler {
	return &AutocompleteHandler{
		autocompleteService: autocompleteService,
	}
}

// Autocomplete provides search suggestions
// @Summary Get autocomplete suggestions
// @Description Get search suggestions as user types
// @Tags autocomplete
// @Produce json
// @Param query query string true "Search query (min 2 characters)"
// @Param type query string false "Search type" default(all)
// @Param limit query int false "Limit" default(10)
// @Success 200 {object} model.AutocompleteResponse
// @Router /autocomplete [get]
func (h *AutocompleteHandler) Autocomplete(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
		return
	}

	if len(query) < 2 {
		c.JSON(http.StatusOK, &model.AutocompleteResponse{
			Query:       query,
			Suggestions: []model.AutocompleteSuggestion{},
			Took:        0,
		})
		return
	}

	searchType := model.SearchType(c.DefaultQuery("type", "all"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	req := &model.AutocompleteRequest{
		Query: query,
		Type:  searchType,
		Limit: limit,
	}

	if userID := c.GetHeader("X-User-ID"); userID != "" {
		req.UserID = userID
	}

	result, err := h.autocompleteService.Autocomplete(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetRecentSearches gets user's recent searches for autocomplete
// @Summary Get recent searches
// @Description Get user's recent search queries for autocomplete
// @Tags autocomplete
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Success 200 {object} []model.AutocompleteSuggestion
// @Router /autocomplete/recent [get]
func (h *AutocompleteHandler) GetRecentSearches(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID required"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	suggestions, err := h.autocompleteService.GetRecentSearches(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"suggestions": suggestions,
		"count":       len(suggestions),
	})
}
