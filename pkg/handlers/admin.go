package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"entgo.io/ent/entc/gen"
	"entgo.io/ent/entc/load"
	"github.com/labstack/echo/v4"
	"github.com/mikestefanello/backlite/ui"
	"github.com/r-scheele/zero/ent"
	"github.com/r-scheele/zero/ent/admin"
	"github.com/r-scheele/zero/pkg/context"
	"github.com/r-scheele/zero/pkg/htmx"
	"github.com/r-scheele/zero/pkg/middleware"
	"github.com/r-scheele/zero/pkg/msg"
	"github.com/r-scheele/zero/pkg/pager"
	"github.com/r-scheele/zero/pkg/redirect"
	"github.com/r-scheele/zero/pkg/routenames"
	"github.com/r-scheele/zero/pkg/services"
	"github.com/r-scheele/zero/pkg/ui/pages"
)

type Admin struct {
	orm      *ent.Client
	graph    *gen.Graph
	admin    *admin.Handler
	backlite *ui.Handler
}

func init() {
	Register(new(Admin))
}

func (h *Admin) Init(c *services.Container) error {
	var err error
	h.graph = c.Graph
	h.orm = c.ORM
	h.admin = admin.NewHandler(h.orm, admin.HandlerConfig{
		ItemsPerPage: 25,
		PageQueryKey: pager.QueryKey,
		TimeFormat:   time.DateTime,
	})
	h.backlite, err = ui.NewHandler(ui.Config{
		DB:           c.Database,
		BasePath:     "/admin/tasks",
		ItemsPerPage: 25,
		ReleaseAfter: c.Config.Tasks.ReleaseAfter,
	})
	return err
}

func (h *Admin) Routes(g *echo.Group) {
	ag := g.Group("/admin", middleware.RequireAdmin)

	// Admin overview/dashboard
	ag.GET("", h.Overview()).Name = "admin:overview"

	entities := ag.Group("/entity")
	for _, n := range h.graph.Nodes {
		// Skip PasswordToken entity for security reasons
		if n.Name == "PasswordToken" {
			continue
		}

		ng := entities.Group(fmt.Sprintf("/%s", strings.ToLower(n.Name)))
		ng.GET("", h.EntityList(n)).
			Name = routenames.AdminEntityList(n.Name)
		ng.GET("/add", h.EntityAdd(n)).
			Name = routenames.AdminEntityAdd(n.Name)
		ng.POST("/add", h.EntityAddSubmit(n)).
			Name = routenames.AdminEntityAddSubmit(n.Name)
		ng.GET("/:id", h.EntityView(n), h.middlewareEntityLoad(n)).
			Name = routenames.AdminEntityView(n.Name)
		ng.GET("/:id/edit", h.EntityEdit(n), h.middlewareEntityLoad(n)).
			Name = routenames.AdminEntityEdit(n.Name)
		ng.POST("/:id/edit", h.EntityEditSubmit(n), h.middlewareEntityLoad(n)).
			Name = routenames.AdminEntityEditSubmit(n.Name)
		ng.GET("/:id/delete", h.EntityDelete(n), h.middlewareEntityLoad(n)).
			Name = routenames.AdminEntityDelete(n.Name)
		ng.POST("/:id/delete", h.EntityDeleteSubmit(n), h.middlewareEntityLoad(n)).
			Name = routenames.AdminEntityDeleteSubmit(n.Name)
	}

	tasks := ag.Group("/tasks")
	tasks.GET("", h.AdminTasks).Name = routenames.AdminTasks
	tasks.GET("/succeeded", h.Backlite(h.backlite.Succeeded))
	tasks.GET("/failed", h.Backlite(h.backlite.Failed))
	tasks.GET("/upcoming", h.Backlite(h.backlite.Upcoming))
	tasks.GET("/task/:id", h.Backlite(h.backlite.Task))
	tasks.GET("/completed/:id", h.Backlite(h.backlite.TaskCompleted))
}

// Overview displays the admin dashboard
func (h *Admin) Overview() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return pages.AdminOverview(ctx, h.orm)
	}
}

// AdminTasks displays the admin tasks page with proper layout
func (h *Admin) AdminTasks(ctx echo.Context) error {
	return pages.AdminTasks(ctx)
}

// middlewareEntityLoad is middleware to extract the entity ID and attempt to load the given entity.
func (h *Admin) middlewareEntityLoad(n *gen.Type) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			id, err := strconv.Atoi(ctx.Param("id"))
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid entity ID")
			}

			entity, err := h.admin.Get(ctx, n.Name, id)
			switch {
			case err == nil:
				ctx.Set(context.AdminEntityIDKey, id)
				ctx.Set(context.AdminEntityKey, map[string][]string(entity))
				return next(ctx)
			case ent.IsNotFound(err):
				return echo.NewHTTPError(http.StatusNotFound, "entity not found")
			default:
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
		}
	}
}

func (h *Admin) EntityList(n *gen.Type) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		list, err := h.admin.List(ctx, n.Name)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		// Check if this is a targeted HTMX request for the search functionality
		// Only return partial content for HTMX requests that specifically target the table container
		htmxReq := htmx.GetRequest(ctx)
		if htmxReq.Enabled && !htmxReq.Boosted && htmxReq.Target == "user-table-container" && n.Name == "User" {
			// For HTMX search requests on User entity, return only the table content
			return pages.AdminEntityListTable(ctx, n.Name, list)
		}

		// For all other requests (including boosted navigation), return the full page
		return pages.AdminEntityList(ctx, n.Name, list)
	}
}

func (h *Admin) EntityView(n *gen.Type) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		entity := ctx.Get(context.AdminEntityKey).(map[string][]string)
		id := ctx.Get(context.AdminEntityIDKey).(int)
		return pages.AdminEntityView(ctx, n.Name, entity, id)
	}
}

func (h *Admin) EntityAdd(n *gen.Type) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return pages.AdminEntityInput(ctx, h.getEntitySchema(n), nil)
	}
}

func (h *Admin) EntityAddSubmit(n *gen.Type) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		err := h.admin.Create(ctx, n.Name)
		if err != nil {
			msg.Error(ctx, err.Error())
			return h.EntityAdd(n)(ctx)
		}

		msg.Success(ctx, fmt.Sprintf("Successfully added %s.", n.Name))

		return redirect.
			New(ctx).
			Route(routenames.AdminEntityList(n.Name)).
			StatusCode(http.StatusFound).
			Go()
	}
}

func (h *Admin) EntityEdit(n *gen.Type) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		v := ctx.Get(context.AdminEntityKey).(map[string][]string)
		return pages.AdminEntityInput(ctx, h.getEntitySchema(n), v)
	}
}

func (h *Admin) EntityEditSubmit(n *gen.Type) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id := ctx.Get(context.AdminEntityIDKey).(int)
		err := h.admin.Update(ctx, n.Name, id)
		if err != nil {
			msg.Error(ctx, err.Error())
			return h.EntityEdit(n)(ctx)
		}

		msg.Success(ctx, fmt.Sprintf("Updated %s.", n.Name))

		return redirect.
			New(ctx).
			Route(routenames.AdminEntityList(n.Name)).
			StatusCode(http.StatusFound).
			Go()
	}
}

func (h *Admin) EntityDelete(n *gen.Type) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return pages.AdminEntityDelete(ctx, n.Name)
	}
}

func (h *Admin) EntityDeleteSubmit(n *gen.Type) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id := ctx.Get(context.AdminEntityIDKey).(int)
		if err := h.admin.Delete(ctx, n.Name, id); err != nil {
			msg.Error(ctx, err.Error())
			return h.EntityDelete(n)(ctx)
		}

		msg.Success(ctx, fmt.Sprintf("Successfully deleted %s (ID %d).", n.Name, id))

		return redirect.
			New(ctx).
			Route(routenames.AdminEntityList(n.Name)).
			StatusCode(http.StatusFound).
			Go()
	}
}

func (h *Admin) getEntitySchema(n *gen.Type) *load.Schema {
	for _, s := range h.graph.Schemas {
		if s.Name == n.Name {
			return s
		}
	}
	return nil
}

func (h *Admin) Backlite(handler func(http.ResponseWriter, *http.Request) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		if id := c.Param("id"); id != "" {
			c.Request().SetPathValue("task", id)
		}
		return handler(c.Response().Writer, c.Request())
	}
}
