// cli.go - Command-line interface for the translate tool.
//
// This file provides the CLI entry point and command routing.
// The main.go in cmd/translate just imports and calls Run().
package translator

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// CLIOptions holds global CLI flags
type CLIOptions struct {
	GithubIssue bool
	Force       bool
	Version     bool
}

// Run is the main entry point for the translate CLI.
// Returns exit code (0 = success, 1 = error).
func Run(args []string, version string, stdout, stderr io.Writer, stdin io.Reader) int {
	// Parse flags
	fs := flag.NewFlagSet("translate", flag.ContinueOnError)
	fs.SetOutput(stderr)

	opts := &CLIOptions{}
	fs.BoolVar(&opts.GithubIssue, "github-issue", false, "Output markdown for GitHub Issue")
	fs.BoolVar(&opts.Force, "force", false, "Skip confirmation prompts (for CI)")
	fs.BoolVar(&opts.Version, "version", false, "Print version and exit")

	if err := fs.Parse(args[1:]); err != nil {
		return 1
	}

	if opts.Version {
		fmt.Fprintf(stdout, "translate %s\n", version)
		return 0
	}

	if fs.NArg() < 1 {
		printUsage(stderr)
		return 1
	}

	// Create checker instance
	checker, err := NewChecker()
	if err != nil {
		fmt.Fprintf(stderr, "Error: %v\n", err)
		return 1
	}

	namespace := fs.Arg(0)
	subCmd := fs.Arg(1)

	// Create CLI context
	ctx := &cliContext{
		checker:     checker,
		opts:        opts,
		stdout:      stdout,
		stderr:      stderr,
		stdin:       stdin,
		fs:          fs,
		allArgs:     args,
	}

	switch namespace {
	case "content":
		return ctx.runContentCommand(subCmd)
	case "menu":
		return ctx.runMenuCommand(subCmd)
	case "lang":
		return ctx.runLangCommand(subCmd)
	default:
		fmt.Fprintf(stderr, "Unknown namespace: %s\n", namespace)
		fmt.Fprintf(stderr, "Available: content, menu, lang\n")
		printUsage(stderr)
		return 1
	}
}

// cliContext holds state for CLI execution
type cliContext struct {
	checker *Checker
	opts    *CLIOptions
	stdout  io.Writer
	stderr  io.Writer
	stdin   io.Reader
	fs      *flag.FlagSet
	allArgs []string
}

// ============================================================================
// Content Commands - Track English source changes and translation problems
// ============================================================================

func (ctx *cliContext) runContentCommand(subCmd string) int {
	switch subCmd {
	case "status":
		return ctx.runStatus()
	case "diff":
		file := ctx.fs.Arg(2)
		if file == "" {
			fmt.Fprintln(ctx.stderr, "Error: content diff requires a file argument")
			fmt.Fprintln(ctx.stderr, "Usage: translate content diff <file>")
			return 1
		}
		return ctx.runDiff(file)
	case "changed":
		return ctx.runChanged()
	case "next":
		return ctx.runNext()
	case "done":
		return ctx.runDone()
	case "missing":
		return ctx.runMissing()
	case "orphans":
		return ctx.runOrphans()
	case "stale":
		return ctx.runStale()
	case "clean":
		return ctx.runClean()
	case "":
		fmt.Fprintln(ctx.stderr, "Error: content requires a subcommand")
		printContentUsage(ctx.stderr)
		return 1
	default:
		fmt.Fprintf(ctx.stderr, "Unknown content command: %s\n", subCmd)
		printContentUsage(ctx.stderr)
		return 1
	}
}

// ============================================================================
// Menu Commands - Navigation menu management
// ============================================================================

func (ctx *cliContext) runMenuCommand(subCmd string) int {
	switch subCmd {
	case "check":
		return ctx.runMenuCheck()
	case "sync":
		return ctx.runMenuSync()
	case "":
		fmt.Fprintln(ctx.stderr, "Error: menu requires a subcommand")
		printMenuUsage(ctx.stderr)
		return 1
	default:
		fmt.Fprintf(ctx.stderr, "Unknown menu command: %s\n", subCmd)
		printMenuUsage(ctx.stderr)
		return 1
	}
}

// ============================================================================
// Lang Commands - Language management
// ============================================================================

func (ctx *cliContext) runLangCommand(subCmd string) int {
	switch subCmd {
	case "list":
		return ctx.runLangs()
	case "validate":
		return ctx.runValidate()
	case "add":
		return ctx.runLangAdd()
	case "remove":
		return ctx.runLangRemove()
	case "init":
		return ctx.runLangInit()
	case "":
		fmt.Fprintln(ctx.stderr, "Error: lang requires a subcommand")
		printLangUsage(ctx.stderr)
		return 1
	default:
		fmt.Fprintf(ctx.stderr, "Unknown lang command: %s\n", subCmd)
		printLangUsage(ctx.stderr)
		return 1
	}
}

// ============================================================================
// Run Functions - Wire checker functions to presenters
// ============================================================================

func (ctx *cliContext) runStatus() int {
	result := ctx.checker.CheckStatus()

	if ctx.opts.GithubIssue {
		p := NewMarkdownPresenter()
		p.Status(result)
		if result.HasIssues() {
			return 1
		}
		return 0
	}

	p := NewTerminalPresenter()
	p.Status(result)
	return 0
}

func (ctx *cliContext) runDiff(file string) int {
	result := ctx.checker.CheckDiff(file)

	if result.Error != nil {
		fmt.Fprintf(ctx.stderr, "ERROR: File not found: %s\n", file)
		return 1
	}

	p := NewTerminalPresenter()
	p.Diff(result)
	return 0
}

func (ctx *cliContext) runMissing() int {
	result := ctx.checker.CheckMissing()

	if ctx.opts.GithubIssue {
		p := NewMarkdownPresenter()
		p.Missing(result)
		if result.HasIssues() {
			return 1
		}
		return 0
	}

	p := NewTerminalPresenter()
	p.Missing(result)
	return 0
}

func (ctx *cliContext) runStale() int {
	result := ctx.checker.CheckStale()

	if ctx.opts.GithubIssue {
		p := NewMarkdownPresenter()
		p.Stale(result)
		if result.HasIssues() {
			return 1
		}
		return 0
	}

	p := NewTerminalPresenter()
	p.Stale(result)
	return 0
}

func (ctx *cliContext) runOrphans() int {
	result := ctx.checker.CheckOrphans()

	if ctx.opts.GithubIssue {
		p := NewMarkdownPresenter()
		p.Orphans(result)
		if result.HasIssues() {
			return 1
		}
		return 0
	}

	p := NewTerminalPresenter()
	p.Orphans(result)
	return 0
}

func (ctx *cliContext) runClean() int {
	// First pass: get what would be deleted
	result := ctx.checker.DoClean(ctx.opts.Force, false)

	if result.TotalCount == 0 {
		fmt.Fprintln(ctx.stdout, "OK: No orphaned files to delete")
		return 0
	}

	// Show what will be deleted
	p := NewTerminalPresenter()
	p.Clean(result)

	// If force, skip confirmation
	if !ctx.opts.Force {
		fmt.Fprint(ctx.stdout, "\nDelete these files? [y/N]: ")
		reader := bufio.NewReader(ctx.stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(response)
		if response != "y" && response != "Y" {
			fmt.Fprintln(ctx.stdout, "Cancelled")
			return 0
		}
	}

	// Second pass: actually delete
	result = ctx.checker.DoClean(ctx.opts.Force, true)
	if result.Error != nil {
		fmt.Fprintf(ctx.stderr, "Error: %v\n", result.Error)
		return 1
	}

	fmt.Fprintf(ctx.stdout, "\nOK: Deleted %d orphaned files\n", result.TotalCount)
	return 0
}

func (ctx *cliContext) runDone() int {
	result := ctx.checker.DoDone()

	p := NewTerminalPresenter()
	p.Done(result)

	if result.Error != nil {
		return 1
	}
	return 0
}

func (ctx *cliContext) runNext() int {
	result := ctx.checker.CheckNext()

	p := NewTerminalPresenter()
	p.Next(result)
	return 0
}

func (ctx *cliContext) runChanged() int {
	result := ctx.checker.CheckChanged()

	p := NewTerminalPresenter()
	p.Changed(result)
	return 0
}

func (ctx *cliContext) runValidate() int {
	result := ctx.checker.CheckValidate()
	config := ctx.checker.GetConfig()

	// For validate, we need more context in output
	fmt.Fprintln(ctx.stdout, "========================================")
	fmt.Fprintln(ctx.stdout, "Validating Translator Configuration")
	fmt.Fprintln(ctx.stdout, "========================================")
	fmt.Fprintln(ctx.stdout)

	// Check if this is a Hugo project
	if !IsHugoProject() {
		fmt.Fprintln(ctx.stdout, "Mode: Standalone (no Hugo config found)")
		fmt.Fprintln(ctx.stdout)
		fmt.Fprintln(ctx.stdout, "Current configuration:")
		fmt.Fprintf(ctx.stdout, "  Source: %s -> content/%s\n", config.SourceLang, config.SourceDir)
		for _, lang := range config.TargetLangs {
			fmt.Fprintf(ctx.stdout, "  Target: %s (%s) -> content/%s\n", lang.Code, lang.Name, lang.DirName)
		}
		fmt.Fprintln(ctx.stdout)
		fmt.Fprintln(ctx.stdout, "========================================")
		fmt.Fprintln(ctx.stdout, "OK: Using default configuration")
		fmt.Fprintln(ctx.stdout, "========================================")
		return 0
	}

	fmt.Fprintln(ctx.stdout, "Mode: Hugo project detected")
	fmt.Fprintln(ctx.stdout)

	// Show current config (auto-loaded from Hugo)
	fmt.Fprintf(ctx.stdout, "Source: %s -> content/%s\n", config.SourceLang, config.SourceDir)
	for _, lang := range config.TargetLangs {
		fmt.Fprintf(ctx.stdout, "Target: %s (%s) -> content/%s\n", lang.Code, lang.Name, lang.DirName)
	}
	fmt.Fprintln(ctx.stdout)

	if result.HasIssues() {
		fmt.Fprintln(ctx.stdout, "========================================")
		fmt.Fprintf(ctx.stdout, "WARNING: %d mismatch(es) found\n", len(result.Mismatches))
		fmt.Fprintln(ctx.stdout, "========================================")
		for _, m := range result.Mismatches {
			fmt.Fprintf(ctx.stdout, "  - %s\n", m)
		}
		fmt.Fprintln(ctx.stdout)
		fmt.Fprintln(ctx.stdout, "This shouldn't happen - languages are auto-loaded from Hugo config.")
		fmt.Fprintln(ctx.stdout, "Check if config/_default/languages.toml changed after binary was built.")
		return 1
	}

	fmt.Fprintln(ctx.stdout, "========================================")
	fmt.Fprintln(ctx.stdout, "OK: Configuration loaded from Hugo")
	fmt.Fprintln(ctx.stdout, "========================================")
	return 0
}

func (ctx *cliContext) runLangs() int {
	result := ctx.checker.CheckLangs()

	p := NewTerminalPresenter()
	p.Langs(result)

	if result.HasIssues() {
		return 1
	}
	return 0
}

func (ctx *cliContext) runMenuCheck() int {
	result := ctx.checker.CheckMenu()

	if ctx.opts.GithubIssue {
		p := NewMarkdownPresenter()
		p.MenuCheck(result)
		if result.HasIssues() {
			return 1
		}
		return 0
	}

	p := NewTerminalPresenter()
	p.MenuCheck(result)
	return 0
}

func (ctx *cliContext) runMenuSync() int {
	// Show header with source info
	enMenuPath := GetMenuFilePath("en")
	enMenu, err := ParseMenuFile(enMenuPath)
	if err != nil {
		fmt.Fprintf(ctx.stderr, "Error reading English menu: %v\n", err)
		return 1
	}

	fmt.Fprintln(ctx.stdout, "========================================")
	fmt.Fprintln(ctx.stdout, "Menu Sync")
	fmt.Fprintln(ctx.stdout, "========================================")
	fmt.Fprintln(ctx.stdout)
	fmt.Fprintf(ctx.stdout, "Source: %s (%d main items, %d footer items)\n", enMenuPath, len(enMenu.Main), len(enMenu.Footer))
	fmt.Fprintln(ctx.stdout)

	result := ctx.checker.DoMenuSync()

	if result.Error != nil {
		fmt.Fprintf(ctx.stderr, "Error: %v\n", result.Error)
		return 1
	}

	for _, path := range result.FilesWritten {
		fmt.Fprintf(ctx.stdout, "Generated: %s\n", path)
	}

	fmt.Fprintln(ctx.stdout)
	fmt.Fprintln(ctx.stdout, "========================================")
	fmt.Fprintln(ctx.stdout, "OK: Menu files regenerated from English")
	fmt.Fprintln(ctx.stdout, "========================================")
	return 0
}

// ============================================================================
// Language Management Implementations
// ============================================================================

func (ctx *cliContext) runLangAdd() int {
	if ctx.fs.NArg() < 5 {
		fmt.Fprintln(ctx.stderr, "Error: lang add requires <code> <name> <dirname>")
		fmt.Fprintln(ctx.stderr, "Usage: translate lang add fr \"Francais\" french")
		return 1
	}

	code := ctx.fs.Arg(2)
	name := ctx.fs.Arg(3)
	dirname := ctx.fs.Arg(4)

	fmt.Fprintln(ctx.stdout, "========================================")
	fmt.Fprintf(ctx.stdout, "Adding language: %s (%s)\n", code, name)
	fmt.Fprintln(ctx.stdout, "========================================")
	fmt.Fprintln(ctx.stdout)

	result := ctx.checker.DoLangAdd(code, name, dirname)

	if result.Error != nil {
		fmt.Fprintf(ctx.stderr, "Error: %v\n", result.Error)
		return 1
	}

	fmt.Fprintf(ctx.stdout, "1. Adding to %s...\n", result.ConfigPath)
	fmt.Fprintln(ctx.stdout, "   OK")
	fmt.Fprintf(ctx.stdout, "2. Creating %s...\n", result.ContentPath)
	fmt.Fprintln(ctx.stdout, "   OK")
	if result.MenuPath != "" {
		fmt.Fprintf(ctx.stdout, "3. Generating %s...\n", result.MenuPath)
		fmt.Fprintln(ctx.stdout, "   OK")
	}

	fmt.Fprintln(ctx.stdout)
	fmt.Fprintln(ctx.stdout, "========================================")
	fmt.Fprintf(ctx.stdout, "OK: Language '%s' added\n", code)
	fmt.Fprintln(ctx.stdout, "========================================")
	fmt.Fprintln(ctx.stdout)
	fmt.Fprintln(ctx.stdout, "Next steps:")
	fmt.Fprintln(ctx.stdout, "  1. Run 'task translate:content:missing' to see what needs translating")
	fmt.Fprintf(ctx.stdout, "  2. Translate content to %s\n", name)
	fmt.Fprintln(ctx.stdout, "  3. Run 'task translate:content:done' when complete")
	return 0
}

func (ctx *cliContext) runLangRemove() int {
	if ctx.fs.NArg() < 3 {
		fmt.Fprintln(ctx.stderr, "Error: lang remove requires <code>")
		fmt.Fprintln(ctx.stderr, "Usage: translate lang remove fr [-force]")
		return 1
	}

	code := ctx.fs.Arg(2)
	config := ctx.checker.GetConfig()

	fmt.Fprintln(ctx.stdout, "========================================")
	fmt.Fprintf(ctx.stdout, "Removing language: %s\n", code)
	fmt.Fprintln(ctx.stdout, "========================================")
	fmt.Fprintln(ctx.stdout)

	// Check if language exists
	existing, err := GetLanguageByCode(code)
	if err != nil {
		fmt.Fprintf(ctx.stderr, "Error reading config: %v\n", err)
		return 1
	}
	if existing == nil {
		fmt.Fprintf(ctx.stdout, "Language '%s' not found in config\n", code)
		return 1
	}

	// Get content directory
	dirname := strings.TrimPrefix(existing.ContentDir, "content/")
	contentPath := filepath.Join(config.ContentDir, dirname)

	// Count files in content directory
	fileCount := 0
	if _, err := os.Stat(contentPath); err == nil {
		filepath.Walk(contentPath, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && strings.HasSuffix(path, ".md") {
				fileCount++
			}
			return nil
		})
	}

	// Show what will be deleted
	fmt.Fprintf(ctx.stdout, "Content directory: %s (%d .md files)\n", contentPath, fileCount)
	fmt.Fprintf(ctx.stdout, "Config: config/_default/languages.toml [%s] section\n", code)
	fmt.Fprintf(ctx.stdout, "Menu: config/_default/menus.%s.toml\n", code)
	fmt.Fprintln(ctx.stdout)

	// Confirm unless force
	if !ctx.opts.Force && fileCount > 0 {
		fmt.Fprintf(ctx.stdout, "Delete %d files and remove language '%s'? [y/N]: ", fileCount, code)
		reader := bufio.NewReader(ctx.stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(response)
		if response != "y" && response != "Y" {
			fmt.Fprintln(ctx.stdout, "Cancelled")
			return 0
		}
	}

	result := ctx.checker.DoLangRemove(code, ctx.opts.Force, true)

	if result.Error != nil {
		fmt.Fprintf(ctx.stderr, "Error: %v\n", result.Error)
		return 1
	}

	fmt.Fprintln(ctx.stdout, "1. Removing from languages.toml...")
	fmt.Fprintln(ctx.stdout, "   OK")
	if result.FilesRemoved > 0 {
		fmt.Fprintf(ctx.stdout, "2. Deleting %s...\n", contentPath)
		fmt.Fprintln(ctx.stdout, "   OK")
	}
	fmt.Fprintf(ctx.stdout, "3. Deleting menus.%s.toml...\n", code)
	fmt.Fprintln(ctx.stdout, "   OK")

	fmt.Fprintln(ctx.stdout)
	fmt.Fprintln(ctx.stdout, "========================================")
	fmt.Fprintf(ctx.stdout, "OK: Language '%s' removed\n", code)
	fmt.Fprintln(ctx.stdout, "========================================")
	return 0
}

func (ctx *cliContext) runLangInit() int {
	if ctx.fs.NArg() < 3 {
		fmt.Fprintln(ctx.stderr, "Error: lang init requires <code>")
		fmt.Fprintln(ctx.stderr, "Usage: translate lang init fr")
		return 1
	}

	code := ctx.fs.Arg(2)

	fmt.Fprintln(ctx.stdout, "========================================")
	fmt.Fprintf(ctx.stdout, "Initializing content for language: %s\n", code)
	fmt.Fprintln(ctx.stdout, "========================================")
	fmt.Fprintln(ctx.stdout)

	result := ctx.checker.DoLangInit(code)

	if result.Error != nil {
		fmt.Fprintf(ctx.stderr, "Error: %v\n", result.Error)
		return 1
	}

	if result.AlreadyExists {
		fmt.Fprintf(ctx.stdout, "Directory %s already exists (%d .md files)\n", result.Path, result.FileCount)
		fmt.Fprintln(ctx.stdout, "========================================")
		return 0
	}

	fmt.Fprintf(ctx.stdout, "Creating %s...\n", result.Path)
	fmt.Fprintln(ctx.stdout, "OK")

	fmt.Fprintln(ctx.stdout)
	fmt.Fprintln(ctx.stdout, "========================================")
	fmt.Fprintf(ctx.stdout, "OK: Content directory initialized for '%s'\n", code)
	fmt.Fprintln(ctx.stdout, "========================================")
	fmt.Fprintln(ctx.stdout)
	fmt.Fprintln(ctx.stdout, "Next steps:")
	fmt.Fprintln(ctx.stdout, "  Run 'task translate:content:missing' to see what needs translating")
	return 0
}

// ============================================================================
// Usage
// ============================================================================

func printUsage(w io.Writer) {
	fmt.Fprintf(w, `translate - Translation workflow for Hugo multilingual content

Usage:
  translate <namespace> <command> [args] [flags]

Namespaces:
  content   Track English source changes and find translation problems
  menu      Manage navigation menus per language
  lang      Add, remove, and configure languages

Flags:
  -github-issue  Output markdown for GitHub Issue (exit 1 if action needed)
  -force         Skip confirmation prompts (for CI)
  -version       Print version and exit

Run 'translate <namespace>' for namespace-specific help.

Examples:
  translate content status              # See what English files changed
  translate content missing             # See what's missing in translations
  translate menu check                  # Validate menu files
  translate lang list                   # Show configured languages
  translate lang add fr Francais french # Add French language

`)
}

func printContentUsage(w io.Writer) {
	fmt.Fprintf(w, `translate content - Track English source changes and translation problems

Commands:
  status            Show what English files changed since last translation
  diff <file>       Show git diff for specific file since checkpoint
  changed           Show detailed changes for all files
  next              Show next file to translate with progress
  done              Mark translations complete (update checkpoint)
  missing           Show files missing in target languages
  orphans           Show target files with no English source
  stale             Show potentially outdated translations (target < 50%% of source)
  clean             Delete orphaned files (prompts unless -force)

Examples:
  translate content status
  translate content diff blog/my-post.md
  translate content missing -github-issue
  translate content clean -force

`)
}

func printMenuUsage(w io.Writer) {
	fmt.Fprintf(w, `translate menu - Manage navigation menus per language

Commands:
  check             Validate menu files for broken links and sync issues
  sync              Generate translated menu files from English

Examples:
  translate menu check
  translate menu check -github-issue
  translate menu sync

`)
}

func printLangUsage(w io.Writer) {
	fmt.Fprintf(w, `translate lang - Add, remove, and configure languages

Commands:
  list              Show configured languages and detect stray directories
  validate          Check translator config matches Hugo config
  add <code> <name> <dirname>   Add a new target language
  remove <code>     Remove a language (prompts unless -force)
  init <code>       Initialize content directory for configured language

Examples:
  translate lang list
  translate lang add fr "Francais" french
  translate lang add ko "Korean" korean
  translate lang remove fr
  translate lang remove fr -force
  translate lang init fr

`)
}
