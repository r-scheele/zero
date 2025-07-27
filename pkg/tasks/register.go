package tasks

import (
	"github.com/r-scheele/zero/pkg/services"
)

// Register registers all task queues with the task client.
func Register(c *services.Container) {
	c.Tasks.Register(NewExampleTaskQueue(c))
	c.Tasks.Register(NewPhoneVerificationTaskQueue(c))
	c.Tasks.Register(NewPasswordResetTaskQueue(c))
}
