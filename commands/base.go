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
		case "queue":
			Queue(s, i)
		default:
			DefaultResponse(s, i)
		}
	}
	if i.Type == 3 {
		switch i.MessageComponentData().CustomID {
		case "startbid":
			AuctionButton(s, i)
		case "placebid":
			Bid(s, i)
		default:
			DefaultResponse(s, i)
		}
	}
}

func DefaultResponse (s *discordgo.Session, i *discordgo.InteractionCreate) {
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

func ParseSlashCommand(i *discordgo.InteractionCreate) map[string]interface{} {
    var options map[string]interface{} = make(map[string]interface{})
    for _, option := range i.ApplicationCommandData().Options {
        options[option.Name] = option.Value
    }

    return options
}