package services

import (
	"context"
	"fmt"

	"github.com/geowa4/rara-membership/ent"
	"github.com/geowa4/rara-membership/ent/event"
	"github.com/geowa4/rara-membership/ent/member"
	"github.com/geowa4/rara-membership/ent/pointallocation"
)

// PointAllocationService handles business logic for point allocations
type PointAllocationService struct {
	client *ent.Client
}

// NewPointAllocationService creates a new point allocation service
func NewPointAllocationService(client *ent.Client) *PointAllocationService {
	return &PointAllocationService{client: client}
}

// AllocatePoints allocates points to a member for an event
func (s *PointAllocationService) AllocatePoints(ctx context.Context, eventID, memberID int, points *int, notes string) (*ent.PointAllocation, error) {
	// First, check if the event exists
	evt, err := s.client.Event.Get(ctx, eventID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("event with ID %d not found", eventID)
		}
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

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

	// Use default points if not specified
	pointsToAllocate := int(evt.DefaultPointsAllocated)
	if points != nil {
		pointsToAllocate = *points
	}

	// Check if allocation already exists
	exists, err := s.client.PointAllocation.Query().
		Where(
			pointallocation.HasEventWith(event.ID(eventID)),
			pointallocation.HasMemberWith(member.ID(memberID)),
		).
		Exist(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing allocation: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("point allocation already exists for member ID %d and event ID %d", memberID, eventID)
	}

	// Create the point allocation
	builder := s.client.PointAllocation.Create().
		SetEventID(eventID).
		SetMemberID(memberID).
		SetPoints(pointsToAllocate)

	if notes != "" {
		builder.SetNotes(notes)
	}

	allocation, err := builder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create point allocation: %w", err)
	}

	// Load the edges for the response
	allocation, err = s.client.PointAllocation.Query().
		Where(pointallocation.ID(allocation.ID)).
		WithEvent().
		WithMember().
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load point allocation with edges: %w", err)
	}

	return allocation, nil
}

// GetAllocationsByEvent returns all point allocations for an event
func (s *PointAllocationService) GetAllocationsByEvent(ctx context.Context, eventID int) ([]*ent.PointAllocation, error) {
	allocations, err := s.client.PointAllocation.Query().
		Where(pointallocation.HasEventWith(event.ID(eventID))).
		WithMember().
		WithEvent().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get allocations for event: %w", err)
	}
	return allocations, nil
}

// GetAllocationsByMember returns all point allocations for a member
func (s *PointAllocationService) GetAllocationsByMember(ctx context.Context, memberID int) ([]*ent.PointAllocation, error) {
	allocations, err := s.client.PointAllocation.Query().
		Where(pointallocation.HasMemberWith(member.ID(memberID))).
		WithEvent().
		WithMember().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get allocations for member: %w", err)
	}
	return allocations, nil
}

// GetTotalPointsForMember returns the total points earned by a member
func (s *PointAllocationService) GetTotalPointsForMember(ctx context.Context, memberID int) (int, error) {
	allocations, err := s.client.PointAllocation.Query().
		Where(pointallocation.HasMemberWith(member.ID(memberID))).
		All(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get allocations for member: %w", err)
	}

	total := 0
	for _, allocation := range allocations {
		total += allocation.Points
	}
	return total, nil
}

// GetVolunteerHistoryForMember returns all volunteer history for a member, sorted by event date (newest first)
func (s *PointAllocationService) GetVolunteerHistoryForMember(ctx context.Context, memberID int) ([]*ent.PointAllocation, error) {
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

	allocations, err := s.client.PointAllocation.Query().
		Where(pointallocation.HasMemberWith(member.ID(memberID))).
		WithEvent().
		WithMember().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get volunteer history for member: %w", err)
	}

	// Sort by event date (newest first)
	// Since we can't easily sort by event.date in the query, we'll sort in memory
	for i := 0; i < len(allocations)-1; i++ {
		for j := i + 1; j < len(allocations); j++ {
			if allocations[i].Edges.Event.Date.Before(allocations[j].Edges.Event.Date) {
				allocations[i], allocations[j] = allocations[j], allocations[i]
			}
		}
	}

	return allocations, nil
}