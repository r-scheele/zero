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

type ForgotPassword struct {
	PhoneNumber string `form:"phone_number" validate:"required"`
	form.Submission
}

func (f *ForgotPassword) Render(r *ui.Request) Node {
	return Form(
		ID("forgot-password"),
		Method(http.MethodPost),
		HxBoost(),
		Action(r.Path(routenames.ForgotPasswordSubmit)),
		Class("space-y-6"),
		ForgotPasswordIllustration(), // Add forgot password illustration
		InputField(InputFieldParams{
			Form:        f,
			FormField:   "PhoneNumber",
			Name:        "phone_number",
			InputType:   "tel",
			Label:       "Phone Number",
			Placeholder: "+1234567890",
			Value:       f.PhoneNumber,
		}),
		Div(
			Class("pt-4"), // Removed flex gap since only one button now
			Button(
				Class("btn-modern-primary w-full"), // Full width button
				Type("submit"),
				Text("Reset Password via WhatsApp �"),
			),
		),
		CSRF(r),
	)
}
