package commands

import "github.com/bwmarrin/discordgo"

var ShopCommand = discordgo.ApplicationCommand{
	Name:        "shop",
	Description: "Buy, Sell and Trade items with others!",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "setup",
			Description: "Setup shop settings on your server",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "sell",
			Description: "Select an item to sell",
			Options:     []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "Item",
					Description: "The item you wish to sell",
					Required:    false,
					Choices:     []*discordgo.ApplicationCommandOptionChoice{},
					Options:     []*discordgo.ApplicationCommandOption{},
				},
			},
		},
	},
}


