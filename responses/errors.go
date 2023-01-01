package responses

import (
	"fmt"
	"runtime"

	"github.com/bwmarrin/discordgo"
)

func Errorfmt(err error) error {
	_, file, line, _ := runtime.Caller(1)
	return fmt.Errorf("%s:%d %w", file, line, err)
}

func ErrorComponentfmt(err error) []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				SupportButton,
			},
		},
	}
}

func ErrorEmbedfmt(err error) []*discordgo.MessageEmbed {
	return []*discordgo.MessageEmbed{
		{
			Title:       "Error",
			Description: "There was an error. Please contact support if you need help.",
			Color:       ResponseColorRed,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "**Error Message:**",
					Value: err.Error(),
				},
			},
		},
	}
}

func ErrorResponse(s *discordgo.Session, i *discordgo.InteractionCreate, err error) error {

	fmt.Println(Errorfmt(err).Error())

	return respond(s, i, &discordgo.InteractionResponseData{
		Components: ErrorComponentfmt(err),
		Embeds:     ErrorEmbedfmt(err),
		Flags:      64,
	})
}

func FollowUpErrorResponse(s *discordgo.Session, i *discordgo.InteractionCreate, err error) (*discordgo.Message, error) {

	fmt.Println(Errorfmt(err).Error())

	return followupRespond(s, i, &discordgo.WebhookParams{
		Components: ErrorComponentfmt(err),
		Embeds:     ErrorEmbedfmt(err),
	})
}

func ErrorMessage(s *discordgo.Session, channelID string, err error) (*discordgo.Message, error) {

	fmt.Println(Errorfmt(err).Error())

	return s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Embeds:     ErrorEmbedfmt(err),
		Components: ErrorComponentfmt(err),
	})
}

func PremiumError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) error {

	return respond(s, i, &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "__**Premium Users Only**__",
				Description: "This function is restricted to premium users. Please check out `/premium info` for more information.",
				Color:       ResponseColorYellow,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "**Message**",
						Value:  message,
						Inline: false,
					},
				},
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: "https://i.imgur.com/SjJT2xp.png",
				},
			},
		},
	})
}
