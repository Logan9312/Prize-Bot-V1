package commands

import (
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
			Name:        "bid",
			Description: "Starting Bid Amount",
			Required:    true,
		},
	},
}

func Auction(s *discordgo.Session, i *discordgo.InteractionCreate) {

	details := i.ApplicationCommandData().Options[0].StringValue()
	initialBid := i.ApplicationCommandData().Options[1].StringValue()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						&discordgo.Button{
							Label:    "Lower your bid",
							Style:    2,
							Disabled: false,
							Emoji:    discordgo.ButtonEmoji{},
							CustomID: "auction1",
						},
						&discordgo.Button{
							Label:    "Bid: " + initialBid + " 🍓",
							Style:    3,
							Disabled: false,
							Emoji: discordgo.ButtonEmoji{
								Name:     "check",
								ID:       "623703744592347146",
								Animated: false,
							},
							Link:     "",
							CustomID: "auction2",
						},
						&discordgo.Button{
							Label:    "Raise your bid",
							Style:    2,
							Disabled: false,
							Emoji:    discordgo.ButtonEmoji{},
							CustomID: "auction3",
						},
					},
				},
			},
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Item: " + details,
					Description: "Current Highest Bid: " + initialBid + " 🍓",
					Timestamp:   "",
					Color:       0,
				},
			},
		},
	})
}

func AuctionButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	content := ""
	switch i.MessageComponentData().CustomID {
	case "auction1":
		content += "Current Bid"
	case "auction2":
		content += "Bid Successful"
	case "auction3":
		content += "Current Bid"

	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    content,
			Components: []discordgo.MessageComponent{},
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Item: Test Item",
					Description: content,
				},
			},
		},
	})

}
