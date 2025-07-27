package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/ent"
	"github.com/r-scheele/zero/pkg/services"
	"github.com/r-scheele/zero/pkg/tasks"
)

// SendPhoneVerification generates a verification token and queues a phone verification task
func SendPhoneVerification(ctx echo.Context, container *services.Container, user *ent.User, method string) error {
	var token string
	var err error

	if method == "web" {
		// For web registration, generate a simple code for WhatsApp confirmation
		token, err = container.Auth.GenerateSimpleVerificationCode()
	} else {
		// For WhatsApp registration, generate a JWT token for web login
		token, err = container.Auth.GeneratePhoneVerificationToken(user.PhoneNumber)
	}

	if err != nil {
		return err
	}

	// Get verification code or empty string if nil
	verificationCode := ""
	if user.VerificationCode != nil {
		verificationCode = *user.VerificationCode
	}

	// Create and queue the phone verification task
	task := tasks.PhoneVerificationTask{
		UserID:           user.ID,
		PhoneNumber:      user.PhoneNumber,
		Username:         user.Name,
		Token:            token,
		Method:           method,
		VerificationCode: verificationCode,
	}

	// Queue the task
	return container.Tasks.Add(task).Save()
}
