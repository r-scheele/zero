package handlers

import (
	"net/http"
	
	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/ent"
	"github.com/r-scheele/zero/pkg/context"
	"github.com/r-scheele/zero/pkg/routenames"
	"github.com/r-scheele/zero/pkg/services"
	"github.com/r-scheele/zero/pkg/ui/pages"
)

type Pages struct{
	orm *ent.Client
}

func init() {
	Register(new(Pages))
}

func (h *Pages) Init(c *services.Container) error {
	h.orm = c.ORM
	return nil
}

func (h *Pages) Routes(g *echo.Group) {
	// Root path for non-authenticated users (public landing page)
	g.GET("/", h.PublicHome).Name = routenames.Home
	// Authenticated user home page
	g.GET("/home", h.AuthenticatedHome).Name = "authenticated_home"
	g.GET("/dashboard", h.Dashboard).Name = routenames.Dashboard
	g.GET("/about", h.About).Name = routenames.About
	// Quizzes page (redirect to quiz list)
	g.GET("/quizzes", h.QuizzesPage).Name = "quizzes"
}

// PublicHome serves the public landing page for non-authenticated users
func (h *Pages) PublicHome(ctx echo.Context) error {
	// Check if user is already authenticated
	currentUser := ctx.Get(context.AuthenticatedUserKey)
	if currentUser != nil {
		// Redirect authenticated users to their home page
		return ctx.Redirect(http.StatusFound, "/home")
	}
	
	// Serve public landing page
	return pages.Home(ctx, nil)
}

// AuthenticatedHome serves the authenticated user's home page
func (h *Pages) AuthenticatedHome(ctx echo.Context) error {
	// Safety check: ensure current user is not nil
	currentUser := ctx.Get(context.AuthenticatedUserKey)
	if currentUser == nil {
		// Redirect unauthenticated users to login page
		return ctx.Redirect(http.StatusFound, "/user/login")
	}
	
	return pages.Home(ctx, nil)
}

func (h *Pages) Dashboard(ctx echo.Context) error {
	// Safety check: ensure current user is not nil
	currentUser := ctx.Get(context.AuthenticatedUserKey)
	if currentUser == nil {
		// Redirect unauthenticated users to login page
		return ctx.Redirect(http.StatusFound, "/user/login")
	}
	
	return pages.Dashboard(ctx, nil, h.orm)
}

func (h *Pages) About(ctx echo.Context) error {
	return pages.About(ctx)
}

func (h *Pages) QuizzesPage(ctx echo.Context) error {
	return pages.Quizzes(ctx)
}
