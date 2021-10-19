package commands

import "github.com/bwmarrin/discordgo"

type PresetMessageComplex struct {
	Content         string                            `json:"content,omitempty"`
	Embeds          []*discordgo.MessageEmbed         `json:"embeds,omitempty"`
	TTS             bool                              `json:"tts"`
	Components      []discordgo.MessageComponent      `json:"components"`
	Files           []*discordgo.File                 `json:"-"`
	AllowedMentions *discordgo.MessageAllowedMentions `json:"allowed_mentions,omitempty"`
	Reference       *discordgo.MessageReference       `json:"message_reference,omitempty"`
}

type PresetResponse struct {
	Content     string `json:"content"`
	Title       string
	Description string
	Fields      []*discordgo.MessageEmbedField
	Thumbnail   *discordgo.MessageEmbedThumbnail
	//Components      []discordgo.MessageComponent      `json:"components"`
	//Embeds          []*discordgo.MessageEmbed         `json:"embeds,omitempty"`

	//Flags uint64 `json:"flags,omitempty"`

	//Files []*File `json:"-"`
}

func ErrorResponse(s *discordgo.Session, i *discordgo.InteractionCreate, errorText string) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
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
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Error",
					Description: "There was an error. Please contact support if you need help.",
					Color:       0xff0000,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "**Error Message:**",
							Value: errorText,
						},
					},
					//Footer: &discordgo.MessageEmbedFooter{
					//	Text:    "Aftermath Auction Bot",
					//	IconURL: "https://media.discordapp.net/attachments/859271759920234518/859951197795123261/final_logo.png?width=473&height=473",
					//},
				},
			},
			Flags: 64,
		},
	})
}

func ErrorMessage(s *discordgo.Session, channelID, err string) (*discordgo.Message, error) {

	return s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title:       "Error",
			Description: "There was an error. Please create a support ticket if you need help.",
			Color:       0xff0000,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "**Error Message:**",
					Value: err,
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Aftermath Auction Bot",
				IconURL: "https://media.discordapp.net/attachments/859271759920234518/859951197795123261/final_logo.png?width=473&height=473",
			},
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

func SuccessResponse(s *discordgo.Session, i *discordgo.InteractionCreate, r PresetResponse) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    r.Content,
			Components: []discordgo.MessageComponent{},
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       r.Title,
					Description: r.Description,
					Color:       0x8073ff,
					Footer:      &discordgo.MessageEmbedFooter{},
					Image:       &discordgo.MessageEmbedImage{},
					Thumbnail:   &discordgo.MessageEmbedThumbnail{},
					Video:       &discordgo.MessageEmbedVideo{},
					Author:      &discordgo.MessageEmbedAuthor{},
					Fields:      r.Fields,
				},
			},
			AllowedMentions: &discordgo.MessageAllowedMentions{},
			Flags:           64,
			Files:           []*discordgo.File{},
		},
	})
}
