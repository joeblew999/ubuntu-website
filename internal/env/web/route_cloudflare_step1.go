package web

import (
	"github.com/go-via/via"
	"github.com/go-via/via/h"
	"github.com/joeblew999/ubuntu-website/internal/env"
)

// cloudflareStep1Page - API Token setup (Step 1 of 4)
func cloudflareStep1Page(c *via.Context, cfg *env.EnvConfig, mockMode bool) {
	svc := env.NewService(mockMode)

	// Create fields for API token and token name
	fields := CreateFormFields(c, cfg, []string{
		env.KeyCloudflareAPIToken,
		env.KeyCloudflareAPITokenName,
	})

	saveMessage := c.Signal("")

	// Next action - validate and go to step 2
	nextAction := c.Action(func() {
		saveMessage.SetValue("")

		fieldUpdates := map[string]string{
			env.KeyCloudflareAPIToken:     fields[0].ValueSignal.String(),
			env.KeyCloudflareAPITokenName: fields[1].ValueSignal.String(),
		}

		results, err := svc.ValidateAndUpdateFields(fieldUpdates)
		UpdateValidationStatus(results, fields, c)

		if err != nil {
			saveMessage.SetValue("error:" + err.Error())
			c.Sync()
			return
		}

		// Check for validation errors
		if HasValidationErrors(results, fieldUpdates) {
			saveMessage.SetValue("error:Please fix validation errors before continuing")
			c.Sync()
			return
		}

		// Validation passed - redirect to step 2
		saveMessage.SetValue("success:Token validated! Moving to step 2...")
		c.Sync()
		c.ExecScript("window.location.href = '/cloudflare/step2'")
	})

	c.View(func() h.H {
		return h.Main(
			h.Class("container"),
			h.H1(h.Text("Cloudflare Setup - Step 1 of 4")),
			h.P(h.Text("API Token")),

			RenderNavigation("cloudflare"),

			h.H2(h.Text("Create API Token")),
			h.Ol(
				h.Li(h.Text("Visit: "), h.A(h.Href(env.CloudflareAPITokensURL), h.Attr("target", "_blank"), h.Text("Cloudflare API Tokens ↗"))),
				h.Li(h.Text("Click 'Create Token'")),
				h.Li(h.Text("Under 'Custom Token', click 'Get started'")),
				h.Li(h.Text("Give your token a descriptive name (e.g., 'Production Deploy Token')")),
				h.Li(h.Text("Under Permissions, add: Account → Cloudflare Pages → Edit")),
				h.Li(h.Text("Click 'Continue to summary' → 'Create Token'")),
				h.Li(h.Text("Copy the token value and paste below (save securely - you won't see it again!)")),
			),

			h.H3(h.Text("Enter your API Token:")),
			RenderFormField(fields[0]),

			h.H3(h.Text("Token Name:")),
			h.P(h.Text("Enter the name you gave this token in Cloudflare (helps you remember which token this is).")),
			RenderFormField(fields[1]),

			h.Div(
				h.Style("margin-top: 2rem;"),
				h.Button(h.Text("Next: Account ID →"), nextAction.OnClick()),
				h.Text(" or "),
				h.A(h.Href("/cloudflare/step2"), h.Text("Skip validation")),
			),

			RenderSaveMessage(saveMessage)[0],
			RenderSaveMessage(saveMessage)[1],
		)
	})
}
