package commands

import "github.com/bwmarrin/discordgo"

var InventoryCommand =	discordgo.ApplicationCommand{
	Name:          "inventory",
	Description:   "Displays a user's inventory.",
	Options:       []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "username",
			Description: "Chose who's inventory to display",
			Required:    false,
			Choices:     []*discordgo.ApplicationCommandOptionChoice{},
			Options:     []*discordgo.ApplicationCommandOption{},
		},
	},
}

func Inventory(s *discordgo.Session, i *discordgo.InteractionCreate) {
	
	userID := i.ApplicationCommandData().Options[0].UserValue(nil).ID
	username := i.ApplicationCommandData().Options[0].UserValue(nil).Username

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			TTS:             false,
			Content:         "",
			Components:      []discordgo.MessageComponent{
			},
			Embeds:          []*discordgo.MessageEmbed{
				{
					URL:   "",
					Type:  "",
					Title: "Profile: " + username + "test",
					Description: "Inventory For: " + "<@" + userID + ">",
				},
			},
			AllowedMentions: &discordgo.MessageAllowedMentions{},
		},
	})
}