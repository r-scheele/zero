package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/r-scheele/zero/pkg/services"
)

type API struct {
	container *services.Container
}

func init() {
	Register(new(API))
}

// Init initializes the WhatsApp API handler with the service container
func (h *API) Init(c *services.Container) error {
	h.container = c
	return nil
}

// Routes registers all WhatsApp API routes
func (h *API) Routes(g *echo.Group) {

	api_str := "/api/v1/"
	// WhatsApp Bot API endpoints
	apiGroup := g.Group(api_str)

	// Health check
	apiGroup.GET("/health", h.HealthCheck)

	// WhatsApp webhook endpoints (360dialog integration)
	webhookGroup := g.Group(api_str + "whatsapp")
	webhook := &WhatsAppWebhook{Container: h.container}

	webhookGroup.GET("/webhook", webhook.VerifyWebhook)
	webhookGroup.POST("/webhook", webhook.HandleWebhook)
}

// HealthCheck for WhatsApp bot to verify API availability
func (h *API) HealthCheck(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"status":    "healthy",
		"service":   "whatsapp-bot-api",
		"timestamp": ctx.Request().Context().Value("timestamp"),
	})
}
