package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
)

var HelpCommand = discordgo.ApplicationCommand{
	Name:        "help",
	Description: "Shows the list of commands.",
}

var PermissionNames = map[string]int64{
	"Administrator": discordgo.PermissionAdministrator,
}

func Help(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	commands, err := s.ApplicationCommands(s.State.User.ID, i.GuildID)
	if err != nil {
		return fmt.Errorf("Error fetching command list: %w", err)
	}

	fields := []*discordgo.MessageEmbedField{}

	for _, c := range commands {
		fmt.Println(c.Name, c.DefaultMemberPermissions)
		perms := ParsePerms(c.DefaultMemberPermissions)
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("**%s**", strings.Title(c.Name)),
			Value:  fmt.Sprintf("**Description**: %s\n**Default Permissions:** %s\n\u200b", c.Description, perms),
			Inline: true,
		})
	}

	return h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "Auction Bot Help",
		Description: "Please select a command below for more details",
		Fields:      fields,
	})
}

func ParsePerms(permission *int64) (list string) {
	if permission == nil {
		list += "\n\u3000- Everyone"
	} else {
		for name, value := range PermissionNames {
			if *permission&value != 0 {
				list += fmt.Sprintf("\n\u3000- %s", name)
			}
		}
	}
	return
}
