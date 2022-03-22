package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var FeedbackCommand = discordgo.ApplicationCommand{
	Name:        "feedback",
	Description: "Send me feedback for the bot, or report a bug!",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "suggestion",
			Description: "Suggest an improvement for the bot!",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "bug_report",
			Description: "Report a bug!",
		},
	},
}

func Feedback(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Options[0].Name {
	case "suggestion":
		FeedbackSuggestion(s, i)
	}
}

func FeedbackSuggestion(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			TTS:     false,
			Content: "",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "feedback",
							Label:       "Modal Test",
							Style:       discordgo.TextInputParagraph,
							Placeholder: "Placeholder text",
							Value:       "Test Value",
							Required:    false,
							MinLength:   10,
							MaxLength:   1000,
						},
					},
				},
			},
			Embeds:          []*discordgo.MessageEmbed{},
			AllowedMentions: &discordgo.MessageAllowedMentions{},
			Flags:           0,
			Choices:         []*discordgo.ApplicationCommandOptionChoice{},
			CustomID:        "modal_test",
			Title:           "**THis is a modal**",
		},
	})
	if err != nil {
		fmt.Println(err)
	}
}
