package web

import (
	"strings"

	"github.com/go-via/via"
	"github.com/go-via/via/h"
	"github.com/joeblew999/ubuntu-website/internal/env"
)

// deployPage - Build and deploy Hugo site to Cloudflare Pages
func deployPage(c *via.Context, cfg *env.EnvConfig, mockMode bool) {
	// Get project name from config
	projectName := cfg.Get(env.KeyCloudflarePageProject)

	// Deployment signals
	deployOutput := c.Signal("")
	deployInProgress := c.Signal(false)
	localURL := c.Signal("")
	deploymentURL := c.Signal("")

	// Build & Deploy action
	buildDeployAction := c.Action(func() {
		currentProject := cfg.Get(env.KeyCloudflarePageProject)
		if currentProject == "" || env.IsPlaceholder(currentProject) {
			deployOutput.SetValue("error:No project configured. Please complete Step 4 of the Cloudflare setup first.")
			c.Sync()
			return
		}

		deployInProgress.SetValue(true)
		deployOutput.SetValue("Starting build and deployment...\n")
		c.Sync()

		// Run build and deploy
		result := env.BuildAndDeploy(currentProject, mockMode)

		deployInProgress.SetValue(false)
		if result.Error != nil {
			deployOutput.SetValue("error:" + result.Output + "\nError: " + result.Error.Error())
		} else {
			deployOutput.SetValue("success:" + result.Output)
		}
		localURL.SetValue(result.LocalURL)
		deploymentURL.SetValue(result.DeploymentURL)
		c.Sync()
	})

	// Build only action
	buildOnlyAction := c.Action(func() {
		deployInProgress.SetValue(true)
		deployOutput.SetValue("Building Hugo site...\n")
		c.Sync()

		result := env.BuildHugoSite(mockMode)

		deployInProgress.SetValue(false)
		if result.Error != nil {
			deployOutput.SetValue("error:" + result.Output + "\nError: " + result.Error.Error())
		} else {
			deployOutput.SetValue("success:" + result.Output)
		}
		localURL.SetValue(result.LocalURL)
		c.Sync()
	})

	c.View(func() h.H {
		return h.Main(
			h.Class("container"),
			h.H1(h.Text("Deploy to Cloudflare Pages")),

			RenderNavigation("deploy"),

			// Project info section
			h.Article(
				h.Style("background-color: var(--pico-card-background-color); padding: 1rem; margin-bottom: 2rem;"),
				h.H3(h.Text("Current Configuration")),
				h.P(
					h.Strong(h.Text("Project: ")),
					h.If(projectName != "" && !env.IsPlaceholder(projectName),
						h.Text(projectName),
					),
					h.If(projectName == "" || env.IsPlaceholder(projectName),
						h.Span(
							h.Style("color: var(--pico-del-color);"),
							h.Text("Not configured"),
						),
					),
				),
				h.If(projectName == "" || env.IsPlaceholder(projectName),
					h.P(
						h.Style("margin-top: 1rem;"),
						h.Text("⚠️ Please "),
						h.A(h.Href("/cloudflare/step4"), h.Text("complete Cloudflare setup")),
						h.Text(" and select a project before deploying."),
					),
				),
			),

			// Deployment section
			h.Div(
				h.H2(h.Text("Build & Deploy")),
				h.P(h.Text("Build your Hugo site and deploy it to Cloudflare Pages.")),

				// Build and deploy buttons
				h.Div(
					h.Style("display: flex; gap: 1rem; margin-bottom: 1rem;"),
					h.Button(
						h.Attr("class", "secondary"),
						h.Text("Build Site Only"),
						h.If(deployInProgress.String() == "true", h.Attr("aria-busy", "true")),
						h.If(deployInProgress.String() == "true", h.Attr("disabled", "disabled")),
						h.If(projectName == "" || env.IsPlaceholder(projectName), h.Attr("disabled", "disabled")),
						buildOnlyAction.OnClick(),
					),
					h.Button(
						h.Text("Build & Deploy"),
						h.If(deployInProgress.String() == "true", h.Attr("aria-busy", "true")),
						h.If(deployInProgress.String() == "true", h.Attr("disabled", "disabled")),
						h.If(projectName == "" || env.IsPlaceholder(projectName), h.Attr("disabled", "disabled")),
						buildDeployAction.OnClick(),
					),
				),

				// Output display
				h.If(deployOutput.String() != "",
					h.Div(
						h.Style("margin-top: 1.5rem;"),
						// Success output
						h.If(strings.HasPrefix(deployOutput.String(), "success:"),
							h.Article(
								h.Style("background-color: var(--pico-ins-background); border-left: 4px solid var(--pico-ins-color); padding: 1rem;"),
								h.Pre(
									h.Style("margin: 0; white-space: pre-wrap; font-size: 0.875rem; color: var(--pico-ins-color);"),
									h.Text(strings.TrimPrefix(deployOutput.String(), "success:")),
								),
							),
						),
						// Error output
						h.If(strings.HasPrefix(deployOutput.String(), "error:"),
							h.Article(
								h.Style("background-color: var(--pico-del-background); border-left: 4px solid var(--pico-del-color); padding: 1rem;"),
								h.Pre(
									h.Style("margin: 0; white-space: pre-wrap; font-size: 0.875rem; color: var(--pico-del-color);"),
									h.Text(strings.TrimPrefix(deployOutput.String(), "error:")),
								),
							),
						),
						// In-progress output
						h.If(!strings.HasPrefix(deployOutput.String(), "success:") && !strings.HasPrefix(deployOutput.String(), "error:"),
							h.Article(
								h.Style("background-color: var(--pico-card-background-color); border-left: 4px solid var(--pico-primary); padding: 1rem;"),
								h.Pre(
									h.Style("margin: 0; white-space: pre-wrap; font-size: 0.875rem;"),
									h.Text(deployOutput.String()),
								),
							),
						),

						// Preview URLs section
						h.If(localURL.String() != "" || deploymentURL.String() != "",
							h.Div(
								h.Style("margin-top: 1.5rem;"),
								h.H3(h.Text("Preview URLs")),

								// Local preview URL
								h.If(localURL.String() != "",
									RenderURLLink(localURL.String(), "Local Preview", "🌐"),
								),

								// Deployment URL
								h.If(deploymentURL.String() != "",
									RenderURLLink(deploymentURL.String(), "Live Deployment", "☁️"),
								),
							),
						),
					),
				),
			),

			h.Div(
				h.Style("margin-top: 2rem;"),
				h.A(h.Href("/"), h.Text("← Back to Overview")),
			),
		)
	})
}
