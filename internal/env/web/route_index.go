package web

import (
	"github.com/go-via/via"
	"github.com/go-via/via/h"
	"github.com/joeblew999/ubuntu-website/internal/env"
)

// homePage creates a static welcome page with navigation (no reactive forms)
func homePage(c *via.Context, cfg *env.EnvConfig, mockMode bool) {
	c.View(func() h.H {
		return h.Main(
			h.Class("container"),
			h.H1(h.Text("Environment Setup")),
			h.P(h.Text("Configure your Cloudflare and Claude credentials for deployment and translation")),

			// Status overview (non-interactive)
			h.H2(h.Text("Configuration Status")),
			h.P(h.Text("Click on a section below to configure credentials:")),

			// Navigation cards
			h.Div(
				h.Style("display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 1rem; margin: 2rem 0;"),

				// Cloudflare Card
				h.Article(
					h.A(
						h.Href("/cloudflare"),
						h.Style("text-decoration: none; color: inherit;"),
						h.H3(h.Text("Cloudflare Setup")),
						h.P(h.Text("Configure Cloudflare credentials for deployment to Cloudflare Pages")),
					),
				),

				// Claude Card
				h.Article(
					h.A(
						h.Href("/claude"),
						h.Style("text-decoration: none; color: inherit;"),
						h.H3(h.Text("Claude AI Setup")),
						h.P(h.Text("Configure Claude AI credentials for content translation")),
					),
				),
			),

			// Current Configuration Summary (non-interactive)
			h.H2(h.Text("Current Configuration")),
			h.Dl(
				h.Dt(h.Strong(h.Text(env.GetDisplayName(env.KeyCloudflareAPIToken)))),
				h.Dd(renderConfigValue(cfg.CloudflareToken)),

				h.Dt(h.Strong(h.Text(env.GetDisplayName(env.KeyCloudflareAPITokenName)))),
				h.Dd(renderConfigValue(cfg.CloudflareTokenName)),

				h.Dt(h.Strong(h.Text(env.GetDisplayName(env.KeyCloudflareAccountID)))),
				h.Dd(renderConfigValue(cfg.CloudflareAccount)),

				h.Dt(h.Strong(h.Text(env.GetDisplayName(env.KeyCloudflarePageProject)))),
				h.Dd(renderConfigValue(cfg.CloudflareProject)),

				h.Dt(h.Strong(h.Text(env.GetDisplayName(env.KeyClaudeAPIKey)))),
				h.Dd(renderConfigValue(cfg.ClaudeAPIKey)),

				h.Dt(h.Strong(h.Text(env.GetDisplayName(env.KeyClaudeWorkspaceName)))),
				h.Dd(renderConfigValue(cfg.ClaudeWorkspace)),
			),
		)
	})
}

// renderConfigValue displays a config value or "(not set)" if it's a placeholder
func renderConfigValue(value string) h.H {
	if env.IsPlaceholder(value) || value == "" {
		return h.Text("(not set)")
	}
	// Truncate long values for display
	if len(value) > 24 {
		return h.Text(value[:10] + "..." + value[len(value)-10:])
	}
	return h.Text(value)
}
