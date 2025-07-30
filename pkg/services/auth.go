package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/r-scheele/zero/config"
	"github.com/r-scheele/zero/ent"
	"github.com/r-scheele/zero/ent/passwordtoken"
	"github.com/r-scheele/zero/ent/user"
	"github.com/r-scheele/zero/pkg/context"
	"github.com/r-scheele/zero/pkg/session"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	// authSessionName stores the name of the session which contains authentication data
	authSessionName = "ua"

	// authSessionKeyUserID stores the key used to store the user ID in the session
	authSessionKeyUserID = "user_id"

	// authSessionKeyAuthenticated stores the key used to store the authentication status in the session
	authSessionKeyAuthenticated = "authenticated"
)

// NotAuthenticatedError is an error returned when a user is not authenticated
type NotAuthenticatedError struct{}

// Error implements the error interface.
func (e NotAuthenticatedError) Error() string {
	return "user not authenticated"
}

// InvalidPasswordTokenError is an error returned when an invalid token is provided
type InvalidPasswordTokenError struct{}

// Error implements the error interface.
func (e InvalidPasswordTokenError) Error() string {
	return "invalid password token"
}

// AuthClient is the client that handles authentication requests
type AuthClient struct {
	config *config.Config
	orm    *ent.Client
	cache  *CacheClient
}

// NewAuthClient creates a new authentication client
func NewAuthClient(cfg *config.Config, orm *ent.Client, cache *CacheClient) *AuthClient {
	return &AuthClient{
		config: cfg,
		orm:    orm,
		cache:  cache,
	}
}

// Login logs in a user of a given ID
func (c *AuthClient) Login(ctx echo.Context, userID int) error {
	// Clear any existing user cache before login
	if existingUserID, _ := c.GetAuthenticatedUserID(ctx); existingUserID > 0 {
		cacheKey := fmt.Sprintf("user:%d", existingUserID)
		c.cache.Flush().Key(cacheKey).Execute(ctx.Request().Context())
	}

	sess, err := session.Get(ctx, authSessionName)
	if err != nil {
		return err
	}

	// Clear existing session values first
	delete(sess.Values, authSessionKeyUserID)
	delete(sess.Values, authSessionKeyAuthenticated)

	// Set new session values
	sess.Values[authSessionKeyUserID] = userID
	sess.Values[authSessionKeyAuthenticated] = true

	// Clear cache for the new user to ensure fresh data
	cacheKey := fmt.Sprintf("user:%d", userID)
	c.cache.Flush().Key(cacheKey).Execute(ctx.Request().Context())

	return sess.Save(ctx.Request(), ctx.Response())
}

// Logout logs the requesting user out
func (c *AuthClient) Logout(ctx echo.Context) error {
	// Fast logout - get user ID for cache clearing
	userID, _ := c.GetAuthenticatedUserID(ctx)

	// Get and clear session quickly
	sess, err := session.Get(ctx, authSessionName)
	if err != nil {
		return nil // Don't fail if session doesn't exist
	}

	// Clear session data
	sess.Values = make(map[interface{}]interface{})
	sess.Options.MaxAge = -1
	sess.Options.Path = "/"

	// Save session - ignore errors for speed
	sess.Save(ctx.Request(), ctx.Response())

	// Clear user cache asynchronously
	if userID > 0 {
		go func() {
			cacheKey := fmt.Sprintf("user:%d", userID)
			c.cache.Flush().Key(cacheKey).Execute(ctx.Request().Context())
		}()
	}

	return nil
}

// GetAuthenticatedUserID returns the authenticated user's ID, if the user is logged in
func (c *AuthClient) GetAuthenticatedUserID(ctx echo.Context) (int, error) {
	sess, err := session.Get(ctx, authSessionName)
	if err != nil {
		return 0, err
	}

	if sess.Values[authSessionKeyAuthenticated] == true {
		return sess.Values[authSessionKeyUserID].(int), nil
	}

	return 0, NotAuthenticatedError{}
}

// GetAuthenticatedUser returns the authenticated user if the user is logged in
func (c *AuthClient) GetAuthenticatedUser(ctx echo.Context) (*ent.User, error) {
	userID, err := c.GetAuthenticatedUserID(ctx)
	if err != nil {
		return nil, NotAuthenticatedError{}
	}

	// Try to get user from cache first
	cacheKey := fmt.Sprintf("user:%d", userID)
	if cachedUser, err := c.cache.Get().Key(cacheKey).Fetch(ctx.Request().Context()); err == nil {
		if user, ok := cachedUser.(*ent.User); ok {
			return user, nil
		}
	}

	// If not in cache, fetch from database
	user, err := c.orm.User.Query().
		Where(user.ID(userID)).
		Only(ctx.Request().Context())
	if err != nil {
		return nil, err
	}

	// Cache the user for future requests
	c.cache.Set().
		Key(cacheKey).
		Data(user).
		Expiration(c.config.Cache.Expiration.UserSession).
		Save(ctx.Request().Context())

	return user, nil
}

// CheckPassword check if a given password matches a given hash
func (c *AuthClient) CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// HashPassword hashes a plain text password using bcrypt
func (c *AuthClient) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// GeneratePasswordResetToken generates a password reset token for a given user.
// For security purposes, the token itself is not stored in the database but rather
// a hash of the token, exactly how passwords are handled. This method returns both
// the generated token and the token entity which only contains the hash.
func (c *AuthClient) GeneratePasswordResetToken(ctx echo.Context, userID int) (string, *ent.PasswordToken, error) {
	// Generate the token, which is what will go in the URL, but not the database
	token, err := c.RandomToken(c.config.App.PasswordToken.Length)
	if err != nil {
		return "", nil, err
	}

	// Hash the token before storing it in the database
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, err
	}

	// Create and save the password reset token
	pt, err := c.orm.PasswordToken.
		Create().
		SetToken(string(hashedToken)).
		SetUserID(userID).
		Save(ctx.Request().Context())

	return token, pt, err
}

// GetValidPasswordToken returns a valid, non-expired password token entity for a given user, token ID and token.
// Since the actual token is not stored in the database for security purposes, if a matching password token entity is
// found a hash of the provided token is compared with the hash stored in the database in order to validate.
func (c *AuthClient) GetValidPasswordToken(ctx echo.Context, userID, tokenID int, token string) (*ent.PasswordToken, error) {
	// Ensure expired tokens are never returned
	expiration := time.Now().Add(-c.config.App.PasswordToken.Expiration)

	// Query to find a password token entity that matches the given user and token ID
	pt, err := c.orm.PasswordToken.
		Query().
		Where(passwordtoken.ID(tokenID)).
		Where(passwordtoken.HasUserWith(user.ID(userID))).
		Where(passwordtoken.CreatedAtGTE(expiration)).
		Only(ctx.Request().Context())

	switch err.(type) {
	case *ent.NotFoundError:
	case nil:
		// Check the token for a hash match
		if err := c.CheckPassword(token, pt.Token); err == nil {
			return pt, nil
		}
	default:
		if !context.IsCanceledError(err) {
			return nil, err
		}
	}

	return nil, InvalidPasswordTokenError{}
}

// DeletePasswordTokens deletes all password tokens in the database for a belonging to a given user.
// This should be called after a successful password reset.
func (c *AuthClient) DeletePasswordTokens(ctx echo.Context, userID int) error {
	_, err := c.orm.PasswordToken.
		Delete().
		Where(passwordtoken.HasUserWith(user.ID(userID))).
		Exec(ctx.Request().Context())

	return err
}

// RandomToken generates a random token string of a given length
func (c *AuthClient) RandomToken(length int) (string, error) {
	b := make([]byte, (length/2)+1)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)
	return token[:length], nil
}

// GenerateEmailVerificationToken generates an email verification token for a given email address using JWT which
// is set to expire based on the duration stored in configuration
func (c *AuthClient) GenerateEmailVerificationToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(c.config.App.EmailVerificationTokenExpiration).Unix(),
	})

	return token.SignedString([]byte(c.config.App.EncryptionKey))
}

// ValidateEmailVerificationToken validates an email verification token and returns the associated email address if
// the token is valid and has not expired
func (c *AuthClient) ValidateEmailVerificationToken(token string) (string, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(c.config.App.EncryptionKey), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		return claims["email"].(string), nil
	}

	return "", errors.New("invalid or expired token")
}

// GeneratePhoneVerificationToken generates a JWT token for phone verification
func (c *AuthClient) GeneratePhoneVerificationToken(phoneNumber string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"phone_number": phoneNumber,
		"exp":          time.Now().Add(c.config.App.EmailVerificationTokenExpiration).Unix(), // Reuse email config for now
	})

	return token.SignedString([]byte(c.config.App.EncryptionKey))
}

// ValidatePhoneVerificationToken validates a phone verification token and returns the associated phone number
func (c *AuthClient) ValidatePhoneVerificationToken(token string) (string, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(c.config.App.EncryptionKey), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		return claims["phone_number"].(string), nil
	}

	return "", errors.New("invalid or expired token")
}

// GenerateSimpleVerificationCode generates a simple 6-digit verification code
func (c *AuthClient) GenerateSimpleVerificationCode() (string, error) {
	// Generate a simple 6-digit code for WhatsApp verification
	b := make([]byte, 3)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// Convert to 6-digit string
	code := fmt.Sprintf("%06d", int(b[0])<<16|int(b[1])<<8|int(b[2]))
	return code[:6], nil
}

// GenerateWhatsAppPasswordResetToken generates a JWT token for WhatsApp password reset
func (c *AuthClient) GenerateWhatsAppPasswordResetToken(userID int, phoneNumber string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":      userID,
		"phone_number": phoneNumber,
		"type":         "whatsapp_password_reset",
		"exp":          time.Now().Add(1 * time.Hour).Unix(), // 1 hour expiry
	})

	return token.SignedString([]byte(c.config.App.EncryptionKey))
}

// ValidateWhatsAppPasswordResetToken validates a WhatsApp password reset token and returns user ID and phone number
func (c *AuthClient) ValidateWhatsAppPasswordResetToken(tokenString string) (int, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(c.config.App.EncryptionKey), nil
	})

	if err != nil {
		return 0, "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if tokenType, ok := claims["type"].(string); !ok || tokenType != "whatsapp_password_reset" {
			return 0, "", errors.New("invalid token type")
		}

		userIDFloat, userIDOk := claims["user_id"].(float64)
		phoneNumber, phoneOk := claims["phone_number"].(string)

		if !userIDOk || !phoneOk {
			return 0, "", errors.New("invalid token claims")
		}

		return int(userIDFloat), phoneNumber, nil
	}

	return 0, "", errors.New("invalid token")
}
