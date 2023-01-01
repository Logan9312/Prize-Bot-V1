package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
	r "gitlab.com/logan9312/discord-auction-bot/responses"
)

var AuctionCommand = discordgo.ApplicationCommand{
	Name:        "auction",
	Description: "Put an item up for auction!",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "schedule",
			Description: "Schedule an auction",
			Options: []*discordgo.ApplicationCommandOption{
				Require(CommandOptionTimeUntil),
				Require(CommandOptionItem),
				Require(CommandOptionBid),
				Require(CommandOptionDuration),
				CommandOptionDescription,
				CommandOptionCurrency,
				CommandOptionUseCurrency,
				CommandOptionIntegerOnly,
				CommandOptionIncrementMax,
				CommandOptionIncrementMin,
				CommandOptionBuyout,
				CommandOptionTargetPrice,
				CommandOptionCategory,
				CommandOptionImage,
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "create",
			Description: "Create an Auction",
			Options: []*discordgo.ApplicationCommandOption{
				Require(CommandOptionItem),
				Require(CommandOptionBid),
				Require(CommandOptionDuration),
				CommandOptionDescription,
				CommandOptionCurrency,
				CommandOptionUseCurrency,
				CommandOptionIntegerOnly,
				CommandOptionIncrementMax,
				CommandOptionIncrementMin,
				CommandOptionBuyout,
				CommandOptionTargetPrice,
				CommandOptionCategory,
				CommandOptionImage,
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "queue",
			Description: "Display the current auction Queue",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "edit",
			Description: "Edit any auction details",
			Options: []*discordgo.ApplicationCommandOption{
				//CommandOptionQueueNumber
				CommandOptionItem,
				CommandOptionHost,
				CommandOptionBid,
				CommandOptionHost,
				CommandOptionWinner,
				CommandOptionExtend,
				CommandOptionDescription,
				CommandOptionCurrency,
				CommandOptionUseCurrency,
				CommandOptionCurrencySide,
				CommandOptionIntegerOnly,
				CommandOptionIncrementMax,
				CommandOptionIncrementMin,
				CommandOptionBuyout,
				CommandOptionTargetPrice,
				CommandOptionImage,
			},
		},
	},
}

var BidCommand = discordgo.ApplicationCommand{
	Name:        "bid",
	Description: "Bid on an Auction",
	Options: []*discordgo.ApplicationCommandOption{
		Require(CommandOptionBid),
	},
}

func Auction(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	switch i.ApplicationCommandData().Options[0].Name {
	case "schedule":
		return AuctionSchedule(s, i)
	case "create":
		return AuctionCreate(s, i)
	case "queue":
		return AuctionQueue(s, i)
	case "edit":
		return AuctionEdit(s, i)
	}
	return fmt.Errorf("Unknown Auction command, please contact support")
}

func AuctionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	options := h.ParseSubCommand(i)

	if options["image"] != nil {
		options["image_url"] = i.ApplicationCommandData().Resolved.Attachments[options["image"].(string)].URL
	}

	options["channel_id"] = i.ChannelID

	duration, err := h.ParseTime(strings.ToLower(options["duration"].(string)))
	if err != nil {
		return fmt.Errorf("Error parsing duration input: %w", err)
	}

	auctions := strings.Split(options["item"].(string), ";")

	if len(auctions) > 5 && !CheckPremiumGuild(i.GuildID) {
		r.PremiumError(s, i, "Free users can only start 5 auctions in bulk. Upgrade to premium to start up to 100 in bulk.")
	}

	if len(auctions) > 100 {
		return fmt.Errorf("You can only start 100 auctions in bulk at once. You attempted to start: %d.", len(multiAuctions))
	}

	for _, item := range auctions {
		auctionMap := map[string]any{}
		for k, v := range options {
			auctionMap[k] = v
		}
		auctionMap["item"] = item

		channelID, err := AuctionHandler(s, auctionMap, i.Member, i.GuildID, duration)
		if err != nil {
			return err
		}

		if channelID != "" {
			err = h.SuccessResponse(s, i, h.PresetResponse{
				Title:       "**Auction Starting**",
				Description: fmt.Sprintf("Auction has successfully been started in <#%s>!", channelID),
			})
			if err != nil {
				_, err = h.FollowUpSuccessResponse(s, i, h.PresetResponse{
					Title:       "**Auction Starting**",
					Description: fmt.Sprintf("Auction has successfully been started in <#%s>!", channelID),
				})
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
