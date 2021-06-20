package commands

import (
	"github.com/bwmarrin/discordgo"
)


var InventoryCommand =	discordgo.ApplicationCommand{
	Name:          "inventory",
	Description:   "Displays a user's inventory.",
	Options:       []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "username",
			Description: "Chose who's inventory to display",
			Required:    true,
			Choices:     []*discordgo.ApplicationCommandOptionChoice{},
			Options:     []*discordgo.ApplicationCommandOption{},
		},
	},
}

func Inventory(s *discordgo.Session, i *discordgo.InteractionCreate) {

	userID := i.ApplicationCommandData().Options[0].UserValue(s).ID
	username := i.ApplicationCommandData().Options[0].UserValue(s).Username
	Discriminator := i.ApplicationCommandData().Options[0].UserValue(s).Discriminator

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			TTS:             false,
			Content:         "",
			Components:      []discordgo.MessageComponent{
			},
			Embeds:          []*discordgo.MessageEmbed{
				{
					Title: "**__" + username + "__**" + "#" + Discriminator,
					Description: "Inventory For: " + "<@" + userID + ">",
					Timestamp: "",
					Color: 0x8073ff,
					Fields:    []*discordgo.MessageEmbedField{
						{
							Name:   "**Super Breeding Ticket**" + " x1",
							Value:  "Use this ticket to breed any dinos in Ark, including non-breedables",
							Inline: false,
						},
						{
							Name:   "**Auction Prize: Max dino of Choice**" + " x4",
							Value:  "Claim this to spawn any max dino",
							Inline: false,
						},
					},
				},
			},
			AllowedMentions: &discordgo.MessageAllowedMentions{},
		},
	})
}