package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/geowa4/rara-membership/database"
	"github.com/geowa4/rara-membership/services"
	"github.com/geowa4/rara-membership/validation"
	"github.com/spf13/cobra"
)

var createEventCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new event",
	Long:  "Create a new event using an interactive form.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			name             string
			description      string
			dateStr          string
			timezone         string = "America/New_York"
			latitudeStr      string
			longitudeStr     string
			defaultPointsStr string = "10"
		)

		// Get current date as default
		defaultDate := time.Now().AddDate(0, 0, 7).Format("2006-01-02 15:04")

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
					Placeholder(defaultDate).
					Validate(validation.ValidateEventDateString),

				huh.NewSelect[string]().
					Title("Time Zone").
					Description("Select the timezone for this event").
					Value(&timezone).
					Options(
						huh.NewOption("Eastern Time", "America/New_York"),
						huh.NewOption("Central Time", "America/Chicago"),
						huh.NewOption("Mountain Time", "America/Denver"),
						huh.NewOption("Pacific Time", "America/Los_Angeles"),
						huh.NewOption("Arizona Time", "America/Phoenix"),
						huh.NewOption("Alaska Time", "America/Anchorage"),
						huh.NewOption("Hawaii Time", "Pacific/Honolulu"),
						huh.NewOption("UTC", "UTC"),
					),

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
					Placeholder("10").
					Validate(validation.ValidatePointsString),
			),
		)

		err := form.Run()
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

		// Validate coordinates if provided
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

		// Create the event using the service
		ctx := cmd.Context()
		eventService := services.NewEventService(database.Client)
		event, err := eventService.CreateEvent(ctx, name, description, &eventDate, &timezone, &latitude, &longitude, defaultPoints)
		if err != nil {
			return fmt.Errorf("failed to create event: %w", err)
		}

		// Display success message
		fmt.Printf("\n✅ Event created successfully!\n")
		fmt.Printf("ID: %d\n", event.ID)
		fmt.Printf("Name: %s\n", event.Name)
		fmt.Printf("Description: %s\n", event.Description)
		fmt.Printf("Date: %s\n", event.Date.Format("2006-01-02 15:04"))
		fmt.Printf("Time Zone: %s\n", event.Timezone)
		if latitude != 0 || longitude != 0 {
			fmt.Printf("Location: %.4f, %.4f\n", event.Latitude, event.Longitude)
		}
		fmt.Printf("Default Points: %d\n", event.DefaultPointsAllocated)

		return nil
	},
}
