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

func (info InfoType) Help(s *discordgo.Session, i *discordgo.InteractionCreate) {

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: 64,
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "Discord Bot Help",
					Description: "Developed by Logan. Thank you for using my bot!",
					Color:     0x8073ff,
					Footer:    &discordgo.MessageEmbedFooter{
						Text:         "Powered by Aftermath Gaming",
					},
					Fields:    fields,
				},
			},
		},
	})

}