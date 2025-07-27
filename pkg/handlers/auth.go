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
	var input forms.Login

	authFailed := func() error {
		input.SetFieldError("PhoneNumber", "")
		input.SetFieldError("Password", "")
		msg.Error(ctx, "Invalid credentials. Please try again.")
		return h.LoginPage(ctx)
	}

	err := form.Submit(ctx, &input)

	switch err.(type) {
	case nil:
	case validator.ValidationErrors:
		return h.LoginPage(ctx)
	default:
		return err
	}

	// Attempt to load the user.
	u, err := h.orm.User.
		Query().
		Where(user.PhoneNumber(strings.TrimSpace(input.PhoneNumber))).
		Only(ctx.Request().Context())

	switch err.(type) {
	case *ent.NotFoundError:
		return authFailed()
	case nil:
	default:
		return fail(err, "error querying user during login")
	}

	// Check if the password is correct.
	err = h.auth.CheckPassword(input.Password, u.Password)
	if err != nil {
		return authFailed()
	}

	// Check if the user's phone number is verified
	if !u.Verified {
		// Log the user in temporarily so they can access the verification page
		err = h.auth.Login(ctx, u.ID)
		if err != nil {
			return fail(err, "unable to log in unverified user")
		}

		msg.Warning(ctx, "Your phone number is not yet verified. Please verify your phone to access all features.")
		return redirect.New(ctx).Route(routenames.VerificationNotice).Go()
	}

	// Log the user in.
	err = h.auth.Login(ctx, u.ID)
	if err != nil {
		return fail(err, "unable to log in user")
	}

	msg.Success(ctx, fmt.Sprintf("Welcome back, %s. You are now logged in.", u.Name))

	// Redirect admin users to admin panel, regular users to home
	if u.Admin {
		return redirect.New(ctx).
			Route("admin:overview").
			Go()
	}

	return redirect.New(ctx).
		Route(routenames.Home).
		Go()
}

func (h *Auth) Logout(ctx echo.Context) error {
	// Always try to logout, even if there are errors - be resilient
	h.auth.Logout(ctx)

	// Clear the authenticated user from the current request context
	ctx.Set(context.AuthenticatedUserKey, nil)

	// Add essential headers for security
	ctx.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Response().Header().Set("Pragma", "no-cache")
	ctx.Response().Header().Set("Expires", "0")

	// If this is an HTMX request, force a full page reload
	if ctx.Request().Header.Get("HX-Request") == "true" {
		ctx.Response().Header().Set("HX-Redirect", "/")
		return nil
	}

	// Use direct redirect for non-HTMX requests
	return ctx.Redirect(302, "/")
}

func (h *Auth) RegisterPage(ctx echo.Context) error {
	return pages.Register(ctx, form.Get[forms.Register](ctx))
}

func (h *Auth) RegisterSubmit(ctx echo.Context) error {
	var input forms.Register

	err := form.Submit(ctx, &input)

	switch err.(type) {
	case nil:
	case validator.ValidationErrors:
		return h.RegisterPage(ctx)
	default:
		return err
	}

	// Generate a 2-digit verification code for WhatsApp verification
	verificationCode := generateTwoDigitCode()

	// Attempt creating the user.
	u, err := h.orm.User.
		Create().
		SetName(input.Name).
		SetPhoneNumber(input.PhoneNumber).
		SetPassword(input.Password).
		SetRegistrationMethod("web").
		SetVerificationCode(verificationCode).
		Save(ctx.Request().Context())

	switch err.(type) {
	case nil:
		log.Ctx(ctx).Info("user created",
			"user_name", u.Name,
			"user_id", u.ID,
			"phone_number", u.PhoneNumber,
		)
	case *ent.ConstraintError:
		msg.Warning(ctx, "A user with this phone number already exists. Please log in.")
		return redirect.New(ctx).
			Route(routenames.Login).
			Go()
	default:
		return fail(err, "unable to create user")
	}

	// Log the user in.
	err = h.auth.Login(ctx, u.ID)
	if err != nil {
		log.Ctx(ctx).Error("unable to log user in",
			"error", err,
			"user_id", u.ID,
		)
		msg.Info(ctx, "Your account has been created.")
		return redirect.New(ctx).
			Route(routenames.Login).
			Go()
	}

	// Get the verification code for display
	displayCode := ""
	if u.VerificationCode != nil {
		displayCode = *u.VerificationCode
	}

	msg.Success(ctx, fmt.Sprintf("🎉 Account created! Your WhatsApp verification code is: %s. Check your WhatsApp and select the button with this code to complete verification.", displayCode))

	// Send the phone verification.
	h.sendPhoneVerification(ctx, u, "web")

	return redirect.New(ctx).
		Route(routenames.Home).
		Go()
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

	// Update the user.
	_, err = usr.
		Update().
		SetPassword(input.Password).
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
	return redirect.New(ctx).
		Route(routenames.Home).
		Go()
}

func (h *Auth) VerificationNotice(ctx echo.Context) error {
	return pages.VerificationNotice(ctx)
}

func (h *Auth) ResendVerification(ctx echo.Context) error {
	if u := ctx.Get(context.AuthenticatedUserKey); u != nil {
		if user, ok := u.(*ent.User); ok {
			if user.Verified {
				msg.Info(ctx, "Your phone number is already verified.")
				return redirect.New(ctx).Route(routenames.Home).Go()
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
