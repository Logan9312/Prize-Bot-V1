package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/database"
)

var CurrencyCommand = &discordgo.ApplicationCommand{
	Name:        "currency",
	Description: "Shows the current currency of the bot",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "add",
			Description: "Add Currency to a user",
			Options:     []*discordgo.ApplicationCommandOption{},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "set",
			Description: "Set the exact amount of currency for a user.",
			Options:     []*discordgo.ApplicationCommandOption{},
		},
	},
}

func Currency(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	switch i.ApplicationCommandData().Options[0].Name {
	case "":

	}
	return fmt.Errorf("Unknown Currency command, please contact support")
}

func PriceFormat(price float64, guildID string, override interface{}) string {

	currencyMap := map[string]interface{}{}
	priceString := strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", price), "0"), ".")

	result := database.DB.Model(database.CurrencySetup{}).First(&currencyMap, guildID)
	if result.Error != nil {
		fmt.Println("Error getting currency setup: " + result.Error.Error())
	}

	if override != nil {
		currencyMap["currency"] = override
	}

	if currencyMap["currency"] == nil {
		return priceString
	} else {
		if currencyMap["side"] == "right" {
			return fmt.Sprintf("%s %s", priceString, currencyMap["currency"])
		} else {
			return fmt.Sprintf("%s %s", currencyMap["currency"], priceString)
		}
	}

}
