package pages

import (
	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/pkg/ui"
	"github.com/r-scheele/zero/pkg/ui/layouts"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// QuizList displays the list of quizzes
func QuizList(ctx echo.Context) error {
	r := ui.NewRequest(ctx)
	r.Title = "Quizzes"
	r.Metatags.Description = "Browse and manage your quizzes"

	page := Div(
		Class("container mx-auto px-4 py-8"),
		Div(
			Class("flex justify-between items-center mb-6"),
			H1(Class("text-3xl font-bold text-gray-900"), Text("Quizzes")),
			A(
				Href("/quiz/create"),
				Class("bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg font-medium"),
				Text("Create Quiz"),
			),
		),
		Div(
			Class("bg-white rounded-lg shadow p-6"),
			P(Class("text-gray-600 text-center py-8"), Text("Quiz functionality is coming soon! You'll be able to create, manage, and take interactive quizzes.")),
		),
	)

	return r.Render(layouts.Primary, page)
}

// CreateQuiz displays the quiz creation form
func CreateQuiz(ctx echo.Context) error {
	r := ui.NewRequest(ctx)
	r.Title = "Create Quiz"
	r.Metatags.Description = "Create a new interactive quiz"

	page := Div(
		Class("container mx-auto px-4 py-8"),
		Div(
			Class("max-w-2xl mx-auto"),
			H1(Class("text-3xl font-bold text-gray-900 mb-6"), Text("Create Quiz")),
			Div(
				Class("bg-white rounded-lg shadow p-6"),
				Div(
					Class("text-center py-12"),
					Div(
						Class("text-6xl mb-4"),
						Text("🧠"),
					),
					H2(Class("text-2xl font-semibold text-gray-900 mb-4"), Text("Quiz Creation Coming Soon!")),
					P(Class("text-gray-600 mb-6"), Text("We're working on an amazing quiz creation tool that will allow you to build interactive quizzes with multiple question types, automatic grading, and detailed analytics.")),

				),
			),
		),
	)

	return r.Render(layouts.Primary, page)
}

// ViewQuiz displays a specific quiz
func ViewQuiz(ctx echo.Context) error {
	r := ui.NewRequest(ctx)
	r.Title = "View Quiz"
	r.Metatags.Description = "View quiz details"

	page := Div(
		Class("container mx-auto px-4 py-8"),
		Div(
			Class("bg-white rounded-lg shadow p-6"),
			P(Class("text-gray-600 text-center py-8"), Text("Quiz viewing functionality is coming soon!")),
		),
	)

	return r.Render(layouts.Primary, page)
}

// EditQuiz displays the quiz edit form
func EditQuiz(ctx echo.Context) error {
	r := ui.NewRequest(ctx)
	r.Title = "Edit Quiz"
	r.Metatags.Description = "Edit quiz details"

	page := Div(
		Class("container mx-auto px-4 py-8"),
		Div(
			Class("bg-white rounded-lg shadow p-6"),
			P(Class("text-gray-600 text-center py-8"), Text("Quiz editing functionality is coming soon!")),
		),
	)

	return r.Render(layouts.Primary, page)
}

// TakeQuiz displays the quiz taking interface
func TakeQuiz(ctx echo.Context) error {
	r := ui.NewRequest(ctx)
	r.Title = "Take Quiz"
	r.Metatags.Description = "Take an interactive quiz"

	page := Div(
		Class("container mx-auto px-4 py-8"),
		Div(
			Class("bg-white rounded-lg shadow p-6"),
			P(Class("text-gray-600 text-center py-8"), Text("Quiz taking functionality is coming soon!")),
		),
	)

	return r.Render(layouts.Primary, page)
}