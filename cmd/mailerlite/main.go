// mailerlite - CLI tool for MailerLite subscriber management
//
// Usage:
//
//	mailerlite subscribers list         List all subscribers
//	mailerlite subscribers count        Count total subscribers
//	mailerlite subscribers get EMAIL    Get subscriber by email
//	mailerlite subscribers add EMAIL    Add/update subscriber
//	mailerlite subscribers delete EMAIL Delete subscriber
//	mailerlite groups list              List all groups
//	mailerlite groups create NAME       Create a group
//	mailerlite groups subscribers ID    List subscribers in a group
//	mailerlite groups assign ID EMAIL   Assign subscriber to group
//	mailerlite groups unassign ID EMAIL Remove subscriber from group
//	mailerlite forms list               List all forms
//	mailerlite forms subscribers ID     List subscribers for a form
//	mailerlite webhooks list            List all webhooks
//	mailerlite webhooks create URL      Create a webhook
//	mailerlite webhooks delete ID       Delete a webhook
//	mailerlite automations list         List all automations
//	mailerlite automations get ID       Get automation details
//	mailerlite stats                    Show account statistics
//	mailerlite open [TARGET]            Open dashboard in browser
//	mailerlite releases latest          Show latest GitHub release
//	mailerlite releases list            List all GitHub releases
//	mailerlite releases urls            Show release URLs for emails
//	mailerlite server [PORT=8086]       Start webhook server for Web3Forms
//
// Flags:
//
//	-github-issue    Output markdown for GitHub issue
//	-v               Verbose output
//	-version         Show version
//
// Environment:
//
//	MAILERLITE_API_KEY    API key (required)
//
// Web3Forms Integration:
//
//	The server command starts an HTTP server that receives webhooks from
//	Web3Forms and adds subscribers to MailerLite. For production, expose
//	the server using a tunnel service (ngrok, cloudflared).
//
//	Flow: Web3Forms form → webhook POST → mailerlite server → MailerLite API
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/mailerlite/mailerlite-go"
)

var version = "dev"

// ============================================================================
// Constants - MailerLite URLs and Endpoints
// ============================================================================

const (
	// Dashboard URLs
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

	// Webhook Event Types
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

// AllWebhookEvents is a list of all available webhook event types
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

// ============================================================================
// GitHub Release URLs - for software delivery emails
// ============================================================================

const (
	// GitHub URLs
	GitHubAPIBase     = "https://api.github.com"
	GitHubReleasesURL = "https://github.com/%s/%s/releases"
	GitHubLatestURL   = "https://github.com/%s/%s/releases/latest"
	GitHubDownloadURL = "https://github.com/%s/%s/releases/download/%s/%s"

	// Default repository (Ubuntu Software)
	DefaultGitHubOwner = "joeblew999"
	DefaultGitHubRepo  = "ubuntu-website"
)

// ReleaseAsset represents a downloadable asset from a GitHub release
type ReleaseAsset struct {
	Name        string `json:"name"`
	DownloadURL string `json:"browser_download_url"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
}

// Release represents a GitHub release
type Release struct {
	TagName     string         `json:"tag_name"`
	Name        string         `json:"name"`
	Body        string         `json:"body"`
	Draft       bool           `json:"draft"`
	Prerelease  bool           `json:"prerelease"`
	CreatedAt   string         `json:"created_at"`
	PublishedAt string         `json:"published_at"`
	Assets      []ReleaseAsset `json:"assets"`
	HTMLURL     string         `json:"html_url"`
}

// GetReleasesURL returns the URL for a repo's releases page
func GetReleasesURL(owner, repo string) string {
	return fmt.Sprintf(GitHubReleasesURL, owner, repo)
}

// GetLatestReleaseURL returns the URL for the latest release
func GetLatestReleaseURL(owner, repo string) string {
	return fmt.Sprintf(GitHubLatestURL, owner, repo)
}

// GetDownloadURL returns the direct download URL for a release asset
func GetDownloadURL(owner, repo, tag, filename string) string {
	return fmt.Sprintf(GitHubDownloadURL, owner, repo, tag, filename)
}

func main() {
	// Global flags
	var (
		githubIssue = flag.Bool("github-issue", false, "Output markdown for GitHub issue")
		verbose     = flag.Bool("v", false, "Verbose output")
		showVersion = flag.Bool("version", false, "Show version")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("mailerlite %s\n", version)
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	// Commands that don't require API key
	noAPIKeyCommands := map[string]bool{
		"releases": true,
		"open":     true,
	}

	// Get API key (only required for MailerLite API commands)
	apiKey := os.Getenv("MAILERLITE_API_KEY")
	cmd := args[0]

	if apiKey == "" && !noAPIKeyCommands[cmd] {
		fmt.Fprintln(os.Stderr, "ERROR: MAILERLITE_API_KEY environment variable is required")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintf(os.Stderr, "Get your API key from: %s\n", APIKeyURL)
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Commands that work without API key: releases, open")
		os.Exit(1)
	}

	var client *mailerlite.Client
	if apiKey != "" {
		client = mailerlite.NewClient(apiKey)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	app := &App{
		client:      client,
		ctx:         ctx,
		githubIssue: *githubIssue,
		verbose:     *verbose,
	}

	if err := app.Run(args); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}

type App struct {
	client      *mailerlite.Client
	ctx         context.Context
	githubIssue bool
	verbose     bool
}

func (a *App) Run(args []string) error {
	cmd := args[0]
	subArgs := args[1:]

	switch cmd {
	case "subscribers":
		return a.handleSubscribers(subArgs)
	case "groups":
		return a.handleGroups(subArgs)
	case "forms":
		return a.handleForms(subArgs)
	case "webhooks":
		return a.handleWebhooks(subArgs)
	case "automations":
		return a.handleAutomations(subArgs)
	case "stats":
		return a.handleStats()
	case "open":
		return a.handleOpen(subArgs)
	case "releases":
		return a.handleReleases(subArgs)
	case "server":
		return a.handleServer(subArgs)
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}

func (a *App) handleSubscribers(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("subscribers subcommand required: list, count, get, add, delete")
	}

	switch args[0] {
	case "list":
		return a.subscribersList()
	case "count":
		return a.subscribersCount()
	case "get":
		if len(args) < 2 {
			return fmt.Errorf("email required: subscribers get EMAIL")
		}
		return a.subscribersGet(args[1])
	case "add":
		if len(args) < 2 {
			return fmt.Errorf("email required: subscribers add EMAIL [NAME]")
		}
		name := ""
		if len(args) >= 3 {
			name = strings.Join(args[2:], " ")
		}
		return a.subscribersAdd(args[1], name)
	case "delete":
		if len(args) < 2 {
			return fmt.Errorf("email required: subscribers delete EMAIL")
		}
		return a.subscribersDelete(args[1])
	default:
		return fmt.Errorf("unknown subscribers subcommand: %s", args[0])
	}
}

func (a *App) subscribersList() error {
	options := &mailerlite.ListSubscriberOptions{
		Limit: 100,
		Page:  1,
	}

	subscribers, _, err := a.client.Subscriber.List(a.ctx, options)
	if err != nil {
		return fmt.Errorf("list subscribers: %w", err)
	}

	if a.githubIssue {
		fmt.Println("## Subscribers")
		fmt.Println()
		fmt.Printf("Total: **%d**\n\n", subscribers.Meta.Total)
		if len(subscribers.Data) > 0 {
			fmt.Println("| Email | Status | Subscribed |")
			fmt.Println("|-------|--------|------------|")
			for _, s := range subscribers.Data {
				fmt.Printf("| %s | %s | %s |\n", s.Email, s.Status, s.SubscribedAt)
			}
		}
		return nil
	}

	fmt.Printf("Total subscribers: %d\n\n", subscribers.Meta.Total)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "EMAIL\tSTATUS\tOPENS\tCLICKS\tSUBSCRIBED")
	for _, s := range subscribers.Data {
		subscribed := ""
		if s.SubscribedAt != "" {
			if t, err := time.Parse(time.RFC3339, s.SubscribedAt); err == nil {
				subscribed = t.Format("2006-01-02")
			}
		}
		fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%s\n", s.Email, s.Status, s.OpensCount, s.ClicksCount, subscribed)
	}
	w.Flush()

	return nil
}

func (a *App) subscribersCount() error {
	count, _, err := a.client.Subscriber.Count(a.ctx)
	if err != nil {
		return fmt.Errorf("count subscribers: %w", err)
	}

	if a.githubIssue {
		fmt.Printf("**Total Subscribers:** %d\n", count.Total)
		return nil
	}

	fmt.Printf("Total subscribers: %d\n", count.Total)
	return nil
}

func (a *App) subscribersGet(email string) error {
	options := &mailerlite.GetSubscriberOptions{
		Email: email,
	}

	subscriber, _, err := a.client.Subscriber.Get(a.ctx, options)
	if err != nil {
		return fmt.Errorf("get subscriber: %w", err)
	}

	s := subscriber.Data
	if a.githubIssue {
		fmt.Printf("## Subscriber: %s\n\n", s.Email)
		fmt.Printf("- **Status:** %s\n", s.Status)
		fmt.Printf("- **Opens:** %d\n", s.OpensCount)
		fmt.Printf("- **Clicks:** %d\n", s.ClicksCount)
		fmt.Printf("- **Subscribed:** %s\n", s.SubscribedAt)
		if len(s.Groups) > 0 {
			fmt.Println("- **Groups:**")
			for _, g := range s.Groups {
				fmt.Printf("  - %s\n", g.Name)
			}
		}
		return nil
	}

	fmt.Printf("Email:      %s\n", s.Email)
	fmt.Printf("Status:     %s\n", s.Status)
	fmt.Printf("Opens:      %d\n", s.OpensCount)
	fmt.Printf("Clicks:     %d\n", s.ClicksCount)
	fmt.Printf("Open Rate:  %.1f%%\n", s.OpenRate*100)
	fmt.Printf("Click Rate: %.1f%%\n", s.ClickRate*100)
	fmt.Printf("Subscribed: %s\n", s.SubscribedAt)
	fmt.Printf("Created:    %s\n", s.CreatedAt)

	if len(s.Groups) > 0 {
		fmt.Println("\nGroups:")
		for _, g := range s.Groups {
			fmt.Printf("  - %s (ID: %s)\n", g.Name, g.ID)
		}
	}

	if len(s.Fields) > 0 && a.verbose {
		fmt.Println("\nCustom Fields:")
		for k, v := range s.Fields {
			if v != nil && v != "" {
				fmt.Printf("  %s: %v\n", k, v)
			}
		}
	}

	return nil
}

func (a *App) subscribersAdd(email, name string) error {
	subscriber := &mailerlite.UpsertSubscriber{
		Email: email,
	}

	// Set name if provided
	if name != "" {
		subscriber.Fields = map[string]interface{}{
			"name": name,
		}
	}

	result, _, err := a.client.Subscriber.Upsert(a.ctx, subscriber)
	if err != nil {
		return fmt.Errorf("add subscriber: %w", err)
	}

	s := result.Data
	if a.githubIssue {
		fmt.Printf("## Subscriber Added\n\n")
		fmt.Printf("- **Email:** %s\n", s.Email)
		fmt.Printf("- **Status:** %s\n", s.Status)
		return nil
	}

	fmt.Printf("Subscriber added/updated:\n")
	fmt.Printf("  Email:  %s\n", s.Email)
	fmt.Printf("  Status: %s\n", s.Status)
	fmt.Printf("  ID:     %s\n", s.ID)

	return nil
}

func (a *App) subscribersDelete(email string) error {
	// First, get the subscriber to find their ID
	options := &mailerlite.GetSubscriberOptions{
		Email: email,
	}

	subscriber, _, err := a.client.Subscriber.Get(a.ctx, options)
	if err != nil {
		return fmt.Errorf("subscriber not found: %w", err)
	}

	// Delete by ID
	_, err = a.client.Subscriber.Delete(a.ctx, subscriber.Data.ID)
	if err != nil {
		return fmt.Errorf("delete subscriber: %w", err)
	}

	if a.githubIssue {
		fmt.Printf("## Subscriber Deleted\n\n")
		fmt.Printf("- **Email:** %s\n", email)
		return nil
	}

	fmt.Printf("Subscriber deleted: %s\n", email)

	return nil
}

func (a *App) handleGroups(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("groups subcommand required: list, create, subscribers, assign, unassign")
	}

	switch args[0] {
	case "list":
		return a.groupsList()
	case "create":
		if len(args) < 2 {
			return fmt.Errorf("group name required: groups create NAME")
		}
		return a.groupsCreate(strings.Join(args[1:], " "))
	case "subscribers":
		if len(args) < 2 {
			return fmt.Errorf("group ID required: groups subscribers ID")
		}
		return a.groupsSubscribers(args[1])
	case "assign":
		if len(args) < 3 {
			return fmt.Errorf("group ID and email required: groups assign GROUP_ID EMAIL")
		}
		return a.groupsAssign(args[1], args[2])
	case "unassign":
		if len(args) < 3 {
			return fmt.Errorf("group ID and email required: groups unassign GROUP_ID EMAIL")
		}
		return a.groupsUnassign(args[1], args[2])
	default:
		return fmt.Errorf("unknown groups subcommand: %s", args[0])
	}
}

func (a *App) groupsList() error {
	options := &mailerlite.ListGroupOptions{
		Page:  1,
		Limit: 100,
		Sort:  mailerlite.SortByName,
	}

	groups, _, err := a.client.Group.List(a.ctx, options)
	if err != nil {
		return fmt.Errorf("list groups: %w", err)
	}

	if a.githubIssue {
		fmt.Println("## Groups")
		fmt.Println()
		fmt.Printf("Total: **%d**\n\n", groups.Meta.Total)
		if len(groups.Data) > 0 {
			fmt.Println("| Name | Active | Sent |")
			fmt.Println("|------|--------|------|")
			for _, g := range groups.Data {
				fmt.Printf("| %s | %d | %d |\n", g.Name, g.ActiveCount, g.SentCount)
			}
		}
		return nil
	}

	fmt.Printf("Total groups: %d\n\n", groups.Meta.Total)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tACTIVE\tSENT")
	for _, g := range groups.Data {
		fmt.Fprintf(w, "%s\t%s\t%d\t%d\n", g.ID, g.Name, g.ActiveCount, g.SentCount)
	}
	w.Flush()

	return nil
}

func (a *App) groupsSubscribers(groupID string) error {
	options := &mailerlite.ListGroupSubscriberOptions{
		GroupID: groupID,
		Page:    1,
		Limit:   100,
	}

	subscribers, _, err := a.client.Group.Subscribers(a.ctx, options)
	if err != nil {
		return fmt.Errorf("list group subscribers: %w", err)
	}

	if a.githubIssue {
		fmt.Printf("## Group Subscribers (ID: %s)\n\n", groupID)
		fmt.Printf("Total: **%d**\n\n", subscribers.Meta.Total)
		if len(subscribers.Data) > 0 {
			fmt.Println("| Email | Status |")
			fmt.Println("|-------|--------|")
			for _, s := range subscribers.Data {
				fmt.Printf("| %s | %s |\n", s.Email, s.Status)
			}
		}
		return nil
	}

	fmt.Printf("Group %s - %d subscribers\n\n", groupID, subscribers.Meta.Total)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "EMAIL\tSTATUS")
	for _, s := range subscribers.Data {
		fmt.Fprintf(w, "%s\t%s\n", s.Email, s.Status)
	}
	w.Flush()

	return nil
}

func (a *App) groupsCreate(name string) error {
	result, _, err := a.client.Group.Create(a.ctx, name)
	if err != nil {
		return fmt.Errorf("create group: %w", err)
	}

	g := result.Data
	if a.githubIssue {
		fmt.Printf("## Group Created\n\n")
		fmt.Printf("- **Name:** %s\n", g.Name)
		fmt.Printf("- **ID:** %s\n", g.ID)
		return nil
	}

	fmt.Printf("Group created:\n")
	fmt.Printf("  Name: %s\n", g.Name)
	fmt.Printf("  ID:   %s\n", g.ID)

	return nil
}

func (a *App) groupsAssign(groupID, email string) error {
	// First, get the subscriber to find their ID
	options := &mailerlite.GetSubscriberOptions{
		Email: email,
	}

	subscriber, _, err := a.client.Subscriber.Get(a.ctx, options)
	if err != nil {
		return fmt.Errorf("subscriber not found: %w", err)
	}

	_, _, err = a.client.Group.Assign(a.ctx, groupID, subscriber.Data.ID)
	if err != nil {
		return fmt.Errorf("assign to group: %w", err)
	}

	if a.githubIssue {
		fmt.Printf("## Subscriber Assigned to Group\n\n")
		fmt.Printf("- **Email:** %s\n", email)
		fmt.Printf("- **Group ID:** %s\n", groupID)
		return nil
	}

	fmt.Printf("Subscriber %s assigned to group %s\n", email, groupID)

	return nil
}

func (a *App) groupsUnassign(groupID, email string) error {
	// First, get the subscriber to find their ID
	options := &mailerlite.GetSubscriberOptions{
		Email: email,
	}

	subscriber, _, err := a.client.Subscriber.Get(a.ctx, options)
	if err != nil {
		return fmt.Errorf("subscriber not found: %w", err)
	}

	_, err = a.client.Group.UnAssign(a.ctx, groupID, subscriber.Data.ID)
	if err != nil {
		return fmt.Errorf("unassign from group: %w", err)
	}

	if a.githubIssue {
		fmt.Printf("## Subscriber Removed from Group\n\n")
		fmt.Printf("- **Email:** %s\n", email)
		fmt.Printf("- **Group ID:** %s\n", groupID)
		return nil
	}

	fmt.Printf("Subscriber %s removed from group %s\n", email, groupID)

	return nil
}

func (a *App) handleForms(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("forms subcommand required: list, subscribers")
	}

	switch args[0] {
	case "list":
		return a.formsList()
	case "subscribers":
		if len(args) < 2 {
			return fmt.Errorf("form ID required: forms subscribers ID")
		}
		return a.formsSubscribers(args[1])
	default:
		return fmt.Errorf("unknown forms subcommand: %s", args[0])
	}
}

func (a *App) formsList() error {
	// List all form types
	formTypes := []string{"popup", "embedded", "promotion"}

	var allForms []mailerlite.Form
	for _, formType := range formTypes {
		options := &mailerlite.ListFormOptions{
			Type:  formType,
			Page:  1,
			Limit: 100,
		}

		forms, _, err := a.client.Form.List(a.ctx, options)
		if err != nil {
			if a.verbose {
				fmt.Fprintf(os.Stderr, "Warning: failed to list %s forms: %v\n", formType, err)
			}
			continue
		}
		allForms = append(allForms, forms.Data...)
	}

	if a.githubIssue {
		fmt.Println("## Forms")
		fmt.Println()
		fmt.Printf("Total: **%d**\n\n", len(allForms))
		if len(allForms) > 0 {
			fmt.Println("| Name | Type | Opens | Conversions |")
			fmt.Println("|------|------|-------|-------------|")
			for _, f := range allForms {
				fmt.Printf("| %s | %s | %d | %d |\n", f.Name, f.Type, f.OpensCount, f.ConversionsCount)
			}
		}
		return nil
	}

	fmt.Printf("Total forms: %d\n\n", len(allForms))

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tTYPE\tOPENS\tCONVERSIONS")
	for _, f := range allForms {
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%d\n", f.Id, f.Name, f.Type, f.OpensCount, f.ConversionsCount)
	}
	w.Flush()

	return nil
}

func (a *App) formsSubscribers(formID string) error {
	options := &mailerlite.ListFormSubscriberOptions{
		FormID: formID,
		Page:   1,
		Limit:  100,
	}

	subscribers, _, err := a.client.Form.Subscribers(a.ctx, options)
	if err != nil {
		return fmt.Errorf("list form subscribers: %w", err)
	}

	if a.githubIssue {
		fmt.Printf("## Form Subscribers (ID: %s)\n\n", formID)
		fmt.Printf("Total: **%d**\n\n", subscribers.Meta.Total)
		if len(subscribers.Data) > 0 {
			fmt.Println("| Email | Status |")
			fmt.Println("|-------|--------|")
			for _, s := range subscribers.Data {
				fmt.Printf("| %s | %s |\n", s.Email, s.Status)
			}
		}
		return nil
	}

	fmt.Printf("Form %s - %d subscribers\n\n", formID, subscribers.Meta.Total)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "EMAIL\tSTATUS")
	for _, s := range subscribers.Data {
		fmt.Fprintf(w, "%s\t%s\n", s.Email, s.Status)
	}
	w.Flush()

	return nil
}

func (a *App) handleWebhooks(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("webhooks subcommand required: list, create, delete")
	}

	switch args[0] {
	case "list":
		return a.webhooksList()
	case "create":
		if len(args) < 2 {
			return fmt.Errorf("URL required: webhooks create URL [EVENT...]")
		}
		events := []string{EventSubscriberCreated}
		if len(args) > 2 {
			events = args[2:]
		}
		return a.webhooksCreate(args[1], events)
	case "delete":
		if len(args) < 2 {
			return fmt.Errorf("webhook ID required: webhooks delete ID")
		}
		return a.webhooksDelete(args[1])
	default:
		return fmt.Errorf("unknown webhooks subcommand: %s", args[0])
	}
}

func (a *App) webhooksList() error {
	options := &mailerlite.ListWebhookOptions{
		Page:  1,
		Limit: 100,
	}

	webhooks, _, err := a.client.Webhook.List(a.ctx, options)
	if err != nil {
		return fmt.Errorf("list webhooks: %w", err)
	}

	if a.githubIssue {
		fmt.Println("## Webhooks")
		fmt.Println()
		fmt.Printf("Total: **%d**\n\n", len(webhooks.Data))
		if len(webhooks.Data) > 0 {
			fmt.Println("| Name | URL | Events | Enabled |")
			fmt.Println("|------|-----|--------|---------|")
			for _, w := range webhooks.Data {
				events := strings.Join(w.Events, ", ")
				fmt.Printf("| %s | %s | %s | %v |\n", w.Name, w.Url, events, w.Enabled)
			}
		}
		return nil
	}

	fmt.Printf("Total webhooks: %d\n\n", len(webhooks.Data))

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tURL\tEVENTS\tENABLED")
	for _, wh := range webhooks.Data {
		events := strings.Join(wh.Events, ", ")
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%v\n", wh.Id, wh.Name, wh.Url, events, wh.Enabled)
	}
	w.Flush()

	return nil
}

func (a *App) webhooksCreate(url string, events []string) error {
	webhook := &mailerlite.CreateWebhookOptions{
		Name:   "CLI Webhook",
		Events: events,
		Url:    url,
	}

	result, _, err := a.client.Webhook.Create(a.ctx, webhook)
	if err != nil {
		return fmt.Errorf("create webhook: %w", err)
	}

	w := result.Data
	if a.githubIssue {
		fmt.Printf("## Webhook Created\n\n")
		fmt.Printf("- **ID:** %s\n", w.Id)
		fmt.Printf("- **URL:** %s\n", w.Url)
		fmt.Printf("- **Events:** %s\n", strings.Join(w.Events, ", "))
		return nil
	}

	fmt.Printf("Webhook created:\n")
	fmt.Printf("  ID:     %s\n", w.Id)
	fmt.Printf("  URL:    %s\n", w.Url)
	fmt.Printf("  Events: %s\n", strings.Join(w.Events, ", "))

	return nil
}

func (a *App) webhooksDelete(webhookID string) error {
	_, err := a.client.Webhook.Delete(a.ctx, webhookID)
	if err != nil {
		return fmt.Errorf("delete webhook: %w", err)
	}

	if a.githubIssue {
		fmt.Printf("## Webhook Deleted\n\n")
		fmt.Printf("- **ID:** %s\n", webhookID)
		return nil
	}

	fmt.Printf("Webhook deleted: %s\n", webhookID)

	return nil
}

func (a *App) handleAutomations(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("automations subcommand required: list, get")
	}

	switch args[0] {
	case "list":
		return a.automationsList()
	case "get":
		if len(args) < 2 {
			return fmt.Errorf("automation ID required: automations get ID")
		}
		return a.automationsGet(args[1])
	default:
		return fmt.Errorf("unknown automations subcommand: %s", args[0])
	}
}

func (a *App) automationsList() error {
	options := &mailerlite.ListAutomationOptions{
		Page:  1,
		Limit: 100,
	}

	automations, _, err := a.client.Automation.List(a.ctx, options)
	if err != nil {
		return fmt.Errorf("list automations: %w", err)
	}

	if a.githubIssue {
		fmt.Println("## Automations")
		fmt.Println()
		fmt.Printf("Total: **%d**\n\n", len(automations.Data))
		if len(automations.Data) > 0 {
			fmt.Println("| Name | Enabled | Emails | Opens | Clicks |")
			fmt.Println("|------|---------|--------|-------|--------|")
			for _, auto := range automations.Data {
				fmt.Printf("| %s | %v | %d | %d | %d |\n",
					auto.Name, auto.Enabled, auto.Stats.Sent, auto.Stats.OpensCount, auto.Stats.ClicksCount)
			}
		}
		return nil
	}

	fmt.Printf("Total automations: %d\n\n", len(automations.Data))

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tENABLED\tEMAILS\tOPENS\tCLICKS")
	for _, auto := range automations.Data {
		fmt.Fprintf(w, "%s\t%s\t%v\t%d\t%d\t%d\n",
			auto.ID, auto.Name, auto.Enabled, auto.Stats.Sent, auto.Stats.OpensCount, auto.Stats.ClicksCount)
	}
	w.Flush()

	return nil
}

func (a *App) automationsGet(automationID string) error {
	automation, _, err := a.client.Automation.Get(a.ctx, automationID)
	if err != nil {
		return fmt.Errorf("get automation: %w", err)
	}

	auto := automation.Data
	if a.githubIssue {
		fmt.Printf("## Automation: %s\n\n", auto.Name)
		fmt.Printf("- **ID:** %s\n", auto.ID)
		fmt.Printf("- **Enabled:** %v\n", auto.Enabled)
		fmt.Printf("- **Emails Sent:** %d\n", auto.Stats.Sent)
		fmt.Printf("- **Opens:** %d\n", auto.Stats.OpensCount)
		fmt.Printf("- **Clicks:** %d\n", auto.Stats.ClicksCount)
		return nil
	}

	fmt.Printf("Automation: %s\n", auto.Name)
	fmt.Println(strings.Repeat("=", 40))
	fmt.Printf("ID:          %s\n", auto.ID)
	fmt.Printf("Enabled:     %v\n", auto.Enabled)
	fmt.Printf("Emails Sent: %d\n", auto.Stats.Sent)
	fmt.Printf("Opens:       %d\n", auto.Stats.OpensCount)
	fmt.Printf("Clicks:      %d\n", auto.Stats.ClicksCount)

	return nil
}

func (a *App) handleStats() error {
	// Get subscriber count
	count, _, err := a.client.Subscriber.Count(a.ctx)
	if err != nil {
		return fmt.Errorf("get subscriber count: %w", err)
	}

	// Get groups
	groupOptions := &mailerlite.ListGroupOptions{
		Page:  1,
		Limit: 100,
	}
	groups, _, err := a.client.Group.List(a.ctx, groupOptions)
	if err != nil {
		return fmt.Errorf("list groups: %w", err)
	}

	// Count active subscribers across groups
	totalActive := 0
	for _, g := range groups.Data {
		totalActive += g.ActiveCount
	}

	if a.githubIssue {
		fmt.Println("## MailerLite Statistics")
		fmt.Println()
		fmt.Printf("| Metric | Value |\n")
		fmt.Printf("|--------|-------|\n")
		fmt.Printf("| Total Subscribers | %d |\n", count.Total)
		fmt.Printf("| Total Groups | %d |\n", groups.Meta.Total)
		fmt.Println()
		fmt.Printf("*Generated: %s*\n", time.Now().Format("2006-01-02 15:04:05"))
		return nil
	}

	fmt.Println("MailerLite Statistics")
	fmt.Println(strings.Repeat("=", 40))
	fmt.Printf("Total Subscribers: %d\n", count.Total)
	fmt.Printf("Total Groups:      %d\n", groups.Meta.Total)
	fmt.Println()
	fmt.Printf("Generated: %s\n", time.Now().Format("2006-01-02 15:04:05"))

	return nil
}

func (a *App) handleOpen(args []string) error {
	target := "dashboard"
	if len(args) > 0 {
		target = args[0]
	}

	var url string
	switch target {
	case "dashboard", "home":
		url = DashboardURL
	case "subscribers":
		url = SubscribersURL
	case "groups":
		url = GroupsURL
	case "forms":
		url = FormsURL
	case "automations":
		url = AutomationsURL
	case "campaigns":
		url = CampaignsURL
	case "webhooks":
		url = WebhooksURL
	case "integrations":
		url = IntegrationsURL
	case "api", "apikey":
		url = APIKeyURL
	default:
		return fmt.Errorf("unknown target: %s (use: dashboard, subscribers, groups, forms, automations, campaigns, webhooks, integrations, api)", target)
	}

	fmt.Printf("Opening %s...\n", url)
	return openBrowser(url)
}

// openBrowser opens a URL in the default browser
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	default: // linux, freebsd, etc.
		cmd = "xdg-open"
		args = []string{url}
	}

	return exec.Command(cmd, args...).Start()
}

// handleReleases handles the releases command for GitHub release information
func (a *App) handleReleases(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("releases subcommand required: latest, list, urls")
	}

	// Parse optional owner/repo flags
	owner := DefaultGitHubOwner
	repo := DefaultGitHubRepo

	// Check for OWNER=x REPO=y in args
	var subCmd string
	for _, arg := range args {
		if strings.HasPrefix(arg, "OWNER=") {
			owner = strings.TrimPrefix(arg, "OWNER=")
		} else if strings.HasPrefix(arg, "REPO=") {
			repo = strings.TrimPrefix(arg, "REPO=")
		} else if subCmd == "" {
			subCmd = arg
		}
	}

	switch subCmd {
	case "latest":
		return a.releasesLatest(owner, repo)
	case "list":
		return a.releasesList(owner, repo)
	case "urls":
		return a.releasesURLs(owner, repo)
	default:
		return fmt.Errorf("unknown releases subcommand: %s (use: latest, list, urls)", subCmd)
	}
}

// releasesLatest fetches and displays the latest release
func (a *App) releasesLatest(owner, repo string) error {
	apiURL := fmt.Sprintf("%s/repos/%s/%s/releases/latest", GitHubAPIBase, owner, repo)

	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("failed to fetch release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return fmt.Errorf("no releases found for %s/%s", owner, repo)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return fmt.Errorf("failed to parse release: %w", err)
	}

	if a.githubIssue {
		fmt.Printf("## Latest Release: %s\n\n", release.Name)
		fmt.Printf("- **Tag:** %s\n", release.TagName)
		fmt.Printf("- **URL:** %s\n", release.HTMLURL)
		if len(release.Assets) > 0 {
			fmt.Println("\n### Downloads")
			for _, asset := range release.Assets {
				fmt.Printf("- [%s](%s) (%.2f MB)\n", asset.Name, asset.DownloadURL, float64(asset.Size)/1024/1024)
			}
		}
		return nil
	}

	fmt.Printf("Latest Release: %s\n", release.Name)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Tag:       %s\n", release.TagName)
	fmt.Printf("Published: %s\n", release.PublishedAt)
	fmt.Printf("URL:       %s\n", release.HTMLURL)

	if len(release.Assets) > 0 {
		fmt.Println("\nDownloads:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "FILENAME\tSIZE\tURL")
		for _, asset := range release.Assets {
			fmt.Fprintf(w, "%s\t%.2f MB\t%s\n", asset.Name, float64(asset.Size)/1024/1024, asset.DownloadURL)
		}
		w.Flush()
	}

	return nil
}

// releasesList fetches and displays all releases
func (a *App) releasesList(owner, repo string) error {
	apiURL := fmt.Sprintf("%s/repos/%s/%s/releases", GitHubAPIBase, owner, repo)

	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("failed to fetch releases: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var releases []Release
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return fmt.Errorf("failed to parse releases: %w", err)
	}

	if len(releases) == 0 {
		fmt.Printf("No releases found for %s/%s\n", owner, repo)
		return nil
	}

	if a.githubIssue {
		fmt.Printf("## Releases for %s/%s\n\n", owner, repo)
		fmt.Println("| Tag | Name | Published | Assets |")
		fmt.Println("|-----|------|-----------|--------|")
		for _, r := range releases {
			fmt.Printf("| %s | %s | %s | %d |\n", r.TagName, r.Name, r.PublishedAt[:10], len(r.Assets))
		}
		return nil
	}

	fmt.Printf("Releases for %s/%s\n", owner, repo)
	fmt.Println(strings.Repeat("=", 50))

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TAG\tNAME\tPUBLISHED\tASSETS")
	for _, r := range releases {
		published := r.PublishedAt
		if len(published) > 10 {
			published = published[:10]
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\n", r.TagName, r.Name, published, len(r.Assets))
	}
	w.Flush()

	return nil
}

// releasesURLs prints the release URLs for email templates
func (a *App) releasesURLs(owner, repo string) error {
	fmt.Printf("GitHub Release URLs for %s/%s\n", owner, repo)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println()
	fmt.Println("Use these URLs in your MailerLite email templates:")
	fmt.Println()
	fmt.Printf("  Releases Page:   %s\n", GetReleasesURL(owner, repo))
	fmt.Printf("  Latest Release:  %s\n", GetLatestReleaseURL(owner, repo))
	fmt.Println()
	fmt.Println("For direct download links, use:")
	fmt.Printf("  %s\n", fmt.Sprintf(GitHubDownloadURL, owner, repo, "{TAG}", "{FILENAME}"))
	fmt.Println()
	fmt.Println("Example:")
	fmt.Printf("  %s\n", GetDownloadURL(owner, repo, "v1.0.0", "software-darwin-arm64.tar.gz"))
	fmt.Printf("  %s\n", GetDownloadURL(owner, repo, "v1.0.0", "software-windows-amd64.zip"))
	fmt.Printf("  %s\n", GetDownloadURL(owner, repo, "v1.0.0", "software-linux-amd64.tar.gz"))

	return nil
}

// ============================================================================
// Web3Forms Webhook Server
// ============================================================================

// Web3FormPayload represents the incoming webhook payload from Web3Forms
type Web3FormPayload struct {
	// Standard Web3Forms fields
	Name      string `json:"name"`
	Email     string `json:"email"`
	Message   string `json:"message"`
	Subject   string `json:"subject"`
	AccessKey string `json:"access_key"`

	// Custom fields from our Get Started form
	Company  string `json:"company"`
	Platform string `json:"platform"`
	Industry string `json:"industry"`
	UseCase  string `json:"usecase"`
}

// handleServer starts a webhook server to receive Web3Forms submissions
func (a *App) handleServer(args []string) error {
	port := 8086  // Default port
	groupID := "" // Optional: auto-assign to group

	// Parse args for PORT= and GROUP_ID=
	for _, arg := range args {
		if strings.HasPrefix(arg, "PORT=") {
			fmt.Sscanf(strings.TrimPrefix(arg, "PORT="), "%d", &port)
		} else if strings.HasPrefix(arg, "GROUP_ID=") {
			groupID = strings.TrimPrefix(arg, "GROUP_ID=")
		}
	}

	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║         MailerLite Webhook Server                            ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("Starting server on port %d...\n", port)
	fmt.Println()
	fmt.Println("Webhook URL (for Web3Forms):")
	fmt.Printf("  http://localhost:%d/webhook\n", port)
	fmt.Println()
	fmt.Println("For production, use a tunnel service like:")
	fmt.Println("  ngrok http " + fmt.Sprintf("%d", port))
	fmt.Println("  cloudflared tunnel --url http://localhost:" + fmt.Sprintf("%d", port))
	fmt.Println()
	if groupID != "" {
		fmt.Printf("Auto-assigning to group: %s\n", groupID)
	}
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println(strings.Repeat("─", 60))

	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	// Webhook endpoint
	mux.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse the form data (Web3Forms sends form-urlencoded or JSON)
		var payload Web3FormPayload

		contentType := r.Header.Get("Content-Type")
		if strings.Contains(contentType, "application/json") {
			// JSON payload
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				fmt.Printf("[ERROR] Failed to parse JSON: %v\n", err)
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}
		} else {
			// Form-encoded payload
			if err := r.ParseForm(); err != nil {
				fmt.Printf("[ERROR] Failed to parse form: %v\n", err)
				http.Error(w, "Invalid form data", http.StatusBadRequest)
				return
			}
			payload = Web3FormPayload{
				Name:     r.FormValue("name"),
				Email:    r.FormValue("email"),
				Message:  r.FormValue("message"),
				Subject:  r.FormValue("subject"),
				Company:  r.FormValue("company"),
				Platform: r.FormValue("platform"),
				Industry: r.FormValue("industry"),
				UseCase:  r.FormValue("usecase"),
			}
		}

		// Validate required fields
		if payload.Email == "" {
			fmt.Println("[WARN] Received webhook without email")
			http.Error(w, "Email required", http.StatusBadRequest)
			return
		}

		// Log the submission
		fmt.Printf("\n[%s] New submission\n", time.Now().Format("15:04:05"))
		fmt.Printf("  Email:    %s\n", payload.Email)
		if payload.Name != "" {
			fmt.Printf("  Name:     %s\n", payload.Name)
		}
		if payload.Company != "" {
			fmt.Printf("  Company:  %s\n", payload.Company)
		}
		if payload.Platform != "" {
			fmt.Printf("  Platform: %s\n", payload.Platform)
		}
		if payload.Industry != "" {
			fmt.Printf("  Industry: %s\n", payload.Industry)
		}

		// Add subscriber to MailerLite
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		subscriber := &mailerlite.UpsertSubscriber{
			Email: payload.Email,
			Fields: map[string]interface{}{
				"name":    payload.Name,
				"company": payload.Company,
			},
		}

		result, _, err := a.client.Subscriber.Upsert(ctx, subscriber)
		if err != nil {
			fmt.Printf("  [ERROR] Failed to add to MailerLite: %v\n", err)
			// Still return success to Web3Forms (don't retry)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "Received (MailerLite error logged)")
			return
		}

		fmt.Printf("  [OK] Added to MailerLite: ID=%s, Status=%s\n", result.Data.ID, result.Data.Status)

		// Auto-assign to group if specified
		if groupID != "" {
			_, _, err := a.client.Group.Assign(ctx, groupID, result.Data.ID)
			if err != nil {
				fmt.Printf("  [WARN] Failed to assign to group: %v\n", err)
			} else {
				fmt.Printf("  [OK] Assigned to group: %s\n", groupID)
			}
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	return server.ListenAndServe()
}

func printUsage() {
	fmt.Println("mailerlite - CLI tool for MailerLite subscriber management")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  mailerlite [flags] <command> [subcommand] [args]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  subscribers list              List all subscribers")
	fmt.Println("  subscribers count             Count total subscribers")
	fmt.Println("  subscribers get EMAIL         Get subscriber by email")
	fmt.Println("  subscribers add EMAIL [NAME]  Add or update subscriber")
	fmt.Println("  subscribers delete EMAIL      Delete subscriber")
	fmt.Println("  groups list                   List all groups")
	fmt.Println("  groups create NAME            Create a new group")
	fmt.Println("  groups subscribers ID         List subscribers in a group")
	fmt.Println("  groups assign ID EMAIL        Assign subscriber to group")
	fmt.Println("  groups unassign ID EMAIL      Remove subscriber from group")
	fmt.Println("  forms list                    List all forms")
	fmt.Println("  forms subscribers ID          List subscribers for a form")
	fmt.Println("  webhooks list                 List all webhooks")
	fmt.Println("  webhooks create URL [EVENTS]  Create a webhook (default: subscriber.created)")
	fmt.Println("  webhooks delete ID            Delete a webhook")
	fmt.Println("  automations list              List all automations")
	fmt.Println("  automations get ID            Get automation details")
	fmt.Println("  stats                         Show account statistics")
	fmt.Println("  open [TARGET]                 Open dashboard in browser")
	fmt.Println("  releases latest               Show latest GitHub release")
	fmt.Println("  releases list                 List all GitHub releases")
	fmt.Println("  releases urls                 Show release URLs for email templates")
	fmt.Println("  server [PORT=8086]            Start webhook server for Web3Forms")
	fmt.Println()
	fmt.Println("Open Targets:")
	fmt.Println("  dashboard (default), subscribers, groups, forms, automations")
	fmt.Println("  campaigns, webhooks, integrations, api")
	fmt.Println()
	fmt.Println("Webhook Events:")
	fmt.Printf("  %s, %s, %s\n", EventSubscriberCreated, EventSubscriberUpdated, EventSubscriberUnsubscribed)
	fmt.Printf("  %s, %s\n", EventSubscriberAddedToGroup, EventSubscriberRemovedFromGroup)
	fmt.Printf("  %s, %s, %s\n", EventSubscriberBounced, EventSubscriberAutomationTriggered, EventSubscriberAutomationComplete)
	fmt.Printf("  %s, %s, %s\n", EventCampaignSent, EventCampaignOpened, EventCampaignClicked)
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -github-issue    Output markdown for GitHub issue")
	fmt.Println("  -v               Verbose output")
	fmt.Println("  -version         Show version")
	fmt.Println()
	fmt.Println("Environment:")
	fmt.Println("  MAILERLITE_API_KEY    API key (required)")
	fmt.Println()
	fmt.Printf("Get API key from: %s\n", APIKeyURL)
	fmt.Println()
	fmt.Println("Dashboard URLs:")
	fmt.Printf("  Dashboard:     %s\n", DashboardURL)
	fmt.Printf("  Subscribers:   %s\n", SubscribersURL)
	fmt.Printf("  Groups:        %s\n", GroupsURL)
	fmt.Printf("  Forms:         %s\n", FormsURL)
	fmt.Printf("  Automations:   %s\n", AutomationsURL)
	fmt.Printf("  Webhooks:      %s\n", WebhooksURL)
}
