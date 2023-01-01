package responses

import "github.com/bwmarrin/discordgo"

const (
	ResponseColorRed    = 0xff0000
	ResponseColorPurple = 0x8073ff
	ResponseColorOrange = 0xff6700
	ResponseColorYellow = 0xffd700
)

var SupportButton = discordgo.Button{
	Label: "Support Server",
	Style: discordgo.LinkButton,
	Emoji: discordgo.ComponentEmoji{
		Name:     "logo",
		ID:       "889025400120950804",
		Animated: false,
	},
	URL: "https://discord.gg/RxP2z5NGtj",
}

func respond(s *discordgo.Session, i *discordgo.InteractionCreate, resp *discordgo.InteractionResponseData) error {
	resp.Flags = 64
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: resp,
	})
}

func followupRespond(s *discordgo.Session, i *discordgo.InteractionCreate, resp *discordgo.WebhookParams) (*discordgo.Message, error) {
	resp.Flags = 64
	return s.FollowupMessageCreate(i.Interaction, true, resp)
}
