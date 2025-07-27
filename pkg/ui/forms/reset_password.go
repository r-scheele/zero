package forms

import (
	"net/http"

	"github.com/r-scheele/zero/pkg/form"
	"github.com/r-scheele/zero/pkg/ui"
	. "github.com/r-scheele/zero/pkg/ui/components"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type ResetPassword struct {
	Password        string `form:"password" validate:"required,min=8"`
	ConfirmPassword string `form:"password-confirm" validate:"required,eqfield=Password"`
	form.Submission
}

func (f *ResetPassword) Render(r *ui.Request) Node {
	return Form(
		ID("reset-password"),
		Method(http.MethodPost),
		HxBoost(),
		Action(r.CurrentPath),
		Class("space-y-6"),
		InputField(InputFieldParams{
			Form:        f,
			FormField:   "Password",
			Name:        "password",
			InputType:   "password",
			Label:       "New Password",
			Placeholder: "••••••••",
		}),
		InputField(InputFieldParams{
			Form:        f,
			FormField:   "ConfirmPassword",
			Name:        "password-confirm",
			InputType:   "password",
			Label:       "Confirm password",
			Placeholder: "••••••••",
		}),
		Div(
			Class("pt-4"),
			Button(
				Class("btn-modern-primary w-full"),
				Type("submit"),
				Text("Update Password"),
			),
		),
		CSRF(r),
	)
}
