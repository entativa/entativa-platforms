package handler

import (
	"net/http"

	"socialink/post-service/internal/model"
	"socialink/post-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LikeHandler struct {
	likeService *service.LikeService
}

func NewLikeHandler(likeService *service.LikeService) *LikeHandler {
	return &LikeHandler{
		likeService: likeService,
	}
}

// LikePost adds a like/reaction to a post
// @Summary Like post
// @Description Add a reaction to a post
// @Tags likes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Param reaction body model.LikeRequest false "Reaction type (defaults to 'like')"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /posts/{post_id}/like [post]
func (h *LikeHandler) LikePost(c *gin.Context) {
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

	// Default to "like" reaction
	reactionType := model.ReactionLike

	// Check if specific reaction type provided
	var req model.LikeRequest
	if err := c.ShouldBindJSON(&req); err == nil {
		reactionType = req.ReactionType
	}

	if err := h.likeService.LikePost(c.Request.Context(), postID, userUUID, reactionType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to like post",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Post liked successfully",
	})
}

// UnlikePost removes a like from a post
// @Summary Unlike post
// @Description Remove a like/reaction from a post
// @Tags likes
// @Security BearerAuth
// @Produce json
// @Param post_id path string true "Post ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /posts/{post_id}/like [delete]
func (h *LikeHandler) UnlikePost(c *gin.Context) {
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

	if err := h.likeService.UnlikePost(c.Request.Context(), postID, userUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to unlike post",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Post unliked successfully",
	})
}

// LikeComment adds a like/reaction to a comment
// @Summary Like comment
// @Description Add a reaction to a comment
// @Tags likes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param comment_id path string true "Comment ID"
// @Param reaction body model.LikeRequest false "Reaction type"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /comments/{comment_id}/like [post]
func (h *LikeHandler) LikeComment(c *gin.Context) {
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

	commentIDStr := c.Param("comment_id")
	commentID, err := uuid.Parse(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid comment ID",
			"message": "The provided comment ID is not valid",
		})
		return
	}

	// Default to "like" reaction
	reactionType := model.ReactionLike

	var req model.LikeRequest
	if err := c.ShouldBindJSON(&req); err == nil {
		reactionType = req.ReactionType
	}

	if err := h.likeService.LikeComment(c.Request.Context(), commentID, userUUID, reactionType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to like comment",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Comment liked successfully",
	})
}

// UnlikeComment removes a like from a comment
// @Summary Unlike comment
// @Description Remove a like/reaction from a comment
// @Tags likes
// @Security BearerAuth
// @Produce json
// @Param comment_id path string true "Comment ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /comments/{comment_id}/like [delete]
func (h *LikeHandler) UnlikeComment(c *gin.Context) {
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

	commentIDStr := c.Param("comment_id")
	commentID, err := uuid.Parse(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid comment ID",
			"message": "The provided comment ID is not valid",
		})
		return
	}

	if err := h.likeService.UnlikeComment(c.Request.Context(), commentID, userUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to unlike comment",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Comment unliked successfully",
	})
}

// GetPostLikers retrieves users who liked a post
// @Summary Get post likers
// @Description Get users who liked a post
// @Tags likes
// @Produce json
// @Param post_id path string true "Post ID"
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} map[string]interface{}
// @Router /posts/{post_id}/likes [get]
func (h *LikeHandler) GetPostLikers(c *gin.Context) {
	postIDStr := c.Param("post_id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid post ID",
			"message": "The provided post ID is not valid",
		})
		return
	}

	limit := 50
	offset := 0
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if o := c.Query("offset"); o != "" {
		fmt.Sscanf(o, "%d", &offset)
	}

	userIDs, err := h.likeService.GetPostLikers(c.Request.Context(), postID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get likers",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    userIDs,
		"count":   len(userIDs),
	})
}
