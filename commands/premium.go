package commands

import (
	"fmt"

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
			Name:        "option",
			Description: "Get information or activate a premium subscription",
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "info",
					Value: "info",
				},
				{
					Name:  "activate",
					Value: "activate",
				},
			},
			Type:     discordgo.ApplicationCommandOptionString,
			Required: true,
		},
	},
	DMPermission: new(bool),
}

func Premium(s *discordgo.Session, i *discordgo.InteractionCreate) {

	options := h.ParseSlashCommand(i)
	switch options["option"] {
	case "info":
		PremiumInfo(s, i)
	case "activate":
		PremiumActivate(s, i)
	}
}

func PremiumInfo(s *discordgo.Session, i *discordgo.InteractionCreate) {

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
		h.ErrorResponse(s, i, err.Error())
		return
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
			ReturnURL: stripe.String("https://aftmgaming.wordpress.com/"),
		})
		if err != nil {
			h.ErrorResponse(s, i, err.Error())
			return
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
		fmt.Println(err)
		return
	}

	fmt.Println("Current Subscriptions")
	ListSubscriptions(s)
}

func PremiumActivate(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if CheckPremiumGuild(i.GuildID) {
		h.SuccessResponse(s, i, h.PresetResponse{
			Title:       "**Premium Server**",
			Description: "Premium has already been activated for this server!",
		})
		return
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
				return
			}
		}
	}

	err := h.ErrorResponse(s, i, "No valid subscription found! Please subscribe to premium first using `/premium info`")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func PremiumSession(userID, customerID string) (*stripe.CheckoutSession, error) {

	params := &stripe.CheckoutSessionParams{
		CancelURL:  stripe.String("https://discord.gg/YBRvZ3mRtb"),
		SuccessURL: stripe.String("https://aftmgaming.wordpress.com/"),
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

func ListSubscriptions(s *discordgo.Session) {
	params := &stripe.SubscriptionListParams{}
	activeMap := map[string]bool{}
	i := sub.List(params)
	for i.Next() {
		subscription := i.Subscription()
		fmt.Println(subscription.Metadata)
		fmt.Println(subscription.Status)
		fmt.Println("Metadata:", subscription.Metadata["discord_id"])
		if subscription.Status == stripe.SubscriptionStatusActive {
			activeMap[subscription.Metadata["discord_id"]] = true
		} else if activeMap[subscription.Metadata["discord_id"]] != true {
			activeMap[subscription.Metadata["discord_id"]] = false
		}
	}

	fmt.Println(activeMap)

	for userID, active := range activeMap {
		if active {
			fmt.Println("Adding role for:", userID)
			err := s.GuildMemberRoleAdd("885228283573178408", userID, "942927890100682752")
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Removing role for:", userID)
			err := s.GuildMemberRoleRemove("885228283573178408", userID, "942927890100682752")
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func CheckPremiumGuild(guildID string) (status bool) {

	params := &stripe.SubscriptionSearchParams{}
	params.Query = *stripe.String(fmt.Sprintf("status:'active' AND metadata['guild_id']:'%s'", guildID))
	iter := sub.Search(params)

	for iter.Next() {
		subscription := iter.Subscription()
		if subscription.Status == stripe.SubscriptionStatusActive {
			status = true
		}
	}

	return
}
