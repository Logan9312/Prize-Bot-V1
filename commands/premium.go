package commands

import (
	"fmt"
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/billingportal/session"
	checkout "github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/stripe/stripe-go/v72/sub"
	"gitlab.com/logan9312/discord-auction-bot/config"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
	"gitlab.com/logan9312/discord-auction-bot/logger"
)

var PremiumCommand = discordgo.ApplicationCommand{
	Name:        "premium",
	Description: "Info",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "info",
			Description: "Get information about a premium subscription",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
		},
		{
			Name:        "activate",
			Description: "Activate a premium subscription",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
		},
	},
	DMPermission: new(bool),
}

// discordIDRegex validates Discord snowflake IDs (numeric strings)
var discordIDRegex = regexp.MustCompile(`^\d{17,20}$`)

// validateDiscordID ensures a Discord ID is safe for use in queries
func validateDiscordID(id string) error {
	if !discordIDRegex.MatchString(id) {
		return fmt.Errorf("invalid Discord ID format")
	}
	return nil
}

// buildStripeQuery safely builds a Stripe search query with validated Discord ID
func buildStripeQuery(field, discordID string) (string, error) {
	if err := validateDiscordID(discordID); err != nil {
		return "", err
	}
	return fmt.Sprintf("status:'active' AND metadata['%s']:'%s'", field, discordID), nil
}

func Premium(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	switch i.ApplicationCommandData().Options[0].Name {
	case "info":
		return PremiumInfo(s, i)
	case "activate":
		return PremiumActivate(s, i)
	}
	return fmt.Errorf("Unknown Premium command, please contact support")
}

func PremiumInfo(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	if i.Member == nil || i.Member.User == nil {
		return fmt.Errorf("this command cannot be used in DMs")
	}

	customerID := ""

	// Build query with validated Discord ID to prevent injection
	query, err := buildStripeQuery("discord_id", i.Member.User.ID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	params := &stripe.SubscriptionSearchParams{}
	params.Query = *stripe.String(query)
	iter := sub.Search(params)

	for iter.Next() {
		subscription := iter.Subscription()
		customerID = subscription.Customer.ID
	}

	pSession, err := PremiumSession(i.Member.User.ID, customerID)
	if err != nil {
		return err
	}
	buttons := []discordgo.MessageComponent{
		discordgo.Button{
			Label: "Purchase!",
			Style: discordgo.LinkButton,
			URL:   pSession.URL,
		},
	}

	if customerID != "" {
		portal, err := session.New(&stripe.BillingPortalSessionParams{
			Customer:  &customerID,
			ReturnURL: stripe.String(config.C.SuccessURL),
		})
		if err != nil {
			return err
		}
		buttons = append(buttons, discordgo.Button{
			Label: "Customer Portal",
			Style: discordgo.LinkButton,
			Emoji: &discordgo.ComponentEmoji{},
			URL:   portal.URL,
		})
	}

	err = h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "Auction Bot Premium",
		Description: "The premium version of the bot allows for extra customization and helps support development. It's greatly appreciated but not required",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "🏅 **Show Support**",
				Value:  "Help cover the costs of hosting and maintaining the bot.",
				Inline: true,
			},
			{
				Name:   "💸 **Auctions**",
				Value:  "✔️ Queue auctions in advance.",
				Inline: true,
			},
			{
				Name:   "🎉 **Giveaways**",
				Value:  "✔️ Auto Enter users by role (Coming Soon)",
				Inline: true,
			},
			{
				Name:   "🎁 **Claims**",
				Value:  "✔️ Give a prize to everyone in a role at once.",
				Inline: true,
			},
			{
				Name:   "🤖 **Whitelabelling**",
				Value:  "✔️ Use a custom bot logo and username",
				Inline: true,
			},
			{
				Name:   "🎨 **Customization**",
				Value:  "✔️ Custom Thumbnails (Coming Soon).\n✔️ Custom embed colors (Coming Soon)",
				Inline: true,
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: buttons,
			},
		},
	})
	if err != nil {
		return err
	}

	logger.Sugar.Debug("checking current subscriptions")
	return nil
}

func PremiumActivate(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	if i.Member == nil || i.Member.User == nil {
		return fmt.Errorf("this command cannot be used in DMs")
	}

	if CheckPremiumGuild(i.GuildID) {
		h.SuccessResponse(s, i, h.PresetResponse{
			Title:       "**Premium Server**",
			Description: "Premium has already been activated for this server!",
		})
		return nil
	}

	// Build query with validated Discord ID to prevent injection
	query, err := buildStripeQuery("discord_id", i.Member.User.ID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	params := &stripe.SubscriptionSearchParams{}
	params.Query = *stripe.String(query)
	iter := sub.Search(params)

	for iter.Next() {
		subscription := iter.Subscription()
		if subscription.Status == stripe.SubscriptionStatusActive {
			if _, ok := subscription.Metadata["guild_id"]; !ok {
				subscription.Metadata["guild_id"] = i.GuildID
				sub.Update(subscription.ID, &stripe.SubscriptionParams{
					Params: stripe.Params{
						Metadata: subscription.Metadata,
					},
				})
				h.SuccessResponse(s, i, h.PresetResponse{
					Title:       "**Server Activated!**",
					Description: "Your server has successfully been activated and should now work with premium features!",
				})
				return nil
			}
		}
	}

	err = h.ErrorResponse(s, i, "No unlinked subscription found! Please subscribe to premium first using `/premium info`")
	if err != nil {
		logger.Sugar.Warnw("premium operation error", "error", err)
		return err
	}
	return nil
}

func PremiumSession(userID, customerID string) (*stripe.CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{
		CancelURL:  stripe.String(config.C.CancelURL),
		SuccessURL: stripe.String(config.C.SuccessURL),
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(config.C.StripePriceID),
				Quantity: stripe.Int64(1),
			},
		},
	}
	if customerID != "" {
		params.Customer = stripe.String(customerID)
	}
	params.SubscriptionData = &stripe.CheckoutSessionSubscriptionDataParams{
		Metadata: map[string]string{
			"discord_id": userID,
		},
	}

	return checkout.New(params)
}

func SetRoles(s *discordgo.Session) {
	defer func() {
		if r := recover(); r != nil {
			logger.Sugar.Errorw("panic in SetRoles goroutine",
				"panic", r,
			)
		}
	}()

	premiumRole := config.C.PremiumRoleID
	supportServer := config.C.SupportServerID

	for {
		params := &stripe.SubscriptionListParams{}
		activeMap := map[string]bool{}
		i := sub.List(params)
		for i.Next() {
			subscription := i.Subscription()
			if subscription.Status == stripe.SubscriptionStatusActive {
				activeMap[subscription.Metadata["discord_id"]] = true
			} else if !activeMap[subscription.Metadata["discord_id"]] {
				activeMap[subscription.Metadata["discord_id"]] = false
			}
		}

		for userID, active := range activeMap {
			hasRole := false

			m, err := s.GuildMember(supportServer, userID)
			if err != nil {
				continue
			}
			for _, role := range m.Roles {
				if role == premiumRole {
					hasRole = true
				}
			}

			if active && !hasRole {
				logger.Sugar.Debugw("adding premium role", "user_id", userID)
				err := s.GuildMemberRoleAdd(supportServer, userID, premiumRole)
				if err != nil {
					logger.Sugar.Warnw("premium operation error", "error", err)
				}
			} else if !active && hasRole {
				logger.Sugar.Debugw("removing premium role", "user_id", userID)
				err := s.GuildMemberRoleRemove(supportServer, userID, premiumRole)
				if err != nil {
					logger.Sugar.Warnw("premium operation error", "error", err)
				}
			}
		}
		time.Sleep(5 * time.Minute)
	}
}

func CheckPremiumUser(userID string) bool {
	// Check hardcoded premium users from config
	if config.IsPremiumUser(userID) {
		return true
	}

	// Validate user ID before using in query
	query, err := buildStripeQuery("discord_id", userID)
	if err != nil {
		logger.Sugar.Warnw("invalid user ID for premium check", "user_id", userID, "error", err)
		return false
	}

	params := &stripe.SubscriptionSearchParams{}
	params.Query = *stripe.String(query)
	iter := sub.Search(params)

	for iter.Next() {
		subscription := iter.Subscription()
		if subscription.Status == stripe.SubscriptionStatusActive {
			return true
		}
	}

	return false
}

func CheckPremiumGuild(guildID string) bool {
	// Check hardcoded premium servers from config
	if config.IsPremiumServer(guildID) {
		return true
	}

	// Validate guild ID before using in query
	query, err := buildStripeQuery("guild_id", guildID)
	if err != nil {
		logger.Sugar.Warnw("invalid guild ID for premium check", "guild_id", guildID, "error", err)
		return false
	}

	params := &stripe.SubscriptionSearchParams{}
	params.Query = *stripe.String(query)
	iter := sub.Search(params)

	for iter.Next() {
		subscription := iter.Subscription()
		if subscription.Status == stripe.SubscriptionStatusActive {
			return true
		}
	}

	return false
}
