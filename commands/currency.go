package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/database"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gorm.io/gorm/clause"
)

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
	DefaultMemberPermissions: h.Ptr(int64(discordgo.PermissionManageServer)),
	DMPermission:             new(bool),
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

func CurrencyEdit(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	currencyMap := h.ParseSubCommand(i)
	action := currencyMap["action"].(string)

	currencyMap["guild_id"] = i.GuildID
	currencyMap["interaction"] = i
	if i.ApplicationCommandData().Resolved.Roles[currencyMap["target"].(string)] != nil {
		if !CheckPremiumGuild(i.GuildID) {
			err := h.PremiumError(s, i, "Premium is needed to edit the currency of an entire role. Please select only an user or purchase premium to use this function.")
			if err != nil {
				fmt.Println(err)
			}
			return nil
		}
		currencyMap["role"] = i.ApplicationCommandData().Resolved.Roles[currencyMap["target"].(string)]

		h.SaveChunkData(i.ID, currencyMap)
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
			return CurrencyAddUser(i.GuildID, i.ApplicationCommandData().Resolved.Users[currencyMap["target"].(string)].ID, currencyMap["amount"].(float64))
		case "subtract":
			return CurrencySubtractUser(i.GuildID, i.ApplicationCommandData().Resolved.Users[currencyMap["target"].(string)].ID, currencyMap["amount"].(float64))
		case "set":
			return CurrencySetUser(i.GuildID, i.ApplicationCommandData().Resolved.Users[currencyMap["target"].(string)].ID, currencyMap["amount"].(float64))
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

	currencyMap := h.ReadChunkData(details[1])

	amount := currencyMap["amount"].(float64)

	successCount, userCount, err := CurrencyEditRole(s, g, currencyMap["role"].(*discordgo.Role).ID, amount, currencyMap["action"].(string))
	if err != nil {
		return err
	}

	_, err = h.FollowUpSuccessResponse(s, currencyMap["interaction"].(*discordgo.InteractionCreate), h.PresetResponse{
		Title:       "__**Currency Add Role**__",
		Description: fmt.Sprintf("Currency is currently being added to all users in <@&%s>", currencyMap["role"].(*discordgo.Role).ID),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "**Progress**",
				Value:  fmt.Sprintf("`%d`/`%d` chunks completed\n%d/%d user's currency applied successfully in this chunk.", g.ChunkIndex+1, g.ChunkCount, successCount, userCount),
				Inline: false,
			},
		},
	})
	return err
}

func CurrencyEditRole(s *discordgo.Session, g *discordgo.GuildMembersChunk, roleID string, amount float64, action string) (int, int, error) {
	errCount := 0
	userCount := 0
	successCount := 0
	var err error
	for _, v := range g.Members {
		if roleID != g.GuildID && !HasRole(v, roleID) {
			continue
		}
		if v.User.Bot {
			continue
		}

		userCount++
		switch action {
		case "add":
			err = CurrencyAddUser(g.GuildID, v.User.ID, amount)
		case "subtract":
			err = CurrencySubtractUser(g.GuildID, v.User.ID, amount)
		case "set":
			err = CurrencySetUser(g.GuildID, v.User.ID, amount)
		}

		if err != nil && errCount < 5 {
			data := h.ReadChunkData(strings.Split(g.Nonce, ":")[1])
			fmt.Println(err)
			_, err = h.FollowUpErrorResponse(s, data["interaction"].(*discordgo.InteractionCreate), fmt.Sprintf("There was an issue adding currency for <@%s>. Error Message: %s", v.User.ID, err))
			if err != nil {
				fmt.Println(err)
			}
			errCount++
			if errCount == 5 {
				_, err = h.FollowUpErrorResponse(s, data["interaction"].(*discordgo.InteractionCreate), "**5 or more users have failed to update their currency.** Refer to previous errors for more information.")
				if err != nil {
					fmt.Println(err)
				}
			}
		} else {
			successCount++
		}
	}

	return successCount, userCount, nil
}

func CurrencySubtractUser(guildID, userID string, amount float64) error {
	return CurrencyAddUser(guildID, userID, -1*amount)
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

	if amount < 0 {
		return fmt.Errorf("<@%s> does not have enough currency to complete this request. Resulting balance would be: %s", userID, PriceFormat(amount, guildID, nil))
	}

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

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		fmt.Println(err)
	}

	currencyList := ""

	for k, v := range userMapSlice {
		member, err := s.State.Member(i.GuildID, v["user_id"].(string))
		if err != nil {
			fmt.Println(err)
			member, err = s.GuildMember(i.GuildID, v["user_id"].(string))
			if err != nil {
				fmt.Println(err)
				continue
			}
		}

		currencyList += fmt.Sprintf("%d. %s#%s %s\n", k, member.User.Username, member.User.Discriminator, strings.TrimRight(strings.TrimRight(p.Sprintf("%f", v["balance"].(float64)), "0"), "."))
	}

	_, err = h.DeferredResponse(s, i, h.PresetResponse{
		Content:     "",
		Title:       "Currency List",
		Description: "Successfully generated a list of all users in the server",
		Fields:      []*discordgo.MessageEmbedField{},
		Thumbnail:   &discordgo.MessageEmbedThumbnail{},
		Image:       &discordgo.MessageEmbedImage{},
		Components:  []discordgo.MessageComponent{},
		Embeds:      []*discordgo.MessageEmbed{},
		Files:       []*discordgo.File{{Name: "currency_list.txt", ContentType: "txt", Reader: strings.NewReader(currencyList)}},
		Reference:   &discordgo.MessageReference{},
	})

	return err
}
