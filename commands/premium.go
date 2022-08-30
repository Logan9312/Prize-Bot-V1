package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/billingportal/session"
	checkout "github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/stripe/stripe-go/v72/sub"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
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

var PremiumServers = map[string]string{
	"915767892467920967": "Testing",
	"626094990984216586": "Aftermath",
}

var PremiumUsers = map[string]string{
	"280812467775471627": "Logan",
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

	customerID := ""

	params := &stripe.SubscriptionSearchParams{}
	params.Query = *stripe.String(fmt.Sprintf("status:'active' AND metadata['discord_id']:'%s'", i.Member.User.ID))
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
			ReturnURL: stripe.String("http://www.aftmgaming.com/auction-bot/success"),
		})
		if err != nil {
			return err
		}
		buttons = append(buttons, discordgo.Button{
			Label: "Customer Portal",
			Style: discordgo.LinkButton,
			Emoji: discordgo.ComponentEmoji{},
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

	fmt.Println("Current Subscriptions")
	return nil
}

func PremiumActivate(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	if CheckPremiumGuild(i.GuildID) {
		h.SuccessResponse(s, i, h.PresetResponse{
			Title:       "**Premium Server**",
			Description: "Premium has already been activated for this server!",
		})
		return nil
	}

	params := &stripe.SubscriptionSearchParams{}
	params.Query = *stripe.String(fmt.Sprintf("status:'active' AND metadata['discord_id']:'%s'", i.Member.User.ID))
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
					Content:     "",
					Title:       "**Server Activated!**",
					Description: "Your server has successfully been activated and should now work with premium features!",
				})
				return nil
			}
		}
	}

	err := h.ErrorResponse(s, i, "No valid subscription found! Please subscribe to premium first using `/premium info`")
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func PremiumSession(userID, customerID string) (*stripe.CheckoutSession, error) {

	params := &stripe.CheckoutSessionParams{
		CancelURL:  stripe.String("https://discord.gg/YBRvZ3mRtb"),
		SuccessURL: stripe.String("http://www.aftmgaming.com/auction-bot/success"),
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String("price_1KYE0EKpOiJyve6nT9Qo9IfN"),
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

	const premiumRole = "942927890100682752"
	const supportServer = "885228283573178408"

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
				fmt.Println("Adding role for:", userID)
				err := s.GuildMemberRoleAdd(supportServer, userID, premiumRole)
				if err != nil {
					fmt.Println(err)
				}
			} else if !active && hasRole {
				fmt.Println("Removing role for:", userID)
				err := s.GuildMemberRoleRemove(supportServer, userID, premiumRole)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
		time.Sleep(5 * time.Minute)
	}
}

func CheckPremiumUser(userID string) bool {
	if PremiumUsers[userID] != "" {
		return true
	}

	params := &stripe.SubscriptionSearchParams{}
	params.Query = *stripe.String(fmt.Sprintf("status:'active' AND metadata['user_id']:'%s'", userID))
	iter := sub.Search(params)

	for iter.Next() {
		subscription := iter.Subscription()
		if subscription.Status == stripe.SubscriptionStatusActive {
			return true
		}
	}

	return false
}

func CheckPremiumGuild(guildID string) (status bool) {

	if PremiumServers[guildID] != "" {
		return true
	}

	params := &stripe.SubscriptionSearchParams{}
	params.Query = *stripe.String(fmt.Sprintf("status:'active' AND metadata['guild_id']:'%s'", guildID))
	iter := sub.Search(params)

	for iter.Next() {
		subscription := iter.Subscription()
		if subscription.Status == stripe.SubscriptionStatusActive {
			return true
		}
	}

	return false
}
