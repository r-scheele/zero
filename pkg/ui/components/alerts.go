package components

import (
	"github.com/r-scheele/zero/pkg/msg"
	"github.com/r-scheele/zero/pkg/ui"
	"github.com/r-scheele/zero/pkg/ui/icons"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func FlashMessages(r *ui.Request) Node {
	var g Group
	var color Color

	for _, typ := range []msg.Type{
		msg.TypeSuccess,
		msg.TypeInfo,
		msg.TypeWarning,
		msg.TypeError,
	} {
		for _, str := range msg.Get(r.Context, typ) {
			switch typ {
			case msg.TypeSuccess:
				color = ColorSuccess
			case msg.TypeInfo:
				color = ColorInfo
			case msg.TypeWarning:
				color = ColorWarning
			case msg.TypeError:
				color = ColorError
			}

			g = append(g, Alert(color, str))
		}
	}

	return g
}

func Alert(color Color, text string) Node {
	var bgClass, borderClass, textClass string

	switch color {
	case ColorSuccess:
		bgClass = "bg-green-50"
		borderClass = "border-green-400"
		textClass = "text-green-800"
	case ColorInfo:
		bgClass = "bg-blue-50"
		borderClass = "border-blue-400"
		textClass = "text-blue-800"
	case ColorWarning:
		bgClass = "bg-yellow-50"
		borderClass = "border-yellow-400"
		textClass = "text-yellow-800"
	case ColorError:
		bgClass = "bg-red-50"
		borderClass = "border-red-400"
		textClass = "text-red-800"
	}

	return Div(
		Role("alert"),
		Class("mb-4 px-4 py-3 rounded-lg border-l-4 shadow-sm flex items-center justify-between " + bgClass + " " + borderClass),
		Attr("x-data", "{show: true}"),
		Attr("x-show", "show"),
		Attr("x-init", "setTimeout(() => show = false, 5000)"), // Auto-dismiss after 5 seconds
		Attr("x-transition:leave", "transition ease-in duration-300"),
		Attr("x-transition:leave-start", "opacity-100 transform translate-y-0"),
		Attr("x-transition:leave-end", "opacity-0 transform -translate-y-2"),
		Span(
			Class("font-medium " + textClass),
			Text(text),
		),
		Span(
			Attr("@click", "show = false"),
			Class("cursor-pointer " + textClass),
			icons.XCircle(),
		),
	)
}
