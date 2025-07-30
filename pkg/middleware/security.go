package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"github.com/r-scheele/zero/config"
	"golang.org/x/time/rate"
)

// CORS returns a CORS middleware configured from the application configuration.
func CORS(cfg *config.Config) echo.MiddlewareFunc {
	if !cfg.Security.CORS.Enabled {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return next
		}
	}

	return echomw.CORSWithConfig(echomw.CORSConfig{
		AllowOrigins: cfg.Security.CORS.AllowedOrigins,
		AllowMethods: cfg.Security.CORS.AllowedMethods,
		AllowHeaders: cfg.Security.CORS.AllowedHeaders,
	})
}

// RateLimit returns a rate limiting middleware configured from the application configuration.
func RateLimit(cfg *config.Config) echo.MiddlewareFunc {
	if !cfg.Security.RateLimit.Enabled {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return next
		}
	}

	// Convert requests per minute to requests per second
	requestsPerSecond := float64(cfg.Security.RateLimit.RequestsPerMinute) / 60.0
	limiter := rate.NewLimiter(rate.Limit(requestsPerSecond), cfg.Security.RateLimit.BurstSize)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !limiter.Allow() {
				return echo.NewHTTPError(http.StatusTooManyRequests, "Rate limit exceeded")
			}
			return next(c)
		}
	}
}

// CSP returns a Content Security Policy middleware configured from the application configuration.
func CSP(cfg *config.Config) echo.MiddlewareFunc {
	if !cfg.Security.CSP.Enabled {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return next
		}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Build CSP header value
			var cspParts []string
			
			if cfg.Security.CSP.Directives.DefaultSrc != "" {
				cspParts = append(cspParts, fmt.Sprintf("default-src %s", cfg.Security.CSP.Directives.DefaultSrc))
			}
			if cfg.Security.CSP.Directives.ScriptSrc != "" {
				cspParts = append(cspParts, fmt.Sprintf("script-src %s", cfg.Security.CSP.Directives.ScriptSrc))
			}
			if cfg.Security.CSP.Directives.StyleSrc != "" {
				cspParts = append(cspParts, fmt.Sprintf("style-src %s", cfg.Security.CSP.Directives.StyleSrc))
			}
			if cfg.Security.CSP.Directives.ImgSrc != "" {
				cspParts = append(cspParts, fmt.Sprintf("img-src %s", cfg.Security.CSP.Directives.ImgSrc))
			}

			if len(cspParts) > 0 {
				c.Response().Header().Set("Content-Security-Policy", strings.Join(cspParts, "; "))
			}

			return next(c)
		}
	}
}

// RequestLogging returns a request logging middleware configured from the application configuration.
func RequestLogging(cfg *config.Config) echo.MiddlewareFunc {
	if !cfg.Monitoring.RequestLogging.Enabled {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return next
		}
	}

	return echomw.LoggerWithConfig(echomw.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			path := c.Request().URL.Path
			for _, excludePath := range cfg.Monitoring.RequestLogging.ExcludePaths {
				if strings.HasPrefix(path, excludePath) {
					return true
				}
			}
			return false
		},
		Format: "method=${method}, uri=${uri}, status=${status}, latency=${latency_human}\n",
	})
}

// HealthCheck returns a health check middleware that responds to health check requests.
func HealthCheck(cfg *config.Config) echo.MiddlewareFunc {
	if !cfg.Monitoring.Health.Enabled {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return next
		}
	}

	healthEndpoint := cfg.Monitoring.Health.Endpoint
	if healthEndpoint == "" {
		healthEndpoint = "/health"
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().URL.Path == healthEndpoint {
				return c.JSON(http.StatusOK, map[string]interface{}{
					"status": "healthy",
					"timestamp": time.Now().UTC().Format(time.RFC3339),
				})
			}
			return next(c)
		}
	}
}

// Metrics returns a metrics middleware that responds to metrics requests.
func Metrics(cfg *config.Config) echo.MiddlewareFunc {
	if !cfg.Monitoring.Metrics.Enabled {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return next
		}
	}

	metricsEndpoint := cfg.Monitoring.Metrics.Endpoint
	if metricsEndpoint == "" {
		metricsEndpoint = "/metrics"
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().URL.Path == metricsEndpoint {
				// Basic metrics response - in a real application, you'd integrate with
				// a metrics library like Prometheus
				return c.String(http.StatusOK, "# Basic metrics placeholder\n# Integrate with Prometheus or similar for production use\n")
			}
			return next(c)
		}
	}
}