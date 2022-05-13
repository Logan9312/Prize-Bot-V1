package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/database"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
)

var CurrencyCreateRolesChunk = []map[string]interface{}{}

var CurrencyCommand = &discordgo.ApplicationCommand{
	Name:        "currency",
	Description: "Shows the current currency of the bot",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "add",
			Description: "Add Currency to a user or role.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionMentionable,
					Name:        "select",
					Description: "Select a user or a role* (premium) add currency to.",
				},
			},
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
	case "add":
		return CurrencyAdd(s, i)
	case "set":
		//return CurrencySet(s, i)
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

func CurrencyAdd(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	if !CheckPremiumGuild(i.GuildID) {
		err := h.PremiumError(s, i)
		if err != nil {
			fmt.Println(err)
		}
		return nil
	}

	currencyMap := h.ParseSubCommand(i)
	currencyMap["guild_id"] = i.GuildID

	if i.ApplicationCommandData().Resolved.Roles[currencyMap["select"].(string)] != nil {
		for _, role := range i.ApplicationCommandData().Resolved.Roles {
			currencyMap["interaction"] = i

			currencyMap["role"] = role
			id := len(CurrencyCreateRolesChunk)
			CurrencyCreateRolesChunk = append(CurrencyCreateRolesChunk, currencyMap)

			err := s.RequestGuildMembers(i.GuildID, "", 0, "currency_add:"+fmt.Sprint(id), false)
			if err != nil {
				return err
			}

			err = h.ExperimentalResponse(s, i, h.PresetResponse{
				Title:       "Adding currency to roles!",
				Description: "This might take a while.",
			})
			if err != nil {
				fmt.Println(err)
			}
			err = s.RequestGuildMembers(i.GuildID, "", 0, "currency_add:"+fmt.Sprint(id), false)
			if err != nil {
				return err
			}
		}
	} else {
		//Add currency user here
	}

	return nil
}

func CurrencyAddRole(s *discordgo.Session, g *discordgo.GuildMembersChunk) error {

	details := strings.Split(g.Nonce, ":")

	id, err := strconv.Atoi(details[1])
	if err != nil {
		fmt.Println(err)
	}

	currencyMap := CurrencyCreateRolesChunk[id]

	for _, v := range g.Members {
		for _, role := range v.Roles {
			if role == currencyMap["role"].(string) {
				currencyMap["user_id"] = v.User.ID
				err = CurrencyAddUser(s, currencyMap)
				if err != nil {
					h.FollowUpErrorResponse(s, currencyMap["interaction"].(*discordgo.InteractionCreate), fmt.Sprintf("There was an issue creating a claim for <@%s>. Error Message: %s", v.User.ID, err))
				}
			}
		}
	}

	h.FollowUpSuccessResponse(s, currencyMap["interaction"].(*discordgo.InteractionCreate), h.PresetResponse{
		Title:       "__**Claim Create Role**__",
		Description: fmt.Sprintf("Claims are currently being created for all users in <@&%s>", currencyMap["role"]),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "**Progress**",
				Value:  fmt.Sprintf("`%d`/`%d` chunks completed", g.ChunkIndex+1, g.ChunkCount),
				Inline: false,
			},
		},
	})
	return nil
}

func CurrencyAddUser(s *discordgo.Session, currencyMap map[string]interface{}) error {

	result := database.DB.Model(database.UserProfile{}).FirstOrCreate(currencyMap, currencyMap)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
