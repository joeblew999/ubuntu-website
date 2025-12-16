package cli

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// PackageMeta holds metadata for generating package documentation.
type PackageMeta struct {
	// Required
	Name        string // Package name (e.g., "mailerlite")
	ImportPath  string // Full import path (e.g., "www.ubuntusoftware.net/pkg/mailerlite")
	Description string // Short description
	Version     string // Semantic version (e.g., "v1.0.0")

	// Optional
	RepoURL          string   // GitHub repository URL
	DocumentationURL string   // pkg.go.dev URL (auto-generated if empty)
	License          string   // License (default: MIT)
	Author           string   // Author name (default: Gerard Webb)
	HasBinary        bool     // Has downloadable binaries
	InstallCommand   string   // Alternative install command (e.g., "brew install ...")
	Features         []string // Feature bullet points
	Commands         []Command // CLI commands
	Examples         []Example // Code examples
}

// Command represents a CLI command for documentation.
type Command struct {
	Name        string // Command name
	Description string // What it does
	Example     string // Example usage
}

// Example represents a code example.
type Example struct {
	Title    string // Example title
	Language string // Code language (go, bash, etc.)
	Code     string // The code
}

// GenerateFrontmatter generates Hugo frontmatter YAML.
func (m *PackageMeta) GenerateFrontmatter(w io.Writer) {
	now := time.Now().Format(time.RFC3339)

	// Set defaults
	if m.License == "" {
		m.License = "MIT"
	}
	if m.Author == "" {
		m.Author = "Gerard Webb"
	}
	if m.DocumentationURL == "" && m.ImportPath != "" {
		m.DocumentationURL = fmt.Sprintf("https://pkg.go.dev/%s", m.ImportPath)
	}
	if m.RepoURL == "" {
		m.RepoURL = "https://github.com/joeblew999/ubuntu-website"
	}

	fmt.Fprintln(w, "---")
	fmt.Fprintf(w, "title: %s\n", m.Name)
	fmt.Fprintf(w, "import_path: %s\n", m.ImportPath)
	fmt.Fprintf(w, "repo_url: %s\n", m.RepoURL)
	fmt.Fprintf(w, "description: %s\n", m.Description)
	fmt.Fprintf(w, "version: %s\n", m.Version)
	fmt.Fprintf(w, "documentation_url: %s\n", m.DocumentationURL)
	fmt.Fprintf(w, "license: %s\n", m.License)
	fmt.Fprintf(w, "author: %s\n", m.Author)
	fmt.Fprintf(w, "created_at: %s\n", now)
	fmt.Fprintf(w, "updated_at: %s\n", now)
	if m.HasBinary {
		fmt.Fprintln(w, "has_binary: true")
	}
	if m.InstallCommand != "" {
		fmt.Fprintf(w, "install_command: %s\n", m.InstallCommand)
	}
	fmt.Fprintln(w, "---")
}

// GenerateMarkdown generates the full package markdown documentation.
func (m *PackageMeta) GenerateMarkdown(w io.Writer) {
	m.GenerateFrontmatter(w)
	fmt.Fprintln(w)

	// Features section
	if len(m.Features) > 0 {
		fmt.Fprintln(w, "## Features")
		fmt.Fprintln(w)
		for _, f := range m.Features {
			fmt.Fprintf(w, "- %s\n", f)
		}
		fmt.Fprintln(w)
	}

	// CLI Commands section
	if len(m.Commands) > 0 {
		fmt.Fprintln(w, "## CLI Usage")
		fmt.Fprintln(w)
		fmt.Fprintln(w, "```bash")
		for _, cmd := range m.Commands {
			if cmd.Description != "" {
				fmt.Fprintf(w, "# %s\n", cmd.Description)
			}
			fmt.Fprintln(w, cmd.Example)
			fmt.Fprintln(w)
		}
		fmt.Fprintln(w, "```")
		fmt.Fprintln(w)
	}

	// Examples section
	for _, ex := range m.Examples {
		if ex.Title != "" {
			fmt.Fprintf(w, "## %s\n\n", ex.Title)
		}
		lang := ex.Language
		if lang == "" {
			lang = "go"
		}
		fmt.Fprintf(w, "```%s\n", lang)
		fmt.Fprintln(w, strings.TrimSpace(ex.Code))
		fmt.Fprintln(w, "```")
		fmt.Fprintln(w)
	}
}

// String returns the full markdown as a string.
func (m *PackageMeta) String() string {
	var b strings.Builder
	m.GenerateMarkdown(&b)
	return b.String()
}

// DocBuilder helps build package documentation fluently.
type DocBuilder struct {
	meta PackageMeta
}

// NewDoc creates a new documentation builder.
func NewDoc(name, version string) *DocBuilder {
	return &DocBuilder{
		meta: PackageMeta{
			Name:       name,
			ImportPath: fmt.Sprintf("www.ubuntusoftware.net/pkg/%s", name),
			Version:    version,
		},
	}
}

// Description sets the package description.
func (b *DocBuilder) Description(desc string) *DocBuilder {
	b.meta.Description = desc
	return b
}

// Repo sets the repository URL.
func (b *DocBuilder) Repo(url string) *DocBuilder {
	b.meta.RepoURL = url
	return b
}

// Author sets the author.
func (b *DocBuilder) Author(author string) *DocBuilder {
	b.meta.Author = author
	return b
}

// License sets the license.
func (b *DocBuilder) License(license string) *DocBuilder {
	b.meta.License = license
	return b
}

// HasBinary marks the package as having downloadable binaries.
func (b *DocBuilder) HasBinary() *DocBuilder {
	b.meta.HasBinary = true
	return b
}

// InstallCommand sets an alternative install command.
func (b *DocBuilder) InstallCommand(cmd string) *DocBuilder {
	b.meta.InstallCommand = cmd
	return b
}

// Feature adds a feature bullet point.
func (b *DocBuilder) Feature(feature string) *DocBuilder {
	b.meta.Features = append(b.meta.Features, feature)
	return b
}

// Command adds a CLI command example.
func (b *DocBuilder) Command(name, description, example string) *DocBuilder {
	b.meta.Commands = append(b.meta.Commands, Command{
		Name:        name,
		Description: description,
		Example:     example,
	})
	return b
}

// Example adds a code example.
func (b *DocBuilder) Example(title, language, code string) *DocBuilder {
	b.meta.Examples = append(b.meta.Examples, Example{
		Title:    title,
		Language: language,
		Code:     code,
	})
	return b
}

// Build returns the PackageMeta.
func (b *DocBuilder) Build() *PackageMeta {
	return &b.meta
}

// String returns the markdown string.
func (b *DocBuilder) String() string {
	return b.meta.String()
}

// Write writes the markdown to a writer.
func (b *DocBuilder) Write(w io.Writer) {
	b.meta.GenerateMarkdown(w)
}

// ImportPath sets a custom import path.
func (b *DocBuilder) ImportPath(path string) *DocBuilder {
	b.meta.ImportPath = path
	return b
}

// DocumentationURL sets the documentation URL.
func (b *DocBuilder) DocumentationURL(url string) *DocBuilder {
	b.meta.DocumentationURL = url
	return b
}

// InstallInstruction returns a formatted install instruction.
func (m *PackageMeta) InstallInstruction() string {
	if m.InstallCommand != "" {
		return m.InstallCommand
	}
	return "go get " + m.ImportPath
}

// Version is the package version.
const DocVersion = "v0.1.0"

// PackageDoc returns the package documentation for registry publishing.
func PackageDoc() *DocBuilder {
	return NewDoc("cli", DocVersion).
		Description("Shared CLI framework for Ubuntu Software tools. Standard flags, output formatting, and GitHub issue markdown.").
		Repo("https://github.com/joeblew999/ubuntu-website").
		Feature("**Standard flags** - `--github-issue`, `--verbose`, `--version` out of the box").
		Feature("**Output formatting** - Tables, lists, headers, code blocks").
		Feature("**Markdown mode** - GitHub issue-friendly output with `--github-issue`").
		Feature("**Cross-platform** - Browser opening for macOS, Linux, Windows").
		Feature("**Context handling** - Timeouts and cancellation").
		Example("Usage", "go", `
package main

import (
    "fmt"
    "os"

    "www.ubuntusoftware.net/pkg/cli"
)

func main() {
    app := cli.New("myapp", "v1.0.0")

    err := app.Run(os.Args[1:], func(c *cli.Context) error {
        if len(c.Args) == 0 {
            return fmt.Errorf("command required")
        }

        switch c.Args[0] {
        case "list":
            c.Header("Items")
            table := c.NewTable("NAME", "VALUE")
            table.Row("foo", "bar")
            table.Row("baz", "qux")
            table.Flush()

        case "open":
            return c.Open("https://example.com")
        }

        return nil
    })

    if err != nil {
        fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
        os.Exit(1)
    }
}
`)
}
