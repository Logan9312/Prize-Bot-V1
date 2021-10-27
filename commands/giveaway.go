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
			Required:    true,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "The channel where the giveaway will take place",
					Required:    true,
					ChannelTypes: []discordgo.ChannelType{
						discordgo.ChannelTypeGuildText,
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "channel",
					Description: "The channel where the giveaway will take place",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "winners",
					Description: "The number of winners.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "time",
					Description: "The amount of time that the giveaway will run for. (Example: 1d, 12h, etc.)",
					Required:    true,
				},
			},
		},
	},
}

func Giveaway(s *discordgo.Session, m *discordgo.MessageCreate) {

}
