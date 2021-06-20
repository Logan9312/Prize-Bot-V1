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

var initialBid int64 = 500

func Auction(s *discordgo.Session, i *discordgo.InteractionCreate) {

	details := i.ApplicationCommandData().Options[0].StringValue()
	initialBid := i.ApplicationCommandData().Options[1].IntValue()
	bidMessage := &discordgo.MessageSend{
		Content: "",
		Embed:   &discordgo.MessageEmbed{},
		TTS:     false,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					&discordgo.Button{
						Label:    "Raise Bid by 500",
						Style:    2,
						Disabled: false,
						Emoji:    discordgo.ButtonEmoji{},
						CustomID: "auction1",
					},
					&discordgo.Button{
						Label:    "Raise Bid by 1000: " + fmt.Sprintf("%d", initialBid) + " 🍓",
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
						Label:    "Raise bid by 10000",
						Style:    2,
						Disabled: false,
						Emoji:    discordgo.ButtonEmoji{},
						CustomID: "auction3",
					},
				},
			},
		},

		Files:           []*discordgo.File{},
		AllowedMentions: &discordgo.MessageAllowedMentions{},
		Reference:       &discordgo.MessageReference{},
		File:            &discordgo.File{},
	}

	if i.ApplicationCommandData().Name == "auction" {
		s.ChannelMessageSendComplex(i.ChannelID, bidMessage)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Item: " + details,
						Description: "Current Highest Bid: " + fmt.Sprintf("%d", initialBid) + " 🍓",
						Timestamp:   "",
						Color:       0x8073ff,
					},
				},
			},
		})

		var BidCommand = discordgo.ApplicationCommand{Name: "auction",
			Description: "Bid in an auction!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "bid",
					Description: "Starting Bid Amount",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "display_username",
					Description: "Display your username, or keep it private?",
					Required:    false,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "display_name",
							Value: i.Member.Nick,
						},
						{
							Name:  "hide_name",
							Value: "User has chosen to keep their identity secret.",
						},
					},
					Options: []*discordgo.ApplicationCommandOption{},
				},
			},
		}

		_, err := s.ApplicationCommandCreate(s.State.User.ID, i.GuildID, &BidCommand)
		if err != nil {
			fmt.Println(err)
		}
	}

	if i.ApplicationCommandData().Name == "bid" {

	}
}

func AuctionButtons(s *discordgo.Session, i *discordgo.InteractionCreate) {

	fmt.Println(initialBid)

}
