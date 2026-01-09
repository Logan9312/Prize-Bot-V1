package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/labstack/echo/v4"
)

// BotSession holds the Discord bot session for fetching guild data
var BotSession *discordgo.Session

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

	return c.JSON(http.StatusOK, map[string]interface{}{
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

	return c.JSON(http.StatusOK, map[string]interface{}{
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

	return c.JSON(http.StatusOK, map[string]interface{}{
		"roles": roles,
	})
}

// verifyGuildAccess checks if the user has admin access to the guild
func verifyGuildAccess(c echo.Context, guildID string) error {
	accessToken := c.Get("access_token").(string)

	// Fetch user's guilds from Discord
	req, err := http.NewRequest("GET", "https://discord.com/api/v10/users/@me/guilds", nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to verify access"})
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to verify access"})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Discord API error: %s", string(body))})
	}

	var guilds []DiscordGuild
	if err := json.NewDecoder(resp.Body).Decode(&guilds); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse guilds"})
	}

	for _, g := range guilds {
		if g.ID == guildID {
			var perms int64
			fmt.Sscanf(g.Permissions, "%d", &perms)
			if (perms&0x8) == 0x8 || (perms&0x20) == 0x20 || g.Owner {
				return nil
			}
		}
	}

	return c.JSON(http.StatusForbidden, map[string]string{"error": "You don't have access to this guild"})
}
