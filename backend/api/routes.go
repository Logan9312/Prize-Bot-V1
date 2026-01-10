package api

import (
	"github.com/bwmarrin/discordgo"
	"github.com/labstack/echo/v4"
	"gitlab.com/logan9312/discord-auction-bot/api/handlers"
	"gitlab.com/logan9312/discord-auction-bot/api/middleware"
)

// RegisterRoutes registers all API routes on the Echo instance
func RegisterRoutes(e *echo.Echo, botSession *discordgo.Session) {
	// Set the bot session for handlers that need it
	handlers.BotSession = botSession

	// API group with CORS
	api := e.Group("/api")
	api.Use(middleware.CORSConfig())

	// Auth routes (no auth required)
	auth := api.Group("/auth")
	auth.GET("/discord", handlers.DiscordOAuthRedirect)
	auth.GET("/discord/callback", handlers.DiscordOAuthCallback)
	auth.POST("/logout", handlers.Logout)

	// Protected routes (auth required)
	protected := api.Group("")
	protected.Use(middleware.JWTAuth)

	// User info
	protected.GET("/auth/me", handlers.GetCurrentUser)

	// Guilds
	protected.GET("/guilds", handlers.ListUserGuilds)
	protected.GET("/guilds/:guildId/channels", handlers.GetGuildChannels)
	protected.GET("/guilds/:guildId/roles", handlers.GetGuildRoles)
	protected.GET("/guilds/:guildId/stats", handlers.GetGuildStats)

	// List endpoints
	protected.GET("/guilds/:guildId/auctions/list", handlers.ListAuctions)
	protected.GET("/guilds/:guildId/giveaways/list", handlers.ListGiveaways)
	protected.GET("/guilds/:guildId/claims/list", handlers.ListClaims)

	// Auction Settings
	protected.GET("/guilds/:guildId/settings/auction", handlers.GetAuctionSettings)
	protected.PUT("/guilds/:guildId/settings/auction", handlers.UpdateAuctionSettings)
	protected.DELETE("/guilds/:guildId/settings/auction", handlers.DeleteAuctionSettings)

	// Claim Settings
	protected.GET("/guilds/:guildId/settings/claim", handlers.GetClaimSettings)
	protected.PUT("/guilds/:guildId/settings/claim", handlers.UpdateClaimSettings)
	protected.DELETE("/guilds/:guildId/settings/claim", handlers.DeleteClaimSettings)

	// Giveaway Settings
	protected.GET("/guilds/:guildId/settings/giveaway", handlers.GetGiveawaySettings)
	protected.PUT("/guilds/:guildId/settings/giveaway", handlers.UpdateGiveawaySettings)
	protected.DELETE("/guilds/:guildId/settings/giveaway", handlers.DeleteGiveawaySettings)

	// Currency Settings
	protected.GET("/guilds/:guildId/settings/currency", handlers.GetCurrencySettings)
	protected.PUT("/guilds/:guildId/settings/currency", handlers.UpdateCurrencySettings)
	protected.DELETE("/guilds/:guildId/settings/currency", handlers.DeleteCurrencySettings)

	// Shop Settings
	protected.GET("/guilds/:guildId/settings/shop", handlers.GetShopSettings)
	protected.PUT("/guilds/:guildId/settings/shop", handlers.UpdateShopSettings)
	protected.DELETE("/guilds/:guildId/settings/shop", handlers.DeleteShopSettings)
}
