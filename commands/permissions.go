package commands

import "github.com/bwmarrin/discordgo"

//The command will save perms by role for each category. Setup, create, etc.

var PermissionCommand = discordgo.ApplicationCommand{
	Name:              "permission",
	Description:       "Set the permissions for each role",
	Options:           []*discordgo.ApplicationCommandOption{
		{
			Type:         discordgo.ApplicationCommandOptionRole,
			Name:         "role",
			Description:  "The role you want to set the permissions for.",
			ChannelTypes: []discordgo.ChannelType{},
			Required:     false,
			Options:      []*discordgo.ApplicationCommandOption{
				{
					Type:         0,
					Name:         "",
					Description:  "",
					ChannelTypes: []discordgo.ChannelType{},
					Required:     false,
					Options:      []*discordgo.ApplicationCommandOption{},
					Autocomplete: false,
					Choices:      []*discordgo.ApplicationCommandOptionChoice{},
				},
			},
			Autocomplete: false,
			Choices:      []*discordgo.ApplicationCommandOptionChoice{},
		},
	},
}