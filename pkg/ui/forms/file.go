package forms

import (
	"net/http"

	"github.com/r-scheele/zero/pkg/routenames"
	"github.com/r-scheele/zero/pkg/ui"
	. "github.com/r-scheele/zero/pkg/ui/components"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type File struct{}

func (f File) Render(r *ui.Request) Node {
	return Form(
		ID("files"),
		Method(http.MethodPost),
		Action(r.Path(routenames.FilesSubmit)),
		EncType("multipart/form-data"),
		FileField(FileFieldParams{
			Name:  "file",
			Label: "Test file",
			Help:  "Pick a file to upload.",
		}),
		ControlGroup(
			FormButton(ColorPrimary, "Upload"),
		),
		CSRF(r),
	)
}
