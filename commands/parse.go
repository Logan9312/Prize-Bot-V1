package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func ParseSlashCommand(i *discordgo.InteractionCreate) map[string]interface{} {
	var options map[string]interface{} = make(map[string]interface{})
	for _, option := range i.ApplicationCommandData().Options {
		options[option.Name] = option.Value
	}

	return options
}

func ParseSubCommand(i *discordgo.InteractionCreate) map[string]interface{} {
	var options map[string]interface{} = make(map[string]interface{})
	for _, option := range i.ApplicationCommandData().Options[0].Options {
		options[option.Name] = option.Value
	}

	return options
}

func ErrorResponse(s *discordgo.Session, i *discordgo.InteractionCreate, errorText string) {

	command := "/" + i.ApplicationCommandData().Name

	fmt.Println(errorText)

	if i.ApplicationCommandData().Options[0].Type == discordgo.ApplicationCommandOptionSubCommand {
		command += " " + i.ApplicationCommandData().Options[0].Name
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label: "Support Server",
							Style: discordgo.LinkButton,
							Emoji: discordgo.ComponentEmoji{
								Name:     "logo",
								ID:       "889025400120950804",
								Animated: false,
							},
							URL: "https://discord.gg/RxP2z5NGtj",
						},
					},
				},
			},
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Error",
					Description: fmt.Sprintf("There was an error trying to run `%s`. Please create a support ticket if you need help.", command),
					Color:       0xff0000,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "**Error Message:**",
							Value: errorText,
						},
					},
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "Aftermath Auction Bot",
						IconURL: "https://media.discordapp.net/attachments/859271759920234518/859951197795123261/final_logo.png?width=473&height=473",
					},
				},
			},
			Flags: 64,
		},
	})

	if err != nil {
		fmt.Println(err)
	}
}

func ErrorMessage(s *discordgo.Session, channelID, err string) {

	message := discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title:       "Error",
			Description: "There was an error. Please create a support ticket if you need help.",
			Color:       0xff0000,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "**Error Message:**",
					Value: err,
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Aftermath Auction Bot",
				IconURL: "https://media.discordapp.net/attachments/859271759920234518/859951197795123261/final_logo.png?width=473&height=473",
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label: "Support Server",
						Style: discordgo.LinkButton,
						Emoji: discordgo.ComponentEmoji{
							Name:     "logo",
							ID:       "889025400120950804",
							Animated: false,
						},
						URL: "https://discord.gg/RxP2z5NGtj",
					},
				},
			},
		},
	}
	s.ChannelMessageSendComplex(channelID, &message)
}
