package handlers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/pkg/form"
	"github.com/r-scheele/zero/pkg/routenames"
	"github.com/r-scheele/zero/pkg/services"
	"github.com/r-scheele/zero/pkg/ui/forms"
	"github.com/r-scheele/zero/pkg/ui/pages"
)

type Contact struct {
	mail *services.MailClient
}

func init() {
	Register(new(Contact))
}

func (h *Contact) Init(c *services.Container) error {
	h.mail = c.Mail
	return nil
}

func (h *Contact) Routes(g *echo.Group) {
	contact := g.Group("/contact")
	contact.GET("", h.Page).Name = routenames.Contact
	contact.POST("", h.Submit).Name = routenames.ContactSubmit
}

func (h *Contact) Page(ctx echo.Context) error {
	return pages.ContactUs(ctx, form.Get[forms.Contact](ctx))
}

func (h *Contact) Submit(ctx echo.Context) error {
	var input forms.Contact

	err := form.Submit(ctx, &input)

	switch err.(type) {
	case nil:
	case validator.ValidationErrors:
		return h.Page(ctx)
	default:
		return err
	}

	err = h.mail.
		Compose().
		To(input.Email).
		Subject("Contact form submitted").
		Body(fmt.Sprintf("The message is: %s", input.Message)).
		Send(ctx)

	if err != nil {
		return fail(err, "unable to send email")
	}

	return h.Page(ctx)
}
