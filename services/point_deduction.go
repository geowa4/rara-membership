package services

import (
	"context"
	"fmt"

	"github.com/geowa4/rara-membership/ent"
	"github.com/geowa4/rara-membership/ent/member"
	"github.com/geowa4/rara-membership/ent/pointallocation"
	"github.com/geowa4/rara-membership/ent/pointdeduction"
)

// PointDeductionService handles business logic for point deductions/redemptions
type PointDeductionService struct {
	client *ent.Client
}

// NewPointDeductionService creates a new point deduction service
func NewPointDeductionService(client *ent.Client) *PointDeductionService {
	return &PointDeductionService{client: client}
}

// DeductPoints creates a point deduction for a member, ensuring they have sufficient points
func (s *PointDeductionService) DeductPoints(ctx context.Context, memberID int, points int, notes string) (*ent.PointDeduction, error) {
	// Check if the member exists
	mbr, err := s.client.Member.Get(ctx, memberID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("member with ID %d not found", memberID)
		}
		return nil, fmt.Errorf("failed to get member: %w", err)
	}

	// Check if member is active
	if !mbr.IsActive {
		return nil, fmt.Errorf("member %s is not active", mbr.Name)
	}

	// Calculate current point balance
	balance, err := s.GetPointBalanceForMember(ctx, memberID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate point balance: %w", err)
	}

	// Check if member has enough points
	if balance < points {
		return nil, fmt.Errorf("insufficient points: member has %d points but trying to redeem %d", balance, points)
	}

	// Create the point deduction
	builder := s.client.PointDeduction.Create().
		SetMemberID(memberID).
		SetPoints(points)

	if notes != "" {
		builder.SetNotes(notes)
	}

	deduction, err := builder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create point deduction: %w", err)
	}

	// Load the edges for the response
	deduction, err = s.client.PointDeduction.Query().
		Where(pointdeduction.ID(deduction.ID)).
		WithMember().
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load point deduction with edges: %w", err)
	}

	return deduction, nil
}

// GetDeductionsByMember returns all point deductions for a member
func (s *PointDeductionService) GetDeductionsByMember(ctx context.Context, memberID int) ([]*ent.PointDeduction, error) {
	deductions, err := s.client.PointDeduction.Query().
		Where(pointdeduction.HasMemberWith(member.ID(memberID))).
		Order(ent.Desc(pointdeduction.FieldCreatedAt)).
		WithMember().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get deductions for member: %w", err)
	}
	return deductions, nil
}

// GetTotalDeductionsForMember returns the total points deducted for a member
func (s *PointDeductionService) GetTotalDeductionsForMember(ctx context.Context, memberID int) (int, error) {
	deductions, err := s.client.PointDeduction.Query().
		Where(pointdeduction.HasMemberWith(member.ID(memberID))).
		All(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get deductions for member: %w", err)
	}

	total := 0
	for _, deduction := range deductions {
		total += deduction.Points
	}
	return total, nil
}

// GetPointBalanceForMember returns the current point balance (allocations - deductions)
func (s *PointDeductionService) GetPointBalanceForMember(ctx context.Context, memberID int) (int, error) {
	// Get total allocated points
	allocations, err := s.client.PointAllocation.Query().
		Where(pointallocation.HasMemberWith(member.ID(memberID))).
		All(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get allocations for member: %w", err)
	}

	totalAllocated := 0
	for _, allocation := range allocations {
		totalAllocated += allocation.Points
	}

	// Get total deducted points
	deductions, err := s.client.PointDeduction.Query().
		Where(pointdeduction.HasMemberWith(member.ID(memberID))).
		All(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get deductions for member: %w", err)
	}

	totalDeducted := 0
	for _, deduction := range deductions {
		totalDeducted += deduction.Points
	}

	return totalAllocated - totalDeducted, nil
}

// GetRedemptionHistoryForMember returns all redemption history for a member, sorted by created date (newest first)
func (s *PointDeductionService) GetRedemptionHistoryForMember(ctx context.Context, memberID int) ([]*ent.PointDeduction, error) {
	// First verify the member exists
	exists, err := s.client.Member.Query().
		Where(member.ID(memberID)).
		Exist(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check if member exists: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("member with ID %d not found", memberID)
	}

	deductions, err := s.client.PointDeduction.Query().
		Where(pointdeduction.HasMemberWith(member.ID(memberID))).
		Order(ent.Desc(pointdeduction.FieldCreatedAt)).
		WithMember().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get redemption history for member: %w", err)
	}

	return deductions, nil
}
