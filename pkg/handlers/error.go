package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/pkg/context"
	"github.com/r-scheele/zero/pkg/log"
	"github.com/r-scheele/zero/pkg/ui/pages"
)

type Error struct{}

func (e *Error) Page(err error, ctx echo.Context) {
	// Safety check for nil context
	if ctx == nil {
		return
	}

	if ctx.Response().Committed || context.IsCanceledError(err) {
		return
	}

	// Determine the error status code.
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	// Safely log the error - check if logger context is available
	if logger := log.Ctx(ctx); logger != nil {
		switch {
		case code >= 500:
			logger.Error("ERROR DETAILS", "error", err.Error(), "path", ctx.Request().URL.Path, "method", ctx.Request().Method)
		case code >= 400:
			logger.Warn("WARNING DETAILS", "error", err.Error(), "path", ctx.Request().URL.Path, "method", ctx.Request().Method)
		}
	}

	// Set the status code.
	ctx.Response().WriteHeader(code)

	// Render the error page with additional safety
	if renderErr := pages.Error(ctx, code); renderErr != nil {
		// Fallback to simple HTTP error if error page rendering fails
		if logger := log.Ctx(ctx); logger != nil {
			logger.Error("failed to render error page",
				"error", renderErr,
			)
		}
		// Write a simple fallback response
		ctx.String(code, "Internal Server Error")
	}
}
