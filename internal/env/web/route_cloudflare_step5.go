package web

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-via/via"
	"github.com/go-via/via/h"
	"github.com/joeblew999/ubuntu-website/internal/env"
)

// cloudflareStep5Page - Custom Domain setup (Step 5 of 5)
func cloudflareStep5Page(c *via.Context, cfg *env.EnvConfig, mockMode bool) {
	// Get required values from config
	apiToken := cfg.Get(env.KeyCloudflareAPIToken)
	accountID := cfg.Get(env.KeyCloudflareAccountID)
	projectName := cfg.Get(env.KeyCloudflarePageProject)
	customDomain := cfg.Get(env.KeyCloudflareDomain)

	// Signals
	status := c.Signal("")
	isAttaching := c.Signal(false)
	isRemoving := c.Signal(false)

	// Domains loader - populated lazily when first accessed
	domainsLoader := NewLazyLoader(func() ([]env.PagesDomain, error) {
		if apiToken == "" || env.IsPlaceholder(apiToken) ||
			accountID == "" || env.IsPlaceholder(accountID) ||
			projectName == "" || env.IsPlaceholder(projectName) {
			return []env.PagesDomain{}, fmt.Errorf("configuration incomplete")
		}

		if mockMode {
			// Mock data for testing
			return []env.PagesDomain{
				{Name: "ubuntusoftware.net", Status: "active"},
				{Name: "example.com", Status: "pending"},
			}, nil
		}

		return env.ListPagesDomains(apiToken, accountID, projectName)
	})

	// Attach domain action
	attachDomainAction := c.Action(func() {
		// Validation
		if apiToken == "" || env.IsPlaceholder(apiToken) {
			status.SetValue("error:Cloudflare API Token is not configured (complete Step 1)")
			c.Sync()
			return
		}
		if accountID == "" || env.IsPlaceholder(accountID) {
			status.SetValue("error:Account ID is not configured (complete Step 2)")
			c.Sync()
			return
		}
		if projectName == "" || env.IsPlaceholder(projectName) {
			status.SetValue("error:Project Name is not configured (complete Step 4)")
			c.Sync()
			return
		}
		if customDomain == "" || env.IsPlaceholder(customDomain) {
			status.SetValue("error:Custom Domain is not configured (complete Step 3)")
			c.Sync()
			return
		}

		isAttaching.SetValue(true)
		status.SetValue("info:Attaching custom domain...")
		c.Sync()

		// Add domain via Cloudflare Pages API
		err := env.AddPagesDomain(apiToken, accountID, projectName, customDomain)

		isAttaching.SetValue(false)
		if err != nil {
			status.SetValue("error:Failed to attach domain: " + err.Error())
		} else {
			status.SetValue("success:Successfully attached " + customDomain + " - Reloading page to show updated domains...")
			c.Sync()
			// Reload page to refresh domains list
			c.ExecScript("setTimeout(function() { window.location.reload(); }, 1500);")
			return
		}
		c.Sync()
	})

	c.View(func() h.H {
		// Load domains using LazyLoader
		domainsCache, domainsErr := domainsLoader.Get()
		if domainsErr != nil {
			log.Printf("Failed to fetch Pages domains: %v", domainsErr)
		}

		// Build domain list UI elements
		domainListElements := make([]h.H, 0, len(domainsCache))
		for _, domain := range domainsCache {
			domainName := domain.Name     // Capture in closure
			domainStatus := domain.Status // Capture in closure

			// Create remove action for this specific domain
			removeAction := c.Action(func() {
				if apiToken == "" || env.IsPlaceholder(apiToken) ||
					accountID == "" || env.IsPlaceholder(accountID) ||
					projectName == "" || env.IsPlaceholder(projectName) {
					status.SetValue("error:Configuration incomplete")
					c.Sync()
					return
				}

				isRemoving.SetValue(true)
				status.SetValue("info:Removing domain " + domainName + "...")
				c.Sync()

				err := env.DeletePagesDomain(apiToken, accountID, projectName, domainName)

				isRemoving.SetValue(false)
				if err != nil {
					status.SetValue("error:Failed to remove domain: " + err.Error())
				} else {
					status.SetValue("success:Successfully removed " + domainName + " - Reloading page to show updated domains...")
					c.Sync()
					// Reload page to refresh domains list
					c.ExecScript("setTimeout(function() { window.location.reload(); }, 1500);")
					return
				}
				c.Sync()
			})

			domainListElements = append(domainListElements, h.Div(
				h.Style("display: flex; justify-content: space-between; align-items: center; padding: 0.75rem; background: var(--pico-card-background-color); border-radius: 0.25rem; margin-bottom: 0.5rem;"),
				h.Div(
					h.Strong(
						h.Style("font-family: monospace; font-size: 1.1em;"),
						h.Text(domainName),
					),
					h.Small(
						h.Style("margin-left: 1rem; color: var(--pico-muted-color);"),
						h.Text("Status: "),
						h.Span(
							h.Style(func() string {
								if domainStatus == "active" {
									return "color: var(--pico-ins-color); font-weight: bold;"
								}
								return "color: var(--pico-muted-color); font-weight: bold;"
							}()),
							h.Text(domainStatus),
						),
					),
				),
				h.Button(
					h.Attr("class", "secondary outline"),
					h.If(isRemoving.String() == "true", h.Attr("disabled", "disabled")),
					h.Text("Remove"),
					removeAction.OnClick(),
				),
			))
		}

		return h.Main(
			h.Class("container"),
			h.H1(h.Text("Cloudflare Setup - Step 5 of 5")),
			h.P(h.Text("Custom Domain Setup")),

			RenderNavigation("cloudflare"),

			// Instructions
			h.Article(
				h.Style("background-color: var(--pico-card-background-color); border-left: 4px solid var(--pico-primary); padding: 1rem; margin-bottom: 1rem;"),
				h.H4(h.Text("📖 Instructions")),
				h.Ul(
					h.Style("margin: 0.5rem 0 0 1.5rem;"),
					h.Li(h.Text("This step attaches your custom domain to the Cloudflare Pages project via the API")),
					h.Li(h.Text("This resolves CNAME Cross-User Banned errors (Error 1014)")),
					h.Li(h.Text("Your domain's DNS records in Cloudflare will automatically be configured")),
					h.Li(h.Text("HTTPS certificate provisioning may take a few minutes")),
					h.Li(h.Text("After attaching, your site will be live at both the preview URL and custom domain")),
				),
			),

			// Current Configuration
			h.Article(
				h.Style("background-color: var(--pico-card-background-color); padding: 1rem; margin-bottom: 1rem;"),
				h.H4(h.Text("🔧 Current Configuration")),
				h.Table(
					h.Tr(
						h.Td(h.Strong(h.Text("Account ID:"))),
						h.Td(
							h.Code(
								h.Style("color: var(--pico-color);"),
								h.Text(func() string {
									if env.IsPlaceholder(accountID) {
										return "Not configured"
									}
									return accountID
								}()),
							),
						),
					),
					h.Tr(
						h.Td(h.Strong(h.Text("Project Name:"))),
						h.Td(
							h.Code(
								h.Style("color: var(--pico-color);"),
								h.Text(func() string {
									if env.IsPlaceholder(projectName) {
										return "Not configured"
									}
									return projectName
								}()),
							),
						),
					),
					h.Tr(
						h.Td(h.Strong(h.Text("Custom Domain:"))),
						h.Td(
							h.Code(
								h.Style("color: var(--pico-color);"),
								h.Text(func() string {
									if env.IsPlaceholder(customDomain) {
										return "Not configured"
									}
									return customDomain
								}()),
							),
						),
					),
				),
			),

			// Attach Domain Button
			h.If(customDomain != "" && !env.IsPlaceholder(customDomain),
				h.Div(
					h.Style("margin-bottom: 1rem;"),
					h.Button(
						h.If(isAttaching.String() == "true", h.Attr("aria-busy", "true")),
						h.If(isAttaching.String() == "true", h.Attr("disabled", "disabled")),
						h.Text(func() string {
							if isAttaching.String() == "true" {
								return "Attaching Domain..."
							}
							return "🔗 Attach " + customDomain
						}()),
						attachDomainAction.OnClick(),
					),
				),
			),

			// Status Messages - use helper functions for proper PicoCSS styling
			RenderErrorMessage(status),
			RenderSuccessMessage(status),
			// Info message (for in-progress states)
			h.If(strings.HasPrefix(status.String(), "info:"),
				h.Article(
					h.Style("background-color: var(--pico-card-background-color); border-left: 4px solid var(--pico-primary); padding: 1rem; margin-top: 1rem;"),
					h.P(
						h.Style("margin: 0; color: var(--pico-color);"),
						h.Text(strings.TrimPrefix(status.String(), "info:")),
					),
				),
			),

			// Attached Domains List
			h.If(len(domainsCache) > 0,
				h.Div(
					h.Style("margin-top: 2rem;"),
					h.H3(h.Text("🌐 Attached Custom Domains")),
					h.Div(domainListElements...),
				),
			),

			// Navigation
			h.Div(
				h.Style("margin-top: 2rem;"),
				h.A(h.Href("/cloudflare/step4"), h.Text("← Back: Project Selection")),
				h.Text(" "),
				h.A(
					h.Href("/deploy"),
					h.Attr("role", "button"),
					h.Text("✅ Complete Setup - Go to Deploy →"),
				),
			),
		)
	})
}
