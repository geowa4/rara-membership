package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/geowa4/rara-membership/database"
	"github.com/geowa4/rara-membership/services"
)

func ListEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := context.Background()
	eventService := services.NewEventService(database.Client)

	events, err := eventService.ListEvents(ctx)
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
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Name          string `json:"name"`
		Description   string `json:"description"`
		DefaultPoints int    `json:"default_points"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	eventService := services.NewEventService(database.Client)

	event, err := eventService.CreateEvent(ctx, input.Name, input.Description, input.DefaultPoints)
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
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	eventIDStr := r.URL.Query().Get("id")
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
		Name          *string `json:"name"`
		Description   *string `json:"description"`
		DefaultPoints *int    `json:"default_points"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	eventService := services.NewEventService(database.Client)

	event, err := eventService.UpdateEvent(ctx, eventID, input.Name, input.Description, input.DefaultPoints)
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
