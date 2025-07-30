package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/config"
	"github.com/r-scheele/zero/ent"
	"github.com/r-scheele/zero/ent/user"
	"github.com/r-scheele/zero/pkg/context"
	"github.com/r-scheele/zero/pkg/form"
	"github.com/r-scheele/zero/pkg/log"
	"github.com/r-scheele/zero/pkg/middleware"
	"github.com/r-scheele/zero/pkg/msg"
	"github.com/r-scheele/zero/pkg/redirect"
	"github.com/r-scheele/zero/pkg/routenames"
	"github.com/r-scheele/zero/pkg/services"
	"github.com/r-scheele/zero/pkg/tasks"
	"github.com/r-scheele/zero/pkg/ui/forms"
	"github.com/r-scheele/zero/pkg/ui/pages"
)

type Auth struct {
	config    *config.Config
	auth      *services.AuthClient
	mail      *services.MailClient
	orm       *ent.Client
	container *services.Container
}

func init() {
	Register(new(Auth))
}

func (h *Auth) Init(c *services.Container) error {
	h.config = c.Config
	h.orm = c.ORM
	h.auth = c.Auth
	h.mail = c.Mail
	h.container = c
	return nil
}

func (h *Auth) Routes(g *echo.Group) {
	g.GET("/logout", h.Logout).Name = routenames.Logout // Remove RequireAuthentication middleware for faster logout
	g.GET("/email/verify/:token", h.VerifyEmail).Name = routenames.VerifyEmail
	g.GET("/verification-notice", h.VerificationNotice, middleware.RequireAuthentication).Name = routenames.VerificationNotice
	g.POST("/resend-verification", h.ResendVerification, middleware.RequireAuthentication).Name = routenames.ResendVerification

	noAuth := g.Group("/user", middleware.RequireNoAuthentication)
	noAuth.GET("/login", h.LoginPage).Name = routenames.Login
	noAuth.POST("/login", h.LoginSubmit).Name = routenames.LoginSubmit
	noAuth.GET("/register", h.RegisterPage).Name = routenames.Register
	noAuth.POST("/register", h.RegisterSubmit).Name = routenames.RegisterSubmit
	noAuth.GET("/password", h.ForgotPasswordPage).Name = routenames.ForgotPassword
	noAuth.POST("/password", h.ForgotPasswordSubmit).Name = routenames.ForgotPasswordSubmit

	resetGroup := noAuth.Group("/password/reset",
		middleware.LoadUser(h.orm),
		middleware.LoadValidPasswordToken(h.auth),
	)
	resetGroup.GET("/token/:user/:password_token/:token", h.ResetPasswordPage).Name = routenames.ResetPassword
	resetGroup.POST("/token/:user/:password_token/:token", h.ResetPasswordSubmit).Name = routenames.ResetPasswordSubmit
}

func (h *Auth) ForgotPasswordPage(ctx echo.Context) error {
	return pages.ForgotPassword(ctx, form.Get[forms.ForgotPassword](ctx))
}

func (h *Auth) ForgotPasswordSubmit(ctx echo.Context) error {
	var input forms.ForgotPassword

	err := form.Submit(ctx, &input)

	switch err.(type) {
	case nil:
	case validator.ValidationErrors:
		return h.ForgotPasswordPage(ctx)
	default:
		return err
	}

	// Find user by phone number
	user, err := h.orm.User.Query().
		Where(user.PhoneNumber(input.PhoneNumber)).
		First(ctx.Request().Context())

	if err != nil {
		// Don't reveal if user exists or not for security
		msg.Success(ctx, "If your phone number is registered, you will receive a password reset message on WhatsApp.")
		return h.ForgotPasswordPage(ctx)
	}

	// Generate WhatsApp password reset token
	resetToken, err := h.auth.GenerateWhatsAppPasswordResetToken(user.ID, user.PhoneNumber)
	if err != nil {
		msg.Error(ctx, "Failed to generate reset token. Please try again.")
		return h.ForgotPasswordPage(ctx)
	}

	// Queue password reset task
	task := tasks.PasswordResetTask{
		UserID:      user.ID,
		PhoneNumber: user.PhoneNumber,
		Username:    user.Name,
		ResetToken:  resetToken,
	}

	err = h.container.Tasks.Add(task).Save()
	if err != nil {
		msg.Error(ctx, "Failed to send password reset message. Please try again.")
		return h.ForgotPasswordPage(ctx)
	}

	msg.Success(ctx, "Password reset instructions have been sent to your WhatsApp. Please check your messages.")
	return redirect.New(ctx).Route(routenames.Login).Go()
}

func (h *Auth) LoginPage(ctx echo.Context) error {
	return pages.Login(ctx, form.Get[forms.Login](ctx))
}

func (h *Auth) LoginSubmit(ctx echo.Context) error {
	// Clear any existing authentication context to ensure fresh login attempt
	ctx.Set(context.AuthenticatedUserKey, nil)

	var input forms.Login

	// Handle form submission errors
	err := form.Submit(ctx, &input)
	if err != nil {
		switch err.(type) {
		case validator.ValidationErrors:
			return h.LoginPage(ctx)
		default:
			log.Ctx(ctx).Error("form submission error", "error", err)
			msg.Error(ctx, "An error occurred. Please try again.")
			return h.LoginPage(ctx)
		}
	}

	// Find user by phone number
	u, err := h.orm.User.
		Query().
		Where(user.PhoneNumber(strings.TrimSpace(input.PhoneNumber))).
		Only(ctx.Request().Context())

	if err != nil {
		log.Ctx(ctx).Warn("login attempt with invalid phone", "phone", input.PhoneNumber)
		input.SetFieldError("PhoneNumber", "")
		input.SetFieldError("Password", "")
		msg.Error(ctx, "Invalid phone number or password.")
		return h.LoginPage(ctx)
	}

	// Verify password
	err = h.auth.CheckPassword(input.Password, u.Password)
	if err != nil {
		log.Ctx(ctx).Warn("login attempt with invalid password", "user_id", u.ID)
		input.SetFieldError("PhoneNumber", "")
		input.SetFieldError("Password", "")
		msg.Error(ctx, "Invalid phone number or password.")
		return h.LoginPage(ctx)
	}

	// Log the user in
	err = h.auth.Login(ctx, u.ID)
	if err != nil {
		log.Ctx(ctx).Error("failed to create session", "error", err, "user_id", u.ID)
		msg.Error(ctx, "Login failed. Please try again.")
		return h.LoginPage(ctx)
	}

	log.Ctx(ctx).Info("user logged in successfully", "user_id", u.ID, "name", u.Name)

	// Check verification status after successful login
	if !u.Verified {
		msg.Warning(ctx, "Please verify your phone number to access all features.")
		return redirect.New(ctx).Route(routenames.VerificationNotice).Go()
	}

	msg.Success(ctx, fmt.Sprintf("Welcome back, %s!", u.Name))

	// Redirect based on user role
	if u.Admin {
		return redirect.New(ctx).Route("admin:overview").Go()
	}

	return ctx.Redirect(http.StatusFound, "/home")
}

func (h *Auth) Logout(ctx echo.Context) error {
	// Fast logout with minimal processing
	log.Ctx(ctx).Info("logout attempt started")

	// Clear session immediately
	h.auth.Logout(ctx)

	// Clear context
	ctx.Set(context.AuthenticatedUserKey, nil)

	// Set minimal cache headers for faster redirect
	ctx.Response().Header().Set("Cache-Control", "no-cache")
	ctx.Response().Header().Set("Location", "/")

	log.Ctx(ctx).Info("logout completed, redirecting")

	// Use faster redirect method
	return ctx.Redirect(http.StatusSeeOther, "/")
}

func (h *Auth) RegisterPage(ctx echo.Context) error {
	return pages.Register(ctx, form.Get[forms.Register](ctx))
}

func (h *Auth) RegisterSubmit(ctx echo.Context) error {
	var input forms.Register

	// Handle form submission errors
	err := form.Submit(ctx, &input)
	if err != nil {
		switch err.(type) {
		case validator.ValidationErrors:
			return h.RegisterPage(ctx)
		default:
			log.Ctx(ctx).Error("registration form submission error", "error", err)
			msg.Error(ctx, "An error occurred. Please try again.")
			return h.RegisterPage(ctx)
		}
	}

	// Clear any existing session before registration
	if existingUser := ctx.Get(context.AuthenticatedUserKey); existingUser != nil {
		h.auth.Logout(ctx)
	}

	// Generate verification code
	verificationCode := generateTwoDigitCode()

	// Hash the password before storing
	hashedPassword, err := h.auth.HashPassword(input.Password)
	if err != nil {
		log.Ctx(ctx).Error("failed to hash password", "error", err)
		msg.Error(ctx, "Registration failed. Please try again.")
		return h.RegisterPage(ctx)
	}

	// Create the user
	u, err := h.orm.User.
		Create().
		SetName(strings.TrimSpace(input.Name)).
		SetPhoneNumber(strings.TrimSpace(input.PhoneNumber)).
		SetPassword(hashedPassword).
		SetRegistrationMethod("web").
		SetVerificationCode(verificationCode).
		Save(ctx.Request().Context())

	if err != nil {
		switch err.(type) {
		case *ent.ConstraintError:
			log.Ctx(ctx).Warn("registration attempt with existing phone", "phone", input.PhoneNumber)
			input.SetFieldError("PhoneNumber", "This phone number is already registered")
			msg.Warning(ctx, "A user with this phone number already exists. Please log in instead.")
			return h.RegisterPage(ctx)
		default:
			log.Ctx(ctx).Error("failed to create user", "error", err)
			msg.Error(ctx, "Registration failed. Please try again.")
			return h.RegisterPage(ctx)
		}
	}

	log.Ctx(ctx).Info("user registered successfully", "user_id", u.ID, "name", u.Name, "phone", u.PhoneNumber)

	// Log the user in immediately after registration
	err = h.auth.Login(ctx, u.ID)
	if err != nil {
		log.Ctx(ctx).Error("failed to login after registration", "error", err, "user_id", u.ID)
		msg.Success(ctx, "Account created successfully! Please log in.")
		return redirect.New(ctx).Route(routenames.Login).Go()
	}

	// Get verification code for display
	displayCode := ""
	if u.VerificationCode != nil {
		displayCode = *u.VerificationCode
	}

	msg.Success(ctx, fmt.Sprintf("🎉 Welcome, %s! Your account has been created. Verification code: %s", u.Name, displayCode))

	// Send phone verification
	h.sendPhoneVerification(ctx, u, "web")

	// Redirect to verification notice since user needs to verify
	return redirect.New(ctx).Route(routenames.VerificationNotice).Go()
}

func (h *Auth) sendPhoneVerification(ctx echo.Context, usr *ent.User, method string) {
	// Use the async phone verification task instead of sending directly
	err := SendPhoneVerification(ctx, h.container, usr, method)
	if err != nil {
		log.Ctx(ctx).Error("unable to queue phone verification task",
			"user_id", usr.ID,
			"error", err,
		)
		// Don't fail the registration, just log the error
		return
	}

	if method == "web" {
		// Message already shown in registration success - no need for additional message
	} else {
		msg.Info(ctx, "Welcome! You can now log in to our web platform.")
	}
}

func (h *Auth) ResetPasswordPage(ctx echo.Context) error {
	return pages.ResetPassword(ctx, form.Get[forms.ResetPassword](ctx))
}

func (h *Auth) ResetPasswordSubmit(ctx echo.Context) error {
	var input forms.ResetPassword

	err := form.Submit(ctx, &input)

	switch err.(type) {
	case nil:
	case validator.ValidationErrors:
		return h.ResetPasswordPage(ctx)
	default:
		return err
	}

	// Get the requesting user.
	usr := ctx.Get(context.UserKey).(*ent.User)

	// Hash the new password before storing
	hashedPassword, err := h.auth.HashPassword(input.Password)
	if err != nil {
		log.Ctx(ctx).Error("failed to hash password", "error", err)
		msg.Error(ctx, "Password reset failed. Please try again.")
		return h.ResetPasswordPage(ctx)
	}

	// Update the user.
	_, err = usr.
		Update().
		SetPassword(hashedPassword).
		Save(ctx.Request().Context())

	if err != nil {
		return fail(err, "unable to update password")
	}

	// Delete all password tokens for this user.
	err = h.auth.DeletePasswordTokens(ctx, usr.ID)
	if err != nil {
		return fail(err, "unable to delete password tokens")
	}

	msg.Success(ctx, "Your password has been updated.")
	return redirect.New(ctx).
		Route(routenames.Login).
		Go()
}

func (h *Auth) VerifyEmail(ctx echo.Context) error {
	// TODO: Implement phone verification via WhatsApp confirmation
	msg.Warning(ctx, "Phone verification is handled via WhatsApp. Please check your WhatsApp messages.")
	return ctx.Redirect(http.StatusFound, "/home")
}

func (h *Auth) VerificationNotice(ctx echo.Context) error {
	return pages.VerificationNotice(ctx)
}

func (h *Auth) ResendVerification(ctx echo.Context) error {
	if u := ctx.Get(context.AuthenticatedUserKey); u != nil {
		if user, ok := u.(*ent.User); ok {
			if user.Verified {
				msg.Info(ctx, "Your phone number is already verified.")
				return ctx.Redirect(http.StatusFound, "/home")
			}

			// Send phone verification
			h.sendPhoneVerification(ctx, user, "web")
			msg.Success(ctx, "Verification message has been resent to your WhatsApp. Please check your messages.")
			return redirect.New(ctx).Route(routenames.VerificationNotice).Go()
		}
	}

	return echo.NewHTTPError(http.StatusUnauthorized)
}

// generateTwoDigitCode generates a random 2-digit verification code (10-99)
func generateTwoDigitCode() string {
	code := rand.Intn(90) + 10 // generates numbers from 10 to 99
	return fmt.Sprintf("%02d", code)
}
