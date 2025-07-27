package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/mikestefanello/backlite"
	"github.com/r-scheele/zero/pkg/log"
	"github.com/r-scheele/zero/pkg/services"
)

// PhoneVerificationTask represents a task to send phone verification via WhatsApp
type PhoneVerificationTask struct {
	UserID           int    `json:"user_id"`
	PhoneNumber      string `json:"phone_number"`
	Username         string `json:"username"`
	Token            string `json:"token"`
	Method           string `json:"method"`            // "web" or "whatsapp"
	VerificationCode string `json:"verification_code"` // 2-digit code for button verification
}

// Config satisfies the backlite.Task interface by providing configuration for the queue
func (t PhoneVerificationTask) Config() backlite.QueueConfig {
	return backlite.QueueConfig{
		Name:        "PhoneVerificationTask",
		MaxAttempts: 3,
		Timeout:     30 * time.Second,
		Backoff:     30 * time.Second,
		Retention: &backlite.Retention{
			Duration:   24 * time.Hour,
			OnlyFailed: false,
			Data: &backlite.RetainData{
				OnlyFailed: false,
			},
		},
	}
}

// NewPhoneVerificationTaskQueue provides a Queue that can process PhoneVerificationTask tasks
func NewPhoneVerificationTaskQueue(c *services.Container) backlite.Queue {
	return backlite.NewQueue[PhoneVerificationTask](func(ctx context.Context, task PhoneVerificationTask) error {
		log.Default().Info("Processing phone verification task",
			"user_id", task.UserID,
			"phone_number", task.PhoneNumber,
			"method", task.Method,
		)

		var err error
		if task.Method == "web" {
			// User registered on web, send verification message with buttons
			err = c.API.SendVerificationMessage(ctx, task.PhoneNumber, task.Username, task.Token, task.VerificationCode)
		} else {
			// User registered via WhatsApp, send welcome message with buttons
			err = c.API.SendWelcomeMessage(ctx, task.PhoneNumber, task.Username)
		}

		if err != nil {
			log.Default().Error("Failed to send WhatsApp verification",
				"user_id", task.UserID,
				"phone_number", task.PhoneNumber,
				"error", err,
			)
			return fmt.Errorf("failed to send WhatsApp verification: %w", err)
		}

		log.Default().Info("Phone verification sent successfully via WhatsApp",
			"user_id", task.UserID,
			"phone_number", task.PhoneNumber,
			"method", task.Method,
		)

		return nil
	})
}
