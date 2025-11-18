package web

import (
	"github.com/go-via/via"
	"github.com/go-via/via/h"
	"github.com/joeblew999/ubuntu-website/internal/env"
)

// cloudflarePage creates the Cloudflare-only setup page
func cloudflarePage(c *via.Context, cfg *env.EnvConfig, mockMode bool) {
	// Create service for config operations
	svc := env.NewService(mockMode)

	// Create form fields using helper
	fields := CreateFormFields(c, cfg, []string{
		env.KeyCloudflareAPIToken,
		env.KeyCloudflareAPITokenName,
		env.KeyCloudflareAccountID,
		env.KeyCloudflarePageProject,
	})

	// Save message signal
	saveMessage := c.Signal("")

	// Save and validate action using helper
	saveAction := c.Action(CreateSaveAction(c, svc, fields, saveMessage))

	c.View(func() h.H {
		return h.Main(
			h.Class("container"),
			h.H1(h.Text("Cloudflare Setup")),
			h.P(h.Text("Configure your Cloudflare credentials for deployment to Cloudflare Pages")),

			// Navigation using helper
			RenderNavigation("cloudflare"),

			// Setup instructions
			h.H2(h.Text("Setup Instructions")),
			h.P(h.Strong(h.Text("Step 1: Create API Token"))),
			h.Ul(
				h.Li(h.Text("Visit: "), h.A(h.Href(env.CloudflareAPITokensURL), h.Attr("target", "_blank"), h.Text("Cloudflare API Tokens"))),
				h.Li(h.Text("Click 'Create Token' button")),
				h.Li(h.Text("Select 'Create Custom Token' template")),
				h.Li(h.Text("Set permissions: Account > Cloudflare Pages > Edit")),
				h.Li(h.Text("Copy the token (save it securely - you won't see it again)")),
			),
			h.P(h.Strong(h.Text("Step 2: Find Account ID"))),
			h.Ul(
				h.Li(h.Text("Visit: "), h.A(h.Href(env.CloudflareDashboardURL), h.Attr("target", "_blank"), h.Text("Cloudflare Dashboard"))),
				h.Li(h.Text("Find 'Account Home' with your email address in the left sidebar")),
			h.Li(h.Text("Right-click on the icon next to 'Account Home' (this reveals the Account ID)")),
				h.Li(h.Text("The Account ID is a 32-character hex string - copy it from the right sidebar or click the copy icon next to it")),
			),
			h.P(h.Strong(h.Text("Step 3: Cloudflare Pages Project Name"))),
		h.Ul(
			h.Li(h.Text("If you already have a Cloudflare Pages project, just enter its name in the field below")),
			h.Li(h.Text("To create a new project: Visit "), h.A(h.Href(env.CloudflarePagesURL), h.Attr("target", "_blank"), h.Text("Workers & Pages")), h.Text(" > 'Create application' > 'Pages' > 'Connect to Git'")),
			h.Li(h.Text("Project names must be lowercase letters, numbers, and hyphens only (1-63 characters)")),
			h.Li(h.Text("Example: 'ubuntusoftware-net' or 'my-hugo-site'")),
			h.Li(h.Text("Note: You can leave this blank and create/connect your project later")),
		),
		h.P(h.Strong(h.Text("Step 4: Enter Configuration Below"))),
			h.Ul(
				h.Li(h.Text("Paste your API token and Account ID into the fields below")),
				h.Li(h.Text("Optionally set a token name to remember which token this is")),
				h.Li(h.Text("Enter your Cloudflare Pages project name (or leave blank if creating later)")),
				h.Li(h.Text("Click 'Save Cloudflare Configuration' to validate and save")),
			),

			// Cloudflare Section - render form fields using helpers
			h.H2(h.Text("Cloudflare Credentials")),
			RenderFormField(fields[0]),
			RenderFormField(fields[1]),
			RenderFormField(fields[2]),
			RenderFormField(fields[3]),

			// Action buttons
			h.Div(
				h.Button(h.Text("Save Cloudflare Configuration"), saveAction.OnClick()),
			),

			// Save message using helper
			RenderSaveMessage(saveMessage)[0],
			RenderSaveMessage(saveMessage)[1],
		)
	})
}
