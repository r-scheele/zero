package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/mikestefanello/backlite"
	"github.com/r-scheele/zero/ent"
	"github.com/spf13/afero"
)

type APIService struct {
	orm      *ent.Client
	auth     *AuthClient
	cache    *CacheClient
	files    afero.Fs
	mail     *MailClient
	tasks    *backlite.Client
	whatsapp *WhatsAppAPI
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

type ProcessMessageRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
	MessageID   string `json:"message_id" validate:"required"`
	Content     string `json:"content" validate:"required"`
	MessageType string `json:"message_type" validate:"required,oneof=text image voice document"`
}

type ProcessMessageResponse struct {
	Success    bool   `json:"success"`
	Response   string `json:"response"`
	MessageID  string `json:"message_id"`
	ActionType string `json:"action_type"` // "reply", "verification", "help"
}

type UserStatsResponse struct {
	PhoneNumber  string    `json:"phone_number"`
	MessageCount int       `json:"message_count"`
	LastActive   time.Time `json:"last_active"`
}

func NewAPIService(orm *ent.Client, auth *AuthClient, cache *CacheClient, files afero.Fs, mail *MailClient, tasks *backlite.Client) *APIService {
	// Initialize WhatsApp API with environment variables
	whatsappAPI := &WhatsAppAPI{
		baseURL: "https://waba.360dialog.io/v1",
		apiKey:  getEnvWithDefault("WHATSAPP_API_KEY", ""),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	return &APIService{
		orm:      orm,
		auth:     auth,
		cache:    cache,
		files:    files,
		mail:     mail,
		tasks:    tasks,
		whatsapp: whatsappAPI,
	}
}

// getEnvWithDefault gets environment variable with default fallback
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
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
	buttonCodes := s.generateThreeButtonCodes(verificationCode)

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

func (s *APIService) ProcessMessage(ctx context.Context, req ProcessMessageRequest) (*ProcessMessageResponse, error) {
	// For now, just return a basic response until we have the database schemas
	response := s.generateResponse(req.Content, req.MessageType)
	response.Success = true
	response.MessageID = fmt.Sprintf("resp_%s_%d", req.MessageID, time.Now().Unix())

	return response, nil
}

func (s *APIService) GetUserStats(ctx context.Context, phoneNumber string) (*UserStatsResponse, error) {
	// Return basic stats until we have the database schemas
	return &UserStatsResponse{
		PhoneNumber:  phoneNumber,
		MessageCount: 0,
		LastActive:   time.Now(),
	}, nil
}

func (s *APIService) GetConversationHistory(ctx context.Context, phoneNumber string, limit int) ([]map[string]interface{}, error) {
	// Return empty history until we have the database schemas
	return []map[string]interface{}{}, nil
}

func (s *APIService) generateResponse(content, messageType string) *ProcessMessageResponse {
	response := &ProcessMessageResponse{}

	switch {
	case content == "/help":
		response.Response = `🤖 Zero Bot Commands:

/help - Show this help
/verify - Verify your account

Welcome to Zero! Send me a message to get started. 📱`
		response.ActionType = "reply"

	default:
		// Basic echo for now - will be enhanced with AI processing
		response.Response = fmt.Sprintf("I received your %s message. Processing...", messageType)
		response.ActionType = "reply"
	}

	return response
}

// generateTwoDigitCode generates a random 2-digit verification code (10-99)
func (s *APIService) generateTwoDigitCode() string {
	code := rand.Intn(90) + 10 // generates numbers from 10 to 99
	return fmt.Sprintf("%02d", code)
}

// generateThreeButtonCodes generates 3 different 2-digit codes including the correct one
func (s *APIService) generateThreeButtonCodes(correctCode string) []string {
	codes := []string{correctCode}
	used := map[string]bool{correctCode: true}

	for len(codes) < 3 {
		newCode := s.generateTwoDigitCode()
		if !used[newCode] {
			codes = append(codes, newCode)
			used[newCode] = true
		}
	}

	// Shuffle the codes so the correct one isn't always first
	for i := len(codes) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		codes[i], codes[j] = codes[j], codes[i]
	}

	return codes
}
