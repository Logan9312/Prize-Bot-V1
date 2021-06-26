package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

var AuctionCommand = discordgo.ApplicationCommand{Name: "auction",
	Description: "Put an item up for auction!",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "item",
			Description: "Choose an Item to put up for auction",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "initial_bid",
			Description: "Starting Bid Amount",
			Required:    true,
		},
	},
}

var initialBid int = 0

func Auction(s *discordgo.Session, i *discordgo.InteractionCreate) {

	details := i.ApplicationCommandData().Options[0].StringValue()
	initialBid := int(i.ApplicationCommandData().Options[1].IntValue())

	bidMessage := &discordgo.MessageSend{
		Content: "",
		Embed:  &discordgo.MessageEmbed{
				Title:       "Item: " + details,
				Description: "Current Highest Bid: " + fmt.Sprint(initialBid) + " 🍓",
				Color:       0x00bfff,
		},
	}

		fmt.Println("About to send message")

		_, err := s.ChannelMessageSendComplex(i.ChannelID, bidMessage)

		fmt.Println("Message Sent")

		if err != nil {
		fmt.Println(err)
		}

}

func AuctionButtons(s *discordgo.Session, i *discordgo.InteractionCreate) {

	fmt.Println(initialBid)

}
