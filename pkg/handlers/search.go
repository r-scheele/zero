package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/pkg/routenames"
	"github.com/r-scheele/zero/pkg/services"
	"github.com/r-scheele/zero/pkg/ui/models"
	"github.com/r-scheele/zero/pkg/ui/pages"
)

type Search struct{}

func init() {
	Register(new(Search))
}

func (h *Search) Init(c *services.Container) error {
	return nil
}

func (h *Search) Routes(g *echo.Group) {
	g.GET("/search", h.Page).Name = routenames.Search
}

func (h *Search) Page(ctx echo.Context) error {
	// Return empty search results for now
	results := make([]*models.SearchResult, 0)
	return pages.SearchResults(ctx, results)
}
