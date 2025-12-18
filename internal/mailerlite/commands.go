package mailerlite

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/mailerlite/mailerlite-go"
)

// ============================================================================
// Subscribers Commands
// ============================================================================

func (c *CLI) handleSubscribers(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("subscribers subcommand required: list, count, get, add, delete")
	}

	switch args[0] {
	case "list":
		return c.subscribersList()
	case "count":
		return c.subscribersCount()
	case "get":
		if len(args) < 2 {
			return fmt.Errorf("email required: subscribers get EMAIL")
		}
		return c.subscribersGet(args[1])
	case "add":
		if len(args) < 2 {
			return fmt.Errorf("email required: subscribers add EMAIL [NAME]")
		}
		name := ""
		if len(args) >= 3 {
			name = strings.Join(args[2:], " ")
		}
		return c.subscribersAdd(args[1], name)
	case "delete":
		if len(args) < 2 {
			return fmt.Errorf("email required: subscribers delete EMAIL")
		}
		return c.subscribersDelete(args[1])
	default:
		return fmt.Errorf("unknown subscribers subcommand: %s", args[0])
	}
}

func (c *CLI) subscribersList() error {
	options := &mailerlite.ListSubscriberOptions{
		Limit: 100,
		Page:  1,
	}

	subscribers, _, err := c.client.sdk.Subscriber.List(c.ctx, options)
	if err != nil {
		return fmt.Errorf("list subscribers: %w", err)
	}

	if c.githubIssue {
		c.println("## Subscribers")
		c.println()
		c.printf("Total: **%d**\n\n", subscribers.Meta.Total)
		if len(subscribers.Data) > 0 {
			c.println("| Email | Status | Subscribed |")
			c.println("|-------|--------|------------|")
			for _, s := range subscribers.Data {
				c.printf("| %s | %s | %s |\n", s.Email, s.Status, s.SubscribedAt)
			}
		}
		return nil
	}

	c.printf("Total subscribers: %d\n\n", subscribers.Meta.Total)

	w := tabwriter.NewWriter(c.out, 0, 0, 2, ' ', 0)
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

func (c *CLI) subscribersCount() error {
	count, _, err := c.client.sdk.Subscriber.Count(c.ctx)
	if err != nil {
		return fmt.Errorf("count subscribers: %w", err)
	}

	if c.githubIssue {
		c.printf("**Total Subscribers:** %d\n", count.Total)
		return nil
	}

	c.printf("Total subscribers: %d\n", count.Total)
	return nil
}

func (c *CLI) subscribersGet(email string) error {
	options := &mailerlite.GetSubscriberOptions{
		Email: email,
	}

	subscriber, _, err := c.client.sdk.Subscriber.Get(c.ctx, options)
	if err != nil {
		return fmt.Errorf("get subscriber: %w", err)
	}

	s := subscriber.Data
	if c.githubIssue {
		c.printf("## Subscriber: %s\n\n", s.Email)
		c.printf("- **Status:** %s\n", s.Status)
		c.printf("- **Opens:** %d\n", s.OpensCount)
		c.printf("- **Clicks:** %d\n", s.ClicksCount)
		c.printf("- **Subscribed:** %s\n", s.SubscribedAt)
		if len(s.Groups) > 0 {
			c.println("- **Groups:**")
			for _, g := range s.Groups {
				c.printf("  - %s\n", g.Name)
			}
		}
		return nil
	}

	c.printf("Email:      %s\n", s.Email)
	c.printf("Status:     %s\n", s.Status)
	c.printf("Opens:      %d\n", s.OpensCount)
	c.printf("Clicks:     %d\n", s.ClicksCount)
	c.printf("Open Rate:  %.1f%%\n", s.OpenRate*100)
	c.printf("Click Rate: %.1f%%\n", s.ClickRate*100)
	c.printf("Subscribed: %s\n", s.SubscribedAt)
	c.printf("Created:    %s\n", s.CreatedAt)

	if len(s.Groups) > 0 {
		c.println("\nGroups:")
		for _, g := range s.Groups {
			c.printf("  - %s (ID: %s)\n", g.Name, g.ID)
		}
	}

	if len(s.Fields) > 0 && c.verbose {
		c.println("\nCustom Fields:")
		for k, v := range s.Fields {
			if v != nil && v != "" {
				c.printf("  %s: %v\n", k, v)
			}
		}
	}

	return nil
}

func (c *CLI) subscribersAdd(email, name string) error {
	subscriber := &mailerlite.UpsertSubscriber{
		Email: email,
	}

	if name != "" {
		subscriber.Fields = map[string]interface{}{
			"name": name,
		}
	}

	result, _, err := c.client.sdk.Subscriber.Upsert(c.ctx, subscriber)
	if err != nil {
		return fmt.Errorf("add subscriber: %w", err)
	}

	s := result.Data
	if c.githubIssue {
		c.println("## Subscriber Added")
		c.println()
		c.printf("- **Email:** %s\n", s.Email)
		c.printf("- **Status:** %s\n", s.Status)
		return nil
	}

	c.println("Subscriber added/updated:")
	c.printf("  Email:  %s\n", s.Email)
	c.printf("  Status: %s\n", s.Status)
	c.printf("  ID:     %s\n", s.ID)

	return nil
}

func (c *CLI) subscribersDelete(email string) error {
	options := &mailerlite.GetSubscriberOptions{
		Email: email,
	}

	subscriber, _, err := c.client.sdk.Subscriber.Get(c.ctx, options)
	if err != nil {
		return fmt.Errorf("subscriber not found: %w", err)
	}

	_, err = c.client.sdk.Subscriber.Delete(c.ctx, subscriber.Data.ID)
	if err != nil {
		return fmt.Errorf("delete subscriber: %w", err)
	}

	if c.githubIssue {
		c.println("## Subscriber Deleted")
		c.println()
		c.printf("- **Email:** %s\n", email)
		return nil
	}

	c.printf("Subscriber deleted: %s\n", email)
	return nil
}

// ============================================================================
// Groups Commands
// ============================================================================

func (c *CLI) handleGroups(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("groups subcommand required: list, create, subscribers, assign, unassign")
	}

	switch args[0] {
	case "list":
		return c.groupsList()
	case "create":
		if len(args) < 2 {
			return fmt.Errorf("group name required: groups create NAME")
		}
		return c.groupsCreate(strings.Join(args[1:], " "))
	case "subscribers":
		if len(args) < 2 {
			return fmt.Errorf("group ID required: groups subscribers ID")
		}
		return c.groupsSubscribers(args[1])
	case "assign":
		if len(args) < 3 {
			return fmt.Errorf("group ID and email required: groups assign GROUP_ID EMAIL")
		}
		return c.groupsAssign(args[1], args[2])
	case "unassign":
		if len(args) < 3 {
			return fmt.Errorf("group ID and email required: groups unassign GROUP_ID EMAIL")
		}
		return c.groupsUnassign(args[1], args[2])
	default:
		return fmt.Errorf("unknown groups subcommand: %s", args[0])
	}
}

func (c *CLI) groupsList() error {
	options := &mailerlite.ListGroupOptions{
		Page:  1,
		Limit: 100,
		Sort:  mailerlite.SortByName,
	}

	groups, _, err := c.client.sdk.Group.List(c.ctx, options)
	if err != nil {
		return fmt.Errorf("list groups: %w", err)
	}

	if c.githubIssue {
		c.println("## Groups")
		c.println()
		c.printf("Total: **%d**\n\n", groups.Meta.Total)
		if len(groups.Data) > 0 {
			c.println("| Name | Active | Sent |")
			c.println("|------|--------|------|")
			for _, g := range groups.Data {
				c.printf("| %s | %d | %d |\n", g.Name, g.ActiveCount, g.SentCount)
			}
		}
		return nil
	}

	c.printf("Total groups: %d\n\n", groups.Meta.Total)

	w := tabwriter.NewWriter(c.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tACTIVE\tSENT")
	for _, g := range groups.Data {
		fmt.Fprintf(w, "%s\t%s\t%d\t%d\n", g.ID, g.Name, g.ActiveCount, g.SentCount)
	}
	w.Flush()

	return nil
}

func (c *CLI) groupsSubscribers(groupID string) error {
	options := &mailerlite.ListGroupSubscriberOptions{
		GroupID: groupID,
		Page:    1,
		Limit:   100,
	}

	subscribers, _, err := c.client.sdk.Group.Subscribers(c.ctx, options)
	if err != nil {
		return fmt.Errorf("list group subscribers: %w", err)
	}

	if c.githubIssue {
		c.printf("## Group Subscribers (ID: %s)\n\n", groupID)
		c.printf("Total: **%d**\n\n", subscribers.Meta.Total)
		if len(subscribers.Data) > 0 {
			c.println("| Email | Status |")
			c.println("|-------|--------|")
			for _, s := range subscribers.Data {
				c.printf("| %s | %s |\n", s.Email, s.Status)
			}
		}
		return nil
	}

	c.printf("Group %s - %d subscribers\n\n", groupID, subscribers.Meta.Total)

	w := tabwriter.NewWriter(c.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "EMAIL\tSTATUS")
	for _, s := range subscribers.Data {
		fmt.Fprintf(w, "%s\t%s\n", s.Email, s.Status)
	}
	w.Flush()

	return nil
}

func (c *CLI) groupsCreate(name string) error {
	result, _, err := c.client.sdk.Group.Create(c.ctx, name)
	if err != nil {
		return fmt.Errorf("create group: %w", err)
	}

	g := result.Data
	if c.githubIssue {
		c.println("## Group Created")
		c.println()
		c.printf("- **Name:** %s\n", g.Name)
		c.printf("- **ID:** %s\n", g.ID)
		return nil
	}

	c.println("Group created:")
	c.printf("  Name: %s\n", g.Name)
	c.printf("  ID:   %s\n", g.ID)

	return nil
}

func (c *CLI) groupsAssign(groupID, email string) error {
	options := &mailerlite.GetSubscriberOptions{
		Email: email,
	}

	subscriber, _, err := c.client.sdk.Subscriber.Get(c.ctx, options)
	if err != nil {
		return fmt.Errorf("subscriber not found: %w", err)
	}

	_, _, err = c.client.sdk.Group.Assign(c.ctx, groupID, subscriber.Data.ID)
	if err != nil {
		return fmt.Errorf("assign to group: %w", err)
	}

	if c.githubIssue {
		c.println("## Subscriber Assigned to Group")
		c.println()
		c.printf("- **Email:** %s\n", email)
		c.printf("- **Group ID:** %s\n", groupID)
		return nil
	}

	c.printf("Subscriber %s assigned to group %s\n", email, groupID)
	return nil
}

func (c *CLI) groupsUnassign(groupID, email string) error {
	options := &mailerlite.GetSubscriberOptions{
		Email: email,
	}

	subscriber, _, err := c.client.sdk.Subscriber.Get(c.ctx, options)
	if err != nil {
		return fmt.Errorf("subscriber not found: %w", err)
	}

	_, err = c.client.sdk.Group.UnAssign(c.ctx, groupID, subscriber.Data.ID)
	if err != nil {
		return fmt.Errorf("unassign from group: %w", err)
	}

	if c.githubIssue {
		c.println("## Subscriber Removed from Group")
		c.println()
		c.printf("- **Email:** %s\n", email)
		c.printf("- **Group ID:** %s\n", groupID)
		return nil
	}

	c.printf("Subscriber %s removed from group %s\n", email, groupID)
	return nil
}

// ============================================================================
// Forms Commands
// ============================================================================

func (c *CLI) handleForms(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("forms subcommand required: list, subscribers")
	}

	switch args[0] {
	case "list":
		return c.formsList()
	case "subscribers":
		if len(args) < 2 {
			return fmt.Errorf("form ID required: forms subscribers ID")
		}
		return c.formsSubscribers(args[1])
	default:
		return fmt.Errorf("unknown forms subcommand: %s", args[0])
	}
}

func (c *CLI) formsList() error {
	formTypes := []string{"popup", "embedded", "promotion"}

	var allForms []mailerlite.Form
	for _, formType := range formTypes {
		options := &mailerlite.ListFormOptions{
			Type:  formType,
			Page:  1,
			Limit: 100,
		}

		forms, _, err := c.client.sdk.Form.List(c.ctx, options)
		if err != nil {
			if c.verbose {
				c.printf("Warning: failed to list %s forms: %v\n", formType, err)
			}
			continue
		}
		allForms = append(allForms, forms.Data...)
	}

	if c.githubIssue {
		c.println("## Forms")
		c.println()
		c.printf("Total: **%d**\n\n", len(allForms))
		if len(allForms) > 0 {
			c.println("| Name | Type | Opens | Conversions |")
			c.println("|------|------|-------|-------------|")
			for _, f := range allForms {
				c.printf("| %s | %s | %d | %d |\n", f.Name, f.Type, f.OpensCount, f.ConversionsCount)
			}
		}
		return nil
	}

	c.printf("Total forms: %d\n\n", len(allForms))

	w := tabwriter.NewWriter(c.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tTYPE\tOPENS\tCONVERSIONS")
	for _, f := range allForms {
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%d\n", f.Id, f.Name, f.Type, f.OpensCount, f.ConversionsCount)
	}
	w.Flush()

	return nil
}

func (c *CLI) formsSubscribers(formID string) error {
	options := &mailerlite.ListFormSubscriberOptions{
		FormID: formID,
		Page:   1,
		Limit:  100,
	}

	subscribers, _, err := c.client.sdk.Form.Subscribers(c.ctx, options)
	if err != nil {
		return fmt.Errorf("list form subscribers: %w", err)
	}

	if c.githubIssue {
		c.printf("## Form Subscribers (ID: %s)\n\n", formID)
		c.printf("Total: **%d**\n\n", subscribers.Meta.Total)
		if len(subscribers.Data) > 0 {
			c.println("| Email | Status |")
			c.println("|-------|--------|")
			for _, s := range subscribers.Data {
				c.printf("| %s | %s |\n", s.Email, s.Status)
			}
		}
		return nil
	}

	c.printf("Form %s - %d subscribers\n\n", formID, subscribers.Meta.Total)

	w := tabwriter.NewWriter(c.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "EMAIL\tSTATUS")
	for _, s := range subscribers.Data {
		fmt.Fprintf(w, "%s\t%s\n", s.Email, s.Status)
	}
	w.Flush()

	return nil
}

// ============================================================================
// Webhooks Commands
// ============================================================================

func (c *CLI) handleWebhooks(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("webhooks subcommand required: list, create, delete")
	}

	switch args[0] {
	case "list":
		return c.webhooksList()
	case "create":
		if len(args) < 2 {
			return fmt.Errorf("URL required: webhooks create URL [EVENT...]")
		}
		events := []string{EventSubscriberCreated}
		if len(args) > 2 {
			events = args[2:]
		}
		return c.webhooksCreate(args[1], events)
	case "delete":
		if len(args) < 2 {
			return fmt.Errorf("webhook ID required: webhooks delete ID")
		}
		return c.webhooksDelete(args[1])
	default:
		return fmt.Errorf("unknown webhooks subcommand: %s", args[0])
	}
}

func (c *CLI) webhooksList() error {
	options := &mailerlite.ListWebhookOptions{
		Page:  1,
		Limit: 100,
	}

	webhooks, _, err := c.client.sdk.Webhook.List(c.ctx, options)
	if err != nil {
		return fmt.Errorf("list webhooks: %w", err)
	}

	if c.githubIssue {
		c.println("## Webhooks")
		c.println()
		c.printf("Total: **%d**\n\n", len(webhooks.Data))
		if len(webhooks.Data) > 0 {
			c.println("| Name | URL | Events | Enabled |")
			c.println("|------|-----|--------|---------|")
			for _, w := range webhooks.Data {
				events := strings.Join(w.Events, ", ")
				c.printf("| %s | %s | %s | %v |\n", w.Name, w.Url, events, w.Enabled)
			}
		}
		return nil
	}

	c.printf("Total webhooks: %d\n\n", len(webhooks.Data))

	w := tabwriter.NewWriter(c.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tURL\tEVENTS\tENABLED")
	for _, wh := range webhooks.Data {
		events := strings.Join(wh.Events, ", ")
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%v\n", wh.Id, wh.Name, wh.Url, events, wh.Enabled)
	}
	w.Flush()

	return nil
}

func (c *CLI) webhooksCreate(url string, events []string) error {
	webhook := &mailerlite.CreateWebhookOptions{
		Name:   "CLI Webhook",
		Events: events,
		Url:    url,
	}

	result, _, err := c.client.sdk.Webhook.Create(c.ctx, webhook)
	if err != nil {
		return fmt.Errorf("create webhook: %w", err)
	}

	w := result.Data
	if c.githubIssue {
		c.println("## Webhook Created")
		c.println()
		c.printf("- **ID:** %s\n", w.Id)
		c.printf("- **URL:** %s\n", w.Url)
		c.printf("- **Events:** %s\n", strings.Join(w.Events, ", "))
		return nil
	}

	c.println("Webhook created:")
	c.printf("  ID:     %s\n", w.Id)
	c.printf("  URL:    %s\n", w.Url)
	c.printf("  Events: %s\n", strings.Join(w.Events, ", "))

	return nil
}

func (c *CLI) webhooksDelete(webhookID string) error {
	_, err := c.client.sdk.Webhook.Delete(c.ctx, webhookID)
	if err != nil {
		return fmt.Errorf("delete webhook: %w", err)
	}

	if c.githubIssue {
		c.println("## Webhook Deleted")
		c.println()
		c.printf("- **ID:** %s\n", webhookID)
		return nil
	}

	c.printf("Webhook deleted: %s\n", webhookID)
	return nil
}

// ============================================================================
// Automations Commands
// ============================================================================

func (c *CLI) handleAutomations(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("automations subcommand required: list, get")
	}

	switch args[0] {
	case "list":
		return c.automationsList()
	case "get":
		if len(args) < 2 {
			return fmt.Errorf("automation ID required: automations get ID")
		}
		return c.automationsGet(args[1])
	default:
		return fmt.Errorf("unknown automations subcommand: %s", args[0])
	}
}

func (c *CLI) automationsList() error {
	options := &mailerlite.ListAutomationOptions{
		Page:  1,
		Limit: 100,
	}

	automations, _, err := c.client.sdk.Automation.List(c.ctx, options)
	if err != nil {
		return fmt.Errorf("list automations: %w", err)
	}

	if c.githubIssue {
		c.println("## Automations")
		c.println()
		c.printf("Total: **%d**\n\n", len(automations.Data))
		if len(automations.Data) > 0 {
			c.println("| Name | Enabled | Emails | Opens | Clicks |")
			c.println("|------|---------|--------|-------|--------|")
			for _, auto := range automations.Data {
				c.printf("| %s | %v | %d | %d | %d |\n",
					auto.Name, auto.Enabled, auto.Stats.Sent, auto.Stats.OpensCount, auto.Stats.ClicksCount)
			}
		}
		return nil
	}

	c.printf("Total automations: %d\n\n", len(automations.Data))

	w := tabwriter.NewWriter(c.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tENABLED\tEMAILS\tOPENS\tCLICKS")
	for _, auto := range automations.Data {
		fmt.Fprintf(w, "%s\t%s\t%v\t%d\t%d\t%d\n",
			auto.ID, auto.Name, auto.Enabled, auto.Stats.Sent, auto.Stats.OpensCount, auto.Stats.ClicksCount)
	}
	w.Flush()

	return nil
}

func (c *CLI) automationsGet(automationID string) error {
	automation, _, err := c.client.sdk.Automation.Get(c.ctx, automationID)
	if err != nil {
		return fmt.Errorf("get automation: %w", err)
	}

	auto := automation.Data
	if c.githubIssue {
		c.printf("## Automation: %s\n\n", auto.Name)
		c.printf("- **ID:** %s\n", auto.ID)
		c.printf("- **Enabled:** %v\n", auto.Enabled)
		c.printf("- **Emails Sent:** %d\n", auto.Stats.Sent)
		c.printf("- **Opens:** %d\n", auto.Stats.OpensCount)
		c.printf("- **Clicks:** %d\n", auto.Stats.ClicksCount)
		return nil
	}

	c.printf("Automation: %s\n", auto.Name)
	c.println(strings.Repeat("=", 40))
	c.printf("ID:          %s\n", auto.ID)
	c.printf("Enabled:     %v\n", auto.Enabled)
	c.printf("Emails Sent: %d\n", auto.Stats.Sent)
	c.printf("Opens:       %d\n", auto.Stats.OpensCount)
	c.printf("Clicks:      %d\n", auto.Stats.ClicksCount)

	return nil
}

// ============================================================================
// Stats Command
// ============================================================================

func (c *CLI) handleStats() error {
	count, _, err := c.client.sdk.Subscriber.Count(c.ctx)
	if err != nil {
		return fmt.Errorf("get subscriber count: %w", err)
	}

	groupOptions := &mailerlite.ListGroupOptions{
		Page:  1,
		Limit: 100,
	}
	groups, _, err := c.client.sdk.Group.List(c.ctx, groupOptions)
	if err != nil {
		return fmt.Errorf("list groups: %w", err)
	}

	if c.githubIssue {
		c.println("## MailerLite Statistics")
		c.println()
		c.println("| Metric | Value |")
		c.println("|--------|-------|")
		c.printf("| Total Subscribers | %d |\n", count.Total)
		c.printf("| Total Groups | %d |\n", groups.Meta.Total)
		c.println()
		c.printf("*Generated: %s*\n", time.Now().Format("2006-01-02 15:04:05"))
		return nil
	}

	c.println("MailerLite Statistics")
	c.println(strings.Repeat("=", 40))
	c.printf("Total Subscribers: %d\n", count.Total)
	c.printf("Total Groups:      %d\n", groups.Meta.Total)
	c.println()
	c.printf("Generated: %s\n", time.Now().Format("2006-01-02 15:04:05"))

	return nil
}

// ============================================================================
// Open Command
// ============================================================================

func (c *CLI) handleOpen(args []string) error {
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

	c.printf("Opening %s...\n", url)
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
	default:
		cmd = "xdg-open"
		args = []string{url}
	}

	return exec.Command(cmd, args...).Start()
}
