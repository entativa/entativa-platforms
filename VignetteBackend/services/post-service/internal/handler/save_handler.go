package handler

import (
	"fmt"
	"net/http"

	"vignette/post-service/internal/model"
	"vignette/post-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SaveHandler struct {
	postService *service.PostService
}

func NewSaveHandler(postService *service.PostService) *SaveHandler {
	return &SaveHandler{
		postService: postService,
	}
}

// SavePost saves a post to user's saved collection
// @Summary Save post
// @Description Save/bookmark a post
// @Tags saves
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Param request body model.SavePostRequest false "Save options"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /posts/{post_id}/save [post]
func (h *SaveHandler) SavePost(c *gin.Context) {
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
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Invalid user ID",
			"message": "The user ID is not valid",
		})
		return
	}

	postIDStr := c.Param("post_id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid post ID",
			"message": "The provided post ID is not valid",
		})
		return
	}

	var req model.SavePostRequest
	c.ShouldBindJSON(&req)

	if err := h.postService.SavePost(c.Request.Context(), userUUID, postID, req.Collection); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to save post",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Post saved successfully",
	})
}

// UnsavePost removes a post from saved collection
// @Summary Unsave post
// @Description Remove a post from saved collection
// @Tags saves
// @Security BearerAuth
// @Produce json
// @Param post_id path string true "Post ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /posts/{post_id}/save [delete]
func (h *SaveHandler) UnsavePost(c *gin.Context) {
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
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Invalid user ID",
			"message": "The user ID is not valid",
		})
		return
	}

	postIDStr := c.Param("post_id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid post ID",
			"message": "The provided post ID is not valid",
		})
		return
	}

	if err := h.postService.UnsavePost(c.Request.Context(), userUUID, postID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to unsave post",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Post unsaved successfully",
	})
}

// GetSavedPosts retrieves user's saved posts
// @Summary Get saved posts
// @Description Get all saved/bookmarked posts
// @Tags saves
// @Security BearerAuth
// @Produce json
// @Param collection query string false "Collection name"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} map[string]interface{}
// @Router /saved [get]
func (h *SaveHandler) GetSavedPosts(c *gin.Context) {
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
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Invalid user ID",
			"message": "The user ID is not valid",
		})
		return
	}

	var collection *string
	if c := c.Query("collection"); c != "" {
		collection = &c
	}

	limit := 20
	offset := 0
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if o := c.Query("offset"); o != "" {
		fmt.Sscanf(o, "%d", &offset)
	}

	saves, err := h.postService.GetSavedPosts(c.Request.Context(), userUUID, collection, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get saved posts",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    saves,
		"count":   len(saves),
	})
}
