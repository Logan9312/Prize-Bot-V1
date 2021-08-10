package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var DynamicVotingCommand = discordgo.ApplicationCommand{
	Name:          "dynamic-voting",
	Description:   "Setup dynamic voting",
	Options:       []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "create",
			Description: "Creates a new Dynamic Stats Voting page",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "add",
			Description: "Adds a new dynamic vote option",
			Required:    false,
			Choices:     []*discordgo.ApplicationCommandOptionChoice{},
			Options:     []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "The name of the stat you want to add",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "The channel where the dynamic vote is.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message-id",
					Description: "The ID of the dynamic voting message",
					Required:    true,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "remove",
			Description: "Removes an entry from Dynamic Stat Voting",
			Required:    false,
			Options:     []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "The name of the stat you want to remove",
					Required:    true,
				},
			},
		},
	},
}

func DynamicCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {

message := discordgo.MessageSend{
	Embed:           &discordgo.MessageEmbed{
		Title:       "__**Dynamic Stats Voting**__",
		Description: "Please vote below using the provided select menu.",
		Footer:      &discordgo.MessageEmbedFooter{
			Text:         "DSV | Sent By: " + i.Member.Nick,
		},
	},
	Components: []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:    "defaultstat",
					Placeholder: "No Stats Added!",
					MinValues:   0,
					MaxValues:   0,
				},
			},
		},
	},
}
s.ChannelMessageSendComplex("837391772736159754", &message)
	
}

func DynamicEdit(s *discordgo.Session, i *discordgo.InteractionCreate){

	var statName, channelID, messageID string

for _, v := range i.ApplicationCommandData().Options[0].Options {


	switch v.Name {
	case "name":
		statName = v.StringValue()
	case "channel":
		channelID = v.ChannelValue(s).ID
	case "message-id":
		messageID = v.StringValue()
	}
}

returnMessage, err := s.ChannelMessage(channelID, messageID)
if err != nil {
	fmt.Println(err)
	return
}

newfield := []*discordgo.MessageEmbedField{
	{
	Name:   statName,
	Value:  "Current Vote: No Vote",
	Inline: false,
	},
}

messageEdit := discordgo.NewMessageEdit(channelID, messageID)
messageEdit.Embed = &discordgo.MessageEmbed{
	Fields:      append(returnMessage.Embeds[0].Fields, newfield[0]),
	Title:       "__**Dynamic Stats Voting**__",
	Description: "Please vote below using the provided select menu.",
	Footer:      &discordgo.MessageEmbedFooter{
		Text:         "DSV | Sent By: " + i.Member.Nick,
	},
}

messageEdit.Components = []discordgo.MessageComponent{
	discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.SelectMenu{
				CustomID:    "defaultstat",
				Placeholder: "No Stats Added!",
				MinValues:   0,
				MaxValues:   0,
			},
		},
	},
}

	_, err = s.ChannelMessageEditComplex(messageEdit)
	if err != nil {
		fmt.Println(err)
		return
	}
	
}