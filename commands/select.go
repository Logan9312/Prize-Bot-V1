package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var SelectCommand = discordgo.ApplicationCommand{
	Name:        "select-test",
	Description: "Testing Selects",
}

func Select(s *discordgo.Session, i *discordgo.InteractionCreate) {

err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	Type: discordgo.InteractionResponseChannelMessageWithSource,
	Data: &discordgo.InteractionResponseData{
		Components:      []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:    "test",
					Placeholder: "Test Menu",
					MinValues:   1,
					MaxValues:   1,
					Options:     []discordgo.SelectMenuOption{
						{
							Label:       "Option1",
							Value:       "one",
							Description: "This is Option 1",
							Emoji:       discordgo.ComponentEmoji{
								Name:     "😁",
								ID:       "",
								Animated: false,
							},
							Default:     false,
						},
						{
							Label:       "Option2",
							Value:       "two",
							Description: "This is Option 2",
							Emoji:       discordgo.ComponentEmoji{
								Name:     "😁",
								ID:       "",
								Animated: false,
							},
							Default:     false,
						},
					},
				},
			},
		},
		},
		Embeds:          []*discordgo.MessageEmbed{
			{
				Title:       "Select Test",
			},
		},
		AllowedMentions: &discordgo.MessageAllowedMentions{},
		Flags:           0,
	},
})

fmt.Println(err)

}