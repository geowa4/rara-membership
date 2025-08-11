package services

import (
	"context"
	"fmt"
	"time"

	"github.com/geowa4/rara-membership/database"
	"github.com/geowa4/rara-membership/ent"
	"github.com/geowa4/rara-membership/ent/event"
)

type EventService struct {
	client *ent.Client
}

func NewEventService(client *ent.Client) *EventService {
	return &EventService{client: client}
}

// GetAllEvents returns all events ordered by date
func GetAllEvents() ([]*ent.Event, error) {
	ctx := context.Background()
	return database.Client.Event.Query().Order(ent.Desc(event.FieldDate)).All(ctx)
}

// GetUpcomingEvents returns events happening after the current date
func GetUpcomingEvents() ([]*ent.Event, error) {
	ctx := context.Background()
	return database.Client.Event.Query().
		Where(event.DateGTE(time.Now())).
		Order(ent.Asc(event.FieldDate)).
		All(ctx)
}

// GetPastEvents returns events that have already happened
func GetPastEvents() ([]*ent.Event, error) {
	ctx := context.Background()
	return database.Client.Event.Query().
		Where(event.DateLT(time.Now())).
		Order(ent.Desc(event.FieldDate)).
		All(ctx)
}

// CreateEvent creates a new event
func CreateEvent(name, description string, date time.Time, latitude, longitude float64, defaultPoints int8) (*ent.Event, error) {
	ctx := context.Background()
	return database.Client.Event.Create().
		SetName(name).
		SetDescription(description).
		SetDate(date).
		SetLatitude(latitude).
		SetLongitude(longitude).
		SetDefaultPointsAllocated(defaultPoints).
		Save(ctx)
}

// GetEventByID retrieves an event by its ID
func GetEventByID(id int) (*ent.Event, error) {
	ctx := context.Background()
	return database.Client.Event.Get(ctx, id)
}

// UpdateEvent updates an existing event
func UpdateEvent(id int, name, description string, date time.Time, latitude, longitude float64, defaultPoints int8) (*ent.Event, error) {
	ctx := context.Background()
	return database.Client.Event.UpdateOneID(id).
		SetName(name).
		SetDescription(description).
		SetDate(date).
		SetLatitude(latitude).
		SetLongitude(longitude).
		SetDefaultPointsAllocated(defaultPoints).
		Save(ctx)
}

// DeleteEvent deletes an event by ID
func DeleteEvent(id int) error {
	ctx := context.Background()
	return database.Client.Event.DeleteOneID(id).Exec(ctx)
}

// ListEvents returns events with pagination support, defaulting to last 50 events ordered by date descending
func (s *EventService) ListEvents(ctx context.Context) ([]*ent.Event, error) {
	return s.client.Event.Query().
		Order(ent.Desc(event.FieldDate)).
		Limit(50).
		All(ctx)
}

// ListEventsPaginated returns events with pagination support
func (s *EventService) ListEventsPaginated(ctx context.Context, perPage int, beforeID *int) ([]*ent.Event, error) {
	query := s.client.Event.Query().Order(ent.Desc(event.FieldDate))
	
	// Apply beforeID filter for cursor-based pagination
	if beforeID != nil {
		// Get the date of the beforeID event for proper cursor pagination
		beforeEvent, err := s.client.Event.Get(ctx, *beforeID)
		if err != nil {
			return nil, fmt.Errorf("failed to get before event: %w", err)
		}
		// Return events with date less than the before event's date, or same date but higher ID
		query = query.Where(
			event.Or(
				event.DateLT(beforeEvent.Date),
				event.And(
					event.DateEQ(beforeEvent.Date),
					event.IDLT(*beforeID),
				),
			),
		)
	}
	
	// Apply pagination limit, with sensible defaults
	if perPage <= 0 || perPage > 1000 {
		perPage = 50 // Default to 50, max 1000
	}
	
	return query.Limit(perPage).All(ctx)
}

// ListUpcomingEvents returns events happening after the current date
func (s *EventService) ListUpcomingEvents(ctx context.Context) ([]*ent.Event, error) {
	return s.client.Event.Query().
		Where(event.DateGTE(time.Now())).
		Order(ent.Desc(event.FieldDate)).
		All(ctx)
}

// ListPastEvents returns events that have already happened
func (s *EventService) ListPastEvents(ctx context.Context) ([]*ent.Event, error) {
	return s.client.Event.Query().
		Where(event.DateLT(time.Now())).
		Order(ent.Desc(event.FieldDate)).
		All(ctx)
}

// CreateEvent creates a new event with comprehensive parameters
func (s *EventService) CreateEvent(ctx context.Context, name, description string, date *time.Time, timezone *string, latitude, longitude *float64, defaultPoints int) (*ent.Event, error) {
	create := s.client.Event.Create().
		SetName(name).
		SetDescription(description).
		SetDefaultPointsAllocated(int8(defaultPoints))
	
	if date != nil {
		create = create.SetDate(*date)
	} else {
		create = create.SetDate(time.Now())
	}
	
	if timezone != nil && *timezone != "" {
		create = create.SetTimezone(*timezone)
	}
	
	if latitude != nil {
		create = create.SetLatitude(*latitude)
	}
	
	if longitude != nil {
		create = create.SetLongitude(*longitude)
	}
	
	return create.Save(ctx)
}

// UpdateEvent updates an existing event with comprehensive parameters
func (s *EventService) UpdateEvent(ctx context.Context, id int, name, description *string, date *time.Time, timezone *string, latitude, longitude *float64, defaultPoints *int) (*ent.Event, error) {
	update := s.client.Event.UpdateOneID(id)

	if name != nil {
		update = update.SetName(*name)
	}

	if description != nil {
		update = update.SetDescription(*description)
	}

	if date != nil {
		update = update.SetDate(*date)
	}

	if timezone != nil && *timezone != "" {
		update = update.SetTimezone(*timezone)
	}

	if latitude != nil {
		update = update.SetLatitude(*latitude)
	}

	if longitude != nil {
		update = update.SetLongitude(*longitude)
	}

	if defaultPoints != nil {
		update = update.SetDefaultPointsAllocated(int8(*defaultPoints))
	}

	return update.Save(ctx)
}
