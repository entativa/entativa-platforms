package service

import (
	"context"
	"fmt"
	"time"

	"github.com/entativa/socialink/event-service/internal/model"
	"github.com/entativa/socialink/event-service/internal/repository"
	"github.com/google/uuid"
)

type EventService struct {
	eventRepo *repository.EventRepository
	rsvpRepo  *repository.RSVPRepository
	kafka     *KafkaProducer
	redis     *RedisClient
}

func NewEventService(
	eventRepo *repository.EventRepository,
	rsvpRepo *repository.RSVPRepository,
	kafka *KafkaProducer,
	redis *RedisClient,
) *EventService {
	return &EventService{
		eventRepo: eventRepo,
		rsvpRepo:  rsvpRepo,
		kafka:     kafka,
		redis:     redis,
	}
}

// CreateEvent creates a new event
func (s *EventService) CreateEvent(ctx context.Context, req *model.CreateEventRequest, creatorID uuid.UUID) (*model.Event, error) {
	// Validate event type and location
	if req.Type == model.EventTypeInPerson {
		if req.LocationName == nil || *req.LocationName == "" {
			return nil, fmt.Errorf("location is required for in-person events")
		}
	} else if req.Type == model.EventTypeOnline {
		if req.OnlineLink == nil || *req.OnlineLink == "" {
			return nil, fmt.Errorf("online link is required for virtual events")
		}
	}

	// Validate start time is in future
	if req.StartTime.Before(time.Now()) {
		return nil, fmt.Errorf("start time must be in the future")
	}

	// Validate end time is after start time
	if req.EndTime != nil && req.EndTime.Before(req.StartTime) {
		return nil, fmt.Errorf("end time must be after start time")
	}

	// Validate recurring settings
	if req.IsRecurring {
		if req.RecurrenceRule == nil || *req.RecurrenceRule == "" {
			return nil, fmt.Errorf("recurrence rule is required for recurring events")
		}
	}

	event := &model.Event{
		ID:                uuid.New(),
		CreatorID:         creatorID,
		Title:             req.Title,
		Description:       req.Description,
		CoverPhoto:        req.CoverPhoto,
		Type:              req.Type,
		Category:          req.Category,
		Privacy:           req.Privacy,
		LocationName:      req.LocationName,
		Address:           req.Address,
		City:              req.City,
		Country:           req.Country,
		Latitude:          req.Latitude,
		Longitude:         req.Longitude,
		OnlineLink:        req.OnlineLink,
		StartTime:         req.StartTime,
		EndTime:           req.EndTime,
		Timezone:          req.Timezone,
		IsRecurring:       req.IsRecurring,
		RecurrenceRule:    req.RecurrenceRule,
		RecurrenceEndDate: req.RecurrenceEndDate,
		AllowGuestInvites: req.AllowGuestInvites,
		RequireApproval:   req.RequireApproval,
		MaxAttendees:      req.MaxAttendees,
		CoHosts:           req.CoHosts,
		GoingCount:        0,
		InterestedCount:   0,
		ViewCount:         0,
		IsCancelled:       false,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	err := s.eventRepo.Create(ctx, event)
	if err != nil {
		return nil, err
	}

	// Publish event created (for notifications)
	s.kafka.PublishEventCreated(event.ID, event.CreatorID)

	return event, nil
}

// GetEvent gets event by ID
func (s *EventService) GetEvent(ctx context.Context, eventID uuid.UUID, viewerID *uuid.UUID) (*model.Event, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	// Check privacy
	if event.Privacy == model.PrivacyPrivate && viewerID != nil {
		// Check if viewer is creator, co-host, or invited
		if event.CreatorID != *viewerID {
			isCoHost, _ := s.eventRepo.IsCoHost(ctx, eventID, *viewerID)
			if !isCoHost {
				// Check if invited (simplified - in production check invites table)
				return nil, fmt.Errorf("event not found")
			}
		}
	}

	// Increment view count (async)
	go s.eventRepo.IncrementViewCount(context.Background(), eventID)

	return event, nil
}

// UpdateEvent updates an event
func (s *EventService) UpdateEvent(ctx context.Context, eventID, userID uuid.UUID, req *model.UpdateEventRequest) error {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	// Check authorization (creator or co-host)
	if event.CreatorID != userID {
		isCoHost, _ := s.eventRepo.IsCoHost(ctx, eventID, userID)
		if !isCoHost {
			return fmt.Errorf("unauthorized")
		}
	}

	// Update fields
	if req.Title != nil {
		event.Title = *req.Title
	}
	if req.Description != nil {
		event.Description = *req.Description
	}
	if req.CoverPhoto != nil {
		event.CoverPhoto = req.CoverPhoto
	}
	if req.StartTime != nil {
		event.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		event.EndTime = req.EndTime
	}
	if req.OnlineLink != nil {
		event.OnlineLink = req.OnlineLink
	}

	event.UpdatedAt = time.Now()

	err = s.eventRepo.Update(ctx, event)
	if err != nil {
		return err
	}

	// Publish event updated
	s.kafka.PublishEventUpdated(eventID, userID)

	return nil
}

// CancelEvent cancels an event
func (s *EventService) CancelEvent(ctx context.Context, eventID, userID uuid.UUID) error {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	// Only creator can cancel
	if event.CreatorID != userID {
		return fmt.Errorf("only event creator can cancel event")
	}

	if event.IsCancelled {
		return fmt.Errorf("event already cancelled")
	}

	err = s.eventRepo.Cancel(ctx, eventID)
	if err != nil {
		return err
	}

	// Publish event cancelled (notify all attendees)
	s.kafka.PublishEventCancelled(eventID, userID)

	return nil
}

// RSVP to event
func (s *EventService) RSVP(ctx context.Context, eventID, userID uuid.UUID, req *model.RSVPRequest) error {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	if event.IsCancelled {
		return fmt.Errorf("cannot RSVP to cancelled event")
	}

	// Check max attendees
	if req.Status == model.RSVPGoing && event.MaxAttendees != nil {
		if event.GoingCount >= *event.MaxAttendees {
			return fmt.Errorf("event is full")
		}
	}

	// Check if already RSVPed
	existing, _ := s.rsvpRepo.Get(ctx, eventID, userID)

	rsvp := &model.EventRSVP{
		ID:         uuid.New(),
		EventID:    eventID,
		UserID:     userID,
		Status:     req.Status,
		GuestCount: req.GuestCount,
		CheckedIn:  false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if existing != nil {
		rsvp.ID = existing.ID
	}

	err = s.rsvpRepo.Upsert(ctx, rsvp)
	if err != nil {
		return err
	}

	// Publish RSVP event
	s.kafka.PublishEventRSVP(eventID, userID, string(req.Status))

	return nil
}

// RemoveRSVP removes user's RSVP
func (s *EventService) RemoveRSVP(ctx context.Context, eventID, userID uuid.UUID) error {
	return s.rsvpRepo.Delete(ctx, eventID, userID)
}

// CheckIn user to event
func (s *EventService) CheckIn(ctx context.Context, eventID, userID uuid.UUID) error {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	// Check if event is happening now
	now := time.Now()
	if now.Before(event.StartTime.Add(-1 * time.Hour)) {
		return fmt.Errorf("check-in not available yet (opens 1 hour before event)")
	}

	if event.EndTime != nil && now.After(*event.EndTime) {
		return fmt.Errorf("event has ended")
	}

	// Check if user RSVP'd as going
	rsvp, err := s.rsvpRepo.Get(ctx, eventID, userID)
	if err != nil {
		return err
	}
	if rsvp == nil || rsvp.Status != model.RSVPGoing {
		return fmt.Errorf("must RSVP as 'going' to check in")
	}

	if rsvp.CheckedIn {
		return fmt.Errorf("already checked in")
	}

	return s.rsvpRepo.CheckIn(ctx, eventID, userID)
}

// GetUpcomingEvents gets upcoming public events
func (s *EventService) GetUpcomingEvents(ctx context.Context, limit, offset int, category string) ([]*model.Event, error) {
	return s.eventRepo.GetUpcoming(ctx, limit, offset, category)
}

// SearchEvents searches events
func (s *EventService) SearchEvents(ctx context.Context, query string, limit, offset int) ([]*model.Event, error) {
	return s.eventRepo.Search(ctx, query, limit, offset)
}

// GetNearbyEvents gets events near a location
func (s *EventService) GetNearbyEvents(ctx context.Context, lat, lng float64, radiusKm, limit int) ([]*model.Event, error) {
	return s.eventRepo.GetNearby(ctx, lat, lng, radiusKm, limit)
}

// GetUserEvents gets events user is attending/interested in
func (s *EventService) GetUserEvents(ctx context.Context, userID uuid.UUID, status model.RSVPStatus, upcoming bool, limit, offset int) ([]*model.Event, error) {
	// Get event IDs
	eventIDs, err := s.rsvpRepo.GetUserEvents(ctx, userID, status, upcoming, limit, offset)
	if err != nil {
		return nil, err
	}

	// Get full event details
	events := []*model.Event{}
	for _, eventID := range eventIDs {
		event, err := s.eventRepo.GetByID(ctx, *eventID)
		if err == nil {
			events = append(events, event)
		}
	}

	return events, nil
}

// GetEventAttendees gets list of attendees
func (s *EventService) GetEventAttendees(ctx context.Context, eventID uuid.UUID, status model.RSVPStatus, limit, offset int) ([]*model.EventRSVP, error) {
	return s.rsvpRepo.GetByEvent(ctx, eventID, status, limit, offset)
}

// GetEventStats gets event statistics
func (s *EventService) GetEventStats(ctx context.Context, eventID uuid.UUID) (*model.EventStats, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	stats := &model.EventStats{
		EventID:         eventID,
		GoingCount:      event.GoingCount,
		InterestedCount: event.InterestedCount,
		ViewCount:       event.ViewCount,
	}

	return stats, nil
}

// Stub types for Kafka and Redis (to be implemented)
type KafkaProducer struct{}

func (k *KafkaProducer) PublishEventCreated(eventID, creatorID uuid.UUID)          {}
func (k *KafkaProducer) PublishEventUpdated(eventID, userID uuid.UUID)             {}
func (k *KafkaProducer) PublishEventCancelled(eventID, userID uuid.UUID)           {}
func (k *KafkaProducer) PublishEventRSVP(eventID, userID uuid.UUID, status string) {}

type RedisClient struct{}
