package forms

import (
	"net/http"

	"github.com/r-scheele/zero/pkg/form"
	"github.com/r-scheele/zero/pkg/routenames"
	"github.com/r-scheele/zero/pkg/ui"
	. "github.com/r-scheele/zero/pkg/ui/components"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type Register struct {
	Name            string `form:"name" validate:"required"`
	PhoneNumber     string `form:"phone_number" validate:"required,e164"`
	Password        string `form:"password" validate:"required,min=8"` 
	ConfirmPassword string `form:"password-confirm" validate:"required,eqfield=Password"`
	form.Submission
}

func (f *Register) Render(r *ui.Request) Node {
	return Form(
		ID("register"),
		Method(http.MethodPost),
		HxBoost(),
		Action(r.Path(routenames.RegisterSubmit)),
		Class("space-y-4"), // Reduced spacing
		InputField(InputFieldParams{
			Form:      f,
			FormField: "Name",
			Name:      "name",
			InputType: "text",
			Label:     "Name",
			Value:     f.Name,
		}),
		InputField(InputFieldParams{
			Form:        f,
			FormField:   "PhoneNumber",
			Name:        "phone_number",
			InputType:   "tel",
			Label:       "Phone Number",
			Value:       f.PhoneNumber,
			Placeholder: "+1234567890",
		}),
		InputField(InputFieldParams{
			Form:        f,
			FormField:   "Password",
			Name:        "password",
			InputType:   "password",
			Label:       "Password",
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
			Class("pt-2"), // Removed flex gap since only one button now
			Button(
				Type("submit"),
				Class("w-full bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-md transition-colors"),
				Text("Sign up!"),
			),
		),
		CSRF(r),
		Div(
			Class("text-center text-slate-600 pt-4 border-t border-slate-200"),
			Text("Already have an account? "),
			A(
				Class("text-blue-600 hover:text-blue-800 font-medium transition-colors"),
				Href(r.Path(routenames.Login)),
				Text("Sign in!"),
			),
		),
	)
}
