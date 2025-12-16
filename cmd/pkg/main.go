package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/joeblew999/ubuntu-website/internal/vanity"
	"github.com/joeblew999/ubuntu-website/pkg/cli"
	"github.com/spf13/cobra"
)

const (
	vanityDomain = "www.ubuntusoftware.net"
	contentDir   = "content/english/pkg"
)

var githubIssue bool

func main() {
	rootCmd := &cobra.Command{
		Use:   "pkg",
		Short: "Manage Go vanity import packages",
		Long:  "CLI tool for managing Go packages with vanity imports on www.ubuntusoftware.net/pkg/",
	}

	rootCmd.PersistentFlags().BoolVar(&githubIssue, "github-issue", false, "Output markdown for GitHub issue")

	rootCmd.AddCommand(addCmd())
	rootCmd.AddCommand(listCmd())
	rootCmd.AddCommand(updateCmd())
	rootCmd.AddCommand(infoCmd())
	rootCmd.AddCommand(docsCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func addCmd() *cobra.Command {
	var (
		repoURL     string
		description string
		author      string
		version     string
		license     string
	)

	cmd := &cobra.Command{
		Use:   "add <package-name>",
		Short: "Add a new package",
		Long:  "Add a new Go package with vanity import to the website",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			packageName := args[0]
			return addPackage(packageName, repoURL, description, author, version, license)
		},
	}

	cmd.Flags().StringVar(&repoURL, "repo", "", "GitHub repository URL (default: infer from package name)")
	cmd.Flags().StringVar(&description, "description", "", "Package description")
	cmd.Flags().StringVar(&author, "author", "Gerard Webb", "Package author")
	cmd.Flags().StringVar(&version, "version", "", "Package version (default: fetch from GitHub)")
	cmd.Flags().StringVar(&license, "license", "MIT", "Package license")

	return cmd
}

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all packages",
		RunE: func(cmd *cobra.Command, args []string) error {
			packages, err := vanity.ListPackages(contentDir)
			if err != nil {
				return fmt.Errorf("failed to list packages: %w", err)
			}

			if len(packages) == 0 {
				fmt.Println("No packages found.")
				return nil
			}

			if githubIssue {
				fmt.Println("## Open Source Packages")
				fmt.Println()
				fmt.Printf("Total: **%d**\n\n", len(packages))
				fmt.Println("| Package | Import Path | Version | License |")
				fmt.Println("|---------|-------------|---------|---------|")
				for _, pkgPath := range packages {
					pkg, err := vanity.ReadPackage(pkgPath)
					if err != nil {
						continue
					}
					fmt.Printf("| [%s](%s) | `%s` | %s | %s |\n",
						pkg.Title, pkg.RepoURL, pkg.ImportPath, pkg.Version, pkg.License)
				}
				return nil
			}

			fmt.Printf("Found %d package(s):\n\n", len(packages))
			for _, pkgPath := range packages {
				pkg, err := vanity.ReadPackage(pkgPath)
				if err != nil {
					log.Printf("Warning: could not read %s: %v", pkgPath, err)
					continue
				}

				fmt.Printf("  %s\n", pkg.Title)
				fmt.Printf("    Import: %s\n", pkg.ImportPath)
				fmt.Printf("    Version: %s\n", pkg.Version)
				fmt.Println()
			}

			return nil
		},
	}
}

func updateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update [package-name]",
		Short: "Update package metadata from GitHub",
		Long:  "Fetch latest version and metadata from GitHub for one or all packages",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return updatePackage(args[0])
			}
			return updateAllPackages()
		},
	}
}

func infoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info <package-name>",
		Short: "Show package information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pkg, err := vanity.GetPackage(contentDir, args[0])
			if err != nil {
				return fmt.Errorf("failed to read package: %w", err)
			}

			if githubIssue {
				fmt.Printf("## Package: %s\n\n", pkg.Title)
				fmt.Printf("- **Import:** `go get %s`\n", pkg.ImportPath)
				fmt.Printf("- **Repository:** [%s](%s)\n", pkg.RepoURL, pkg.RepoURL)
				fmt.Printf("- **Version:** %s\n", pkg.Version)
				fmt.Printf("- **License:** %s\n", pkg.License)
				fmt.Printf("- **Author:** %s\n", pkg.Author)
				if pkg.Description != "" {
					fmt.Printf("\n### Description\n\n%s\n", pkg.Description)
				}
				if pkg.DocumentationURL != "" {
					fmt.Printf("\n### Links\n\n- [Documentation](%s)\n", pkg.DocumentationURL)
				}
				return nil
			}

			fmt.Printf("Package: %s\n", pkg.Title)
			fmt.Printf("  Import Path: %s\n", pkg.ImportPath)
			fmt.Printf("  Repository: %s\n", pkg.RepoURL)
			fmt.Printf("  Version: %s\n", pkg.Version)
			fmt.Printf("  Description: %s\n", pkg.Description)
			fmt.Printf("  License: %s\n", pkg.License)
			fmt.Printf("  Author: %s\n", pkg.Author)
			fmt.Printf("  Created: %s\n", pkg.CreatedAt.Format(time.RFC3339))
			fmt.Printf("  Updated: %s\n", pkg.UpdatedAt.Format(time.RFC3339))
			if pkg.DocumentationURL != "" {
				fmt.Printf("  Docs: %s\n", pkg.DocumentationURL)
			}

			return nil
		},
	}
}

func docsCmd() *cobra.Command {
	var (
		output   string
		features []string
	)

	cmd := &cobra.Command{
		Use:   "docs <package-name>",
		Short: "Generate markdown documentation for a package",
		Long: `Generate Hugo-compatible markdown documentation for a package.

This uses the package metadata from the registry and can be customized
with additional features and examples.

Examples:
  # Generate docs to stdout
  pkg docs mailerlite

  # Generate and save to content directory
  pkg docs mailerlite --output content/english/pkg/mailerlite.md

  # Add custom features
  pkg docs mailerlite --feature "**Fast** - Optimized API calls"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			packageName := args[0]
			return generateDocs(packageName, output, features)
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file (default: stdout)")
	cmd.Flags().StringArrayVarP(&features, "feature", "f", nil, "Add feature bullet points")

	return cmd
}

func generateDocs(packageName, output string, extraFeatures []string) error {
	// Try to read existing package metadata
	filePath := filepath.Join(contentDir, packageName+".md")
	pkg, err := vanity.ReadPackage(filePath)

	var doc *cli.DocBuilder

	if err != nil {
		// Package doesn't exist yet - create from scratch
		fmt.Fprintf(os.Stderr, "Note: Package %s not in registry, creating template\n", packageName)
		doc = cli.NewDoc(packageName, "v0.1.0").
			Description("Package description here")
	} else {
		// Build from existing package metadata
		doc = cli.NewDoc(pkg.Title, pkg.Version).
			Description(pkg.Description).
			Repo(pkg.RepoURL).
			Author(pkg.Author).
			License(pkg.License)

		if pkg.HasBinary {
			doc.HasBinary()
		}
	}

	// Add extra features from flags
	for _, f := range extraFeatures {
		doc.Feature(f)
	}

	// Generate markdown
	markdown := doc.String()

	// Output to file or stdout
	if output != "" {
		if err := os.WriteFile(output, []byte(markdown), 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
		fmt.Printf("✓ Generated %s\n", output)
	} else {
		fmt.Print(markdown)
	}

	return nil
}

func addPackage(packageName, repoURL, description, author, version, license string) error {
	// Set default repo URL if not provided
	if repoURL == "" {
		repoURL = "https://github.com/joeblew999/ubuntu-website"
	}

	fmt.Printf("Adding package '%s'...\n", packageName)

	// Try to fetch metadata from GitHub
	var pkg *vanity.Package
	var err error

	owner, repo, parseErr := vanity.ParseRepoURL(repoURL)
	if parseErr == nil {
		pkg, err = vanity.CreatePackageFromRepo(vanityDomain, packageName, repoURL)
		if err != nil {
			fmt.Printf("Warning: could not fetch GitHub metadata: %v\n", err)
		}
	}

	// Create package manually if GitHub fetch failed
	if pkg == nil {
		now := time.Now()
		pkg = &vanity.Package{
			Title:            packageName,
			ImportPath:       fmt.Sprintf("%s/pkg/%s", vanityDomain, packageName),
			RepoURL:          repoURL,
			Description:      description,
			Version:          version,
			DocumentationURL: fmt.Sprintf("https://pkg.go.dev/%s/pkg/%s", vanityDomain, packageName),
			License:          license,
			Author:           author,
			CreatedAt:        now,
			UpdatedAt:        now,
		}
	}

	// Override with provided values
	if description != "" {
		pkg.Description = description
	}
	if author != "" {
		pkg.Author = author
	}
	if version != "" {
		pkg.Version = version
	}
	if license != "" {
		pkg.License = license
	}

	// Set defaults if still empty
	if pkg.Version == "" {
		pkg.Version = "v0.1.0"
	}
	if pkg.License == "" {
		pkg.License = "MIT"
	}
	if pkg.Author == "" {
		pkg.Author = "Gerard Webb"
	}

	// Check if file already exists
	filePath := filepath.Join(contentDir, packageName+".md")
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("package already exists: %s", filePath)
	}

	// Write package file
	if err := vanity.WritePackage(filePath, pkg); err != nil {
		return fmt.Errorf("failed to write package: %w", err)
	}

	fmt.Printf("✓ Created %s\n", filePath)

	// Validate site build
	fmt.Println("Validating Hugo build...")
	cmd := exec.Command("hugo", "--gc", "--minify")
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("Warning: build validation failed: %v\n%s\n", err, output)
	} else {
		fmt.Println("✓ Site builds successfully")
	}

	// Print summary
	fmt.Println("\n=== Package Added ===")
	fmt.Printf("Name: %s\n", pkg.Title)
	fmt.Printf("Import: go get %s\n", pkg.ImportPath)
	fmt.Printf("Repository: %s\n", pkg.RepoURL)
	fmt.Printf("Version: %s\n", pkg.Version)

	fmt.Println("\nNext steps:")
	fmt.Printf("1. Create pkg/%s/ directory with Go code\n", packageName)
	fmt.Printf("2. Run: git add %s pkg/%s/\n", filePath, packageName)
	fmt.Printf("3. Commit and push to deploy\n")

	// Suppress unused variable warning
	_ = owner
	_ = repo

	return nil
}

func updatePackage(packageName string) error {
	filePath := filepath.Join(contentDir, packageName+".md")
	pkg, err := vanity.ReadPackage(filePath)
	if err != nil {
		return fmt.Errorf("failed to read package: %w", err)
	}

	fmt.Printf("Updating %s...\n", packageName)

	owner, repo, err := vanity.ParseRepoURL(pkg.RepoURL)
	if err != nil {
		return fmt.Errorf("invalid repo URL: %w", err)
	}

	// Fetch latest version
	version, err := vanity.GetLatestVersion(owner, repo)
	if err != nil {
		fmt.Printf("Warning: could not fetch version: %v\n", err)
	} else if version != "" && version != pkg.Version {
		fmt.Printf("  Version: %s → %s\n", pkg.Version, version)
		pkg.Version = version
	} else {
		fmt.Println("  Version: up to date")
	}

	pkg.UpdatedAt = time.Now()

	if err := vanity.WritePackage(filePath, pkg); err != nil {
		return fmt.Errorf("failed to write package: %w", err)
	}

	fmt.Println("✓ Updated")
	return nil
}

func updateAllPackages() error {
	packages, err := vanity.ListPackages(contentDir)
	if err != nil {
		return fmt.Errorf("failed to list packages: %w", err)
	}

	fmt.Printf("Updating %d package(s)...\n\n", len(packages))

	for _, pkgPath := range packages {
		packageName := filepath.Base(pkgPath)
		packageName = packageName[:len(packageName)-3] // Remove .md
		if err := updatePackage(packageName); err != nil {
			log.Printf("Error updating %s: %v", packageName, err)
		}
		fmt.Println()
	}

	return nil
}
