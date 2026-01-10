package commands

import (
	"fmt"
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/billingportal/session"
	checkout "github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/subscription"
	"gitlab.com/logan9312/discord-auction-bot/config"
	"gitlab.com/logan9312/discord-auction-bot/database"
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

// BuildStripeQuery safely builds a Stripe search query with validated Discord ID
func BuildStripeQuery(field, discordID string) (string, error) {
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
	query, err := BuildStripeQuery("discord_id", i.Member.User.ID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	params := &stripe.SubscriptionSearchParams{}
	params.Query = *stripe.String(query)
	iter := subscription.Search(params)

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
	query, err := BuildStripeQuery("discord_id", i.Member.User.ID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	params := &stripe.SubscriptionSearchParams{}
	params.Query = *stripe.String(query)
	iter := subscription.Search(params)

	for iter.Next() {
		sub := iter.Subscription()
		if sub.Status == stripe.SubscriptionStatusActive {
			if _, ok := sub.Metadata["guild_id"]; !ok {
				sub.Metadata["guild_id"] = i.GuildID
				subscription.Update(sub.ID, &stripe.SubscriptionParams{
					Params: stripe.Params{
						Metadata: sub.Metadata,
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
			// Restart the goroutine after a delay to avoid tight panic loops
			time.Sleep(30 * time.Second)
			go SetRoles(s)
		}
	}()

	premiumRole := config.C.PremiumRoleID
	supportServer := config.C.SupportServerID

	for {
		// Query local database cache instead of Stripe API
		// This is kept up-to-date by webhooks and sync functions
		activeMap := map[string]bool{}

		// Get all unique discord_user_ids and their best status
		var subscriptions []database.Subscription
		if err := database.DB.Find(&subscriptions).Error; err != nil {
			logger.Sugar.Warnw("failed to query subscriptions for role sync", "error", err)
			time.Sleep(5 * time.Minute)
			continue
		}

		// Build activeMap: user is active if ANY of their subscriptions is active
		for _, sub := range subscriptions {
			if sub.DiscordUserID == "" {
				continue
			}
			if sub.Status == "active" {
				activeMap[sub.DiscordUserID] = true
			} else if !activeMap[sub.DiscordUserID] {
				// Only set to false if not already true (active takes precedence)
				activeMap[sub.DiscordUserID] = false
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

	// Check local cache first
	var count int64
	database.DB.Model(&database.Subscription{}).
		Where("discord_user_id = ? AND status = ?", userID, "active").
		Count(&count)
	if count > 0 {
		return true
	}

	// Fallback: Query Stripe API (for race condition edge cases)
	query, err := BuildStripeQuery("discord_id", userID)
	if err != nil {
		logger.Sugar.Warnw("invalid user ID for premium check", "user_id", userID, "error", err)
		return false
	}

	params := &stripe.SubscriptionSearchParams{}
	params.Query = *stripe.String(query)
	iter := subscription.Search(params)

	for iter.Next() {
		subscription := iter.Subscription()
		if subscription.Status == stripe.SubscriptionStatusActive {
			// Found in Stripe but not in local cache - trigger async sync
			go func() {
				if err := SyncUserSubscriptions(userID); err != nil {
					logger.Sugar.Warnw("async sync failed", "user_id", userID, "error", err)
				}
			}()
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

	// Check local cache first
	var count int64
	database.DB.Model(&database.Subscription{}).
		Where("guild_id = ? AND status = ?", guildID, "active").
		Count(&count)
	if count > 0 {
		return true
	}

	// Fallback: Query Stripe API (for race condition edge cases)
	query, err := BuildStripeQuery("guild_id", guildID)
	if err != nil {
		logger.Sugar.Warnw("invalid guild ID for premium check", "guild_id", guildID, "error", err)
		return false
	}

	params := &stripe.SubscriptionSearchParams{}
	params.Query = *stripe.String(query)
	iter := subscription.Search(params)

	for iter.Next() {
		subscription := iter.Subscription()
		if subscription.Status == stripe.SubscriptionStatusActive {
			// Found in Stripe but not in local cache - trigger async sync
			go func(discordID string) {
				if discordID != "" {
					if err := SyncUserSubscriptions(discordID); err != nil {
						logger.Sugar.Warnw("async sync failed", "discord_id", discordID, "error", err)
					}
				}
			}(subscription.Metadata["discord_id"])
			return true
		}
	}

	return false
}

// SyncUserSubscriptions fetches all subscriptions for a Discord user from Stripe
// and updates the local cache. This is the SINGLE source of truth sync function.
// All webhook events should call this function rather than updating fields selectively.
func SyncUserSubscriptions(discordUserID string) error {
	if err := validateDiscordID(discordUserID); err != nil {
		return fmt.Errorf("invalid discord user ID: %w", err)
	}

	log := logger.Sugar.With("discord_user_id", discordUserID)
	log.Debug("syncing user subscriptions from Stripe")

	// Query Stripe for all subscriptions with this discord_id
	// Note: We query ALL statuses, not just active, to properly sync cancellations
	params := &stripe.SubscriptionSearchParams{}
	params.Query = *stripe.String(fmt.Sprintf("metadata['discord_id']:'%s'", discordUserID))
	iter := subscription.Search(params)

	// Collect all subscription IDs from Stripe
	stripeSubIDs := make(map[string]bool)
	var subscriptionsToUpsert []database.Subscription

	for iter.Next() {
		stripeSub := iter.Subscription()
		stripeSubIDs[stripeSub.ID] = true

		// Get price ID and period info from first item (v82 moved these to SubscriptionItem)
		priceID := ""
		var currentPeriodStart, currentPeriodEnd int64
		if stripeSub.Items != nil && len(stripeSub.Items.Data) > 0 {
			item := stripeSub.Items.Data[0]
			priceID = item.Price.ID
			currentPeriodStart = item.CurrentPeriodStart
			currentPeriodEnd = item.CurrentPeriodEnd
		}

		subscriptionsToUpsert = append(subscriptionsToUpsert, database.Subscription{
			ID:                 stripeSub.ID,
			CustomerID:         stripeSub.Customer.ID,
			DiscordUserID:      discordUserID,
			GuildID:            stripeSub.Metadata["guild_id"],
			Status:             string(stripeSub.Status),
			PriceID:            priceID,
			CurrentPeriodStart: time.Unix(currentPeriodStart, 0),
			CurrentPeriodEnd:   time.Unix(currentPeriodEnd, 0),
			CancelAtPeriodEnd:  stripeSub.CancelAtPeriodEnd,
		})
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("stripe search error: %w", err)
	}

	// Delete local subscriptions that no longer exist in Stripe for this user
	if err := database.DB.Where("discord_user_id = ? AND id NOT IN ?", discordUserID, keysFromMap(stripeSubIDs)).
		Delete(&database.Subscription{}).Error; err != nil {
		log.Warnw("failed to delete stale subscriptions", "error", err)
	}

	// Upsert all subscriptions from Stripe
	for _, sub := range subscriptionsToUpsert {
		if err := database.DB.Save(&sub).Error; err != nil {
			log.Warnw("failed to upsert subscription", "subscription_id", sub.ID, "error", err)
		}
	}

	log.Debugw("sync completed", "subscription_count", len(subscriptionsToUpsert))
	return nil
}

// keysFromMap extracts keys from a map for use in SQL IN clause
func keysFromMap(m map[string]bool) []string {
	if len(m) == 0 {
		return []string{""} // Return placeholder to avoid empty IN clause
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
