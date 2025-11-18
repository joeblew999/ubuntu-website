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
			status.SetValue("❌ Cloudflare API Token is not configured (complete Step 1)")
			c.Sync()
			return
		}
		if accountID == "" || env.IsPlaceholder(accountID) {
			status.SetValue("❌ Account ID is not configured (complete Step 2)")
			c.Sync()
			return
		}
		if projectName == "" || env.IsPlaceholder(projectName) {
			status.SetValue("❌ Project Name is not configured (complete Step 4)")
			c.Sync()
			return
		}
		if customDomain == "" || env.IsPlaceholder(customDomain) {
			status.SetValue("❌ Custom Domain is not configured (complete Step 3)")
			c.Sync()
			return
		}

		isAttaching.SetValue(true)
		status.SetValue("🔄 Attaching custom domain...")
		c.Sync()

		// Add domain via Cloudflare Pages API
		err := env.AddPagesDomain(apiToken, accountID, projectName, customDomain)

		isAttaching.SetValue(false)
		if err != nil {
			status.SetValue("❌ Failed to attach domain: " + err.Error())
		} else {
			status.SetValue("✅ Successfully attached " + customDomain + " - Reloading page to show updated domains...")
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
					status.SetValue("❌ Configuration incomplete")
					c.Sync()
					return
				}

				isRemoving.SetValue(true)
				status.SetValue("🔄 Removing domain " + domainName + "...")
				c.Sync()

				err := env.DeletePagesDomain(apiToken, accountID, projectName, domainName)

				isRemoving.SetValue(false)
				if err != nil {
					status.SetValue("❌ Failed to remove domain: " + err.Error())
				} else {
					status.SetValue("✅ Successfully removed " + domainName + " - Reloading page to show updated domains...")
					c.Sync()
					// Reload page to refresh domains list
					c.ExecScript("setTimeout(function() { window.location.reload(); }, 1500);")
					return
				}
				c.Sync()
			})

			domainListElements = append(domainListElements, h.Tr(
				h.Style("border-bottom: 1px solid #eee;"),
				h.Td(
					h.Style("padding: 0.75rem; font-family: monospace; font-weight: bold; font-size: 1.1em;"),
					h.Text(domainName),
				),
				h.Td(
					h.Style("padding: 0.75rem;"),
					h.Span(
						h.Style(func() string {
							if domainStatus == "active" {
								return "color: #28a745; font-weight: bold;"
							}
							return "color: #ffc107; font-weight: bold;"
						}()),
						h.Text(domainStatus),
					),
				),
				h.Td(
					h.Style("padding: 0.75rem; text-align: right;"),
					h.Button(
						h.Style("background: #dc3545; color: white; padding: 0.25rem 1rem; border: none; border-radius: 4px; cursor: pointer;"),
						h.If(isRemoving.String() == "true", h.Attr("disabled", "disabled")),
						h.Text("Remove"),
						removeAction.OnClick(),
					),
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
				h.Style("margin-bottom: 2rem; background: #e7f3ff; border-left: 4px solid #0066cc; padding: 1rem;"),
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
				h.Style("margin-bottom: 2rem; background: #f5f5f5; padding: 1rem;"),
				h.H4(h.Text("🔧 Current Configuration")),
				h.Table(
					h.Style("width: 100%;"),
					h.TBody(
						h.Tr(
							h.Td(h.Style("padding: 0.5rem; font-weight: bold; width: 180px;"), h.Text("Account ID:")),
							h.Td(
								h.Style("padding: 0.5rem; font-family: monospace; color: #333;"),
								h.Text(func() string {
									if env.IsPlaceholder(accountID) {
										return "Not configured"
									}
									return accountID
								}()),
							),
						),
						h.Tr(
							h.Td(h.Style("padding: 0.5rem; font-weight: bold;"), h.Text("Project Name:")),
							h.Td(
								h.Style("padding: 0.5rem; font-family: monospace; color: #333;"),
								h.Text(func() string {
									if env.IsPlaceholder(projectName) {
										return "Not configured"
									}
									return projectName
								}()),
							),
						),
						h.Tr(
							h.Td(h.Style("padding: 0.5rem; font-weight: bold;"), h.Text("Custom Domain:")),
							h.Td(
								h.Style("padding: 0.5rem; font-family: monospace; color: #333;"),
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
					h.Style("margin-bottom: 2rem;"),
					h.Button(
						h.If(isAttaching.String() == "true", h.Attr("disabled", "disabled")),
						h.Text(func() string {
							if isAttaching.String() == "true" {
								return "⏳ Attaching Domain..."
							}
							return "🔗 Attach " + customDomain
						}()),
						attachDomainAction.OnClick(),
					),
				),
			),

			// Status Message
			h.If(status.String() != "",
				h.Article(
					h.Style(func() string {
						statusVal := status.String()
						baseStyle := "padding: 1rem; margin-bottom: 2rem; border-radius: 4px;"
						if strings.HasPrefix(statusVal, "✅") {
							return baseStyle + " background: #d4edda; border: 1px solid #c3e6cb; color: #155724;"
						} else if strings.HasPrefix(statusVal, "❌") {
							return baseStyle + " background: #f8d7da; border: 1px solid #f5c6cb; color: #721c24;"
						} else if strings.HasPrefix(statusVal, "⚠️") {
							return baseStyle + " background: #fff3cd; border: 1px solid #ffeeba; color: #856404;"
						}
						return baseStyle + " background: #d1ecf1; border: 1px solid #bee5eb; color: #0c5460;"
					}()),
					h.Text(status.String()),
				),
			),

			// Attached Domains List
			h.If(len(domainsCache) > 0,
				h.Div(
					h.Style("margin-bottom: 2rem;"),
					h.H4(h.Text("🌐 Attached Custom Domains")),
					h.Table(
						h.Style("width: 100%; border-collapse: collapse;"),
						h.THead(
							h.Tr(
								h.Th(h.Style("text-align: left; padding: 0.75rem; border-bottom: 2px solid #ddd;"), h.Text("Domain")),
								h.Th(h.Style("text-align: left; padding: 0.75rem; border-bottom: 2px solid #ddd;"), h.Text("Status")),
								h.Th(h.Style("text-align: right; padding: 0.75rem; border-bottom: 2px solid #ddd;"), h.Text("Actions")),
							),
						),
						h.TBody(domainListElements...),
					),
				),
			),

			// Navigation
			h.Div(
				h.Style("margin-top: 2rem;"),
				h.A(h.Href("/cloudflare/step4"), h.Text("← Back: Project Selection")),
				h.Text(" "),
				h.A(
					h.Href("/deploy"),
					h.Style("background: #28a745; color: white; padding: 0.5rem 1rem; border-radius: 4px; text-decoration: none; font-weight: bold;"),
					h.Text("✅ Complete Setup - Go to Deploy →"),
				),
			),
		)
	})
}
