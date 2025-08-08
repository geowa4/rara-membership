package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/geowa4/rara-membership/ent"
	"github.com/geowa4/rara-membership/services"
	"github.com/spf13/cobra"
)

var (
	showPast     bool
	showUpcoming bool
)

var listEventsCmd = &cobra.Command{
	Use:   "list",
	Short: "List events",
	Long:  "Display a list of events. By default shows all events. Use flags to filter.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var events []*ent.Event
		var err error
		var title string

		// Determine which events to fetch based on flags
		if showUpcoming && showPast {
			// Both flags set, show all events
			events, err = services.GetAllEvents()
			title = "All Events"
		} else if showUpcoming {
			events, err = services.GetUpcomingEvents()
			title = "Upcoming Events"
		} else if showPast {
			events, err = services.GetPastEvents()
			title = "Past Events"
		} else {
			// No flags set, show all events
			events, err = services.GetAllEvents()
			title = "All Events"
		}

		if err != nil {
			return fmt.Errorf("failed querying events: %w", err)
		}

		if len(events) == 0 {
			fmt.Printf("No %s found.\n", title)
			return nil
		}

		// Define colors and styles
		var (
			purple    = lipgloss.Color("99")
			gray      = lipgloss.Color("245")
			lightGray = lipgloss.Color("241")
			green     = lipgloss.Color("42")

			headerStyle  = lipgloss.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
			cellStyle    = lipgloss.NewStyle().Padding(0, 1)
			oddRowStyle  = cellStyle.Foreground(gray)
			evenRowStyle = cellStyle.Foreground(lightGray)
		)

		// Prepare table data
		rows := make([][]string, len(events))
		now := time.Now()

		for i, e := range events {
			// Format date
			dateStr := e.Date.Format("2006-01-02 15:04")

			// Determine if event is upcoming or past
			status := ""
			if e.Date.After(now) {
				status = "Upcoming"
			} else {
				status = "Past"
			}

			// Format coordinates
			coords := fmt.Sprintf("%.4f, %.4f", e.Latitude, e.Longitude)
			if e.Latitude == 0 && e.Longitude == 0 {
				coords = "-"
			}

			rows[i] = []string{
				strconv.Itoa(e.ID),
				e.Name,
				e.Description,
				dateStr,
				coords,
				strconv.Itoa(int(e.DefaultPointsAllocated)),
				status,
			}
		}

		// Create and style the table
		t := table.New().
			Border(lipgloss.NormalBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(purple)).
			StyleFunc(func(row, col int) lipgloss.Style {
				switch {
				case row == table.HeaderRow:
					return headerStyle
				case row%2 == 0:
					return evenRowStyle
				default:
					return oddRowStyle
				}
			}).
			Headers("ID", "NAME", "DESCRIPTION", "DATE", "LOCATION", "POINTS", "STATUS").
			Rows(rows...)

		// Print title and table
		titleStyle := lipgloss.NewStyle().
			Foreground(green).
			Bold(true).
			MarginBottom(1)

		// Add emoji based on filter
		emoji := "📅 "
		if showPast && !showUpcoming {
			emoji = "📚 "
		} else if showUpcoming && !showPast {
			emoji = "🔮 "
		}

		fmt.Println(titleStyle.Render(fmt.Sprintf("%s %s (%d)", emoji, title, len(events))))
		fmt.Println(t)

		return nil
	},
}

func init() {
	listEventsCmd.Flags().BoolVar(&showPast, "past", false, "Show only past events")
	listEventsCmd.Flags().BoolVar(&showUpcoming, "upcoming", false, "Show only upcoming events")
}
