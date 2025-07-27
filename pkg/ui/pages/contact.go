package pages

import (
	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/pkg/ui"
	"github.com/r-scheele/zero/pkg/ui/forms"
	"github.com/r-scheele/zero/pkg/ui/layouts"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func ContactUs(ctx echo.Context, form *forms.Contact) error {
	r := ui.NewRequest(ctx)
	r.Metatags.Description = "Get in touch with the Zero team. We'd love to hear from you and help you dominate your learning journey."

	content := Div(
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
						Text("Let's Connect"),
					),
					// Main heading
					H1(
						Class("text-5xl sm:text-6xl lg:text-7xl font-black text-white mb-8 leading-tight"),
						Text("We're Here "),
						Br(),
						Span(
							Class("bg-gradient-to-r from-emerald-400 to-cyan-400 bg-clip-text text-transparent"),
							Text("to Help"),
						),
					),
					// Subheading
					P(
						Class("text-xl sm:text-2xl text-slate-300 mb-12 max-w-4xl mx-auto leading-relaxed"),
						Text("Got questions? Feedback? Ideas that could make Zero even better? We're all ears and respond fast."),
					),
				),
			),
		),
		// Contact Options Section
		Div(
			Class("py-24 bg-white relative"),
			Div(
				Class("max-w-7xl mx-auto px-4 sm:px-6 lg:px-8"),
				// Quick Contact Options
				Iff(!form.IsDone(), func() Node {
					return Div(
						Class("grid grid-cols-1 md:grid-cols-3 gap-8 mb-20"),
						// Email
						Div(
							Class("group text-center p-8 rounded-3xl bg-gradient-to-br from-blue-50 to-indigo-50 hover:shadow-xl transition-all duration-300 hover:-translate-y-2"),
							Div(
								Class("w-16 h-16 bg-blue-500 rounded-2xl flex items-center justify-center mx-auto mb-6 group-hover:scale-110 transition-transform"),
								Span(Class("text-white text-2xl"), Text("📧")),
							),
							H3(Class("text-2xl font-bold text-slate-900 mb-4"), Text("Email Us")),
							P(Class("text-slate-600 mb-6"), Text("Drop us a line anytime. We typically respond within a few hours.")),
							A(
								Href("mailto:hello@zero.com"),
								Class("text-blue-600 hover:text-blue-700 font-semibold"),
								Text("hello@zero.com"),
							),
						),
						// Social
						Div(
							Class("group text-center p-8 rounded-3xl bg-gradient-to-br from-purple-50 to-pink-50 hover:shadow-xl transition-all duration-300 hover:-translate-y-2"),
							Div(
								Class("w-16 h-16 bg-purple-500 rounded-2xl flex items-center justify-center mx-auto mb-6 group-hover:scale-110 transition-transform"),
								Span(Class("text-white text-2xl"), Text("💬")),
							),
							H3(Class("text-2xl font-bold text-slate-900 mb-4"), Text("Social Media")),
							P(Class("text-slate-600 mb-6"), Text("Follow us for updates, tips, and behind-the-scenes content.")),
							Div(Class("flex justify-center gap-4"),
								A(Href("#"), Class("text-purple-600 hover:text-purple-700"), Text("Twitter")),
								A(Href("#"), Class("text-purple-600 hover:text-purple-700"), Text("LinkedIn")),
							),
						),
						// Form teaser
						Div(
							Class("group text-center p-8 rounded-3xl bg-gradient-to-br from-emerald-50 to-green-50 hover:shadow-xl transition-all duration-300 hover:-translate-y-2"),
							Div(
								Class("w-16 h-16 bg-emerald-500 rounded-2xl flex items-center justify-center mx-auto mb-6 group-hover:scale-110 transition-transform"),
								Span(Class("text-white text-2xl"), Text("✍️")),
							),
							H3(Class("text-2xl font-bold text-slate-900 mb-4"), Text("Quick Message")),
							P(Class("text-slate-600 mb-6"), Text("Use the form below for detailed inquiries or feedback.")),
							Span(Class("text-emerald-600 font-semibold"), Text("Scroll down ↓")),
						),
					)
				}),
				// Demo Notice - only show if not in HTMX target and form not done
				Iff(r.Htmx.Target != "contact" && !form.IsDone(), func() Node {
					return Div(
						Class("bg-gradient-to-br from-amber-50 to-orange-50 rounded-3xl p-8 lg:p-12 shadow-lg border border-amber-200/60 mb-12"),
						Div(
							Class("flex items-start gap-6"),
							Div(
								Class("w-16 h-16 bg-amber-500 rounded-2xl flex items-center justify-center flex-shrink-0"),
								Span(Class("text-white text-2xl"), Text("⚡")),
							),
							Div(
								H3(
									Class("text-2xl font-bold text-amber-900 mb-4"),
									Text("See Our Tech in Action"),
								),
								P(
									Class("text-amber-800 leading-relaxed mb-4 text-lg"),
									Text("This contact form demonstrates Zero's lightning-fast, server-side validation and HTMX-powered submissions. No JavaScript fatigue, just smooth UX."),
								),
								Div(
									Class("flex flex-wrap gap-2"),
									Span(Class("px-3 py-1 bg-amber-200 text-amber-800 rounded-full text-sm font-medium"), Text("Server-side validation")),
									Span(Class("px-3 py-1 bg-amber-200 text-amber-800 rounded-full text-sm font-medium"), Text("HTMX AJAX")),
									Span(Class("px-3 py-1 bg-amber-200 text-amber-800 rounded-full text-sm font-medium"), Text("Zero JavaScript")),
								),
							),
						),
					)
				}),
			),
		),
		// Success Message
		Iff(form.IsDone(), func() Node {
			return Div(
				Class("py-24 bg-white"),
				Div(
					Class("max-w-4xl mx-auto text-center px-4 sm:px-6 lg:px-8"),
					Div(
						Class("bg-gradient-to-br from-emerald-50 to-teal-50 rounded-3xl p-12 lg:p-16 shadow-xl border border-emerald-200/60"),
						Div(
							Class("w-24 h-24 bg-emerald-500 rounded-full flex items-center justify-center mx-auto mb-8"),
							Span(Class("text-white text-4xl"), Text("🎉")),
						),
						H2(
							Class("text-4xl font-black text-emerald-900 mb-6"),
							Text("Message Sent!"),
						),
						P(
							Class("text-emerald-800 text-xl leading-relaxed mb-8"),
							Text("Thanks for reaching out! While this is just a demo (no actual email was sent), you just experienced Zero's smooth, server-side form handling."),
						),
						Div(
							Class("bg-emerald-100 rounded-2xl p-6 mb-8"),
							P(
								Class("text-emerald-700 font-medium"),
								Text("💡 Cool fact: This entire interaction happened without any client-side JavaScript. Everything was handled server-side for maximum performance and reliability."),
							),
						),
						A(
							Href("/contact"),
							Class("inline-flex items-center bg-emerald-600 hover:bg-emerald-700 text-white px-8 py-4 rounded-2xl font-bold text-lg transition-all duration-300 shadow-lg hover:shadow-xl"),
							Text("Send Another Message"),
							Span(Class("ml-2"), Text("→")),
						),
					),
				),
			)
		}),
		// Form Section
		Iff(!form.IsDone(), func() Node {
			return Div(
				Class("py-24 bg-slate-50"),
				Div(
					Class("max-w-4xl mx-auto px-4 sm:px-6 lg:px-8"),
					Div(
						Class("text-center mb-12"),
						H2(
							Class("text-4xl font-black text-slate-900 mb-6"),
							Text("Send Us a "),
							Span(Class("text-emerald-600"), Text("Message")),
						),
						P(
							Class("text-xl text-slate-600"),
							Text("Tell us what's on your mind. We read every message and respond personally."),
						),
					),
					Div(
						Class("bg-white rounded-3xl p-8 lg:p-12 shadow-xl border border-slate-200/60"),
						form.Render(r),
					),
				),
			)
		}),
		// CTA Section (only show if form not done)
		Iff(!form.IsDone(), func() Node {
			return Div(
				Class("py-24 bg-gradient-to-r from-slate-900 to-slate-800 relative overflow-hidden"),
				Div(
					Class("absolute inset-0 bg-[url('data:image/svg+xml,%3csvg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 32 32\" width=\"32\" height=\"32\" fill=\"none\" stroke=\"rgb(148 163 184 / 0.1)\"%%3e%3cpath d=\"m0 2 2-2 2 2-2 2-2-2\" stroke-width=\"0.5\"/%3e%3c/svg%3e')]"),
				),
				Div(
					Class("relative max-w-4xl mx-auto text-center px-4 sm:px-6 lg:px-8"),
					H2(
						Class("text-4xl sm:text-5xl font-black text-white mb-6"),
						Text("Ready to "),
						Span(Class("bg-gradient-to-r from-emerald-400 to-cyan-400 bg-clip-text text-transparent"), Text("Get Started?")),
					),
					P(
						Class("text-xl text-slate-300 mb-12 leading-relaxed"),
						Text("Don't wait for the perfect moment. Your learning transformation starts with one click."),
					),
					A(
						Href("/user/register"),
						Class("group inline-flex items-center bg-gradient-to-r from-emerald-500 to-cyan-500 hover:from-emerald-600 hover:to-cyan-600 text-white px-12 py-6 rounded-2xl font-bold text-xl transition-all duration-300 shadow-2xl hover:shadow-emerald-500/25 hover:scale-105"),
						Text("Start Learning Now"),
						Span(Class("ml-3 text-2xl group-hover:translate-x-2 transition-transform"), Text("🚀")),
					),
				),
			)
		}),
	)

	return r.Render(layouts.Primary, content)
}
