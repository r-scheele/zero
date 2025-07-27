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

type Login struct {
	PhoneNumber string `form:"phone_number" validate:"required,e164"`
	Password    string `form:"password" validate:"required"`
	form.Submission
}

func (f *Login) Render(r *ui.Request) Node {
	return Form(
		ID("login"),
		Method(http.MethodPost),
		HxBoost(),
		Action(r.Path(routenames.LoginSubmit)),
		Class("space-y-4"), // Reduced spacing
		FlashMessages(r),
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
		Div(
			Class("text-right"),
			A(
				Class("text-blue-600 hover:text-blue-800 text-sm font-medium transition-colors"),
				Href(r.Path(routenames.ForgotPassword)),
				Text("Forgot password? 🤔"),
			),
		),
		Div(
			Class("pt-2"), // Removed flex gap since only one button now
			Button(
				Class("btn-modern-primary w-full"), // Full width button
				Type("submit"),
				Text("Sign In"),
			),
		),
		CSRF(r),
		Div(
			Class("text-center text-slate-600 pt-4 border-t border-slate-200"),
			Text("Are you new here? "),
			A(
				Class("text-blue-600 hover:text-blue-800 font-medium transition-colors"),
				Href(r.Path(routenames.Register)),
				Text("Join us! 💫"),
			),
		),
	)
}
