package services

import (
	"context"
	"fmt"

	"github.com/r-scheele/zero/ent"
	"github.com/r-scheele/zero/ent/user"
)

// AdminService handles admin-specific operations
type AdminService struct {
	orm *ent.Client
}

// NewAdminService creates a new admin service
func NewAdminService(orm *ent.Client) *AdminService {
	return &AdminService{
		orm: orm,
	}
}

// AdminStats represents admin dashboard statistics
type AdminStats struct {
	TotalUsers    int `json:"total_users"`
	VerifiedUsers int `json:"verified_users"`
	AdminUsers    int `json:"admin_users"`
}

// GetOverview returns basic admin statistics
func (s *AdminService) GetOverview(ctx context.Context) (*AdminStats, error) {
	// Get basic stats
	userCount, err := s.orm.User.Query().Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user count: %w", err)
	}

	verifiedCount, err := s.orm.User.Query().Where(user.Verified(true)).Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get verified user count: %w", err)
	}

	adminCount, err := s.orm.User.Query().Where(user.Admin(true)).Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin count: %w", err)
	}

	return &AdminStats{
		TotalUsers:    userCount,
		VerifiedUsers: verifiedCount,
		AdminUsers:    adminCount,
	}, nil
}

// PaginationInfo represents pagination metadata
type PaginationInfo struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
	Pages int `json:"pages"`
}

// UserListResult represents paginated user list result
type UserListResult struct {
	Users      []*ent.User      `json:"users"`
	Pagination *PaginationInfo `json:"pagination"`
}

// ListUsers returns a paginated list of users
func (s *AdminService) ListUsers(ctx context.Context, page, limit int) (*UserListResult, error) {
	// Validate and set defaults
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 25
	}

	offset := (page - 1) * limit

	// Get users with pagination
	users, err := s.orm.User.Query().
		Offset(offset).
		Limit(limit).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	// Get total count
	total, err := s.orm.User.Query().Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user count: %w", err)
	}

	return &UserListResult{
		Users: users,
		Pagination: &PaginationInfo{
			Page:  page,
			Limit: limit,
			Total: total,
			Pages: (total + limit - 1) / limit,
		},
	}, nil
}

// GetUser returns a specific user by ID
func (s *AdminService) GetUser(ctx context.Context, userID int) (*ent.User, error) {
	u, err := s.orm.User.Get(ctx, userID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return u, nil
}

// VerifyUser manually verifies a user (admin action)
func (s *AdminService) VerifyUser(ctx context.Context, userID int) error {
	// Update user as verified
	_, err := s.orm.User.UpdateOneID(userID).
		SetVerified(true).
		ClearVerificationCode().
		Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to verify user: %w", err)
	}

	return nil
}