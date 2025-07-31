package pages

import (
	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/pkg/ui"
	"github.com/r-scheele/zero/pkg/ui/layouts"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// Quizzes displays the main quizzes page
func Quizzes(ctx echo.Context) error {
	r := ui.NewRequest(ctx)
	r.Title = "Quizzes"
	r.Metatags.Description = "Interactive quizzes and assessments"
	r.Metatags.Keywords = []string{"Quiz", "Assessment", "Learning", "Test", "Education"}

	page := Div(
		Class("max-w-4xl mx-auto px-4 py-8"),
		// Header Section
		Div(
			Class("text-center mb-8"),
			H1(
				Class("text-3xl font-bold text-gray-900 mb-3"),
				Text("Quizzes"),
			),
			P(
				Class("text-gray-600"),
				Text("Create and take interactive quizzes"),
			),
		),

		// Quick Actions
		Div(
			Class("grid grid-cols-1 md:grid-cols-3 gap-4 mb-8"),
			// Create Quiz Card
			A(
				Href("/quiz/create"),
				Class("block bg-white border border-gray-200 p-6 rounded-lg hover:shadow-md transition-shadow"),
				Div(
					Class("text-center"),
					Div(
						Class("w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center mx-auto mb-3"),
						Span(Class("text-blue-600 text-xl"), Text("✏️")),
					),
					H3(Class("font-semibold text-gray-900 mb-1"), Text("Create Quiz")),
					P(Class("text-sm text-gray-600"), Text("Build interactive quizzes")),
				),
			),
			// Browse Quizzes Card
			A(
				Href("/quiz"),
				Class("block bg-white border border-gray-200 p-6 rounded-lg hover:shadow-md transition-shadow"),
				Div(
					Class("text-center"),
					Div(
						Class("w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center mx-auto mb-3"),
						Span(Class("text-green-600 text-xl"), Text("📚")),
					),
					H3(Class("font-semibold text-gray-900 mb-1"), Text("Browse Quizzes")),
					P(Class("text-sm text-gray-600"), Text("Explore available quizzes")),
				),
			),
			// My Results Card
			A(
				Href("/progress"),
				Class("block bg-white border border-gray-200 p-6 rounded-lg hover:shadow-md transition-shadow"),
				Div(
					Class("text-center"),
					Div(
						Class("w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center mx-auto mb-3"),
						Span(Class("text-purple-600 text-xl"), Text("📊")),
					),
					H3(Class("font-semibold text-gray-900 mb-1"), Text("My Results")),
					P(Class("text-sm text-gray-600"), Text("Track your progress")),
				),
			),
		),

		// Coming Soon Section
		Div(
			Class("bg-white rounded-xl shadow-lg p-8 text-center"),
			Div(
				Class("text-6xl mb-6"),
				Text("🧠"),
			),
			H2(
				Class("text-3xl font-bold text-gray-900 mb-4"),
				Text("Advanced Quiz Features Coming Soon!"),
			),
			P(
				Class("text-lg text-gray-600 mb-8 max-w-3xl mx-auto"),
				Text("We're developing an comprehensive quiz system with multiple question types, timed assessments, collaborative quizzes, and AI-powered question generation. Stay tuned for these exciting features!"),
			),
			Div(
				Class("grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 text-sm text-gray-500"),
				Div(
					Class("flex items-center justify-center p-3 bg-gray-50 rounded-lg"),
					Span(Class("mr-2"), Text("✅")),
					Text("Multiple Choice"),
				),
				Div(
					Class("flex items-center justify-center p-3 bg-gray-50 rounded-lg"),
					Span(Class("mr-2"), Text("✅")),
					Text("True/False"),
				),
				Div(
					Class("flex items-center justify-center p-3 bg-gray-50 rounded-lg"),
					Span(Class("mr-2"), Text("✅")),
					Text("Fill in the Blank"),
				),
				Div(
					Class("flex items-center justify-center p-3 bg-gray-50 rounded-lg"),
					Span(Class("mr-2"), Text("✅")),
					Text("Essay Questions"),
				),
			),

		),
	)

	return r.Render(layouts.Primary, page)
}
