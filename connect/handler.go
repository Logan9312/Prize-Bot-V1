package connect

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	c "gitlab.com/logan9312/discord-auction-bot/commands"
	"gitlab.com/logan9312/discord-auction-bot/database"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
)

func CommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	switch i.Type {
	case 2:
		switch i.ApplicationCommandData().Name {
		case "help":
			c.Help(s, i)
		case "auction":
			c.Auction(s, i)
		case "bid":
			c.AuctionBid(s, i)
		case "profile":
			c.Profile(s, i)
		case "giveaway":
			c.Giveaway(s, i)
		default:
			h.ErrorResponse(s, i, "Command response has not been set properly, please contact Logan to fix")
		}
	case 3:
		switch i.MessageComponentData().CustomID {
		case "endauction":
			c.AuctionEndButton(s, i)
		case "claim_prize":
			c.ClaimPrizeButton(s, i)
		case "clearauction":
			c.ClearAuctionButton(s, i)
		case "delete_auction_queue":
			c.DeleteAuctionQueue(s, i)
		case "enter_giveaway":
			c.GiveawayEnter(s, i)
		case "claim_giveaway":
			c.ClaimGiveawayButton(s, i)
		case "reroll_giveaway":
			c.RerollGiveawayButton(s, i)
		case "clear_auction_setup":
			c.AuctionSetupClearButton(s, i)
		case "clear_giveaway_setup":
			c.GiveawaySetupClearButton(s, i)
		default:
			h.ErrorResponse(s, i, "Command response has not been set properly, please contact Logan to fix")
		}
	case 4:
		switch i.ApplicationCommandData().Name {
		case "auction":
			c.AuctionAutoComplete(s, i)
		case "giveaway":
			c.GiveawayAutoComplete(s, i)
		}
	}
}

func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	bidValue := ""
	mentioned := false
	if m.Content == "" {
		return
	}

	for _, v := range m.Message.Mentions {
		if v.ID == s.State.User.ID {
			mentioned = true
		}
	}

	if !mentioned {
		return
	}

	splitString := strings.Split(m.Content, " ")

	for n, v := range splitString {
		if strings.ToLower(v) == "bid" {
			if len(splitString) <= n+1 {
				h.ErrorMessage(s, m.ChannelID, "Error Bidding. Your message must contain the bid amount after the word bid. Ex: Bid 500")
				return
			}
			bidValue = splitString[n+1]
			bidAmount, err := strconv.ParseFloat(bidValue, 64)
			if err != nil {
				h.ErrorMessage(s, m.ChannelID, err.Error())
				fmt.Println(err)
				return
			}
			c.AuctionBidAlternate(s, database.Auction{
				ChannelID: m.ChannelID,
				Bid:       bidAmount,
				Winner:    m.Author.ID,
				GuildID:   m.GuildID,
			})
		}
	}

}
