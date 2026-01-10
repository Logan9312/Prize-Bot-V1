package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
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
