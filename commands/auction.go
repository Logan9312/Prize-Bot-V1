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
			Name:        "startingbid",
			Description: "Starting Bid Amount",
			Required:    true,
		},
	},
}

var initialBid int = 500

var AuctionData *discordgo.Interaction
var Details string
var InitialBid int64

func Auction(s *discordgo.Session, i *discordgo.InteractionCreate, AppID string) {

	AuctionData = i.Interaction

	Details := i.ApplicationCommandData().Options[0].StringValue()
	InitialBid := i.ApplicationCommandData().Options[1].IntValue()
	bidder := i.Interaction.Member.User.Username

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: 4,
		Data: &discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{},
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Item: " + Details,
					Description: "Current Highest Bid: " + fmt.Sprint(InitialBid) + " 🍓",
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

	newbid := i.ApplicationCommandData().Options[0].IntValue()
	bidder := i.Interaction.Member.User.Username

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: 4,
		Data: &discordgo.InteractionResponseData{
			Content: "Successful Bid",
			Flags:   64,
		},
	})

	bidupdate := discordgo.WebhookEdit{
		Content: "",
		Embeds: []*discordgo.MessageEmbed{
			{Title: "Item: " + Details,
				Description: "Current Highest Bid: " + fmt.Sprint(newbid) + " 🍓",
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
	}

	s.InteractionResponseEdit(appID, AuctionData, &bidupdate)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: 4,
		Data: &discordgo.InteractionResponseData{
			Content: "Bid failed. Please select a larger value",
		},
	})

}

func AuctionButtons(s *discordgo.Session, i *discordgo.InteractionCreate) {

	fmt.Println(initialBid)

}
