package connect

import (
	"github.com/bwmarrin/discordgo"
	c "gitlab.com/logan9312/discord-auction-bot/commands"
)

func CommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == 2 {
		switch i.ApplicationCommandData().Name {
		case "help":
			c.Help(s, i)
		case "auction":
			c.Auction(s, i)
		case "bid":
			c.AuctionBid(s, i)
		case "profile":
			c.Profile(s, i)
		case "queue":
			c.Queue(s, i)
		default:
			c.ErrorResponse(s, i, "Command response has not been set properly, please contact Logan to fix")
		}
	}
	if i.Type == 3 {
		switch i.MessageComponentData().CustomID {
		case "endauction":
			c.AuctionEndButton(s, i)
		case "claim_prize":
			c.ClaimPrizeButton(s, i)
		case "clearauction":
			c.ClearAuctionButton(s, i)
		case "delete_auction_queue":
			c.DeleteAuctionQueue(s, i)
		default:
			c.ErrorResponse(s, i, "Command response has not been set properly, please contact Logan to fix")
		}
	}
}
