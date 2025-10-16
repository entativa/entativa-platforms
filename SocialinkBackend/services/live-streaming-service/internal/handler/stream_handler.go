package handler

import (
	"net/http"

	"socialink/live-streaming-service/internal/model"
	"socialink/live-streaming-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StreamHandler struct {
	service *service.StreamingService
}

func NewStreamHandler(service *service.StreamingService) *StreamHandler {
	return &StreamHandler{service: service}
}

// CreateStream creates a new stream
func (h *StreamHandler) CreateStream(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))

	var req model.CreateStreamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stream, err := h.service.CreateStream(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, stream)
}

// StartStream starts a stream
func (h *StreamHandler) StartStream(c *gin.Context) {
	streamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stream ID"})
		return
	}

	userID, _ := uuid.Parse(c.GetString("user_id"))

	err = h.service.StartStream(c.Request.Context(), streamID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "stream started", "stream_id": streamID})
}

// EndStream ends a stream
func (h *StreamHandler) EndStream(c *gin.Context) {
	streamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stream ID"})
		return
	}

	userID, _ := uuid.Parse(c.GetString("user_id"))

	err = h.service.EndStream(c.Request.Context(), streamID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "stream ended"})
}

// GetStream gets stream details
func (h *StreamHandler) GetStream(c *gin.Context) {
	streamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stream ID"})
		return
	}

	stream, err := h.service.GetStream(c.Request.Context(), streamID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "stream not found"})
		return
	}

	c.JSON(http.StatusOK, stream)
}

// GetLiveStreams gets all live streams
func (h *StreamHandler) GetLiveStreams(c *gin.Context) {
	limit := 20
	offset := 0
	category := c.Query("category")

	streams, err := h.service.GetLiveStreams(c.Request.Context(), limit, offset, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"streams": streams, "count": len(streams)})
}

// CheckEligibility checks if user can go live
func (h *StreamHandler) CheckEligibility(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))

	eligibility, err := h.service.CheckEligibility(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, eligibility)
}

// PostComment posts a comment on a stream
func (h *StreamHandler) PostComment(c *gin.Context) {
	streamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stream ID"})
		return
	}

	userID, _ := uuid.Parse(c.GetString("user_id"))

	var req model.StreamCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment, err := h.service.PostComment(c.Request.Context(), streamID, userID, req.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// AddReaction adds a reaction to a stream
func (h *StreamHandler) AddReaction(c *gin.Context) {
	streamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stream ID"})
		return
	}

	userID, _ := uuid.Parse(c.GetString("user_id"))

	var req model.StreamReactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.AddReaction(c.Request.Context(), streamID, userID, req.Type)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "reaction added"})
}

// GetAnalytics gets stream analytics
func (h *StreamHandler) GetAnalytics(c *gin.Context) {
	streamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stream ID"})
		return
	}

	analytics, err := h.service.GetStreamAnalytics(c.Request.Context(), streamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analytics)
}
