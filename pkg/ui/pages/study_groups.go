package pages

import (
	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/pkg/ui"
	"github.com/r-scheele/zero/pkg/ui/layouts"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// StudyGroupsList displays the list of study groups
func StudyGroupsList(ctx echo.Context) error {
	r := ui.NewRequest(ctx)
	r.Title = "Study Groups"
	r.Metatags.Description = "Browse and join collaborative study groups"

	page := Div(
		Class("container mx-auto px-4 py-8"),
		Div(
			Class("flex justify-between items-center mb-6"),
			H1(Class("text-3xl font-bold text-gray-900"), Text("Study Groups")),
			A(
				Href("/study-groups/create"),
				Class("bg-pink-600 hover:bg-pink-700 text-white px-4 py-2 rounded-lg font-medium"),
				Text("Create Study Group"),
			),
		),
		Div(
			Class("bg-white rounded-lg shadow p-6"),
			P(Class("text-gray-600 text-center py-8"), Text("Study groups functionality is coming soon! You'll be able to create and join collaborative study sessions with other learners.")),
		),
	)

	return r.Render(layouts.Primary, page)
}

// CreateStudyGroup displays the study group creation form
func CreateStudyGroup(ctx echo.Context) error {
	r := ui.NewRequest(ctx)
	r.Title = "Create Study Group"
	r.Metatags.Description = "Create a new collaborative study group"

	page := Div(
		Class("container mx-auto px-4 py-8"),
		Div(
			Class("max-w-2xl mx-auto"),
			H1(Class("text-3xl font-bold text-gray-900 mb-6"), Text("Create Study Group")),
			Div(
				Class("bg-white rounded-lg shadow p-6"),
				Div(
					Class("text-center py-12"),
					Div(
						Class("text-6xl mb-4"),
						Text("👥"),
					),
					H2(Class("text-2xl font-semibold text-gray-900 mb-4"), Text("Study Groups Coming Soon!")),
					P(Class("text-gray-600 mb-6"), Text("We're developing collaborative study groups where you can connect with other learners, share resources, and study together in real-time.")),

				),
			),
		),
	)

	return r.Render(layouts.Primary, page)
}

// ViewStudyGroup displays a specific study group
func ViewStudyGroup(ctx echo.Context) error {
	r := ui.NewRequest(ctx)
	r.Title = "View Study Group"
	r.Metatags.Description = "View study group details"

	page := Div(
		Class("container mx-auto px-4 py-8"),
		Div(
			Class("bg-white rounded-lg shadow p-6"),
			P(Class("text-gray-600 text-center py-8"), Text("Study group viewing functionality is coming soon!")),
		),
	)

	return r.Render(layouts.Primary, page)
}