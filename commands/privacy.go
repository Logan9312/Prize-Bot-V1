package commands

import (
	"github.com/bwmarrin/discordgo"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
)

var PrivacyCommand = discordgo.ApplicationCommand{
	Name:        "privacy_policy",
	Description: "The privacy policy of the bot.",
}

func Privacy(s *discordgo.Session, i *discordgo.InteractionCreate) {
	h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "__**Privacy Policy**__",
		Description: "By using the AFTM Prize Manager bot, you agree to this policy. All collected data is stored in an RDS database on AWS to ensure safety.",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "**POLICY**",
				Value: "The bots will be collecting the following information and only when necessary:\n• User ID (ex: your ID will be saved as the winner if you win an auction).\n• The content of your messages posted in an auction channel.\n• Your Server ID to track auction settings.\n• Any data you send through an interaction with the bot (ex: bid value).\n• The ID of any roles that the bot needs to mention/use.\n•Any other data that is necessary for the bot's functionality can be stored.",
			},
			{
				Name:  "**PRIVACY**",
				Value: "If you feel that you have private data stored on the bot that needs to be deleted, please message me on discord at <@280812467775471627> (Logan#3088), or join my support server for help: [HERE](https://discord.gg/RxP2z5NGtj)",
			},
			{
				Name:  "**CHANGES**",
				Value: "This policy is subject to change at any time, and I will keep it updated whenever needed to ensure the safety and fairness of the data being used. Please contact me if you have any issues with the policy, or the bot.",
			},
			{
				Name:  "**THIRD PARTY POLICIES**",
				Value: "Anything used in conjunction with the bot may have their own privacy policies, and it is encouraged that you read through them before using their services.",
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: s.State.User.AvatarURL(""),
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label: "Support Server",
						Style: discordgo.LinkButton,
						Emoji: discordgo.ComponentEmoji{
							Name:     "logo",
							ID:       "889025400120950804",
							Animated: false,
						},
						URL: "https://discord.gg/RxP2z5NGtj",
					},
				},
			},
		},
	})
}
