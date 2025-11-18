package env

import (
	"fmt"
)

// RunBuild runs Hugo build only (no deployment) for CLI
func RunBuild() error {
	fmt.Println("Building Hugo site...")
	fmt.Println()

	result := BuildHugoSite(false)

	fmt.Println(result.Output)

	if result.Error != nil {
		return fmt.Errorf("build failed: %w", result.Error)
	}

	if result.LocalURL != "" {
		fmt.Printf("\n✓ Build complete!\n")
		fmt.Printf("\nLocal preview available at:\n  %s\n", result.LocalURL)
		if result.LANURL != "" {
			fmt.Printf("\nMobile/LAN preview available at:\n  %s\n", result.LANURL)
		}
	}

	return nil
}

// RunDeployPreview runs build + deploy to Cloudflare Pages preview for CLI
func RunDeployPreview() error {
	// Load config to get project name
	svc := NewService(false)
	cfg, err := svc.GetCurrentConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	projectName := cfg.Get(KeyCloudflarePageProject)
	if projectName == "" || IsPlaceholder(projectName) {
		return fmt.Errorf("no Cloudflare Pages project configured. Run 'web-gui' and complete Step 4 first")
	}

	fmt.Printf("Building and deploying to Cloudflare Pages (preview)...\n")
	fmt.Printf("Project: %s\n", projectName)
	fmt.Println()

	// Run build and deploy (no branch = preview only)
	result := BuildAndDeploy(projectName, "", false)

	fmt.Println(result.Output)

	if result.Error != nil {
		return fmt.Errorf("deployment failed: %w", result.Error)
	}

	// Print URLs
	fmt.Printf("\n✓ Deployment complete!\n")
	if result.PreviewURL != "" {
		fmt.Printf("\nPreview URL:\n  %s\n", result.PreviewURL)
	}

	return nil
}

// RunDeployProduction runs build + deploy to Cloudflare Pages production for CLI
func RunDeployProduction() error {
	// Load config to get project name and custom domain
	svc := NewService(false)
	cfg, err := svc.GetCurrentConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	projectName := cfg.Get(KeyCloudflarePageProject)
	if projectName == "" || IsPlaceholder(projectName) {
		return fmt.Errorf("no Cloudflare Pages project configured. Run 'web-gui' and complete Step 4 first")
	}

	customDomain := cfg.Get(KeyCloudflareDomain)

	fmt.Printf("Building and deploying to Cloudflare Pages (production)...\n")
	fmt.Printf("Project: %s\n", projectName)
	if customDomain != "" && !IsPlaceholder(customDomain) {
		fmt.Printf("Custom domain: %s\n", customDomain)
	}
	fmt.Println()

	// Run build and deploy (branch=main = production)
	result := BuildAndDeploy(projectName, "main", false)

	fmt.Println(result.Output)

	if result.Error != nil {
		return fmt.Errorf("deployment failed: %w", result.Error)
	}

	// Print URLs
	fmt.Printf("\n✓ Deployment complete!\n")
	if result.PreviewURL != "" {
		fmt.Printf("\nPreview URL:\n  %s\n", result.PreviewURL)
	}

	// Show custom domain if configured
	if customDomain != "" && !IsPlaceholder(customDomain) {
		fmt.Printf("\nProduction URL (custom domain):\n  https://%s\n", customDomain)
	} else if result.DeploymentURL != "" {
		fmt.Printf("\nProduction URL:\n  %s\n", result.DeploymentURL)
	}

	return nil
}
