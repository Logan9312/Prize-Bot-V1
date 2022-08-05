package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
)

var QuestCommand = discordgo.ApplicationCommand{
	Name:        "quest",
	Description: "Quests",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "create",
			Description: "Create a Quest",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "item",
					Description: "The item you wish to auction off",
					Required:    true,
				},
			},
		},
	},
}

var QuestContextMenu = discordgo.ApplicationCommand{
	Type:                     discordgo.MessageApplicationCommand,
	Name:                     "Add New Object",
}

func Quest(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	switch i.ApplicationCommandData().Options[0].Name {
	case "create":
		return QuestCreate(s, i)
	}
	return fmt.Errorf("Unknown Auction command, please contact support")
}

func QuestCreate(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	err := h.SuccessResponse(s, i, h.PresetResponse{
		Content:     "",
		Title:       "QUEST",
		Description: "Claim your quest here",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "STARTER KIT",
				Value:  "REWARD: starter kit",
				Inline: false,
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{},
		Image:     &discordgo.MessageEmbedImage{},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Claim",
						Style:    discordgo.SuccessButton,
						Disabled: false,
						Emoji: discordgo.ComponentEmoji{
							Name: "🐛",
						},
						CustomID: "questbutton2",
					},
					discordgo.Button{
						Label:    "Add New Response",
						Style:    discordgo.SecondaryButton,
						Disabled: false,
						CustomID: "questbutton",
					},
				},
			},
		},
	})
	return err
}

func QuestClaim(s *discordgo.Session, i *discordgo.InteractionCreate) {

}

func QuestButton(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "modal1",
							Label:       "Description",
							Style:       discordgo.TextInputParagraph,
							Placeholder: "Test Placeholder",
							Value:       "Test Value",
							Required:    true,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "modal",
							Label:       "Vault Code",
							Placeholder: "Test Placeholder",
							Value:       "Test Value",
							Required:    true,
							Style:       discordgo.TextInputShort,
						},
					},
				},
			},
			CustomID: "modalsubmit",
			Title:    "Enter The Information",
		},
	})
	return err
}
