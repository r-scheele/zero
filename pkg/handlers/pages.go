package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/pkg/routenames"
	"github.com/r-scheele/zero/pkg/services"
	"github.com/r-scheele/zero/pkg/ui/pages"
)

type Pages struct{}

func init() {
	Register(new(Pages))
}

func (h *Pages) Init(c *services.Container) error {
	return nil
}

func (h *Pages) Routes(g *echo.Group) {
	g.GET("/", h.Home).Name = routenames.Home
	g.GET("/dashboard", h.Dashboard).Name = routenames.Dashboard
	g.GET("/about", h.About).Name = routenames.About
}

func (h *Pages) Home(ctx echo.Context) error {
	return pages.Home(ctx, nil)
}

func (h *Pages) Dashboard(ctx echo.Context) error {
	return pages.Dashboard(ctx, nil)
}

func (h *Pages) About(ctx echo.Context) error {
	return pages.About(ctx)
}
