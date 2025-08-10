package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/geowa4/rara-membership/database"
	"github.com/geowa4/rara-membership/services"
)

type PointHistoryResponse struct {
	Member         interface{}         `json:"member"`
	CurrentBalance int                 `json:"current_balance"`
	TotalEarned    int                 `json:"total_earned"`
	TotalRedeemed  int                 `json:"total_redeemed"`
	Transactions   []TransactionDetail `json:"transactions"`
}

type TransactionDetail struct {
	Type        string `json:"type"`
	Points      int    `json:"points"`
	Date        string `json:"date"`
	Description string `json:"description"`
	Notes       string `json:"notes"`
}

func GetMemberPoints(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	callSign := r.URL.Query().Get("call_sign")
	if callSign == "" {
		http.Error(w, "Call sign is required", http.StatusBadRequest)
		return
	}

	callSign = strings.ToUpper(strings.TrimSpace(callSign))
	ctx := context.Background()

	// Find the member by call sign
	mbr, err := services.GetActiveMemberByCallSign(callSign)
	if err != nil {
		http.Error(w, fmt.Sprintf("Member with call sign %s not found or is not active", callSign), http.StatusNotFound)
		return
	}

	// Create services
	allocationService := services.NewPointAllocationService(database.Client)
	deductionService := services.NewPointDeductionService(database.Client)

	// Get allocations (points earned)
	allocations, err := allocationService.GetAllocationsByMember(ctx, mbr.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting point allocations: %v", err), http.StatusInternalServerError)
		return
	}

	// Get deductions (points redeemed)
	deductions, err := deductionService.GetDeductionsByMember(ctx, mbr.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting point deductions: %v", err), http.StatusInternalServerError)
		return
	}

	// Build response
	var transactions []TransactionDetail
	totalEarned := 0
	totalRedeemed := 0

	for _, allocation := range allocations {
		totalEarned += allocation.Points
		transactions = append(transactions, TransactionDetail{
			Type:        "EARNED",
			Points:      allocation.Points,
			Date:        allocation.CreatedAt.Format("2006-01-02T15:04:05Z"),
			Description: allocation.Edges.Event.Name,
			Notes:       allocation.Notes,
		})
	}

	for _, deduction := range deductions {
		totalRedeemed += deduction.Points
		transactions = append(transactions, TransactionDetail{
			Type:        "REDEEMED",
			Points:      -deduction.Points,
			Date:        deduction.CreatedAt.Format("2006-01-02T15:04:05Z"),
			Description: "Points Redeemed",
			Notes:       deduction.Notes,
		})
	}

	response := PointHistoryResponse{
		Member:         mbr,
		CurrentBalance: totalEarned - totalRedeemed,
		TotalEarned:    totalEarned,
		TotalRedeemed:  totalRedeemed,
		Transactions:   transactions,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

func AllocatePoints(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		EventID  int    `json:"event_id"`
		MemberID int    `json:"member_id"`
		Points   *int   `json:"points"`
		Notes    string `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	if input.EventID <= 0 || input.MemberID <= 0 {
		http.Error(w, "Event ID and Member ID must be positive numbers", http.StatusBadRequest)
		return
	}

	if input.Points != nil && (*input.Points < 0 || *input.Points > 100) {
		http.Error(w, "Points must be between 0 and 100", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	pointService := services.NewPointAllocationService(database.Client)

	allocation, err := pointService.AllocatePoints(ctx, input.EventID, input.MemberID, input.Points, input.Notes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error allocating points: %v", err), http.StatusInternalServerError)
		return
	}

	// Get total points for the member
	totalPoints, _ := pointService.GetTotalPointsForMember(ctx, input.MemberID)

	response := map[string]interface{}{
		"allocation":   allocation,
		"total_points": totalPoints,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

func RedeemPoints(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		MemberID int    `json:"member_id"`
		Points   int    `json:"points"`
		Notes    string `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	if input.MemberID <= 0 {
		http.Error(w, "Member ID must be a positive number", http.StatusBadRequest)
		return
	}

	if input.Points <= 0 {
		http.Error(w, "Points must be a positive number", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	deductionService := services.NewPointDeductionService(database.Client)

	// Check current balance before attempting redemption
	balance, err := deductionService.GetPointBalanceForMember(ctx, input.MemberID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error checking member point balance: %v", err), http.StatusInternalServerError)
		return
	}

	if balance < input.Points {
		http.Error(w, fmt.Sprintf("Insufficient points: member has %d points but trying to redeem %d", balance, input.Points), http.StatusBadRequest)
		return
	}

	// Redeem points
	deduction, err := deductionService.DeductPoints(ctx, input.MemberID, input.Points, input.Notes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error redeeming points: %v", err), http.StatusInternalServerError)
		return
	}

	// Get new balance
	newBalance, _ := deductionService.GetPointBalanceForMember(ctx, input.MemberID)

	response := map[string]interface{}{
		"deduction":        deduction,
		"previous_balance": balance,
		"new_balance":      newBalance,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

func GetMemberPointBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	memberIDStr := r.URL.Query().Get("member_id")
	if memberIDStr == "" {
		http.Error(w, "Member ID is required", http.StatusBadRequest)
		return
	}

	memberID, err := strconv.Atoi(memberIDStr)
	if err != nil || memberID <= 0 {
		http.Error(w, "Invalid member ID", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	deductionService := services.NewPointDeductionService(database.Client)

	balance, err := deductionService.GetPointBalanceForMember(ctx, memberID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting member point balance: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"member_id": memberID,
		"balance":   balance,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}
