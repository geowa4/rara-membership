package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/geowa4/rara-membership/services"
)

func Members(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	members, err := services.GetActiveMembers()
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to query members: %v", err),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(members); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to encode response: %v", err),
			http.StatusInternalServerError,
		)
		return
	}
}
