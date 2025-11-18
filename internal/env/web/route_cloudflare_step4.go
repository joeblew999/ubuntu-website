package web

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-via/via"
	"github.com/go-via/via/h"
	"github.com/joeblew999/ubuntu-website/internal/env"
)

// cloudflareStep4Page - Project details (Step 4 of 4)
func cloudflareStep4Page(c *via.Context, cfg *env.EnvConfig, mockMode bool) {
	svc := env.NewService(mockMode)

	// All fields for final save (token fields already set in previous steps)
	fields := CreateFormFields(c, cfg, []string{
		env.KeyCloudflareAPIToken,
		env.KeyCloudflareAPITokenName,
		env.KeyCloudflareAccountID,
		env.KeyCloudflareDomain,
		env.KeyCloudflareZoneID,
		env.KeyCloudflarePageProject,
	})

	saveMessage := c.Signal("")
	projectsMessage := c.Signal("") // For projects loading status
	deleteMessage := c.Signal("")   // For delete operation feedback
	projectToDelete := c.Signal("") // Holds project name pending deletion confirmation
	showDeleteConfirm := c.Signal(false) // Controls visibility of delete confirmation dialog

	// Load projects from Cloudflare API
	// Read directly from config (not from form signals which may be cleared by placeholder detection)
	token := cfg.Get(env.KeyCloudflareAPIToken)
	accountID := cfg.Get(env.KeyCloudflareAccountID)

	var projects []env.PagesProject
	var projectsErr error

	if !mockMode && token != "" && accountID != "" && !env.IsPlaceholder(token) && !env.IsPlaceholder(accountID) {
		projects, projectsErr = env.ListPagesProjects(token, accountID)
		if projectsErr != nil {
			log.Printf("Failed to fetch Pages projects: %v", projectsErr)
			projectsMessage.SetValue("error:Failed to load projects: " + projectsErr.Error())
		} else if len(projects) == 0 {
			projectsMessage.SetValue("info:No projects found in this account")
		}
	} else if mockMode {
		// Mock data for testing
		projects = []env.PagesProject{
			{Name: "my-hugo-site", CreatedOn: "2024-01-15T10:00:00Z"},
			{Name: "ubuntusoftware-net", CreatedOn: "2024-02-20T14:30:00Z"},
			{Name: "test-project", CreatedOn: "2024-03-10T09:15:00Z"},
		}
	}

	// Build smart Pages URL with account ID if available
	pagesURL := BuildCloudflareURL(env.CloudflarePagesURL, accountID)

	// Build dropdown options from projects
	projectOptions := make([]SelectOption, 0, len(projects)+1)
	projectOptions = append(projectOptions, SelectOption{Value: "", Label: "-- Select a project --"})
	for _, project := range projects {
		projectOptions = append(projectOptions, SelectOption{Value: project.Name, Label: project.Name})
	}

	// Build project list UI elements with delete buttons
	projectListElements := make([]h.H, 0, len(projects))
	for _, project := range projects {
		projectName := project.Name       // Capture in closure
		createdOn := project.CreatedOn    // Capture in closure
		deleteAction := c.Action(func() {
			projectToDelete.SetValue(projectName)
			showDeleteConfirm.SetValue(true)
			deleteMessage.SetValue("")
			c.Sync()
		})

		projectListElements = append(projectListElements, h.Div(
			h.Style("display: flex; justify-content: space-between; align-items: center; padding: 0.75rem; background: var(--pico-card-background-color); border-radius: 0.25rem;"),
			h.Div(
				h.Strong(h.Text(projectName)),
				h.Small(
					h.Style("margin-left: 1rem; color: var(--pico-muted-color);"),
					h.Text("Created: "+createdOn),
				),
			),
			h.Button(
				h.Attr("class", "secondary outline"),
				h.Text("Delete"),
				deleteAction.OnClick(),
			),
		))
	}

	// Cancel delete operation
	cancelDeleteAction := c.Action(func() {
		projectToDelete.SetValue("")
		showDeleteConfirm.SetValue(false)
		deleteMessage.SetValue("")
		c.Sync()
	})

	// Confirm and execute delete
	confirmDeleteAction := c.Action(func() {
		projectName := projectToDelete.String()
		if projectName == "" {
			deleteMessage.SetValue("error:No project selected for deletion")
			c.Sync()
			return
		}

		// Get credentials from config
		token := cfg.Get(env.KeyCloudflareAPIToken)
		accountID := cfg.Get(env.KeyCloudflareAccountID)

		if mockMode {
			// Mock mode - simulate successful deletion
			deleteMessage.SetValue("success:Project '" + projectName + "' deleted successfully (mock mode)")
			showDeleteConfirm.SetValue(false)
			projectToDelete.SetValue("")
			c.Sync()
			// In real app, would reload page to refresh project list
			c.ExecScript("setTimeout(() => window.location.reload(), 1500)")
			return
		}

		// Call delete API with automatic custom domain cleanup
		removedDomains, err := env.DeletePagesProjectWithCleanup(token, accountID, projectName)
		if err != nil {
			deleteMessage.SetValue("error:Failed to delete project: " + err.Error())
			c.Sync()
			return
		}

		// Success message includes info about removed domains
		successMsg := "Project '" + projectName + "' deleted successfully!"
		if len(removedDomains) > 0 {
			successMsg += " (Removed " + fmt.Sprintf("%d", len(removedDomains)) + " custom domain(s) first)"
		}
		successMsg += " Refreshing..."

		deleteMessage.SetValue("success:" + successMsg)
		showDeleteConfirm.SetValue(false)
		projectToDelete.SetValue("")
		c.Sync()

		// Reload page to refresh project list
		c.ExecScript("setTimeout(() => window.location.reload(), 1500)")
	})

	// Deployment signals
	deployOutput := c.Signal("")
	deployInProgress := c.Signal(false)
	newProjectName := c.Signal("")

	// Build & Deploy action
	buildDeployAction := c.Action(func() {
		projectName := fields[5].ValueSignal.String()
		if projectName == "" || env.IsPlaceholder(projectName) {
			deployOutput.SetValue("error:Please select or create a project first")
			c.Sync()
			return
		}

		deployInProgress.SetValue(true)
		deployOutput.SetValue("Starting build and deployment...\n")
		c.Sync()

		// Run build and deploy
		result := env.BuildAndDeploy(projectName, mockMode)

		deployInProgress.SetValue(false)
		if result.Error != nil {
			deployOutput.SetValue("error:" + result.Output + "\nError: " + result.Error.Error())
		} else {
			deployOutput.SetValue("success:" + result.Output)
		}
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
		c.Sync()
	})

	// Create project action
	createProjectAction := c.Action(func() {
		projectName := newProjectName.String()
		if projectName == "" {
			deployOutput.SetValue("error:Please enter a project name")
			c.Sync()
			return
		}

		deployInProgress.SetValue(true)
		deployOutput.SetValue("Creating project '" + projectName + "'...\n")
		c.Sync()

		result := env.CreatePagesProject(projectName, mockMode)

		deployInProgress.SetValue(false)
		if result.Error != nil {
			deployOutput.SetValue("error:" + result.Output + "\nError: " + result.Error.Error())
		} else {
			deployOutput.SetValue("success:" + result.Output + "\n\nRefreshing page to show new project...")
			c.Sync()
			// Reload page to refresh project list
			c.ExecScript("setTimeout(() => window.location.reload(), 2000)")
			return
		}
		c.Sync()
	})

	// Finish action - save everything
	finishAction := c.Action(func() {
		saveMessage.SetValue("")

		fieldUpdates := map[string]string{
			env.KeyCloudflareAPIToken:     fields[0].ValueSignal.String(),
			env.KeyCloudflareAPITokenName: fields[1].ValueSignal.String(),
			env.KeyCloudflareAccountID:    fields[2].ValueSignal.String(),
			env.KeyCloudflareDomain:       fields[3].ValueSignal.String(),
			env.KeyCloudflareZoneID:       fields[4].ValueSignal.String(),
			env.KeyCloudflarePageProject:  fields[5].ValueSignal.String(),
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
			h.H1(h.Text("Cloudflare Setup - Step 4 of 4")),
			h.P(h.Text("Project Name (Optional)")),

			RenderNavigation("cloudflare"),

			h.H2(h.Text("Cloudflare Pages Project")),
			h.P(h.Text("Select a Pages project to deploy your Hugo site to.")),

			// Show projects loading status - info message
			h.If(projectsMessage.String() == "info:No projects found in this account",
				h.Article(
					h.Style("background-color: var(--pico-ins-background); border-left: 4px solid var(--pico-ins-color); padding: 1rem; margin-bottom: 1rem;"),
					h.P(
						h.Style("margin: 0;"),
						h.Text("No projects found. You can leave this blank and create one later via "),
						h.A(h.Href(pagesURL), h.Attr("target", "_blank"), h.Text("Workers & Pages ↗")),
						h.Text("."),
					),
				),
			),
			// Show projects loading status - error message
			h.If(strings.HasPrefix(projectsMessage.String(), "error:"),
				h.Article(
					h.Style("background-color: var(--pico-del-background); border-left: 4px solid var(--pico-del-color); padding: 1rem; margin-bottom: 1rem;"),
					h.P(
						h.Style("margin: 0; color: var(--pico-del-color);"),
						h.Text(strings.TrimPrefix(projectsMessage.String(), "error:")),
					),
				),
			),

			// Project dropdown (always shown, even if empty - will have placeholder option)
			h.H3(h.Text("Choose Project:")),
			RenderSelectField("Project", fields[5].ValueSignal, projectOptions),
			h.Small(
				h.Style("color: var(--pico-muted-color);"),
				h.Text("Select a Pages project to deploy your Hugo site to"),
			),

			// Show creation instructions if user wants to create a new project
			h.Details(
				h.Style("margin-top: 1.5rem;"),
				h.Summary(h.Text("How to create a new Cloudflare Pages project")),
				h.P(h.Text("To create a new project:")),
				h.Ol(
					h.Li(h.Text("Visit "), h.A(h.Href(pagesURL), h.Attr("target", "_blank"), h.Text("Workers & Pages ↗"))),
					h.Li(h.Text("Click 'Create application' → 'Pages' → 'Connect to Git'")),
					h.Li(h.Text("Follow the setup wizard to connect your repository")),
					h.Li(h.Text("Once created, return here and refresh to see it in the dropdown")),
				),
				h.P(
					h.Style("margin-top: 1rem;"),
					h.Strong(h.Text("Project naming rules: ")),
					h.Text("Lowercase letters, numbers, and hyphens only (1-63 chars). Examples: 'ubuntusoftware-net' or 'my-hugo-site'"),
				),
			),

			// Manage Projects section - delete existing projects
			h.If(len(projects) > 0,
				h.Div(
					h.Style("margin-top: 3rem; padding-top: 2rem; border-top: 1px solid var(--pico-muted-border-color);"),
					h.H2(h.Text("Manage Existing Projects")),
					h.P(h.Text("Delete projects you no longer need:")),

					// Delete confirmation dialog
					h.If(showDeleteConfirm.String() == "true",
						h.Dialog(
							h.Attr("open", "open"),
							h.Article(
								h.H3(h.Text("Confirm Deletion")),
								h.P(
									h.Text("Are you sure you want to delete the project "),
									h.Strong(h.Text(projectToDelete.String())),
									h.Text("?"),
								),
								h.P(
									h.Style("color: var(--pico-del-color);"),
									h.Text("⚠️ This action cannot be undone. All deployments and settings for this project will be permanently deleted."),
								),
								h.Div(
									h.Style("display: flex; gap: 1rem; justify-content: flex-end;"),
									h.Button(
										h.Attr("class", "secondary"),
										h.Text("Cancel"),
										cancelDeleteAction.OnClick(),
									),
									h.Button(
										h.Attr("class", "contrast"),
										h.Text("Delete Project"),
										confirmDeleteAction.OnClick(),
									),
								),
							),
						),
					),

					// List of projects with delete buttons
					func() h.H {
						listChildren := []h.H{h.Style("display: grid; gap: 0.5rem;")}
						listChildren = append(listChildren, projectListElements...)
						return h.Div(listChildren...)
					}(),

					// Show delete messages
					h.If(strings.HasPrefix(deleteMessage.String(), "success:"),
						h.Article(
							h.Style("background-color: var(--pico-ins-background); border-left: 4px solid var(--pico-ins-color); padding: 1rem; margin-top: 1rem;"),
							h.P(
								h.Style("margin: 0; color: var(--pico-ins-color);"),
								h.Text(strings.TrimPrefix(deleteMessage.String(), "success:")),
							),
						),
					),
					h.If(strings.HasPrefix(deleteMessage.String(), "error:"),
						h.Article(
							h.Style("background-color: var(--pico-del-background); border-left: 4px solid var(--pico-del-color); padding: 1rem; margin-top: 1rem;"),
							h.P(
								h.Style("margin: 0; color: var(--pico-del-color);"),
								h.Text(strings.TrimPrefix(deleteMessage.String(), "error:")),
							),
						),
					),
				),
			),

			// Deployment section
			h.Div(
				h.Style("margin-top: 3rem; padding-top: 2rem; border-top: 1px solid var(--pico-muted-border-color);"),
				h.H2(h.Text("Build & Deploy")),
				h.P(h.Text("Build your Hugo site and deploy it to Cloudflare Pages.")),

				// Create new project section
				h.Details(
					h.Style("margin-bottom: 2rem;"),
					h.Summary(h.Text("Create New Project via Wrangler")),
					h.P(h.Text("Enter a project name and create it directly using Wrangler CLI.")),
					h.Div(
						h.Style("display: flex; gap: 1rem; align-items: flex-end;"),
						h.Div(
							h.Style("flex: 1;"),
							h.Label(h.Text("Project Name")),
							h.Input(
								h.Attr("type", "text"),
								h.Attr("placeholder", "my-hugo-site"),
								newProjectName.Bind(),
							),
							h.Small(
								h.Style("color: var(--pico-muted-color);"),
								h.Text("Lowercase letters, numbers, and hyphens only (1-63 chars)"),
							),
						),
						h.Button(
							h.Text("Create Project"),
							h.If(deployInProgress.String() == "true", h.Attr("aria-busy", "true")),
							h.If(deployInProgress.String() == "true", h.Attr("disabled", "disabled")),
							createProjectAction.OnClick(),
						),
					),
				),

				// Build and deploy buttons
				h.Div(
					h.Style("display: flex; gap: 1rem; margin-bottom: 1rem;"),
					h.Button(
						h.Attr("class", "secondary"),
						h.Text("Build Site Only"),
						h.If(deployInProgress.String() == "true", h.Attr("aria-busy", "true")),
						h.If(deployInProgress.String() == "true", h.Attr("disabled", "disabled")),
						buildOnlyAction.OnClick(),
					),
					h.Button(
						h.Text("Build & Deploy"),
						h.If(deployInProgress.String() == "true", h.Attr("aria-busy", "true")),
						h.If(deployInProgress.String() == "true", h.Attr("disabled", "disabled")),
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
					),
				),
			),

			h.Div(
				h.Style("margin-top: 2rem;"),
				h.A(h.Href("/cloudflare/step3"), h.Text("← Back: Domain Selection")),
				h.Text(" "),
				h.Button(h.Text("Finish & Save"), finishAction.OnClick()),
			),

			RenderSaveMessage(saveMessage)[0],
			RenderSaveMessage(saveMessage)[1],
		)
	})
}
