package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/ent/user"
	"github.com/r-scheele/zero/pkg/log"
	"github.com/r-scheele/zero/pkg/services"
)

type WhatsAppWebhook struct {
	Container *services.Container
}

// 360dialog webhook payload structures
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

// VerifyWebhook handles webhook verification from 360dialog
func (h *WhatsAppWebhook) VerifyWebhook(c echo.Context) error {
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
func (h *WhatsAppWebhook) HandleWebhook(c echo.Context) error {
	var payload WebhookPayload
	if err := json.NewDecoder(c.Request().Body).Decode(&payload); err != nil {
		log.Default().Error("Failed to decode webhook payload", "error", err)
		return c.String(http.StatusBadRequest, "Invalid payload")
	}

	ctx := c.Request().Context()

	// Process incoming messages
	for _, message := range payload.Messages {
		if err := h.processMessage(ctx, message); err != nil {
			log.Default().Error("Failed to process message", "error", err, "message_id", message.ID)
		}
	}

	// Process message statuses (delivered, read, etc.)
	for _, status := range payload.Statuses {
		log.Default().Info("Message status update",
			"message_id", status.ID,
			"status", status.Status,
			"recipient", status.RecipientID)
	}

	return c.String(http.StatusOK, "OK")
}

// processMessage handles individual WhatsApp messages
func (h *WhatsAppWebhook) processMessage(ctx context.Context, msg Message) error {
	phoneNumber := msg.From

	// Handle button responses
	if msg.Button != nil {
		return h.handleButtonResponse(ctx, phoneNumber, msg.Button)
	}

	// Handle text messages
	if msg.Text != nil {
		return h.handleTextMessage(ctx, phoneNumber, msg.Text.Body)
	}

	return nil
}

// handleButtonResponse processes button click responses
func (h *WhatsAppWebhook) handleButtonResponse(ctx context.Context, phoneNumber string, button *ButtonMessage) error {
	payload := button.Payload

	log.Default().Info("Button response received",
		"phone", phoneNumber,
		"payload", payload,
		"text", button.Text)

	// Handle password reset button
	if strings.HasPrefix(payload, "reset_password_") {
		token := strings.TrimPrefix(payload, "reset_password_")
		return h.handlePasswordReset(ctx, phoneNumber, token)
	}

	// Handle help reset button
	if payload == "help_reset" {
		return h.sendPasswordResetHelp(ctx, phoneNumber)
	}

	// Handle verification button
	if strings.HasPrefix(payload, "verify_") {
		payloadParts := strings.TrimPrefix(payload, "verify_")
		// New format: verify_{token}_{code} or old format: verify_{token}
		parts := strings.Split(payloadParts, "_")

		if len(parts) == 2 {
			// New format with verification code
			token := parts[0]
			selectedCode := parts[1]
			return h.verifyUserWithTokenAndCode(ctx, phoneNumber, token, selectedCode)
		} else {
			// Old format (backwards compatibility)
			token := payloadParts
			return h.verifyUserWithToken(ctx, phoneNumber, token)
		}
	}

	// Handle help button
	if payload == "help_verification" {
		return h.sendHelpMessage(ctx, phoneNumber)
	}

	// Handle login web button
	if payload == "login_web" {
		return h.sendLoginInstructions(ctx, phoneNumber)
	}

	// Handle get started button
	if payload == "get_started" {
		return h.sendGettingStartedMessage(ctx, phoneNumber)
	}

	return nil
}

// handleTextMessage processes text messages for registration
func (h *WhatsAppWebhook) handleTextMessage(ctx context.Context, phoneNumber, text string) error {
	originalText := strings.TrimSpace(text)
	text = strings.ToLower(originalText)

	// Check for direct password setting (case-sensitive check on original text)
	if strings.HasPrefix(strings.ToUpper(originalText), "NEW PASSWORD:") {
		return h.handleDirectPasswordReset(ctx, phoneNumber, originalText)
	}

	// Check if user already exists
	existingUser, err := h.Container.ORM.User.Query().
		Where(user.PhoneNumber(phoneNumber)).
		First(ctx)

	if err == nil {
		// User exists, send welcome back message
		return h.Container.API.SendWhatsAppMessage(ctx, phoneNumber,
			fmt.Sprintf("Hi %s! You're already registered. You can log in to our web platform using your phone number.", existingUser.Name))
	}

	// New user - check if they want to register
	if strings.Contains(text, "register") || strings.Contains(text, "sign up") || strings.Contains(text, "join") {
		return h.registerUserFromWhatsApp(ctx, phoneNumber, text)
	}

	// Send welcome/help message
	return h.sendWelcomeMessage(ctx, phoneNumber)
}

// verifyUserWithToken verifies a user using the verification token
func (h *WhatsAppWebhook) verifyUserWithToken(ctx context.Context, phoneNumber, token string) error {
	// Verify the token - this returns the phone number from the token
	tokenPhoneNumber, err := h.Container.Auth.ValidatePhoneVerificationToken(token)
	if err != nil {
		return h.Container.API.SendWhatsAppMessage(ctx, phoneNumber,
			"❌ Invalid or expired verification code. Please try registering again on our website.")
	}

	// Verify that the phone number matches
	if tokenPhoneNumber != phoneNumber {
		return h.Container.API.SendWhatsAppMessage(ctx, phoneNumber,
			"❌ Phone number mismatch. Please use the same phone number you registered with.")
	}

	// Find and update user as verified
	userResult, err := h.Container.ORM.User.Query().
		Where(user.PhoneNumber(phoneNumber)).
		First(ctx)
	if err != nil {
		return h.Container.API.SendWhatsAppMessage(ctx, phoneNumber,
			"❌ User not found. Please register first on our website.")
	}

	// Update user as verified
	_, err = userResult.Update().
		SetVerified(true).
		Save(ctx)
	if err != nil {
		return err
	}

	// Send success message
	return h.Container.API.SendWhatsAppMessage(ctx, phoneNumber,
		fmt.Sprintf("🎉 Congratulations! Your account has been verified successfully.\n\nYou can now log in to our web platform using your phone number: %s", phoneNumber))
}

// verifyUserWithTokenAndCode verifies a user using the verification token and selected code
func (h *WhatsAppWebhook) verifyUserWithTokenAndCode(ctx context.Context, phoneNumber, token, selectedCode string) error {
	// Verify the token - this returns the phone number from the token
	tokenPhoneNumber, err := h.Container.Auth.ValidatePhoneVerificationToken(token)
	if err != nil {
		return h.Container.API.SendWhatsAppMessage(ctx, phoneNumber,
			"❌ Invalid or expired verification code. Please try registering again on our website.")
	}

	// Verify that the phone number matches
	if tokenPhoneNumber != phoneNumber {
		return h.Container.API.SendWhatsAppMessage(ctx, phoneNumber,
			"❌ Phone number mismatch. Please use the same phone number you registered with.")
	}

	// Find user and check verification code
	userResult, err := h.Container.ORM.User.Query().
		Where(user.PhoneNumber(phoneNumber)).
		First(ctx)
	if err != nil {
		return h.Container.API.SendWhatsAppMessage(ctx, phoneNumber,
			"❌ User not found. Please register first on our website.")
	}

	// Check if the selected code matches the stored verification code
	if userResult.VerificationCode == nil || *userResult.VerificationCode != selectedCode {
		return h.Container.API.SendWhatsAppMessage(ctx, phoneNumber,
			"❌ Incorrect verification code selected. Please try again with the code shown on your registration page.")
	}

	// Update user as verified and clear the verification code
	_, err = userResult.Update().
		SetVerified(true).
		ClearVerificationCode().
		Save(ctx)
	if err != nil {
		return err
	}

	// Send success message
	return h.Container.API.SendWhatsAppMessage(ctx, phoneNumber,
		fmt.Sprintf("🎉 Congratulations! Your account has been verified successfully.\n\nYou can now log in to our web platform using your phone number: %s", phoneNumber))
}

// registerUserFromWhatsApp registers a new user who contacted via WhatsApp
func (h *WhatsAppWebhook) registerUserFromWhatsApp(ctx context.Context, phoneNumber, message string) error {
	// Extract name from message or use phone number
	name := extractNameFromMessage(message)
	if name == "" {
		name = "WhatsApp User"
	}

	// Create user with WhatsApp registration method (auto-verified)
	user, err := h.Container.ORM.User.Create().
		SetPhoneNumber(phoneNumber).
		SetName(name).
		SetVerified(true). // Auto-verify WhatsApp users
		SetRegistrationMethod(user.RegistrationMethodWhatsapp).
		Save(ctx)

	if err != nil {
		return h.Container.API.SendWhatsAppMessage(ctx, phoneNumber,
			"❌ Sorry, there was an error creating your account. Please try again later.")
	}

	log.Default().Info("User registered via WhatsApp",
		"user_id", user.ID,
		"phone", phoneNumber,
		"name", name)

	// Send welcome message with account details
	return h.Container.API.SendWelcomeMessage(ctx, phoneNumber, user.Name)
}

// Helper functions
func extractNameFromMessage(message string) string {
	// Simple name extraction logic
	words := strings.Fields(message)
	for i, word := range words {
		if strings.ToLower(word) == "register" || strings.ToLower(word) == "name" {
			if i+1 < len(words) {
				return strings.Title(words[i+1])
			}
		}
	}
	return ""
}

func (h *WhatsAppWebhook) sendHelpMessage(ctx context.Context, phoneNumber string) error {
	helpText := `❓ Need help with verification?

1. Make sure you clicked the "✅ Verify Account" button
2. Check that the verification code matches
3. Ensure you're using the same phone number you registered with

If you're still having issues, you can:
• Try registering again on our website
• Contact our support team

Type "register" to create a new account directly on WhatsApp.`

	return h.Container.API.SendWhatsAppMessage(ctx, phoneNumber, helpText)
}

func (h *WhatsAppWebhook) sendLoginInstructions(ctx context.Context, phoneNumber string) error {
	loginText := fmt.Sprintf(`🌐 To log in to our web platform:

1. Go to our website login page
2. Enter your phone number: %s
3. Enter your password (if you set one)
4. Click "Sign In"

If you registered via WhatsApp and don't have a password, you can set one after logging in with a verification code.`, phoneNumber)

	return h.Container.API.SendWhatsAppMessage(ctx, phoneNumber, loginText)
}

func (h *WhatsAppWebhook) sendGettingStartedMessage(ctx context.Context, phoneNumber string) error {
	startText := `🚀 Getting Started with Zero:

1. 🌐 Log in to our web platform
2. 📝 Complete your profile
3. 🔧 Explore our features
4. 💬 Join our community

You can always message us here on WhatsApp for quick support!

What would you like to do next?`

	return h.Container.API.SendWhatsAppMessage(ctx, phoneNumber, startText)
}

func (h *WhatsAppWebhook) sendWelcomeMessage(ctx context.Context, phoneNumber string) error {
	welcomeText := `👋 Welcome to Zero!

I can help you:
• Register a new account (just type "register")
• Get login instructions
• Answer questions about our platform

Type "register" to create an account, or let me know how I can help you!`

	return h.Container.API.SendWhatsAppMessage(ctx, phoneNumber, welcomeText)
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// handlePasswordReset processes password reset requests
func (h *WhatsAppWebhook) handlePasswordReset(ctx context.Context, phoneNumber, token string) error {
	// Validate the password reset token
	_, tokenPhone, err := h.Container.Auth.ValidateWhatsAppPasswordResetToken(token)
	if err != nil {
		return h.Container.API.SendWhatsAppMessage(ctx, phoneNumber,
			"❌ Invalid or expired password reset link. Please request a new password reset from our website.")
	}

	// Verify phone number matches
	if tokenPhone != phoneNumber {
		return h.Container.API.SendWhatsAppMessage(ctx, phoneNumber,
			"❌ Phone number mismatch. Please use the same phone number you registered with.")
	}

	// Send password setup instructions
	return h.Container.API.SendWhatsAppMessage(ctx, phoneNumber,
		`🔐 Password Reset Confirmed!

To complete your password reset:

1. Go to our website login page
2. Click "Forgot Password"
3. Enter your phone number: `+phoneNumber+`
4. Use the web form to set your new password

Or simply create a new password and we'll update it for you. Reply with:
"NEW PASSWORD: your_new_password_here"

⚠️ This reset session expires in 1 hour for security.`)
}

// sendPasswordResetHelp sends help information for password reset
func (h *WhatsAppWebhook) sendPasswordResetHelp(ctx context.Context, phoneNumber string) error {
	helpText := `🔐 Password Reset Help:

Having trouble resetting your password?

1. Make sure you clicked the "🔐 Reset Password" button
2. Ensure you're using the same phone number you registered with
3. Check that the reset link hasn't expired (valid for 1 hour)

Alternative methods:
• Visit our website and use "Forgot Password"
• Contact our support team
• Reply "NEW PASSWORD: your_password" to set it directly

Need more help? Just ask!`

	return h.Container.API.SendWhatsAppMessage(ctx, phoneNumber, helpText)
}

// handleDirectPasswordReset handles password reset via direct WhatsApp message
func (h *WhatsAppWebhook) handleDirectPasswordReset(ctx context.Context, phoneNumber, message string) error {
	// Extract password from message: "NEW PASSWORD: mypassword123"
	parts := strings.SplitN(message, ":", 2)
	if len(parts) != 2 {
		return h.Container.API.SendWhatsAppMessage(ctx, phoneNumber,
			`❌ Invalid format. Please use:
"NEW PASSWORD: your_new_password_here"

Example:
NEW PASSWORD: MySecurePass123`)
	}

	newPassword := strings.TrimSpace(parts[1])
	if len(newPassword) < 8 {
		return h.Container.API.SendWhatsAppMessage(ctx, phoneNumber,
			"❌ Password must be at least 8 characters long. Please try again with a stronger password.")
	}

	// Find user by phone number
	userResult, err := h.Container.ORM.User.Query().
		Where(user.PhoneNumber(phoneNumber)).
		First(ctx)
	if err != nil {
		return h.Container.API.SendWhatsAppMessage(ctx, phoneNumber,
			"❌ User not found. Please make sure you're using the registered phone number.")
	}

	// Update user password (password will be automatically hashed by the User schema hook)
	_, err = userResult.Update().
		SetPassword(newPassword).
		Save(ctx)
	if err != nil {
		return h.Container.API.SendWhatsAppMessage(ctx, phoneNumber,
			"❌ Error updating password. Please try again.")
	}

	log.Default().Info("Password reset via WhatsApp",
		"user_id", userResult.ID,
		"phone", phoneNumber)

	// Send success message
	return h.Container.API.SendWhatsAppMessage(ctx, phoneNumber,
		`✅ Password Updated Successfully!

Your password has been changed. You can now log in to our website using:
📱 Phone: `+phoneNumber+`
🔐 Your new password

For security, this conversation will be cleared. Your account is now ready to use!`)
}
