package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gitlab.com/logan9312/discord-auction-bot/commands"
	"gitlab.com/logan9312/discord-auction-bot/database"
)

// AuctionListItem represents an auction in the list response
type AuctionListItem struct {
	ChannelID    string    `json:"channel_id"`
	Item         string    `json:"item"`
	Bid          float64   `json:"bid"`
	Winner       string    `json:"winner"`
	Host         string    `json:"host"`
	EndTime      time.Time `json:"end_time"`
	Currency     string    `json:"currency"`
	CurrencySide string    `json:"currency_side"`
}

// GiveawayListItem represents a giveaway in the list response
type GiveawayListItem struct {
	MessageID string    `json:"message_id"`
	ChannelID string    `json:"channel_id"`
	Item      string    `json:"item"`
	EndTime   time.Time `json:"end_time"`
	Host      string    `json:"host"`
	Winners   int       `json:"winners"`
}

// ClaimListItem represents a claim in the list response
type ClaimListItem struct {
	MessageID string  `json:"message_id"`
	ChannelID string  `json:"channel_id"`
	Item      string  `json:"item"`
	Winner    string  `json:"winner"`
	Cost      float64 `json:"cost"`
	Status    string  `json:"status"`
	TicketID  string  `json:"ticket_id"`
}

// ListAuctions returns active auctions for a guild
func ListAuctions(c echo.Context) error {
	guildID := c.Param("guildId")

	if err := verifyGuildAccess(c, guildID); err != nil {
		return err
	}

	var auctions []database.Auction
	result := database.DB.Where("guild_id = ? AND end_time > ?", guildID, time.Now()).
		Order("end_time ASC").
		Find(&auctions)

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch auctions"})
	}

	items := make([]AuctionListItem, len(auctions))
	for i, a := range auctions {
		items[i] = AuctionListItem{
			ChannelID:    a.ChannelID,
			Item:         a.Item,
			Bid:          a.Bid,
			Winner:       a.Winner,
			Host:         a.Host,
			EndTime:      a.EndTime,
			Currency:     a.Currency,
			CurrencySide: a.CurrencySide,
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"auctions": items,
	})
}

// ListGiveaways returns running giveaways for a guild
func ListGiveaways(c echo.Context) error {
	guildID := c.Param("guildId")

	if err := verifyGuildAccess(c, guildID); err != nil {
		return err
	}

	var giveaways []database.Giveaway
	result := database.DB.Where("guild_id = ? AND finished = ?", guildID, false).
		Order("end_time ASC").
		Find(&giveaways)

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch giveaways"})
	}

	items := make([]GiveawayListItem, len(giveaways))
	for i, g := range giveaways {
		items[i] = GiveawayListItem{
			MessageID: g.MessageID,
			ChannelID: g.ChannelID,
			Item:      g.Item,
			EndTime:   g.EndTime,
			Host:      g.Host,
			Winners:   g.Winners,
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"giveaways": items,
	})
}

// ListClaims returns open claims for a guild
func ListClaims(c echo.Context) error {
	guildID := c.Param("guildId")

	if err := verifyGuildAccess(c, guildID); err != nil {
		return err
	}

	var claims []database.Claim
	result := database.DB.Where("guild_id = ? AND status = ?", guildID, database.ClaimStatusPending).
		Find(&claims)

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch claims"})
	}

	items := make([]ClaimListItem, len(claims))
	for i, cl := range claims {
		items[i] = ClaimListItem{
			MessageID: cl.MessageID,
			ChannelID: cl.ChannelID,
			Item:      cl.Item,
			Winner:    cl.Winner,
			Cost:      cl.Cost,
			Status:    cl.Status,
			TicketID:  cl.TicketID,
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"claims": items,
	})
}

// UpdateClaimRequest represents the request body for updating a claim
type UpdateClaimRequest struct {
	Item   string  `json:"item"`
	Winner string  `json:"winner"`
	Cost   float64 `json:"cost"`
}

// UpdateClaim updates a claim's details
func UpdateClaim(c echo.Context) error {
	guildID := c.Param("guildId")
	messageID := c.Param("messageId")

	if err := verifyGuildAccess(c, guildID); err != nil {
		return err
	}

	var req UpdateClaimRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Verify the claim exists and belongs to this guild
	var claim database.Claim
	result := database.DB.Where("message_id = ? AND guild_id = ?", messageID, guildID).First(&claim)
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Claim not found"})
	}

	// Update the claim
	updates := map[string]interface{}{}
	if req.Item != "" {
		updates["item"] = req.Item
	}
	if req.Winner != "" {
		updates["winner"] = req.Winner
	}
	if req.Cost >= 0 {
		updates["cost"] = req.Cost
	}

	if len(updates) > 0 {
		result = database.DB.Model(&claim).Updates(updates)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update claim"})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Claim updated successfully"})
}

// ResendClaim resends a claim message to the log channel
func ResendClaim(c echo.Context) error {
	guildID := c.Param("guildId")
	messageID := c.Param("messageId")

	if err := verifyGuildAccess(c, guildID); err != nil {
		return err
	}

	if BotSession == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Bot session not available"})
	}

	// Look up the claim
	claimMap := map[string]interface{}{}
	result := database.DB.Model(database.Claim{}).First(&claimMap, messageID)
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Claim not found"})
	}

	// Verify guild ownership
	if claimMap["guild_id"] != guildID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Claim does not belong to this guild"})
	}

	// Look up the claim setup to get the log channel
	claimSetup := map[string]interface{}{}
	result = database.DB.Model(database.ClaimSetup{}).First(&claimSetup, guildID)
	if result.Error != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No claim settings found. Please configure claim settings first."})
	}

	logChannel, ok := claimSetup["log_channel"].(string)
	if !ok || logChannel == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No log channel configured. Please set a log channel in claim settings."})
	}

	// Prepare the claim data for resending
	claimMap["log_channel"] = logChannel
	claimMap["old_id"] = claimMap["message_id"]

	eventType := "Auction"
	if claimMap["type"] != nil && claimMap["type"] != "" {
		eventType = claimMap["type"].(string)
	}

	// Resend the claim
	_, err := commands.ClaimOutput(BotSession, claimMap, eventType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to resend claim: %s", err.Error())})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Claim resent successfully"})
}

// CancelClaim cancels a claim by setting its status to cancelled
func CancelClaim(c echo.Context) error {
	guildID := c.Param("guildId")
	messageID := c.Param("messageId")

	if err := verifyGuildAccess(c, guildID); err != nil {
		return err
	}

	// Verify the claim exists and belongs to this guild
	var claim database.Claim
	result := database.DB.Where("message_id = ? AND guild_id = ?", messageID, guildID).First(&claim)
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Claim not found"})
	}

	// Update the status to cancelled
	result = database.DB.Model(&claim).Updates(map[string]interface{}{
		"status":    database.ClaimStatusCancelled,
		"ticket_id": "",
	})
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to cancel claim"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Claim cancelled successfully"})
}
