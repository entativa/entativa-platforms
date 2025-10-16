package handler

import (
	"fmt"
	"net/http"

	"socialink/post-service/internal/model"
	"socialink/post-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CommentHandler struct {
	commentService *service.CommentService
}

func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

// CreateComment creates a new comment on a post
// @Summary Create comment
// @Description Add a comment to a post
// @Tags comments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Param comment body model.CreateCommentRequest true "Comment data"
// @Success 201 {object} model.CommentResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /posts/{post_id}/comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
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

	var req model.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	comment, err := h.commentService.CreateComment(c.Request.Context(), postID, userUUID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create comment",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    comment,
	})
}

// GetComments retrieves comments for a post
// @Summary Get comments
// @Description Get all comments for a post
// @Tags comments
// @Produce json
// @Param post_id path string true "Post ID"
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} []model.CommentResponse
// @Router /posts/{post_id}/comments [get]
func (h *CommentHandler) GetComments(c *gin.Context) {
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

	comments, err := h.commentService.GetComments(c.Request.Context(), postID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get comments",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    comments,
		"count":   len(comments),
	})
}

// GetReplies retrieves replies to a comment
// @Summary Get comment replies
// @Description Get all replies to a specific comment
// @Tags comments
// @Produce json
// @Param comment_id path string true "Comment ID"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} []model.CommentResponse
// @Router /comments/{comment_id}/replies [get]
func (h *CommentHandler) GetReplies(c *gin.Context) {
	commentIDStr := c.Param("comment_id")
	commentID, err := uuid.Parse(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid comment ID",
			"message": "The provided comment ID is not valid",
		})
		return
	}

	limit := 20
	offset := 0
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if o := c.Query("offset"); o != "" {
		fmt.Sscanf(o, "%d", &offset)
	}

	replies, err := h.commentService.GetReplies(c.Request.Context(), commentID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get replies",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    replies,
		"count":   len(replies),
	})
}

// UpdateComment updates a comment
// @Summary Update comment
// @Description Update an existing comment
// @Tags comments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param comment_id path string true "Comment ID"
// @Param comment body model.UpdateCommentRequest true "Update data"
// @Success 200 {object} model.CommentResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /comments/{comment_id} [put]
func (h *CommentHandler) UpdateComment(c *gin.Context) {
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

	var req model.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	comment, err := h.commentService.UpdateComment(c.Request.Context(), commentID, userUUID, &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "permission denied: not the comment owner" {
			statusCode = http.StatusForbidden
		}

		c.JSON(statusCode, gin.H{
			"error":   "Failed to update comment",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    comment,
	})
}

// DeleteComment deletes a comment
// @Summary Delete comment
// @Description Delete a comment
// @Tags comments
// @Security BearerAuth
// @Produce json
// @Param comment_id path string true "Comment ID"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /comments/{comment_id} [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
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

	if err := h.commentService.DeleteComment(c.Request.Context(), commentID, userUUID); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "permission denied: not the comment owner" {
			statusCode = http.StatusForbidden
		}

		c.JSON(statusCode, gin.H{
			"error":   "Failed to delete comment",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Comment deleted successfully",
	})
}
