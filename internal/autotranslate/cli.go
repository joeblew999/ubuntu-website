// cli.go - CLI entry point for autotranslate command.
//
// This file contains the CLI entry point. The main.go in cmd/autotranslate
// just imports and calls Run().
package autotranslate

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/joeblew999/ubuntu-website/internal/translate"

	appcli "github.com/joeblew999/ubuntu-website/internal/cli"
)

// CLIOptions holds the parsed CLI flags
type CLIOptions struct {
	ProviderName string
	DryRun       bool
	Verbose      bool
	APIKey       string
}

// Run is the main entry point for the autotranslate CLI.
// Returns exit code (0 = success, 1 = error).
func Run(args []string, version string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("autotranslate", flag.ContinueOnError)
	fs.SetOutput(stderr)

	opts := &CLIOptions{}
	fs.StringVar(&opts.ProviderName, "provider", "deepl", "Translation provider (deepl, claude, claude-cli)")
	fs.BoolVar(&opts.DryRun, "dry-run", false, "Show what would be translated without actually translating")
	fs.BoolVar(&opts.Verbose, "verbose", false, "Verbose output")
	fs.StringVar(&opts.APIKey, "api-key", "", "API key (or use DEEPL_API_KEY/CLAUDE_API_KEY env var)")

	fs.Usage = func() {
		fmt.Fprintf(stderr, "Usage: autotranslate [flags] <command> [args]\n\n")
		fmt.Fprintf(stderr, "Commands:\n")
		fmt.Fprintf(stderr, "  file <source-file> <target-lang>   Translate a single file\n")
		fmt.Fprintf(stderr, "  missing <target-lang>              Translate all missing files for language\n")
		fmt.Fprintf(stderr, "  languages                          List supported languages\n")
		fmt.Fprintf(stderr, "  status                             Show API status and usage\n")
		fmt.Fprintf(stderr, "\nFlags:\n")
		fs.PrintDefaults()
		fmt.Fprintf(stderr, "\nEnvironment:\n")
		fmt.Fprintf(stderr, "  DEEPL_API_KEY    DeepL API key (free tier ends with :fx)\n")
		fmt.Fprintf(stderr, "  CLAUDE_API_KEY   Claude API key (uses your subscription)\n")
	}

	if err := fs.Parse(appcli.Ordered(fs, args[1:])); err != nil {
		return 1
	}

	if fs.NArg() < 1 {
		fs.Usage()
		return 1
	}

	cmd := fs.Arg(0)
	cli := &cliRunner{opts: opts, stdout: stdout, stderr: stderr, args: fs.Args()}

	switch cmd {
	case "file":
		if fs.NArg() < 3 {
			fmt.Fprintf(stderr, "Usage: autotranslate file <source-file> <target-lang>\n")
			return 1
		}
		return cli.runFileTranslation(fs.Arg(1), fs.Arg(2))

	case "missing":
		if fs.NArg() < 2 {
			fmt.Fprintf(stderr, "Usage: autotranslate missing <target-lang>\n")
			return 1
		}
		return cli.runMissingTranslation(fs.Arg(1))

	case "languages":
		return cli.runListLanguages()

	case "status":
		return cli.runStatus()

	default:
		fmt.Fprintf(stderr, "Unknown command: %s\n", cmd)
		fs.Usage()
		return 1
	}
}

type cliRunner struct {
	opts   *CLIOptions
	stdout io.Writer
	stderr io.Writer
	args   []string
}

func (c *cliRunner) getProvider() (Provider, error) {
	key := c.opts.APIKey

	switch c.opts.ProviderName {
	case "deepl":
		if key == "" {
			key = os.Getenv("DEEPL_API_KEY")
		}
		if key == "" {
			return nil, fmt.Errorf("DEEPL_API_KEY environment variable not set")
		}
		return NewDeepLProvider(key)

	case "claude":
		if key == "" {
			key = os.Getenv("CLAUDE_API_KEY")
		}
		if key == "" {
			return nil, fmt.Errorf("CLAUDE_API_KEY environment variable not set")
		}
		return NewClaudeProvider(key)

	case "claude-cli":
		// Uses logged-in Claude CLI session - no API key needed
		return NewClaudeCLIProvider()

	default:
		return nil, fmt.Errorf("unsupported provider: %s (available: deepl, claude, claude-cli)", c.opts.ProviderName)
	}
}

func (c *cliRunner) runFileTranslation(sourcePath, targetLang string) int {
	ctx := context.Background()

	// Get provider
	provider, err := c.getProvider()
	if err != nil {
		fmt.Fprintf(c.stderr, "Error: %v\n", err)
		return 1
	}

	// Check language support
	if !provider.SupportsLanguage(targetLang) {
		fmt.Fprintf(c.stderr, "Error: language '%s' not supported by %s\n", targetLang, provider.Name())
		return 1
	}

	// Read source file
	content, err := os.ReadFile(sourcePath)
	if err != nil {
		fmt.Fprintf(c.stderr, "Error reading file: %v\n", err)
		return 1
	}

	// Determine target path
	targetPath := c.sourceToTargetPath(sourcePath, targetLang)

	if c.opts.DryRun {
		fmt.Fprintf(c.stdout, "Would translate:\n")
		fmt.Fprintf(c.stdout, "  Source: %s\n", sourcePath)
		fmt.Fprintf(c.stdout, "  Target: %s\n", targetPath)
		fmt.Fprintf(c.stdout, "  Lang:   en → %s\n", targetLang)
		fmt.Fprintf(c.stdout, "  Chars:  %d\n", len(content))
		return 0
	}

	// Translate
	mt := NewMarkdownTranslator(provider)
	translated, err := mt.TranslateFile(ctx, string(content), "en", targetLang)
	if err != nil {
		fmt.Fprintf(c.stderr, "Error translating: %v\n", err)
		return 1
	}

	// Ensure target directory exists
	targetDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		fmt.Fprintf(c.stderr, "Error creating directory: %v\n", err)
		return 1
	}

	// Write translated file
	if err := os.WriteFile(targetPath, []byte(translated), 0644); err != nil {
		fmt.Fprintf(c.stderr, "Error writing file: %v\n", err)
		return 1
	}

	fmt.Fprintf(c.stdout, "✓ Translated: %s → %s\n", sourcePath, targetPath)
	return 0
}

func (c *cliRunner) runMissingTranslation(targetLang string) int {
	ctx := context.Background()

	// Load Hugo config to get languages
	config := translate.DefaultConfig()
	if err := translate.TryLoadHugoConfig(config); err != nil {
		fmt.Fprintf(c.stderr, "Error loading config: %v\n", err)
		return 1
	}

	// Find target language directory
	var targetDir string
	for _, lang := range config.TargetLangs {
		if lang.Code == targetLang {
			targetDir = lang.DirName
			break
		}
	}
	if targetDir == "" {
		fmt.Fprintf(c.stderr, "Error: language '%s' not configured in Hugo\n", targetLang)
		fmt.Fprintf(c.stderr, "Available: ")
		for _, lang := range config.TargetLangs {
			fmt.Fprintf(c.stderr, "%s ", lang.Code)
		}
		fmt.Fprintf(c.stderr, "\n")
		return 1
	}

	// Find all English files
	sourceDir := filepath.Join("content", config.SourceDir)
	var missingFiles []string
	var totalChars int

	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		// Calculate target path
		relPath, _ := filepath.Rel(sourceDir, path)
		targetPath := filepath.Join("content", targetDir, relPath)

		// Check if target exists
		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			missingFiles = append(missingFiles, path)
			content, _ := os.ReadFile(path)
			totalChars += len(content)
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(c.stderr, "Error scanning files: %v\n", err)
		return 1
	}

	if len(missingFiles) == 0 {
		fmt.Fprintf(c.stdout, "✓ All files already translated to %s\n", targetLang)
		return 0
	}

	fmt.Fprintf(c.stdout, "Found %d files missing translation to %s\n", len(missingFiles), targetLang)
	fmt.Fprintf(c.stdout, "Total characters: ~%d\n\n", totalChars)

	if c.opts.DryRun {
		fmt.Fprintf(c.stdout, "Files to translate:\n")
		for _, f := range missingFiles {
			info, _ := os.Stat(f)
			size := int64(0)
			if info != nil {
				size = info.Size()
			}
			fmt.Fprintf(c.stdout, "  %s (%s chars)\n", f, formatNumber(size))
		}
		fmt.Fprintf(c.stdout, "\nEstimated total: ~%s characters\n", formatNumber(int64(totalChars)))
		fmt.Fprintf(c.stdout, "(Actual usage may be lower - front matter, code blocks, etc. are not translated)\n")
		return 0
	}

	// Get provider (only needed when actually translating)
	provider, err := c.getProvider()
	if err != nil {
		fmt.Fprintf(c.stderr, "Error: %v\n", err)
		return 1
	}

	if !provider.SupportsLanguage(targetLang) {
		fmt.Fprintf(c.stderr, "Error: language '%s' not supported by %s\n", targetLang, provider.Name())
		return 1
	}

	// Get usage before translation
	deeplProvider, _ := provider.(*DeepLProvider)
	var usageBefore *Usage
	if deeplProvider != nil {
		usageBefore, _ = deeplProvider.GetUsage(ctx)
		if usageBefore != nil {
			fmt.Fprintf(c.stdout, "API Usage before: %s / %s characters (%.1f%%)\n\n",
				formatNumber(usageBefore.CharacterCount),
				formatNumber(usageBefore.CharacterLimit),
				float64(usageBefore.CharacterCount)/float64(usageBefore.CharacterLimit)*100)
		}
	}

	// Translate each file
	mt := NewMarkdownTranslator(provider)
	successCount := 0
	errorCount := 0

	for i, sourcePath := range missingFiles {
		content, err := os.ReadFile(sourcePath)
		if err != nil {
			fmt.Fprintf(c.stderr, "✗ Error reading %s: %v\n", sourcePath, err)
			errorCount++
			continue
		}

		relPath, _ := filepath.Rel(sourceDir, sourcePath)
		targetPath := filepath.Join("content", targetDir, relPath)

		if c.opts.Verbose {
			fmt.Fprintf(c.stdout, "[%d/%d] Translating %s...\n", i+1, len(missingFiles), relPath)
		}

		translated, err := mt.TranslateFile(ctx, string(content), "en", targetLang)
		if err != nil {
			fmt.Fprintf(c.stderr, "✗ Error translating %s: %v\n", sourcePath, err)
			errorCount++
			continue
		}

		// Ensure target directory exists
		targetDirPath := filepath.Dir(targetPath)
		if err := os.MkdirAll(targetDirPath, 0755); err != nil {
			fmt.Fprintf(c.stderr, "✗ Error creating directory for %s: %v\n", targetPath, err)
			errorCount++
			continue
		}

		if err := os.WriteFile(targetPath, []byte(translated), 0644); err != nil {
			fmt.Fprintf(c.stderr, "✗ Error writing %s: %v\n", targetPath, err)
			errorCount++
			continue
		}

		successCount++
		if !c.opts.Verbose {
			fmt.Fprintf(c.stdout, "✓ [%d/%d] %s\n", i+1, len(missingFiles), relPath)
		} else {
			fmt.Fprintf(c.stdout, "✓ Wrote %s\n", targetPath)
		}
	}

	fmt.Fprintf(c.stdout, "\nComplete: %d translated, %d errors\n", successCount, errorCount)

	// Show usage after translation
	if deeplProvider != nil {
		usageAfter, err := deeplProvider.GetUsage(ctx)
		if err == nil && usageAfter != nil {
			charsUsed := usageAfter.CharacterCount
			if usageBefore != nil {
				charsUsed = usageAfter.CharacterCount - usageBefore.CharacterCount
			}
			fmt.Fprintf(c.stdout, "\nAPI Usage:\n")
			fmt.Fprintf(c.stdout, "  Characters used this run: %s\n", formatNumber(charsUsed))
			fmt.Fprintf(c.stdout, "  Total used this period:   %s / %s (%.1f%%)\n",
				formatNumber(usageAfter.CharacterCount),
				formatNumber(usageAfter.CharacterLimit),
				float64(usageAfter.CharacterCount)/float64(usageAfter.CharacterLimit)*100)
			fmt.Fprintf(c.stdout, "  Remaining:                %s characters\n",
				formatNumber(usageAfter.CharacterLimit-usageAfter.CharacterCount))
		}
	}

	return 0
}

func (c *cliRunner) runListLanguages() int {
	provider, err := c.getProvider()
	if err != nil {
		// Show DeepL languages even without API key
		fmt.Fprintln(c.stdout, "Supported languages (DeepL):")
		langs := []string{"en", "de", "fr", "es", "it", "nl", "pl", "pt", "ru", "ja", "zh", "ko", "vi", "id", "tr", "uk", "cs", "da", "fi", "el", "hu", "lt", "lv", "nb", "ro", "sk", "sl", "sv", "bg", "et"}
		for _, l := range langs {
			fmt.Fprintf(c.stdout, "  %s\n", l)
		}
		return 0
	}

	fmt.Fprintf(c.stdout, "Supported languages (%s):\n", provider.Name())
	for _, lang := range provider.SupportedLanguages() {
		fmt.Fprintf(c.stdout, "  %s\n", lang)
	}
	return 0
}

func (c *cliRunner) runStatus() int {
	fmt.Fprintln(c.stdout, "========================================")
	fmt.Fprintln(c.stdout, "Translation Provider Status")
	fmt.Fprintln(c.stdout, "========================================")
	fmt.Fprintln(c.stdout)

	// Show DeepL status
	deeplKey := os.Getenv("DEEPL_API_KEY")
	if deeplKey != "" {
		isFree := strings.HasSuffix(deeplKey, ":fx")
		fmt.Fprintln(c.stdout, "--- DeepL ---")
		fmt.Fprintln(c.stdout, "Status: Configured")
		if isFree {
			fmt.Fprintln(c.stdout, "Plan: Free (500k chars/month)")
		} else {
			fmt.Fprintln(c.stdout, "Plan: Pro")
		}
		fmt.Fprintf(c.stdout, "API Key: %s...%s\n", deeplKey[:8], deeplKey[len(deeplKey)-4:])

		// Get actual usage from API
		provider, err := NewDeepLProvider(deeplKey)
		if err == nil {
			usage, err := provider.GetUsage(context.Background())
			if err == nil {
				fmt.Fprintf(c.stdout, "Usage: %s / %s (%.1f%%)\n",
					formatNumber(usage.CharacterCount),
					formatNumber(usage.CharacterLimit),
					float64(usage.CharacterCount)/float64(usage.CharacterLimit)*100)
			}
		}
		fmt.Fprintln(c.stdout)
	} else {
		fmt.Fprintln(c.stdout, "--- DeepL ---")
		fmt.Fprintln(c.stdout, "Status: Not configured")
		fmt.Fprintln(c.stdout, "Set DEEPL_API_KEY for 500k free chars/month")
		fmt.Fprintln(c.stdout)
	}

	// Show Claude API status
	claudeKey := os.Getenv("CLAUDE_API_KEY")
	if claudeKey != "" {
		fmt.Fprintln(c.stdout, "--- Claude (API) ---")
		fmt.Fprintln(c.stdout, "Status: Configured")
		fmt.Fprintln(c.stdout, "Plan: Requires API credits")
		if len(claudeKey) > 12 {
			fmt.Fprintf(c.stdout, "API Key: %s...%s\n", claudeKey[:8], claudeKey[len(claudeKey)-4:])
		}
		fmt.Fprintln(c.stdout)
	} else {
		fmt.Fprintln(c.stdout, "--- Claude (API) ---")
		fmt.Fprintln(c.stdout, "Status: Not configured")
		fmt.Fprintln(c.stdout, "Set CLAUDE_API_KEY for API-based translation")
		fmt.Fprintln(c.stdout)
	}

	// Show Claude CLI status
	fmt.Fprintln(c.stdout, "--- Claude (CLI) ---")
	if _, err := NewClaudeCLIProvider(); err == nil {
		fmt.Fprintln(c.stdout, "Status: Available")
		fmt.Fprintln(c.stdout, "Plan: Uses logged-in session (your subscription)")
		fmt.Fprintln(c.stdout, "Binary: claude")
	} else {
		fmt.Fprintln(c.stdout, "Status: Not available")
		fmt.Fprintln(c.stdout, "Install: mise run claude-cli:ensure")
	}
	fmt.Fprintln(c.stdout)

	fmt.Fprintln(c.stdout, "========================================")
	fmt.Fprintf(c.stdout, "Current provider: %s\n", c.opts.ProviderName)
	fmt.Fprintln(c.stdout, "Use --provider=deepl, --provider=claude, or --provider=claude-cli")
	fmt.Fprintln(c.stdout, "========================================")
	return 0
}

func (c *cliRunner) sourceToTargetPath(sourcePath, targetLang string) string {
	// Load Hugo config to get target directory name
	config := translate.DefaultConfig()
	translate.TryLoadHugoConfig(config)

	var targetDir string
	for _, lang := range config.TargetLangs {
		if lang.Code == targetLang {
			targetDir = lang.DirName
			break
		}
	}
	if targetDir == "" {
		targetDir = targetLang // Fallback to lang code
	}

	// Replace source directory with target
	sourceDir := filepath.Join("content", config.SourceDir)
	if strings.HasPrefix(sourcePath, sourceDir) {
		relPath, _ := filepath.Rel(sourceDir, sourcePath)
		return filepath.Join("content", targetDir, relPath)
	}

	// Fallback: just replace "english" with target dir
	return strings.Replace(sourcePath, "/english/", "/"+targetDir+"/", 1)
}

// formatNumber adds comma separators to large numbers
func formatNumber(n int64) string {
	str := fmt.Sprintf("%d", n)
	if len(str) <= 3 {
		return str
	}

	var result []byte
	for i, c := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(c))
	}
	return string(result)
}
