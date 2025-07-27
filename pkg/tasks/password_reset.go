package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/mikestefanello/backlite"
	"github.com/r-scheele/zero/pkg/log"
	"github.com/r-scheele/zero/pkg/services"
)

// PasswordResetTask represents a task to send password reset via WhatsApp
type PasswordResetTask struct {
	UserID      int    `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
	Username    string `json:"username"`
	ResetToken  string `json:"reset_token"`
}

// Config satisfies the backlite.Task interface by providing configuration for the queue
func (t PasswordResetTask) Config() backlite.QueueConfig {
	return backlite.QueueConfig{
		Name:        "PasswordResetTask",
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

// NewPasswordResetTaskQueue provides a Queue that can process PasswordResetTask tasks
func NewPasswordResetTaskQueue(c *services.Container) backlite.Queue {
	return backlite.NewQueue[PasswordResetTask](func(ctx context.Context, task PasswordResetTask) error {
		log.Default().Info("Processing password reset task",
			"user_id", task.UserID,
			"phone_number", task.PhoneNumber,
		)

		// Send WhatsApp password reset message with buttons
		err := c.API.SendPasswordResetMessage(ctx, task.PhoneNumber, task.Username, task.ResetToken)
		if err != nil {
			log.Default().Error("Failed to send WhatsApp password reset message",
				"user_id", task.UserID,
				"phone_number", task.PhoneNumber,
				"error", err,
			)
			return fmt.Errorf("failed to send WhatsApp password reset: %w", err)
		}

		log.Default().Info("Password reset sent successfully via WhatsApp",
			"user_id", task.UserID,
			"phone_number", task.PhoneNumber,
		)

		return nil
	})
}
