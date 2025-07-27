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

type Profile struct {
	Name               string `form:"name" validate:"required"`
	PhoneNumber        string `form:"phone_number" validate:"required"`
	Email              string `form:"email" validate:"omitempty,email"`
	Bio                string `form:"bio"`
	DarkMode           bool   `form:"dark_mode"`
	EmailNotifications bool   `form:"email_notifications"`
	SmsNotifications   bool   `form:"sms_notifications"`
	form.Submission
}

func (f *Profile) Render(r *ui.Request) Node {
	return Form(
		ID("profile"),
		Method(http.MethodPost),
		HxBoost(),
		Action(r.Path(routenames.ProfileUpdate)),
		Class("space-y-6"),
		FlashMessages(r),

		// Basic Information
		Div(
			Class("space-y-4"),

			InputField(InputFieldParams{
				Form:        f,
				FormField:   "Name",
				Name:        "name",
				InputType:   "text",
				Label:       "Full Name",
				Value:       f.Name,
				Placeholder: "Enter your full name",
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
				FormField:   "Email",
				Name:        "email",
				InputType:   "email",
				Label:       "Email Address (Optional)",
				Value:       f.Email,
				Placeholder: "your@email.com",
			}),
			TextareaField(TextareaFieldParams{
				Form:      f,
				FormField: "Bio",
				Name:      "bio",
				Label:     "Bio (Optional)",
				Value:     f.Bio,
			}),
		),

		// Preferences
		Div(
			Class("border-t border-gray-200 pt-6"),
			H3(
				Class("text-sm font-medium text-gray-900 mb-4"),
				Text("Preferences"),
			),
			Div(
				Class("space-y-3"),
				Checkbox(CheckboxParams{
					Form:      f,
					FormField: "DarkMode",
					Name:      "dark_mode",
					Label:     "Dark Mode",
					Checked:   f.DarkMode,
				}),
				Checkbox(CheckboxParams{
					Form:      f,
					FormField: "EmailNotifications",
					Name:      "email_notifications",
					Label:     "Email Notifications",
					Checked:   f.EmailNotifications,
				}),
				Checkbox(CheckboxParams{
					Form:      f,
					FormField: "SmsNotifications",
					Name:      "sms_notifications",
					Label:     "SMS Notifications",
					Checked:   f.SmsNotifications,
				}),
			),
		),

		// Submit Button
		Div(
			Class("flex justify-end pt-6 border-t border-gray-200"),
			Button(
				Type("submit"),
				Class("bg-blue-600 hover:bg-blue-700 text-white font-medium px-6 py-2 rounded-lg transition-colors"),
				Text("Update Profile"),
			),
		),
	)
}

type ChangePassword struct {
	CurrentPassword string `form:"current_password" validate:"required"`
	NewPassword     string `form:"new_password" validate:"required,min=8"`
	ConfirmPassword string `form:"confirm_password" validate:"required,eqfield=NewPassword"`
	form.Submission
}

func (f *ChangePassword) Render(r *ui.Request) Node {
	return Form(
		ID("change-password"),
		Method(http.MethodPost),
		HxBoost(),
		Action(r.Path(routenames.ProfileUpdate)+"/password"),
		Class("space-y-6"),
		FlashMessages(r),

		Div(
			Class("bg-white rounded-2xl p-6 border border-slate-200"),

			Div(
				Class("space-y-4"),
				InputField(InputFieldParams{
					Form:        f,
					FormField:   "CurrentPassword",
					Name:        "current_password",
					InputType:   "password",
					Label:       "Current Password",
					Placeholder: "••••••••",
				}),
				InputField(InputFieldParams{
					Form:        f,
					FormField:   "NewPassword",
					Name:        "new_password",
					InputType:   "password",
					Label:       "New Password",
					Placeholder: "••••••••",
				}),
				InputField(InputFieldParams{
					Form:        f,
					FormField:   "ConfirmPassword",
					Name:        "confirm_password",
					InputType:   "password",
					Label:       "Confirm New Password",
					Placeholder: "••••••••",
				}),
			),
		),

		Div(
			Class("flex justify-end space-x-3"),
			A(
				Href(r.Path(routenames.Profile)),
				Class("bg-slate-600 hover:bg-slate-700 text-white font-medium px-6 py-3 rounded-xl transition-colors duration-300"),
				Text("Cancel"),
			),
			Button(
				Type("submit"),
				Class("bg-emerald-600 hover:bg-emerald-700 text-white font-semibold px-6 py-3 rounded-xl transition-colors duration-300"),
				Text("Update Password"),
			),
		),
	)
}

type ProfilePicture struct {
	Picture string `form:"picture" validate:"required"`
	form.Submission
}

func (f *ProfilePicture) Render(r *ui.Request) Node {
	return Form(
		ID("profile-picture"),
		Method(http.MethodPost),
		HxBoost(),
		Action(r.Path(routenames.ProfilePicture)),
		EncType("multipart/form-data"),
		Class("space-y-6"),
		FlashMessages(r),

		Div(
			Class("space-y-4"),
			FileField(FileFieldParams{
				Name:  "picture",
				Label: "Choose Profile Picture",
			}),
		),

		Div(
			Class("flex justify-end space-x-3 pt-4"),
			A(
				Href(r.Path(routenames.Profile)),
				Class("bg-gray-600 hover:bg-gray-700 text-white font-medium px-6 py-2 rounded-lg transition-colors"),
				Text("Cancel"),
			),
			Button(
				Type("submit"),
				Class("bg-blue-600 hover:bg-blue-700 text-white font-medium px-6 py-2 rounded-lg transition-colors"),
				Text("Upload"),
			),
		),
	)
}

type DeactivateAccount struct {
	Password string `form:"password" validate:"required"`
	Reason   string `form:"reason" validate:"required"`
	form.Submission
}

func (f *DeactivateAccount) Render(r *ui.Request) Node {
	return Form(
		ID("deactivate-account"),
		Method(http.MethodPost),
		HxBoost(),
		Action(r.Path(routenames.ProfileDeactivate)),
		Class("space-y-6"),
		FlashMessages(r),

		Div(
			Class("bg-red-50 border border-red-200 rounded-2xl p-6"),
			Div(
				Class("flex items-center gap-3 mb-4"),
				Span(Class("text-red-600 text-2xl"), Text("⚠️")),
				H3(
					Class("text-lg font-semibold text-red-900"),
					Text("Deactivate Account"),
				),
			),
			P(
				Class("text-red-700 mb-4"),
				Text("This action will deactivate your account. You will be logged out and won't be able to access your account until it's reactivated by an administrator."),
			),
			Div(
				Class("space-y-4"),
				InputField(InputFieldParams{
					Form:        f,
					FormField:   "Password",
					Name:        "password",
					InputType:   "password",
					Label:       "Confirm your password",
					Placeholder: "••••••••",
				}),
				TextareaField(TextareaFieldParams{
					Form:      f,
					FormField: "Reason",
					Name:      "reason",
					Label:     "Reason for deactivation",
					Value:     f.Reason,
				}),
			),
		),

		Div(
			Class("flex justify-end space-x-3"),
			A(
				Href(r.Path(routenames.Profile)),
				Class("bg-slate-600 hover:bg-slate-700 text-white font-medium px-6 py-3 rounded-xl transition-colors duration-300"),
				Text("Cancel"),
			),
			Button(
				Type("submit"),
				Class("bg-red-600 hover:bg-red-700 text-white font-semibold px-6 py-3 rounded-xl transition-colors duration-300"),
				Text("Deactivate Account"),
			),
		),
	)
}
