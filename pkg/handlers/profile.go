package handlers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/ent"
	"github.com/r-scheele/zero/pkg/context"
	"github.com/r-scheele/zero/pkg/form"
	"github.com/r-scheele/zero/pkg/middleware"
	"github.com/r-scheele/zero/pkg/msg"
	"github.com/r-scheele/zero/pkg/redirect"
	"github.com/r-scheele/zero/pkg/routenames"
	"github.com/r-scheele/zero/pkg/services"
	"github.com/r-scheele/zero/pkg/ui/forms"
	"github.com/r-scheele/zero/pkg/ui/pages"
	"golang.org/x/crypto/bcrypt"
)

type Profile struct {
	orm       *ent.Client
	container *services.Container
}

func init() {
	Register(new(Profile))
}

func (h *Profile) Init(c *services.Container) error {
	h.orm = c.ORM
	h.container = c
	return nil
}

func (h *Profile) Routes(g *echo.Group) {
	profileGroup := g.Group("/profile", middleware.RequireAuthentication, middleware.RequireVerification)
	profileGroup.GET("", h.ProfilePage).Name = routenames.Profile
	profileGroup.GET("/edit", h.ProfileEditPage).Name = routenames.ProfileEdit
	profileGroup.POST("/update", h.ProfileUpdate).Name = routenames.ProfileUpdate
	profileGroup.GET("/picture", h.ProfilePicturePage).Name = routenames.ProfilePicture
	profileGroup.POST("/picture", h.ProfilePictureSubmit)
	profileGroup.GET("/change-password", h.ChangePasswordPage).Name = routenames.ProfileChangePassword
	profileGroup.POST("/change-password", h.ChangePasswordSubmit)
	profileGroup.GET("/deactivate", h.DeactivateAccountPage).Name = routenames.ProfileDeactivate
	profileGroup.POST("/deactivate", h.DeactivateAccountSubmit)
}

func (h *Profile) ProfilePage(ctx echo.Context) error {
	userValue := ctx.Get(context.AuthenticatedUserKey)
	if userValue == nil {
		return echo.NewHTTPError(401, "User not authenticated")
	}

	u, ok := userValue.(*ent.User)
	if !ok || u == nil {
		return echo.NewHTTPError(401, "Invalid user data")
	}

	// Create and populate the profile form with current user data
	profileForm := form.Get[forms.Profile](ctx)
	if profileForm.Name == "" { // Check if form is empty instead of WasSubmitted
		profileForm.Name = u.Name
		profileForm.PhoneNumber = u.PhoneNumber
		if u.Email != nil {
			profileForm.Email = *u.Email
		}
		if u.Bio != nil {
			profileForm.Bio = *u.Bio
		}
		profileForm.DarkMode = u.DarkMode
		profileForm.EmailNotifications = u.EmailNotifications
		profileForm.SmsNotifications = u.SmsNotifications
	}

	return pages.Profile(ctx, profileForm, u)
}

func (h *Profile) ProfileEditPage(ctx echo.Context) error {
	userValue := ctx.Get(context.AuthenticatedUserKey)
	if userValue == nil {
		return echo.NewHTTPError(401, "User not authenticated")
	}

	u, ok := userValue.(*ent.User)
	if !ok || u == nil {
		return echo.NewHTTPError(401, "Invalid user data")
	}

	// Create and populate the profile form with current user data
	profileForm := form.Get[forms.Profile](ctx)
	if profileForm.Name == "" { // Check if form is empty instead of WasSubmitted
		profileForm.Name = u.Name
		profileForm.PhoneNumber = u.PhoneNumber
		if u.Email != nil {
			profileForm.Email = *u.Email
		}
		if u.Bio != nil {
			profileForm.Bio = *u.Bio
		}
		profileForm.DarkMode = u.DarkMode
		profileForm.EmailNotifications = u.EmailNotifications
		profileForm.SmsNotifications = u.SmsNotifications
	}

	return pages.ProfileEdit(ctx, profileForm)
}

func (h *Profile) ProfileUpdate(ctx echo.Context) error {
	userValue := ctx.Get(context.AuthenticatedUserKey)
	if userValue == nil {
		return echo.NewHTTPError(401, "User not authenticated")
	}

	u, ok := userValue.(*ent.User)
	if !ok || u == nil {
		return echo.NewHTTPError(401, "Invalid user data")
	}

	var input forms.Profile
	err := form.Submit(ctx, &input)

	switch err.(type) {
	case nil:
	case validator.ValidationErrors:
		return h.ProfileEditPage(ctx)
	default:
		return err
	}

	// Prepare update fields
	updateBuilder := h.orm.User.UpdateOneID(u.ID).
		SetName(input.Name).
		SetPhoneNumber(input.PhoneNumber).
		SetDarkMode(input.DarkMode).
		SetEmailNotifications(input.EmailNotifications).
		SetSmsNotifications(input.SmsNotifications)

	// Handle optional email
	if input.Email != "" {
		updateBuilder = updateBuilder.SetEmail(input.Email)
	} else {
		updateBuilder = updateBuilder.ClearEmail()
	}

	// Handle optional bio
	if input.Bio != "" {
		updateBuilder = updateBuilder.SetBio(input.Bio)
	} else {
		updateBuilder = updateBuilder.ClearBio()
	}

	// Execute the update
	updatedUser, err := updateBuilder.Save(ctx.Request().Context())
	if err != nil {
		msg.Error(ctx, "Failed to update your profile. Please try again.")
		return h.ProfileEditPage(ctx)
	}

	// Update the user in context
	ctx.Set(context.AuthenticatedUserKey, updatedUser)

	msg.Success(ctx, "Your profile has been updated successfully!")
	return redirect.New(ctx).Route(routenames.Profile).Go()
}

func (h *Profile) ProfilePicturePage(ctx echo.Context) error {
	userValue := ctx.Get(context.AuthenticatedUserKey)
	if userValue == nil {
		return echo.NewHTTPError(401, "User not authenticated")
	}

	u, ok := userValue.(*ent.User)
	if !ok || u == nil {
		return echo.NewHTTPError(401, "Invalid user data")
	}

	return pages.ProfilePicture(ctx, form.Get[forms.ProfilePicture](ctx), u)
}

func (h *Profile) ProfilePictureSubmit(ctx echo.Context) error {
	userValue := ctx.Get(context.AuthenticatedUserKey)
	if userValue == nil {
		return echo.NewHTTPError(401, "User not authenticated")
	}

	u, ok := userValue.(*ent.User)
	if !ok || u == nil {
		return echo.NewHTTPError(401, "Invalid user data")
	}

	var input forms.ProfilePicture
	err := form.Submit(ctx, &input)

	switch err.(type) {
	case nil:
	case validator.ValidationErrors:
		return h.ProfilePicturePage(ctx)
	default:
		return err
	}

	// Handle file upload
	file, err := ctx.FormFile("picture")
	if err != nil {
		input.SetFieldError("Picture", "Please select a valid image file")
		return h.ProfilePicturePage(ctx)
	}

	// Validate file type
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		input.SetFieldError("Picture", "Please upload a valid image file")
		return h.ProfilePicturePage(ctx)
	}

	// Validate file size (5MB max)
	if file.Size > 5*1024*1024 {
		input.SetFieldError("Picture", "Image file must be smaller than 5MB")
		return h.ProfilePicturePage(ctx)
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		msg.Error(ctx, "Failed to process the uploaded image. Please try again.")
		return h.ProfilePicturePage(ctx)
	}
	defer src.Close()

	// Create uploads directory if it doesn't exist
	uploadsDir := "uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		msg.Error(ctx, "Failed to save the image. Please try again.")
		return h.ProfilePicturePage(ctx)
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("profile_%d_%d%s", u.ID, time.Now().Unix(), ext)
	filePath := filepath.Join(uploadsDir, filename)

	// Create the destination file
	dst, err := os.Create(filePath)
	if err != nil {
		msg.Error(ctx, "Failed to save the image. Please try again.")
		return h.ProfilePicturePage(ctx)
	}
	defer dst.Close()

	// Copy the uploaded file to the destination
	if _, err = io.Copy(dst, src); err != nil {
		msg.Error(ctx, "Failed to save the image. Please try again.")
		return h.ProfilePicturePage(ctx)
	}

	// Delete old profile picture if exists
	if u.ProfilePicture != nil && *u.ProfilePicture != "" {
		oldPath := filepath.Join(uploadsDir, *u.ProfilePicture)
		if err := os.Remove(oldPath); err != nil {
			// Log but don't fail - old file might not exist
			fmt.Printf("Warning: failed to delete old profile picture: %v\n", err)
		}
	}

	// Update user profile picture in database
	_, err = h.orm.User.UpdateOneID(u.ID).
		SetProfilePicture(filename).
		Save(ctx.Request().Context())
	if err != nil {
		// Try to clean up the uploaded file
		os.Remove(filePath)
		msg.Error(ctx, "Failed to update profile picture. Please try again.")
		return h.ProfilePicturePage(ctx)
	}

	msg.Success(ctx, "Profile picture updated successfully!")
	return redirect.New(ctx).Route(routenames.Profile).Go()
}

func (h *Profile) ChangePasswordPage(ctx echo.Context) error {
	return pages.ChangePassword(ctx, form.Get[forms.ChangePassword](ctx))
}

func (h *Profile) ChangePasswordSubmit(ctx echo.Context) error {
	userValue := ctx.Get(context.AuthenticatedUserKey)
	if userValue == nil {
		return echo.NewHTTPError(401, "User not authenticated")
	}

	u, ok := userValue.(*ent.User)
	if !ok || u == nil {
		return echo.NewHTTPError(401, "Invalid user data")
	}

	var input forms.ChangePassword
	err := form.Submit(ctx, &input)

	switch err.(type) {
	case nil:
	case validator.ValidationErrors:
		return h.ChangePasswordPage(ctx)
	default:
		return err
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.CurrentPassword)); err != nil {
		input.SetFieldError("CurrentPassword", "Current password is incorrect")
		return h.ChangePasswordPage(ctx)
	}

	// Check if new passwords match
	if input.NewPassword != input.ConfirmPassword {
		input.SetFieldError("ConfirmPassword", "Passwords do not match")
		return h.ChangePasswordPage(ctx)
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		msg.Error(ctx, "Failed to update password. Please try again.")
		return h.ChangePasswordPage(ctx)
	}

	// Update password in database
	_, err = h.orm.User.UpdateOneID(u.ID).
		SetPassword(string(hashedPassword)).
		Save(ctx.Request().Context())
	if err != nil {
		msg.Error(ctx, "Failed to update password. Please try again.")
		return h.ChangePasswordPage(ctx)
	}

	msg.Success(ctx, "Password updated successfully!")
	return redirect.New(ctx).Route(routenames.Profile).Go()
}

func (h *Profile) DeactivateAccountPage(ctx echo.Context) error {
	return pages.DeactivateAccount(ctx, form.Get[forms.DeactivateAccount](ctx))
}

func (h *Profile) DeactivateAccountSubmit(ctx echo.Context) error {
	userValue := ctx.Get(context.AuthenticatedUserKey)
	if userValue == nil {
		return echo.NewHTTPError(401, "User not authenticated")
	}

	u, ok := userValue.(*ent.User)
	if !ok || u == nil {
		return echo.NewHTTPError(401, "Invalid user data")
	}

	var input forms.DeactivateAccount
	err := form.Submit(ctx, &input)

	switch err.(type) {
	case nil:
	case validator.ValidationErrors:
		return h.DeactivateAccountPage(ctx)
	default:
		return err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password)); err != nil {
		input.SetFieldError("Password", "Password is incorrect")
		return h.DeactivateAccountPage(ctx)
	}

	// Check confirmation field
	if input.Reason == "" {
		input.SetFieldError("Reason", "Please provide a reason for deactivation")
		return h.DeactivateAccountPage(ctx)
	}

	// Deactivate the account
	_, err = h.orm.User.UpdateOneID(u.ID).
		SetIsActive(false).
		Save(ctx.Request().Context())
	if err != nil {
		msg.Error(ctx, "Failed to deactivate account. Please try again.")
		return h.DeactivateAccountPage(ctx)
	}

	// Clear session - use auth client to logout
	err = h.container.Auth.Logout(ctx)
	if err != nil {
		msg.Error(ctx, "Failed to clear session. Please try again.")
		return h.DeactivateAccountPage(ctx)
	}

	msg.Success(ctx, "Your account has been deactivated successfully.")
	return redirect.New(ctx).Route(routenames.Login).Go()
}
