package helpers

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/config"
	"gitlab.com/logan9312/discord-auction-bot/logger"
)

// GetSupportButton returns the standard support server button
func GetSupportButton() discordgo.Button {
	return discordgo.Button{
		Label: "Support Server",
		Style: discordgo.LinkButton,
		Emoji: &discordgo.ComponentEmoji{
			Name:     config.C.SupportButtonEmojiName,
			ID:       config.C.SupportButtonEmojiID,
			Animated: false,
		},
		URL: config.C.SupportServerURL,
	}
}

// GetSupportButtonRow returns an action row containing the support button
func GetSupportButtonRow() discordgo.ActionsRow {
	return discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			GetSupportButton(),
		},
	}
}

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
	Reference *discordgo.MessageReference
}

func ErrorResponse(s *discordgo.Session, i *discordgo.InteractionCreate, errorText string) error {
	logger.Sugar.Errorw("error response sent to user", "error", errorText, "guild_id", i.GuildID)

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{
				GetSupportButtonRow(),
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
				},
			},
			Flags: 64,
		},
	})
}

func FollowUpErrorResponse(s *discordgo.Session, i *discordgo.InteractionCreate, errorText string) (*discordgo.Message, error) {
	logger.Sugar.Errorw("followup error response sent to user", "error", errorText, "guild_id", i.GuildID)
	return s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		TTS:   false,
		Files: []*discordgo.File{},
		Components: []discordgo.MessageComponent{
			GetSupportButtonRow(),
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
			},
		},
		Flags: 64,
	})
}

func ErrorMessage(s *discordgo.Session, channelID, err string) (*discordgo.Message, error) {
	logger.Sugar.Errorw("error message sent to channel", "error", err, "channel_id", channelID)
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
		},
		Components: []discordgo.MessageComponent{
			GetSupportButtonRow(),
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

	return s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
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
	return s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content:    r.Content,
		Components: r.Components,
		Embeds:     embed,
		Files:      r.Files,
		Reference:  r.Reference,
	})
}

func SuccessMessageEdit(s *discordgo.Session, channelID, messageID string, r PresetResponse) (*discordgo.Message, error) {

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

	return s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Content:    &r.Content,
		Components: &r.Components,
		Embeds:     &embed,
		ID:         messageID,
		Channel:    channelID,
	})
}

func DeferredResponse(s *discordgo.Session, i *discordgo.InteractionCreate, r PresetResponse) (*discordgo.Message, error) {

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
	return s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content:    &r.Content,
		Components: &r.Components,
		Embeds:     &embed,
		Files:      r.Files,
	})
}

func DeferredErrorResponse(s *discordgo.Session, i *discordgo.InteractionCreate, errorText string) (*discordgo.Message, error) {
	supportRow := GetSupportButtonRow()
	return s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Components: &[]discordgo.MessageComponent{
			supportRow,
		},
		Embeds: &[]*discordgo.MessageEmbed{
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

func PremiumError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) error {

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{},
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "__**Premium Users Only**__",
					Description: "This function is restricted to premium users. Please check out `/premium info` for more information.",
					Color:       0xffd700,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "**Message**",
							Value:  message,
							Inline: false,
						},
					},
				},
			},
			Flags: 64,
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
				GetSupportButtonRow(),
			},
			Embeds: embed,
			Flags:  64,
		},
	})
}
