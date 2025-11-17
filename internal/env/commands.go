package env

import (
	"fmt"
)

// List displays unified view of local .env and GitHub secrets status
func List() error {
	// Get absolute path to .env file
	envPath, err := GetEnvPath()
	if err != nil {
		envPath = ".env"
	}

	// Use service to load config
	svc := NewService(false)
	cfg, err := svc.GetCurrentConfig()
	if err != nil {
		return err
	}

	// Try to get GitHub secrets (may fail if gh CLI not available)
	var githubSecrets []GitHubSecret
	var repoInfo string

	if err := ValidateGitHubSetup(); err == nil {
		owner, name, err := GetRepositoryInfo()
		if err == nil {
			repoInfo = fmt.Sprintf(" • GitHub: %s/%s", owner, name)
			secrets, err := ListGitHubSecrets()
			if err == nil {
				githubSecrets = secrets
			}
		}
	}

	printHeader("Environment Configuration", envPath+repoInfo)

	// Build unified table showing local + GitHub status
	rows := buildUnifiedRows(cfg, envFieldsInOrder, githubSecrets)
	renderCredentialTable(rows)

	// Show GitHub management link if available
	if len(repoInfo) > 0 {
		fmt.Println()
		repoURL, _ := GetRepositoryURL()
		fmt.Println("  → Manage GitHub secrets: " + repoURL + "/settings/secrets/actions")
	}

	printFooter("")

	return nil
}

// PushGithub syncs environment variables to GitHub secrets
func PushGithub(dryRun, force, validate bool) error {
	// Validate GitHub setup
	if err := ValidateGitHubSetup(); err != nil {
		return err
	}

	// Get repository info and print header
	repoName, err := printSyncHeader(dryRun, validate)
	if err != nil {
		return err
	}

	// Perform the sync
	opts := SyncOptions{
		DryRun:   dryRun,
		Force:    force,
		Validate: validate,
	}

	results, err := SyncSecretsToGitHub(opts)
	if err != nil {
		return err
	}

	// Display results and return error if any failed
	return displaySyncResults(results, dryRun, repoName)
}

// printSyncHeader prints the sync command header
func printSyncHeader(dryRun, validate bool) (string, error) {
	owner, name, err := GetRepositoryInfo()
	if err != nil {
		return "", err
	}

	repoName := fmt.Sprintf("%s/%s", owner, name)
	mode := "Push to GitHub"
	if dryRun {
		mode = "Push to GitHub (Dry Run)"
	}
	printHeader(mode, repoName)

	if dryRun {
		fmt.Println("DRY RUN MODE - No secrets will be modified")
		fmt.Println()
	}

	if validate {
		fmt.Println("Validating credentials before push...")
		fmt.Println()
	}

	return repoName, nil
}

// displaySyncResults displays sync results in a formatted table
func displaySyncResults(results []SyncResult, dryRun bool, repoName string) error {
	created, updated, skipped, failed := countSyncResults(results)

	// Find longest secret name for alignment
	maxNameLen := 0
	for _, result := range results {
		if len(result.Name) > maxNameLen {
			maxNameLen = len(result.Name)
		}
	}

	printSection("Secrets Status")

	// Display each result
	for _, result := range results {
		displaySyncResult(result, maxNameLen)
	}

	fmt.Println()

	// Build and print summary
	displaySyncSummary(created, updated, skipped, failed, dryRun)

	if failed > 0 {
		// Check if failures were due to validation and collect failed fields
		var failedFields []string
		for _, result := range results {
			if result.Status == "failed" && result.Reason == "validation failed" {
				failedFields = append(failedFields, result.Name)
			}
		}

		if len(failedFields) > 0 {
			fmt.Println()
			fmt.Println("⚠ Validation failed. Run 'local-setup' to fix credentials.")
			fmt.Printf("   Failed: %s\n", joinParts(failedFields))
			fmt.Println()
		}

		return fmt.Errorf("failed to sync %d secrets", failed)
	}

	return nil
}

// countSyncResults counts results by status
func countSyncResults(results []SyncResult) (created, updated, skipped, failed int) {
	for _, result := range results {
		switch result.Status {
		case SyncStatusSynced:
			if result.Reason == SyncReasonCreated {
				created++
			} else {
				updated++
			}
		case SyncStatusWouldSync:
			if result.Reason == SyncReasonWouldCreateNew {
				created++
			} else {
				updated++
			}
		case SyncStatusSkipped:
			skipped++
		case SyncStatusFailed:
			failed++
		}
	}
	return created, updated, skipped, failed
}

// displaySyncResult displays a single sync result line
func displaySyncResult(result SyncResult, maxNameLen int) {
	var icon, status string

	switch result.Status {
	case "synced":
		icon = "✓"
		if result.Reason == "created" {
			status = "Created new"
		} else {
			status = "Updated"
		}
	case "would-sync":
		icon = "→"
		if result.Reason == "would create new" {
			status = "Would create"
		} else {
			status = "Would update"
		}
	case "skipped":
		icon = "○"
		status = result.Reason
	case "failed":
		icon = "✗"
		status = fmt.Sprintf("%s: %v", result.Reason, result.Error)
	}

	nameDisplay := fmt.Sprintf("%-*s", maxNameLen, result.Name)
	fmt.Printf("  %s  %s  %s\n", icon, nameDisplay, status)
}

// displaySyncSummary displays the sync summary footer
func displaySyncSummary(created, updated, skipped, failed int, dryRun bool) {
	synced := created + updated
	parts := buildSummaryParts(created, updated, skipped, failed)

	var summary string
	if dryRun {
		summary = fmt.Sprintf("Would sync: %s\n→ Run without --check to actually sync", joinParts(parts))
	} else {
		summary = fmt.Sprintf("Synced: %s", joinParts(parts))
		if synced > 0 {
			repoURL, _ := GetRepositoryURL()
			summary += fmt.Sprintf("\n→ Verify at: %s/settings/secrets/actions", repoURL)
		}
	}

	printFooter(summary)
}
