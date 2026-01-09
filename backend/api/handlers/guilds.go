package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/labstack/echo/v4"
	"gitlab.com/logan9312/discord-auction-bot/logger"
)

// BotSession holds the Discord bot session for fetching guild data
var BotSession *discordgo.Session

// accessVerificationCache caches OAuth guild access verification to prevent rate limiting
var accessVerificationCache = &sync.Map{} // map[string]cacheEntry

type cacheEntry struct {
	hasAccess bool
	expiresAt time.Time
}

// getCacheKey generates a cache key from the access token and guild ID
func getCacheKey(accessToken, guildID string) string {
	// Hash the token to avoid storing full tokens in memory
	hash := sha256.Sum256([]byte(accessToken + ":" + guildID))
	return hex.EncodeToString(hash[:16]) // Use first 16 bytes for shorter key
}

// Guild represents a Discord guild for the API response
type Guild struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Icon     string `json:"icon"`
	IconURL  string `json:"icon_url"`
	Owner    bool   `json:"owner"`
	IsAdmin  bool   `json:"is_admin"`
	BotIn    bool   `json:"bot_in"`
}

// Channel represents a Discord channel for the API response
type Channel struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     int    `json:"type"`
	Position int    `json:"position"`
	ParentID string `json:"parent_id"`
}

// Role represents a Discord role for the API response
type Role struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Color    int    `json:"color"`
	Position int    `json:"position"`
	Managed  bool   `json:"managed"`
}

// DiscordGuild represents a guild from Discord's API
type DiscordGuild struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Owner       bool   `json:"owner"`
	Permissions string `json:"permissions"`
}

// ListUserGuilds returns the user's guilds where they have admin permissions
func ListUserGuilds(c echo.Context) error {
	accessToken := c.Get("access_token").(string)

	// Fetch user's guilds from Discord
	req, err := http.NewRequest("GET", "https://discord.com/api/v10/users/@me/guilds", nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create request"})
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch guilds"})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Discord API error: %s", string(body))})
	}

	var discordGuilds []DiscordGuild
	if err := json.NewDecoder(resp.Body).Decode(&discordGuilds); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse guilds"})
	}

	// Filter for guilds where user has admin or manage_guild permission
	var guilds []Guild
	for _, dg := range discordGuilds {
		// Parse permissions as int64
		var perms int64
		fmt.Sscanf(dg.Permissions, "%d", &perms)

		// Check for ADMINISTRATOR (0x8) or MANAGE_GUILD (0x20)
		isAdmin := (perms&0x8) == 0x8 || (perms&0x20) == 0x20 || dg.Owner

		if !isAdmin {
			continue
		}

		iconURL := ""
		if dg.Icon != "" {
			iconURL = fmt.Sprintf("https://cdn.discordapp.com/icons/%s/%s.png", dg.ID, dg.Icon)
		}

		// Check if bot is in this guild
		botIn := false
		if BotSession != nil {
			_, err := BotSession.State.Guild(dg.ID)
			botIn = err == nil
		}

		guilds = append(guilds, Guild{
			ID:      dg.ID,
			Name:    dg.Name,
			Icon:    dg.Icon,
			IconURL: iconURL,
			Owner:   dg.Owner,
			IsAdmin: isAdmin,
			BotIn:   botIn,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"guilds": guilds,
	})
}

// GetGuildChannels returns the channels for a guild using the bot session
func GetGuildChannels(c echo.Context) error {
	guildID := c.Param("guildId")

	if BotSession == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Bot session not available"})
	}

	// Verify user has access to this guild
	if err := verifyGuildAccess(c, guildID); err != nil {
		return err
	}

	// Get channels from bot's state or fetch from Discord
	guildChannels, err := BotSession.GuildChannels(guildID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch channels"})
	}

	var channels []Channel
	for _, ch := range guildChannels {
		// Only include text channels (0) and categories (4)
		if ch.Type != discordgo.ChannelTypeGuildText && ch.Type != discordgo.ChannelTypeGuildCategory {
			continue
		}

		channels = append(channels, Channel{
			ID:       ch.ID,
			Name:     ch.Name,
			Type:     int(ch.Type),
			Position: ch.Position,
			ParentID: ch.ParentID,
		})
	}

	// Sort by position
	sort.Slice(channels, func(i, j int) bool {
		return channels[i].Position < channels[j].Position
	})

	return c.JSON(http.StatusOK, map[string]any{
		"channels": channels,
	})
}

// GetGuildRoles returns the roles for a guild using the bot session
func GetGuildRoles(c echo.Context) error {
	guildID := c.Param("guildId")

	if BotSession == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Bot session not available"})
	}

	// Verify user has access to this guild
	if err := verifyGuildAccess(c, guildID); err != nil {
		return err
	}

	// Get roles from Discord
	guildRoles, err := BotSession.GuildRoles(guildID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch roles"})
	}

	var roles []Role
	for _, r := range guildRoles {
		// Skip @everyone role and managed roles (bots)
		if r.Name == "@everyone" {
			continue
		}

		roles = append(roles, Role{
			ID:       r.ID,
			Name:     r.Name,
			Color:    r.Color,
			Position: r.Position,
			Managed:  r.Managed,
		})
	}

	// Sort by position (highest first)
	sort.Slice(roles, func(i, j int) bool {
		return roles[i].Position > roles[j].Position
	})

	return c.JSON(http.StatusOK, map[string]any{
		"roles": roles,
	})
}

// verifyGuildAccess checks if the user has admin access to the guild
func verifyGuildAccess(c echo.Context, guildID string) error {
	accessToken := c.Get("access_token").(string)
	cacheKey := getCacheKey(accessToken, guildID)

	// Check cache first to avoid rate limiting
	if cached, ok := accessVerificationCache.Load(cacheKey); ok {
		entry := cached.(cacheEntry)
		if time.Now().Before(entry.expiresAt) {
			// Cache hit and still valid
			if !entry.hasAccess {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "You don't have access to this guild"})
			}
			return nil
		}
		// Cache expired, remove it
		accessVerificationCache.Delete(cacheKey)
	}

	// Retry logic for handling transient rate limits
	maxRetries := 2
	client := &http.Client{Timeout: 10 * time.Second}
	var guilds []DiscordGuild

	for attempt := range maxRetries {
		// Fetch user's guilds from Discord
		req, err := http.NewRequest("GET", "https://discord.com/api/v10/users/@me/guilds", nil)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to verify access"})
		}
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := client.Do(req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to verify access"})
		}

		// Handle rate limiting with automatic retry
		if resp.StatusCode == http.StatusTooManyRequests {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			// Log the actual Discord rate limit response
			logger.Sugar.Warnw("Discord OAuth rate limit hit",
				"guild_id", guildID,
				"attempt", attempt+1,
				"max_retries", maxRetries,
				"response_body", string(body),
			)

			// Parse retry_after from response
			var rateLimitResponse struct {
				RetryAfter float64 `json:"retry_after"`
				Global     bool    `json:"global"`
				Message    string  `json:"message"`
			}

			if json.Unmarshal(body, &rateLimitResponse) == nil && attempt < maxRetries-1 {
				// Wait for the retry_after duration plus a small buffer, then retry
				waitDuration := time.Duration(rateLimitResponse.RetryAfter*1000) * time.Millisecond
				logger.Sugar.Infow("Retrying after rate limit wait",
					"guild_id", guildID,
					"wait_duration_ms", waitDuration.Milliseconds(),
					"retry_after", rateLimitResponse.RetryAfter,
					"global", rateLimitResponse.Global,
				)
				time.Sleep(waitDuration + 100*time.Millisecond)
				continue
			}

			// All retries exhausted
			logger.Sugar.Errorw("Rate limit retries exhausted",
				"guild_id", guildID,
				"attempts", attempt+1,
			)
			return c.JSON(http.StatusTooManyRequests, map[string]string{"error": "Rate limited, please try again in a moment"})
		}

		// Handle other non-success status codes
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Discord API error: %s", string(body))})
		}

		// Success - parse the guilds
		if err := json.NewDecoder(resp.Body).Decode(&guilds); err != nil {
			resp.Body.Close()
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse guilds"})
		}
		resp.Body.Close()

		// Log successful verification (especially useful if it was after a retry)
		if attempt > 0 {
			logger.Sugar.Infow("Guild access verification succeeded after retry",
				"guild_id", guildID,
				"attempt", attempt+1,
			)
		}
		break
	}

	// Check if user has access to the requested guild
	for _, g := range guilds {
		if g.ID == guildID {
			var perms int64
			fmt.Sscanf(g.Permissions, "%d", &perms)
			hasAccess := (perms&0x8) == 0x8 || (perms&0x20) == 0x20 || g.Owner

			// Cache the result for 10 seconds to prevent rate limiting on subsequent requests
			accessVerificationCache.Store(cacheKey, cacheEntry{
				hasAccess: hasAccess,
				expiresAt: time.Now().Add(10 * time.Second),
			})

			if hasAccess {
				return nil
			}
			return c.JSON(http.StatusForbidden, map[string]string{"error": "You don't have access to this guild"})
		}
	}

	// Guild not found in user's guilds - cache as no access
	accessVerificationCache.Store(cacheKey, cacheEntry{
		hasAccess: false,
		expiresAt: time.Now().Add(10 * time.Second),
	})
	return c.JSON(http.StatusForbidden, map[string]string{"error": "You don't have access to this guild"})
}
