package commands

import "github.com/bwmarrin/discordgo"

var RolesCommand = discordgo.ApplicationCommand{
	Name:        "roles",
	Description: "Manage roles in your server",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:         discordgo.ApplicationCommandOptionSubCommand,
			Name:         "panel",
			Description:  "",
			Options:      []*discordgo.ApplicationCommandOption{},
		},
	},
}
