package validation

import (
	"fmt"
	"strconv"
	"time"
)

// ValidateEventName validates the event name
func ValidateEventName(name string) error {
	if name == "" {
		return fmt.Errorf("event name is required")
	}
	return nil
}

// ValidateEventDescription validates the event description
func ValidateEventDescription(description string) error {
	if description == "" {
		return fmt.Errorf("event description is required")
	}
	return nil
}

// ValidateEventDate validates the event date
func ValidateEventDate(date time.Time) error {
	// Allow events to be created for past dates (for historical records)
	// but warn if the date is more than 1 year in the past
	oneYearAgo := time.Now().AddDate(-1, 0, 0)
	if date.Before(oneYearAgo) {
		return fmt.Errorf("event date is more than 1 year in the past")
	}
	return nil
}

// ValidateEventDateString validates and parses a date string
func ValidateEventDateString(dateStr string) error {
	if dateStr == "" {
		return fmt.Errorf("date is required")
	}
	// Try to parse the date
	_, err := time.Parse("2006-01-02 15:04", dateStr)
	if err != nil {
		return fmt.Errorf("please use format: YYYY-MM-DD HH:MM")
	}
	return nil
}

// ValidateCoordinates validates latitude and longitude
func ValidateCoordinates(latitude, longitude float64) error {
	if latitude < -90 || latitude > 90 {
		return fmt.Errorf("latitude must be between -90 and 90")
	}
	if longitude < -180 || longitude > 180 {
		return fmt.Errorf("longitude must be between -180 and 180")
	}
	return nil
}

// ValidateLatitudeString validates a latitude string input
func ValidateLatitudeString(str string) error {
	if str == "" {
		return nil // Allow empty
	}
	lat, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return fmt.Errorf("invalid number format")
	}
	if lat < -90 || lat > 90 {
		return fmt.Errorf("latitude must be between -90 and 90")
	}
	return nil
}

// ValidateLongitudeString validates a longitude string input
func ValidateLongitudeString(str string) error {
	if str == "" {
		return nil // Allow empty
	}
	lon, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return fmt.Errorf("invalid number format")
	}
	if lon < -180 || lon > 180 {
		return fmt.Errorf("longitude must be between -180 and 180")
	}
	return nil
}

// ValidatePoints validates the default points allocated
func ValidatePoints(points int8) error {
	if points < 0 {
		return fmt.Errorf("points must be a positive number")
	}
	if points > 100 {
		return fmt.Errorf("points cannot exceed 100")
	}
	return nil
}

// ValidatePointsString validates a points string input
func ValidatePointsString(str string) error {
	if str == "" {
		return fmt.Errorf("points are required")
	}
	points, err := strconv.Atoi(str)
	if err != nil {
		return fmt.Errorf("invalid number format")
	}
	if points < 0 || points > 100 {
		return fmt.Errorf("points must be between 0 and 100")
	}
	return nil
}

