package handlers

import (
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"github.com/r-scheele/zero/pkg/context"
	mw "github.com/r-scheele/zero/pkg/middleware"
	"github.com/r-scheele/zero/pkg/services"
	files "github.com/r-scheele/zero/public"
)

// BuildRouter builds the router.
func BuildRouter(c *services.Container) error {
	// Force HTTPS, if enabled.
	if c.Config.HTTP.TLS.Enabled {
		c.Web.Use(echomw.HTTPSRedirect())
	}

	// Serve public files with cache control.
	c.Web.Group("", mw.CacheControl(c.Config.Cache.Expiration.PublicFile)).
		Static("files", "public/files")

	// Serve static files.
	// ui.StaticFile() should be used in ui components to append a cache key to the URL to break cache
	// after each server reboot.
	c.Web.Group(
		"",
		echomw.GzipWithConfig(echomw.GzipConfig{
			Skipper: func(c echo.Context) bool {
				for _, ext := range []string{
					".js",
					".css",
				} {
					if strings.HasSuffix(c.Request().URL.Path, ext) {
						return false
					}
				}
				return true
			},
		}),
		mw.CacheControl(c.Config.Cache.Expiration.PublicFile),
	).StaticFS("static", echo.MustSubFS(files.Static, "static"))

	// Non-static file route group.
	g := c.Web.Group("")

	// Create a cookie store for session data.
	cookieStore := sessions.NewCookieStore([]byte(c.Config.App.EncryptionKey))
	cookieStore.Options.HttpOnly = true
	// Use SameSiteLaxMode instead of Strict for Safari compatibility
	// Safari has issues with SameSiteStrictMode in certain authentication flows
	cookieStore.Options.SameSite = http.SameSiteLaxMode
	// Set Secure flag for HTTPS environments
	cookieStore.Options.Secure = c.Config.HTTP.TLS.Enabled

	g.Use(
		echomw.RemoveTrailingSlashWithConfig(echomw.TrailingSlashConfig{
			RedirectCode: http.StatusMovedPermanently,
		}),
		echomw.RecoverWithConfig(echomw.RecoverConfig{
			DisableErrorHandler: false,
			DisablePrintStack:   false,
		}),
		mw.CORS(c.Config),
		mw.RateLimit(c.Config),
		mw.CSP(c.Config),
		mw.RequestLogging(c.Config),
		mw.HealthCheck(c.Config),
		mw.Metrics(c.Config),
		echomw.RequestID(),
		mw.SetLogger(),
		mw.LogRequest(),
		echomw.Gzip(),
		// Temporarily removed timeout middleware due to Go stdlib panic
		// echomw.TimeoutWithConfig(echomw.TimeoutConfig{
		//	Timeout: c.Config.App.Timeout,
		//	Skipper: func(ctx echo.Context) bool {
		//		// Skip timeout for health checks and static files to prevent issues
		//		path := ctx.Request().URL.Path
		//		return path == "/health" ||
		//			path == "/ping" ||
		//			strings.HasPrefix(path, "/static/") ||
		//			strings.HasPrefix(path, "/files/")
		//	},
		// }),
		mw.Config(c.Config),
		mw.Session(cookieStore),
		mw.LoadAuthenticatedUser(c.Auth),
		mw.ResponseCache(mw.ResponseCacheConfig{
			Cache:      c.Cache,
			Config:     c.Config,
			Expiration: c.Config.Cache.Expiration.PublicNotes,
		}),
		echomw.CSRFWithConfig(echomw.CSRFConfig{
			TokenLookup:    "form:csrf",
			CookieName:     "csrf",
			CookieHTTPOnly: true,
			CookieSameSite: func() http.SameSite {
				if c.Config.Security.CSP.Enabled {
					return http.SameSiteStrictMode
				}
				return http.SameSiteLaxMode
			}(),
			ContextKey:     context.CSRFKey,
		}),
	)

	// Error handler.
	c.Web.HTTPErrorHandler = new(Error).Page

	// Initialize and register all handlers.
	for _, h := range GetHandlers() {
		if err := h.Init(c); err != nil {
			return err
		}

		h.Routes(g)
	}

	return nil
}
