package commands

import (
	"github.com/bwmarrin/discordgo"
)

var SelectCommand = discordgo.ApplicationCommand{
	Name:        "select-test",
	Description: "Testing Selects",
}

func Select(s *discordgo.Session, i *discordgo.InteractionCreate) {
s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	Type: discordgo.InteractionResponseChannelMessageWithSource,
	Data: &discordgo.InteractionResponseData{
		Content:         "This is testing selects",
		Components:      []discordgo.MessageComponent{
		discordgo.ActionsRow{},
		discordgo.ActionsRow{},
		discordgo.ActionsRow{},
		discordgo.ActionsRow{},
		discordgo.ActionsRow{},
		},
		Embeds:          []*discordgo.MessageEmbed{},
		AllowedMentions: &discordgo.MessageAllowedMentions{},
		Flags:           0,
	},
})
}