package web

import (
	"strings"
	"time"

	"github.com/go-via/via"
	"github.com/go-via/via/h"
	"github.com/joeblew999/ubuntu-website/internal/env"
)

// cloudflarePage creates the Cloudflare-only setup page
func cloudflarePage(c *via.Context, cfg *env.EnvConfig, mockMode bool) {
	// Create service for config operations
	svc := env.NewService(mockMode)

	// Create signals for Cloudflare fields - use passed config
	cloudflareToken := c.Signal(cfg.CloudflareToken)
	cloudflareTokenName := c.Signal(cfg.CloudflareTokenName)
	cloudflareAccount := c.Signal(cfg.CloudflareAccount)
	cloudflareProject := c.Signal(cfg.CloudflareProject)

	// Validation status signals
	cloudflareTokenStatus := c.Signal("")
	cloudflareAccountStatus := c.Signal("")

	// Save message signal
	saveMessage := c.Signal("")

	// Validate function - validates current field values using service
	validateFields := func() {
		// Build current config for validation
		currentCfg := &env.EnvConfig{
			CloudflareToken:     cloudflareToken.String(),
			CloudflareTokenName: cloudflareTokenName.String(),
			CloudflareAccount:   cloudflareAccount.String(),
			CloudflareProject:   cloudflareProject.String(),
		}

		// Validate using service
		results := svc.ValidateConfig(currentCfg)

		// Update status signals
		result := results[env.KeyCloudflareAPIToken]
		if result.Skipped {
			cloudflareTokenStatus.SetValue("")
		} else if result.Valid {
			cloudflareTokenStatus.SetValue("valid")
		} else {
			cloudflareTokenStatus.SetValue("invalid")
		}

		result = results[env.KeyCloudflareAccountID]
		if result.Skipped {
			cloudflareAccountStatus.SetValue("")
		} else if result.Valid {
			cloudflareAccountStatus.SetValue("valid")
		} else {
			cloudflareAccountStatus.SetValue("invalid")
		}

		c.SyncSignals()
	}

	// Run initial validation on page load
	validateFields()

	// Save and validate action
	saveAction := c.Action(func() {
		// Prepare field updates
		fieldUpdates := map[string]string{
			env.KeyCloudflareAPIToken:     cloudflareToken.String(),
			env.KeyCloudflareAPITokenName: cloudflareTokenName.String(),
			env.KeyCloudflareAccountID:   cloudflareAccount.String(),
			env.KeyCloudflarePageProject:   cloudflareProject.String(),
		}

		// Use service to validate and save atomically
		results, err := svc.ValidateAndUpdateFields(fieldUpdates)

		// Update validation status from results
		if result, ok := results[env.KeyCloudflareAPIToken]; ok {
			if result.Skipped {
				cloudflareTokenStatus.SetValue("")
			} else if result.Valid {
				cloudflareTokenStatus.SetValue("valid")
			} else {
				cloudflareTokenStatus.SetValue("invalid")
			}
		}

		if result, ok := results[env.KeyCloudflareAccountID]; ok {
			if result.Skipped {
				cloudflareAccountStatus.SetValue("")
			} else if result.Valid {
				cloudflareAccountStatus.SetValue("valid")
			} else {
				cloudflareAccountStatus.SetValue("invalid")
			}
		}

		c.SyncSignals()

		// Handle save result
		if err != nil {
			saveMessage.SetValue("error:" + err.Error())
		} else {
			// Check if there were validation errors
			hasErrors := false
			for _, result := range results {
				if !result.Skipped && !result.Valid {
					hasErrors = true
					break
				}
			}

			if hasErrors {
				saveMessage.SetValue("error:Please fix validation errors before saving")
			} else {
				saveMessage.SetValue("success:Cloudflare configuration saved successfully!")
			}
		}

		c.SyncSignals()

		// Clear message after 5 seconds
		time.AfterFunc(5*time.Second, func() {
			saveMessage.SetValue("")
			c.SyncSignals()
		})
	})

	c.View(func() h.H {
		return h.Main(
			h.Class("container"),
			h.H1(h.Text("Cloudflare Setup")),
			h.P(h.Text("Configure your Cloudflare credentials for deployment to Cloudflare Pages")),

			// Navigation
			h.Nav(
				h.Ul(
					h.Li(h.A(h.Href("/"), h.Text("All Settings"))),
					h.Li(h.Strong(h.Text("Cloudflare Only"))),
					h.Li(h.A(h.Href("/claude"), h.Text("Claude Only"))),
				),
			),

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
				h.Li(h.Text("Select your account from the left sidebar")),
				h.Li(h.Text("Copy the Account ID from the right sidebar under 'Account ID'")),
			),
			h.P(h.Strong(h.Text("Step 3: Enter Credentials Below"))),
			h.Ul(
				h.Li(h.Text("Paste your API token and Account ID into the fields below")),
				h.Li(h.Text("Optionally set a token name to remember which token this is")),
				h.Li(h.Text("Optionally set your Cloudflare Pages project name")),
				h.Li(h.Text("Click 'Save Cloudflare Configuration' to validate and save")),
			),

			// Cloudflare Section
			h.H2(h.Text("Cloudflare Credentials")),
			h.Div(
				h.Label(h.Text(env.GetFieldLabel(env.KeyCloudflareAPIToken))),
				h.Input(h.Type("text"), h.Value(cloudflareToken.String()), cloudflareToken.Bind()),
				h.If(cloudflareTokenStatus.String() == "valid", h.Span(h.Text("✓"))),
				h.If(cloudflareTokenStatus.String() == "invalid", h.Span(h.Text("✗"))),
			),
			h.Div(
				h.Label(h.Text(env.GetFieldLabel(env.KeyCloudflareAPITokenName))),
				h.Input(h.Type("text"), h.Value(cloudflareTokenName.String()), cloudflareTokenName.Bind()),
			),
			h.Div(
				h.Label(h.Text(env.GetFieldLabel(env.KeyCloudflareAccountID))),
				h.Input(h.Type("text"), h.Value(cloudflareAccount.String()), cloudflareAccount.Bind()),
				h.If(cloudflareAccountStatus.String() == "valid", h.Span(h.Text("✓"))),
				h.If(cloudflareAccountStatus.String() == "invalid", h.Span(h.Text("✗"))),
			),
			h.Div(
				h.Label(h.Text(env.GetFieldLabel(env.KeyCloudflarePageProject))),
				h.Input(h.Type("text"), h.Value(cloudflareProject.String()), cloudflareProject.Bind()),
			),

			// Action buttons
			h.Div(
				h.Button(h.Text("Save Cloudflare Configuration"), saveAction.OnClick()),
			),

			// Save message - rendered inline using h.If for reactivity
			// Note: Check string prefix using strings.HasPrefix to avoid slice bounds issues
			h.If(strings.HasPrefix(saveMessage.String(), "error:"),
				h.Div(h.Text("❌ "+strings.TrimPrefix(saveMessage.String(), "error:"))),
			),
			h.If(strings.HasPrefix(saveMessage.String(), "success:"),
				h.Div(h.Text("✅ "+strings.TrimPrefix(saveMessage.String(), "success:"))),
			),
		)
	})
}
