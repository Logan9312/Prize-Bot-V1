package commands

import (
	"fmt"
	"strings"
	"github.com/bwmarrin/discordgo"
)

var fields []*discordgo.MessageEmbedField

var HelpCommand = discordgo.ApplicationCommand{
	Name:        "help",
	Description: "Basic bot functionality",
}

func HelpBuilder (slashCommands []*discordgo.ApplicationCommand) {
for _, command := range slashCommands {
	
	if command.Name == "help" {
		continue
	}

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
							Emoji:    discordgo.ComponentEmoji{
									Name:     "",
									ID:       "",
									Animated: false,
								},
							URL:      "",
							CustomID: "Help",
						},
					},
				},
			},
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "Discord Bot Help",
					Description: "Developed by Logan. Thank you for using my bot!",
					Color:     0x8073ff,
					Footer:    &discordgo.MessageEmbedFooter{
						Text:         "Powered by Aftermath Gaming",
					},
					Image:     &discordgo.MessageEmbedImage{},
					Thumbnail: &discordgo.MessageEmbedThumbnail{},
					Fields:    fields,
				},
			},
		},
	})

}

func HelpButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	
	err := s.ChannelMessageDelete(i.ChannelID, i.Message.ID)

	if err != nil {
		fmt.Println(err)
	}
}