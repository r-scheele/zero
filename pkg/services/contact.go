package services

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/pkg/ui/forms"
)

// ContactService handles contact form operations
type ContactService struct {
	mail *MailClient
}

// NewContactService creates a new contact service
func NewContactService(mail *MailClient) *ContactService {
	return &ContactService{
		mail: mail,
	}
}

// SubmitContactForm processes a contact form submission
func (s *ContactService) SubmitContactForm(ctx echo.Context, input forms.Contact) error {
	// Send email
	err := s.mail.
		Compose().
		To(input.Email).
		Subject("Contact form submitted").
		Body(fmt.Sprintf("The message is: %s", input.Message)).
		Send(ctx)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}