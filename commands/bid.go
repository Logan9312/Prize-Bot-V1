package commands

import (

	"github.com/bwmarrin/discordgo"
)

var BidCommand = discordgo.ApplicationCommand{
	Name:        "bid",
	Description: "Bid on an Auction!",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "bidprice",
			Description: "Starting Bid Amount",
			Required:    true,
		},
	},
}

func Bid (s *discordgo.Session, i *discordgo.InteractionCreate) {

}