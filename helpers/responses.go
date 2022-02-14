package helpers

import (
	"fmt"
	"runtime"

	"github.com/bwmarrin/discordgo"
)

/*type PresetMessageComplex struct {
	Content     string `json:"content,omitempty"`
	Title       string
	Description string
	Fields      []*discordgo.MessageEmbedField
	Thumbnail   *discordgo.MessageEmbedThumbnail
	Image       *discordgo.MessageEmbedImage
	Components  []discordgo.MessageComponent `json:"components"`

	//Files           []*discordgo.File                 `json:"-"`
	//AllowedMentions *discordgo.MessageAllowedMentions `json:"allowed_mentions,omitempty"`
	//Reference       *discordgo.MessageReference       `json:"message_reference,omitempty"`
}*/

type PresetResponse struct {
	Content     string `json:"content"`
	Title       string
	Description string
	Fields      []*discordgo.MessageEmbedField
	Thumbnail   *discordgo.MessageEmbedThumbnail
	Image       *discordgo.MessageEmbedImage
	Components  []discordgo.MessageComponent
	Embeds      []*discordgo.MessageEmbed `json:"embeds,omitempty"`
	Files       []*discordgo.File         `json:"-"`
	//Flags uint64 `json:"flags,omitempty"`

}

func ErrorResponse(s *discordgo.Session, i *discordgo.InteractionCreate, errorText string) error {
	_, file, line, _ := runtime.Caller(1)
	fmt.Println(file, line, errorText)
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

func FollowUpErrorResponse(s *discordgo.Session, i *discordgo.InteractionCreate, errorText string) (*discordgo.Message, error) {
	_, file, line, _ := runtime.Caller(1)
	fmt.Println(file, line, errorText)
	return s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
		TTS:   false,
		Files: []*discordgo.File{},
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
	})
}

func ErrorMessage(s *discordgo.Session, channelID, err string) (*discordgo.Message, error) {
	_, file, line, _ := runtime.Caller(1)
	fmt.Println(file, line, err)
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
			//Footer: &discordgo.MessageEmbedFooter{
			//	Text:    "Aftermath Auction Bot",
			//	IconURL: "https://media.discordapp.net/attachments/859271759920234518/859951197795123261/final_logo.png?width=473&height=473",
			//},
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

	embed := []*discordgo.MessageEmbed{
		{
			Title:       r.Title,
			Description: r.Description,
			Color:       0x8073ff,
			Image:       r.Image,
			Thumbnail:   r.Thumbnail,
			Fields:      r.Fields,
		},
	}

	embed = append(embed, r.Embeds...)

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    r.Content,
			Components: r.Components,
			Embeds:     embed,
			Flags:      64,
			Files:      r.Files,
		},
	})
}
func FollowUpSuccessResponse(s *discordgo.Session, i *discordgo.InteractionCreate, r PresetResponse) (*discordgo.Message, error) {

	embed := []*discordgo.MessageEmbed{
		{
			Title:       r.Title,
			Description: r.Description,
			Color:       0x8073ff,
			Image:       r.Image,
			Thumbnail:   r.Thumbnail,
			Fields:      r.Fields,
		},
	}

	embed = append(embed, r.Embeds...)

	return s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
		Content:         r.Content,
		Username:        "",
		AvatarURL:       "",
		TTS:             false,
		Files:           r.Files,
		Components:      r.Components,
		Embeds:          embed,
		AllowedMentions: &discordgo.MessageAllowedMentions{},
		Flags:           64,
	})
}
func SuccessMessage(s *discordgo.Session, channelID string, r PresetResponse) (*discordgo.Message, error) {

	return s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content:    r.Content,
		Components: r.Components,
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       r.Title,
				Description: r.Description,
				Color:       0x8073ff,
				Image:       r.Image,
				Thumbnail:   r.Thumbnail,
				Fields:      r.Fields,
			},
		},
		Files: r.Files,
	})
}

func SuccessMessageEdit(s *discordgo.Session, channelID, messageID string, r PresetResponse) (*discordgo.Message, error) {

	return s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Content:    &r.Content,
		Components: r.Components,
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       r.Title,
				Description: r.Description,
				Color:       0x8073ff,
				Image:       r.Image,
				Thumbnail:   r.Thumbnail,
				Fields:      r.Fields,
			},
		},
		ID:              messageID,
		Channel:         channelID,
	})
}

func DeferredResponse(s *discordgo.Session, i *discordgo.InteractionCreate, r PresetResponse) (*discordgo.Message, error) {
	return s.InteractionResponseEdit(s.State.User.ID, i.Interaction, &discordgo.WebhookEdit{
		Content:    r.Content,
		Components: r.Components,
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       r.Title,
				Description: r.Description,
				Color:       0x8073ff,
				Image:       r.Image,
				Thumbnail:   r.Thumbnail,
				Fields:      r.Fields,
			},
		},
	})
}

func DeferredErrorResponse(s *discordgo.Session, i *discordgo.InteractionCreate, errorText string) (*discordgo.Message, error) {
	return s.InteractionResponseEdit(s.State.User.ID, i.Interaction, &discordgo.WebhookEdit{
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
	})
}

func PresetMessageSend(s *discordgo.Session, channelID string, m PresetResponse) (*discordgo.Message, error) {

	embed := []*discordgo.MessageEmbed{
		{
			Title:       m.Title,
			Description: m.Description,
			Color:       0x8073ff,
			Fields:      m.Fields,
			Thumbnail:   m.Thumbnail,
			Image:       m.Image,
		},
	}

	embed = append(embed, m.Embeds...)

	return s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content:    m.Content,
		Embeds:     embed,
		Components: m.Components,
	})
}

func PremiumResponse(s *discordgo.Session, i *discordgo.InteractionCreate, r PresetResponse) error {

	r.Fields = append(r.Fields, &discordgo.MessageEmbedField{
		Name:   "**Free Premium Feature!**",
		Value:  "Eventually this is gonna be a premium feature, since it's just for convenience and it took more effort to create. For now though, its free for all users!",
		Inline: false,
	})

	embed := []*discordgo.MessageEmbed{
		{
			Title:       r.Title,
			Description: r.Description,
			Color:       0xffd700,
			Footer:      &discordgo.MessageEmbedFooter{},
			Image:       r.Image,
			Thumbnail:   r.Thumbnail,
			Video:       &discordgo.MessageEmbedVideo{},
			Author:      &discordgo.MessageEmbedAuthor{},
			Fields:      r.Fields,
		},
	}

	embed = append(embed, r.Embeds...)

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: r.Content,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label: "Coming Soon",
							Style: discordgo.LinkButton,
							Emoji: discordgo.ComponentEmoji{
								Name:     "logo",
								ID:       "889025400120950804",
								Animated: false,
							},
							URL:      "https://discord.gg/RxP2z5NGtj",
							Disabled: true,
						},
					},
				},
			},
			Embeds: embed,
			Flags:  64,
		},
	})
}

func ExperimentalResponse(s *discordgo.Session, i *discordgo.InteractionCreate, r PresetResponse) error {

	r.Fields = append(r.Fields, &discordgo.MessageEmbedField{
		Name:   "**Experimental Feature!**",
		Value:  "This feature is highly experimental, and for that reason may not function perfectly as intended. Often problems with this command are more frequent in large servers.",
		Inline: false,
	})

	embed := []*discordgo.MessageEmbed{
		{
			Title:       r.Title,
			Description: r.Description,
			Color:       0xff6700,
			Footer:      &discordgo.MessageEmbedFooter{},
			Image:       r.Image,
			Thumbnail:   r.Thumbnail,
			Video:       &discordgo.MessageEmbedVideo{},
			Author:      &discordgo.MessageEmbedAuthor{},
			Fields:      r.Fields,
		},
	}

	embed = append(embed, r.Embeds...)

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: r.Content,
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
			Embeds: embed,
			Flags:  64,
		},
	})
}
