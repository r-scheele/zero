package pages

import (
	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/pkg/ui"
	"github.com/r-scheele/zero/pkg/ui/forms"
	"github.com/r-scheele/zero/pkg/ui/layouts"
)

func UpdateCache(ctx echo.Context, form *forms.Cache) error {
	r := ui.NewRequest(ctx)
	r.Title = "Set a cache entry"

	// Use admin layout for admin users, regular layout for others
	if r.IsAdmin {
		return r.Render(layouts.Admin, form.Render(r))
	}
	return r.Render(layouts.Primary, form.Render(r))
}
