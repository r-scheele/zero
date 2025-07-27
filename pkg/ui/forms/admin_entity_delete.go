package forms

import (
	"net/http"

	"github.com/r-scheele/zero/pkg/routenames"
	"github.com/r-scheele/zero/pkg/ui"
	. "github.com/r-scheele/zero/pkg/ui/components"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func AdminEntityDelete(r *ui.Request, entityTypeName string) Node {
	return Form(
		Method(http.MethodPost),
		P(
			Textf("Are you sure you want to delete this %s?", entityTypeName),
		),
		ControlGroup(
			FormButton(ColorError, "Delete"),
			ButtonLink(
				ColorNone,
				r.Path(routenames.AdminEntityList(entityTypeName)),
				"Cancel",
			),
		),
		CSRF(r),
	)
}
