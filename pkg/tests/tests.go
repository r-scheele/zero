package tests

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/r-scheele/zero/ent"
	"github.com/r-scheele/zero/pkg/session"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

// NewContext creates a new Echo context for tests using an HTTP test request and response recorder
func NewContext(e *echo.Echo, url string) (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(http.MethodGet, url, strings.NewReader(""))
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

// InitSession initializes a session for a given Echo context
func InitSession(ctx echo.Context) {
	session.Store(ctx, sessions.NewCookieStore([]byte("secret")))
}

// ExecuteMiddleware executes a middleware function on a given Echo context
func ExecuteMiddleware(ctx echo.Context, mw echo.MiddlewareFunc) error {
	handler := mw(func(c echo.Context) error {
		return nil
	})
	return handler(ctx)
}

// ExecuteHandler executes a handler with an optional stack of middleware
func ExecuteHandler(ctx echo.Context, handler echo.HandlerFunc, mw ...echo.MiddlewareFunc) error {
	return ExecuteMiddleware(ctx, func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			run := handler

			for _, w := range mw {
				run = w(run)
			}

			return run(ctx)
		}
	})
}

// AssertHTTPErrorCode asserts an HTTP status code on a given Echo HTTP error
func AssertHTTPErrorCode(t *testing.T, err error, code int) {
	httpError, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	assert.Equal(t, code, httpError.Code)
}

// CreateUser creates a random user entity
func CreateUser(orm *ent.Client) (*ent.User, error) {
	seed := rand.Intn(9999999) // Generate 7-digit number
	return orm.User.
		Create().
		SetPhoneNumber(fmt.Sprintf("+1555%07d", seed)).
		SetPassword("password").
		SetName(fmt.Sprintf("Test User %d", seed)).
		Save(context.Background())
}
