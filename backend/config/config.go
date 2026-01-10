package config

import (
	"os"
	"strings"
)

// Config holds all configurable values for the application
type Config struct {
	// Discord IDs
	DevUserIDs        []string // Users with developer access
	PremiumServerIDs  []string // Servers with permanent premium
	PremiumUserIDs    []string // Users with permanent premium
	SupportServerID   string   // Support server guild ID
	PremiumRoleID     string   // Premium role ID in support server
	DevCommandsGuild  string   // Guild ID for dev commands registration
	RestartChannelID  string   // Channel to send restart notifications
	GuildJoinChannel  string   // Channel to send guild join notifications
	WhitelabelBotID   string   // Bot ID for special whitelabel handling

	// URLs
	SupportServerURL string // Discord invite URL for support
	SuccessURL       string // Stripe success redirect URL
	CancelURL        string // Stripe cancel redirect URL

	// Stripe
	StripePriceID        string // Stripe price ID for premium
	StripeWebhookSecret  string // Stripe webhook signing secret

	// Support Button
	SupportButtonEmojiName string
	SupportButtonEmojiID   string

	// OAuth2 / Dashboard
	DiscordClientID     string // Discord OAuth2 client ID
	DiscordClientSecret string // Discord OAuth2 client secret
	JWTSecret           string // Secret for signing JWT tokens
	FrontendURL         string // Frontend URL for CORS and redirects
	APIBaseURL          string // API base URL (should include /api prefix, e.g., https://api.prizebot.dev/api)
	SecureCookies       bool   // Whether to set Secure flag on cookies (true in prod)
}

// Global config instance
var C Config

// Init initializes the config from environment variables
func Init() {
	C = Config{
		// Discord IDs - comma-separated for multiple values
		DevUserIDs:        parseCSV(getEnv("DEV_USER_IDS", "280812467775471627")),
		PremiumServerIDs:  parseCSV(getEnv("PREMIUM_SERVER_IDS", "915767892467920967,626094990984216586")),
		PremiumUserIDs:    parseCSV(getEnv("PREMIUM_USER_IDS", "280812467775471627")),
		SupportServerID:   getEnv("SUPPORT_SERVER_ID", "885228283573178408"),
		PremiumRoleID:     getEnv("PREMIUM_ROLE_ID", "942927890100682752"),
		DevCommandsGuild:  getEnv("DEV_COMMANDS_GUILD", "915767892467920967"),
		RestartChannelID:  getEnv("RESTART_CHANNEL_ID", "943175605858496602"),
		GuildJoinChannel:  getEnv("GUILD_JOIN_CHANNEL", "1005255087200948234"),
		WhitelabelBotID:   getEnv("WHITELABEL_BOT_ID", "995022149226082324"),

		// URLs
		SupportServerURL: getEnv("SUPPORT_SERVER_URL", "https://discord.gg/RxP2z5NGtj"),
		SuccessURL:       getEnv("STRIPE_SUCCESS_URL", "https://www.aftmgaming.com/auction-bot/success"),
		CancelURL:        getEnv("STRIPE_CANCEL_URL", "https://discord.gg/YBRvZ3mRtb"),

		// Stripe
		StripePriceID:       getEnv("STRIPE_PRICE_ID", "price_1KYE0EKpOiJyve6nT9Qo9IfN"),
		StripeWebhookSecret: getEnv("STRIPE_WEBHOOK_SECRET", ""),

		// Support Button
		SupportButtonEmojiName: getEnv("SUPPORT_BUTTON_EMOJI_NAME", "logo"),
		SupportButtonEmojiID:   getEnv("SUPPORT_BUTTON_EMOJI_ID", "889025400120950804"),

		// OAuth2 / Dashboard
		DiscordClientID:     getEnv("DISCORD_CLIENT_ID", ""),
		DiscordClientSecret: getEnv("DISCORD_CLIENT_SECRET", ""),
		JWTSecret:           getEnv("JWT_SECRET", "change-me-in-production-32chars"),
		FrontendURL:         getEnv("FRONTEND_URL", "http://localhost:5173"),
		APIBaseURL:          getEnv("API_BASE_URL", "http://localhost:8080/api"),
		SecureCookies:       getEnv("ENVIRONMENT", "") == "prod",
	}
}

// IsDevUser checks if a user ID is in the dev users list
func IsDevUser(userID string) bool {
	for _, id := range C.DevUserIDs {
		if id == userID {
			return true
		}
	}
	return false
}

// IsPremiumServer checks if a server ID is in the premium servers list
func IsPremiumServer(serverID string) bool {
	for _, id := range C.PremiumServerIDs {
		if id == serverID {
			return true
		}
	}
	return false
}

// IsPremiumUser checks if a user ID is in the premium users list
func IsPremiumUser(userID string) bool {
	for _, id := range C.PremiumUserIDs {
		if id == userID {
			return true
		}
	}
	return false
}

// getEnv returns the environment variable value or a default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseCSV splits a comma-separated string into a slice
func parseCSV(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
