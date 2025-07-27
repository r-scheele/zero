package pages

import (
	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/ent"
	"github.com/r-scheele/zero/pkg/ui"
	"github.com/r-scheele/zero/pkg/ui/components"
	"github.com/r-scheele/zero/pkg/ui/forms"
	"github.com/r-scheele/zero/pkg/ui/layouts"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Profile(ctx echo.Context, profileForm *forms.Profile, user *ent.User) error {
	r := ui.NewRequest(ctx)

	// Safety check
	if user == nil {
		return echo.NewHTTPError(500, "User not found")
	}

	return r.Render(layouts.Primary, Group{
		Div(
			Class("max-w-2xl mx-auto px-4 py-12"),

			// Profile navigation
			components.ProfileNav(r),

			// Profile form
			Div(
				ID("profile"),
				Class("max-w-2xl mx-auto"),
				profileForm.Render(r),
			),
		),
	})
}

func ProfileEdit(ctx echo.Context, profileForm *forms.Profile) error {
	r := ui.NewRequest(ctx)
	// No title for cleaner design

	return r.Render(layouts.Primary, Group{
		Div(
			Class("max-w-2xl mx-auto px-4 py-12"),
			// Profile navigation
			components.ProfileNav(r),

			// Profile form
			Div(
				Class("max-w-2xl mx-auto"),
				profileForm.Render(r),
			),
		),
	})
}

func ChangePassword(ctx echo.Context, form *forms.ChangePassword) error {
	r := ui.NewRequest(ctx)
	// No title for cleaner design

	return r.Render(layouts.Primary, Group{
		Div(
			Class("max-w-2xl mx-auto px-4 py-12"),
			// Profile navigation
			components.ProfileNav(r),

			// Change password form
			Div(
				Class("max-w-md mx-auto"),
				form.Render(r),
			),
		),
	})
}

func ProfilePicture(ctx echo.Context, form *forms.ProfilePicture, user *ent.User) error {
	r := ui.NewRequest(ctx)
	// No title for cleaner design

	// Safety check
	if user == nil {
		return echo.NewHTTPError(500, "User not found")
	}

	return r.Render(layouts.Primary, Group{
		Div(
			Class("max-w-2xl mx-auto px-4 py-12"),

			// Profile navigation
			components.ProfileNav(r),

			// Current profile picture preview
			Div(
				Class("text-center mb-8"),
				func() Node {
					if user.ProfilePicture != nil && *user.ProfilePicture != "" {
						return Img(
							Src("/files/"+*user.ProfilePicture),
							Alt("Current Profile Picture"),
							Class("w-32 h-32 rounded-full mx-auto border-4 border-gray-200 object-cover"),
						)
					} else {
						return Div(
							Class("w-32 h-32 rounded-full mx-auto border-4 border-gray-200 bg-blue-500 flex items-center justify-center"),
							Span(
								Class("text-white text-4xl font-bold"),
								Text(func() string {
									if len(user.Name) > 0 {
										return string(user.Name[0])
									}
									return "U"
								}()),
							),
						)
					}
				}(),
			),

			// Upload form
			Div(
				Class("bg-white rounded-lg shadow-sm border border-gray-200 p-8"),
				form.Render(r),
			),
		),
	})
}

func DeactivateAccount(ctx echo.Context, form *forms.DeactivateAccount) error {
	r := ui.NewRequest(ctx)
	// No title for cleaner design

	return r.Render(layouts.Primary, Group{
		Main(
			Class("max-w-4xl mx-auto px-4 py-12"),

			// Profile navigation
			components.ProfileNav(r),

			// Deactivate form
			Div(
				Class("bg-white rounded-lg shadow-sm border border-gray-200 p-8 max-w-2xl mx-auto"),
				Div(
					Class("text-center mb-8"),
					H2(
						Class("text-2xl font-bold text-gray-900 mb-2"),
						Text("Deactivate Account"),
					),
					P(
						Class("text-gray-600"),
						Text("Once you deactivate your account, your profile and data will be permanently removed. This action cannot be undone."),
					),
				),

				Div(
					Class("bg-red-50 border border-red-200 rounded-lg p-6 mb-6"),
					Div(
						Class("flex"),
						Div(
							Class("flex-shrink-0"),
							Div(
								Class("h-5 w-5 text-red-400"),
								Text("⚠️"),
							),
						),
						Div(
							Class("ml-3"),
							H3(
								Class("text-sm font-medium text-red-800"),
								Text("Warning"),
							),
							Div(
								Class("mt-2 text-sm text-red-700"),
								P(Text("This will permanently delete your account and all associated data.")),
							),
						),
					),
				),

				// Form content
				form.Render(r),
			),
		),
	})
}
