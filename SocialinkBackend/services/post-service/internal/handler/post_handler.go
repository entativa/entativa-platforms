package handler

import (
	"net/http"

	"socialink/post-service/internal/model"
	"socialink/post-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// CreatePost creates a new post
// @Summary Create post
// @Description Create a new social media post
// @Tags posts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param post body model.CreatePostRequest true "Post data"
// @Success 201 {object} model.PostResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
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

	var req model.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	post, err := h.postService.CreatePost(c.Request.Context(), userUUID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create post",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    post,
	})
}

// GetPost retrieves a single post
// @Summary Get post
// @Description Get a post by ID
// @Tags posts
// @Security BearerAuth
// @Produce json
// @Param post_id path string true "Post ID"
// @Success 200 {object} model.PostResponse
// @Failure 404 {object} map[string]interface{}
// @Router /posts/{post_id} [get]
func (h *PostHandler) GetPost(c *gin.Context) {
	postIDStr := c.Param("post_id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid post ID",
			"message": "The provided post ID is not valid",
		})
		return
	}

	// Get requesting user ID (optional for public posts)
	var requestingUserID uuid.UUID
	if userID, exists := c.Get("user_id"); exists {
		requestingUserID, _ = uuid.Parse(userID.(string))
	}

	post, err := h.postService.GetPost(c.Request.Context(), postID, requestingUserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Post not found",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    post,
	})
}

// UpdatePost updates a post
// @Summary Update post
// @Description Update an existing post
// @Tags posts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Param post body model.UpdatePostRequest true "Update data"
// @Success 200 {object} model.PostResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /posts/{post_id} [put]
func (h *PostHandler) UpdatePost(c *gin.Context) {
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

	var req model.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	post, err := h.postService.UpdatePost(c.Request.Context(), postID, userUUID, &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "permission denied: not the post owner" {
			statusCode = http.StatusForbidden
		}

		c.JSON(statusCode, gin.H{
			"error":   "Failed to update post",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    post,
	})
}

// DeletePost deletes a post
// @Summary Delete post
// @Description Delete a post
// @Tags posts
// @Security BearerAuth
// @Produce json
// @Param post_id path string true "Post ID"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /posts/{post_id} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
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

	if err := h.postService.DeletePost(c.Request.Context(), postID, userUUID); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "permission denied: not the post owner" {
			statusCode = http.StatusForbidden
		}

		c.JSON(statusCode, gin.H{
			"error":   "Failed to delete post",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Post deleted successfully",
	})
}

// GetUserPosts retrieves posts by a user
// @Summary Get user posts
// @Description Get all posts by a specific user
// @Tags posts
// @Security BearerAuth
// @Produce json
// @Param user_id path string true "User ID"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} model.PostListResponse
// @Router /posts/user/{user_id} [get]
func (h *PostHandler) GetUserPosts(c *gin.Context) {
	targetUserIDStr := c.Param("user_id")
	targetUserID, err := uuid.Parse(targetUserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"message": "The provided user ID is not valid",
		})
		return
	}

	// Get requesting user ID (optional)
	var requestingUserID uuid.UUID
	if userID, exists := c.Get("user_id"); exists {
		requestingUserID, _ = uuid.Parse(userID.(string))
	}

	limit := 20
	offset := 0
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if o := c.Query("offset"); o != "" {
		fmt.Sscanf(o, "%d", &offset)
	}

	posts, err := h.postService.GetUserPosts(c.Request.Context(), targetUserID, requestingUserID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get posts",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    posts,
		"count":   len(posts),
	})
}

// GetFeed retrieves the user's personalized feed
// @Summary Get feed
// @Description Get personalized feed for authenticated user
// @Tags posts
// @Security BearerAuth
// @Produce json
// @Param cursor query string false "Cursor for pagination"
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} model.PostListResponse
// @Router /posts/feed [get]
func (h *PostHandler) GetFeed(c *gin.Context) {
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

	cursor := c.Query("cursor")
	limit := 20
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	posts, nextCursor, err := h.postService.GetFeed(c.Request.Context(), userUUID, cursor, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get feed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"data":        posts,
		"next_cursor": nextCursor,
		"has_more":    nextCursor != nil,
	})
}

// GetTrendingPosts retrieves trending posts
// @Summary Get trending posts
// @Description Get currently trending posts
// @Tags posts
// @Produce json
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} model.PostListResponse
// @Router /posts/trending [get]
func (h *PostHandler) GetTrendingPosts(c *gin.Context) {
	limit := 20
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	posts, err := h.postService.GetTrendingPosts(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get trending posts",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    posts,
		"count":   len(posts),
	})
}
