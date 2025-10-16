package handler

import (
	"bytes"
	"io"
	"net/http"

	"socialink/user-service/pkg/media"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	pb "socialink/user-service/proto/media"
)

type MediaHandler struct {
	mediaClient *media.Client
}

func NewMediaHandler(mediaClient *media.Client) *MediaHandler {
	return &MediaHandler{
		mediaClient: mediaClient,
	}
}

// UploadProfilePicture handles profile picture upload
// @Summary Upload profile picture
// @Description Upload a profile picture via media service
// @Tags media
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Profile Picture"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /media/profile-picture [post]
func (h *MediaHandler) UploadProfilePicture(c *gin.Context) {
	// Get authenticated user ID
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

	// Get file from request
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "File required",
			"message": "Please provide a file to upload",
		})
		return
	}
	defer file.Close()

	// Validate file size (max 10MB for profile picture)
	if header.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "File too large",
			"message": "Profile picture must be less than 10MB",
		})
		return
	}

	// Read file data
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to read file",
			"message": err.Error(),
		})
		return
	}

	// Upload via media service
	resp, err := h.mediaClient.UploadProfilePicture(
		c.Request.Context(),
		buf.Bytes(),
		header.Filename,
		userUUID.String(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Upload failed",
			"message": err.Error(),
		})
		return
	}

	// TODO: Update user profile with new profile picture URL
	// This would involve calling the profile service to update the profile_picture_url field

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"media_id":      resp.MediaId,
			"url":           resp.Url,
			"thumbnail_url": resp.ThumbnailUrl,
			"width":         resp.Width,
			"height":        resp.Height,
			"blurhash":      resp.Blurhash,
		},
	})
}

// UploadCoverPhoto handles cover photo upload
// @Summary Upload cover photo
// @Description Upload a cover photo via media service
// @Tags media
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Cover Photo"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /media/cover-photo [post]
func (h *MediaHandler) UploadCoverPhoto(c *gin.Context) {
	// Get authenticated user ID
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

	// Get file from request
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "File required",
			"message": "Please provide a file to upload",
		})
		return
	}
	defer file.Close()

	// Validate file size (max 20MB for cover photo)
	if header.Size > 20*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "File too large",
			"message": "Cover photo must be less than 20MB",
		})
		return
	}

	// Read file data
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to read file",
			"message": err.Error(),
		})
		return
	}

	// Upload via media service
	resp, err := h.mediaClient.UploadCoverPhoto(
		c.Request.Context(),
		buf.Bytes(),
		header.Filename,
		userUUID.String(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Upload failed",
			"message": err.Error(),
		})
		return
	}

	// TODO: Update user profile with new cover photo URL
	// This would involve calling the profile service to update the cover_photo_url field

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"media_id":      resp.MediaId,
			"url":           resp.Url,
			"thumbnail_url": resp.ThumbnailUrl,
			"width":         resp.Width,
			"height":        resp.Height,
			"blurhash":      resp.Blurhash,
		},
	})
}

// GetMedia retrieves media information
// @Summary Get media info
// @Description Get media information by ID
// @Tags media
// @Security BearerAuth
// @Produce json
// @Param media_id path string true "Media ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /media/{media_id} [get]
func (h *MediaHandler) GetMedia(c *gin.Context) {
	mediaID := c.Param("media_id")

	resp, err := h.mediaClient.GetMedia(c.Request.Context(), mediaID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Media not found",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
	})
}

// DeleteMedia deletes a media file
// @Summary Delete media
// @Description Delete a media file
// @Tags media
// @Security BearerAuth
// @Produce json
// @Param media_id path string true "Media ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /media/{media_id} [delete]
func (h *MediaHandler) DeleteMedia(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User not authenticated",
		})
		return
	}

	mediaID := c.Param("media_id")

	err := h.mediaClient.DeleteMedia(c.Request.Context(), mediaID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Delete failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Media deleted successfully",
	})
}
