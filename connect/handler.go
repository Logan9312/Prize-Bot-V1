package connect

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"regexp"

	"github.com/bwmarrin/discordgo"
	c "gitlab.com/logan9312/discord-auction-bot/commands"
	"gitlab.com/logan9312/discord-auction-bot/database"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
	"gitlab.com/logan9312/discord-auction-bot/logger"
)

// Move these to commands package
var commandMap = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) error{
	"auction":        c.Auction,
	"bid":            c.AuctionBid,
	"profile":        c.Profile,
	"giveaway":       c.Giveaway,
	"shop":           c.Shop,
	"claim":          c.Claim,
	"privacy_policy": c.Privacy,
	"dev":            c.Dev,
	"premium":        c.Premium,
	"settings":       c.Settings,
	"quest":          c.Quest,
	"currency":       c.Currency,
	"stats":          c.Stats,
	"help":           Help,
	"whitelabel":     Whitelabel,
}

var buttonMap = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) error{
	"endauction":             c.AuctionEndButton,
	"claim_prize":            c.ClaimPrizeButton,
	"clearauction":           c.ClearAuctionButton,
	"delete_auction_queue":   c.DeleteAuctionQueue,
	"delete_auction_channel": c.DeleteAuctionChannel,
	"reroll_giveaway":        c.RerollGiveawayButton,
	"clear_auction_setup":    c.AuctionSetupClearButton,
	"clear_giveaway_setup":   c.GiveawaySetupClearButton,
	"clear_claim_setup":      c.ClaimSetupClearButton,
	"clear_currency_setup":   c.CurrencySetupClearButton,
	"claim_cancel":           c.CancelButton,
	"claim_complete":         c.CompleteButton,
	"reopen_ticket":          c.ReopenTicket,
	"additem":                c.AddItem,
	"bid_history":            c.AuctionBidHistory,
	"questbutton":            c.QuestButton,
	"helpmenu":               HelpMenu,
}

var autoCompleteMap = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) error{
	"auction":  c.AuctionAutoComplete,
	"giveaway": c.GiveawayAutoComplete,
}

var modalSubmitMap = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) error{}

var guildMembersMap = map[string]func(s *discordgo.Session, g *discordgo.GuildMembersChunk) error{
	"claim_create": c.ClaimCreateRole,
	"$":            c.CurrencyRoleHandler,
}

func RegisterHandlers(s *discordgo.Session) {
	s.AddHandler(ReadyHandler)
	s.AddHandler(InteractionHandler)
	s.AddHandler(MessageHandler)
	s.AddHandler(GuildMemberChunkHandler)
	s.AddHandler(GuildCreateHandler)
}

func InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if i.Member == nil {
			h.ErrorResponse(s, i, "Commands cannot be run in a DM. Please contact support if you're not currently in a DM with the bot.")
			return
		}
		commandName := i.ApplicationCommandData().Name
		log := logger.Command(commandName, i.GuildID, i.Member.User.ID, i.Member.User.Username)
		log.Info("command invoked")

		if f, ok := commandMap[commandName]; ok {
			err := f(s, i)
			if err != nil {
				log.Errorw("command failed", "error", err)
				err2 := h.ErrorResponse(s, i, err.Error())
				if err2 != nil {
					log.Warnw("failed to send error response, trying followup", "error", err2)
					_, err = h.FollowUpErrorResponse(s, i, err.Error())
					if err != nil {
						log.Errorw("failed to send followup error response", "error", err)
					}
				}
			}
		} else {
			log.Error("command handler not found")
			h.ErrorResponse(s, i, "Command response has not been set properly, please contact Logan to fix")
		}

	case discordgo.InteractionMessageComponent:
		buttonID := strings.Split(i.MessageComponentData().CustomID, ":")[0]
		log := logger.Button(buttonID, i.GuildID, i.Member.User.ID, i.Member.User.Username)
		log.Debug("button clicked")

		if f, ok := buttonMap[buttonID]; ok {
			err := f(s, i)
			if err != nil {
				log.Errorw("button handler failed", "error", err)
				err2 := h.ErrorResponse(s, i, err.Error())
				if err2 != nil {
					log.Warnw("failed to send error response, trying followup", "error", err2)
					_, err = h.FollowUpErrorResponse(s, i, err.Error())
					if err != nil {
						log.Errorw("failed to send followup error response", "error", err)
					}
				}
			}
		} else {
			log.Error("button handler not found")
			h.ErrorResponse(s, i, "Button response has not been set properly, please contact Logan to fix")
		}

	case discordgo.InteractionApplicationCommandAutocomplete:
		commandName := i.ApplicationCommandData().Name
		if f, ok := autoCompleteMap[commandName]; ok {
			err := f(s, i)
			if err != nil {
				logger.Sugar.Warnw("autocomplete failed",
					"command", commandName,
					"error", err,
				)
				err = h.ErrorResponse(s, i, err.Error())
				if err != nil {
					logger.Sugar.Errorw("failed to send autocomplete error response", "error", err)
				}
			}
		} else {
			h.ErrorResponse(s, i, "AutoComplete response has not been set properly, please contact Logan to fix")
		}

	case discordgo.InteractionModalSubmit:
		modalID := i.ModalSubmitData().CustomID
		log := logger.Sugar.With("modal_id", modalID, "guild_id", i.GuildID)

		if modalID == "whitelabel_token" {
			err := WhitelabelTokenModal(s, i)
			if err != nil {
				log.Errorw("whitelabel modal failed", "error", err)
				err2 := h.ErrorResponse(s, i, err.Error())
				if err2 != nil {
					log.Warnw("failed to send error response, trying followup", "error", err2)
					_, err = h.FollowUpErrorResponse(s, i, err.Error())
					if err != nil {
						log.Errorw("failed to send followup error response", "error", err)
					}
				}
			}
		} else if f, ok := modalSubmitMap[modalID]; ok {
			err := f(s, i)
			if err != nil {
				log.Errorw("modal handler failed", "error", err)
				err2 := h.ErrorResponse(s, i, err.Error())
				if err2 != nil {
					log.Warnw("failed to send error response, trying followup", "error", err2)
					_, err = h.FollowUpErrorResponse(s, i, err.Error())
					if err != nil {
						log.Errorw("failed to send followup error response", "error", err)
					}
				}
			}
		} else {
			log.Error("modal handler not found")
			h.ErrorResponse(s, i, "Modal Submit response has not been set properly, please contact Logan to fix")
		}
	}
}

func ReadyHandler(s *discordgo.Session, i *discordgo.Ready) {
	log := logger.Bot(s.State.User.ID, s.State.User.Username)
	_, err := s.ChannelMessageSend("943175605858496602", "Bot has finished restarting")
	if err != nil {
		log.Warnw("failed to send restart notification", "error", err)
	}
	log.Info("bot is ready")
}

func GuildCreateHandler(s *discordgo.Session, g *discordgo.GuildCreate) {
	channelID := "1005255087200948234"

	fields := []*discordgo.MessageEmbedField{
		{
			Name:  "Name",
			Value: g.Name,
		},
		{
			Name:  "Member Count",
			Value: fmt.Sprint(g.MemberCount),
		},
	}
	if g.VanityURLCode != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Vanity URL",
			Value: g.VanityURLCode,
		})
	}

	if time.Since(g.JoinedAt) < 5*time.Minute {
		logger.Sugar.Infow("bot joined new server",
			"guild_name", g.Name,
			"guild_id", g.ID,
			"member_count", g.MemberCount,
		)

		_, err := h.SuccessMessage(s, channelID, h.PresetResponse{
			Title:       "New Server Joined!",
			Description: "The bot has joined a new discord within the last 5 minutes.",
			Fields:      fields,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: g.IconURL(""),
			},
			Image: &discordgo.MessageEmbedImage{
				URL: g.BannerURL(""),
			},
		})
		if err != nil {
			logger.Sugar.Warnw("failed to send guild join notification", "error", err)
		}
	}
}

func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	var message *discordgo.Message
	auctionMap := map[string]any{}
	re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)

	if !strings.HasPrefix(m.Content, "/") {
		return
	}

	args := strings.Split(strings.TrimPrefix(m.Content, "/"), " ")

	switch args[0] {
	case "bid":
		log := logger.Sugar.With(
			"channel_id", m.ChannelID,
			"guild_id", m.GuildID,
			"user_id", m.Author.ID,
		)

		// Fetch auction data to verify ownership
		result := database.DB.Model(database.Auction{}).First(&auctionMap, m.ChannelID)
		if result.Error != nil {
			log.Warnw("failed to fetch auction data for bid", "error", result.Error)
			return
		}

		// Verify this bot owns the auction (FIX: now also returns if we can't verify)
		m2, err := s.ChannelMessage(m.ChannelID, auctionMap["message_id"].(string))
		if err != nil {
			log.Warnw("could not fetch auction message for bid verification", "error", err)
			return // Don't process bid if we can't verify ownership
		}
		if m2.Author.ID != s.State.User.ID {
			log.Debugw("ignoring bid - auction owned by different bot instance",
				"owner_bot_id", m2.Author.ID,
				"our_bot_id", s.State.User.ID,
			)
			return
		}

		if len(args) != 2 {
			_, err := h.ErrorMessage(s, m.ChannelID, fmt.Sprintf("Invalid number of arguments passed. Need 2, used %d", len(args)))
			if err != nil {
				log.Warnw("failed to send argument error message", "error", err)
			}
			return
		}

		bidAmount, err := strconv.ParseFloat(strings.Replace(re.FindString(args[1]), ",", "", -1), 64)
		if err != nil {
			log.Warnw("failed to parse bid amount", "error", err, "input", args[1])
			h.ErrorMessage(s, m.ChannelID, "Invalid bid amount. Please enter a valid number.")
			return
		}

		member, err := s.GuildMember(m.GuildID, m.Author.ID)
		if err != nil {
			log.Errorw("failed to get guild member", "error", err)
			_, err = h.ErrorMessage(s, m.ChannelID, "Unable to verify your membership. Please try again.")
			if err != nil {
				log.Warnw("failed to send member error message", "error", err)
			}
			return
		}

		err = c.AuctionBidPlace(s, bidAmount, member, m.ChannelID, m.GuildID)
		if err != nil {
			log.Warnw("bid placement failed", "error", err, "amount", bidAmount)
			message, err = h.ErrorMessage(s, m.ChannelID, err.Error())
			if err != nil {
				log.Warnw("failed to send bid error message", "error", err)
			}
			time.Sleep(30 * time.Second)
			if message != nil {
				err = s.ChannelMessageDelete(m.ChannelID, message.ID)
				if err != nil {
					log.Debugw("failed to delete error message", "error", err)
				}
			}
			err = s.ChannelMessageDelete(m.ChannelID, m.ID)
			if err != nil {
				log.Debugw("failed to delete user bid message", "error", err)
			}
		} else {
			log.Infow("bid placed successfully", "amount", bidAmount, "username", m.Author.Username)
			message, err = h.SuccessMessage(s, m.ChannelID, h.PresetResponse{
				Title:     "Bid has been successfully placed!",
				Reference: m.Reference(),
			})
			if err != nil {
				log.Warnw("failed to send bid success message", "error", err)
			}
			time.Sleep(30 * time.Second)
			if message != nil {
				err = s.ChannelMessageDelete(m.ChannelID, message.ID)
				if err != nil {
					log.Debugw("failed to delete success message", "error", err)
				}
			}
			err = s.ChannelMessageDelete(m.ChannelID, m.ID)
			if err != nil {
				log.Debugw("failed to delete user bid message", "error", err)
			}
		}

	default:
		return
	}
}

func GuildMemberChunkHandler(s *discordgo.Session, g *discordgo.GuildMembersChunk) {
	log := logger.Sugar.With("nonce", g.Nonce, "guild_id", g.GuildID)
	log.Debug("guild member chunk received")

	if f, ok := guildMembersMap[strings.Split(g.Nonce, ":")[0]]; ok {
		err := f(s, g)
		if err != nil {
			log.Errorw("guild member chunk handler failed", "error", err)
			_, err = h.FollowUpErrorResponse(s, h.ReadChunkData(strings.Split(g.Nonce, ":")[1])["interaction"].(*discordgo.InteractionCreate), "An error occurred processing member data. Please try again or contact support.")
			if err != nil {
				log.Warnw("failed to send followup error response", "error", err)
			}
		}
	} else {
		log.Error("guild member chunk handler not found")
	}
}
