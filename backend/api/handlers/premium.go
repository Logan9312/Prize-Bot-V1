package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/billingportal/session"
	"github.com/stripe/stripe-go/v82/subscription"
	"gitlab.com/logan9312/discord-auction-bot/commands"
	"gitlab.com/logan9312/discord-auction-bot/config"
	"gitlab.com/logan9312/discord-auction-bot/database"
)

// UserPremiumStatusResponse represents the user's premium status
type UserPremiumStatusResponse struct {
	IsPremium     bool               `json:"is_premium"`
	Subscriptions []SubscriptionInfo `json:"subscriptions"`
}

// SubscriptionInfo contains details about a subscription
type SubscriptionInfo struct {
	ID               string `json:"id"`
	Status           string `json:"status"`
	CurrentPeriodEnd int64  `json:"current_period_end"`
	GuildID          string `json:"guild_id,omitempty"`
	PlanName         string `json:"plan_name"`
}

// GuildPremiumStatusResponse represents a guild's premium status
type GuildPremiumStatusResponse struct {
	IsPremium bool   `json:"is_premium"`
	GuildID   string `json:"guild_id"`
}

// BillingPortalResponse contains the portal URL
type BillingPortalResponse struct {
	URL string `json:"url"`
}

// GetUserPremiumStatus returns the premium status and subscriptions for the current user
func GetUserPremiumStatus(c echo.Context) error {
	userID := c.Get("user_id").(string)

	// Check if user has premium (uses existing function - now checks local cache first)
	isPremium := commands.CheckPremiumUser(userID)

	// Try local cache first for subscription list
	var localSubs []database.Subscription
	database.DB.Where("discord_user_id = ?", userID).Find(&localSubs)

	var subscriptions []SubscriptionInfo

	if len(localSubs) > 0 {
		// Use cached subscriptions
		for _, sub := range localSubs {
			subscriptions = append(subscriptions, SubscriptionInfo{
				ID:               sub.ID,
				Status:           sub.Status,
				CurrentPeriodEnd: sub.CurrentPeriodEnd.Unix(),
				GuildID:          sub.GuildID,
				PlanName:         "Prize Bot Premium",
			})
		}
	} else {
		// Fallback to Stripe API if no local cache
		query, err := commands.BuildStripeQuery("discord_id", userID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		}

		params := &stripe.SubscriptionSearchParams{}
		params.Query = *stripe.String(query)
		iter := subscription.Search(params)

		for iter.Next() {
			sub := iter.Subscription()
			// In v82, CurrentPeriodEnd moved to SubscriptionItem
			var periodEnd int64
			if sub.Items != nil && len(sub.Items.Data) > 0 {
				periodEnd = sub.Items.Data[0].CurrentPeriodEnd
			}
			subscriptions = append(subscriptions, SubscriptionInfo{
				ID:               sub.ID,
				Status:           string(sub.Status),
				CurrentPeriodEnd: periodEnd,
				GuildID:          sub.Metadata["guild_id"],
				PlanName:         "Prize Bot Premium",
			})
		}

		// Trigger async sync if we had to use Stripe API
		if len(subscriptions) > 0 {
			go commands.SyncUserSubscriptions(userID)
		}
	}

	return c.JSON(http.StatusOK, UserPremiumStatusResponse{
		IsPremium:     isPremium,
		Subscriptions: subscriptions,
	})
}

// GetGuildPremiumStatus returns the premium status for a specific guild
func GetGuildPremiumStatus(c echo.Context) error {
	guildID := c.Param("guildId")

	if err := verifyGuildAccess(c, guildID); err != nil {
		return err
	}

	// Use existing function from commands/premium.go
	isPremium := commands.CheckPremiumGuild(guildID)

	return c.JSON(http.StatusOK, GuildPremiumStatusResponse{
		IsPremium: isPremium,
		GuildID:   guildID,
	})
}

// SyncAfterCheckout is called by the frontend after successful checkout
// This handles the race condition where user returns before webhook fires
func SyncAfterCheckout(c echo.Context) error {
	userID := c.Get("user_id").(string)

	// Sync subscriptions from Stripe
	if err := commands.SyncUserSubscriptions(userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to sync subscription data",
		})
	}

	// Return updated premium status
	return GetUserPremiumStatus(c)
}

// CreateBillingPortalSession creates a Stripe billing portal session
func CreateBillingPortalSession(c echo.Context) error {
	userID := c.Get("user_id").(string)

	// Find the customer ID from Stripe subscriptions
	query, err := commands.BuildStripeQuery("discord_id", userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	params := &stripe.SubscriptionSearchParams{}
	params.Query = *stripe.String(query)
	iter := subscription.Search(params)

	var customerID string
	for iter.Next() {
		subscription := iter.Subscription()
		customerID = subscription.Customer.ID
		break
	}

	if customerID == "" {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "No subscription found. Please subscribe first.",
		})
	}

	// Create billing portal session
	portal, err := session.New(&stripe.BillingPortalSessionParams{
		Customer:  &customerID,
		ReturnURL: stripe.String(config.C.FrontendURL + "/dashboard/subscription"),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create billing portal session",
		})
	}

	return c.JSON(http.StatusOK, BillingPortalResponse{
		URL: portal.URL,
	})
}
