package handlers

import (
	"fmt"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/labstack/echo/v4"
	"gitlab.com/logan9312/discord-auction-bot/commands"
	"gitlab.com/logan9312/discord-auction-bot/config"
	"gitlab.com/logan9312/discord-auction-bot/crypto"
	"gitlab.com/logan9312/discord-auction-bot/database"
)

// WhitelabelResponse represents a whitelabel in API responses
type WhitelabelResponse struct {
	BotID     string `json:"bot_id"`
	BotName   string `json:"bot_name"`
	BotAvatar string `json:"bot_avatar"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

// WhitelabelListResponse is the response for listing whitelabels
type WhitelabelListResponse struct {
	Whitelabels []WhitelabelResponse `json:"whitelabels"`
	IsAdmin     bool                 `json:"is_admin"`
}

// WhitelabelCreateRequest for creating whitelabels
type WhitelabelCreateRequest struct {
	BotToken string `json:"bot_token"`
}

// ValidateTokenResponse for token validation
type ValidateTokenResponse struct {
	Valid     bool   `json:"valid"`
	BotID     string `json:"bot_id,omitempty"`
	BotName   string `json:"bot_name,omitempty"`
	BotAvatar string `json:"bot_avatar,omitempty"`
	Error     string `json:"error,omitempty"`
}

// ListWhitelabels returns whitelabels based on user access level
func ListWhitelabels(c echo.Context) error {
	userID := c.Get("user_id").(string)
	isAdmin := config.IsDevUser(userID)

	// Check premium status for non-admin users
	if !isAdmin && !commands.CheckPremiumUser(userID) {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "Premium subscription required for whitelabels",
		})
	}

	var whitelabels []database.WhiteLabels
	query := database.DB.Model(&database.WhiteLabels{})

	// Filter by user if not admin
	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&whitelabels).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch whitelabels",
		})
	}

	// Build response with bot info
	var response []WhitelabelResponse
	for _, wl := range whitelabels {
		// Only include CreatedAt if it's not the zero time
		createdAt := ""
		if !wl.CreatedAt.IsZero() {
			createdAt = wl.CreatedAt.Format("2006-01-02")
		}

		item := WhitelabelResponse{
			BotID:     wl.BotID,
			UserID:    wl.UserID,
			CreatedAt: createdAt,
		}

		// Try to get bot info from Discord
		botInfo, err := getBotInfo(wl)
		if err == nil {
			item.BotName = botInfo.Username
			item.BotAvatar = getBotAvatarURL(botInfo.ID, botInfo.Avatar)
		}

		response = append(response, item)
	}

	return c.JSON(http.StatusOK, WhitelabelListResponse{
		Whitelabels: response,
		IsAdmin:     isAdmin,
	})
}

// CreateWhitelabel adds or updates a whitelabel
func CreateWhitelabel(c echo.Context) error {
	userID := c.Get("user_id").(string)
	isAdmin := config.IsDevUser(userID)

	// Check premium for non-admin users
	if !isAdmin && !commands.CheckPremiumUser(userID) {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "Premium subscription required for whitelabels",
		})
	}

	var req WhitelabelCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if req.BotToken == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Bot token is required",
		})
	}

	// Validate token by connecting
	session, err := discordgo.New("Bot " + req.BotToken)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid bot token format",
		})
	}

	// Try to get bot user
	user, err := session.User("@me")
	if err != nil {
		session.Close()
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Could not authenticate with Discord. Please verify your token.",
		})
	}
	session.Close()

	// Encrypt token before storing
	tokenToStore := req.BotToken
	encrypted := false
	if crypto.IsInitialized() {
		encryptedToken, err := crypto.Encrypt(req.BotToken)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to secure token",
			})
		}
		tokenToStore = encryptedToken
		encrypted = true
	}

	// Save to database (upsert)
	wl := database.WhiteLabels{
		BotID:     user.ID,
		UserID:    userID,
		BotToken:  tokenToStore,
		Encrypted: encrypted,
	}

	if err := database.DB.Save(&wl).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to save whitelabel",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":    "Whitelabel saved successfully. Restart required to activate.",
		"bot_id":     user.ID,
		"bot_name":   user.Username,
		"bot_avatar": getBotAvatarURL(user.ID, user.Avatar),
	})
}

// DeleteWhitelabel removes a whitelabel
func DeleteWhitelabel(c echo.Context) error {
	userID := c.Get("user_id").(string)
	botID := c.Param("botId")
	isAdmin := config.IsDevUser(userID)

	query := database.DB.Where("bot_id = ?", botID)
	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}

	result := query.Delete(&database.WhiteLabels{})
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete whitelabel",
		})
	}

	if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Whitelabel not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Whitelabel deleted successfully",
	})
}

// ValidateToken checks if a token is valid without saving
func ValidateToken(c echo.Context) error {
	userID := c.Get("user_id").(string)
	isAdmin := config.IsDevUser(userID)

	// Check premium for non-admin users
	if !isAdmin && !commands.CheckPremiumUser(userID) {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "Premium subscription required for whitelabels",
		})
	}

	var req WhitelabelCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusOK, ValidateTokenResponse{
			Valid: false,
			Error: "Invalid request body",
		})
	}

	if req.BotToken == "" {
		return c.JSON(http.StatusOK, ValidateTokenResponse{
			Valid: false,
			Error: "Bot token is required",
		})
	}

	session, err := discordgo.New("Bot " + req.BotToken)
	if err != nil {
		return c.JSON(http.StatusOK, ValidateTokenResponse{
			Valid: false,
			Error: "Invalid token format",
		})
	}

	user, err := session.User("@me")
	if err != nil {
		session.Close()
		return c.JSON(http.StatusOK, ValidateTokenResponse{
			Valid: false,
			Error: "Could not authenticate with Discord. Please verify your token.",
		})
	}
	session.Close()

	return c.JSON(http.StatusOK, ValidateTokenResponse{
		Valid:     true,
		BotID:     user.ID,
		BotName:   user.Username,
		BotAvatar: getBotAvatarURL(user.ID, user.Avatar),
	})
}

// getBotInfo attempts to get bot info by decrypting the token and querying Discord
func getBotInfo(wl database.WhiteLabels) (*discordgo.User, error) {
	token := wl.BotToken
	if wl.Encrypted && crypto.IsInitialized() {
		decrypted, err := crypto.Decrypt(wl.BotToken)
		if err != nil {
			return nil, err
		}
		token = decrypted
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	user, err := session.User("@me")
	if err != nil {
		return nil, err
	}

	return user, nil
}

// getBotAvatarURL constructs the Discord CDN URL for a bot avatar
func getBotAvatarURL(botID, avatar string) string {
	if avatar == "" {
		return ""
	}
	return fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", botID, avatar)
}
