package handler

import (
	"net/http"

	"vignette/user-service/internal/middleware"
	"vignette/user-service/internal/repository"
	"vignette/user-service/internal/service"

	"github.com/gin-gonic/gin"
)

type MediaHandler struct {
	mediaService *service.MediaService
	userRepo     *repository.UserRepository
}

func NewMediaHandler(mediaService *service.MediaService, userRepo *repository.UserRepository) *MediaHandler {
	return &MediaHandler{
		mediaService: mediaService,
		userRepo:     userRepo,
	}
}

// UploadProfilePicture handles profile picture upload
// @Summary Upload Profile Picture
// @Description Upload and set profile picture
// @Tags media
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Profile picture"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /media/profile-picture [post]
func (h *MediaHandler) UploadProfilePicture(c *gin.Context) {
	if h.mediaService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":   "Media upload not configured",
			"message": "Please configure S3/MinIO to enable media uploads",
		})
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "No file provided",
			"message": "Please provide a file to upload",
		})
		return
	}

	// Upload file
	url, err := h.mediaService.UploadProfilePicture(userID, file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Upload failed",
			"message": err.Error(),
		})
		return
	}

	// Update user profile
	user, err := h.userRepo.FindByID(userID)
	if err == nil {
		user.ProfilePictureURL = &url
		_ = h.userRepo.Update(user)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile picture uploaded successfully",
		"data": gin.H{
			"url": url,
		},
	})
}

// UploadCoverPhoto handles cover photo upload
// @Summary Upload Cover Photo
// @Description Upload and set cover photo
// @Tags media
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Cover photo"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /media/cover-photo [post]
func (h *MediaHandler) UploadCoverPhoto(c *gin.Context) {
	if h.mediaService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":   "Media upload not configured",
			"message": "Please configure S3/MinIO to enable media uploads",
		})
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "No file provided",
			"message": "Please provide a file to upload",
		})
		return
	}

	// Upload file
	url, err := h.mediaService.UploadCoverPhoto(userID, file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Upload failed",
			"message": err.Error(),
		})
		return
	}

	// Update user profile
	user, err := h.userRepo.FindByID(userID)
	if err == nil {
		user.CoverPhotoURL = &url
		_ = h.userRepo.Update(user)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cover photo uploaded successfully",
		"data": gin.H{
			"url": url,
		},
	})
}
