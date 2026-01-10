package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/billingportal/session"
	"github.com/stripe/stripe-go/v72/sub"
	"gitlab.com/logan9312/discord-auction-bot/commands"
	"gitlab.com/logan9312/discord-auction-bot/config"
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

	// Check if user has premium (uses existing function)
	isPremium := commands.CheckPremiumUser(userID)

	// Query Stripe for user's subscriptions
	var subscriptions []SubscriptionInfo

	query, err := commands.BuildStripeQuery("discord_id", userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	params := &stripe.SubscriptionSearchParams{}
	params.Query = *stripe.String(query)
	iter := sub.Search(params)

	for iter.Next() {
		subscription := iter.Subscription()
		subscriptions = append(subscriptions, SubscriptionInfo{
			ID:               subscription.ID,
			Status:           string(subscription.Status),
			CurrentPeriodEnd: subscription.CurrentPeriodEnd,
			GuildID:          subscription.Metadata["guild_id"],
			PlanName:         "Prize Bot Premium",
		})
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
	iter := sub.Search(params)

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
