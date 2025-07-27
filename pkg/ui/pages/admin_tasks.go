package pages

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/pkg/ui"
	"github.com/r-scheele/zero/pkg/ui/icons"
	"github.com/r-scheele/zero/pkg/ui/layouts"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// AdminTasks renders a custom Background Tasks page
func AdminTasks(ctx echo.Context) error {
	r := ui.NewRequest(ctx)
	r.Title = ""

	// Simplified - just show a placeholder for now
	backliteOutput := "No tasks to display"
	hasTasks := !strings.Contains(backliteOutput, "No tasks to display")

	// Create our own clean task interface
	content := Div(
		Class("space-y-6"),

		// Header Section
		Div(
			Class("bg-gradient-to-r from-blue-50 to-cyan-50 rounded-2xl p-6 border border-blue-200"),
			Div(
				Class("flex items-center gap-3 mb-3"),
				Div(
					Class("w-10 h-10 bg-blue-100 rounded-xl flex items-center justify-center"),
					Div(Class("w-5 h-5 text-blue-600"), icons.CircleStack()),
				),
			),
			P(
				Class("text-blue-700"),
				Text("Monitor and manage background processes and scheduled tasks."),
			),
		),

		// Task Status Cards
		Div(
			Class("grid grid-cols-1 md:grid-cols-4 gap-4 mb-6"),

			// Running Tasks
			Div(
				Class("bg-white rounded-lg p-4 shadow-sm border border-gray-200"),
				Div(
					Class("flex items-center justify-between"),
					Div(
						H3(Class("text-sm font-medium text-gray-600"), Text("Running")),
						P(Class("text-2xl font-bold text-blue-600"), Text("0")),
					),
					Div(Class("w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center")),
				),
			),

			// Upcoming Tasks
			Div(
				Class("bg-white rounded-lg p-4 shadow-sm border border-gray-200"),
				Div(
					Class("flex items-center justify-between"),
					Div(
						H3(Class("text-sm font-medium text-gray-600"), Text("Upcoming")),
						P(Class("text-2xl font-bold text-indigo-600"), Text("0")),
					),
					Div(Class("w-8 h-8 bg-indigo-100 rounded-full flex items-center justify-center")),
				),
			),

			// Succeeded Tasks
			Div(
				Class("bg-white rounded-lg p-4 shadow-sm border border-gray-200"),
				Div(
					Class("flex items-center justify-between"),
					Div(
						H3(Class("text-sm font-medium text-gray-600"), Text("Succeeded")),
						P(Class("text-2xl font-bold text-green-600"), Text("0")),
					),
					Div(Class("w-8 h-8 bg-green-100 rounded-full flex items-center justify-center")),
				),
			),

			// Failed Tasks
			Div(
				Class("bg-white rounded-lg p-4 shadow-sm border border-gray-200"),
				Div(
					Class("flex items-center justify-between"),
					Div(
						H3(Class("text-sm font-medium text-gray-600"), Text("Failed")),
						P(Class("text-2xl font-bold text-red-600"), Text("0")),
					),
					Div(Class("w-8 h-8 bg-red-100 rounded-full flex items-center justify-center")),
				),
			),
		),

		// Recent Tasks Section
		Div(
			Class("bg-white rounded-lg shadow-sm border border-gray-200"),
			Div(
				Class("px-6 py-4 border-b border-gray-200"),
				H2(Class("text-lg font-semibold text-gray-900"), Text("Recent Tasks")),
			),
			Div(
				Class("p-6"),
				If(hasTasks,
					Div(
						Class("text-center py-8"),
						P(Class("text-gray-500"), Text("Tasks available - data parsing in progress")),
					),
				),
				If(!hasTasks,
					Div(
						Class("text-center py-8"),
						P(Class("text-gray-500"), Text("No tasks to display")),
					),
				),
			),
		),

		// Action Buttons
		Div(
			Class("flex gap-4 pt-6"),
			Button(
				Class("inline-flex items-center px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition-colors"),
				Type("button"),
				Text("Refresh Tasks"),
			),
		),

	)

	return r.Render(layouts.Admin, content)
}
