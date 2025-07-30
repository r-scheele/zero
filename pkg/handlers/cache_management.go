package handlers

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/pkg/form"
	"github.com/r-scheele/zero/pkg/middleware"
	"github.com/r-scheele/zero/pkg/services"
	"github.com/r-scheele/zero/pkg/ui/forms"
	"github.com/r-scheele/zero/pkg/ui/pages"
)

type CacheManagement struct {
	cache *services.CacheClient
}

func init() {
	Register(new(CacheManagement))
}

func (h *CacheManagement) Init(c *services.Container) error {
	h.cache = c.Cache
	return nil
}

func (h *CacheManagement) Routes(g *echo.Group) {
	// Cache management routes (require authentication, verification, and admin)
	cacheManagement := g.Group("/admin/cache", 
		middleware.RequireAuthentication, 
		middleware.RequireVerification,
		middleware.RequireAdmin)
	
	cacheManagement.GET("", h.Page).Name = "admin.cache"
	cacheManagement.POST("/clear", h.ClearCache).Name = "admin.cache.clear"
	cacheManagement.POST("/clear-pattern", h.ClearCachePattern).Name = "admin.cache.clear_pattern"
}

func (h *CacheManagement) Page(ctx echo.Context) error {
	f := form.Get[forms.Cache](ctx)
	return pages.UpdateCache(ctx, f)
}

func (h *CacheManagement) ClearCache(ctx echo.Context) error {
	// Clear all cache entries
	// Note: This is a simplified implementation. In a production system,
	// you might want to implement a more sophisticated cache clearing mechanism
	
	// For now, we'll just clear specific cache patterns
	patterns := []string{
		"user:",
		"note_likes_count:",
		"note_reposts_count:",
		"response_cache:",
		"pages.",
		"layout.",
		"icon.",
	}
	
	for _, pattern := range patterns {
		h.cache.Flush().Tags(pattern).Execute(ctx.Request().Context())
	}
	
	return ctx.JSON(200, map[string]interface{}{
		"success": true,
		"message": "Cache cleared successfully",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *CacheManagement) ClearCachePattern(ctx echo.Context) error {
	pattern := ctx.FormValue("pattern")
	if pattern == "" {
		return ctx.JSON(400, map[string]interface{}{
			"success": false,
			"message": "Pattern is required",
		})
	}
	
	// Clear cache entries matching the pattern
	h.cache.Flush().Tags(pattern).Execute(ctx.Request().Context())
	
	return ctx.JSON(200, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Cache entries matching pattern '%s' cleared successfully", pattern),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}