package responses

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func SuccessResponse(s *discordgo.Session, i *discordgo.InteractionCreate, data *discordgo.InteractionResponseData) error {
	for _, embed := range data.Embeds {
		if embed.Color == 0 {
			embed.Color = ResponseColorPurple
		}
	}
	return respond(s, i, data)
}

func FollowUpSuccessResponse(s *discordgo.Session, i *discordgo.InteractionCreate, data *discordgo.WebhookParams) (*discordgo.Message, error) {
	for _, embed := range data.Embeds {
		if embed.Color == 0 {
			embed.Color = ResponseColorPurple
		}
	}
	return followupRespond(s, i, data)
}

func SuccessMessage(s *discordgo.Session, channelID string, data *discordgo.MessageSend) (*discordgo.Message, error) {

	for _, embed := range data.Embeds {
		if embed.Color == 0 {
			embed.Color = ResponseColorPurple
		}
	}
	return s.ChannelMessageSendComplex(channelID, data)
}

func ExperimentalResponse(s *discordgo.Session, i *discordgo.InteractionCreate, data *discordgo.InteractionResponseData) error {
	if len(data.Embeds) == 0 {
		return Errorfmt(fmt.Errorf("no embeds provided"))
	}
	data.Embeds = append(data.Embeds, &discordgo.MessageEmbed{
		Title:       "**Experimental Feature!**",
		Description: "This feature is highly experimental, and for that reason may not function perfectly as intended. Often problems with this command are more frequent in large servers.",
	})

	for _, embed := range data.Embeds {
		if embed.Color == 0 {
			embed.Color = ResponseColorOrange
		}
	}
	if data.Components == nil {
		data.Components = []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					SupportButton,
				},
			},
		}
	} else {
		data.Components = append(data.Components, discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				SupportButton,
			},
		})
	}

	return respond(s, i, data)
}