package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var HelpCommand =	discordgo.ApplicationCommand{
	Name:          "help",
	Description:   "Basic bot functionality",
}

var fields []*discordgo.MessageEmbedField

func HelpBuilder (slashCommands []*discordgo.ApplicationCommand) {
for _, command := range slashCommands {
	field := &discordgo.MessageEmbedField{
	  Name: "/" + strings.Title(fmt.Sprintf("**%s**", command.Name)),
	  Value: fmt.Sprintf("```%s```", command.Description),
	}
  
	fields = append(fields, field)
  }
}

func Help(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Close Help menu",
							Style:    3,
							Disabled: false,
							Emoji: discordgo.ButtonEmoji{
								Name:     "",
								ID:       "",
								Animated: false,
							},
							CustomID: "Help",
						},
					},
				},
			},
			Embeds: []*discordgo.MessageEmbed{
				{
					URL:         "",
					Type:        "",
					Title:       "Discord Bot Help",
					Description: "With slash commands you can use any of the following features. In the future I will be adding Inventory tracking and other related commands.",
					Color:       0x8073ff,
					Fields: fields,
				},
			},
			AllowedMentions: &discordgo.MessageAllowedMentions{},
		},
	})

}

func HelpButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.ChannelMessageDelete(i.ChannelID, i.Message.ID)

	if err != nil {
		fmt.Println(err)
	}
}


//Help Command Plans.

// Use slashCommands to create an auto updating help command. 