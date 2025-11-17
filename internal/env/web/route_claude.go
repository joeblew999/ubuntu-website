package web

import (
	"strings"
	"time"

	"github.com/go-via/via"
	"github.com/go-via/via/h"
	"github.com/joeblew999/ubuntu-website/internal/env"
)

// claudePage creates the Claude-only setup page
func claudePage(c *via.Context, cfg *env.EnvConfig, mockMode bool) {
	// Create service for config operations
	svc := env.NewService(mockMode)

	// Create signals for Claude fields - use passed config
	claudeAPIKey := c.Signal(cfg.ClaudeAPIKey)
	claudeWorkspace := c.Signal(cfg.ClaudeWorkspace)

	// Validation status signal
	claudeAPIKeyStatus := c.Signal("")

	// Save message signal
	saveMessage := c.Signal("")

	// Validate function - validates current field values using service
	validateFields := func() {
		// Build current config for validation
		currentCfg := &env.EnvConfig{
			ClaudeAPIKey:    claudeAPIKey.String(),
			ClaudeWorkspace: claudeWorkspace.String(),
		}

		// Validate using service
		results := svc.ValidateConfig(currentCfg)

		// Update status signals
		result := results[env.KeyClaudeAPIKey]
		if result.Skipped {
			claudeAPIKeyStatus.SetValue("")
		} else if result.Valid {
			claudeAPIKeyStatus.SetValue("valid")
		} else {
			claudeAPIKeyStatus.SetValue("invalid")
		}

		c.SyncSignals()
	}

	// Run initial validation on page load
	validateFields()

	// Save and validate action
	saveAction := c.Action(func() {
		// Prepare field updates
		fieldUpdates := map[string]string{
			env.KeyClaudeAPIKey:    claudeAPIKey.String(),
			env.KeyClaudeWorkspaceName: claudeWorkspace.String(),
		}

		// Use service to validate and save atomically
		results, err := svc.ValidateAndUpdateFields(fieldUpdates)

		// Update validation status from results
		if result, ok := results[env.KeyClaudeAPIKey]; ok {
			if result.Skipped {
				claudeAPIKeyStatus.SetValue("")
			} else if result.Valid {
				claudeAPIKeyStatus.SetValue("valid")
			} else {
				claudeAPIKeyStatus.SetValue("invalid")
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
				saveMessage.SetValue("success:Claude configuration saved successfully!")
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
			h.H1(h.Text("Claude AI Setup")),
			h.P(h.Text("Configure your Claude AI credentials for content translation")),

			// Navigation
			h.Nav(
				h.Ul(
					h.Li(h.A(h.Href("/"), h.Text("All Settings"))),
					h.Li(h.A(h.Href("/cloudflare"), h.Text("Cloudflare Only"))),
					h.Li(h.Strong(h.Text("Claude Only"))),
				),
			),

			// Setup instructions
			h.H2(h.Text("Setup Instructions")),
			h.P(h.Strong(h.Text("Step 1: Sign up for Claude API"))),
			h.Ul(
				h.Li(h.Text("Visit: "), h.A(h.Href(env.AnthropicConsoleURL), h.Attr("target", "_blank"), h.Text("Claude Console"))),
				h.Li(h.Text("Create an account if you don't have one")),
				h.Li(h.Text("Verify your email address")),
				h.Li(h.Text("Add billing information at: "), h.A(h.Href(env.AnthropicBillingURL), h.Attr("target", "_blank"), h.Text("Billing Settings"))),
			),
			h.P(h.Strong(h.Text("Step 2: Create API Key"))),
			h.Ul(
				h.Li(h.Text("Visit: "), h.A(h.Href(env.AnthropicAPIKeysURL), h.Attr("target", "_blank"), h.Text("API Keys"))),
				h.Li(h.Text("Click 'Create Key' button")),
				h.Li(h.Text("Give your key a descriptive name")),
				h.Li(h.Text("Copy the API key (save it securely - you won't see it again)")),
			),
			h.P(h.Strong(h.Text("Step 3: Find Workspace Name (Optional)"))),
			h.Ul(
				h.Li(h.Text("Visit: "), h.A(h.Href(env.AnthropicWorkspacesURL), h.Attr("target", "_blank"), h.Text("Workspaces"))),
				h.Li(h.Text("Find your workspace name in the list")),
				h.Li(h.Text("This is optional but helps organize your API usage")),
			),
			h.P(h.Strong(h.Text("Step 4: Enter Credentials Below"))),
			h.Ul(
				h.Li(h.Text("Paste your API key into the field below")),
				h.Li(h.Text("Optionally enter your workspace name")),
				h.Li(h.Text("Click 'Save Claude Configuration' to validate and save")),
			),

			// Claude Section
			h.H2(h.Text("Claude API Credentials")),
			h.Div(
				h.Label(h.Text(env.GetFieldLabel(env.KeyClaudeAPIKey))),
				h.Input(h.Type("text"), h.Value(claudeAPIKey.String()), claudeAPIKey.Bind()),
				h.If(claudeAPIKeyStatus.String() == "valid", h.Span(h.Text("✓"))),
				h.If(claudeAPIKeyStatus.String() == "invalid", h.Span(h.Text("✗"))),
			),
			h.Div(
				h.Label(h.Text(env.GetFieldLabel(env.KeyClaudeWorkspaceName))),
				h.Input(h.Type("text"), h.Value(claudeWorkspace.String()), claudeWorkspace.Bind()),
			),

			// Action buttons
			h.Div(
				h.Button(h.Text("Save Claude Configuration"), saveAction.OnClick()),
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
