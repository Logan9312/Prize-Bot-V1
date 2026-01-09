package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gitlab.com/logan9312/discord-auction-bot/config"
)

// JWTClaims represents the claims stored in the JWT token
type JWTClaims struct {
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	Avatar       string `json:"avatar"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenExpiry  int64  `json:"token_expiry"`
	jwt.RegisteredClaims
}

// JWTAuth is middleware that validates JWT tokens from cookies or Authorization header
func JWTAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := ""

		// Try to get token from cookie first
		cookie, err := c.Cookie("auth_token")
		if err == nil {
			tokenString = cookie.Value
		}

		// If no cookie, try Authorization header
		if tokenString == "" {
			authHeader := c.Request().Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "No authentication token provided"})
		}

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.C.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid or expired token"})
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token claims"})
		}

		// Store claims in context for handlers to use
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("avatar", claims.Avatar)
		c.Set("access_token", claims.AccessToken)
		c.Set("refresh_token", claims.RefreshToken)
		c.Set("token_expiry", claims.TokenExpiry)

		return next(c)
	}
}

// GenerateToken creates a new JWT token with the given claims
func GenerateToken(userID, username, avatar, accessToken, refreshToken string, tokenExpiry int64) (string, error) {
	claims := JWTClaims{
		UserID:       userID,
		Username:     username,
		Avatar:       avatar,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenExpiry:  tokenExpiry,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.C.JWTSecret))
}

// SetAuthCookie sets the authentication cookie
func SetAuthCookie(c echo.Context, token string) {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   config.C.SecureCookies,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   7 * 24 * 60 * 60, // 7 days
	}
	c.SetCookie(cookie)
}

// ClearAuthCookie clears the authentication cookie
func ClearAuthCookie(c echo.Context) {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	}
	c.SetCookie(cookie)
}
