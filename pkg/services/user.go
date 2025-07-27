package services

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/r-scheele/zero/ent"
	"github.com/r-scheele/zero/ent/user"
	"github.com/r-scheele/zero/pkg/ui/forms"
)

// UserService handles user-related operations
type UserService struct {
	orm  *ent.Client
	auth *AuthClient
}

// NewUserService creates a new user service
func NewUserService(orm *ent.Client, auth *AuthClient) *UserService {
	return &UserService{
		orm:  orm,
		auth: auth,
	}
}

// CreateUser creates a new user account
func (s *UserService) CreateUser(ctx context.Context, input forms.Register) (*ent.User, error) {
	// Check if user already exists
	existingUser, err := s.orm.User.Query().
		Where(user.PhoneNumber(strings.TrimSpace(input.PhoneNumber))).
		First(ctx)

	if err == nil {
		// User exists
		if existingUser.Verified {
			return nil, fmt.Errorf("phone number is already registered and verified")
		}
		return existingUser, nil
	} else if !ent.IsNotFound(err) {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Create new user
	u, err := s.orm.User.Create().
		SetName(input.Name).
		SetPhoneNumber(strings.TrimSpace(input.PhoneNumber)).
		SetPassword(input.Password).
		SetRegistrationMethod("mobile").
		SetVerificationCode(s.generateTwoDigitCode()).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create user account: %w", err)
	}

	return u, nil
}

// AuthenticateUser authenticates a user by phone number and password
func (s *UserService) AuthenticateUser(ctx context.Context, phoneNumber, password string) (*ent.User, error) {
	// Find user by phone number
	u, err := s.orm.User.
		Query().
		Where(user.PhoneNumber(strings.TrimSpace(phoneNumber))).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Check password
	if err = s.auth.CheckPassword(password, u.Password); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return u, nil
}

// GetUserByPhoneNumber retrieves a user by phone number
func (s *UserService) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*ent.User, error) {
	return s.orm.User.Query().
		Where(user.PhoneNumber(phoneNumber)).
		First(ctx)
}

// UpdateProfile updates user profile information
func (s *UserService) UpdateProfile(ctx context.Context, userID int, input forms.Profile) (*ent.User, error) {
	updateBuilder := s.orm.User.UpdateOneID(userID).
		SetName(input.Name).
		SetPhoneNumber(input.PhoneNumber).
		SetDarkMode(input.DarkMode).
		SetEmailNotifications(input.EmailNotifications).
		SetSmsNotifications(input.SmsNotifications)

	if input.Email != "" {
		updateBuilder = updateBuilder.SetEmail(input.Email)
	} else {
		updateBuilder = updateBuilder.ClearEmail()
	}

	if input.Bio != "" {
		updateBuilder = updateBuilder.SetBio(input.Bio)
	} else {
		updateBuilder = updateBuilder.ClearBio()
	}

	updatedUser, err := updateBuilder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return updatedUser, nil
}

// UpdateProfilePicture handles profile picture upload and update
func (s *UserService) UpdateProfilePicture(ctx context.Context, userID int, file io.Reader, filename string, size int64) (*ent.User, error) {
	// Validate file size (5MB max)
	if size > 5*1024*1024 {
		return nil, fmt.Errorf("image file must be smaller than 5MB")
	}

	// Create uploads directory
	uploadsDir := "uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Generate unique filename
	newFilename := fmt.Sprintf("profile_%d_%d%s", userID, time.Now().Unix(), filepath.Ext(filename))
	filePath := filepath.Join(uploadsDir, newFilename)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to save image: %w", err)
	}
	defer dst.Close()

	// Copy file
	if _, err = io.Copy(dst, file); err != nil {
		return nil, fmt.Errorf("failed to save image: %w", err)
	}

	// Update user profile picture
	updatedUser, err := s.orm.User.UpdateOneID(userID).
		SetProfilePicture(newFilename).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to update profile picture: %w", err)
	}

	return updatedUser, nil
}

// ChangePassword updates user password after validating current password
func (s *UserService) ChangePassword(ctx context.Context, userID int, currentPassword, newPassword string) error {
	// Get user
	u, err := s.orm.User.Get(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Check current password
	if err := s.auth.CheckPassword(currentPassword, u.Password); err != nil {
		return fmt.Errorf("current password is incorrect")
	}

	// Update password
	_, err = s.orm.User.UpdateOneID(userID).
		SetPassword(newPassword).
		Save(ctx)

	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// DeactivateAccount deletes a user account after password verification
func (s *UserService) DeactivateAccount(ctx context.Context, userID int, password string) error {
	// Get user
	u, err := s.orm.User.Get(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Validate password
	if err := s.auth.CheckPassword(password, u.Password); err != nil {
		return fmt.Errorf("password is incorrect")
	}

	// Delete user account
	if err := s.orm.User.DeleteOneID(userID).Exec(ctx); err != nil {
		return fmt.Errorf("failed to deactivate account: %w", err)
	}

	return nil
}

// VerifyUser marks a user as verified
func (s *UserService) VerifyUser(ctx context.Context, userID int) error {
	_, err := s.orm.User.UpdateOneID(userID).
		SetVerified(true).
		ClearVerificationCode().
		Save(ctx)

	if err != nil {
		return fmt.Errorf("failed to verify user: %w", err)
	}

	return nil
}

// VerifyUserByCode verifies a user using their verification code
func (s *UserService) VerifyUserByCode(ctx context.Context, phoneNumber, code string) (*ent.User, error) {
	// Find user by phone number
	u, err := s.orm.User.Query().
		Where(user.PhoneNumber(phoneNumber)).
		First(ctx)

	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Check if verification code matches
	if u.VerificationCode == nil || *u.VerificationCode != strings.TrimSpace(code) {
		return nil, fmt.Errorf("invalid verification code")
	}

	// Verify user
	updatedUser, err := u.Update().
		SetVerified(true).
		ClearVerificationCode().
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to verify user: %w", err)
	}

	return updatedUser, nil
}

// generateTwoDigitCode generates a random 2-digit verification code (10-99)
func (s *UserService) generateTwoDigitCode() string {
	code := 10 + (time.Now().UnixNano() % 90) // generates numbers from 10 to 99
	return fmt.Sprintf("%02d", code)
}