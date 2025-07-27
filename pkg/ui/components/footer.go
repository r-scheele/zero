package components

import (
	"github.com/r-scheele/zero/pkg/routenames"
	"github.com/r-scheele/zero/pkg/ui"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// AppFooter renders the application footer with consistent navigation links
// AppFooter renders the application footer with navigation links
func AppFooter(r *ui.Request) Node {
	return Footer(
		Class("fixed bottom-0 left-0 right-0 z-40 glass-footer"),
		Div(
			Class("container mx-auto px-4 sm:px-6 lg:px-8 py-1 lg:py-2 max-w-7xl"), // Further reduced padding
			Div(
				Class("flex flex-col sm:flex-row items-center justify-center sm:justify-between gap-1 sm:gap-2"), // Further reduced gap
				// Copyright section
				Div(
					Class("flex items-center justify-center sm:justify-start order-2 sm:order-1"),
					P(
						Class("text-xs text-slate-600 text-center sm:text-left font-medium"), // Smaller text
						Text("© 2025 Zero. All rights reserved."),
					),
				),
				// Navigation links
				Div(
					Class("flex items-center gap-4 sm:gap-6 order-1 sm:order-2"), // Reduced gap
					A(
						Href(r.Path(routenames.About)),
						Class("text-slate-600 hover:text-blue-600 transition-colors duration-200 font-medium text-xs"), // Smaller text
						Text("About"),
					),
					A(
						Href(r.Path(routenames.Contact)),
						Class("text-slate-600 hover:text-blue-600 transition-colors duration-200 font-medium text-xs"), // Smaller text
						Text("Contact"),
					),
				),
			),
		),
	)
}
