package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
)

var HelpCommand = discordgo.ApplicationCommand{
	Name:        "help",
	Description: "Basic bot functionality",
}

func Help(s *discordgo.Session, i *discordgo.InteractionCreate, slashCommands []*discordgo.ApplicationCommand) {

	fields := []*discordgo.MessageEmbedField{}

	for _, command := range slashCommands {

		if command.Name == "help" {
			continue
		}

		field := &discordgo.MessageEmbedField{
			Name:  "/" + strings.Title(fmt.Sprintf("**%s**", command.Name)),
			Value: fmt.Sprintf("```%s```", command.Description),
		}

		fields = append(fields, field)
	}

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:  "**Ping**",
		Value: s.HeartbeatLatency().String(),
	})

	err := h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "Discord Bot Help",
		Description: "Developed by Logan. Thank you for using my bot!",
		Fields:      fields,
	})

	if err != nil {
		fmt.Println(err)
	}
}
