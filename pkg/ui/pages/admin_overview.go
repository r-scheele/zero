package pages

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/ent"
	"github.com/r-scheele/zero/ent/user"
	"github.com/r-scheele/zero/pkg/ui"
	"github.com/r-scheele/zero/pkg/ui/layouts"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func AdminOverview(ctx echo.Context, orm *ent.Client) error {
	r := ui.NewRequest(ctx)
	r.Title = "Admin Dashboard"

	// Get some basic statistics
	totalUsers, _ := orm.User.Query().Count(context.Background())
	adminUsers, _ := orm.User.Query().Where(user.Admin(true)).Count(context.Background())
	verifiedUsers, _ := orm.User.Query().Where(user.Verified(true)).Count(context.Background())

	return r.Render(layouts.Admin, Group{
		// Welcome message
		Div(
			Class("mb-8 bg-gradient-to-r from-blue-50 to-cyan-50 rounded-2xl p-6 border border-blue-200"),
			H2(
				Class("text-xl font-semibold text-blue-800 mb-2"),
				Text("Welcome to Admin Dashboard"),
			),
			P(
				Class("text-blue-700"),
				Text("Manage your application, monitor system health, and perform administrative tasks."),
			),
		),

		// Statistics Cards
		Div(
			Class("grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8"),
			// Total Users Card
			Div(
				Class("bg-white rounded-2xl p-6 shadow-lg border border-slate-200"),
				Div(
					Class("flex items-center justify-between"),
					Div(
						H3(
							Class("text-lg font-semibold text-slate-800 mb-2"),
							Text("Total Users"),
						),
						P(
							Class("text-3xl font-bold text-blue-600"),
							Text(fmt.Sprintf("%d", totalUsers)),
						),
					),
					Div(
						Class("w-12 h-12 bg-blue-100 rounded-xl flex items-center justify-center"),
						Span(Class("text-blue-600 text-xl"), Text("👥")),
					),
				),
			),
			// Admin Users Card
			Div(
				Class("bg-white rounded-2xl p-6 shadow-lg border border-slate-200"),
				Div(
					Class("flex items-center justify-between"),
					Div(
						H3(
							Class("text-lg font-semibold text-slate-800 mb-2"),
							Text("Admin Users"),
						),
						P(
							Class("text-3xl font-bold text-indigo-600"),
							Text(fmt.Sprintf("%d", adminUsers)),
						),
					),
					Div(
						Class("w-12 h-12 bg-indigo-100 rounded-xl flex items-center justify-center"),
						Span(Class("text-indigo-600 text-xl"), Text("🛡️")),
					),
				),
			),
			// Verified Users Card
			Div(
				Class("bg-white rounded-2xl p-6 shadow-lg border border-slate-200"),
				Div(
					Class("flex items-center justify-between"),
					Div(
						H3(
							Class("text-lg font-semibold text-slate-800 mb-2"),
							Text("Verified Users"),
						),
						P(
							Class("text-3xl font-bold text-indigo-600"),
							Text(fmt.Sprintf("%d", verifiedUsers)),
						),
					),
					Div(
						Class("w-12 h-12 bg-indigo-100 rounded-xl flex items-center justify-center"),
						Span(Class("text-indigo-600 text-xl"), Text("✅")),
					),
				),
			),
		),

		// Quick Actions
		Div(
			Class("bg-white rounded-2xl p-6 shadow-lg border border-slate-200 mb-8"),
			H2(
				Class("text-xl font-semibold text-slate-800 mb-4"),
				Text("Quick Actions"),
			),
			Div(
				Class("grid grid-cols-1 md:grid-cols-2 gap-4"),
				A(
					Href("/admin/entity/user"),
					Class("flex items-center gap-4 p-4 rounded-xl border border-slate-200 hover:bg-slate-50 transition-colors group"),
					Div(
						Class("w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center group-hover:bg-blue-200 transition-colors"),
						Span(Class("text-blue-600"), Text("👥")),
					),
					Div(
						H3(Class("font-medium text-slate-800"), Text("Manage Users")),
						P(Class("text-sm text-slate-600"), Text("Add, edit, and manage user accounts")),
					),
				),
				A(
					Href("/admin/tasks"),
					Class("flex items-center gap-4 p-4 rounded-xl border border-slate-200 hover:bg-slate-50 transition-colors group"),
					Div(
						Class("w-10 h-10 bg-violet-100 rounded-lg flex items-center justify-center group-hover:bg-violet-200 transition-colors"),
						Span(Class("text-violet-600"), Text("⚙️")),
					),
					Div(
						H3(Class("font-medium text-slate-800"), Text("Background Tasks")),
						P(Class("text-sm text-slate-600"), Text("Monitor and manage system tasks")),
					),
				),
				A(
					Href("/cache"),
					Class("flex items-center gap-4 p-4 rounded-xl border border-slate-200 hover:bg-slate-50 transition-colors group"),
					Div(
						Class("w-10 h-10 bg-red-100 rounded-lg flex items-center justify-center group-hover:bg-red-200 transition-colors"),
						Span(Class("text-red-600"), Text("🗄️")),
					),
					Div(
						H3(Class("font-medium text-slate-800"), Text("Cache Management")),
						P(Class("text-sm text-slate-600"), Text("Clear and manage application cache")),
					),
				),
				A(
					Href("/admin/entity/passwordtoken"),
					Class("flex items-center gap-4 p-4 rounded-xl border border-slate-200 hover:bg-slate-50 transition-colors group"),
					Div(
						Class("w-10 h-10 bg-amber-100 rounded-lg flex items-center justify-center group-hover:bg-amber-200 transition-colors"),
						Span(Class("text-amber-600"), Text("🔑")),
					),
					Div(
						H3(Class("font-medium text-slate-800"), Text("Password Tokens")),
						P(Class("text-sm text-slate-600"), Text("Manage password reset tokens")),
					),
				),
			),
		),

		// System Information
		Div(
			Class("bg-gradient-to-r from-blue-50 to-cyan-50 rounded-2xl p-6 border border-blue-200"),
			H2(
				Class("text-xl font-semibold text-blue-800 mb-4"),
				Text("System Information"),
			),
			Div(
				Class("grid grid-cols-1 md:grid-cols-2 gap-4 text-sm"),
				Div(
					Class("flex justify-between"),
					Span(Class("text-blue-700 font-medium"), Text("Application:")),
					Span(Class("text-blue-800"), Text("Zero Admin Panel")),
				),
				Div(
					Class("flex justify-between"),
					Span(Class("text-blue-700 font-medium"), Text("Version:")),
					Span(Class("text-blue-800"), Text("1.0.0")),
				),
				Div(
					Class("flex justify-between"),
					Span(Class("text-blue-700 font-medium"), Text("Environment:")),
					Span(Class("text-blue-800"), Text("Development")),
				),
				Div(
					Class("flex justify-between"),
					Span(Class("text-blue-700 font-medium"), Text("Database:")),
					Span(Class("text-blue-800"), Text("SQLite")),
				),
			),
		),
	})
}
