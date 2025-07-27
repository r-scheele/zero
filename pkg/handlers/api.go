package handlers

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/config"
	"github.com/r-scheele/zero/ent"
	"github.com/r-scheele/zero/ent/user"
	pkgcontext "github.com/r-scheele/zero/pkg/context"
	"github.com/r-scheele/zero/pkg/services"
	"github.com/r-scheele/zero/pkg/tasks"
	"github.com/r-scheele/zero/pkg/ui/forms"
	"github.com/r-scheele/zero/pkg/ui/models"
	"github.com/spf13/afero"
)

type API struct {
	container *services.Container
	config    *config.Config
	auth      *services.AuthClient
	mail      *services.MailClient
	orm       *ent.Client
	files     afero.Fs
}

func init() {
	Register(new(API))
}

// Init initializes the API handler with the service container
func (h *API) Init(c *services.Container) error {
	h.container = c
	h.config = c.Config
	h.orm = c.ORM
	h.auth = c.Auth
	h.mail = c.Mail
	h.files = c.Files
	return nil
}

// Routes registers all external API routes
func (h *API) Routes(g *echo.Group) {
	api_str := "/api/v1/"
	// Main API endpoints
	apiGroup := g.Group(api_str)

	// Health check
	apiGroup.GET("/health", h.HealthCheck)

	// WhatsApp webhook endpoints (360dialog integration)
	webhookGroup := g.Group(api_str + "whatsapp")
	webhookGroup.GET("/webhook", h.VerifyWebhook)
	webhookGroup.POST("/webhook", h.HandleWebhook)

	// Mobile API v1 routes
	mobileAPI := g.Group("/api/v1/mobile")

	// Authentication endpoints
	auth := mobileAPI.Group("/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
	auth.POST("/logout", h.Logout)
	auth.POST("/forgot-password", h.ForgotPassword)
	auth.POST("/reset-password", h.ResetPassword)
	auth.POST("/resend-verification", h.ResendVerification)

	// Profile endpoints (require authentication)
	profile := mobileAPI.Group("/profile")
	profile.Use(h.requireAuth)
	profile.GET("", h.GetProfile)
	profile.PUT("", h.UpdateProfile)
	profile.POST("/picture", h.UpdateProfilePicture)
	profile.POST("/change-password", h.ChangePassword)
	profile.POST("/deactivate", h.DeactivateAccount)

	// Contact endpoints
	contact := mobileAPI.Group("/contact")
	contact.POST("", h.SubmitContact)

	// File endpoints (require authentication)
	files := mobileAPI.Group("/files")
	files.Use(h.requireAuth)
	files.GET("", h.ListFiles)
	files.POST("/upload", h.UploadFile)

	// Task endpoints (require authentication)
	tasks := mobileAPI.Group("/tasks")
	tasks.Use(h.requireAuth)
	tasks.POST("", h.CreateTask)

	// Search endpoints
	search := mobileAPI.Group("/search")
	search.GET("", h.Search)

	// Admin endpoints (require admin privileges)
	admin := mobileAPI.Group("/admin")
	admin.Use(h.requireAuth, h.requireAdmin)
	admin.GET("/overview", h.AdminOverview)
	admin.GET("/users", h.AdminListUsers)
	admin.GET("/users/:id", h.AdminGetUser)
	admin.POST("/users/:id/verify", h.AdminVerifyUser)
}

// HealthCheck for API availability
func (h *API) HealthCheck(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"status":    "healthy",
		"service":   "external-api",
		"timestamp": ctx.Request().Context().Value("timestamp"),
	})
}

// Mobile API Authentication Methods
func (h *API) Register(ctx echo.Context) error {
	var input forms.Register

	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Validate input
	validate := validator.New()
	if err := validate.Struct(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "Validation failed",
			"fields": err.(validator.ValidationErrors),
		})
	}

	// Check if user already exists
	existingUser, err := h.orm.User.Query().
		Where(user.PhoneNumber(strings.TrimSpace(input.PhoneNumber))).
		First(ctx.Request().Context())

	if err == nil {
		// User exists
		if existingUser.Verified {
			return ctx.JSON(http.StatusConflict, map[string]string{
				"error": "Phone number is already registered and verified",
			})
		} else {
			// User exists but not verified, resend verification
			if err = SendPhoneVerification(ctx, h.container, existingUser, "mobile"); err != nil {
				return ctx.JSON(http.StatusInternalServerError, map[string]string{
					"error": "Failed to send verification message",
				})
			}
			return ctx.JSON(http.StatusOK, map[string]interface{}{
				"message": "Verification message has been resent to your WhatsApp",
				"user_id": existingUser.ID,
			})
		}
	} else if !ent.IsNotFound(err) {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Database error",
		})
	}

	// Create new user
	u, err := h.orm.User.Create().
		SetName(input.Name).
		SetPhoneNumber(strings.TrimSpace(input.PhoneNumber)).
		SetPassword(input.Password).
		SetRegistrationMethod("mobile").
		SetVerificationCode(h.generateTwoDigitCode()).
		Save(ctx.Request().Context())

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create user account",
		})
	}

	// Send phone verification
	if err := SendPhoneVerification(ctx, h.container, u, "mobile"); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Account created but failed to send verification message",
		})
	}

	return ctx.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Account created successfully. Please check your WhatsApp for verification code.",
		"user_id": u.ID,
	})
}

func (h *API) Login(ctx echo.Context) error {
	var input forms.Login

	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Find user by phone number
	u, err := h.orm.User.
		Query().
		Where(user.PhoneNumber(strings.TrimSpace(input.PhoneNumber))).
		Only(ctx.Request().Context())

	if err != nil {
		if ent.IsNotFound(err) {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid credentials",
			})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Database error",
		})
	}

	// Check password
	if err = h.auth.CheckPassword(input.Password, u.Password); err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid credentials",
		})
	}

	// Generate JWT token
	token, err := h.generateJWTToken(u.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate session token",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message":  "Login successful",
		"user_id":  u.ID,
		"token":    token,
		"verified": u.Verified,
		"admin":    u.Admin,
		"name":     u.Name,
	})
}

func (h *API) Logout(ctx echo.Context) error {
	// For mobile API, we just return success since token management is client-side
	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Logout successful",
	})
}

func (h *API) ForgotPassword(ctx echo.Context) error {
	var input forms.ForgotPassword

	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Find user by phone number
	u, err := h.orm.User.Query().
		Where(user.PhoneNumber(input.PhoneNumber)).
		First(ctx.Request().Context())

	if err != nil {
		// Don't reveal if user exists for security
		return ctx.JSON(http.StatusOK, map[string]string{
			"message": "If your phone number is registered, you will receive a password reset message on WhatsApp",
		})
	}

	// Generate reset token
	resetToken, err := h.auth.GenerateWhatsAppPasswordResetToken(u.ID, u.PhoneNumber)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate reset token",
		})
	}

	// Queue password reset task
	task := tasks.PasswordResetTask{
		UserID:      u.ID,
		PhoneNumber: u.PhoneNumber,
		Username:    u.Name,
		ResetToken:  resetToken,
	}

	if err := h.container.Tasks.Add(task).Save(); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to send password reset message",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Password reset instructions have been sent to your WhatsApp",
	})
}

func (h *API) ResetPassword(ctx echo.Context) error {
	var input forms.ResetPassword

	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Get user from context (should be set by middleware)
	u := ctx.Get(pkgcontext.UserKey).(*ent.User)

	// Update password
	_, err := u.Update().
		SetPassword(input.Password).
		Save(ctx.Request().Context())

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update password",
		})
	}

	// Delete all password tokens
	if err := h.auth.DeletePasswordTokens(ctx, u.ID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to clean up tokens",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Password updated successfully",
	})
}

func (h *API) ResendVerification(ctx echo.Context) error {
	u := ctx.Get(pkgcontext.AuthenticatedUserKey).(*ent.User)

	if u.Verified {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Phone number is already verified",
		})
	}

	// Send phone verification
	if err := SendPhoneVerification(ctx, h.container, u, "mobile"); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to send verification message",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Verification message has been resent to your WhatsApp",
	})
}

// Profile endpoints
func (h *API) GetProfile(ctx echo.Context) error {
	u := ctx.Get(pkgcontext.AuthenticatedUserKey).(*ent.User)

	profile := map[string]interface{}{
		"id":                  u.ID,
		"name":                u.Name,
		"phone_number":        u.PhoneNumber,
		"email":               u.Email,
		"bio":                 u.Bio,
		"verified":            u.Verified,
		"admin":               u.Admin,
		"dark_mode":           u.DarkMode,
		"email_notifications": u.EmailNotifications,
		"sms_notifications":   u.SmsNotifications,
		"profile_picture":     u.ProfilePicture,
		"registration_method": u.RegistrationMethod,
		"created_at":          u.CreatedAt,
		"updated_at":          u.UpdatedAt,
	}

	return ctx.JSON(http.StatusOK, profile)
}

func (h *API) UpdateProfile(ctx echo.Context) error {
	u := ctx.Get(pkgcontext.AuthenticatedUserKey).(*ent.User)
	var input forms.Profile

	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Validate input
	validate := validator.New()
	if err := validate.Struct(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "Validation failed",
			"fields": err.(validator.ValidationErrors),
		})
	}

	// Update user
	updateBuilder := h.orm.User.UpdateOneID(u.ID).
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

	updatedUser, err := updateBuilder.Save(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update profile",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "Profile updated successfully",
		"user":    updatedUser,
	})
}

func (h *API) UpdateProfilePicture(ctx echo.Context) error {
	u := ctx.Get(pkgcontext.AuthenticatedUserKey).(*ent.User)

	// Handle file upload
	file, err := ctx.FormFile("picture")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Picture file is required",
		})
	}

	// Validate file type
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Please upload a valid image file",
		})
	}

	// Validate file size (5MB max)
	if file.Size > 5*1024*1024 {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Image file must be smaller than 5MB",
		})
	}

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to process uploaded image",
		})
	}
	defer src.Close()

	// Create uploads directory
	uploadsDir := "uploads"
	if err = os.MkdirAll(uploadsDir, 0755); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create upload directory",
		})
	}

	// Generate unique filename
	filename := fmt.Sprintf("profile_%d_%d%s", u.ID, time.Now().Unix(), filepath.Ext(file.Filename))
	filePath := filepath.Join(uploadsDir, filename)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to save image",
		})
	}
	defer dst.Close()

	// Copy file
	if _, err = io.Copy(dst, src); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to save image",
		})
	}

	// Update user profile picture
	updatedUser, err := h.orm.User.UpdateOneID(u.ID).
		SetProfilePicture(filename).
		Save(ctx.Request().Context())

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update profile picture",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message":         "Profile picture updated successfully",
		"profile_picture": updatedUser.ProfilePicture,
	})
}

func (h *API) ChangePassword(ctx echo.Context) error {
	u := ctx.Get(pkgcontext.AuthenticatedUserKey).(*ent.User)
	var input struct {
		CurrentPassword string `json:"current_password" validate:"required"`
		NewPassword     string `json:"new_password" validate:"required,min=8"`
	}

	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Validate input
	validate := validator.New()
	if err := validate.Struct(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "Validation failed",
			"fields": err.(validator.ValidationErrors),
		})
	}

	// Check current password
	if err := h.auth.CheckPassword(input.CurrentPassword, u.Password); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Current password is incorrect",
		})
	}

	// Update password
	_, err := h.orm.User.UpdateOneID(u.ID).
		SetPassword(input.NewPassword).
		Save(ctx.Request().Context())

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update password",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Password updated successfully",
	})
}

func (h *API) DeactivateAccount(ctx echo.Context) error {
	u := ctx.Get(pkgcontext.AuthenticatedUserKey).(*ent.User)
	var input struct {
		Password string `json:"password" validate:"required"`
		Confirm  bool   `json:"confirm" validate:"required"`
	}

	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Validate password
	if err := h.auth.CheckPassword(input.Password, u.Password); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Password is incorrect",
		})
	}

	if !input.Confirm {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Account deactivation must be confirmed",
		})
	}

	// Delete user account
	if err := h.orm.User.DeleteOneID(u.ID).Exec(ctx.Request().Context()); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to deactivate account",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Account deactivated successfully",
	})
}

// Contact endpoint
func (h *API) SubmitContact(ctx echo.Context) error {
	var input forms.Contact

	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Validate input
	validate := validator.New()
	if err := validate.Struct(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "Validation failed",
			"fields": err.(validator.ValidationErrors),
		})
	}

	// Send email
	err := h.mail.
		Compose().
		To(input.Email).
		Subject("Contact form submitted").
		Body(fmt.Sprintf("The message is: %s", input.Message)).
		Send(ctx)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to send email",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Contact form submitted successfully",
	})
}

// File endpoints
func (h *API) ListFiles(ctx echo.Context) error {
	// Get list of uploaded files
	info, err := afero.ReadDir(h.files, "")
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to list files",
		})
	}

	files := make([]*models.File, 0)
	for _, file := range info {
		files = append(files, &models.File{
			Name:     file.Name(),
			Size:     file.Size(),
			Modified: file.ModTime().Format(time.DateTime),
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"files": files,
	})
}

func (h *API) UploadFile(ctx echo.Context) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "File is required",
		})
	}

	src, err := file.Open()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to open uploaded file",
		})
	}
	defer src.Close()

	dst, err := h.files.Create(file.Filename)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create file",
		})
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to save file",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"message":  "File uploaded successfully",
		"filename": file.Filename,
		"size":     file.Size,
	})
}

// Task endpoint
func (h *API) CreateTask(ctx echo.Context) error {
	var input forms.Task

	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Validate input
	validate := validator.New()
	if err := validate.Struct(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "Validation failed",
			"fields": err.(validator.ValidationErrors),
		})
	}

	// Create task
	err := h.container.Tasks.
		Add(tasks.ExampleTask{
			Message: input.Message,
		}).
		Wait(time.Duration(input.Delay) * time.Second).
		Save()

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create task",
		})
	}

	return ctx.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Task created successfully",
		"delay":   input.Delay,
	})
}

// Search endpoint
func (h *API) Search(ctx echo.Context) error {
	query := ctx.QueryParam("q")
	if query == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Search query is required",
		})
	}

	// For now, return empty results (implement actual search logic as needed)
	results := make([]*models.SearchResult, 0)

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"query":   query,
		"results": results,
	})
}

// Admin endpoints
func (h *API) AdminOverview(ctx echo.Context) error {
	// Get basic stats
	userCount, err := h.orm.User.Query().Count(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get user count",
		})
	}

	verifiedCount, err := h.orm.User.Query().Where(user.Verified(true)).Count(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get verified user count",
		})
	}

	adminCount, err := h.orm.User.Query().Where(user.Admin(true)).Count(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get admin count",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"stats": map[string]int{
			"total_users":    userCount,
			"verified_users": verifiedCount,
			"admin_users":    adminCount,
		},
	})
}

func (h *API) AdminListUsers(ctx echo.Context) error {
	// Parse pagination parameters
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 25
	}

	offset := (page - 1) * limit

	// Get users with pagination
	users, err := h.orm.User.Query().
		Offset(offset).
		Limit(limit).
		All(ctx.Request().Context())

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get users",
		})
	}

	// Get total count
	total, err := h.orm.User.Query().Count(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get user count",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"users": users,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + limit - 1) / limit,
		},
	})
}

func (h *API) AdminGetUser(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}

	u, err := h.orm.User.Get(ctx.Request().Context(), id)
	if err != nil {
		if ent.IsNotFound(err) {
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"error": "User not found",
			})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get user",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"user": u,
	})
}

func (h *API) AdminVerifyUser(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}

	// Update user as verified
	_, err = h.orm.User.UpdateOneID(id).
		SetVerified(true).
		ClearVerificationCode().
		Save(ctx.Request().Context())

	if err != nil {
		if ent.IsNotFound(err) {
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"error": "User not found",
			})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to verify user",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "User verified successfully",
	})
}

// Middleware functions
func (h *API) requireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		token := ctx.Request().Header.Get("Authorization")
		if token == "" {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Authorization token is required",
			})
		}

		// Remove "Bearer " prefix if present
		token = strings.TrimPrefix(token, "Bearer ")

		// Validate JWT token
		userID, err := h.validateJWTToken(token)
		if err != nil {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid or expired token",
			})
		}

		// Get user from database
		u, err := h.orm.User.Get(ctx.Request().Context(), userID)
		if err != nil {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{
				"error": "User not found",
			})
		}

		// Store user in context
		ctx.Set(pkgcontext.AuthenticatedUserKey, u)
		return next(ctx)
	}
}

func (h *API) requireAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		u := ctx.Get(pkgcontext.AuthenticatedUserKey).(*ent.User)
		if !u.Admin {
			return ctx.JSON(http.StatusForbidden, map[string]string{
				"error": "Admin privileges required",
			})
		}
		return next(ctx)
	}
}

// Helper methods for JWT token management
func (h *API) generateJWTToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // 24 hour expiry
		"iat":     time.Now().Unix(),
	})

	return token.SignedString([]byte(h.config.App.EncryptionKey))
}

func (h *API) validateJWTToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(h.config.App.EncryptionKey), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userIDFloat, ok := claims["user_id"].(float64); ok {
			return int(userIDFloat), nil
		}
	}

	return 0, fmt.Errorf("invalid token claims")
}

// Helper function to generate two-digit verification code
func (h *API) generateTwoDigitCode() string {
	code := rand.Intn(90) + 10 // generates numbers from 10 to 99
	return fmt.Sprintf("%02d", code)
}

// WhatsApp webhook payload structures
type WebhookPayload struct {
	Statuses []Status  `json:"statuses,omitempty"`
	Messages []Message `json:"messages,omitempty"`
}

type Status struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	Timestamp   string `json:"timestamp"`
	RecipientID string `json:"recipient_id"`
}

type Message struct {
	From      string         `json:"from"`
	ID        string         `json:"id"`
	Timestamp string         `json:"timestamp"`
	Type      string         `json:"type"`
	Text      *TextMessage   `json:"text,omitempty"`
	Button    *ButtonMessage `json:"button,omitempty"`
}

type TextMessage struct {
	Body string `json:"body"`
}

type ButtonMessage struct {
	Text    string `json:"text"`
	Payload string `json:"payload"`
}

// getEnvWithDefault gets environment variable with default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// VerifyWebhook handles webhook verification from 360dialog
func (h *API) VerifyWebhook(c echo.Context) error {
	mode := c.QueryParam("hub.mode")
	token := c.QueryParam("hub.verify_token")
	challenge := c.QueryParam("hub.challenge")

	expectedToken := getEnvWithDefault("WHATSAPP_VERIFY_TOKEN", "zero_webhook_token")

	if mode == "subscribe" && token == expectedToken {
		return c.String(http.StatusOK, challenge)
	}

	return c.String(http.StatusForbidden, "Forbidden")
}

// HandleWebhook processes incoming WhatsApp messages and button responses
func (h *API) HandleWebhook(c echo.Context) error {
	var payload WebhookPayload
	if err := c.Bind(&payload); err != nil {
		return c.String(http.StatusBadRequest, "Invalid payload")
	}

	ctx := c.Request().Context()

	// Process incoming messages
	for _, message := range payload.Messages {
		if err := h.processMessage(ctx, message); err != nil {
			// Log error but continue processing other messages
		}
	}

	// Process message statuses (delivered, read, etc.)
	for range payload.Statuses {
		// Log status updates
	}

	return c.String(http.StatusOK, "OK")
}

// processMessage handles individual WhatsApp messages
func (h *API) processMessage(ctx context.Context, message Message) error {
	// Handle verification code responses
	if message.Text != nil {
		return h.handleTextMessage(ctx, message)
	}

	// Handle button responses
	if message.Button != nil {
		return h.handleButtonMessage(ctx, message)
	}

	return nil
}

// handleTextMessage processes text messages (verification codes)
func (h *API) handleTextMessage(ctx context.Context, message Message) error {
	// Find user by phone number
	u, err := h.orm.User.Query().
		Where(user.PhoneNumber(message.From)).
		First(ctx)

	if err != nil {
		return err
	}

	// Check if message is a verification code
	if u.VerificationCode != nil && *u.VerificationCode == strings.TrimSpace(message.Text.Body) {
		// Verify user
		_, err := u.Update().
			SetVerified(true).
			ClearVerificationCode().
			Save(ctx)

		if err != nil {
			return err
		}

		// Send confirmation message
		return h.sendWhatsAppMessage(message.From, "✅ Your phone number has been verified successfully!")
	}

	return nil
}

// handleButtonMessage processes button responses
func (h *API) handleButtonMessage(ctx context.Context, message Message) error {
	// Handle different button payloads
	switch message.Button.Payload {
	case "verify_account":
		return h.handleVerifyAccountButton(ctx, message)
	case "reset_password":
		return h.handleResetPasswordButton(ctx, message)
	default:
		return nil
	}
}

// handleVerifyAccountButton processes verify account button clicks
func (h *API) handleVerifyAccountButton(ctx context.Context, message Message) error {
	// Find user by phone number
	u, err := h.orm.User.Query().
		Where(user.PhoneNumber(message.From)).
		First(ctx)

	if err != nil {
		return err
	}

	if u.Verified {
		return h.sendWhatsAppMessage(message.From, "Your account is already verified.")
	}

	// Resend verification code
	return SendPhoneVerification(nil, h.container, u, "mobile")
}

// handleResetPasswordButton processes reset password button clicks
func (h *API) handleResetPasswordButton(ctx context.Context, message Message) error {
	// Find user by phone number
	u, err := h.orm.User.Query().
		Where(user.PhoneNumber(message.From)).
		First(ctx)

	if err != nil {
		return err
	}

	// Generate reset token
	resetToken, err := h.auth.GenerateWhatsAppPasswordResetToken(u.ID, u.PhoneNumber)
	if err != nil {
		return err
	}

	// Queue password reset task
	task := tasks.PasswordResetTask{
		UserID:      u.ID,
		PhoneNumber: u.PhoneNumber,
		Username:    u.Name,
		ResetToken:  resetToken,
	}

	return h.container.Tasks.Add(task).Save()
}

// sendWhatsAppMessage sends a message via WhatsApp API
func (h *API) sendWhatsAppMessage(to, message string) error {
	// Implementation depends on your WhatsApp API client
	// This is a placeholder
	return nil
}
