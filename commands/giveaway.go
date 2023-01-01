package commands

import "github.com/bwmarrin/discordgo"

var GiveawayCommand = discordgo.ApplicationCommand{
	Name:        "giveaway",
	Description: "Start a giveaway!",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "create",
			Description: "Start a giveaway!",
			Options: []*discordgo.ApplicationCommandOption{
				Require(CommandOptionItem),
				Require(CommandOptionWinner),
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "duration",
					Description:  "The amount of time that the giveaway will run for. (Example: 1d, 12h, etc.)",
					Required:     true,
					Autocomplete: true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "The channel where the giveaway will take place",
					Required:    false,
					ChannelTypes: []discordgo.ChannelType{
						discordgo.ChannelTypeGuildText,
						5,
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "description",
					Description: "Description of the giveaway item.",
					Required:    false,
				},
				{
					Type:        11,
					Name:        "image",
					Description: "Attach an image to your giveaway",
					Required:    false,
				},
			},
		},
	},
}
