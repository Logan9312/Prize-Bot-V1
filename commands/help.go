package commands

import (
	"github.com/bwmarrin/discordgo"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
)

var HelpCommand = discordgo.ApplicationCommand{
	Name:        "help",
	Description: "Shows the list of commands.",
}

func Help(s *discordgo.Session, i *discordgo.InteractionCreate) {
	h.SuccessResponse(s, i, h.PresetResponse{
		Content:     "",
		Title:       "",
		Description: "",
		Fields:      []*discordgo.MessageEmbedField{},
		Thumbnail:   &discordgo.MessageEmbedThumbnail{},
		Image:       &discordgo.MessageEmbedImage{},
		Components:  []discordgo.MessageComponent{},
		Embeds:      []*discordgo.MessageEmbed{},
		Files:       []*discordgo.File{},
	})
}
