package pages

import (
	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/pkg/ui"
	"github.com/r-scheele/zero/pkg/ui/layouts"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// ViewProgress displays the progress dashboard
func ViewProgress(ctx echo.Context) error {
	r := ui.NewRequest(ctx)
	r.Title = "View Progress"
	r.Metatags.Description = "Track your learning progress and achievements"

	page := Div(
		Class("container mx-auto px-4 py-8"),
		Div(
			Class("flex justify-between items-center mb-6"),
			H1(Class("text-3xl font-bold text-gray-900"), Text("Learning Progress")),
			A(
				Href("/progress/analytics"),
				Class("bg-purple-600 hover:bg-purple-700 text-white px-4 py-2 rounded-lg font-medium"),
				Text("View Analytics"),
			),
		),
		Div(
			Class("bg-white rounded-lg shadow p-6"),
			Div(
				Class("text-center py-12"),
				Div(
					Class("text-6xl mb-4"),
					Text("📊"),
				),
				H2(Class("text-2xl font-semibold text-gray-900 mb-4"), Text("Progress Tracking Coming Soon!")),
				P(Class("text-gray-600 mb-6"), Text("We're building comprehensive progress tracking that will show your learning journey, achievements, quiz scores, and study patterns with detailed analytics and insights.")),

			),
		),
	)

	return r.Render(layouts.Primary, page)
}

// ViewAnalytics displays detailed analytics
func ViewAnalytics(ctx echo.Context) error {
	r := ui.NewRequest(ctx)
	r.Title = "Learning Analytics"
	r.Metatags.Description = "Detailed analytics of your learning progress"

	page := Div(
		Class("container mx-auto px-4 py-8"),
		Div(
			Class("flex justify-between items-center mb-6"),
			H1(Class("text-3xl font-bold text-gray-900"), Text("Learning Analytics")),
			A(
				Href("/progress"),
				Class("bg-gray-600 hover:bg-gray-700 text-white px-4 py-2 rounded-lg font-medium"),
				Text("Back to Progress"),
			),
		),
		Div(
			Class("bg-white rounded-lg shadow p-6"),
			P(Class("text-gray-600 text-center py-8"), Text("Detailed analytics functionality is coming soon! You'll be able to view comprehensive reports on your learning patterns, strengths, and areas for improvement.")),
		),
	)

	return r.Render(layouts.Primary, page)
}