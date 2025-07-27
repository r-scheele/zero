package pages

import (
	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/pkg/routenames"
	"github.com/r-scheele/zero/pkg/ui"
	. "github.com/r-scheele/zero/pkg/ui/components"
	"github.com/r-scheele/zero/pkg/ui/forms"
	"github.com/r-scheele/zero/pkg/ui/layouts"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Login(ctx echo.Context, form *forms.Login) error {
	r := ui.NewRequest(ctx)
	r.Title = "Login"

	return r.Render(layouts.Auth, form.Render(r))
}

func Register(ctx echo.Context, form *forms.Register) error {
	r := ui.NewRequest(ctx)
	r.Title = "Register"

	return r.Render(layouts.Auth, form.Render(r))
}

func ForgotPassword(ctx echo.Context, form *forms.ForgotPassword) error {
	r := ui.NewRequest(ctx)
	r.Title = "Forgot password"

	content := Div(
		Class("space-y-6"),
		Div(
			Class("text-center mb-8"),
			H2(
				Class("text-2xl font-bold text-slate-900 mb-4"),
				Text("Forgot Password? 🤔"),
			),
			P(
				Class("text-slate-600 leading-relaxed"),
				Text("No worries! It happens to the best of us. Just enter your phone number and we'll send you a reset link on WhatsApp! 📱"),
			),
		),
		form.Render(r),
	)

	return r.Render(layouts.Auth, content)
}

func ResetPassword(ctx echo.Context, form *forms.ResetPassword) error {
	r := ui.NewRequest(ctx)
	r.Title = "Reset your password"

	return r.Render(layouts.Auth, form.Render(r))
}

func VerificationNotice(ctx echo.Context) error {
	r := ui.NewRequest(ctx)
	r.Title = "WhatsApp Verification Required"

	content := Div(
		Class("max-w-2xl mx-auto text-center space-y-8"),
		ForgotPasswordIllustration(), // Reuse the WhatsApp-friendly illustration
		H1(
			Class("text-3xl lg:text-4xl font-bold text-slate-900 mb-4"),
			Text("WhatsApp Verification Required"),
		),
		Div(
			Class("bg-white rounded-2xl p-8 shadow-elegant border border-slate-200/60 space-y-6"),
			P(
				Class("text-lg text-slate-700 leading-relaxed"),
				Text("We need to verify your WhatsApp number to keep things secure! �"),
			),
			P(
				Class("text-slate-600 leading-relaxed"),
				Text("Check your WhatsApp for our verification message with 3 numbered buttons. Select the code that matches the one you saw during registration! ✨"),
			),
			P(
				Class("text-slate-500 text-sm leading-relaxed"),
				Text("No WhatsApp message? We can send another one! �"),
			),
			Form(
				Method("POST"),
				Action(r.Path(routenames.ResendVerification)),
				CSRF(r),
				Button(
					Class("btn btn-primary btn-lg gap-2 bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 border-0 text-white font-medium shadow-lg hover:shadow-xl hover:scale-105 transition-all duration-200"),
					Type("submit"),
					Span(Class("text-lg"), Text("�")),
					Text("Send It Again!"),
				),
			),
		),
	)

	return r.Render(layouts.Auth, content)
}
