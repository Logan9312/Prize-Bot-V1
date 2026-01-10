package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gitlab.com/logan9312/discord-auction-bot/database"
)

// GuildStatsResponse represents the response for guild statistics
type GuildStatsResponse struct {
	ActiveAuctions   int64 `json:"active_auctions"`
	RunningGiveaways int64 `json:"running_giveaways"`
	OpenClaims       int64 `json:"open_claims"`
	ShopItems        int64 `json:"shop_items"`
}

// GetGuildStats returns statistics for a guild
func GetGuildStats(c echo.Context) error {
	guildID := c.Param("guildId")

	if err := verifyGuildAccess(c, guildID); err != nil {
		return err
	}

	var stats GuildStatsResponse

	// Count active auctions (end time in the future)
	database.DB.Model(&database.Auction{}).
		Where("guild_id = ? AND end_time > ?", guildID, time.Now()).
		Count(&stats.ActiveAuctions)

	// Count running giveaways (not finished)
	database.DB.Model(&database.Giveaway{}).
		Where("guild_id = ? AND finished = ?", guildID, false).
		Count(&stats.RunningGiveaways)

	// Count open claims (status = pending, or legacy NULL/empty status)
	database.DB.Model(&database.Claim{}).
		Where("guild_id = ?", guildID).
		Where("status = ? OR status IS NULL OR status = ''", database.ClaimStatusPending).
		Count(&stats.OpenClaims)

	// Shop items - currently no ShopItem model exists, so return 0
	stats.ShopItems = 0

	return c.JSON(http.StatusOK, stats)
}

// AuctionSettingsRequest represents the request body for auction settings
type AuctionSettingsRequest struct {
	Category        string `json:"category"`
	AlertRole       string `json:"alert_role"`
	Currency        string `json:"currency"`
	LogChannel      string `json:"log_channel"`
	HostRole        string `json:"host_role"`
	SnipeExtension  int64  `json:"snipe_extension"` // milliseconds
	SnipeRange      int64  `json:"snipe_range"`
	SnipeLimit      int64  `json:"snipe_limit"`
	SnipeCap        int64  `json:"snipe_cap"`
	CurrencySide    string `json:"currency_side"`
	IntegerOnly     bool   `json:"integer_only"`
	ChannelOverride string `json:"channel_override"`
	ChannelLock     bool   `json:"channel_lock"`
	ChannelPrefix   string `json:"channel_prefix"`
	UseCurrency     bool   `json:"use_currency"`
}

// ClaimSettingsRequest represents the request body for claim settings
type ClaimSettingsRequest struct {
	Category        string `json:"category"`
	StaffRole       string `json:"staff_role"`
	Instructions    string `json:"instructions"`
	LogChannel      string `json:"log_channel"`
	Expiration      string `json:"expiration"`
	DisableClaiming bool   `json:"disable_claiming"`
	ChannelPrefix   string `json:"channel_prefix"`
}

// GiveawaySettingsRequest represents the request body for giveaway settings
type GiveawaySettingsRequest struct {
	HostRole   string `json:"host_role"`
	AlertRole  string `json:"alert_role"`
	LogChannel string `json:"log_channel"`
}

// CurrencySettingsRequest represents the request body for currency settings
type CurrencySettingsRequest struct {
	Currency string `json:"currency"`
	Side     string `json:"side"`
}

// ShopSettingsRequest represents the request body for shop settings
type ShopSettingsRequest struct {
	HostRole   string `json:"host_role"`
	AlertRole  string `json:"alert_role"`
	LogChannel string `json:"log_channel"`
}

// Auction Settings Handlers

func GetAuctionSettings(c echo.Context) error {
	guildID := c.Param("guildId")

	if err := verifyGuildAccess(c, guildID); err != nil {
		return err
	}

	var settings database.AuctionSetup
	result := database.DB.Where("guild_id = ?", guildID).First(&settings)

	if result.Error != nil {
		// Return empty settings if not found
		return c.JSON(http.StatusOK, map[string]interface{}{
			"guild_id": guildID,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"guild_id":         settings.GuildID,
		"category":         settings.Category,
		"alert_role":       settings.AlertRole,
		"currency":         settings.Currency,
		"log_channel":      settings.LogChannel,
		"host_role":        settings.HostRole,
		"snipe_extension":  settings.SnipeExtension.Milliseconds(),
		"snipe_range":      settings.SnipeRange.Milliseconds(),
		"snipe_limit":      settings.SnipeLimit.Milliseconds(),
		"snipe_cap":        settings.SnipeCap.Milliseconds(),
		"currency_side":    settings.CurrencySide,
		"integer_only":     settings.IntegerOnly,
		"channel_override": settings.ChannelOverride,
		"channel_lock":     settings.ChannelLock,
		"channel_prefix":   settings.ChannelPrefix,
		"use_currency":     settings.UseCurrency,
	})
}

func UpdateAuctionSettings(c echo.Context) error {
	guildID := c.Param("guildId")

	if err := verifyGuildAccess(c, guildID); err != nil {
		return err
	}

	var req AuctionSettingsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	settings := database.AuctionSetup{
		GuildID:         guildID,
		Category:        req.Category,
		AlertRole:       req.AlertRole,
		Currency:        req.Currency,
		LogChannel:      req.LogChannel,
		HostRole:        req.HostRole,
		SnipeExtension:  time.Duration(req.SnipeExtension) * time.Millisecond,
		SnipeRange:      time.Duration(req.SnipeRange) * time.Millisecond,
		SnipeLimit:      time.Duration(req.SnipeLimit) * time.Millisecond,
		SnipeCap:        time.Duration(req.SnipeCap) * time.Millisecond,
		CurrencySide:    req.CurrencySide,
		IntegerOnly:     req.IntegerOnly,
		ChannelOverride: req.ChannelOverride,
		ChannelLock:     req.ChannelLock,
		ChannelPrefix:   req.ChannelPrefix,
		UseCurrency:     req.UseCurrency,
	}

	result := database.DB.Save(&settings)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save settings"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Settings updated successfully"})
}

func DeleteAuctionSettings(c echo.Context) error {
	guildID := c.Param("guildId")

	if err := verifyGuildAccess(c, guildID); err != nil {
		return err
	}

	result := database.DB.Delete(&database.AuctionSetup{}, "guild_id = ?", guildID)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete settings"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Settings deleted successfully"})
}

// Claim Settings Handlers

func GetClaimSettings(c echo.Context) error {
	guildID := c.Param("guildId")

	if err := verifyGuildAccess(c, guildID); err != nil {
		return err
	}

	var settings database.ClaimSetup
	result := database.DB.Where("guild_id = ?", guildID).First(&settings)

	if result.Error != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"guild_id": guildID,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"guild_id":         settings.GuildID,
		"category":         settings.Category,
		"staff_role":       settings.StaffRole,
		"instructions":     settings.Instructions,
		"log_channel":      settings.LogChannel,
		"expiration":       settings.Expiration,
		"disable_claiming": settings.DisableClaiming,
		"channel_prefix":   settings.ChannelPrefix,
	})
}

func UpdateClaimSettings(c echo.Context) error {
	guildID := c.Param("guildId")

	if err := verifyGuildAccess(c, guildID); err != nil {
		return err
	}

	var req ClaimSettingsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	settings := database.ClaimSetup{
		GuildID:         guildID,
		Category:        req.Category,
		StaffRole:       req.StaffRole,
		Instructions:    req.Instructions,
		LogChannel:      req.LogChannel,
		Expiration:      req.Expiration,
		DisableClaiming: req.DisableClaiming,
		ChannelPrefix:   req.ChannelPrefix,
	}

	result := database.DB.Save(&settings)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save settings"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Settings updated successfully"})
}

func DeleteClaimSettings(c echo.Context) error {
	guildID := c.Param("guildId")

	if err := verifyGuildAccess(c, guildID); err != nil {
		return err
	}

	result := database.DB.Delete(&database.ClaimSetup{}, "guild_id = ?", guildID)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete settings"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Settings deleted successfully"})
}

// Giveaway Settings Handlers

func GetGiveawaySettings(c echo.Context) error {
	guildID := c.Param("guildId")

	if err := verifyGuildAccess(c, guildID); err != nil {
		return err
	}

	var settings database.GiveawaySetup
	result := database.DB.Where("guild_id = ?", guildID).First(&settings)

	if result.Error != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"guild_id": guildID,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"guild_id":    settings.GuildID,
		"host_role":   settings.HostRole,
		"alert_role":  settings.AlertRole,
		"log_channel": settings.LogChannel,
	})
}

func UpdateGiveawaySettings(c echo.Context) error {
	guildID := c.Param("guildId")

	if err := verifyGuildAccess(c, guildID); err != nil {
		return err
	}

	var req GiveawaySettingsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	settings := database.GiveawaySetup{
		GuildID:    guildID,
		HostRole:   req.HostRole,
		AlertRole:  req.AlertRole,
		LogChannel: req.LogChannel,
	}

	result := database.DB.Save(&settings)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save settings"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Settings updated successfully"})
}

func DeleteGiveawaySettings(c echo.Context) error {
	guildID := c.Param("guildId")

	if err := verifyGuildAccess(c, guildID); err != nil {
		return err
	}

	result := database.DB.Delete(&database.GiveawaySetup{}, "guild_id = ?", guildID)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete settings"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Settings deleted successfully"})
}

// Currency Settings Handlers

func GetCurrencySettings(c echo.Context) error {
	guildID := c.Param("guildId")

	if err := verifyGuildAccess(c, guildID); err != nil {
		return err
	}

	var settings database.CurrencySetup
	result := database.DB.Where("guild_id = ?", guildID).First(&settings)

	if result.Error != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"guild_id": guildID,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"guild_id": settings.GuildID,
		"currency": settings.Currency,
		"side":     settings.Side,
	})
}

func UpdateCurrencySettings(c echo.Context) error {
	guildID := c.Param("guildId")

	if err := verifyGuildAccess(c, guildID); err != nil {
		return err
	}

	var req CurrencySettingsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	settings := database.CurrencySetup{
		GuildID:  guildID,
		Currency: req.Currency,
		Side:     req.Side,
	}

	result := database.DB.Save(&settings)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save settings"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Settings updated successfully"})
}

func DeleteCurrencySettings(c echo.Context) error {
	guildID := c.Param("guildId")

	if err := verifyGuildAccess(c, guildID); err != nil {
		return err
	}

	result := database.DB.Delete(&database.CurrencySetup{}, "guild_id = ?", guildID)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete settings"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Settings deleted successfully"})
}

// Shop Settings Handlers

func GetShopSettings(c echo.Context) error {
	guildID := c.Param("guildId")

	if err := verifyGuildAccess(c, guildID); err != nil {
		return err
	}

	var settings database.ShopSetup
	result := database.DB.Where("guild_id = ?", guildID).First(&settings)

	if result.Error != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"guild_id": guildID,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"guild_id":    settings.GuildID,
		"host_role":   settings.HostRole,
		"alert_role":  settings.AlertRole,
		"log_channel": settings.LogChannel,
	})
}

func UpdateShopSettings(c echo.Context) error {
	guildID := c.Param("guildId")

	if err := verifyGuildAccess(c, guildID); err != nil {
		return err
	}

	var req ShopSettingsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	settings := database.ShopSetup{
		GuildID:    guildID,
		HostRole:   req.HostRole,
		AlertRole:  req.AlertRole,
		LogChannel: req.LogChannel,
	}

	result := database.DB.Save(&settings)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save settings"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Settings updated successfully"})
}

func DeleteShopSettings(c echo.Context) error {
	guildID := c.Param("guildId")

	if err := verifyGuildAccess(c, guildID); err != nil {
		return err
	}

	result := database.DB.Delete(&database.ShopSetup{}, "guild_id = ?", guildID)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete settings"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Settings deleted successfully"})
}
