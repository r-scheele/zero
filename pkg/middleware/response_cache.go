package middleware

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/config"
	"github.com/r-scheele/zero/pkg/services"
)

// ResponseCacheConfig holds configuration for response caching
type ResponseCacheConfig struct {
	Cache      *services.CacheClient
	Config     *config.Config
	Skipper    func(c echo.Context) bool
	Expiration time.Duration
	KeyFunc    func(c echo.Context) string
}

// ResponseCache returns a middleware that caches HTTP responses
func ResponseCache(config ResponseCacheConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = func(c echo.Context) bool {
			// Skip caching for non-GET requests
			if c.Request().Method != http.MethodGet {
				return true
			}
			// Skip caching for authenticated requests (to avoid caching user-specific data)
			if c.Get("authenticated_user") != nil {
				return true
			}
			// Skip caching for admin routes
			if strings.HasPrefix(c.Request().URL.Path, "/admin") {
				return true
			}
			return false
		}
	}

	if config.KeyFunc == nil {
		config.KeyFunc = func(c echo.Context) string {
			// Create cache key from method, path, and query parameters
			key := fmt.Sprintf("%s:%s", c.Request().Method, c.Request().URL.Path)
			if c.Request().URL.RawQuery != "" {
				key += "?" + c.Request().URL.RawQuery
			}
			// Hash the key to ensure it's a valid cache key
			hash := md5.Sum([]byte(key))
			return fmt.Sprintf("response_cache:%x", hash)
		}
	}

	if config.Expiration == 0 {
		config.Expiration = 5 * time.Minute // Default expiration
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			cacheKey := config.KeyFunc(c)

			// Try to get cached response
			if cachedResponse, err := config.Cache.Get().Key(cacheKey).Fetch(c.Request().Context()); err == nil {
				if response, ok := cachedResponse.(CachedResponse); ok {
					// Set cached headers
					for key, value := range response.Headers {
						c.Response().Header().Set(key, value)
					}
					// Set cache hit header
					c.Response().Header().Set("X-Cache", "HIT")
					// Return cached response
					return c.Blob(response.StatusCode, response.ContentType, response.Body)
				}
			}

			// Create a custom response writer to capture the response
			originalWriter := c.Response().Writer
			buf := &bytes.Buffer{}
			c.Response().Writer = &responseCapture{
				ResponseWriter: originalWriter,
				buffer:         buf,
			}

			// Call the next handler
			err := next(c)
			if err != nil {
				return err
			}

			// Only cache successful responses
			if c.Response().Status >= 200 && c.Response().Status < 300 {
				// Capture response data
				cachedResponse := CachedResponse{
					StatusCode:  c.Response().Status,
					ContentType: c.Response().Header().Get("Content-Type"),
					Body:        buf.Bytes(),
					Headers:     make(map[string]string),
				}

				// Copy important headers
				for _, header := range []string{"Content-Type", "Content-Encoding", "Cache-Control"} {
					if value := c.Response().Header().Get(header); value != "" {
						cachedResponse.Headers[header] = value
					}
				}

				// Cache the response
				config.Cache.Set().
					Key(cacheKey).
					Data(cachedResponse).
					Expiration(config.Expiration).
					Save(c.Request().Context())
			}

			// Set cache miss header
			c.Response().Header().Set("X-Cache", "MISS")

			return nil
		}
	}
}

// CachedResponse represents a cached HTTP response
type CachedResponse struct {
	StatusCode  int               `json:"status_code"`
	ContentType string            `json:"content_type"`
	Body        []byte            `json:"body"`
	Headers     map[string]string `json:"headers"`
}

// responseCapture captures the response for caching
type responseCapture struct {
	http.ResponseWriter
	buffer *bytes.Buffer
}

func (rc *responseCapture) Write(b []byte) (int, error) {
	// Write to both the original writer and our buffer
	n, err := rc.ResponseWriter.Write(b)
	if err == nil {
		rc.buffer.Write(b[:n])
	}
	return n, err
}