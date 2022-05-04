package helpers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

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

func ParseSubSubCommand(i *discordgo.InteractionCreate) map[string]interface{} {
	var options map[string]interface{} = make(map[string]interface{})
	for _, option := range i.ApplicationCommandData().Options[0].Options[0].Options {
		options[option.Name] = option.Value
	}

	return options
}

func ParseTime(inputDuration string) (time.Duration, error) {
	if strings.HasSuffix(strings.ToLower(inputDuration), "d") {
		inputDuration = strings.TrimSuffix(strings.ToLower(inputDuration), "d")
		float, err := strconv.ParseFloat(inputDuration, 64)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		inputDuration = fmt.Sprint(float*24) + "h"
	}

	return time.ParseDuration(strings.ToLower(inputDuration))
}

func Ptr[T any](v T) *T {
	return &v
}
