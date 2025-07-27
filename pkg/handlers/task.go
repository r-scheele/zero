package handlers

import (
	"fmt"
	"time"

	"github.com/mikestefanello/backlite"
	"github.com/r-scheele/zero/pkg/middleware"
	"github.com/r-scheele/zero/pkg/msg"
	"github.com/r-scheele/zero/pkg/routenames"
	"github.com/r-scheele/zero/pkg/ui/forms"
	"github.com/r-scheele/zero/pkg/ui/pages"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/pkg/form"
	"github.com/r-scheele/zero/pkg/services"
	"github.com/r-scheele/zero/pkg/tasks"
)

type Task struct {
	tasks *backlite.Client
}

func init() {
	Register(new(Task))
}

func (h *Task) Init(c *services.Container) error {
	h.tasks = c.Tasks
	return nil
}

func (h *Task) Routes(g *echo.Group) {
	task := g.Group("/task", middleware.RequireAuthentication, middleware.RequireVerification)
	task.GET("", h.Page).Name = routenames.Task
	task.POST("", h.Submit).Name = routenames.TaskSubmit
}

// Page renders the task creation page.
func (h *Task) Page(ctx echo.Context) error {
	return pages.AddTask(ctx, form.Get[forms.Task](ctx))
}

// Submit handles the task creation form submission.
func (h *Task) Submit(ctx echo.Context) error {
	var input forms.Task

	err := form.Submit(ctx, &input)

	switch err.(type) {
	case nil:
	case validator.ValidationErrors:
		return h.Page(ctx)
	default:
		return err
	}

	// Insert the task
	err = h.tasks.
		Add(tasks.ExampleTask{
			Message: input.Message,
		}).
		Wait(time.Duration(input.Delay) * time.Second).
		Save()

	if err != nil {
		return fail(err, "unable to create a task")
	}

	msg.Success(ctx, fmt.Sprintf("The task has been created. Check the logs in %d seconds.", input.Delay))
	form.Clear(ctx)

	return h.Page(ctx)
}
