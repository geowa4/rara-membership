package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/geowa4/rara-membership/handlers"
)

func Start() error {
	// Member endpoints
	http.HandleFunc("/api/members", handleMembersEndpoint)
	http.HandleFunc("/api/members/create", handlers.CreateMember)
	http.HandleFunc("/api/members/update", handlers.UpdateMember)
	http.HandleFunc("/api/members/points", handlers.GetMemberPoints)

	// Event endpoints
	http.HandleFunc("/api/events", handlers.ListEvents)
	http.HandleFunc("/api/events/create", handlers.CreateEvent)
	http.HandleFunc("/api/events/update", handlers.UpdateEvent)

	// Point management endpoints
	http.HandleFunc("/api/points/allocate", handlers.AllocatePoints)
	http.HandleFunc("/api/points/redeem", handlers.RedeemPoints)
	http.HandleFunc("/api/points/balance", handlers.GetMemberPointBalance)

	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	fmt.Printf("Server starting on port %s...\n", port)
	fmt.Printf("\nAvailable API endpoints:\n")
	fmt.Printf("  Members:\n")
	fmt.Printf("    GET    /api/members           - List all active members\n")
	fmt.Printf("    POST   /api/members/create    - Create new member\n")
	fmt.Printf("    PUT    /api/members/update    - Update member (requires ?call_sign=)\n")
	fmt.Printf("    GET    /api/members/points    - Get member point history (requires ?call_sign=)\n")
	fmt.Printf("  Events:\n")
	fmt.Printf("    GET    /api/events            - List all events\n")
	fmt.Printf("    POST   /api/events/create     - Create new event\n")
	fmt.Printf("    PUT    /api/events/update     - Update event (requires ?id=)\n")
	fmt.Printf("  Points:\n")
	fmt.Printf("    POST   /api/points/allocate   - Allocate points to member for event\n")
	fmt.Printf("    POST   /api/points/redeem     - Redeem member points\n")
	fmt.Printf("    GET    /api/points/balance    - Get member point balance (requires ?member_id=)\n")
	fmt.Printf("\n")

	return http.ListenAndServe(":"+port, nil)
}

// handleMembersEndpoint routes based on HTTP method for backwards compatibility
func handleMembersEndpoint(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handlers.ListMembers(w, r)
	case http.MethodPost:
		handlers.CreateMember(w, r)
	case http.MethodPut, http.MethodPatch:
		handlers.UpdateMember(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
