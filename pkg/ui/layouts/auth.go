package layouts

import (
	"github.com/r-scheele/zero/pkg/ui"
	. "github.com/r-scheele/zero/pkg/ui/components"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Auth(r *ui.Request, content Node) Node {
	return Doctype(
		HTML(
			Lang("en"),
			Data("theme", "light"),
			Head(
				Metatags(r),
				CSS(),
				JS(),
			),
			Body(
				Class("min-h-screen bg-gradient-to-br from-slate-50 to-gray-100"),
				// Use the same unified header
				unifiedHeader(r),
				Div(
					Class("pt-16 min-h-screen"),
					Main(
						Class("w-full bg-slate-50"),
						Div(
							Class("px-4 sm:px-6 lg:px-8 xl:px-12 py-6 lg:py-8 pb-20 lg:pb-8"),
							Div(
								Class("max-w-md mx-auto"),
								ID("main-content"),
								If(len(r.Title) > 0, H1(
									Class("text-3xl font-bold mb-8 text-slate-900 text-center"),
									Text(r.Title),
								)),
								FlashMessages(r),
								Div(
									Class("bg-white rounded-lg shadow-sm border border-gray-200 p-8"),
									content,
								),
							),
						),
					),
				),
				// Use the same unified footer
				unifiedFooter(r),
				HtmxListeners(r),
			),
		),
	)
}
