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
		fmt.Println(i.ApplicationCommandData().Name, "is being run by:", i.Member.User.Username)
		if f, ok := commandMap[i.ApplicationCommandData().Name]; ok {
			err := f(s, i)
			if err != nil {
				err2 := h.ErrorResponse(s, i, err.Error())
				if err2 != nil {
					fmt.Println(err2)
					_, err = h.FollowUpErrorResponse(s, i, err.Error())
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		} else {
			h.ErrorResponse(s, i, "Command response has not been set properly, please contact Logan to fix")
		}
	case discordgo.InteractionMessageComponent:
		fmt.Println(i.MessageComponentData().CustomID, "is being run by:", i.Member.User.Username)
		if f, ok := buttonMap[strings.Split(i.MessageComponentData().CustomID, ":")[0]]; ok {
			err := f(s, i)
			if err != nil {
				err2 := h.ErrorResponse(s, i, err.Error())
				if err2 != nil {
					fmt.Println(err2)
					_, err = h.FollowUpErrorResponse(s, i, err.Error())
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		} else {
			h.ErrorResponse(s, i, "Button response has not been set properly, please contact Logan to fix")
		}
	case discordgo.InteractionApplicationCommandAutocomplete:
		if f, ok := autoCompleteMap[i.ApplicationCommandData().Name]; ok {
			err := f(s, i)
			if err != nil {
				err = h.ErrorResponse(s, i, err.Error())
				if err != nil {
					fmt.Println(err)
				}
			}
		} else {
			h.ErrorResponse(s, i, "AutoComplete response has not been set properly, please contact Logan to fix")
		}
	case discordgo.InteractionModalSubmit:
		if i.ModalSubmitData().CustomID == "whitelabel_token" {
			err := WhitelabelTokenModal(s, i)
			if err != nil {
				err2 := h.ErrorResponse(s, i, err.Error())
				if err2 != nil {
					fmt.Println(err2)
					_, err = h.FollowUpErrorResponse(s, i, err.Error())
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		} else if f, ok := modalSubmitMap[i.ModalSubmitData().CustomID]; ok {
			err := f(s, i)
			if err != nil {
				err2 := h.ErrorResponse(s, i, err.Error())
				if err2 != nil {
					fmt.Println(err2)
					_, err = h.FollowUpErrorResponse(s, i, err.Error())
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		} else {
			h.ErrorResponse(s, i, "Modal Submit response has not been set properly, please contact Logan to fix")
		}
	}
}

func ReadyHandler(s *discordgo.Session, i *discordgo.Ready) {
	_, err := s.ChannelMessageSend("943175605858496602", "Bot has finished restarting")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Bot is ready")
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
			fmt.Println(err)
		}
	}
}

func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	var message *discordgo.Message
	var err error
	auctionMap := map[string]any{}
	re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)

	if !strings.HasPrefix(m.Content, "/") {
		return
	}

	args := strings.Split(strings.TrimPrefix(m.Content, "/"), " ")

	switch args[0] {
	case "bid":
		//FIXME This was just a quick fix to prevent both bots trying to place a bit. Rework if needed once I can save auctions with bot ID
		result := database.DB.Model(database.Auction{}).First(&auctionMap, m.ChannelID)
		if result.Error != nil {
			fmt.Println("Error fetching auction data from the database. Error Message:", result.Error.Error())
			return
		}
		m2, err := s.ChannelMessage(m.ChannelID, auctionMap["message_id"].(string))
		if err == nil {
			if m2.Author.ID != s.State.User.ID {
				return
			}
		}

		if len(args) != 2 {
			message, err = h.ErrorMessage(s, m.ChannelID, fmt.Sprintf("Invalid number of arguments passed. Need 2, used %d", len(args)))
			if err != nil {
				fmt.Println(err)
				return
			}
			return
		}

		bidAmount, err := strconv.ParseFloat(strings.Replace(re.FindString(args[1]), ",", "", -1), 64)
		if err != nil {
			h.ErrorMessage(s, m.ChannelID, err.Error())
			fmt.Println(err)
			return
		}

		member, err := s.GuildMember(m.GuildID, m.Author.ID)
		if err != nil {
			fmt.Println(err)
			_, err = h.ErrorMessage(s, m.ChannelID, err.Error())
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		err = c.AuctionBidPlace(s, bidAmount, member, m.ChannelID, m.GuildID)
		if err != nil {
			fmt.Println(err)
			message, err = h.ErrorMessage(s, m.ChannelID, err.Error())
			if err != nil {
				fmt.Println(err)
			}
			time.Sleep(30 * time.Second)
			err = s.ChannelMessageDelete(m.ChannelID, message.ID)
			if err != nil {
				fmt.Println(err)
			}
			err = s.ChannelMessageDelete(m.ChannelID, m.ID)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			message, err = h.SuccessMessage(s, m.ChannelID, h.PresetResponse{
				Title:     "Bid has been successfully placed!",
				Reference: m.Reference(),
			})
			if err != nil {
				fmt.Println(err)
			}
			time.Sleep(30 * time.Second)
			err = s.ChannelMessageDelete(m.ChannelID, message.ID)
			if err != nil {
				fmt.Println(err)
			}
			err = s.ChannelMessageDelete(m.ChannelID, m.ID)
			if err != nil {
				fmt.Println(err)
			}
		}

	default:
		return
	}

	if message != nil {
		if err != nil {
			fmt.Println(err)
		}
	}
}

func GuildMemberChunkHandler(s *discordgo.Session, g *discordgo.GuildMembersChunk) {
	fmt.Println(g.Nonce)
	if f, ok := guildMembersMap[strings.Split(g.Nonce, ":")[0]]; ok {
		err := f(s, g)
		if err != nil {
			fmt.Println(err)
			_, err = h.FollowUpErrorResponse(s, h.ReadChunkData(strings.Split(g.Nonce, ":")[1])["interaction"].(*discordgo.InteractionCreate), err.Error())
			if err != nil {
				fmt.Println(err)
			}
		}
	} else {
		fmt.Println("GuildMemberChunkHandler has not been set properly.")
	}
}
