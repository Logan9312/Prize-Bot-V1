package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/database"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
	"gorm.io/gorm/clause"
)

var SettingsCommand = discordgo.ApplicationCommand{
	Name:        "settings",
	Description: "Change any settings for your bot",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "claiming",
			Description: "Edit your claiming settings",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "category",
					Description: "The category to claim prizes in.",
					ChannelTypes: []discordgo.ChannelType{
						4,
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "log_channel",
					Description: "The output channel for completed tickets.",
					ChannelTypes: []discordgo.ChannelType{
						0,
						5,
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "instructions",
					Description: "Leave instructions for whoever opens the ticket.",
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "disable_claiming",
					Description: "Disables the claiming system, only leaves a record of prizes.",
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "auctions",
			Description: "Edit your auction settings",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "category",
					Description: "Sets the category to create auctions in.",
					ChannelTypes: []discordgo.ChannelType{
						4,
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "currency",
					Description: "Sets the auction currency",
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "currency_side",
					Description: "Left/Right currency",
					//Autocomplete: true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Left",
							Value: "left",
						},
						{
							Name:  "Right",
							Value: "right",
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "integer_only",
					Description: "Only allow integer bids (no decimals).",
				},
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "log_channel",
					Description: "Sets the channel where auctions will send outputs when they end",
					Required:    false,
					ChannelTypes: []discordgo.ChannelType{
						0,
						5,
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "alert_role",
					Description: "Set a role to get pinged whenever an auction starts.",
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "host_role",
					Description: "Set a role that can host auctions.",
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "snipe_extension",
					Description: "Set 0 to disable. Duration an auction by when a bid is placed within the snipe range. (Example: 5m)",
					Autocomplete: true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "snipe_range",
					Description: "Set 0 to disable. The remaining time needed to activate Anti-Snipe (Example: 24h, or 1d)",
					Autocomplete: true,
				},
			},
		},
		{
			Type:         discordgo.ApplicationCommandOptionSubCommand,
			Name:         "giveaways",
			Description:  "Edit your giveaway settings",
			ChannelTypes: []discordgo.ChannelType{},
			Required:     false,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "alert_role",
					Description: "Set a role to get pinged whenever an auction starts.",
				},
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "log_channel",
					Description: "Sets the channel where giveaway will send outputs when they end",
					Required:    false,
					ChannelTypes: []discordgo.ChannelType{
						0,
						5,
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "host_role",
					Description: "Set a role that can host auctions.",
				},
				/*{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "auto_enter",
					Description: "Anyone with this role will be automatically entered.",
				},*/
			},
		},
		/*{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "shop",
			Description: "Edit your shop settings",
			Options: []*discordgo.ApplicationCommandOption{},
		},*/
	},
}

func Settings(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if i.Member.Permissions&(1<<3) != 8 {
		h.ErrorResponse(s, i, "User must have administrator permissions to run this command")
		return
	}

	options := h.ParseSubCommand(i)
	responseFields := []*discordgo.MessageEmbedField{}
	menuOptions := []discordgo.SelectMenuOption{}
	var commandData discordgo.ApplicationCommandOption
	var antiSnipeDescription string
	var model interface{}
	var customID string
	var content string
	var err error

	switch i.ApplicationCommandData().Options[0].Name {
	case "claiming":
		model = &database.ClaimSetup{
			GuildID: i.GuildID,
		}
		commandData = *SettingsCommand.Options[0]
		customID = "clear_claim_setup"
	case "auctions":
		model = &database.AuctionSetup{
			GuildID: i.GuildID,
		}
		commandData = *SettingsCommand.Options[1]
		customID = "clear_auction_setup"
	case "giveaways":
		model = &database.GiveawaySetup{
			GuildID: i.GuildID,
		}
		commandData = *SettingsCommand.Options[2]
		customID = "clear_giveaway_setup"
	case "shop":
		model = &database.ShopSetup{
			GuildID: i.GuildID,
		}
		commandData = *SettingsCommand.Options[3]
		customID = "clear_shop_setup"
	default:
		h.ErrorResponse(s, i, fmt.Sprintf("Unknown settings type: %s", i.ApplicationCommandData().Name))
		return
	}

	options["guild_id"] = i.GuildID

	for key := range options {
		content += fmt.Sprintf("• %s has been successfully set.\n", strings.Title(strings.ReplaceAll(key, "_", " ")))

		switch key {
		case "snipe_extension", "snipe_range":
			options[key], err = h.ParseTime(options[key].(string))
		}
		if err != nil {
			fmt.Println(err)
			h.ErrorResponse(s, i, err.Error())
			return
		}
	}

	result := database.DB.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Model(model).Create(options)
	if result.Error != nil {
		h.ErrorResponse(s, i, result.Error.Error())
		return
	}

	result = database.DB.Model(model).Updates(options)
	if result.Error != nil {
		h.ErrorResponse(s, i, result.Error.Error())
		return
	}

	result = database.DB.Model(model).First(options, i.GuildID)
	if result.Error != nil {
		h.ErrorResponse(s, i, result.Error.Error())
		return
	}

	for _, v := range commandData.Options {
		if !strings.Contains(v.Name, "snipe") {
			switch {
			case options[v.Name] == nil || options[v.Name] == "":
				options[v.Name] = "Not Set"
			case strings.Contains(v.Name, "role"):
				options[v.Name] = fmt.Sprintf("<@&%s>", options[v.Name])
			case strings.Contains(v.Name, "channel"):
				options[v.Name] = fmt.Sprintf("<#%s>", options[v.Name])
			case strings.Contains(v.Name, "category"):
				category, err := s.Channel(options[v.Name].(string))
				if err != nil {
					fmt.Println("Category Error:", err)
					options[v.Name] = "Error Displaying Category: " + err.Error()
				} else {
					options[v.Name] = category.Name
				}
			}
			responseFields = append(responseFields, &discordgo.MessageEmbedField{
				Name:  fmt.Sprintf("**%s**", strings.Title(strings.ReplaceAll(v.Name, "_", " "))),
				Value: fmt.Sprint(options[v.Name]),
			})
		}
		menuOptions = append(menuOptions, discordgo.SelectMenuOption{
			Label:       strings.Title(strings.ReplaceAll(v.Name, "_", " ")),
			Value:       v.Name,
			Description: v.Description,
		})
	}

	if options["snipe_range"] != nil && options["snipe_extension"] != nil {
		antiSnipeDescription = fmt.Sprintf("If a bid is placed within %s of the auction ending, it will be extended by %s.", options["snipe_range"].(time.Duration).String(), options["snipe_extension"].(time.Duration).String())
	} else {
		antiSnipeDescription = "Anti Snipe Disabled. To enable, set both snipe_extension and snipe_range"
	}

	responseFields = append(responseFields, &discordgo.MessageEmbedField{
		Name:  "**Anti Snipe**",
		Value: antiSnipeDescription,
	})

	err = h.SuccessResponse(s, i, h.PresetResponse{
		Title:       fmt.Sprintf("%s Setup", strings.Title(i.ApplicationCommandData().Name)),
		Description: content,
		Fields:      responseFields,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID:    customID,
						Placeholder: "Clear Setup Options",
						MinValues:   1,
						MaxValues:   len(commandData.Options),
						Options:     menuOptions,
					},
				},
			},
		},
	})
	if err != nil {
		fmt.Println(err)
	}
}
