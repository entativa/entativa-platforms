package handler

import (
	"fmt"
	"net/http"

	"socialink/post-service/internal/model"
	"socialink/post-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ShareHandler struct {
	shareService *service.ShareService
}

func NewShareHandler(shareService *service.ShareService) *ShareHandler {
	return &ShareHandler{
		shareService: shareService,
	}
}

// SharePost creates a share of a post
// @Summary Share post
// @Description Share a post to your timeline
// @Tags shares
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Param share body model.SharePostRequest true "Share data"
// @Success 201 {object} model.Share
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /posts/{post_id}/share [post]
func (h *ShareHandler) SharePost(c *gin.Context) {
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

	var req model.SharePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	share, err := h.shareService.SharePost(c.Request.Context(), postID, userUUID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to share post",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    share,
	})
}

// GetPostShares retrieves shares of a post
// @Summary Get post shares
// @Description Get all shares of a specific post
// @Tags shares
// @Produce json
// @Param post_id path string true "Post ID"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} []model.Share
// @Router /posts/{post_id}/shares [get]
func (h *ShareHandler) GetPostShares(c *gin.Context) {
	postIDStr := c.Param("post_id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid post ID",
			"message": "The provided post ID is not valid",
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

	shares, err := h.shareService.GetPostShares(c.Request.Context(), postID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get shares",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    shares,
		"count":   len(shares),
	})
}

// DeleteShare removes a share
// @Summary Delete share
// @Description Delete a share
// @Tags shares
// @Security BearerAuth
// @Produce json
// @Param share_id path string true "Share ID"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /shares/{share_id} [delete]
func (h *ShareHandler) DeleteShare(c *gin.Context) {
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

	shareIDStr := c.Param("share_id")
	shareID, err := uuid.Parse(shareIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid share ID",
			"message": "The provided share ID is not valid",
		})
		return
	}

	if err := h.shareService.DeleteShare(c.Request.Context(), shareID, userUUID); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "permission denied: not the share owner" {
			statusCode = http.StatusForbidden
		}

		c.JSON(statusCode, gin.H{
			"error":   "Failed to delete share",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Share deleted successfully",
	})
}
