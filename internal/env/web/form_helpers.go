package web

import (
	"strings"

	"github.com/go-via/via"
	"github.com/go-via/via/h"
	"github.com/joeblew999/ubuntu-website/internal/env"
)

// FormFieldData holds the data needed to render a form field
// We can't reference via.signal directly as it's unexported,
// so we define interfaces matching the actual signal methods
type FormFieldData struct {
	EnvKey       string
	ValueSignal  interface{ String() string; Bind() h.H }
	StatusSignal interface{ String() string; SetValue(any) }
}

// CreateFormFields creates signals for a set of form fields
func CreateFormFields(c *via.Context, cfg *env.EnvConfig, envKeys []string) []FormFieldData {
	fields := make([]FormFieldData, len(envKeys))
	for i, key := range envKeys {
		value := cfg.Get(key)
		// Clear placeholder values for cleaner UX - start with empty fields
		if env.IsPlaceholder(value) {
			value = ""
		}
		fields[i] = FormFieldData{
			EnvKey:       key,
			ValueSignal:  c.Signal(value),
			StatusSignal: c.Signal(""),
		}
	}
	return fields
}

// RenderFormField renders a single form field with label, input, and validation status
func RenderFormField(field FormFieldData) h.H {
	status := field.StatusSignal.String()
	isError := strings.HasPrefix(status, "error:")
	isValid := status == "valid"

	// Build input attributes
	inputAttrs := []h.H{
		h.Type("text"),
		h.Value(field.ValueSignal.String()),
		field.ValueSignal.Bind(),
	}

	// Add PicoCSS validation styling
	if isError {
		inputAttrs = append(inputAttrs, h.Attr("aria-invalid", "true")) // Red border
	} else if isValid {
		inputAttrs = append(inputAttrs, h.Attr("aria-invalid", "false")) // Green border
	}

	return h.Div(
		h.Label(h.Text(env.GetFieldLabel(field.EnvKey))),
		h.Input(inputAttrs...),
		// Error message as <small> helper text (PicoCSS styling)
		h.If(isError,
			h.Small(
				h.Style("color: var(--pico-del-color);"), // Use PicoCSS error color
				h.Text(strings.TrimPrefix(status, "error:")),
			),
		),
	)
}

// UpdateValidationStatus updates validation status signals from results
func UpdateValidationStatus(results map[string]env.ValidationResult, fields []FormFieldData, c *via.Context) {
	for i := range fields {
		result, ok := results[fields[i].EnvKey]
		if !ok {
			continue
		}

		if result.Skipped {
			fields[i].StatusSignal.SetValue("")
		} else if result.Valid {
			fields[i].StatusSignal.SetValue("valid")
		} else {
			// Set error message with "error:" prefix for conditional display
			errorMsg := "error:"
			if result.Error != nil {
				errorMsg += result.Error.Error()
			} else {
				errorMsg += "Invalid value"
			}
			fields[i].StatusSignal.SetValue(errorMsg)
		}
	}
	// Use Sync() instead of SyncSignals() to re-render the view and show validation icons/messages
	c.Sync()
}

// CreateSaveAction creates a save action for form fields
func CreateSaveAction(c *via.Context, svc *env.Service, fields []FormFieldData, saveMessage interface{ String() string; SetValue(any) }) func() {
	return func() {
		// Prepare field updates map
		fieldUpdates := make(map[string]string)
		for _, field := range fields {
			fieldUpdates[field.EnvKey] = field.ValueSignal.String()
		}

		// Use service to validate and save atomically
		results, err := svc.ValidateAndUpdateFields(fieldUpdates)

		// Update validation status from results
		UpdateValidationStatus(results, fields, c)

		// Handle save result
		if err != nil {
			saveMessage.SetValue("error:" + err.Error())
		} else {
			// Check if there were validation errors
			if env.HasInvalidCredentialsMap(results) {
				saveMessage.SetValue("error:Please fix validation errors before saving")
			} else {
				saveMessage.SetValue("success:Configuration saved successfully!")
			}
		}

		// Note: UpdateValidationStatus already called c.Sync() above which re-renders
		// the entire view including the saveMessage, so no need to sync again here
	}
}

// RenderSaveMessage renders the save message with PicoCSS alert styling
func RenderSaveMessage(saveMessage interface{ String() string }) []h.H {
	return []h.H{
		h.If(strings.HasPrefix(saveMessage.String(), "error:"),
			h.Article(
				h.Style("background-color: var(--pico-del-background); border-left: 4px solid var(--pico-del-color); padding: 1rem; margin-top: 1rem;"),
				h.P(
					h.Style("margin: 0; color: var(--pico-del-color);"),
					h.Text(strings.TrimPrefix(saveMessage.String(), "error:")),
				),
			),
		),
		h.If(strings.HasPrefix(saveMessage.String(), "success:"),
			h.Article(
				h.Style("background-color: var(--pico-ins-background); border-left: 4px solid var(--pico-ins-color); padding: 1rem; margin-top: 1rem;"),
				h.P(
					h.Style("margin: 0; color: var(--pico-ins-color);"),
					h.Text(strings.TrimPrefix(saveMessage.String(), "success:")),
				),
			),
		),
	}
}

// RenderNavigation renders the navigation menu
func RenderNavigation(currentPage string) h.H {
	return h.Nav(
		h.Ul(
			h.Li(h.If(currentPage == "all",
				h.Strong(h.Text("All Settings")),
			), h.If(currentPage != "all",
				h.A(h.Href("/"), h.Text("All Settings")),
			)),
			h.Li(h.If(currentPage == "cloudflare",
				h.Strong(h.Text("Cloudflare Only")),
			), h.If(currentPage != "cloudflare",
				h.A(h.Href("/cloudflare"), h.Text("Cloudflare Only")),
			)),
			h.Li(h.If(currentPage == "claude",
				h.Strong(h.Text("Claude Only")),
			), h.If(currentPage != "claude",
				h.A(h.Href("/claude"), h.Text("Claude Only")),
			)),
		),
	)
}
