package web

import (
	"github.com/go-via/via"
	"github.com/go-via/via/h"
	"github.com/joeblew999/ubuntu-website/internal/env"
)

// homePage creates a welcome page with optional validation
func homePage(c *via.Context, cfg *env.EnvConfig, mockMode bool) {
	// Reactive signals for validation state
	isValidated := c.Signal(false) // Track if we've run validation
	validationInProgress := c.Signal(false)
	validationMessage := c.Signal("")

	// Validate action - trigger validation on demand
	validateAction := c.Action(func() {
		validationInProgress.SetValue(true)
		validationMessage.SetValue("Validating credentials...")
		c.Sync()

		// Just set the flag - the view will regenerate with validation
		validationInProgress.SetValue(false)
		isValidated.SetValue(true)
		validationMessage.SetValue("Validation complete!")
		c.Sync()
	})

	c.View(func() h.H {
		// Build table data based on validation state
		// This runs every time the view is rendered
		skipValidation := isValidated.String() != "true"
		tableRows, envPath, err := BuildConfigTableRows(mockMode, skipValidation)

		var configTableElement h.H
		if err != nil {
			configTableElement = h.Article(
				h.Style("background-color: var(--pico-del-background); border-left: 4px solid var(--pico-del-color); padding: 1rem; margin-bottom: 1rem;"),
				h.P(
					h.Style("margin: 0; color: var(--pico-del-color);"),
					h.Text("Error loading configuration: "+err.Error()),
				),
			)
		} else {
			configTableElement = renderConfigTable(tableRows, envPath)
		}

		return h.Main(
			h.Class("container"),
			h.H1(h.Text("Environment Setup")),
			h.P(h.Text("Configure your Cloudflare and Claude credentials for deployment and translation")),

			// Navigation
			RenderNavigation("home"),

			// Configuration Overview Table
			h.H2(h.Text("Configuration Overview")),

			// Validate button and status message
			h.Div(
				h.Style("margin-bottom: 1rem; display: flex; align-items: center; gap: 1rem;"),
				h.Button(
					h.Text("Validate Credentials"),
					h.If(validationInProgress.String() == "true", h.Attr("aria-busy", "true")),
					h.If(validationInProgress.String() == "true", h.Attr("disabled", "disabled")),
					h.If(isValidated.String() == "true", h.Attr("disabled", "disabled")),
					validateAction.OnClick(),
				),
				h.If(validationMessage.String() != "",
					h.Span(
						h.Style("color: var(--pico-ins-color);"),
						h.Text(validationMessage.String()),
					),
				),
			),

			// Render table
			configTableElement,

			// Status overview (non-interactive)
			h.H2(h.Text("Quick Setup")),
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

				// Deploy Card
				h.Article(
					h.A(
						h.Href("/deploy"),
						h.Style("text-decoration: none; color: inherit;"),
						h.H3(h.Text("Deploy Site")),
						h.P(h.Text("Build and deploy your Hugo site to Cloudflare Pages")),
					),
				),
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

// renderConfigTable renders the configuration overview table
func renderConfigTable(rows []ConfigTableRow, envPath string) h.H {
	// Build table header
	tableHeader := h.THead(
		h.Tr(
			h.Th(h.Text("Display")),
			h.Th(h.Text("Key")),
			h.Th(h.Text("Value")),
			h.Th(h.Text("Required")),
			h.Th(h.Text("Validated")),
			h.Th(h.Text("Error")),
		),
	)

	// Build table body rows
	tableBodyRows := make([]h.H, len(rows))
	for i, row := range rows {
		// Color code the validation status
		validatedStyle := ""
		if row.Validated == "✓" {
			validatedStyle = "color: var(--pico-ins-color);"
		} else if row.Validated == "✗" {
			validatedStyle = "color: var(--pico-del-color);"
		}

		// Color code the error column
		errorStyle := ""
		if row.Error != "-" {
			errorStyle = "color: var(--pico-del-color); font-size: 0.875rem;"
		}

		tableBodyRows[i] = h.Tr(
			h.Td(h.Text(row.Display)),
			h.Td(h.Code(h.Text(row.Key))),                                // Monospace for env var names
			h.Td(h.Code(h.Text(row.Value))),                              // Monospace for values
			h.Td(h.Text(row.Required)),
			h.Td(h.Style(validatedStyle), h.Text(row.Validated)),
			h.Td(h.Style(errorStyle), h.Text(row.Error)),
		)
	}

	tableBody := h.TBody(tableBodyRows...)

	return h.Div(
		h.P(
			h.Style("margin-bottom: 1rem; color: var(--pico-muted-color);"),
			h.Text(envPath),
		),
		h.Figure(
			h.Table(
				tableHeader,
				tableBody,
			),
		),
	)
}
