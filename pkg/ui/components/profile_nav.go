package components

import (
	"github.com/r-scheele/zero/pkg/routenames"
	"github.com/r-scheele/zero/pkg/ui"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// ProfileNav creates a consistent navigation bar for all profile-related pages
func ProfileNav(r *ui.Request) Node {
	return Div(
		Class("mb-8"),
		Div(
			Class("border-b border-gray-200"),
			Nav(
				Class("flex space-x-8"),
				// Basic tab
				A(
					Href(r.Path(routenames.Profile)),
					Class(func() string {
						// Check if we're on any profile basic page
						if r.CurrentPath == r.Path(routenames.Profile) || r.CurrentPath == r.Path(routenames.ProfileEdit) || r.CurrentPath == r.Path(routenames.ProfileUpdate) || r.CurrentPath == "/profile" || r.CurrentPath == "/profile/edit" || r.CurrentPath == "/profile/update" {
							return "py-4 px-1 border-b-2 border-blue-500 text-blue-600 font-medium"
						}
						return "py-4 px-1 border-b-2 border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 font-medium"
					}()),
					Text("Basic"),
				),
				// Picture tab
				A(
					Href(r.Path(routenames.ProfilePicture)),
					Class(func() string {
						// Check if we're on any picture-related page
						if r.CurrentPath == r.Path(routenames.ProfilePicture) || r.CurrentPath == "/profile/picture" {
							return "py-4 px-1 border-b-2 border-blue-500 text-blue-600 font-medium"
						}
						return "py-4 px-1 border-b-2 border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 font-medium"
					}()),
					Text("Picture"),
				),
				// Security tab
				A(
					Href(r.Path(routenames.ProfileChangePassword)),
					Class(func() string {
						// Check if we're on any security-related page
						if r.CurrentPath == r.Path(routenames.ProfileChangePassword) || r.CurrentPath == "/profile/change-password" {
							return "py-4 px-1 border-b-2 border-blue-500 text-blue-600 font-medium"
						}
						return "py-4 px-1 border-b-2 border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 font-medium"
					}()),
					Text("Security"),
				),
				// Account tab
				A(
					Href(r.Path(routenames.ProfileDeactivate)),
					Class(func() string {
						// Check if we're on any account-related page
						if r.CurrentPath == r.Path(routenames.ProfileDeactivate) || r.CurrentPath == "/profile/deactivate" {
							return "py-4 px-1 border-b-2 border-blue-500 text-blue-600 font-medium"
						}
						return "py-4 px-1 border-b-2 border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 font-medium"
					}()),
					Text("Account"),
				),
			),
		),
	)
}