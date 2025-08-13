package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/geowa4/rara-membership/handlers"
)

// RequestLogger wraps an http.Handler with request logging
func RequestLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Create a response writer wrapper to capture status code
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		// Call the original handler
		handler.ServeHTTP(rw, r)
		
		// Log the request
		duration := time.Since(start)
		log.Printf("%s %s %d %v %s", 
			r.Method, 
			r.URL.Path, 
			rw.statusCode, 
			duration, 
			r.RemoteAddr,
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func Start() error {
	// Create a new mux for our API routes
	mux := http.NewServeMux()
	
	// Static file serving (no logging for static files)
	mux.Handle("/", http.FileServer(http.Dir("web/")))

	// API endpoints with logging middleware
	// Member endpoints
	mux.Handle("GET /api/members", RequestLogger(http.HandlerFunc(handlers.ListMembers)))
	mux.Handle("POST /api/members", RequestLogger(http.HandlerFunc(handlers.CreateMember)))
	mux.Handle("PUT /api/members/{call_sign}", RequestLogger(http.HandlerFunc(handlers.UpdateMember)))
	mux.Handle("GET /api/members/{call_sign}/points", RequestLogger(http.HandlerFunc(handlers.GetMemberPoints)))

	// Event endpoints
	mux.Handle("GET /api/events", RequestLogger(http.HandlerFunc(handlers.ListEvents)))
	mux.Handle("POST /api/events", RequestLogger(http.HandlerFunc(handlers.CreateEvent)))
	mux.Handle("PUT /api/events/{id}", RequestLogger(http.HandlerFunc(handlers.UpdateEvent)))

	// Point management endpoints
	mux.Handle("POST /api/members/{call_sign}/points", RequestLogger(http.HandlerFunc(handlers.AllocatePoints)))
	mux.Handle("POST /api/members/{call_sign}/redemptions", RequestLogger(http.HandlerFunc(handlers.RedeemPoints)))

	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	fmt.Printf("Server starting on port %s...\n", port)
	fmt.Printf("Web interface: http://localhost:%s\n", port)
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

	return http.ListenAndServe(":"+port, mux)
}
