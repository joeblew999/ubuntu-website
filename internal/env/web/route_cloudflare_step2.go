package web

import (
	"github.com/go-via/via"
	"github.com/go-via/via/h"
	"github.com/joeblew999/ubuntu-website/internal/env"
)

// cloudflareStep2Page - Account ID setup (Step 2 of 5)
func cloudflareStep2Page(c *via.Context, cfg *env.EnvConfig, mockMode bool) {
	svc := env.NewService(mockMode)

	// Create fields for token (needed for validation) and account ID
	fields := CreateFormFields(c, cfg, []string{
		env.KeyCloudflareAPIToken,
		env.KeyCloudflareAccountID,
	})

	saveMessage := c.Signal("")

	// Next action - validate and go to step 3
	nextAction := c.Action(func() {
		saveMessage.SetValue("")

		fieldUpdates := map[string]string{
			env.KeyCloudflareAPIToken:  fields[0].ValueSignal.String(),
			env.KeyCloudflareAccountID: fields[1].ValueSignal.String(),
		}

		results, err := svc.ValidateAndUpdateFields(fieldUpdates)
		UpdateValidationStatus(results, []FormFieldData{fields[1]}, c)

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

		// Validation passed - redirect to step 3
		saveMessage.SetValue("success:Account ID validated! Moving to step 3...")
		c.Sync()
		c.ExecScript("window.location.href = '/cloudflare/step3'")
	})

	c.View(func() h.H {
		return h.Main(
			h.Class("container"),
			h.H1(h.Text("Cloudflare Setup - Step 2 of 5")),
			h.P(h.Text("Account ID")),

			RenderNavigation("cloudflare"),

			h.H2(h.Text("Find Account ID")),
			h.Ol(
				h.Li(h.Text("Visit: "), h.A(h.Href(env.CloudflareDashboardURL), h.Attr("target", "_blank"), h.Text("Cloudflare Dashboard ↗"))),
				h.Li(h.Text("Find 'Account Home' in the left sidebar")),
				h.Li(h.Text("The Account ID is in the right sidebar")),
				h.Li(h.Text("It's a 32-character hex string - click copy icon")),
			),

			h.H3(h.Text("Enter your Account ID:")),
			RenderFormField(fields[1]),

			h.Div(
				h.Style("margin-top: 2rem;"),
				h.A(h.Href("/cloudflare"), h.Text("← Back: API Token")),
				h.Text(" "),
				h.Button(h.Text("Next: Domain Selection →"), nextAction.OnClick()),
				h.Text(" or "),
				h.A(h.Href("/cloudflare/step3"), h.Text("Skip validation")),
			),

			RenderErrorMessage(saveMessage),
			RenderSuccessMessage(saveMessage),
		)
	})
}
