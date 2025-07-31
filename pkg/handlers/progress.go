package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/pkg/middleware"
	"github.com/r-scheele/zero/pkg/services"
	"github.com/r-scheele/zero/pkg/ui/pages"
)

type Progress struct {
	container *services.Container
}

func init() {
	Register(new(Progress))
}

func (h *Progress) Init(c *services.Container) error {
	h.container = c
	return nil
}

func (h *Progress) Routes(g *echo.Group) {
	// Progress routes (require authentication and verification)
	progress := g.Group("/progress", middleware.RequireAuthentication, middleware.RequireVerification)
	
	// View progress dashboard
	progress.GET("", h.ViewProgress).Name = "progress.view"
	
	// View detailed analytics
	progress.GET("/analytics", h.ViewAnalytics).Name = "progress.analytics"
	
	// Export progress report
	progress.GET("/export", h.ExportProgress).Name = "progress.export"
}

// ViewProgress displays the progress dashboard
func (h *Progress) ViewProgress(ctx echo.Context) error {
	return pages.ViewProgress(ctx)
}

// ViewAnalytics displays detailed analytics
func (h *Progress) ViewAnalytics(ctx echo.Context) error {
	return pages.ViewAnalytics(ctx)
}

// ExportProgress handles progress report export
func (h *Progress) ExportProgress(ctx echo.Context) error {
	// TODO: Implement progress export logic
	return ctx.String(200, "Progress export functionality coming soon!")
}