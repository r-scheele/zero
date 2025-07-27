package emails

import (
	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/pkg/routenames"
	"github.com/r-scheele/zero/pkg/ui"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func ConfirmEmailAddress(ctx echo.Context, username, token string) Node {
	url := ui.NewRequest(ctx).
		Url(routenames.VerifyEmail, token)

	return HTML(
		Lang("en"),
		Head(
			Meta(Charset("UTF-8")),
			Meta(Name("viewport"), Content("width=device-width, initial-scale=1.0")),
			TitleEl(Text("Confirm Your Email Address")),
			StyleEl(Text(`
				body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #333; margin: 0; padding: 20px; background-color: #f8fafc; }
				.container { max-width: 600px; margin: 0 auto; background: white; border-radius: 12px; padding: 40px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1); }
				.header { text-align: center; margin-bottom: 32px; }
				.logo { font-size: 28px; font-weight: bold; background: linear-gradient(135deg, #3b82f6, #8b5cf6); -webkit-background-clip: text; -webkit-text-fill-color: transparent; margin-bottom: 8px; }
				.title { font-size: 24px; font-weight: 600; color: #1e293b; margin-bottom: 16px; }
				.content { margin-bottom: 32px; }
				.button { display: inline-block; background: linear-gradient(135deg, #3b82f6, #8b5cf6); color: white; text-decoration: none; padding: 16px 32px; border-radius: 8px; font-weight: 600; text-align: center; margin: 16px 0; }
				.footer { text-align: center; color: #64748b; font-size: 14px; border-top: 1px solid #e2e8f0; padding-top: 24px; margin-top: 32px; }
				.url-fallback { background: #f1f5f9; padding: 12px; border-radius: 6px; margin: 16px 0; word-break: break-all; font-family: monospace; font-size: 14px; }
			`)),
		),
		Body(
			Div(Class("container"),
				Div(Class("header"),
					Div(Class("logo"), Text("Zero")),
					H1(Class("title"), Text("Confirm Your Email Address")),
				),
				Div(Class("content"),
					P(Textf("Hello %s,", username)),
					P(Text("Thank you for creating an account with us! To complete your registration and verify your email address, please click the button below:")),
					Div(Style("text-align: center; margin: 24px 0;"),
						A(
							Href(url),
							Class("button"),
							Text("Verify Email Address"),
						),
					),
					P(Text("If the button above doesn't work, you can copy and paste the following link into your browser:")),
					Div(Class("url-fallback"), Text(url)),
					P(Text("This verification link will expire in 12 hours for security reasons.")),
					P(Text("If you didn't create an account with us, you can safely ignore this email.")),
				),
				Div(Class("footer"),
					P(Text("Best regards,")),
					P(Text("The Zero Team")),
				),
			),
		),
	)
}
