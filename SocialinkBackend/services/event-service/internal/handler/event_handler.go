package handler

import (
	"net/http"
	"strconv"

	"socialink/event-service/internal/model"
	"socialink/event-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EventHandler struct {
	service *service.EventService
}

func NewEventHandler(service *service.EventService) *EventHandler {
	return &EventHandler{service: service}
}

// CreateEvent creates a new event
func (h *EventHandler) CreateEvent(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))

	var req model.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event, err := h.service.CreateEvent(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// GetEvent gets event details
func (h *EventHandler) GetEvent(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	var viewerID *uuid.UUID
	if userIDStr := c.GetString("user_id"); userIDStr != "" {
		uid, _ := uuid.Parse(userIDStr)
		viewerID = &uid
	}

	event, err := h.service.GetEvent(c.Request.Context(), eventID, viewerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	c.JSON(http.StatusOK, event)
}

// UpdateEvent updates an event
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	userID, _ := uuid.Parse(c.GetString("user_id"))

	var req model.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.UpdateEvent(c.Request.Context(), eventID, userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "event updated"})
}

// CancelEvent cancels an event
func (h *EventHandler) CancelEvent(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	userID, _ := uuid.Parse(c.GetString("user_id"))

	err = h.service.CancelEvent(c.Request.Context(), eventID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "event cancelled"})
}

// RSVP to event
func (h *EventHandler) RSVP(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	userID, _ := uuid.Parse(c.GetString("user_id"))

	var req model.RSVPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.RSVP(c.Request.Context(), eventID, userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "RSVP updated", "status": req.Status})
}

// RemoveRSVP removes RSVP
func (h *EventHandler) RemoveRSVP(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	userID, _ := uuid.Parse(c.GetString("user_id"))

	err = h.service.RemoveRSVP(c.Request.Context(), eventID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "RSVP removed"})
}

// CheckIn to event
func (h *EventHandler) CheckIn(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	userID, _ := uuid.Parse(c.GetString("user_id"))

	err = h.service.CheckIn(c.Request.Context(), eventID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "checked in successfully"})
}

// GetUpcomingEvents gets upcoming events
func (h *EventHandler) GetUpcomingEvents(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	category := c.Query("category")

	events, err := h.service.GetUpcomingEvents(c.Request.Context(), limit, offset, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events, "count": len(events)})
}

// SearchEvents searches events
func (h *EventHandler) SearchEvents(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter required"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	events, err := h.service.SearchEvents(c.Request.Context(), query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events, "count": len(events)})
}

// GetNearbyEvents gets events near location
func (h *EventHandler) GetNearbyEvents(c *gin.Context) {
	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	if latStr == "" || lngStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lat and lng required"})
		return
	}

	lat, _ := strconv.ParseFloat(latStr, 64)
	lng, _ := strconv.ParseFloat(lngStr, 64)
	radiusKm, _ := strconv.Atoi(c.DefaultQuery("radius", "50"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	events, err := h.service.GetNearbyEvents(c.Request.Context(), lat, lng, radiusKm, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events, "count": len(events)})
}

// GetUserEvents gets user's events
func (h *EventHandler) GetUserEvents(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))

	statusStr := c.Query("status")
	var status model.RSVPStatus
	if statusStr != "" {
		status = model.RSVPStatus(statusStr)
	}

	upcoming := c.Query("upcoming") == "true"
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	events, err := h.service.GetUserEvents(c.Request.Context(), userID, status, upcoming, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events, "count": len(events)})
}

// GetEventAttendees gets event attendees
func (h *EventHandler) GetEventAttendees(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	statusStr := c.Query("status")
	var status model.RSVPStatus
	if statusStr != "" {
		status = model.RSVPStatus(statusStr)
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	attendees, err := h.service.GetEventAttendees(c.Request.Context(), eventID, status, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"attendees": attendees, "count": len(attendees)})
}

// GetEventStats gets event statistics
func (h *EventHandler) GetEventStats(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	stats, err := h.service.GetEventStats(c.Request.Context(), eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
