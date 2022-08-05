package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
)

var StatsCommand = discordgo.ApplicationCommand{
	Type:        0,
	Name:        "stats",
	Description: "Statistics",
}

func Stats(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	members := 0
	largestServerCount := 0
	largestServerName := ""

	for _, v := range s.State.Guilds {
		members += v.MemberCount
		if v.MemberCount > largestServerCount {
			largestServerCount = v.MemberCount
			largestServerName = v.Name
		}
	}

	

	h.SuccessResponse(s, i, h.PresetResponse{
		Title: "Statistics",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Members",
				Value:  fmt.Sprint(members),
				Inline: false,
			},
			{
				Name:   "Servers",
				Value:  fmt.Sprint(len(s.State.Guilds)),
				Inline: false,
			},
			{
				Name:   largestServerName,
				Value:  fmt.Sprint(largestServerCount),
				Inline: false,
			},
		},
	})

	return nil
}
