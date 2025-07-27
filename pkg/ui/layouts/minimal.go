package layouts

import (
	"github.com/r-scheele/zero/pkg/ui"
	. "github.com/r-scheele/zero/pkg/ui/components"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// Minimal provides a clean layout without navigation for error/access pages
func Minimal(r *ui.Request, content Node) Node {
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
				Class("min-h-screen bg-gradient-to-br from-slate-50 to-gray-100 flex items-center justify-center p-4"),
				Main(
					Class("w-full"),
					content,
				),
				HtmxListeners(r),
			),
		),
	)
}
