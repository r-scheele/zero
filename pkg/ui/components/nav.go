package components

import (
	"fmt"

	"github.com/r-scheele/zero/pkg/pager"
	"github.com/r-scheele/zero/pkg/ui"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func MenuLink(r *ui.Request, icon Node, title, routeName string, routeParams ...any) Node {
	href := r.Path(routeName, routeParams...)
	isActive := href == r.CurrentPath

	return Li(
		Class("mx-0"),
		A(
			Href(href),
			Class("flex items-center gap-3 lg:gap-4 p-3 lg:p-4 rounded-xl transition-all duration-200 font-medium text-slate-700 hover:text-blue-600 hover:bg-blue-50 group"),
			Div(
				Class("flex-shrink-0 w-5 h-5 lg:w-6 lg:h-6 transition-colors duration-200"),
				If(isActive,
					Span(Class("text-blue-600"), icon),
				),
				If(!isActive,
					Span(Class("text-slate-500 group-hover:text-blue-600"), icon),
				),
			),
			Span(
				Class("text-sm lg:text-base transition-colors duration-200"),
				If(isActive,
					Span(Class("text-blue-600 font-semibold"), Text(title)),
				),
				If(!isActive,
					Span(Class("group-hover:text-blue-600"), Text(title)),
				),
			),
			Classes{
				"bg-blue-100 text-blue-700 shadow-sm": isActive,
			},
		),
	)
}

func Pager(page int, path string, hasNext bool, hxTarget string) Node {
	href := func(page int) string {
		return fmt.Sprintf("%s?%s=%d",
			path,
			pager.QueryKey,
			page,
		)
	}

	return Div(
		Class("join"),
		A(
			Class("join-item btn"),
			Text("«"),
			If(page <= 1, Disabled()),
			Href(href(page-1)),
			Iff(len(hxTarget) > 0, func() Node {
				return Group{
					Attr("hx-get", href(page-1)),
					Attr("hx-swap", "outerHTML"),
					Attr("hx-target", hxTarget),
				}
			}),
		),
		Button(
			Class("join-item btn"),
			Textf("Page %d", page),
		),
		A(
			Class("join-item btn"),
			Text("»"),
			If(!hasNext, Disabled()),
			Href(href(page+1)),
			Iff(len(hxTarget) > 0, func() Node {
				return Group{
					Attr("hx-get", href(page+1)),
					Attr("hx-swap", "outerHTML"),
					Attr("hx-target", hxTarget),
				}
			}),
		),
	)
}
