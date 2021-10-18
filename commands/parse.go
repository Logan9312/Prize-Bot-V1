package commands

import (

	"github.com/bwmarrin/discordgo"
)

func ParseSlashCommand(i *discordgo.InteractionCreate) map[string]interface{} {
	var options map[string]interface{} = make(map[string]interface{})
	for _, option := range i.ApplicationCommandData().Options {
		options[option.Name] = option.Value
	}

	return options
}

func ParseSubCommand(i *discordgo.InteractionCreate) map[string]interface{} {
	var options map[string]interface{} = make(map[string]interface{})
	for _, option := range i.ApplicationCommandData().Options[0].Options {
		options[option.Name] = option.Value
	}

	return options
}

