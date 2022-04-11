package subscriptions

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/stripe/stripe-go/v72/sub"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
)

var PremiumCommand = discordgo.ApplicationCommand{
	Name:        "premium",
	Description: "Premium Bot Info",
	Options:     []*discordgo.ApplicationCommandOption{},
}

func Premium(s *discordgo.Session, i *discordgo.InteractionCreate) {

	session, err := PremiumSession(i.Member.User.ID)
	if err != nil {
		h.ErrorResponse(s, i, err.Error())
		return
	}

	h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "Subscribe",
		Description: "This is just a test for now",
		Fields:      []*discordgo.MessageEmbedField{},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label: "Subscribe",
						Style: discordgo.LinkButton,
						URL:   session.URL,
					},
				},
			},
		},
	})

	fmt.Println("Current Subscriptions")
	ListSubscriptions()
}

func PremiumSession(userID string) (*stripe.CheckoutSession, error) {

	fmt.Println(stripe.Key)
	fmt.Println(userID)

	params := &stripe.CheckoutSessionParams{
		CancelURL:  stripe.String("https://discord.gg/YBRvZ3mRtb"),
		SuccessURL: stripe.String("https://discord.gg/YBRvZ3mRtb"),
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String("price_1KYE0EKpOiJyve6nT9Qo9IfN"),
				Quantity: stripe.Int64(1),
			},
		},
	}
	params.SubscriptionData.AddMetadata("discord_id", userID)

	return session.New(params)
}

func ListSubscriptions() {
	params := &stripe.SubscriptionListParams{}
	i := sub.List(params)
	for i.Next() {
		s := i.Subscription()
		fmt.Println("Discord ID:", s.Metadata)
	}
}
