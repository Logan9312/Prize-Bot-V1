package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var AuctionCommand = discordgo.ApplicationCommand{
	Name:        "auction",
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
			Name:        "startingbid",
			Description: "Starting Bid Amount",
			Required:    true,
		},
	},
}

func Auction(s *discordgo.Session, i *discordgo.InteractionCreate, AppID string) {

	Details := i.ApplicationCommandData().Options[0].StringValue()
	initialBid := i.ApplicationCommandData().Options[1].IntValue()
	bidder := i.Interaction.Member.User.Username

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: 4,
		Data: &discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{},
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Item: " + Details,
					Description: "Current Highest Bid: " + fmt.Sprint(initialBid) + " 🍓",
					Color:       0x00bfff,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "**Current Winner:**",
							Value:  bidder,
							Inline: false,
						},
					},
				},
			},
		},
	})

	var BidCommand = discordgo.ApplicationCommand{
		Name:        "bidtest",
		Description: "bid on an existing auction!",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "bidprice",
				Description: "Choose how much you would like to bid",
				Required:    true,
			},
		},
	}

	s.ApplicationCommandCreate(AppID, i.GuildID, &BidCommand)

}

func BidTest(s *discordgo.Session, i *discordgo.InteractionCreate, appID string) {

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: 4,
		Data: &discordgo.InteractionResponseData{
			Content: "Successful Bid",
			Flags:   64,
		},
	})
}
