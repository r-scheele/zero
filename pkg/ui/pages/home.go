package pages

import (
	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/pkg/redirect"
	"github.com/r-scheele/zero/pkg/ui"
	"github.com/r-scheele/zero/pkg/ui/layouts"
	"github.com/r-scheele/zero/pkg/ui/models"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Home(ctx echo.Context, posts *models.Posts) error {
	r := ui.NewRequest(ctx)
	r.Metatags.Description = "Zero - Your comprehensive quiz and document management platform"
	r.Metatags.Keywords = []string{"Quiz", "Documents", "Learning", "Management", "Education"}

	// If user is authenticated, show their home page with content
	if r.IsAuth {
		return authenticatedHomePage(r, posts)
	}

	// Landing page for non-authenticated users
	return landingPage(r)
}

// Dashboard function for the new dashboard route
func Dashboard(ctx echo.Context, posts *models.Posts) error {
	r := ui.NewRequest(ctx)
	r.Title = "Dashboard"
	r.Metatags.Description = "Your learning dashboard"

	// If user is admin, redirect to admin dashboard
	if r.IsAdmin {
		return redirect.New(r.Context).URL("/admin").Go()
	}

	return authenticatedDashboard(r, posts)
}

// Landing page for non-authenticated users
func landingPage(r *ui.Request) error {
	content := Div(
		Class("min-h-screen bg-white"),

		// Main hero section
		Div(
			Class("max-w-6xl mx-auto px-4 py-20"),
			Div(
				Class("flex flex-col lg:flex-row items-center gap-16"),
				// Left side - content
				Div(
					Class("flex-1 text-left"),
					H1(
						Class("text-4xl lg:text-5xl font-bold text-gray-900 mb-6 leading-tight"),
						Text("Learning Platform for "),
						Br(),
						Text("Students & Educators"),
					),
					P(
						Class("text-lg text-gray-600 mb-8 leading-relaxed max-w-lg"),
						Text("Completely free, comprehensive learning platform for students. Covers quizzes, document management, and study tools you are likely to need. Concrete, no-nonsense tools for the learner in a hurry."),
					),
					Div(
						Class("flex gap-4"),
						A(
							Href("/user/register"),
							Class("bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg font-medium transition-colors"),
							Text("Get started"),
						),
						A(
							Href("#features"),
							Class("text-blue-600 hover:text-blue-700 px-6 py-3 font-medium transition-colors"),
							Text("See what we cover"),
							Span(Class("ml-1"), Text("→")),
						),
					),
				),
				// Right side - illustration grid
				Div(
					Class("flex-1 grid grid-cols-2 gap-4 max-w-md"),
					// Study illustration
					Div(
						Class("bg-blue-50 rounded-lg p-6 flex flex-col items-center text-center"),
						Div(
							Class("w-16 h-16 bg-blue-100 rounded-lg mb-4 flex items-center justify-center"),
							Span(Class("text-2xl"), Text("📚")),
						),
						H3(Class("font-semibold text-gray-900 mb-2"), Text("Study Materials")),
						P(Class("text-sm text-gray-600"), Text("Organize and access your documents")),
					),
					// Quiz illustration
					Div(
						Class("bg-green-50 rounded-lg p-6 flex flex-col items-center text-center"),
						Div(
							Class("w-16 h-16 bg-green-100 rounded-lg mb-4 flex items-center justify-center"),
							Span(Class("text-2xl"), Text("🧠")),
						),
						H3(Class("font-semibold text-gray-900 mb-2"), Text("Smart Quizzes")),
						P(Class("text-sm text-gray-600"), Text("Test your knowledge effectively")),
					),
					// Progress illustration
					Div(
						Class("bg-purple-50 rounded-lg p-6 flex flex-col items-center text-center"),
						Div(
							Class("w-16 h-16 bg-purple-100 rounded-lg mb-4 flex items-center justify-center"),
							Span(Class("text-2xl"), Text("📊")),
						),
						H3(Class("font-semibold text-gray-900 mb-2"), Text("Track Progress")),
						P(Class("text-sm text-gray-600"), Text("Monitor your learning journey")),
					),
					// Collaboration illustration
					Div(
						Class("bg-orange-50 rounded-lg p-6 flex flex-col items-center text-center"),
						Div(
							Class("w-16 h-16 bg-orange-100 rounded-lg mb-4 flex items-center justify-center"),
							Span(Class("text-2xl"), Text("👥")),
						),
						H3(Class("font-semibold text-gray-900 mb-2"), Text("Collaborate")),
						P(Class("text-sm text-gray-600"), Text("Learn together with others")),
					),
				),
			),
		),

		// Second section with more details
		Div(
			ID("features"),
			Class("bg-gray-900 text-white py-20"),
			Div(
				Class("max-w-4xl mx-auto px-4 text-center"),
				H2(
					Class("text-3xl lg:text-4xl font-bold mb-8"),
					Text("Make your learning "),
					Br(),
					Text("experience into "),
					Br(),
					Text("academic success"),
				),
				P(
					Class("text-lg text-gray-300 mb-12 max-w-2xl mx-auto leading-relaxed"),
					Text("Our platform provides everything you need to excel in your studies. From interactive quizzes to document management, we've got you covered."),
				),
				A(
					Href("/user/register"),
					Class("bg-blue-600 hover:bg-blue-700 text-white px-8 py-4 rounded-lg font-medium text-lg transition-colors inline-block"),
					Text("Start Learning Today"),
				),
			),
		),
	)

	return r.Render(layouts.Primary, content)
}

// Home page for authenticated users showing their content in social media style
func authenticatedHomePage(r *ui.Request, posts *models.Posts) error {
	// If user is admin, redirect to admin dashboard
	if r.IsAdmin {
		return redirect.New(r.Context).URL("/admin").Go()
	}

	// If user is not verified, show unverified user home page
	if r.AuthUser != nil && !r.AuthUser.Verified {
		return unverifiedUserHomePage(r)
	}

	content := Div(
		Class("max-w-4xl mx-auto space-y-6"),
		
		// Status Update Box
		Div(
			Class("bg-white rounded-lg shadow-sm border border-gray-200 p-4"),
			Div(
				Class("flex items-center gap-3 mb-4"),
				Div(
					Class("w-10 h-10 bg-blue-100 rounded-full flex items-center justify-center"),
					Span(Class("text-blue-600 font-semibold"),
						Text(func() string {
							if r.AuthUser != nil && r.AuthUser.Name != "" {
								return string(r.AuthUser.Name[0])
							}
							return "U"
						}()),
					),
				),
				Form(
					Class("flex-1"),
					Method("GET"),
					Action("/search"),
					Input(
						Type("search"),
						Name("q"),
						Class("w-full bg-gray-50 border border-gray-200 rounded-full px-4 py-2 text-gray-700 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"),
						Placeholder("Search notes, quizzes, files..."),
					),
				),
			),
			Div(
				Class("flex items-center justify-between pt-2 border-t border-gray-100"),
				Div(
					Class("flex gap-4"),
					// Only show action buttons for verified users
					If(r.AuthUser != nil && r.AuthUser.Verified, Group{
						A(
							Href("/notes/create"),
							Class("flex items-center gap-2 px-3 py-2 rounded-lg text-gray-600 hover:text-amber-700 hover:bg-amber-50 transition-all duration-200"),
							Span(Class("text-lg"), Text("📝")),
							Span(Class("font-medium"), Text("Note")),
						),
						A(
							Href("/quizzes/create"),
							Class("flex items-center gap-2 px-3 py-2 rounded-lg text-gray-600 hover:text-blue-700 hover:bg-blue-50 transition-all duration-200"),
							Span(Class("text-lg"), Text("🧠")),
							Span(Class("font-medium"), Text("Quiz")),
						),
						A(
							Href("/files"),
							Class("flex items-center gap-2 px-3 py-2 rounded-lg text-gray-600 hover:text-emerald-700 hover:bg-emerald-50 transition-all duration-200"),
							Span(Class("text-lg"), Text("📄")),
							Span(Class("font-medium"), Text("File")),
						),
					}),
					// Show verification message for unverified users
					If(r.AuthUser != nil && !r.AuthUser.Verified,
						Div(
							Class("flex items-center gap-2 px-3 py-2 rounded-lg bg-yellow-50 border border-yellow-200 text-yellow-700"),
							Span(Class("text-lg"), Text("⚠️")),
							Span(Class("font-medium text-sm"), Text("Please verify your account to access creation tools")),
						),
					),
				),
			),
		),

		// Social Media Feed Posts
		// Note Post
		Div(
			Class("bg-white rounded-lg shadow-sm border-l-4 border-l-amber-400 border-t border-r border-b border-gray-200"),
			// Post Header
			Div(
				Class("flex items-center gap-3 p-4 border-b border-gray-100"),
				Div(
					Class("w-10 h-10 bg-amber-100 rounded-full flex items-center justify-center"),
					Span(Class("text-amber-700 font-semibold"), Text("S")),
				),
				Div(
					Class("flex-1"),
					H3(Class("font-semibold text-gray-900"), Text("Sarah Chen")),
					P(Class("text-sm text-gray-500 flex items-center gap-1"), 
						Text("2 hours ago • "),
						Span(Class("inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium bg-amber-100 text-amber-800"), 
							Text("📝 Note"),
						),
					),
				),
				Button(
					Class("text-gray-400 hover:text-gray-600"),
					Text("⋯"),
				),
			),
			// Post Content
			Div(
				Class("p-4"),
				H4(Class("font-semibold text-gray-900 mb-2"), Text("Advanced Calculus - Integration Techniques")),
				P(Class("text-gray-700 mb-3"), Text("Just finished my notes on integration by parts and substitution methods. The key insight is recognizing patterns in the integrand...")),
				Div(
					Class("bg-gradient-to-r from-amber-50 to-yellow-50 border border-amber-200 rounded-lg p-3"),
					P(Class("text-sm text-amber-800 font-medium"), Text("💡 Pro tip: Always look for the derivative of one part in the other when using integration by parts!")),
				),
			),
			// Post Actions
			Div(
				Class("flex items-center justify-between px-4 py-3 border-t border-gray-100"),
				Div(
					Class("flex items-center gap-6"),
					Button(
						Class("flex items-center gap-2 text-gray-600 hover:text-red-600 transition-colors"),
						Span(Text("❤️")),
						Span(Class("text-sm"), Text("12")),
					),
					Button(
						Class("flex items-center gap-2 text-gray-600 hover:text-blue-600 transition-colors"),
						Span(Text("💬")),
						Span(Class("text-sm"), Text("5")),
					),
					Button(
						Class("flex items-center gap-2 text-gray-600 hover:text-green-600 transition-colors"),
						Span(Text("🔄")),
						Span(Class("text-sm"), Text("Share")),
					),
				),
				Button(
					Class("text-gray-600 hover:text-gray-800 transition-colors"),
					Span(Text("🔖")),
				),
			),
		),

		// Quiz Post
		Div(
			Class("bg-white rounded-lg shadow-sm border-l-4 border-l-blue-400 border-t border-r border-b border-gray-200"),
			// Post Header
			Div(
				Class("flex items-center gap-3 p-4 border-b border-gray-100"),
				Div(
					Class("w-10 h-10 bg-blue-100 rounded-full flex items-center justify-center"),
					Span(Class("text-blue-700 font-semibold"), Text("M")),
				),
				Div(
					Class("flex-1"),
					H3(Class("font-semibold text-gray-900"), Text("Mike Rodriguez")),
					P(Class("text-sm text-gray-500 flex items-center gap-1"), 
						Text("4 hours ago • "),
						Span(Class("inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800"), 
							Text("🧠 Quiz"),
						),
					),
				),
				Button(
					Class("text-gray-400 hover:text-gray-600"),
					Text("⋯"),
				),
			),
			// Post Content
			Div(
				Class("p-4"),
				H4(Class("font-semibold text-gray-900 mb-2"), Text("JavaScript Fundamentals Quiz")),
				P(Class("text-gray-700 mb-3"), Text("Created a comprehensive quiz covering variables, functions, and async programming. Perfect for beginners!")),
				Div(
					Class("bg-gradient-to-r from-blue-50 to-indigo-50 border border-blue-200 rounded-lg p-4"),
					Div(
						Class("flex items-center justify-between mb-2"),
						Span(Class("text-blue-800 font-semibold"), Text("15 Questions")),
						Span(Class("text-blue-600 text-sm font-medium"), Text("~20 min")),
					),
					P(Class("text-blue-700 text-sm font-medium"), Text("Topics: Variables, Functions, Promises, DOM Manipulation")),
				),
			),
			// Post Actions
			Div(
				Class("flex items-center justify-between px-4 py-3 border-t border-gray-100"),
				Div(
					Class("flex items-center gap-6"),
					Button(
						Class("flex items-center gap-2 text-gray-600 hover:text-red-600 transition-colors"),
						Span(Text("❤️")),
						Span(Class("text-sm"), Text("24")),
					),
					Button(
						Class("flex items-center gap-2 text-gray-600 hover:text-blue-600 transition-colors"),
						Span(Text("💬")),
						Span(Class("text-sm"), Text("8")),
					),
					Button(
						Class("flex items-center gap-2 text-gray-600 hover:text-green-600 transition-colors"),
						Span(Text("🔄")),
						Span(Class("text-sm"), Text("Share")),
					),
				),
				// Only show Take Quiz button for verified users
				If(r.AuthUser != nil && r.AuthUser.Verified,
					Button(
						Class("bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 shadow-sm hover:shadow-md"),
						Text("Take Quiz"),
					),
				),
				// Show verification message for unverified users
				If(r.AuthUser != nil && !r.AuthUser.Verified,
					Div(
						Class("bg-yellow-50 border border-yellow-200 text-yellow-700 px-4 py-2 rounded-lg text-sm font-medium"),
						Text("⚠️ Verify account to take quizzes"),
					),
				),
			),
		),

		// File Post
		Div(
			Class("bg-white rounded-lg shadow-sm border-l-4 border-l-emerald-400 border-t border-r border-b border-gray-200"),
			// Post Header
			Div(
				Class("flex items-center gap-3 p-4 border-b border-gray-100"),
				Div(
					Class("w-10 h-10 bg-emerald-100 rounded-full flex items-center justify-center"),
					Span(Class("text-emerald-700 font-semibold"), Text("A")),
				),
				Div(
					Class("flex-1"),
					H3(Class("font-semibold text-gray-900"), Text("Alex Johnson")),
					P(Class("text-sm text-gray-500 flex items-center gap-1"), 
						Text("6 hours ago • "),
						Span(Class("inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium bg-emerald-100 text-emerald-800"), 
							Text("📄 File"),
						),
					),
				),
				Button(
					Class("text-gray-400 hover:text-gray-600"),
					Text("⋯"),
				),
			),
			// Post Content
			Div(
				Class("p-4"),
				H4(Class("font-semibold text-gray-900 mb-2"), Text("Machine Learning Cheat Sheet")),
				P(Class("text-gray-700 mb-3"), Text("Uploaded a comprehensive ML cheat sheet covering algorithms, evaluation metrics, and best practices. Great for quick reference!")),
				Div(
					Class("bg-gradient-to-r from-emerald-50 to-green-50 border border-emerald-200 rounded-lg p-4 flex items-center gap-3"),
					Div(
						Class("w-12 h-12 bg-emerald-100 rounded-lg flex items-center justify-center"),
						Span(Class("text-emerald-600 text-lg"), Text("📄")),
					),
					Div(
						Class("flex-1"),
						P(Class("font-semibold text-emerald-800"), Text("ML_CheatSheet_2024.pdf")),
						P(Class("text-sm text-emerald-600 font-medium"), Text("2.4 MB • PDF Document")),
					),
					// Only show Download button for verified users
					If(r.AuthUser != nil && r.AuthUser.Verified,
						Button(
							Class("bg-emerald-600 hover:bg-emerald-700 text-white px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 shadow-sm hover:shadow-md"),
							Text("Download"),
						),
					),
					// Show verification message for unverified users
					If(r.AuthUser != nil && !r.AuthUser.Verified,
						Div(
							Class("bg-yellow-50 border border-yellow-200 text-yellow-700 px-4 py-2 rounded-lg text-sm font-medium"),
							Text("⚠️ Verify account to download files"),
						),
					),
				),
			),
			// Post Actions
			Div(
				Class("flex items-center justify-between px-4 py-3 border-t border-gray-100"),
				Div(
					Class("flex items-center gap-6"),
					Button(
						Class("flex items-center gap-2 text-gray-600 hover:text-red-600 transition-colors"),
						Span(Text("❤️")),
						Span(Class("text-sm"), Text("18")),
					),
					Button(
						Class("flex items-center gap-2 text-gray-600 hover:text-blue-600 transition-colors"),
						Span(Text("💬")),
						Span(Class("text-sm"), Text("3")),
					),
					Button(
						Class("flex items-center gap-2 text-gray-600 hover:text-green-600 transition-colors"),
						Span(Text("🔄")),
						Span(Class("text-sm"), Text("Share")),
					),
				),
				Button(
					Class("text-gray-600 hover:text-gray-800 transition-colors"),
					Span(Text("🔖")),
				),
			),
		),

		// Study Group Post
		Div(
			Class("bg-white rounded-lg shadow-sm border border-gray-200 hover:shadow-lg hover:shadow-pink-500/20 transition-all duration-200"),
			// Post Header
			Div(
				Class("flex items-center gap-3 p-4 border-b border-gray-100"),
				Div(
					Class("w-10 h-10 bg-pink-100 rounded-full flex items-center justify-center"),
					Span(Class("text-pink-600 font-semibold"), Text("E")),
				),
				Div(
					Class("flex-1"),
					H3(Class("font-semibold text-gray-900"), Text("Emma Wilson")),
					P(Class("text-sm text-gray-500"), Text("1 day ago • 👥 Study Group")),
				),
				Button(
					Class("text-gray-400 hover:text-gray-600"),
					Text("⋯"),
				),
			),
			// Post Content
			Div(
				Class("p-4"),
				H4(Class("font-semibold text-gray-900 mb-2"), Text("Data Structures Study Group - Week 3")),
				P(Class("text-gray-700 mb-3"), Text("Great session today! We covered binary trees and graph algorithms. Thanks everyone for the engaging discussions! 🌟")),
				Div(
					Class("bg-pink-50 border border-pink-200 rounded-lg p-3"),
					Div(
						Class("flex items-center gap-2 mb-2"),
						Span(Class("text-pink-600 font-medium"), Text("📅 Next Session:")),
						Span(Class("text-pink-800"), Text("Friday 3PM - Library Room 204")),
					),
					P(Class("text-pink-700 text-sm"), Text("Topic: Dynamic Programming & Memoization")),
				),
			),
			// Post Actions
			Div(
				Class("flex items-center justify-between px-4 py-3 border-t border-gray-100"),
				Div(
					Class("flex items-center gap-6"),
					Button(
						Class("flex items-center gap-2 text-gray-600 hover:text-red-600 transition-colors"),
						Span(Text("❤️")),
						Span(Class("text-sm"), Text("31")),
					),
					Button(
						Class("flex items-center gap-2 text-gray-600 hover:text-blue-600 transition-colors"),
						Span(Text("💬")),
						Span(Class("text-sm"), Text("12")),
					),
					Button(
						Class("flex items-center gap-2 text-gray-600 hover:text-green-600 transition-colors"),
						Span(Text("🔄")),
						Span(Class("text-sm"), Text("Share")),
					),
				),
				// Only show Join Group button for verified users
				If(r.AuthUser != nil && r.AuthUser.Verified,
					Button(
						Class("bg-pink-600 hover:bg-pink-700 text-white px-4 py-1 rounded-full text-sm transition-colors"),
						Text("Join Group"),
					),
				),
				// Show verification message for unverified users
				If(r.AuthUser != nil && !r.AuthUser.Verified,
					Div(
						Class("bg-yellow-50 border border-yellow-200 text-yellow-700 px-4 py-1 rounded-full text-sm font-medium"),
						Text("⚠️ Verify account to join groups"),
					),
				),
			),
		),
	)

	return r.Render(layouts.Primary, content)
}

// Dashboard for authenticated users
func authenticatedDashboard(r *ui.Request, posts *models.Posts) error {
	// If user is admin, redirect to admin dashboard
	if r.IsAdmin {
		return redirect.New(r.Context).URL("/admin").Go()
	}
	headerMsg := func() Node {
		return Group{
			Div(
				Class("space-y-6"),
				// Simple Welcome Card
				Div(
					Class("bg-white rounded-lg p-6 shadow-sm border border-gray-200"),
					Div(
						Class("flex items-center gap-4"),
						// Profile Avatar
						Div(
							Class("w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center"),
							Span(Class("text-blue-600 font-semibold text-lg"),
								Text(func() string {
									if r.AuthUser != nil && r.AuthUser.Name != "" {
										return string(r.AuthUser.Name[0])
									}
									return "U"
								}()),
							),
						),
						Div(
							Class("flex-1"),
							H2(
								Class("text-lg font-semibold text-gray-900"),
								Text("Welcome back"),
								Iff(r.AuthUser != nil && r.AuthUser.Name != "", func() Node {
									return Text(", " + r.AuthUser.Name)
								}),
								Text("!"),
							),
							P(
								Class("text-sm text-gray-500"),
								Text("Ready to continue your learning journey?"),
							),
						),
					),
				),
			),
		}
	}

	cards := func() Node {
		return Div(
			Class("space-y-6"),
			// Quick Actions Grid - Only show for verified users
			If(r.AuthUser != nil && r.AuthUser.Verified,
				Div(
					Class("grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6"),
					// Create Note Card
					A(
						Href("/notes/create"),
						Class("block bg-white p-6 rounded-lg border-l-4 border-l-amber-400 border-t border-r border-b border-gray-200 hover:border-amber-300 hover:shadow-lg hover:shadow-amber-500/20 transition-all duration-200"),
						Div(
							Class("flex items-center gap-4"),
							Div(
								Class("w-12 h-12 bg-amber-100 rounded-lg flex items-center justify-center"),
								Span(Class("text-amber-600 text-xl"), Text("📝")),
							),
							Div(
								H3(Class("font-semibold text-gray-900"), Text("Create Note")),
								P(Class("text-sm text-amber-600 font-medium"), Text("Write and organize notes")),
							),
						),
					),
					// Create Quiz Card
					A(
						Href("/quiz/create"),
						Class("block bg-white p-6 rounded-lg border-l-4 border-l-blue-400 border-t border-r border-b border-gray-200 hover:border-blue-300 hover:shadow-lg hover:shadow-blue-500/20 transition-all duration-200"),
						Div(
							Class("flex items-center gap-4"),
							Div(
								Class("w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center"),
								Span(Class("text-blue-600 text-xl"), Text("🧠")),
							),
							Div(
								H3(Class("font-semibold text-gray-900"), Text("Create Quiz")),
								P(Class("text-sm text-blue-600 font-medium"), Text("Build interactive quizzes")),
							),
						),
					),
					// Upload Documents Card
					A(
						Href("/documents/upload"),
						Class("block bg-white p-6 rounded-lg border-l-4 border-l-emerald-400 border-t border-r border-b border-gray-200 hover:border-emerald-300 hover:shadow-lg hover:shadow-emerald-500/20 transition-all duration-200"),
						Div(
							Class("flex items-center gap-4"),
							Div(
								Class("w-12 h-12 bg-emerald-100 rounded-lg flex items-center justify-center"),
								Span(Class("text-emerald-600 text-xl"), Text("📄")),
							),
							Div(
								H3(Class("font-semibold text-gray-900"), Text("Upload Files")),
								P(Class("text-sm text-emerald-600 font-medium"), Text("Manage your documents")),
							),
						),
					),
					// Study Group Session Card
					A(
						Href("/study-groups"),
						Class("block bg-white p-6 rounded-lg border-l-4 border-l-pink-400 border-t border-r border-b border-gray-200 hover:border-pink-300 hover:shadow-lg hover:shadow-pink-500/20 transition-all duration-200"),
						Div(
							Class("flex items-center gap-4"),
							Div(
								Class("w-12 h-12 bg-pink-100 rounded-lg flex items-center justify-center"),
								Span(Class("text-pink-600 text-xl"), Text("👥")),
							),
							Div(
								H3(Class("font-semibold text-gray-900"), Text("Study Groups")),
								P(Class("text-sm text-pink-600 font-medium"), Text("Join collaborative sessions")),
							),
						),
					),
					// View Progress Card
					A(
						Href("/progress"),
						Class("block bg-white p-6 rounded-lg border border-gray-200 hover:border-purple-300 hover:shadow-md transition-all"),
						Div(
							Class("flex items-center gap-4"),
							Div(
								Class("w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center"),
								Span(Class("text-purple-600 text-xl"), Text("📊")),
							),
							Div(
								H3(Class("font-semibold text-gray-900"), Text("View Progress")),
								P(Class("text-sm text-gray-500"), Text("Track your learning")),
							),
						),
					),
				),
			),
			// Show verification message for unverified users
			If(r.AuthUser != nil && !r.AuthUser.Verified,
				Div(
					Class("bg-yellow-50 border border-yellow-200 rounded-lg p-6"),
					Div(
						Class("text-center"),
						Div(
							Class("w-16 h-16 bg-yellow-100 rounded-full flex items-center justify-center mx-auto mb-4"),
							Span(Class("text-yellow-600 text-2xl"), Text("⚠️")),
						),
						H3(Class("text-lg font-semibold text-yellow-800 mb-2"), Text("Account Verification Required")),
					P(Class("text-yellow-700 mb-4"), Text("Please contact an administrator to verify your account and access creation tools, quizzes, file downloads, and study groups.")),
					),
				),
			),
			// Recent Activity Card
			Div(
				Class("bg-white rounded-lg p-6 shadow-sm border border-gray-200"),
				H3(Class("text-lg font-semibold text-gray-900 mb-4"), Text("Recent Activity")),
				P(Class("text-gray-500"), Text("Your recent quizzes and documents will appear here.")),
			),
			// Study Statistics Card
			Div(
				Class("bg-white rounded-lg p-6 shadow-sm border border-gray-200"),
				H3(Class("text-lg font-semibold text-gray-900 mb-4"), Text("Study Statistics")),
				P(Class("text-gray-500"), Text("Your learning progress and achievements will be displayed here.")),
			),
		)
	}

	g := Group{
		Iff(r.Htmx.Target != "posts", headerMsg),
		Iff(r.Htmx.Target != "posts", cards),
	}

	return r.Render(layouts.Primary, g)
}

// Home page for unverified users with limited content
func unverifiedUserHomePage(r *ui.Request) error {
	content := Div(
		Class("max-w-4xl mx-auto space-y-6"),
		
		// Simple welcome message for unverified users
		Div(
			Class("bg-gradient-to-r from-amber-50 to-yellow-50 border border-amber-200 rounded-lg p-6 text-center"),
			Div(
				Class("w-16 h-16 bg-amber-100 rounded-full flex items-center justify-center mx-auto mb-4"),
				Span(Class("text-2xl"), Text("📱")),
			),
			H2(
				Class("text-xl font-semibold text-amber-800 mb-2"),
				Text("Welcome to Zero!"),
			),
			P(
				Class("text-amber-700"),
				Text("Please verify your WhatsApp number to access all features and content."),
			),
		),
		
		// Read-only social media feed preview
		// Note Post
		Div(
			Class("bg-white rounded-lg shadow-sm border-l-4 border-l-amber-400 border-t border-r border-b border-gray-200 opacity-75"),
			// Post Header
			Div(
				Class("flex items-center gap-3 p-4 border-b border-gray-100"),
				Div(
					Class("w-10 h-10 bg-amber-100 rounded-full flex items-center justify-center"),
					Span(Class("text-amber-700 font-semibold"), Text("S")),
				),
				Div(
					Class("flex-1"),
					H3(Class("font-semibold text-gray-900"), Text("Sarah Chen")),
					P(Class("text-sm text-gray-500 flex items-center gap-1"), 
						Text("2 hours ago • "),
						Span(Class("inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium bg-amber-100 text-amber-800"), 
							Text("📝 Note"),
						),
					),
				),
			),
			// Post Content
			Div(
				Class("p-4"),
				H4(Class("font-semibold text-gray-900 mb-2"), Text("Advanced Calculus - Integration Techniques")),
				P(Class("text-gray-700 mb-3"), Text("Just finished my notes on integration by parts and substitution methods. The key insight is recognizing patterns in the integrand...")),
				Div(
					Class("bg-gradient-to-r from-amber-50 to-yellow-50 border border-amber-200 rounded-lg p-3"),
					P(Class("text-sm text-amber-800 font-medium"), Text("💡 Pro tip: Always look for the derivative of one part in the other when using integration by parts!")),
				),
			),

		),

		// Quiz Post
		Div(
			Class("bg-white rounded-lg shadow-sm border-l-4 border-l-blue-400 border-t border-r border-b border-gray-200 opacity-75"),
			// Post Header
			Div(
				Class("flex items-center gap-3 p-4 border-b border-gray-100"),
				Div(
					Class("w-10 h-10 bg-blue-100 rounded-full flex items-center justify-center"),
					Span(Class("text-blue-700 font-semibold"), Text("M")),
				),
				Div(
					Class("flex-1"),
					H3(Class("font-semibold text-gray-900"), Text("Mike Rodriguez")),
					P(Class("text-sm text-gray-500 flex items-center gap-1"), 
						Text("4 hours ago • "),
						Span(Class("inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800"), 
							Text("🧠 Quiz"),
						),
					),
				),
			),
			// Post Content
			Div(
				Class("p-4"),
				H4(Class("font-semibold text-gray-900 mb-2"), Text("JavaScript Fundamentals Quiz")),
				P(Class("text-gray-700 mb-3"), Text("Created a comprehensive quiz covering variables, functions, and async programming. Perfect for beginners!")),
				Div(
					Class("bg-gradient-to-r from-blue-50 to-indigo-50 border border-blue-200 rounded-lg p-4"),
					Div(
						Class("flex items-center justify-between mb-2"),
						Span(Class("text-blue-800 font-semibold"), Text("15 Questions")),
						Span(Class("text-blue-600 text-sm font-medium"), Text("~20 min")),
					),
					P(Class("text-blue-700 text-sm font-medium"), Text("Topics: Variables, Functions, Promises, DOM Manipulation")),
				),
			),

		),

		// File Post
		Div(
			Class("bg-white rounded-lg shadow-sm border-l-4 border-l-emerald-400 border-t border-r border-b border-gray-200 opacity-75"),
			// Post Header
			Div(
				Class("flex items-center gap-3 p-4 border-b border-gray-100"),
				Div(
					Class("w-10 h-10 bg-emerald-100 rounded-full flex items-center justify-center"),
					Span(Class("text-emerald-700 font-semibold"), Text("A")),
				),
				Div(
					Class("flex-1"),
					H3(Class("font-semibold text-gray-900"), Text("Alex Johnson")),
					P(Class("text-sm text-gray-500 flex items-center gap-1"), 
						Text("6 hours ago • "),
						Span(Class("inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium bg-emerald-100 text-emerald-800"), 
							Text("📄 File"),
						),
					),
				),
			),
			// Post Content
			Div(
				Class("p-4"),
				H4(Class("font-semibold text-gray-900 mb-2"), Text("Machine Learning Cheat Sheet")),
				P(Class("text-gray-700 mb-3"), Text("Uploaded a comprehensive ML cheat sheet covering algorithms, evaluation metrics, and best practices. Great for quick reference!")),
				Div(
					Class("bg-gradient-to-r from-emerald-50 to-green-50 border border-emerald-200 rounded-lg p-4 flex items-center gap-3"),
					Div(
						Class("w-12 h-12 bg-emerald-100 rounded-lg flex items-center justify-center"),
						Span(Class("text-emerald-600 text-lg"), Text("📄")),
					),
					Div(
						Class("flex-1"),
						P(Class("font-semibold text-emerald-800"), Text("ML_CheatSheet_2024.pdf")),
						P(Class("text-sm text-emerald-600 font-medium"), Text("2.4 MB • PDF Document")),
					),
					Div(
						Class("bg-yellow-50 border border-yellow-200 text-yellow-700 px-4 py-2 rounded-lg text-sm font-medium"),
						Text("⚠️ Verify account to download files"),
					),
				),
			),

		),

		// Study Group Post
		Div(
			Class("bg-white rounded-lg shadow-sm border border-gray-200 opacity-75"),
			// Post Header
			Div(
				Class("flex items-center gap-3 p-4 border-b border-gray-100"),
				Div(
					Class("w-10 h-10 bg-pink-100 rounded-full flex items-center justify-center"),
					Span(Class("text-pink-600 font-semibold"), Text("E")),
				),
				Div(
					Class("flex-1"),
					H3(Class("font-semibold text-gray-900"), Text("Emma Wilson")),
					P(Class("text-sm text-gray-500"), Text("1 day ago • 👥 Study Group")),
				),
			),
			// Post Content
			Div(
				Class("p-4"),
				H4(Class("font-semibold text-gray-900 mb-2"), Text("Data Structures Study Group - Week 3")),
				P(Class("text-gray-700 mb-3"), Text("Great session today! We covered binary trees and graph algorithms. Thanks everyone for the engaging discussions! 🌟")),
				Div(
					Class("bg-pink-50 border border-pink-200 rounded-lg p-3"),
					Div(
						Class("flex items-center gap-2 mb-2"),
						Span(Class("text-pink-600 font-medium"), Text("📅 Next Session:")),
						Span(Class("text-pink-800"), Text("Friday 3PM - Library Room 204")),
					),
					P(Class("text-pink-700 text-sm"), Text("Topic: Dynamic Programming & Memoization")),
				),
			),

		),
	)

	return r.Render(layouts.Primary, content)
}
