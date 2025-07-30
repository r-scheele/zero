package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/r-scheele/zero/ent"
	"github.com/r-scheele/zero/pkg/context"
	"github.com/r-scheele/zero/pkg/log"
	"github.com/r-scheele/zero/pkg/msg"
	"github.com/r-scheele/zero/pkg/routenames"
	"github.com/r-scheele/zero/pkg/services"
	"github.com/r-scheele/zero/pkg/session"

	"github.com/labstack/echo/v4"
)

// LoadAuthenticatedUser loads the authenticated user, if one, and stores in context.
func LoadAuthenticatedUser(authClient *services.AuthClient) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Always start with a clean context for authentication
			c.Set(context.AuthenticatedUserKey, nil)
			
			u, err := authClient.GetAuthenticatedUser(c)
			switch err.(type) {
			case *ent.NotFoundError:
				log.Ctx(c).Warn("auth user not found - clearing session and cache", 
					"path", c.Request().URL.Path)
				// Clear corrupted session data completely
				if sess, err := session.Get(c, "ua"); err == nil {
					// Get user ID before clearing session for cache cleanup
					if userID, ok := sess.Values["user_id"].(int); ok {
						cacheKey := fmt.Sprintf("user:%d", userID)
						authClient.ClearUserCache(c.Request().Context(), cacheKey)
					}
					// Clear session completely
					for key := range sess.Values {
						delete(sess.Values, key)
					}
					sess.Save(c.Request(), c.Response())
				}
			case services.NotAuthenticatedError:
				// User is not authenticated - context remains nil
				log.Ctx(c).Debug("user not authenticated", 
					"path", c.Request().URL.Path)
			case nil:
				// User is authenticated - set in context
				log.Ctx(c).Debug("authenticated user loaded", 
					"user_id", u.ID,
					"path", c.Request().URL.Path)
				c.Set(context.AuthenticatedUserKey, u)
			default:
				log.Ctx(c).Error("error querying for authenticated user", 
					"error", err,
					"path", c.Request().URL.Path)
				// Clear potentially corrupted session and cache completely on any other error
				if sess, err := session.Get(c, "ua"); err == nil {
					// Get user ID before clearing session for cache cleanup
					if userID, ok := sess.Values["user_id"].(int); ok {
						cacheKey := fmt.Sprintf("user:%d", userID)
						authClient.ClearUserCache(c.Request().Context(), cacheKey)
					}
					// Clear session completely
					for key := range sess.Values {
						delete(sess.Values, key)
					}
					sess.Save(c.Request(), c.Response())
				}
				return echo.NewHTTPError(
					http.StatusInternalServerError,
					fmt.Sprintf("error querying for authenticated user: %v", err),
				)
			}

			return next(c)
		}
	}
}

// LoadValidPasswordToken loads a valid password token entity that matches the user and token
// provided in path parameters
// If the token is invalid, the user will be redirected to the forgot password route
// This requires that the user owning the token is loaded in to context.
func LoadValidPasswordToken(authClient *services.AuthClient) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract the user parameter
			if c.Get(context.UserKey) == nil {
				return echo.NewHTTPError(http.StatusInternalServerError)
			}
			usr := c.Get(context.UserKey).(*ent.User)

			// Extract the token ID.
			tokenID, err := strconv.Atoi(c.Param("password_token"))
			if err != nil {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			// Attempt to load a valid password token.
			token, err := authClient.GetValidPasswordToken(
				c,
				usr.ID,
				tokenID,
				c.Param("token"),
			)

			switch err.(type) {
			case nil:
				c.Set(context.PasswordTokenKey, token)
				return next(c)
			case services.InvalidPasswordTokenError:
				msg.Warning(c, "The link is either invalid or has expired. Please request a new one.")
				return c.Redirect(http.StatusFound, c.Echo().Reverse(routenames.ForgotPassword))
			default:
				return echo.NewHTTPError(
					http.StatusInternalServerError,
					fmt.Sprintf("error loading password token: %v", err),
				)
			}
		}
	}
}

// RequireAuthentication requires that the user be authenticated in order to proceed.
func RequireAuthentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if u := c.Get(context.AuthenticatedUserKey); u == nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		return next(c)
	}
}

// RequireNoAuthentication requires that the user not be authenticated in order to proceed.
func RequireNoAuthentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Allow login form submissions even if user is already authenticated
		// This enables users to login with different credentials
		if c.Request().Method == "POST" && (c.Request().URL.Path == "/user/login" || c.Request().URL.Path == "/user/register") {
			return next(c)
		}
		
		if u := c.Get(context.AuthenticatedUserKey); u != nil {
			// Return forbidden error if user is authenticated
			return echo.NewHTTPError(http.StatusForbidden)
		}

		return next(c)
	}
}

// RequireVerification requires that the authenticated user be verified in order to proceed.
func RequireVerification(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if u := c.Get(context.AuthenticatedUserKey); u != nil {
			if user, ok := u.(*ent.User); ok {
				if user.Verified {
					return next(c)
				}
				// User is authenticated but not verified - redirect to verification notice
				return c.Redirect(http.StatusFound, c.Echo().Reverse(routenames.VerificationNotice))
			}
		}

		// User is not authenticated - require authentication first
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
}

// RequireAdmin requires that the authenticated user be an admin in order to proceed.
func RequireAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if u := c.Get(context.AuthenticatedUserKey); u != nil {
			if user, ok := u.(*ent.User); ok {
				if !user.Verified {
					return c.Redirect(http.StatusFound, c.Echo().Reverse(routenames.VerificationNotice))
				}
				if user.Admin {
					return next(c)
				}
			}
		}

		return echo.NewHTTPError(http.StatusUnauthorized)
	}
}
