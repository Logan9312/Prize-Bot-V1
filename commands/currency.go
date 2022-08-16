package commands

import (
	"fmt"
	"strings"

	"sync"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/database"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gorm.io/gorm/clause"
)

type CurrencyEventTracker struct {
	m    sync.Mutex
	data map[string]map[string]any
}

var CurrencyCreateData = CurrencyEventTracker{
	data: map[string]map[string]any{},
}

var CurrencyCommand = discordgo.ApplicationCommand{
	Name:        "currency",
	Description: "Shows the current currency of the bot",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "edit",
			Description: "Change the currency of a user or role.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "action",
					Description: "How to apply currency to the user.",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "add",
							Value: "add",
						},
						{
							Name:  "set",
							Value: "set",
						},
						{
							Name:  "subtract",
							Value: "subtract",
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionMentionable,
					Name:        "target",
					Description: "Select a user or a role to modify the currency of",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionNumber,
					Name:        "amount",
					Description: "The amount of currency",
					Required:    true,
					Options:     []*discordgo.ApplicationCommandOption{},
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "list",
			Description: "Lists the currency of all saved users.",
		},
	},
}

func Currency(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	switch i.ApplicationCommandData().Options[0].Name {
	case "edit":
		return CurrencyEdit(s, i)
	case "list":
		return CurrencyList(s, i)
	}
	return fmt.Errorf("Unknown Currency command, please contact support")
}

func SaveCurrencyData(id string, currencyMap map[string]any) {
	CurrencyCreateData.m.Lock()
	CurrencyCreateData.data[id] = currencyMap
	CurrencyCreateData.m.Unlock()
}

func ReadCurrencyData(id string) map[string]any {
	CurrencyCreateData.m.Lock()
	defer CurrencyCreateData.m.Unlock()
	return CurrencyCreateData.data[id]
}

func CurrencyEdit(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	if !CheckPremiumGuild(i.GuildID) {
		err := h.PremiumError(s, i)
		if err != nil {
			fmt.Println(err)
		}
		return nil
	}

	currencyMap := h.ParseSubCommand(i)
	action := currencyMap["action"].(string)

	currencyMap["guild_id"] = i.GuildID
	currencyMap["interaction"] = i
	if i.ApplicationCommandData().Resolved.Roles[currencyMap["target"].(string)] != nil {

		currencyMap["role"] = i.ApplicationCommandData().Resolved.Roles[currencyMap["target"].(string)]

		SaveCurrencyData(i.ID, currencyMap)
		nonce := "$:" + i.ID

		err := s.RequestGuildMembers(i.GuildID, "", 0, nonce, false)
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

	if i.ApplicationCommandData().Resolved.Users[currencyMap["target"].(string)] != nil {
		err := h.ExperimentalResponse(s, i, h.PresetResponse{
			Title:       "Adding currency to user!",
			Description: "",
		})
		if err != nil {
			fmt.Println(err)
		}
		switch action {
		case "add":
			err = CurrencyAddUser(i.GuildID, i.ApplicationCommandData().Resolved.Users[currencyMap["target"].(string)].ID, currencyMap["amount"].(float64))
			if err != nil {
				return err
			}
		case "subtract":
			err = CurrencyAddUser(i.GuildID, i.ApplicationCommandData().Resolved.Users[currencyMap["target"].(string)].ID, -1*currencyMap["amount"].(float64))
			if err != nil {
				return err
			}
		case "set":
			err = CurrencySetUser(i.GuildID, i.ApplicationCommandData().Resolved.Users[currencyMap["target"].(string)].ID, currencyMap["amount"].(float64))
			if err != nil {
				return err
			}
		}
	}

	return nil
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

func CurrencyRoleHandler(s *discordgo.Session, g *discordgo.GuildMembersChunk) error {

	details := strings.Split(g.Nonce, ":")

	currencyMap := ReadCurrencyData(details[1])

	amount := currencyMap["amount"].(float64)

	defer h.FollowUpSuccessResponse(s, currencyMap["interaction"].(*discordgo.InteractionCreate), h.PresetResponse{
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

	switch currencyMap["action"].(string) {
	case "add":
		return CurrencyAddRole(g, currencyMap["role"].(*discordgo.Role).ID, amount)
	case "subtract":
		return CurrencyAddRole(g, currencyMap["role"].(*discordgo.Role).ID, -1*amount)
	case "set":
		return CurrencySetRole(g, currencyMap["role"].(*discordgo.Role).ID, amount)
	}

	return nil
}

func CurrencyAddRole(g *discordgo.GuildMembersChunk, roleID string, amount float64) error {
	if roleID == g.GuildID {
		for _, v := range g.Members {
			err := CurrencyAddUser(g.GuildID, v.User.ID, amount)
			if err != nil {
				return fmt.Errorf("There was an issue creating a claim for <@%s>. Error Message: %s", v.User.ID, err)
			}
		}
	}
	for _, v := range g.Members {
		for _, role := range v.Roles {
			if role == roleID {
				err := CurrencyAddUser(g.GuildID, v.User.ID, amount)
				if err != nil {
					return fmt.Errorf("There was an issue creating a claim for <@%s>. Error Message: %s", v.User.ID, err)
				}
			}
		}
	}
	return nil
}

func CurrencySetRole(g *discordgo.GuildMembersChunk, roleID string, amount float64) error {
	if roleID == g.GuildID {
		for _, v := range g.Members {
			err := CurrencySetUser(g.GuildID, v.User.ID, amount)
			if err != nil {
				return fmt.Errorf("There was an issue creating a claim for <@%s>. Error Message: %s", v.User.ID, err)
			}
		}
	}
	for _, v := range g.Members {
		for _, role := range v.Roles {
			if role == roleID {
				err := CurrencySetUser(g.GuildID, v.User.ID, amount)
				if err != nil {
					return fmt.Errorf("There was an issue creating a claim for <@%s>. Error Message: %s", v.User.ID, err)
				}
			}
		}
	}
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

	return CurrencySetUser(guildID, userID, amount)
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

	fmt.Println("User Currency: ", userID, amount)

	result = database.DB.Model(database.UserProfile{
		UserID:  userID,
		GuildID: guildID,
	}).Updates(map[string]any{"balance": amount})
	if result.Error != nil {
		fmt.Println(result.Error)
		return result.Error
	}
	return err
}

func CurrencyList(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	p := message.NewPrinter(language.English)
	userMapSlice := []map[string]any{}
	result := database.DB.Model(&database.UserProfile{}).Where(map[string]interface{}{"guild_id": i.GuildID}).Order("balance DESC").Find(&userMapSlice)
	if result.Error != nil {
		return result.Error
	}

	currencyList := ""

	for k, v := range userMapSlice {
		member, err := s.State.Member(i.GuildID, v["user_id"].(string))
		if err != nil {
			fmt.Println(err)
			member, err = s.GuildMember(i.GuildID, v["user_id"].(string))
			if err != nil {
				return err
			}
		}

		currencyList += fmt.Sprintf("%d. %s#%s %s\n", k, member.User.Username, member.User.Discriminator, strings.TrimRight(strings.TrimRight(p.Sprintf("%f", v["balance"].(float64)), "0"), "."))
		fmt.Println(currencyList)
	}

	_, err := s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
		Files: []*discordgo.File{{Name: "currency_list.txt", ContentType: "txt", Reader: strings.NewReader(currencyList)}},
	})

	return err
}
