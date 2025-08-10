package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/geowa4/rara-membership/handlers"
)

func Start() error {
	// Member endpoints
	http.HandleFunc("GET /api/members", handlers.ListMembers)
	http.HandleFunc("POST /api/members", handlers.CreateMember)
	http.HandleFunc("PUT /api/members/{call_sign}", handlers.UpdateMember)
	http.HandleFunc("GET /api/members/{call_sign}/points", handlers.GetMemberPoints)

	// Event endpoints
	http.HandleFunc("GET /api/events", handlers.ListEvents)
	http.HandleFunc("POST /api/events", handlers.CreateEvent)
	http.HandleFunc("PUT /api/events/{id}", handlers.UpdateEvent)

	// Point management endpoints
	http.HandleFunc("POST /api/members/{call_sign}/points", handlers.AllocatePoints)
	http.HandleFunc("POST /api/members/{call_sign}/redemptions", handlers.RedeemPoints)

	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	fmt.Printf("Server starting on port %s...\n", port)
	fmt.Printf("\nAvailable API endpoints:\n")
	fmt.Printf("  Members:\n")
	fmt.Printf("    GET    /api/members                       - List all active members\n")
	fmt.Printf("    POST   /api/members                       - Create new member\n")
	fmt.Printf("    PUT    /api/members/{call_sign}           - Update member by call sign\n")
	fmt.Printf("    GET    /api/members/{call_sign}/points    - Get member point history and balance\n")
	fmt.Printf("  Events:\n")
	fmt.Printf("    GET    /api/events                        - List all events\n")
	fmt.Printf("    POST   /api/events                        - Create new event\n")
	fmt.Printf("    PUT    /api/events/{id}                   - Update event by ID\n")
	fmt.Printf("  Points:\n")
	fmt.Printf("    POST   /api/members/{call_sign}/points        - Allocate points to member for event\n")
	fmt.Printf("    POST   /api/members/{call_sign}/redemptions   - Redeem member points\n")
	fmt.Printf("\n")

	return http.ListenAndServe(":"+port, nil)
}
