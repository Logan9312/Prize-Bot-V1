package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/database"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gorm.io/gorm/clause"
)

var CurrencyCreateRolesChunk = []map[string]interface{}{}

var CurrencyCommand = discordgo.ApplicationCommand{
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
					Name:        "target",
					Description: "Select a user or a role* (premium) add currency to.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionNumber,
					Name:        "amount",
					Description: "The amount of currency to add.",
					Required:    true,
					Options:     []*discordgo.ApplicationCommandOption{},
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "set",
			Description: "Set the exact amount of currency for a user.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionMentionable,
					Name:        "target",
					Description: "Select a user or a role* (premium) add currency to.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionNumber,
					Name:        "amount",
					Description: "The amount of currency to add.",
					Required:    true,
					Options:     []*discordgo.ApplicationCommandOption{},
				},
			},
		},
	},
}

func CurrencyChunkHandler() {

}

func Currency(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	switch i.ApplicationCommandData().Options[0].Name {
	case "add":
		return CurrencyAdd(s, i)
	case "set":
		return CurrencySet(s, i)
	}
	return fmt.Errorf("Unknown Currency command, please contact support")
}

func PriceFormat(price float64, guildID string, override interface{}) string {

	p := message.NewPrinter(language.English)
	currencyMap := map[string]interface{}{}
	priceString := strings.TrimRight(strings.TrimRight(p.Sprintf("%f", price), "0"), ".")

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

	if i.ApplicationCommandData().Resolved.Roles[currencyMap["target"].(string)] != nil {
		for _, role := range i.ApplicationCommandData().Resolved.Roles {
			currencyMap["interaction"] = i

			currencyMap["role"] = role
			id := i.ID
			CurrencyCreateRolesChunk = append(CurrencyCreateRolesChunk, currencyMap)

			err := s.RequestGuildMembers(i.GuildID, "", 0, "currency_add:"+id+":"+fmt.Sprint(currencyMap["amount"]), false)
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

		}
	} else {
		for _, user := range i.ApplicationCommandData().Resolved.Users {

			err := h.ExperimentalResponse(s, i, h.PresetResponse{
				Title:       "Adding currency to user!",
				Description: "",
			})
			if err != nil {
				fmt.Println(err)
			}
			err = CurrencyAddUser(i.GuildID, user.ID, currencyMap["amount"].(float64))
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	return nil
}

func CurrencySet(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	currencyMap := h.ParseSubCommand(i)
	currencyMap["guild_id"] = i.GuildID

	if i.ApplicationCommandData().Resolved.Roles[currencyMap["target"].(string)] != nil {
		for _, role := range i.ApplicationCommandData().Resolved.Roles {
			currencyMap["interaction"] = i

			currencyMap["role"] = role
			id := i.ID
			CurrencyCreateRolesChunk = append(CurrencyCreateRolesChunk, currencyMap)

			err := s.RequestGuildMembers(i.GuildID, "", 0, "currency_set:"+id+":"+fmt.Sprint(currencyMap["amount"]), false)
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

		}
	} else {
		for _, user := range i.ApplicationCommandData().Resolved.Users {

			err := h.ExperimentalResponse(s, i, h.PresetResponse{
				Title:       "Adding currency to user!",
				Description: "",
			})
			if err != nil {
				fmt.Println(err)
			}
			err = CurrencyAddUser(i.GuildID, user.ID, currencyMap["amount"].(float64))
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	return nil
}

func CurrencyAddRole(s *discordgo.Session, g *discordgo.GuildMembersChunk) error {

	details := strings.Split(g.Nonce, ":")

	fmt.Println(details)

	id, err := strconv.Atoi(details[1])
	if err != nil {
		fmt.Println(err)
		return err
	}

	amount, err := strconv.ParseFloat(details[2], 64)
	if err != nil {
		fmt.Println(err)
		return err
	}

	currencyMap := CurrencyCreateRolesChunk[id]

	for _, v := range g.Members {
		for _, role := range v.Roles {
			if role == currencyMap["role"].(*discordgo.Role).ID {
				currencyMap["user_id"] = v.User.ID
				err = CurrencyAddUser(g.GuildID, v.User.ID, amount)
				if err != nil {
					h.FollowUpErrorResponse(s, currencyMap["interaction"].(*discordgo.InteractionCreate), fmt.Sprintf("There was an issue creating a claim for <@%s>. Error Message: %s", v.User.ID, err))
				}
			}
		}
	}

	h.FollowUpSuccessResponse(s, currencyMap["interaction"].(*discordgo.InteractionCreate), h.PresetResponse{
		Title:       "__**Currency Add Role**__",
		Description: fmt.Sprintf("Currency is currently being added to all users in <@&%s>", currencyMap["role"].(*discordgo.Role).ID),
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

func CurrencySetRole(s *discordgo.Session, g *discordgo.GuildMembersChunk) error {

	details := strings.Split(g.Nonce, ":")

	fmt.Println(details)

	id, err := strconv.Atoi(details[1])
	if err != nil {
		fmt.Println(err)
		return err
	}

	amount, err := strconv.ParseFloat(details[2], 64)
	if err != nil {
		fmt.Println(err)
		return err
	}

	currencyMap := CurrencyCreateRolesChunk[id]

	for _, v := range g.Members {
		for _, role := range v.Roles {
			if role == currencyMap["role"].(*discordgo.Role).ID {
				currencyMap["user_id"] = v.User.ID
				err = CurrencySetUser(g.GuildID, v.User.ID, amount)
				if err != nil {
					h.FollowUpErrorResponse(s, currencyMap["interaction"].(*discordgo.InteractionCreate), fmt.Sprintf("There was an issue creating a claim for <@%s>. Error Message: %s", v.User.ID, err))
				}
			}
		}
	}

	h.FollowUpSuccessResponse(s, currencyMap["interaction"].(*discordgo.InteractionCreate), h.PresetResponse{
		Title:       "__**Currency Set Role**__",
		Description: fmt.Sprintf("Currency is being set for all users in <@&%s>", currencyMap["role"].(*discordgo.Role).ID),
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

func CurrencyAddUser(guildID, userID string, amount float64) error {
	result := database.DB.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Model(database.UserProfile{}).Create(map[string]any{
		"user_id":  userID,
		"guild_id": guildID,
		"balance":  0,
	})
	if result.Error != nil {
		return result.Error
	}

	currencyMap := map[string]any{}

	result = database.DB.Model(database.UserProfile{}).First(currencyMap, map[string]any{
		"user_id":  userID,
		"guild_id": guildID,
	})
	if result.Error != nil {
		return result.Error
	}

	amount += currencyMap["balance"].(float64)

	err := SetUserCurrency(amount, guildID, userID)

	return err
}

func CurrencySetUser(guildID, userID string, amount float64) (err error) {
	result := database.DB.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Model(database.UserProfile{}).Create(map[string]any{
		"user_id":  userID,
		"guild_id": guildID,
		"balance":  0,
	})
	if result.Error != nil {
		return result.Error
	}

	currencyMap := map[string]any{}

	result = database.DB.Model(database.UserProfile{}).First(currencyMap, map[string]any{
		"user_id":  userID,
		"guild_id": guildID,
	})
	if result.Error != nil {
		return result.Error
	}

	err = SetUserCurrency(amount, guildID, userID)

	return err
}

func SetUserCurrency(amount float64, guildID, userID string) (err error) {

	result := database.DB.Model(database.UserProfile{
		UserID:  userID,
		GuildID: guildID,
	}).Updates(map[string]any{"balance": amount})
	if result.Error != nil {
		fmt.Println(result.Error)
		return result.Error
	}

	fmt.Println(amount)

	return err

}
