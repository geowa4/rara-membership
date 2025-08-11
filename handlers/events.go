package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/geowa4/rara-membership/database"
	"github.com/geowa4/rara-membership/ent"
	"github.com/geowa4/rara-membership/services"
)

func ListEvents(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	eventService := services.NewEventService(database.Client)

	// Get parameters from query string
	filter := r.URL.Query().Get("filter")
	perPageStr := r.URL.Query().Get("perPage")
	beforeIDStr := r.URL.Query().Get("beforeID")

	// Parse pagination parameters
	perPage := 50 // Default
	if perPageStr != "" {
		if parsed, err := strconv.Atoi(perPageStr); err == nil && parsed > 0 {
			perPage = parsed
		}
	}

	var beforeID *int
	if beforeIDStr != "" {
		if parsed, err := strconv.Atoi(beforeIDStr); err == nil {
			beforeID = &parsed
		}
	}

	var events []*ent.Event
	var err error

	switch filter {
	case "past":
		events, err = eventService.ListPastEvents(ctx)
	case "upcoming":
		events, err = eventService.ListUpcomingEvents(ctx)
	default:
		// Use pagination for default case (all events)
		if beforeID != nil || perPage != 50 {
			events, err = eventService.ListEventsPaginated(ctx, perPage, beforeID)
		} else {
			events, err = eventService.ListEvents(ctx)
		}
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to query events: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(events); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Date        *string  `json:"date"`
		Timezone    *string  `json:"timezone"`
		Points      int      `json:"points"`
		Latitude    *float64 `json:"latitude"`
		Longitude   *float64 `json:"longitude"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Parse date if provided
	var eventDate *time.Time
	if input.Date != nil && *input.Date != "" {
		// Try datetime-local format first (YYYY-MM-DDTHH:mm), then date format (YYYY-MM-DD)
		if parsed, err := time.Parse("2006-01-02T15:04", *input.Date); err == nil {
			eventDate = &parsed
		} else if parsed, err := time.Parse("2006-01-02", *input.Date); err == nil {
			eventDate = &parsed
		}
	}

	ctx := context.Background()
	eventService := services.NewEventService(database.Client)

	event, err := eventService.CreateEvent(ctx, input.Name, input.Description, eventDate, input.Timezone, input.Latitude, input.Longitude, input.Points)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create event: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(event); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

func UpdateEvent(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.PathValue("id")
	if eventIDStr == "" {
		http.Error(w, "Event ID is required", http.StatusBadRequest)
		return
	}

	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	var input struct {
		Name        *string  `json:"name"`
		Description *string  `json:"description"`
		Date        *string  `json:"date"`
		Timezone    *string  `json:"timezone"`
		Points      *int     `json:"points"`
		Latitude    *float64 `json:"latitude"`
		Longitude   *float64 `json:"longitude"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Parse date if provided
	var eventDate *time.Time
	if input.Date != nil && *input.Date != "" {
		// Try datetime-local format first (YYYY-MM-DDTHH:mm), then date format (YYYY-MM-DD)
		if parsed, err := time.Parse("2006-01-02T15:04", *input.Date); err == nil {
			eventDate = &parsed
		} else if parsed, err := time.Parse("2006-01-02", *input.Date); err == nil {
			eventDate = &parsed
		}
	}

	ctx := context.Background()
	eventService := services.NewEventService(database.Client)

	event, err := eventService.UpdateEvent(ctx, eventID, input.Name, input.Description, eventDate, input.Timezone, input.Latitude, input.Longitude, input.Points)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update event: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(event); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}
