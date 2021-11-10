package connect

import (
	"github.com/bwmarrin/discordgo"
	c "gitlab.com/logan9312/discord-auction-bot/commands"
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
			c.ErrorResponse(s, i, "Command response has not been set properly, please contact Logan to fix")
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
		default:
			c.ErrorResponse(s, i, "Command response has not been set properly, please contact Logan to fix")
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
