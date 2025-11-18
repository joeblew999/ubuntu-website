package web

import (
	"log"
	"strings"

	"github.com/go-via/via"
	"github.com/go-via/via/h"
	"github.com/joeblew999/ubuntu-website/internal/env"
)

// cloudflareStep3Page - Domain selection (Step 3 of 4)
func cloudflareStep3Page(c *via.Context, cfg *env.EnvConfig, mockMode bool) {
	svc := env.NewService(mockMode)

	// Fields for all previously entered data plus domain/zone
	fields := CreateFormFields(c, cfg, []string{
		env.KeyCloudflareAPIToken,
		env.KeyCloudflareAccountID,
		env.KeyCloudflareDomain,
		env.KeyCloudflareZoneID,
	})

	// Pre-populate domain dropdown with saved domain+zone if both exist
	// fields[2] is the domain field - we need to set it to "domain|zoneID" format for the dropdown
	savedDomain := cfg.Get(env.KeyCloudflareDomain)
	savedZoneID := cfg.Get(env.KeyCloudflareZoneID)
	if savedDomain != "" && savedZoneID != "" && !env.IsPlaceholder(savedDomain) && !env.IsPlaceholder(savedZoneID) {
		// Replace fields[2] with the combined value for dropdown binding
		fields[2] = FormFieldData{
			EnvKey:       env.KeyCloudflareDomain,
			ValueSignal:  c.Signal(savedDomain + "|" + savedZoneID),
			StatusSignal: fields[2].StatusSignal,
		}
	}

	saveMessage := c.Signal("")
	zonesMessage := c.Signal("") // For zones loading status

	// Zones loader - populated lazily when first accessed
	zonesLoader := NewLazyLoader(func() ([]env.Zone, error) {
		token := cfg.Get(env.KeyCloudflareAPIToken)
		accountID := cfg.Get(env.KeyCloudflareAccountID)

		if mockMode {
			// Mock data for testing
			return []env.Zone{
				{ID: "mock-zone-1", Name: "example.com"},
				{ID: "mock-zone-2", Name: "example.net"},
				{ID: "mock-zone-3", Name: "example.org"},
				{ID: "mock-zone-4", Name: "ubuntusoftware.net"},
				{ID: "mock-zone-5", Name: "mysite.com"},
				{ID: "mock-zone-6", Name: "testdomain.io"},
			}, nil
		}

		if token == "" || accountID == "" || env.IsPlaceholder(token) || env.IsPlaceholder(accountID) {
			return []env.Zone{}, nil
		}

		return env.ListZones(token, accountID)
	})

	// Next action - save domain selection and go to step 4
	nextAction := c.Action(func() {
		saveMessage.SetValue("")

		// Parse selected domain value (format: "domain.com|zone-id")
		selectedValue := fields[2].ValueSignal.String()
		if selectedValue == "" {
			saveMessage.SetValue("error:Please select a domain")
			c.Sync()
			return
		}

		// Split domain|zone-id
		parts := strings.Split(selectedValue, "|")
		if len(parts) != 2 {
			saveMessage.SetValue("error:Invalid domain selection")
			c.Sync()
			return
		}

		domain := parts[0]
		zoneID := parts[1]

		fieldUpdates := map[string]string{
			env.KeyCloudflareAPIToken:  fields[0].ValueSignal.String(),
			env.KeyCloudflareAccountID: fields[1].ValueSignal.String(),
			env.KeyCloudflareDomain:    domain,
			env.KeyCloudflareZoneID:    zoneID,
		}

		results, err := svc.ValidateAndUpdateFields(fieldUpdates)

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

		// Success - redirect to step 4
		saveMessage.SetValue("success:Domain selected! Moving to step 4...")
		c.Sync()
		c.ExecScript("window.location.href = '/cloudflare/step4'")
	})

	c.View(func() h.H {
		// Load zones using LazyLoader
		zonesCache, zonesErr := zonesLoader.Get()
		if zonesErr != nil {
			log.Printf("Failed to fetch zones: %v", zonesErr)
			zonesMessage.SetValue("error:Failed to load domains: " + zonesErr.Error())
		} else if len(zonesCache) == 0 {
			zonesMessage.SetValue("info:No domains found in this account")
		}

		// Build dropdown options from zones
		domainOptions := make([]SelectOption, 0, len(zonesCache)+1)
		domainOptions = append(domainOptions, SelectOption{Value: "", Label: "-- Select a domain --"})
		for _, zone := range zonesCache {
			domainOptions = append(domainOptions, SelectOption{Value: zone.Name + "|" + zone.ID, Label: zone.Name})
		}

		// Build smart "Add Site" URL with account ID if available
		accountID := cfg.Get(env.KeyCloudflareAccountID)
		addSiteURL := BuildCloudflareURL(env.CloudflareAddSiteURL, accountID)

		return h.Main(
			h.Class("container"),
			h.H1(h.Text("Cloudflare Setup - Step 3 of 4")),
			h.P(h.Text("Domain Selection")),

			RenderNavigation("cloudflare"),

			h.H2(h.Text("Select Your Domain")),
			h.P(h.Text("Choose which domain you want to deploy your Hugo site to.")),

			// Show zones loading status - info message
			h.If(zonesMessage.String() == "info:No domains found in this account",
				h.Article(
					h.Style("background-color: var(--pico-ins-background); border-left: 4px solid var(--pico-ins-color); padding: 1rem; margin-bottom: 1rem;"),
					h.P(
						h.Style("margin: 0;"),
						h.Text("No domains found in this account. "),
						h.A(h.Href(addSiteURL), h.Attr("target", "_blank"), h.Text("Add a domain ↗")),
						h.Text(" to Cloudflare first."),
					),
				),
			),
			// Show zones loading status - error message
			h.If(strings.HasPrefix(zonesMessage.String(), "error:"),
				h.Article(
					h.Style("background-color: var(--pico-del-background); border-left: 4px solid var(--pico-del-color); padding: 1rem; margin-bottom: 1rem;"),
					h.P(
						h.Style("margin: 0; color: var(--pico-del-color);"),
						h.Text(strings.TrimPrefix(zonesMessage.String(), "error:")),
					),
				),
			),

			// Domain dropdown
			h.If(len(domainOptions) > 1,
				h.Div(
					h.H3(h.Text("Choose Domain:")),
					RenderSelectField("Domain", fields[2].ValueSignal, domainOptions),
					h.Small(
						h.Style("color: var(--pico-muted-color);"),
						h.Text("The domain where your Hugo site will be deployed"),
					),
				),
			),

			h.Div(
				h.Style("margin-top: 2rem;"),
				h.A(h.Href("/cloudflare/step2"), h.Text("← Back: Account ID")),
				h.Text(" "),
				h.If(len(domainOptions) > 1,
					h.Button(h.Text("Next: Project Details →"), nextAction.OnClick()),
				),
				h.If(len(domainOptions) > 1,
					h.Text(" or "),
				),
				h.A(h.Href("/cloudflare/step4"), h.Text("Skip")),
			),

			RenderErrorMessage(saveMessage),
			RenderSuccessMessage(saveMessage),
		)
	})
}
