package mailerlite

import "github.com/joeblew999/ubuntu-website/internal/cli"

// Version is the package version.
const Version = "v0.1.0"

// PackageDoc returns the package documentation for registry publishing.
func PackageDoc() *cli.DocBuilder {
	return cli.NewDoc("mailerlite", Version).
		Description("Go client library and CLI for the MailerLite API. Manage subscribers, groups, and email campaigns.").
		Repo("https://github.com/joeblew999/ubuntu-website").
		HasBinary().
		Feature("**Full API Client** - Subscribers, groups, forms, webhooks, automations").
		Feature("**CLI Tool** - Manage subscribers from command line").
		Feature("**Webhook Server** - Receive real-time events from MailerLite").
		Feature("**GitHub Releases** - Query release info for email templates").
		Command("subscribers list", "List all subscribers", "mailerlite subscribers list").
		Command("subscribers add", "Add subscriber by email", "mailerlite subscribers add user@example.com").
		Command("groups list", "List subscriber groups", "mailerlite groups list").
		Command("stats", "Show account statistics", "mailerlite stats").
		Command("server", "Start webhook server", "mailerlite server --port 8086").
		Example("Library Usage", "go", `
package main

import "github.com/joeblew999/ubuntu-website/internal/mailerlite"

func main() {
    client := mailerlite.NewClient(os.Getenv("MAILERLITE_API_KEY"))

    // List subscribers
    subs, _ := client.ListSubscribers(ctx, nil)

    // Add subscriber
    client.AddSubscriber(ctx, &mailerlite.CreateSubscriberRequest{
        Email: "user@example.com",
    })
}
`)
}

// Dashboard URLs
const (
	DashboardURL    = "https://dashboard.mailerlite.com"
	APIKeyURL       = "https://dashboard.mailerlite.com/integrations/api"
	SubscribersURL  = "https://dashboard.mailerlite.com/subscribers"
	GroupsURL       = "https://dashboard.mailerlite.com/subscribers/groups"
	FormsURL        = "https://dashboard.mailerlite.com/forms"
	AutomationsURL  = "https://dashboard.mailerlite.com/automations"
	CampaignsURL    = "https://dashboard.mailerlite.com/campaigns"
	WebhooksURL     = "https://dashboard.mailerlite.com/integrations/webhooks"
	IntegrationsURL = "https://dashboard.mailerlite.com/integrations"

	// API Base URL
	APIBaseURL = "https://connect.mailerlite.com/api"
)

// Webhook Event Types
const (
	EventSubscriberCreated             = "subscriber.created"
	EventSubscriberUpdated             = "subscriber.updated"
	EventSubscriberUnsubscribed        = "subscriber.unsubscribed"
	EventSubscriberAddedToGroup        = "subscriber.added_to_group"
	EventSubscriberRemovedFromGroup    = "subscriber.removed_from_group"
	EventSubscriberBounced             = "subscriber.bounced"
	EventSubscriberAutomationTriggered = "subscriber.automation_triggered"
	EventSubscriberAutomationComplete  = "subscriber.automation_complete"
	EventCampaignSent                  = "campaign.sent"
	EventCampaignOpened                = "campaign.opened"
	EventCampaignClicked               = "campaign.clicked"
)

// AllWebhookEvents is a list of all available webhook event types.
var AllWebhookEvents = []string{
	EventSubscriberCreated,
	EventSubscriberUpdated,
	EventSubscriberUnsubscribed,
	EventSubscriberAddedToGroup,
	EventSubscriberRemovedFromGroup,
	EventSubscriberBounced,
	EventSubscriberAutomationTriggered,
	EventSubscriberAutomationComplete,
	EventCampaignSent,
	EventCampaignOpened,
	EventCampaignClicked,
}
