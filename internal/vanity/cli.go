// cli.go - CLI entry point for vanity import package management.
//
// This file contains the CLI entry point. The main.go in cmd/vanityimport
// just imports and calls Run().
package vanity

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	clihelper "github.com/joeblew999/ubuntu-website/pkg/cli"
	"github.com/spf13/cobra"
)

const (
	vanityDomain = "www.ubuntusoftware.net"
	contentDir   = "content/english/pkg"
)

// Run is the main entry point for the vanity import CLI.
// Returns exit code (0 = success, 1 = error).
func Run(args []string, version string, stdout, stderr io.Writer) int {
	var githubIssue bool

	rootCmd := &cobra.Command{
		Use:     "vanityimport",
		Short:   "Manage Go vanity import packages",
		Long:    "CLI tool for managing Go packages with vanity imports on www.ubuntusoftware.net/pkg/",
		Version: version,
	}

	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.PersistentFlags().BoolVar(&githubIssue, "github-issue", false, "Output markdown for GitHub issue")

	cli := &cliRunner{stdout: stdout, stderr: stderr, githubIssue: &githubIssue}

	rootCmd.AddCommand(cli.addCmd())
	rootCmd.AddCommand(cli.listCmd())
	rootCmd.AddCommand(cli.updateCmd())
	rootCmd.AddCommand(cli.infoCmd())
	rootCmd.AddCommand(cli.docsCmd())

	rootCmd.SetArgs(args[1:])
	if err := rootCmd.Execute(); err != nil {
		return 1
	}
	return 0
}

type cliRunner struct {
	stdout      io.Writer
	stderr      io.Writer
	githubIssue *bool
}

func (c *cliRunner) addCmd() *cobra.Command {
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
			return c.addPackage(packageName, repoURL, description, author, version, license)
		},
	}

	cmd.Flags().StringVar(&repoURL, "repo", "", "GitHub repository URL (default: infer from package name)")
	cmd.Flags().StringVar(&description, "description", "", "Package description")
	cmd.Flags().StringVar(&author, "author", "Gerard Webb", "Package author")
	cmd.Flags().StringVar(&version, "version", "", "Package version (default: fetch from GitHub)")
	cmd.Flags().StringVar(&license, "license", "MIT", "Package license")

	return cmd
}

func (c *cliRunner) listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all packages",
		RunE: func(cmd *cobra.Command, args []string) error {
			packages, err := ListPackages(contentDir)
			if err != nil {
				return fmt.Errorf("failed to list packages: %w", err)
			}

			if len(packages) == 0 {
				fmt.Fprintln(c.stdout, "No packages found.")
				return nil
			}

			if *c.githubIssue {
				fmt.Fprintln(c.stdout, "## Open Source Packages")
				fmt.Fprintln(c.stdout)
				fmt.Fprintf(c.stdout, "Total: **%d**\n\n", len(packages))
				fmt.Fprintln(c.stdout, "| Package | Import Path | Version | License |")
				fmt.Fprintln(c.stdout, "|---------|-------------|---------|---------|")
				for _, pkgPath := range packages {
					pkg, err := ReadPackage(pkgPath)
					if err != nil {
						continue
					}
					fmt.Fprintf(c.stdout, "| [%s](%s) | `%s` | %s | %s |\n",
						pkg.Title, pkg.RepoURL, pkg.ImportPath, pkg.Version, pkg.License)
				}
				return nil
			}

			fmt.Fprintf(c.stdout, "Found %d package(s):\n\n", len(packages))
			for _, pkgPath := range packages {
				pkg, err := ReadPackage(pkgPath)
				if err != nil {
					log.Printf("Warning: could not read %s: %v", pkgPath, err)
					continue
				}

				fmt.Fprintf(c.stdout, "  %s\n", pkg.Title)
				fmt.Fprintf(c.stdout, "    Import: %s\n", pkg.ImportPath)
				fmt.Fprintf(c.stdout, "    Version: %s\n", pkg.Version)
				fmt.Fprintln(c.stdout)
			}

			return nil
		},
	}
}

func (c *cliRunner) updateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update [package-name]",
		Short: "Update package metadata from GitHub",
		Long:  "Fetch latest version and metadata from GitHub for one or all packages",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return c.updatePackage(args[0])
			}
			return c.updateAllPackages()
		},
	}
}

func (c *cliRunner) infoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info <package-name>",
		Short: "Show package information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pkg, err := GetPackage(contentDir, args[0])
			if err != nil {
				return fmt.Errorf("failed to read package: %w", err)
			}

			if *c.githubIssue {
				fmt.Fprintf(c.stdout, "## Package: %s\n\n", pkg.Title)
				fmt.Fprintf(c.stdout, "- **Import:** `go get %s`\n", pkg.ImportPath)
				fmt.Fprintf(c.stdout, "- **Repository:** [%s](%s)\n", pkg.RepoURL, pkg.RepoURL)
				fmt.Fprintf(c.stdout, "- **Version:** %s\n", pkg.Version)
				fmt.Fprintf(c.stdout, "- **License:** %s\n", pkg.License)
				fmt.Fprintf(c.stdout, "- **Author:** %s\n", pkg.Author)
				if pkg.Description != "" {
					fmt.Fprintf(c.stdout, "\n### Description\n\n%s\n", pkg.Description)
				}
				if pkg.DocumentationURL != "" {
					fmt.Fprintf(c.stdout, "\n### Links\n\n- [Documentation](%s)\n", pkg.DocumentationURL)
				}
				return nil
			}

			fmt.Fprintf(c.stdout, "Package: %s\n", pkg.Title)
			fmt.Fprintf(c.stdout, "  Import Path: %s\n", pkg.ImportPath)
			fmt.Fprintf(c.stdout, "  Repository: %s\n", pkg.RepoURL)
			fmt.Fprintf(c.stdout, "  Version: %s\n", pkg.Version)
			fmt.Fprintf(c.stdout, "  Description: %s\n", pkg.Description)
			fmt.Fprintf(c.stdout, "  License: %s\n", pkg.License)
			fmt.Fprintf(c.stdout, "  Author: %s\n", pkg.Author)
			fmt.Fprintf(c.stdout, "  Created: %s\n", pkg.CreatedAt.Format(time.RFC3339))
			fmt.Fprintf(c.stdout, "  Updated: %s\n", pkg.UpdatedAt.Format(time.RFC3339))
			if pkg.DocumentationURL != "" {
				fmt.Fprintf(c.stdout, "  Docs: %s\n", pkg.DocumentationURL)
			}

			return nil
		},
	}
}

func (c *cliRunner) docsCmd() *cobra.Command {
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
  vanityimport docs mailerlite

  # Generate and save to content directory
  vanityimport docs mailerlite --output content/english/pkg/mailerlite.md

  # Add custom features
  vanityimport docs mailerlite --feature "**Fast** - Optimized API calls"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			packageName := args[0]
			return c.generateDocs(packageName, output, features)
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file (default: stdout)")
	cmd.Flags().StringArrayVarP(&features, "feature", "f", nil, "Add feature bullet points")

	return cmd
}

func (c *cliRunner) generateDocs(packageName, output string, extraFeatures []string) error {
	// Try to read existing package metadata
	filePath := filepath.Join(contentDir, packageName+".md")
	pkg, err := ReadPackage(filePath)

	var doc *clihelper.DocBuilder

	if err != nil {
		// Package doesn't exist yet - create from scratch
		fmt.Fprintf(c.stderr, "Note: Package %s not in registry, creating template\n", packageName)
		doc = clihelper.NewDoc(packageName, "v0.1.0").
			Description("Package description here")
	} else {
		// Build from existing package metadata
		doc = clihelper.NewDoc(pkg.Title, pkg.Version).
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
		fmt.Fprintf(c.stdout, "✓ Generated %s\n", output)
	} else {
		fmt.Fprint(c.stdout, markdown)
	}

	return nil
}

func (c *cliRunner) addPackage(packageName, repoURL, description, author, version, license string) error {
	// Set default repo URL if not provided
	if repoURL == "" {
		repoURL = "https://github.com/joeblew999/ubuntu-website"
	}

	fmt.Fprintf(c.stdout, "Adding package '%s'...\n", packageName)

	// Try to fetch metadata from GitHub
	var pkg *Package
	var err error

	owner, repo, parseErr := ParseRepoURL(repoURL)
	if parseErr == nil {
		pkg, err = CreatePackageFromRepo(vanityDomain, packageName, repoURL)
		if err != nil {
			fmt.Fprintf(c.stdout, "Warning: could not fetch GitHub metadata: %v\n", err)
		}
	}

	// Create package manually if GitHub fetch failed
	if pkg == nil {
		now := time.Now()
		pkg = &Package{
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
	if err := WritePackage(filePath, pkg); err != nil {
		return fmt.Errorf("failed to write package: %w", err)
	}

	fmt.Fprintf(c.stdout, "✓ Created %s\n", filePath)

	// Validate site build
	fmt.Fprintln(c.stdout, "Validating Hugo build...")
	cmd := exec.Command("hugo", "--gc", "--minify")
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(c.stdout, "Warning: build validation failed: %v\n%s\n", err, output)
	} else {
		fmt.Fprintln(c.stdout, "✓ Site builds successfully")
	}

	// Print summary
	fmt.Fprintln(c.stdout, "\n=== Package Added ===")
	fmt.Fprintf(c.stdout, "Name: %s\n", pkg.Title)
	fmt.Fprintf(c.stdout, "Import: go get %s\n", pkg.ImportPath)
	fmt.Fprintf(c.stdout, "Repository: %s\n", pkg.RepoURL)
	fmt.Fprintf(c.stdout, "Version: %s\n", pkg.Version)

	fmt.Fprintln(c.stdout, "\nNext steps:")
	fmt.Fprintf(c.stdout, "1. Create pkg/%s/ directory with Go code\n", packageName)
	fmt.Fprintf(c.stdout, "2. Run: git add %s pkg/%s/\n", filePath, packageName)
	fmt.Fprintln(c.stdout, "3. Commit and push to deploy")

	// Suppress unused variable warning
	_ = owner
	_ = repo

	return nil
}

func (c *cliRunner) updatePackage(packageName string) error {
	filePath := filepath.Join(contentDir, packageName+".md")
	pkg, err := ReadPackage(filePath)
	if err != nil {
		return fmt.Errorf("failed to read package: %w", err)
	}

	fmt.Fprintf(c.stdout, "Updating %s...\n", packageName)

	owner, repo, err := ParseRepoURL(pkg.RepoURL)
	if err != nil {
		return fmt.Errorf("invalid repo URL: %w", err)
	}

	// Fetch latest version
	version, err := GetLatestVersion(owner, repo)
	if err != nil {
		fmt.Fprintf(c.stdout, "Warning: could not fetch version: %v\n", err)
	} else if version != "" && version != pkg.Version {
		fmt.Fprintf(c.stdout, "  Version: %s → %s\n", pkg.Version, version)
		pkg.Version = version
	} else {
		fmt.Fprintln(c.stdout, "  Version: up to date")
	}

	pkg.UpdatedAt = time.Now()

	if err := WritePackage(filePath, pkg); err != nil {
		return fmt.Errorf("failed to write package: %w", err)
	}

	fmt.Fprintln(c.stdout, "✓ Updated")
	return nil
}

func (c *cliRunner) updateAllPackages() error {
	packages, err := ListPackages(contentDir)
	if err != nil {
		return fmt.Errorf("failed to list packages: %w", err)
	}

	fmt.Fprintf(c.stdout, "Updating %d package(s)...\n\n", len(packages))

	for _, pkgPath := range packages {
		packageName := filepath.Base(pkgPath)
		packageName = packageName[:len(packageName)-3] // Remove .md
		if err := c.updatePackage(packageName); err != nil {
			log.Printf("Error updating %s: %v", packageName, err)
		}
		fmt.Fprintln(c.stdout)
	}

	return nil
}
