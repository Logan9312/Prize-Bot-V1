package commands

import "github.com/bwmarrin/discordgo"

var (
	CommandOptionItem = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "item",
		Description: "Select an item.",
	}

	CommandOptionBid = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionNumber,
		Name:        "bid",
		Description: "The starting/current bid price",
	}

	CommandOptionDescription = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "description",
		Description: "A description of the item or event.",
	}

	CommandOptionDuration = &discordgo.ApplicationCommandOption{
		Type:         discordgo.ApplicationCommandOptionString,
		Name:         "duration",
		Description:  "Duration that the event will run for. (Ex: 24h)",
		Autocomplete: true,
	}

	CommandOptionCurrency = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "currency",
		Description: "Sets the currency to use",
	}

	CommandOptionUseCurrency = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionBoolean,
		Name:        "use_currency",
		Description: "The winner will pay with their currency balance.",
	}

	CommandOptionIntegerOnly = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionBoolean,
		Name:        "integer_only",
		Description: "Only allow integer bids (no decimals).",
	}

	CommandOptionIncrementMax = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionNumber,
		Name:        "increment_max",
		Description: "The max amount someone can bid at once",
	}

	CommandOptionIncrementMin = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionNumber,
		Name:        "increment_min",
		Description: "The minimum amount someone can bid at once",
	}

	CommandOptionBuyout = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionNumber,
		Name:        "buyout",
		Description: "Set a price that someone can immediately win the auction for if they bid it or higher.",
		Required:    false,
	}

	CommandOptionTargetPrice = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionNumber,
		Name:        "target_price",
		Description: "If this hidden price is not reached, no winner will be chosen",
	}

	CommandOptionCategory = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionChannel,
		Name:        "category",
		Description: "Sets the category to create auctions in.",
		ChannelTypes: []discordgo.ChannelType{
			4,
		},
	}

	CommandOptionImage = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionAttachment,
		Name:        "image",
		Description: "Attach an image to your auction",
	}

	CommandOptionSchedule = &discordgo.ApplicationCommandOption{
		Type:         discordgo.ApplicationCommandOptionString,
		Name:         "schedule",
		Description:  "Set how long until the auction starts. (Example: 24h, or 1d)",
		Autocomplete: true,
	}

	CommandOptionHost = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionUser,
		Name:        "host",
		Description: "Changes the host",
		//Autocomplete: true,
	}

	CommandOptionWinner = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionUser,
		Name:        "winner",
		Description: "Set the current winner",
		//Autocomplete: true,
	}

	CommandOptionExtend = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "extend",
		Description: "Extend the length of the auction, use a negative value to reduce the time. (Example: 24h, or 1d)",
		//Autocomplete: true,
	}

	CommandOptionCurrencySide = &discordgo.ApplicationCommandOption{
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
	}

	CommandOptionQueueNumber = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "queue_number",
		Description: "The number of the auction if you are editing one in queue.",
		//Autocomplete: true,
	}

	CommandOptionLogChannel = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionChannel,
		Name:        "log_channel",
		Description: "The output channel for completed tickets.",
		ChannelTypes: []discordgo.ChannelType{
			0,
			5,
		},
	}

	CommandOptionInstructions = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "instructions",
		Description: "Leave instructions for whoever opens the ticket.",
	}

	CommandOptionDisableClaiming = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionBoolean,
		Name:        "disable_claiming",
		Description: "Disables the claiming system, only leaves a record of prizes.",
	}

	CommandOptionWinners = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionInteger,
		Name:        "winners",
		Description: "The number of winners.",
	}

	CommandOptionAlertRole = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionRole,
		Name:        "alert_role",
		Description: "Set a role to get pinged whenever an auction starts.",
	}

	CommandOptionSnipeExtension = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "snipe_extension",
		Description: "Duration to extend an auction by when a bid is placed within the snipe range. (Example: 5m)",
		//Autocomplete: true,
	}

	CommandOptionSnipeRange = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "snipe_range",
		Description: "The remaining time needed to activate Anti-Snipe (Example: 24h)",
		//Autocomplete: true,
	}

	CommandOptionChannelLock = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionBoolean,
		Name:        "channel_lock",
		Description: "Enabling this starts the auction in the channel where you run the auction create command.",
	}

	CommandOptionChannelPrefix = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "channel_prefix",
		Description: "modify the prefix of auction channels. Ex: 💸│",
	}

	CommandOptionAutoEnter = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionRole,
		Name:        "auto_enter",
		Description: "Anyone with this role will be automatically entered.",
	}

	CommandOptionTarget = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionMentionable,
		Name:        "target",
		Description: "Select a group.",
	}

	CommandOptionUser = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionUser,
		Name:        "user",
		Description: "The user who you want to display.",
	}

	CommandOptionChannel = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionChannel,
		Name:        "channel",
		Description: "The output channel for the claims..",
		ChannelTypes: []discordgo.ChannelType{
			0,
			5,
		},
	}

	CommandOptionAmount = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionNumber,
		Name:        "amount",
		Description: "The amount of currency",
	}

	CommandOptionTimeUntil = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "time_until",
		Description: "The time until the start.",
	}
)

func Require(option *discordgo.ApplicationCommandOption) *discordgo.ApplicationCommandOption {
	option.Required = true
	return option
}
