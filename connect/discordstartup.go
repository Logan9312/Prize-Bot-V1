package connect

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var channelID = "863122466485174325"

var message = discordgo.MessageSend{
	Embed:           &discordgo.MessageEmbed{
		Title:        "Yeah but bot can send an embed and you can't",
	},
	Components:      []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "haha loser",
					Style:    2,
					Disabled: true,
					CustomID: "haha",
				},	
			},
		},
	},
	Reference:       &discordgo.MessageReference{
		MessageID: "874820261772394498",
		ChannelID: channelID,
		GuildID:   "863122466485174323",
	},
}

func Random(s *discordgo.Session) {
	if s.State.User.ID == "noID" {
		
		_, err := s.ChannelMessageSendComplex(channelID, &message)
		if err != nil {
			fmt.Println(err)
		}
	}
}