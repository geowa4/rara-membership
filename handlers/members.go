package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/geowa4/rara-membership/database"
	"github.com/geowa4/rara-membership/ent"
	"github.com/geowa4/rara-membership/services"
)

func ListMembers(w http.ResponseWriter, r *http.Request) {

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

func CreateMember(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name           string  `json:"name"`
		CallSign       string  `json:"call_sign"`
		Email          string  `json:"email"`
		Phone          *string `json:"phone"`
		MailingAddress *string `json:"mailing_address"`
		FRN            *string `json:"frn"`
		LicenseClass   *string `json:"license_class"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	if input.Name == "" || input.CallSign == "" || input.Email == "" {
		http.Error(w, "Name, email, and call sign are required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	memberService := services.NewMemberService(database.Client)

	memberInput := services.CreateMemberInput{
		Name:           input.Name,
		CallSign:       strings.ToUpper(input.CallSign),
		Email:          input.Email,
		Phone:          input.Phone,
		MailingAddress: input.MailingAddress,
		FRN:            input.FRN,
		LicenseClass:   input.LicenseClass,
	}

	member, err := memberService.CreateMember(ctx, memberInput)
	if err != nil {
		if ent.IsConstraintError(err) {
			http.Error(w, "A member with this call sign already exists", http.StatusConflict)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to create member: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(member); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

func UpdateMember(w http.ResponseWriter, r *http.Request) {
	callSign := r.PathValue("call_sign")
	if callSign == "" {
		http.Error(w, "Call sign is required", http.StatusBadRequest)
		return
	}

	var input struct {
		Name           *string `json:"name"`
		CallSign       *string `json:"call_sign"`
		Email          *string `json:"email"`
		Phone          *string `json:"phone"`
		MailingAddress *string `json:"mailing_address"`
		FRN            *string `json:"frn"`
		LicenseClass   *string `json:"license_class"`
		IsActive       *bool   `json:"is_active"`
		IsSilentKey    *bool   `json:"is_silent_key"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	memberService := services.NewMemberService(database.Client)

	updateInput := services.UpdateMemberInput{
		Name:           input.Name,
		CallSign:       input.CallSign,
		Email:          input.Email,
		Phone:          input.Phone,
		MailingAddress: input.MailingAddress,
		FRN:            input.FRN,
		LicenseClass:   input.LicenseClass,
		IsActive:       input.IsActive,
		IsSilentKey:    input.IsSilentKey,
	}

	member, err := memberService.UpdateMemberByCallSign(ctx, strings.ToUpper(callSign), updateInput)
	if err != nil {
		if ent.IsNotFound(err) {
			http.Error(w, "Member not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to update member: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(member); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

// Deprecated: Use ListMembers instead
func Members(w http.ResponseWriter, r *http.Request) {
	ListMembers(w, r)
}
