package services

import (
	"context"
	"time"

	"github.com/geowa4/rara-membership/database"
	"github.com/geowa4/rara-membership/ent"
	"github.com/geowa4/rara-membership/ent/event"
)

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

