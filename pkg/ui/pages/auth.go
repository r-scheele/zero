package pages

import (
	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/ent"
	"github.com/r-scheele/zero/pkg/context"
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
	r.Title = ""

	return r.Render(layouts.Auth, form.Render(r))
}

func Register(ctx echo.Context, form *forms.Register) error {
	r := ui.NewRequest(ctx)
	r.Title = ""

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
	r.Title = ""

	// Get the authenticated user to retrieve their verification code
	verificationCode := "--"
	if userValue := ctx.Get(context.AuthenticatedUserKey); userValue != nil {
		if user, ok := userValue.(*ent.User); ok && user.VerificationCode != nil {
			verificationCode = *user.VerificationCode
		}
	}

	content := Div(
		Class("max-w-2xl mx-auto text-center space-y-8"),
		ForgotPasswordIllustration(), // Reuse the WhatsApp-friendly illustration
		Div(
			Class("bg-white rounded-2xl p-8 shadow-elegant border border-slate-200/60 space-y-6"),
			P(
				Class("text-slate-600 leading-relaxed mb-4"),
				Text("Check your WhatsApp for our verification message with 3 numbered buttons. Select the code:"),
			),
			Div(
				Class("text-center mb-4"),
				Span(
					Class("font-bold text-2xl text-blue-600 bg-blue-50 px-3 py-1 rounded-lg border border-blue-200"),
					Text(verificationCode),
				),
			),
			P(
				Class("text-slate-500 text-sm leading-relaxed"), 
				Text("No WhatsApp message? We can send another one!"),
			),
			Div(
				Class("flex flex-col sm:flex-row gap-3 justify-center items-stretch sm:items-center w-full max-w-md mx-auto"),
				Form(
					Class("flex-1"),
					Method("POST"),
					Action(r.Path(routenames.ResendVerification)),
					CSRF(r),
						Button(
						Class("btn btn-primary btn-md gap-2 w-full bg-blue-600 hover:bg-blue-700 border-0 text-white font-medium shadow-sm hover:shadow-md transition-all duration-200 text-sm px-4 py-3"),
						Type("submit"),
						Text("Send It Again!"),
					),
				),
				A(
					Class("flex-1"),
					Href(r.Path(routenames.Home)),
						Button(
						Class("btn btn-outline btn-md gap-2 w-full border-slate-300 text-slate-700 hover:bg-slate-50 hover:border-slate-400 font-medium shadow-sm hover:shadow-md transition-all duration-200 text-sm px-4 py-3"),
						Type("button"),
						Text("Login Instead"),
					),
				),
			),
		),
	)

	return r.Render(layouts.Auth, content)
}
