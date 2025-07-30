package components

import (
	"github.com/r-scheele/zero/pkg/routenames"
	"github.com/r-scheele/zero/pkg/ui"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// AppHeader renders the application header with logo and authentication buttons
func AppHeader(r *ui.Request) Node {
	return Header(
		Class("fixed top-0 left-0 right-0 z-50 glass-header"),
		Div(
			Class("w-full px-4 sm:px-6 lg:px-8 xl:px-12 py-1 lg:py-2"), // Further reduced padding
			Div(
				Class("flex items-center justify-between h-8 lg:h-10 w-full"), // Further reduced height
				// Logo section
				Div(
					Class("flex items-center gap-3 sm:gap-4 lg:gap-5 flex-shrink-0"),
					// Mobile menu toggle button - only show when authenticated
					If(r.IsAuth,
						Label(
							For("sidebar"),
							Class("btn btn-ghost btn-sm p-2 drawer-button lg:hidden hover:bg-slate-100 border-slate-200 text-slate-600 min-h-10"), // Reduced height
							Div(
								Class("w-5 h-5 flex flex-col justify-center space-y-1"), // Smaller icon
								Div(Class("w-5 h-0.5 bg-current rounded-full")),
								Div(Class("w-5 h-0.5 bg-current rounded-full")),
								Div(Class("w-5 h-0.5 bg-current rounded-full")),
							),
						),
					),
					// Logo
					A(
						Href(func() string {
							if r.IsAuth {
								return "/home"
							}
							return "/"
						}()),
						Class("flex items-center gap-2 sm:gap-3 lg:gap-4 text-base sm:text-lg lg:text-xl font-bold text-slate-800 hover:text-indigo-600 transition-all duration-300 ease-out"), // Further reduced text
						Img(
							Src(ui.StaticFile("logo.png")),
							Alt("Zero Logo"),
							Class("h-6 sm:h-7 lg:h-8 w-auto drop-shadow-sm"), // Further reduced logo
						),
						Span(
							Class("hidden xs:block font-black tracking-tight text-indigo-600"),
							Text("Zero"),
						),
					),
				),
				// Authentication buttons
				Div(
					Class("flex items-center gap-2 sm:gap-3 lg:gap-4 flex-shrink-0"), // Reduced gap
					If(r.IsAuth && r.AuthUser != nil && !r.AuthUser.Verified,
						A(
							Href(r.Path(routenames.VerificationNotice)),
							Class("btn btn-warning btn-sm lg:btn-md gap-2 btn-modern shadow-elegant text-white bg-amber-500 hover:bg-amber-600 border-amber-500 font-medium min-h-10 px-3 lg:px-4"), // Smaller sizing
							Style("min-height: 40px; font-size: 14px;"), // Smaller touch targets
							Span(Class("hidden md:inline text-xs lg:text-sm"), Text("Verify WhatsApp")),
						),
					),
					If(r.IsAuth,
						A(
						Href(r.Path(routenames.Logout)),
						Class("btn btn-ghost btn-md lg:btn-lg gap-2 btn-modern text-slate-600 hover:text-red-600 hover:bg-red-50 border-slate-200 font-medium min-h-12 px-4 lg:px-6"),
						Style("min-height: 48px; font-size: 16px;"), // Better touch targets
						Span(Class("hidden md:inline text-sm lg:text-base"), Text("Logout")),
					),
					),
					If(!r.IsAuth,
						Group{
							A(
							Href(r.Path(routenames.Login)),
							Class("btn btn-ghost btn-md lg:btn-lg gap-2 btn-modern text-slate-600 hover:text-indigo-600 hover:bg-indigo-50 border-slate-200 font-medium min-h-12 px-4 lg:px-6"),
							Style("min-height: 48px; font-size: 16px;"), // Better touch targets
							Span(Class("hidden sm:inline text-sm lg:text-base"), Text("Login")),
						),
							A(
								Href(r.Path(routenames.Register)),
								Class("btn btn-primary btn-md lg:btn-lg gap-2 btn-modern shadow-elegant bg-blue-600 hover:bg-blue-700 border-0 text-white font-medium min-h-12 px-4 lg:px-6"),
								Style("min-height: 48px; font-size: 16px;"), // Better touch targets
								Span(Class("hidden sm:inline text-sm lg:text-base"), Text("Register")),
							),
						},
					),
				),
			),
		),
	)
}
