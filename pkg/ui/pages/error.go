package pages

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/pkg/routenames"
	"github.com/r-scheele/zero/pkg/ui"
	"github.com/r-scheele/zero/pkg/ui/layouts"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Error(ctx echo.Context, code int) error {
	r := ui.NewRequest(ctx)
	// Remove title from page
	r.Title = ""

	content := Div(
		Class("max-w-xl mx-auto text-center space-y-4"),
		// Error illustration - slightly smaller for simpler look
		Div(
			Class("flex justify-center mb-2"),
			Div(
				Class("w-32 h-32 sm:w-40 sm:h-40"),
				Raw(getErrorIllustration(code)),
			),
		),
		Div(
			Class("space-y-2"),
			H1(
				Class("text-2xl font-bold text-slate-900"),
				Text(getErrorTitle(code)),
			),
			P(
				Class("text-sm text-slate-600"),
				Text(getErrorMessage(code)),
			),
		),
		getErrorActions(r, code),
	)

	// Use minimal layout for access-related errors (403, 401) and server errors (500, etc.)
	// to avoid issues with corrupted request context during panics
	if code == http.StatusForbidden || code == http.StatusUnauthorized || code >= 500 {
		return r.Render(layouts.Minimal, content)
	}

	return r.Render(layouts.Primary, content)
}

func getErrorTitle(code int) string {
	switch code {
	case http.StatusNotFound:
		return "Not Found"
	case http.StatusInternalServerError:
		return "Server Error"
	case http.StatusForbidden, http.StatusUnauthorized:
		return "Access Required"
	default:
		return "Error"
	}
}

func getErrorMessage(code int) string {
	switch code {
	case http.StatusNotFound:
		return "The page you're looking for seems to have wandered off. Maybe it's taking a coffee break?"
	case http.StatusInternalServerError:
		return "Our server encountered an error. We're working to fix it."
	case http.StatusForbidden, http.StatusUnauthorized:
		return "This content requires authentication. Please sign in to continue."
	default:
		return "An unexpected error occurred. Please try again."
	}
}

func getErrorActions(r *ui.Request, code int) Node {
	switch code {
	case http.StatusNotFound:
		return Div(
			Class("flex justify-center"),
			A(
				Href("/"),
				Class("px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-md shadow-sm transition-colors"),
				Text("Take Me Home"),
			),
		)
	case http.StatusInternalServerError:
		return Div(
			Class("flex gap-4 justify-center"),
			Button(
				Class("px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-md shadow-sm transition-colors"),
				Type("button"),
				Attr("onclick", "window.location.reload()"),
				Text("Try Again"),
			),
			A(
				Href("/"),
				Class("px-4 py-2 border border-slate-300 hover:border-indigo-500 text-slate-700 hover:text-indigo-600 font-medium rounded-md shadow-sm transition-colors"),
				Text("Go Home"),
			),
		)
	case http.StatusForbidden, http.StatusUnauthorized:
		// Check if this is a server-down scenario (middleware not loaded properly)
		// In such cases, only show "Back to Home"
		if r.Config == nil {
			return Div(
				Class("flex justify-center"),
				A(
					Href("/"),
					Class("btn btn-primary btn-lg gap-2 bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 border-0 text-white font-medium shadow-lg hover:shadow-xl hover:scale-105 transition-all duration-200"),
					Span(Class("text-lg"), Text("🏠")),
					Text("Back to Home"),
				),
			)
		}

		// Normal auth required scenario - show login/register options
		return Div(
			Class("space-y-6"),
			// Primary action buttons
			Div(
				Class("flex flex-col sm:flex-row gap-4 justify-center"),
				A(
					Href(r.Path(routenames.Login)),
					Class("btn btn-primary btn-lg gap-3 bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 border-0 text-white font-medium shadow-lg hover:shadow-xl hover:scale-105 transition-all duration-200 min-w-[150px]"),
					Span(Class("text-xl"), Text("🔑")),
					Text("Sign In"),
				),
				A(
					Href(r.Path(routenames.Register)),
					Class("btn btn-success btn-lg gap-3 bg-gradient-to-r from-green-600 to-emerald-600 hover:from-green-700 hover:to-emerald-700 border-0 text-white font-medium shadow-lg hover:shadow-xl hover:scale-105 transition-all duration-200 min-w-[150px]"),
					Span(Class("text-xl"), Text("✨")),
					Text("Join Now"),
				),
			),
			// Secondary action
			Div(
				Class("flex justify-center"),
				A(
					Href("/"),
					Class("btn btn-ghost btn-md gap-2 text-slate-600 hover:text-indigo-600 transition-colors duration-200"),
					Span(Class("text-lg"), Text("🏠")),
					Text("Back to Home"),
				),
			),
			// Additional context
			Div(
				Class("bg-gradient-to-r from-blue-50 to-purple-50 rounded-xl p-6 border border-blue-100"),
				Div(
					Class("text-center space-y-3"),
					P(
						Class("text-sm font-semibold text-blue-700 mb-2"),
						Text("Why sign up? 🎯"),
					),
					Div(
						Class("grid grid-cols-1 sm:grid-cols-3 gap-4 text-xs text-slate-600"),
						Div(
							Class("flex items-center gap-2"),
							Span(Class("text-lg"), Text("📚")),
							Text("Access exclusive content"),
						),
						Div(
							Class("flex items-center gap-2"),
							Span(Class("text-lg"), Text("🚀")),
							Text("Track your progress"),
						),
						Div(
							Class("flex items-center gap-2"),
							Span(Class("text-lg"), Text("🤝")),
							Text("Join the community"),
						),
					),
				),
			),
		)
	default:
		return A(
			Href("/"),
			Class("btn btn-primary btn-lg gap-2 bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 border-0 text-white font-medium shadow-lg hover:shadow-xl hover:scale-105 transition-all duration-200"),
			Span(Class("text-lg"), Text("🏠")),
			Text("Take Me Home"),
		)
	}
}

func getErrorIllustration(code int) string {
	switch code {
	case http.StatusNotFound:
		return `<svg viewBox="0 0 200 200" fill="none" xmlns="http://www.w3.org/2000/svg">
			<defs>
				<linearGradient id="notFoundGrad" x1="0%" y1="0%" x2="100%" y2="100%">
					<stop offset="0%" style="stop-color:#3B82F6;stop-opacity:0.15" />
					<stop offset="100%" style="stop-color:#8B5CF6;stop-opacity:0.15" />
				</linearGradient>
			</defs>
			<circle cx="100" cy="100" r="90" fill="url(#notFoundGrad)" />
			<rect x="80" y="70" width="40" height="35" rx="8" fill="#64748B"/>
			<circle cx="90" cy="85" r="4" fill="#3B82F6"/>
			<circle cx="110" cy="85" r="4" fill="#3B82F6"/>
			<path d="M85 95Q100 105 115 95" stroke="#64748B" stroke-width="2" fill="none"/>
			<circle cx="100" cy="130" r="15" fill="#F59E0B"/>
			<text x="100" y="138" text-anchor="middle" fill="white" font-size="18" font-weight="bold">?</text>
		</svg>`
	case http.StatusInternalServerError:
		return `<svg viewBox="0 0 200 200" fill="none" xmlns="http://www.w3.org/2000/svg">
			<defs>
				<linearGradient id="errorGrad" x1="0%" y1="0%" x2="100%" y2="100%">
					<stop offset="0%" style="stop-color:#EF4444;stop-opacity:0.15" />
					<stop offset="100%" style="stop-color:#F97316;stop-opacity:0.15" />
				</linearGradient>
			</defs>
			<circle cx="100" cy="100" r="90" fill="url(#errorGrad)" />
			<rect x="80" y="70" width="40" height="35" rx="8" fill="#64748B"/>
			<circle cx="90" cy="85" r="4" fill="#EF4444"/>
			<circle cx="110" cy="85" r="4" fill="#EF4444"/>
			<rect x="95" y="92" width="10" height="3" rx="1.5" fill="#EF4444"/>
			<circle cx="100" cy="125" r="12" fill="#EF4444"/>
			<path d="M95 120l10 10M105 120l-10 10" stroke="white" stroke-width="2"/>
		</svg>`
	case http.StatusForbidden, http.StatusUnauthorized:
		return `<svg viewBox="0 0 200 200" fill="none" xmlns="http://www.w3.org/2000/svg">
			<defs>
				<linearGradient id="authGrad" x1="0%" y1="0%" x2="100%" y2="100%">
					<stop offset="0%" style="stop-color:#3B82F6;stop-opacity:0.15" />
					<stop offset="100%" style="stop-color:#8B5CF6;stop-opacity:0.15" />
				</linearGradient>
				<linearGradient id="lockGrad" x1="0%" y1="0%" x2="100%" y2="100%">
					<stop offset="0%" style="stop-color:#4F46E5" />
					<stop offset="100%" style="stop-color:#7C3AED" />
				</linearGradient>
			</defs>
			<circle cx="100" cy="100" r="90" fill="url(#authGrad)" />
			
			<!-- Modern lock icon -->
			<rect x="75" y="110" width="50" height="35" rx="8" fill="url(#lockGrad)" />
			<rect x="80" y="115" width="40" height="25" rx="4" fill="white" fill-opacity="0.1" />
			<circle cx="100" cy="127" r="4" fill="white" />
			<rect x="98" y="127" width="4" height="8" fill="white" />
			
			<!-- Lock shackle -->
			<path d="M85 110V95C85 86.7 91.7 80 100 80C108.3 80 115 86.7 115 95V110" 
				  stroke="url(#lockGrad)" stroke-width="4" fill="none" stroke-linecap="round"/>
			
			<!-- Decorative elements -->
			<circle cx="130" cy="80" r="3" fill="#4F46E5" opacity="0.6"/>
			<circle cx="70" cy="90" r="2" fill="#7C3AED" opacity="0.4"/>
			<circle cx="140" cy="130" r="2" fill="#4F46E5" opacity="0.5"/>
			
			<!-- Key icon floating -->
			<g transform="translate(60, 60)">
				<circle cx="8" cy="8" r="6" fill="#F59E0B" opacity="0.8"/>
				<rect x="12" y="6" width="15" height="4" rx="2" fill="#F59E0B" opacity="0.8"/>
				<rect x="24" y="4" width="3" height="2" fill="#F59E0B" opacity="0.8"/>
				<rect x="24" y="8" width="5" height="2" fill="#F59E0B" opacity="0.8"/>
			</g>
		</svg>`
	default:
		return `<svg viewBox="0 0 200 200" fill="none" xmlns="http://www.w3.org/2000/svg">
			<defs>
				<linearGradient id="defaultGrad" x1="0%" y1="0%" x2="100%" y2="100%">
					<stop offset="0%" style="stop-color:#64748B;stop-opacity:0.15" />
					<stop offset="100%" style="stop-color:#94A3B8;stop-opacity:0.15" />
				</linearGradient>
			</defs>
			<circle cx="100" cy="100" r="90" fill="url(#defaultGrad)" />
			<rect x="80" y="70" width="40" height="35" rx="8" fill="#64748B"/>
			<circle cx="90" cy="85" r="4" fill="#94A3B8"/>
			<circle cx="110" cy="85" r="4" fill="#94A3B8"/>
			<path d="M85 95Q100 105 115 95" stroke="#64748B" stroke-width="2" fill="none"/>
		</svg>`
	}
}
