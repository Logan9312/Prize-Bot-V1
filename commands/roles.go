package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var RolesCommand = discordgo.ApplicationCommand{
	Name:          "roles",
	Description:   "Adds Roles",
}

func Roles(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Awesome Role",
							Style:    3,
							Disabled: false,
							Emoji: discordgo.ButtonEmoji{
								Name:     "😁",
								ID:       "",
								Animated: false,
							},
							CustomID: "roles1",
						},
					},
				},
			},
			Embeds: []*discordgo.MessageEmbed{
				{
					URL: "",
					Type: "image",
					Title: "**__Testing Reaction Roles__**",
					Description: "Select your roles here!",
					Timestamp: "",
				},
			},
		},
	})

}

func RolesButton(s *discordgo.Session, i *discordgo.InteractionCreate) {

	err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, "854948227724214282")

	if err != nil {
		fmt.Println(err)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    "Role Added!",
			Components: []discordgo.MessageComponent{},
		},
	})

}
