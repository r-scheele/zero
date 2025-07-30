package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/r-scheele/zero/config"
	"github.com/r-scheele/zero/ent"
	"github.com/spf13/afero"
)

type APIService struct {
	whatsapp *WhatsAppAPI
	User     *UserService
	Admin    *AdminService
	File     *FileService
	Contact  *ContactService
	JWT      *JWTService
}

// WhatsAppAPI handles 360dialog WhatsApp integration
type WhatsAppAPI struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

// 360dialog API structures
type WhatsAppMessage struct {
	To          string              `json:"to"`
	Type        string              `json:"type"`
	Interactive *InteractiveMessage `json:"interactive,omitempty"`
	Text        *TextMessage        `json:"text,omitempty"`
}

type InteractiveMessage struct {
	Type   string  `json:"type"`
	Header *Header `json:"header,omitempty"`
	Body   *Body   `json:"body"`
	Footer *Footer `json:"footer,omitempty"`
	Action *Action `json:"action"`
}

type TextMessage struct {
	Body string `json:"body"`
}

type Header struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Body struct {
	Text string `json:"text"`
}

type Footer struct {
	Text string `json:"text"`
}

type Action struct {
	Buttons []Button `json:"buttons"`
}

type Button struct {
	Type  string      `json:"type"`
	Reply ReplyButton `json:"reply"`
}

type ReplyButton struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}


func NewAPIService(orm *ent.Client, auth *AuthClient, mail *MailClient, files afero.Fs, config *config.Config) *APIService {
	// Initialize WhatsApp API with configuration values
	whatsappAPI := &WhatsAppAPI{
		baseURL: config.WhatsApp.BaseURL,
		apiKey:  config.WhatsApp.AccessToken,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	return &APIService{
		whatsapp: whatsappAPI,
		User:     NewUserService(orm, auth),
		Admin:    NewAdminService(orm),
		File:     NewFileService(files),
		Contact:  NewContactService(mail),
		JWT:      NewJWTService(config),
	}
}



// SendWhatsAppMessage sends a simple text message via 360dialog
func (s *APIService) SendWhatsAppMessage(ctx context.Context, phoneNumber, message string) error {
	if s.whatsapp.apiKey == "" {
		// Fallback to console logging if no API key
		fmt.Printf("📱 WhatsApp Message to %s: %s\n", phoneNumber, message)
		return nil
	}

	whatsappMsg := WhatsAppMessage{
		To:   phoneNumber,
		Type: "text",
		Text: &TextMessage{
			Body: message,
		},
	}

	return s.whatsapp.sendMessage(ctx, whatsappMsg)
}

// SendVerificationMessage sends a WhatsApp verification message with buttons
func (s *APIService) SendVerificationMessage(ctx context.Context, phoneNumber, username, token, verificationCode string) error {
	if s.whatsapp.apiKey == "" {
		// Fallback to console logging if no API key
		fmt.Printf("📱 WhatsApp Verification to %s: Hi %s! Your verification code is: %s\n", phoneNumber, username, verificationCode)
		return nil
	}

	// Generate 3 different codes including the correct one
	buttonCodes := []string{verificationCode}
	used := map[string]bool{verificationCode: true}
	
	for len(buttonCodes) < 3 {
		newCode := s.generateTwoDigitCode()
		if !used[newCode] {
			buttonCodes = append(buttonCodes, newCode)
			used[newCode] = true
		}
	}
	
	// Shuffle the codes so the correct one isn't always first
	for i := len(buttonCodes) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		buttonCodes[i], buttonCodes[j] = buttonCodes[j], buttonCodes[i]
	}

	message := WhatsAppMessage{
		To:   phoneNumber,
		Type: "interactive",
		Interactive: &InteractiveMessage{
			Type: "button",
			Header: &Header{
				Type: "text",
				Text: "📱 Account Verification",
			},
			Body: &Body{
				Text: fmt.Sprintf("Hi %s! Welcome to Zero.\n\nTo complete your registration, please verify your account by selecting the 2-digit code shown on your registration page.\n\n🔐 Choose the correct verification code:", username),
			},
			Footer: &Footer{
				Text: "This verification will expire in 15 minutes.",
			},
			Action: &Action{
				Buttons: []Button{
					{
						Type: "reply",
						Reply: ReplyButton{
							ID:    fmt.Sprintf("verify_%s_%s", token, buttonCodes[0]),
							Title: fmt.Sprintf("🔢 %s", buttonCodes[0]),
						},
					},
					{
						Type: "reply",
						Reply: ReplyButton{
							ID:    fmt.Sprintf("verify_%s_%s", token, buttonCodes[1]),
							Title: fmt.Sprintf("🔢 %s", buttonCodes[1]),
						},
					},
					{
						Type: "reply",
						Reply: ReplyButton{
							ID:    fmt.Sprintf("verify_%s_%s", token, buttonCodes[2]),
							Title: fmt.Sprintf("🔢 %s", buttonCodes[2]),
						},
					},
				},
			},
		},
	}

	return s.whatsapp.sendMessage(ctx, message)
}

// SendWelcomeMessage sends a welcome message for WhatsApp registrations
func (s *APIService) SendWelcomeMessage(ctx context.Context, phoneNumber, username string) error {
	if s.whatsapp.apiKey == "" {
		// Fallback to console logging if no API key
		fmt.Printf("📱 WhatsApp Welcome to %s: Welcome %s! Your account is ready.\n", phoneNumber, username)
		return nil
	}

	message := WhatsAppMessage{
		To:   phoneNumber,
		Type: "interactive",
		Interactive: &InteractiveMessage{
			Type: "button",
			Body: &Body{
				Text: fmt.Sprintf("🎉 Welcome to Zero, %s!\n\nYour account has been verified successfully.\n\nYou can now log in to our web platform using your phone number: %s", username, phoneNumber),
			},
			Action: &Action{
				Buttons: []Button{
					{
						Type: "reply",
						Reply: ReplyButton{
							ID:    "login_web",
							Title: "🌐 Login to Web",
						},
					},
					{
						Type: "reply",
						Reply: ReplyButton{
							ID:    "get_started",
							Title: "🚀 Get Started",
						},
					},
				},
			},
		},
	}

	return s.whatsapp.sendMessage(ctx, message)
}

// SendPasswordResetMessage sends a WhatsApp password reset message with buttons
func (s *APIService) SendPasswordResetMessage(ctx context.Context, phoneNumber, username, resetToken string) error {
	if s.whatsapp.apiKey == "" {
		// Fallback to console logging if no API key
		fmt.Printf("📱 WhatsApp Password Reset to %s: Hi %s! Click the link to reset your password: %s\n", phoneNumber, username, resetToken)
		return nil
	}

	message := WhatsAppMessage{
		To:   phoneNumber,
		Type: "interactive",
		Interactive: &InteractiveMessage{
			Type: "button",
			Header: &Header{
				Type: "text",
				Text: "🔐 Password Reset",
			},
			Body: &Body{
				Text: fmt.Sprintf("Hi %s!\n\nYou requested a password reset for your Zero account.\n\nClick the button below to set a new password securely.", username),
			},
			Footer: &Footer{
				Text: "This reset link will expire in 1 hour.",
			},
			Action: &Action{
				Buttons: []Button{
					{
						Type: "reply",
						Reply: ReplyButton{
							ID:    fmt.Sprintf("reset_password_%s", resetToken),
							Title: "🔐 Reset Password",
						},
					},
					{
						Type: "reply",
						Reply: ReplyButton{
							ID:    "help_reset",
							Title: "❓ Need Help?",
						},
					},
				},
			},
		},
	}

	return s.whatsapp.sendMessage(ctx, message)
}

// sendMessage sends the actual HTTP request to 360dialog API
func (w *WhatsAppAPI) sendMessage(ctx context.Context, message WhatsAppMessage) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", w.baseURL+"/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("D360-API-KEY", w.apiKey)

	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("WhatsApp API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Success - log the response
	fmt.Printf("📱 WhatsApp message sent successfully to %s\n", message.To)
	return nil
}

// Removed placeholder functions that are not currently implemented

// generateTwoDigitCode generates a random 2-digit verification code (10-99)
func (s *APIService) generateTwoDigitCode() string {
	code := rand.Intn(90) + 10 // generates numbers from 10 to 99
	return fmt.Sprintf("%02d", code)
}

// Removed unused generateThreeButtonCodes function
