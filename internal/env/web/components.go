package web

import (
	"github.com/go-via/via/h"
	"github.com/joeblew999/ubuntu-website/internal/env"
)

// renderNavigation renders the shared navigation menu
// currentPage: "home", "cloudflare", or "claude"
func renderNavigation(currentPage string) h.H {
	// Helper to render a nav item (link or bold text)
	navItem := func(page, label, href string) h.H {
		if currentPage == page {
			return h.Li(h.Strong(h.Text(label)))
		}
		return h.Li(h.A(h.Href(href), h.Text(label)))
	}

	return h.Nav(
		h.Ul(
			navItem("home", "Overview", "/"),
			navItem("cloudflare", "Cloudflare", "/cloudflare"),
			navItem("claude", "Claude AI", "/claude"),
		),
	)
}

// updateValidationError updates error signal based on ValidationResult
// Takes a pointer to an error signal (returned from c.Signal())
func updateValidationError(result env.ValidationResult, errorSignal interface {
	SetValue(any)
}) {
	if result.Skipped || result.Valid {
		errorSignal.SetValue("")
	} else if result.Error != nil {
		errorSignal.SetValue(result.Error.Error())
	}
}
