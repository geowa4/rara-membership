package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/geowa4/rara-membership/services"
	"github.com/geowa4/rara-membership/validation"
	"github.com/spf13/cobra"
)

var updateEventCmd = &cobra.Command{
	Use:   "update [event-id]",
	Short: "Update an existing event",
	Long:  "Update an existing event's information using an interactive form.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Parse event ID
		eventID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid event ID: %w", err)
		}

		// Fetch the existing event
		event, err := services.GetEventByID(eventID)
		if err != nil {
			return fmt.Errorf("failed to fetch event: %w", err)
		}

		// Pre-populate with existing values
		var (
			name             string = event.Name
			description      string = event.Description
			dateStr          string = event.Date.Format("2006-01-02 15:04")
			latitudeStr      string
			longitudeStr     string
			defaultPointsStr string = strconv.Itoa(int(event.DefaultPointsAllocated))
		)

		// Convert existing coordinates to strings if they exist
		if event.Latitude != 0 {
			latitudeStr = fmt.Sprintf("%.6f", event.Latitude)
		}
		if event.Longitude != 0 {
			longitudeStr = fmt.Sprintf("%.6f", event.Longitude)
		}

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Event Name").
					Description("Enter the name of the event").
					Value(&name).
					Validate(validation.ValidateEventName),

				huh.NewText().
					Title("Event Description").
					Description("Enter a description of the event").
					Value(&description).
					Validate(validation.ValidateEventDescription).
					Lines(3),

				huh.NewInput().
					Title("Event Date").
					Description("Enter the date and time of the event").
					Value(&dateStr).
					Validate(validation.ValidateEventDateString),

				huh.NewInput().
					Title("Latitude (optional)").
					Description("Enter the latitude (-90 to 90) or leave empty").
					Value(&latitudeStr).
					Placeholder("43.1234").
					Validate(validation.ValidateLatitudeString),

				huh.NewInput().
					Title("Longitude (optional)").
					Description("Enter the longitude (-180 to 180) or leave empty").
					Value(&longitudeStr).
					Placeholder("-77.5678").
					Validate(validation.ValidateLongitudeString),

				huh.NewInput().
					Title("Default Points").
					Description("Enter the default points allocated for this event (0-100)").
					Value(&defaultPointsStr).
					Validate(validation.ValidatePointsString),
			),
		)

		err = form.Run()
		if err != nil {
			return fmt.Errorf("form error: %w", err)
		}

		// Parse the date
		eventDate, err := time.Parse("2006-01-02 15:04", dateStr)
		if err != nil {
			return fmt.Errorf("failed to parse date: %w", err)
		}

		// Validate the date
		if err := validation.ValidateEventDate(eventDate); err != nil {
			return err
		}

		// Parse coordinates
		var latitude, longitude float64
		if latitudeStr != "" {
			lat, _ := strconv.ParseFloat(latitudeStr, 64)
			latitude = lat
		}
		if longitudeStr != "" {
			lon, _ := strconv.ParseFloat(longitudeStr, 64)
			longitude = lon
		}

		// Validate coordinates if both are provided
		if latitudeStr != "" && longitudeStr != "" {
			if err := validation.ValidateCoordinates(latitude, longitude); err != nil {
				return err
			}
		}


		// Parse points
		defaultPoints, err := strconv.Atoi(defaultPointsStr)
		if err != nil {
			return fmt.Errorf("failed to parse points: %w", err)
		}

		// Validate points
		if err := validation.ValidatePoints(int8(defaultPoints)); err != nil {
			return err
		}

		// Update the event
		updatedEvent, err := services.UpdateEvent(eventID, name, description, eventDate, latitude, longitude, int8(defaultPoints))
		if err != nil {
			return fmt.Errorf("failed to update event: %w", err)
		}

		// Display success message
		fmt.Printf("\n✅ Event updated successfully!\n")
		fmt.Printf("ID: %d\n", updatedEvent.ID)
		fmt.Printf("Name: %s\n", updatedEvent.Name)
		fmt.Printf("Description: %s\n", updatedEvent.Description)
		fmt.Printf("Date: %s\n", updatedEvent.Date.Format("2006-01-02 15:04"))
		if latitude != 0 || longitude != 0 {
			fmt.Printf("Location: %.4f, %.4f\n", updatedEvent.Latitude, updatedEvent.Longitude)
		}
		fmt.Printf("Default Points: %d\n", updatedEvent.DefaultPointsAllocated)

		return nil
	},
}

