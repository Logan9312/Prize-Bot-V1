package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"

)

func Help (s *discordgo.Session, m *discordgo.MessageCreate, details []string) {


	helpmessage := &discordgo.MessageSend{
		Content:         "Test",
		Embed:           &discordgo.MessageEmbed{
			Type:        "",
			Title:       "",
			Description: "",
			Timestamp:   "",
			Color:       0,
			Footer:      &discordgo.MessageEmbedFooter{},
			Image:       &discordgo.MessageEmbedImage{},
			Thumbnail:   &discordgo.MessageEmbedThumbnail{},
			Video:       &discordgo.MessageEmbedVideo{},
			Provider:    &discordgo.MessageEmbedProvider{},
			Author:      &discordgo.MessageEmbedAuthor{},
			Fields:      []*discordgo.MessageEmbedField{},
		},
	}
	_, err := s.ChannelMessageSendComplex(m.ChannelID, helpmessage)
	fmt.Println(err)
}