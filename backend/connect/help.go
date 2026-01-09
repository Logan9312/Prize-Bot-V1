package connect

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	c "gitlab.com/logan9312/discord-auction-bot/commands"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
	"gitlab.com/logan9312/discord-auction-bot/logger"
)

type CommandInfo struct {
	Description    string
	GettingStarted string
	Tips           string
	Settings       string
}

var HelpCommand = discordgo.ApplicationCommand{
	Name:        "help",
	Description: "Shows the list of commands.",
}

var PermissionNames = map[string]int64{
	"Administrator": discordgo.PermissionAdministrator,
	"Manage Server": discordgo.PermissionManageServer,
}

func Help(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	fields := []*discordgo.MessageEmbedField{}

	for _, c := range BotCommands.Prod {
		logger.Sugar.Debugw("loading help for command", "name", c.Name, "permissions", c.DefaultMemberPermissions)
		perms := ParsePerms(c.DefaultMemberPermissions)
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("**%s**", strings.Title(c.Name)),
			Value: fmt.Sprintf("**Description**: %s\n**Default Permissions:** %s\n\u200b", c.Description, perms),
		})
	}

	return h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "Auction Bot Help",
		Description: "Please select a command below for more details",
		Fields:      fields,
		Components: []discordgo.MessageComponent{
			&discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID:    "helpmenu",
						Placeholder: "Select A Command",
						MinValues:   h.Ptr(1),
						MaxValues:   1,
						Options: []discordgo.SelectMenuOption{
							{
								Label:       "Auctions",
								Value:       c.EventTypeAuction,
								Description: "Auction off prizes!",
								Emoji: &discordgo.ComponentEmoji{
									Name: "💸",
								},
							},
						},
					},
				},
			},
		},
	})
}

func HelpMenu(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	name := i.MessageComponentData().Values[0]
	info := CommandInfo{}
	switch name {
	case c.EventTypeAuction:
		info = CommandInfo{
			Description:    "Auctions work by creating a new channel in the selected Category, and self-destructing when completed. You can modify this behaviour in `/settings auctions` by setting a channel override.",
			GettingStarted: ".",
			Tips:           ".",
			Settings:       ".",
		}
	}
	return h.SuccessResponse(s, i, h.PresetResponse{
		Title:       fmt.Sprintf("%s Help", name),
		Description: info.Description,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "**Getting Started**",
				Value: info.GettingStarted,
			},
			{
				Name:  "**Settings**",
				Value: info.Settings,
			},
			{
				Name:  "**Tips**",
				Value: info.Tips,
			},
		},
	})
}

func ParsePerms(permission *int64) (list string) {
	if permission == nil {
		list += "\n\u3000- Everyone"
	} else {
		for name, value := range PermissionNames {
			if *permission&value != 0 {
				list += fmt.Sprintf("\n\u3000- %s", name)
			}
		}
	}
	return
}
