package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"gitlab.com/logan9312/discord-auction-bot/api/middleware"
	"gitlab.com/logan9312/discord-auction-bot/config"
)

const (
	discordAPIBase   = "https://discord.com/api/v10"
	discordAuthURL   = "https://discord.com/api/oauth2/authorize"
	discordTokenURL  = "https://discord.com/api/oauth2/token"
	discordUserURL   = "https://discord.com/api/v10/users/@me"
	discordGuildsURL = "https://discord.com/api/v10/users/@me/guilds"
)

// DiscordTokenResponse represents the response from Discord's token endpoint
type DiscordTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// DiscordUser represents a Discord user
type DiscordUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	GlobalName    string `json:"global_name"`
}

// generateState creates a random state string for OAuth2 CSRF protection
func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// DiscordOAuthRedirect redirects the user to Discord's OAuth2 authorization page
func DiscordOAuthRedirect(c echo.Context) error {
	if config.C.DiscordClientID == "" {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Discord OAuth not configured"})
	}

	state := generateState()

	// Store state in a cookie for verification
	cookie := &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production
		SameSite: http.SameSiteLaxMode,
		MaxAge:   300, // 5 minutes
	}
	c.SetCookie(cookie)

	redirectURI := config.C.APIBaseURL + "/auth/discord/callback"

	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s",
		discordAuthURL,
		config.C.DiscordClientID,
		url.QueryEscape(redirectURI),
		url.QueryEscape("identify guilds"),
		state,
	)

	return c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// DiscordOAuthCallback handles the OAuth2 callback from Discord
func DiscordOAuthCallback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	if code == "" {
		return c.Redirect(http.StatusTemporaryRedirect, config.C.FrontendURL+"/login?error=no_code")
	}

	// Verify state
	stateCookie, err := c.Cookie("oauth_state")
	if err != nil || stateCookie.Value != state {
		return c.Redirect(http.StatusTemporaryRedirect, config.C.FrontendURL+"/login?error=invalid_state")
	}

	// Exchange code for token
	tokenResp, err := exchangeCodeForToken(code)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, config.C.FrontendURL+"/login?error=token_exchange_failed")
	}

	// Get user info
	user, err := getDiscordUser(tokenResp.AccessToken)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, config.C.FrontendURL+"/login?error=user_fetch_failed")
	}

	// Calculate token expiry
	tokenExpiry := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).Unix()

	// Generate JWT
	jwtToken, err := middleware.GenerateToken(
		user.ID,
		user.Username,
		user.Avatar,
		tokenResp.AccessToken,
		tokenResp.RefreshToken,
		tokenExpiry,
	)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, config.C.FrontendURL+"/login?error=jwt_generation_failed")
	}

	// Set auth cookie
	middleware.SetAuthCookie(c, jwtToken)

	// Clear oauth state cookie
	c.SetCookie(&http.Cookie{
		Name:   "oauth_state",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	return c.Redirect(http.StatusTemporaryRedirect, config.C.FrontendURL+"/dashboard")
}

// exchangeCodeForToken exchanges the authorization code for access and refresh tokens
func exchangeCodeForToken(code string) (*DiscordTokenResponse, error) {
	redirectURI := config.C.APIBaseURL + "/auth/discord/callback"

	data := url.Values{}
	data.Set("client_id", config.C.DiscordClientID)
	data.Set("client_secret", config.C.DiscordClientSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)

	req, err := http.NewRequest("POST", discordTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed: %s", string(body))
	}

	var tokenResp DiscordTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// getDiscordUser fetches the user's profile from Discord
func getDiscordUser(accessToken string) (*DiscordUser, error) {
	req, err := http.NewRequest("GET", discordUserURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("user fetch failed: %s", string(body))
	}

	var user DiscordUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// Logout clears the authentication cookie
func Logout(c echo.Context) error {
	middleware.ClearAuthCookie(c)
	return c.JSON(http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

// GetCurrentUser returns the current authenticated user's info
func GetCurrentUser(c echo.Context) error {
	userID := c.Get("user_id").(string)
	username := c.Get("username").(string)
	avatar := c.Get("avatar").(string)

	avatarURL := ""
	if avatar != "" {
		avatarURL = fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", userID, avatar)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":         userID,
		"username":   username,
		"avatar":     avatar,
		"avatar_url": avatarURL,
	})
}
