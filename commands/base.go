package commands 

import (
	"github.com/bwmarrin/discordgo"
)

type InfoType struct {}

func CommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == 2 {
		switch i.ApplicationCommandData().Name {
		case "help":
			Help(s, i)
		case "auction":
			Auction(s, i, s.State.User.ID)
		case "inventory":
			Profile(s, i)
		case "bid":
			Bid(s, i)
		default:
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds:          []*discordgo.MessageEmbed{
						{
							Title:       "Command Selection Error",
							Description: "Command response has not been set properly, please contact Logan to fix",
						},
					},
					Flags:           64,
				},
			})
		}
		switch i.ApplicationCommandData().Options[0].Name {
		case "create":
			DynamicCreate(s, i)
		case "add":
			DynamicEdit(s, i)
		}
	}
	if i.Type == 3 {
		switch i.MessageComponentData().CustomID {
			//auction bid alert button
		}
	}
}
