package web

import (
	"github.com/go-via/via"
	"github.com/go-via/via/h"
	"github.com/joeblew999/ubuntu-website/internal/env"
)

// cloudflareStep3Page - Project details (Step 3 of 3)
func cloudflareStep3Page(c *via.Context, cfg *env.EnvConfig, mockMode bool) {
	svc := env.NewService(mockMode)

	// All fields for final save (token fields already set in previous steps)
	fields := CreateFormFields(c, cfg, []string{
		env.KeyCloudflareAPIToken,
		env.KeyCloudflareAPITokenName,
		env.KeyCloudflareAccountID,
		env.KeyCloudflarePageProject,
	})

	saveMessage := c.Signal("")

	// Finish action - save everything
	finishAction := c.Action(func() {
		saveMessage.SetValue("")

		fieldUpdates := map[string]string{
			env.KeyCloudflareAPIToken:     fields[0].ValueSignal.String(),
			env.KeyCloudflareAPITokenName: fields[1].ValueSignal.String(),
			env.KeyCloudflareAccountID:    fields[2].ValueSignal.String(),
			env.KeyCloudflarePageProject:  fields[3].ValueSignal.String(),
		}

		results, err := svc.ValidateAndUpdateFields(fieldUpdates)
		UpdateValidationStatus(results, fields, c)

		if err != nil {
			saveMessage.SetValue("error:" + err.Error())
			c.Sync()
			return
		}

		// Check only the fields we're updating, not all fields
		hasErrors := false
		for key := range fieldUpdates {
			if result, exists := results[key]; exists {
				if !result.Skipped && !result.Valid {
					hasErrors = true
					break
				}
			}
		}

		if hasErrors {
			saveMessage.SetValue("error:Please fix validation errors before saving")
			c.Sync()
			return
		}

		// Success!
		saveMessage.SetValue("success:✅ Configuration saved successfully!")
		c.Sync()
	})

	c.View(func() h.H {
		return h.Main(
			h.Class("container"),
			h.H1(h.Text("Cloudflare Setup - Step 3 of 3")),
			h.P(h.Text("Project Name (Optional)")),

			RenderNavigation("cloudflare"),

			h.H2(h.Text("Cloudflare Pages Project")),
			h.P(h.Text("Enter your Cloudflare Pages project name, or leave blank to create it later.")),

			h.Ul(
				h.Li(h.Text("If you already have a project, enter its name")),
				h.Li(h.Text("Must be lowercase letters, numbers, and hyphens only (1-63 chars)")),
				h.Li(h.Text("Examples: 'ubuntusoftware-net' or 'my-hugo-site'")),
			),

			h.P(h.Text("To create a new project: "), h.A(h.Href(env.CloudflarePagesURL), h.Attr("target", "_blank"), h.Text("Workers & Pages ↗")),
				h.Text(" → Create application → Pages → Connect to Git")),

			h.H3(h.Text("Project Name (optional):")),
			RenderFormField(fields[3]),

			h.Div(
				h.Style("margin-top: 2rem;"),
				h.A(h.Href("/cloudflare/step2"), h.Text("← Back: Account ID")),
				h.Text(" "),
				h.Button(h.Text("Finish & Save"), finishAction.OnClick()),
			),

			RenderSaveMessage(saveMessage)[0],
			RenderSaveMessage(saveMessage)[1],
		)
	})
}
