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
	var class string

	switch color {
	case ColorSuccess:
		class = "alert-success"
	case ColorInfo:
		class = "alert-info"
	case ColorWarning:
		class = "alert-warning"
	case ColorError:
		class = "alert-error"
	}

	return Div(
		Role("alert"),
		Class("alert mb-2 "+class),
		Attr("x-data", "{show: true}"),
		Attr("x-show", "show"),
		Attr("x-init", "setTimeout(() => show = false, 5000)"), // Auto-dismiss after 5 seconds
		Attr("x-transition:leave", "transition ease-in duration-300"),
		Attr("x-transition:leave-start", "opacity-100 transform translate-x-0"),
		Attr("x-transition:leave-end", "opacity-0 transform translate-x-full"),
		Span(
			Attr("@click", "show = false"),
			Class("cursor-pointer"),
			icons.XCircle(),
		),
		Span(Text(text)),
	)
}
