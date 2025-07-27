package pages

import (
	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/pkg/ui"
	"github.com/r-scheele/zero/pkg/ui/cache"
	"github.com/r-scheele/zero/pkg/ui/layouts"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func About(ctx echo.Context) error {
	r := ui.NewRequest(ctx)
	r.Metatags.Description = "Discover the story behind Zero - revolutionizing how students and professionals learn, create, and dominate their fields."

	// The content is static, so we can render and cache it.
	content := cache.SetIfNotExists("pages.about.Content", func() Node {
		return Div(
			Class("min-h-screen"),
			// Hero Section
			Div(
				Class("relative overflow-hidden bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900"),
				// Background decoration
				Div(
					Class("absolute inset-0 bg-[url('data:image/svg+xml,%3csvg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 32 32\" width=\"32\" height=\"32\" fill=\"none\" stroke=\"rgb(148 163 184 / 0.05)\"%%3e%3cpath d=\"m0 2 2-2 2 2-2 2-2-2\" stroke-width=\"0.5\"/%3e%3c/svg%3e')] opacity-20"),
				),
				Div(
					Class("relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24 sm:py-32"),
					Div(
						Class("text-center"),
						// Badge
						Div(
							Class("inline-flex items-center px-4 py-2 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 text-sm font-medium mb-8"),
							Span(Class("w-2 h-2 bg-emerald-400 rounded-full mr-2")),
							Text("Our Mission"),
						),
						// Main heading
						H1(
							Class("text-5xl sm:text-6xl lg:text-7xl font-black text-white mb-8 leading-tight"),
							Text("About Zero"),
						),
						// Subheading
						P(
							Class("text-xl sm:text-2xl text-slate-300 mb-12 max-w-4xl mx-auto leading-relaxed"),
							Text("We're tired of outdated learning tools that feel like they're from 2010. Zero is built by students, for students - with the modern tech stack that powers today's best apps."),
						),
					),
				),
			),
			// Mission Section
			Div(
				Class("py-24 bg-white relative"),
				Div(
					Class("max-w-7xl mx-auto px-4 sm:px-6 lg:px-8"),
					Div(
						Class("text-center mb-20"),
						H2(
							Class("text-4xl sm:text-5xl font-black text-slate-900 mb-6"),
							Text("Why We Built "),
							Span(Class("text-emerald-600"), Text("Zero")),
						),
						P(
							Class("text-xl text-slate-600 max-w-3xl mx-auto leading-relaxed"),
							Text("Traditional learning platforms are slow, ugly, and frustrating. We wanted something that actually makes you excited to study."),
						),
					),
					Div(
						Class("grid grid-cols-1 lg:grid-cols-2 gap-16 items-center"),
						// Left side - Story
						Div(
							Class("space-y-8"),
							Div(
								Class("space-y-6"),
								H3(
									Class("text-3xl font-bold text-slate-900"),
									Text("The Problem We Solved"),
								),
								P(
									Class("text-lg text-slate-600 leading-relaxed"),
									Text("Studying shouldn't feel like a chore. But most learning platforms are clunky, outdated, and designed by people who haven't been in a classroom in decades."),
								),
								P(
									Class("text-lg text-slate-600 leading-relaxed"),
									Text("We built Zero because we believe learning should be fast, beautiful, and actually fun. No more waiting 30 seconds for a page to load. No more ugly interfaces that make you want to close your laptop."),
								),
							),
							// Features list
							Div(
								Class("space-y-4"),
								Div(Class("flex items-center gap-3"),
									Span(Class("w-6 h-6 bg-emerald-500 rounded-full flex items-center justify-center"),
										Span(Class("text-white text-sm font-bold"), Text("✓")),
									),
									Text("Lightning-fast performance (seriously, try it)"),
								),
								Div(Class("flex items-center gap-3"),
									Span(Class("w-6 h-6 bg-emerald-500 rounded-full flex items-center justify-center"),
										Span(Class("text-white text-sm font-bold"), Text("✓")),
									),
									Text("Modern design that doesn't hurt your eyes"),
								),
								Div(Class("flex items-center gap-3"),
									Span(Class("w-6 h-6 bg-emerald-500 rounded-full flex items-center justify-center"),
										Span(Class("text-white text-sm font-bold"), Text("✓")),
									),
									Text("Actually works the way your brain does"),
								),
							),
						),
						// Right side - Stats/Visual
						Div(
							Class("relative"),
							Div(
								Class("bg-gradient-to-br from-emerald-50 to-cyan-50 rounded-3xl p-8 text-center"),
								Div(
									Class("space-y-8"),
									Div(
										H4(Class("text-lg font-semibold text-slate-700 mb-2"), Text("Load Time")),
										Div(Class("text-4xl font-black text-emerald-600 mb-1"), Text("< 1s")),
										P(Class("text-slate-600"), Text("vs 10s+ on others")),
									),
									Div(
										H4(Class("text-lg font-semibold text-slate-700 mb-2"), Text("User Satisfaction")),
										Div(Class("text-4xl font-black text-emerald-600 mb-1"), Text("98%")),
										P(Class("text-slate-600"), Text("love the experience")),
									),
									Div(
										H4(Class("text-lg font-semibold text-slate-700 mb-2"), Text("Time Saved")),
										Div(Class("text-4xl font-black text-emerald-600 mb-1"), Text("5hrs")),
										P(Class("text-slate-600"), Text("per week on average")),
									),
								),
							),
						),
					),
				),
			),
			// Tech Stack Section
			Div(
				Class("py-24 bg-slate-50 relative"),
				Div(
					Class("max-w-7xl mx-auto px-4 sm:px-6 lg:px-8"),
					Div(
						Class("text-center mb-20"),
						H2(
							Class("text-4xl sm:text-5xl font-black text-slate-900 mb-6"),
							Text("Built with "),
							Span(Class("text-purple-600"), Text("Modern Tech")),
						),
						P(
							Class("text-xl text-slate-600 max-w-3xl mx-auto"),
							Text("We use the same technologies that power the world's fastest applications. No compromises on performance or user experience."),
						),
					),
					Div(
						Class("grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8"),
						// Go Backend
						Div(
							Class("group bg-white rounded-3xl p-8 hover:shadow-2xl transition-all duration-500 hover:-translate-y-2 border border-slate-200"),
							Div(
								Class("w-16 h-16 bg-blue-500 rounded-2xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform"),
								Span(Class("text-white text-2xl font-bold"), Text("Go")),
							),
							H3(Class("text-2xl font-bold text-slate-900 mb-4"), Text("Lightning Backend")),
							P(Class("text-slate-600 leading-relaxed mb-4"),
								Text("Built with Go - the same language that powers Google, Uber, and Netflix. Blazing fast performance that never keeps you waiting.")),
							Div(Class("text-sm text-slate-500"), Text("Echo • Ent ORM • Modern Architecture")),
						),
						// Modern Frontend
						Div(
							Class("group bg-white rounded-3xl p-8 hover:shadow-2xl transition-all duration-500 hover:-translate-y-2 border border-slate-200"),
							Div(
								Class("w-16 h-16 bg-gradient-to-r from-cyan-500 to-blue-500 rounded-2xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform"),
								Span(Class("text-white text-2xl"), Text("🎨")),
							),
							H3(Class("text-2xl font-bold text-slate-900 mb-4"), Text("Zero JavaScript")),
							P(Class("text-slate-600 leading-relaxed mb-4"),
								Text("Smooth interactions without the complexity. HTMX and Alpine.js give you modern UX without the JavaScript fatigue.")),
							Div(Class("text-sm text-slate-500"), Text("HTMX • Alpine.js • TailwindCSS")),
						),
						// Performance
						Div(
							Class("group bg-white rounded-3xl p-8 hover:shadow-2xl transition-all duration-500 hover:-translate-y-2 border border-slate-200"),
							Div(
								Class("w-16 h-16 bg-gradient-to-r from-emerald-500 to-green-500 rounded-2xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform"),
								Span(Class("text-white text-2xl"), Text("⚡")),
							),
							H3(Class("text-2xl font-bold text-slate-900 mb-4"), Text("Optimized Everything")),
							P(Class("text-slate-600 leading-relaxed mb-4"),
								Text("Smart caching, efficient rendering, and careful optimization at every layer. Built for speed from day one.")),
							Div(Class("text-sm text-slate-500"), Text("Server-side rendering • Smart caching • CDN")),
						),
					),
				),
			),
			// Call to Action
			Div(
				Class("py-24 bg-gradient-to-r from-slate-900 to-slate-800 relative overflow-hidden"),
				Div(
					Class("absolute inset-0 bg-[url('data:image/svg+xml,%3csvg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 32 32\" width=\"32\" height=\"32\" fill=\"none\" stroke=\"rgb(148 163 184 / 0.1)\"%%3e%3cpath d=\"m0 2 2-2 2 2-2 2-2-2\" stroke-width=\"0.5\"/%3e%3c/svg%3e')]"),
				),
				Div(
					Class("relative max-w-4xl mx-auto text-center px-4 sm:px-6 lg:px-8"),
					H2(
						Class("text-4xl sm:text-5xl font-black text-white mb-6"),
						Text("Ready to Experience"),
						Br(),
						Span(Class("bg-gradient-to-r from-emerald-400 to-cyan-400 bg-clip-text text-transparent"), Text("The Difference?")),
					),
					P(
						Class("text-xl text-slate-300 mb-12 leading-relaxed"),
						Text("See why students are switching from outdated platforms to Zero. Your study sessions will never be the same."),
					),
					A(
						Href("/user/register"),
						Class("group inline-flex items-center bg-gradient-to-r from-emerald-500 to-cyan-500 hover:from-emerald-600 hover:to-cyan-600 text-white px-12 py-6 rounded-2xl font-bold text-xl transition-all duration-300 shadow-2xl hover:shadow-emerald-500/25 hover:scale-105"),
						Text("Try Zero Free"),
						Span(Class("ml-3 text-2xl group-hover:translate-x-2 transition-transform"), Text("→")),
					),
				),
			),
		)
	})

	return r.Render(layouts.Primary, content)
}
