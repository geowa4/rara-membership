package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/geowa4/rara-membership/handlers"
)

func Start() error {
	http.HandleFunc("/api/members", handlers.Members)

	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	fmt.Printf("Server starting on port %s...\n", port)
	fmt.Printf("API endpoint: http://localhost:%s/api/members\n", port)
	return http.ListenAndServe(":"+port, nil)
}
