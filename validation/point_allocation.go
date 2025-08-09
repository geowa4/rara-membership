package validation

import (
	"fmt"
	"strconv"
	"strings"
)

// ValidateEventID validates an event ID string
func ValidateEventID(eventIDStr string) (int, error) {
	eventIDStr = strings.TrimSpace(eventIDStr)
	if eventIDStr == "" {
		return 0, fmt.Errorf("event ID cannot be empty")
	}

	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		return 0, fmt.Errorf("invalid event ID: %s", eventIDStr)
	}

	if eventID <= 0 {
		return 0, fmt.Errorf("event ID must be positive")
	}

	return eventID, nil
}

// ValidateMemberID validates a member ID string
func ValidateMemberID(memberIDStr string) (int, error) {
	memberIDStr = strings.TrimSpace(memberIDStr)
	if memberIDStr == "" {
		return 0, fmt.Errorf("member ID cannot be empty")
	}

	memberID, err := strconv.Atoi(memberIDStr)
	if err != nil {
		return 0, fmt.Errorf("invalid member ID: %s", memberIDStr)
	}

	if memberID <= 0 {
		return 0, fmt.Errorf("member ID must be positive")
	}

	return memberID, nil
}

// ValidateAllocationPoints validates a points string for allocations
func ValidateAllocationPoints(pointsStr string) (*int, error) {
	pointsStr = strings.TrimSpace(pointsStr)
	if pointsStr == "" {
		// Return nil to use default points
		return nil, nil
	}

	points, err := strconv.Atoi(pointsStr)
	if err != nil {
		return nil, fmt.Errorf("invalid points value: %s", pointsStr)
	}

	if points < 0 {
		return nil, fmt.Errorf("points cannot be negative")
	}

	if points > 100 {
		return nil, fmt.Errorf("points cannot exceed 100")
	}

	return &points, nil
}
