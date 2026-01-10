package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/webhook"
	"gitlab.com/logan9312/discord-auction-bot/commands"
	"gitlab.com/logan9312/discord-auction-bot/config"
	"gitlab.com/logan9312/discord-auction-bot/logger"
)

// StripeWebhook handles incoming Stripe webhook events
// All subscription-related events trigger a full sync for the affected user
func StripeWebhook(c echo.Context) error {
	// Read the raw body for signature verification
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		logger.Sugar.Warnw("failed to read webhook body", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to read request body"})
	}

	// Verify webhook signature
	sigHeader := c.Request().Header.Get("Stripe-Signature")
	if sigHeader == "" {
		logger.Sugar.Warn("missing Stripe-Signature header")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing signature"})
	}

	// If webhook secret is not configured, skip verification (for testing)
	var event stripe.Event
	if config.C.StripeWebhookSecret != "" {
		event, err = webhook.ConstructEvent(body, sigHeader, config.C.StripeWebhookSecret)
		if err != nil {
			logger.Sugar.Warnw("webhook signature verification failed", "error", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid signature"})
		}
	} else {
		// No webhook secret configured - parse event without verification
		// This should only be used in development
		logger.Sugar.Warn("STRIPE_WEBHOOK_SECRET not configured - skipping signature verification")
		if err := json.Unmarshal(body, &event); err != nil {
			logger.Sugar.Warnw("failed to parse webhook event", "error", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid event"})
		}
	}

	log := logger.Sugar.With("event_type", event.Type, "event_id", event.ID)
	log.Info("received Stripe webhook")

	// Extract discord_id from the event and trigger sync
	discordUserID := extractDiscordUserID(event)
	if discordUserID != "" {
		// Sync in background to return 200 quickly
		go func() {
			if err := commands.SyncUserSubscriptions(discordUserID); err != nil {
				logger.Sugar.Warnw("webhook sync failed",
					"event_type", event.Type,
					"discord_user_id", discordUserID,
					"error", err,
				)
			} else {
				logger.Sugar.Infow("webhook sync completed",
					"event_type", event.Type,
					"discord_user_id", discordUserID,
				)
			}
		}()
	} else {
		log.Warn("could not extract discord_user_id from event")
	}

	// Always return 200 OK quickly to acknowledge receipt
	return c.JSON(http.StatusOK, map[string]string{"status": "received"})
}

// extractDiscordUserID attempts to extract the discord_id from various Stripe event types
func extractDiscordUserID(event stripe.Event) string {
	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err == nil {
			// Check subscription metadata first
			if session.Subscription != nil {
				// We need to fetch the full subscription to get metadata
				// For now, return empty and let the next event handle it
				// The subscription.created event will fire right after
			}
		}
		return ""

	case "customer.subscription.created",
		"customer.subscription.updated",
		"customer.subscription.deleted",
		"customer.subscription.paused",
		"customer.subscription.resumed",
		"customer.subscription.pending_update_applied",
		"customer.subscription.pending_update_expired",
		"customer.subscription.trial_will_end":
		var sub stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &sub); err == nil {
			if discordID, ok := sub.Metadata["discord_id"]; ok {
				return discordID
			}
		}

	case "invoice.paid",
		"invoice.payment_failed",
		"invoice.payment_action_required":
		var invoice stripe.Invoice
		if err := json.Unmarshal(event.Data.Raw, &invoice); err == nil {
			// Invoice has subscription attached, get metadata from there
			if invoice.Subscription != nil {
				// The subscription object in invoice might not have full metadata
				// We'd need to fetch it - for now, we rely on subscription events
			}
		}
		return ""
	}

	return ""
}
