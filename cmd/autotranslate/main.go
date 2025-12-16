// Command autotranslate provides automatic translation of Hugo markdown content
// using external translation APIs (DeepL or Claude).
//
// Usage:
//
//	autotranslate [flags] <command> [args]
//
// Commands:
//
//	file     Translate a single file
//	missing  Translate all missing files for a language
//	batch    Translate multiple files
//	status   Show translation quota/usage
//
// Examples:
//
//	# Translate a single file to Vietnamese using DeepL
//	autotranslate file content/english/blog/post.md vi
//
//	# Translate using Claude (no per-character cost if you have subscription)
//	autotranslate --provider=claude missing vi
//
//	# Dry-run to see what would be translated
//	autotranslate missing vi --dry-run
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joeblew999/ubuntu-website/internal/autotranslate"
	"github.com/joeblew999/ubuntu-website/internal/translator"
)

var (
	providerName = flag.String("provider", "deepl", "Translation provider (deepl, claude, claude-cli)")
	dryRun       = flag.Bool("dry-run", false, "Show what would be translated without actually translating")
	verbose      = flag.Bool("verbose", false, "Verbose output")
	apiKey       = flag.String("api-key", "", "API key (or use DEEPL_API_KEY/CLAUDE_API_KEY env var)")
	bundlePath   = flag.String("bundle", "tokibundle", "Path to tokibundle directory for ARB translation")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: autotranslate [flags] <command> [args]\n\n")
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  file <source-file> <target-lang>   Translate a single file\n")
		fmt.Fprintf(os.Stderr, "  missing <target-lang>              Translate all missing files for language\n")
		fmt.Fprintf(os.Stderr, "  arb <target-lang>                  Translate ARB catalog entries (toki workflow)\n")
		fmt.Fprintf(os.Stderr, "  arb-status                         Show ARB translation status\n")
		fmt.Fprintf(os.Stderr, "  languages                          List supported languages\n")
		fmt.Fprintf(os.Stderr, "  status                             Show API status and usage\n")
		fmt.Fprintf(os.Stderr, "\nFlags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEnvironment:\n")
		fmt.Fprintf(os.Stderr, "  DEEPL_API_KEY    DeepL API key (free tier ends with :fx)\n")
		fmt.Fprintf(os.Stderr, "  CLAUDE_API_KEY   Claude API key (uses your subscription)\n")
	}

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	cmd := flag.Arg(0)

	switch cmd {
	case "file":
		if flag.NArg() < 3 {
			fmt.Fprintf(os.Stderr, "Usage: autotranslate file <source-file> <target-lang>\n")
			os.Exit(1)
		}
		runFileTranslation(flag.Arg(1), flag.Arg(2))

	case "missing":
		if flag.NArg() < 2 {
			fmt.Fprintf(os.Stderr, "Usage: autotranslate missing <target-lang>\n")
			os.Exit(1)
		}
		runMissingTranslation(flag.Arg(1))

	case "arb":
		if flag.NArg() < 2 {
			fmt.Fprintf(os.Stderr, "Usage: autotranslate arb <target-lang>\n")
			os.Exit(1)
		}
		runARBTranslation(flag.Arg(1))

	case "arb-status":
		runARBStatus()

	case "languages":
		runListLanguages()

	case "status":
		runStatus()

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		flag.Usage()
		os.Exit(1)
	}
}

func getProvider() (autotranslate.Provider, error) {
	key := *apiKey

	switch *providerName {
	case "deepl":
		if key == "" {
			key = os.Getenv("DEEPL_API_KEY")
		}
		if key == "" {
			return nil, fmt.Errorf("DEEPL_API_KEY environment variable not set")
		}
		return autotranslate.NewDeepLProvider(key)

	case "claude":
		if key == "" {
			key = os.Getenv("CLAUDE_API_KEY")
		}
		if key == "" {
			return nil, fmt.Errorf("CLAUDE_API_KEY environment variable not set")
		}
		return autotranslate.NewClaudeProvider(key)

	case "claude-cli":
		// Uses logged-in Claude CLI session - no API key needed
		return autotranslate.NewClaudeCLIProvider()

	default:
		return nil, fmt.Errorf("unsupported provider: %s (available: deepl, claude, claude-cli)", *providerName)
	}
}

func runFileTranslation(sourcePath, targetLang string) {
	ctx := context.Background()

	// Get provider
	provider, err := getProvider()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Check language support
	if !provider.SupportsLanguage(targetLang) {
		fmt.Fprintf(os.Stderr, "Error: language '%s' not supported by %s\n", targetLang, provider.Name())
		os.Exit(1)
	}

	// Read source file
	content, err := os.ReadFile(sourcePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Determine target path
	targetPath := sourceToTargetPath(sourcePath, targetLang)

	if *dryRun {
		fmt.Printf("Would translate:\n")
		fmt.Printf("  Source: %s\n", sourcePath)
		fmt.Printf("  Target: %s\n", targetPath)
		fmt.Printf("  Lang:   en → %s\n", targetLang)
		fmt.Printf("  Chars:  %d\n", len(content))
		return
	}

	// Translate
	mt := autotranslate.NewMarkdownTranslator(provider)
	translated, err := mt.TranslateFile(ctx, string(content), "en", targetLang)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error translating: %v\n", err)
		os.Exit(1)
	}

	// Ensure target directory exists
	targetDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
		os.Exit(1)
	}

	// Write translated file
	if err := os.WriteFile(targetPath, []byte(translated), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Translated: %s → %s\n", sourcePath, targetPath)
}

func runMissingTranslation(targetLang string) {
	ctx := context.Background()

	// Load Hugo config to get languages
	config := translator.DefaultConfig()
	if err := translator.TryLoadHugoConfig(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
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
		fmt.Fprintf(os.Stderr, "Error: language '%s' not configured in Hugo\n", targetLang)
		fmt.Fprintf(os.Stderr, "Available: ")
		for _, lang := range config.TargetLangs {
			fmt.Fprintf(os.Stderr, "%s ", lang.Code)
		}
		fmt.Fprintf(os.Stderr, "\n")
		os.Exit(1)
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
		fmt.Fprintf(os.Stderr, "Error scanning files: %v\n", err)
		os.Exit(1)
	}

	if len(missingFiles) == 0 {
		fmt.Printf("✓ All files already translated to %s\n", targetLang)
		return
	}

	fmt.Printf("Found %d files missing translation to %s\n", len(missingFiles), targetLang)
	fmt.Printf("Total characters: ~%d\n\n", totalChars)

	if *dryRun {
		fmt.Printf("Files to translate:\n")
		for _, f := range missingFiles {
			info, _ := os.Stat(f)
			size := int64(0)
			if info != nil {
				size = info.Size()
			}
			fmt.Printf("  %s (%s chars)\n", f, formatNumber(size))
		}
		fmt.Printf("\nEstimated total: ~%s characters\n", formatNumber(int64(totalChars)))
		fmt.Printf("(Actual usage may be lower - front matter, code blocks, etc. are not translated)\n")
		return
	}

	// Get provider (only needed when actually translating)
	provider, err := getProvider()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if !provider.SupportsLanguage(targetLang) {
		fmt.Fprintf(os.Stderr, "Error: language '%s' not supported by %s\n", targetLang, provider.Name())
		os.Exit(1)
	}

	// Get usage before translation
	deeplProvider, _ := provider.(*autotranslate.DeepLProvider)
	var usageBefore *autotranslate.Usage
	if deeplProvider != nil {
		usageBefore, _ = deeplProvider.GetUsage(ctx)
		if usageBefore != nil {
			fmt.Printf("API Usage before: %s / %s characters (%.1f%%)\n\n",
				formatNumber(usageBefore.CharacterCount),
				formatNumber(usageBefore.CharacterLimit),
				float64(usageBefore.CharacterCount)/float64(usageBefore.CharacterLimit)*100)
		}
	}

	// Translate each file
	mt := autotranslate.NewMarkdownTranslator(provider)
	successCount := 0
	errorCount := 0

	for i, sourcePath := range missingFiles {
		content, err := os.ReadFile(sourcePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "✗ Error reading %s: %v\n", sourcePath, err)
			errorCount++
			continue
		}

		relPath, _ := filepath.Rel(sourceDir, sourcePath)
		targetPath := filepath.Join("content", targetDir, relPath)

		if *verbose {
			fmt.Printf("[%d/%d] Translating %s...\n", i+1, len(missingFiles), relPath)
		}

		translated, err := mt.TranslateFile(ctx, string(content), "en", targetLang)
		if err != nil {
			fmt.Fprintf(os.Stderr, "✗ Error translating %s: %v\n", sourcePath, err)
			errorCount++
			continue
		}

		// Ensure target directory exists
		targetDirPath := filepath.Dir(targetPath)
		if err := os.MkdirAll(targetDirPath, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "✗ Error creating directory for %s: %v\n", targetPath, err)
			errorCount++
			continue
		}

		if err := os.WriteFile(targetPath, []byte(translated), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "✗ Error writing %s: %v\n", targetPath, err)
			errorCount++
			continue
		}

		successCount++
		if !*verbose {
			fmt.Printf("✓ [%d/%d] %s\n", i+1, len(missingFiles), relPath)
		} else {
			fmt.Printf("✓ Wrote %s\n", targetPath)
		}
	}

	fmt.Printf("\nComplete: %d translated, %d errors\n", successCount, errorCount)

	// Show usage after translation
	if deeplProvider != nil {
		usageAfter, err := deeplProvider.GetUsage(ctx)
		if err == nil && usageAfter != nil {
			charsUsed := usageAfter.CharacterCount
			if usageBefore != nil {
				charsUsed = usageAfter.CharacterCount - usageBefore.CharacterCount
			}
			fmt.Printf("\nAPI Usage:\n")
			fmt.Printf("  Characters used this run: %s\n", formatNumber(charsUsed))
			fmt.Printf("  Total used this period:   %s / %s (%.1f%%)\n",
				formatNumber(usageAfter.CharacterCount),
				formatNumber(usageAfter.CharacterLimit),
				float64(usageAfter.CharacterCount)/float64(usageAfter.CharacterLimit)*100)
			fmt.Printf("  Remaining:                %s characters\n",
				formatNumber(usageAfter.CharacterLimit-usageAfter.CharacterCount))
		}
	}
}

func runListLanguages() {
	provider, err := getProvider()
	if err != nil {
		// Show DeepL languages even without API key
		fmt.Println("Supported languages (DeepL):")
		langs := []string{"en", "de", "fr", "es", "it", "nl", "pl", "pt", "ru", "ja", "zh", "ko", "vi", "id", "tr", "uk", "cs", "da", "fi", "el", "hu", "lt", "lv", "nb", "ro", "sk", "sl", "sv", "bg", "et"}
		for _, l := range langs {
			fmt.Printf("  %s\n", l)
		}
		return
	}

	fmt.Printf("Supported languages (%s):\n", provider.Name())
	for _, lang := range provider.SupportedLanguages() {
		fmt.Printf("  %s\n", lang)
	}
}

func runStatus() {
	fmt.Println("========================================")
	fmt.Println("Translation Provider Status")
	fmt.Println("========================================")
	fmt.Println()

	// Show DeepL status
	deeplKey := os.Getenv("DEEPL_API_KEY")
	if deeplKey != "" {
		isFree := strings.HasSuffix(deeplKey, ":fx")
		fmt.Println("--- DeepL ---")
		fmt.Println("Status: Configured")
		if isFree {
			fmt.Println("Plan: Free (500k chars/month)")
		} else {
			fmt.Println("Plan: Pro")
		}
		fmt.Printf("API Key: %s...%s\n", deeplKey[:8], deeplKey[len(deeplKey)-4:])

		// Get actual usage from API
		provider, err := autotranslate.NewDeepLProvider(deeplKey)
		if err == nil {
			usage, err := provider.GetUsage(context.Background())
			if err == nil {
				fmt.Printf("Usage: %s / %s (%.1f%%)\n",
					formatNumber(usage.CharacterCount),
					formatNumber(usage.CharacterLimit),
					float64(usage.CharacterCount)/float64(usage.CharacterLimit)*100)
			}
		}
		fmt.Println()
	} else {
		fmt.Println("--- DeepL ---")
		fmt.Println("Status: Not configured")
		fmt.Println("Set DEEPL_API_KEY for 500k free chars/month")
		fmt.Println()
	}

	// Show Claude API status
	claudeKey := os.Getenv("CLAUDE_API_KEY")
	if claudeKey != "" {
		fmt.Println("--- Claude (API) ---")
		fmt.Println("Status: Configured")
		fmt.Println("Plan: Requires API credits")
		if len(claudeKey) > 12 {
			fmt.Printf("API Key: %s...%s\n", claudeKey[:8], claudeKey[len(claudeKey)-4:])
		}
		fmt.Println()
	} else {
		fmt.Println("--- Claude (API) ---")
		fmt.Println("Status: Not configured")
		fmt.Println("Set CLAUDE_API_KEY for API-based translation")
		fmt.Println()
	}

	// Show Claude CLI status
	fmt.Println("--- Claude (CLI) ---")
	if _, err := autotranslate.NewClaudeCLIProvider(); err == nil {
		fmt.Println("Status: Available")
		fmt.Println("Plan: Uses logged-in session (your subscription)")
		fmt.Println("Binary: claude")
	} else {
		fmt.Println("Status: Not available")
		fmt.Println("Install: bun add -g @anthropic-ai/claude-code")
	}
	fmt.Println()

	fmt.Println("========================================")
	fmt.Printf("Current provider: %s\n", *providerName)
	fmt.Println("Use --provider=deepl, --provider=claude, or --provider=claude-cli")
	fmt.Println("========================================")
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

func runARBTranslation(targetLang string) {
	ctx := context.Background()

	// Find the tokibundle directory
	bundle := *bundlePath
	if _, err := os.Stat(bundle); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: tokibundle directory not found at %s\n", bundle)
		fmt.Fprintf(os.Stderr, "Run 'toki generate' first to create the ARB catalogs\n")
		os.Exit(1)
	}

	// Try both naming conventions: catalog_en.arb (toki) and en.arb
	sourceARBPath := filepath.Join(bundle, "catalog_en.arb")
	if _, err := os.Stat(sourceARBPath); os.IsNotExist(err) {
		sourceARBPath = filepath.Join(bundle, "en.arb")
	}

	// Load source ARB (English)
	sourceARB, err := autotranslate.LoadARB(sourceARBPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading source ARB: %v\n", err)
		os.Exit(1)
	}

	// Load or create target ARB (try catalog_xx.arb first, then xx.arb)
	targetARBPath := filepath.Join(bundle, "catalog_"+targetLang+".arb")
	if _, err := os.Stat(targetARBPath); os.IsNotExist(err) {
		// Try without catalog_ prefix
		altPath := filepath.Join(bundle, targetLang+".arb")
		if _, err := os.Stat(altPath); err == nil {
			targetARBPath = altPath
		}
		// Otherwise use catalog_ prefix for new files
	}
	var targetARB *autotranslate.ARBFile
	if _, err := os.Stat(targetARBPath); os.IsNotExist(err) {
		// Create new target ARB with same structure but empty messages
		targetARB = &autotranslate.ARBFile{
			Locale:           targetLang,
			Messages:         make(map[string]string),
			Metadata:         make(map[string]any),
			CustomAttributes: sourceARB.CustomAttributes,
		}
		// Copy message IDs with empty values
		for id := range sourceARB.Messages {
			targetARB.Messages[id] = ""
		}
		// Copy metadata
		for k, v := range sourceARB.Metadata {
			targetARB.Metadata[k] = v
		}
	} else {
		targetARB, err = autotranslate.LoadARB(targetARBPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading target ARB: %v\n", err)
			os.Exit(1)
		}
		// Add any new message IDs from source
		for id := range sourceARB.Messages {
			if _, exists := targetARB.Messages[id]; !exists {
				targetARB.Messages[id] = ""
			}
		}
	}

	// Count what needs translation
	sourceStats := autotranslate.GetARBStats(sourceARB)
	targetStats := autotranslate.GetARBStats(targetARB)

	fmt.Printf("ARB Translation: en → %s\n", targetLang)
	fmt.Printf("Source (en):     %d messages\n", sourceStats.TotalMessages)
	fmt.Printf("Target (%s):    %d translated, %d empty (%.1f%% complete)\n",
		targetLang, targetStats.TranslatedCount, targetStats.EmptyCount, targetStats.CompletenessPerc)

	if targetStats.EmptyCount == 0 {
		fmt.Printf("\n✓ All messages already translated!\n")
		return
	}

	if *dryRun {
		fmt.Printf("\n[dry-run] Would translate %d messages\n", targetStats.EmptyCount)
		return
	}

	// Get provider
	provider, err := getProvider()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if !provider.SupportsLanguage(targetLang) {
		fmt.Fprintf(os.Stderr, "Error: language '%s' not supported by %s\n", targetLang, provider.Name())
		os.Exit(1)
	}

	fmt.Printf("\nTranslating %d messages using %s...\n\n", targetStats.EmptyCount, provider.Name())

	// Translate using the ARB translator
	translator := autotranslate.NewARBTranslator(provider)
	translated, err := translator.TranslateARB(ctx, sourceARB, targetARB, targetLang, *verbose,
		func(done, total int) {
			if !*verbose {
				fmt.Printf("\r  Progress: %d/%d messages (%.0f%%)", done, total, float64(done)/float64(total)*100)
			}
		})

	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError during translation: %v\n", err)
		// Save partial progress
		if translated > 0 {
			fmt.Printf("Saving %d translated messages...\n", translated)
			if saveErr := autotranslate.SaveARB(targetARBPath, targetARB); saveErr != nil {
				fmt.Fprintf(os.Stderr, "Error saving ARB: %v\n", saveErr)
			}
		}
		os.Exit(1)
	}

	fmt.Printf("\n\n")

	// Save the translated ARB
	if err := autotranslate.SaveARB(targetARBPath, targetARB); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving ARB: %v\n", err)
		os.Exit(1)
	}

	finalStats := autotranslate.GetARBStats(targetARB)
	fmt.Printf("✓ Translated %d messages\n", translated)
	fmt.Printf("✓ Saved to %s\n", targetARBPath)
	fmt.Printf("  Completeness: %.1f%% (%d/%d)\n",
		finalStats.CompletenessPerc, finalStats.TranslatedCount, finalStats.TotalMessages)
	fmt.Printf("\nNext step: run 'toki apply -t %s' to apply translations to markdown files\n", targetLang)
}

func runARBStatus() {
	bundle := *bundlePath
	if _, err := os.Stat(bundle); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: tokibundle directory not found at %s\n", bundle)
		fmt.Fprintf(os.Stderr, "Run 'toki generate' first to create the ARB catalogs\n")
		os.Exit(1)
	}

	fmt.Println("========================================")
	fmt.Println("ARB Translation Status")
	fmt.Printf("Bundle: %s\n", bundle)
	fmt.Println("========================================")
	fmt.Println()

	// Find all ARB files
	entries, err := os.ReadDir(bundle)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading bundle directory: %v\n", err)
		os.Exit(1)
	}

	var sourceARB *autotranslate.ARBFile
	type langStatus struct {
		code  string
		stats autotranslate.ARBStats
	}
	var statuses []langStatus

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".arb") {
			continue
		}

		arbPath := filepath.Join(bundle, entry.Name())
		arb, err := autotranslate.LoadARB(arbPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not load %s: %v\n", entry.Name(), err)
			continue
		}

		// Handle both naming conventions: catalog_xx.arb and xx.arb
		langCode := strings.TrimSuffix(entry.Name(), ".arb")
		langCode = strings.TrimPrefix(langCode, "catalog_")
		stats := autotranslate.GetARBStats(arb)

		if langCode == "en" {
			sourceARB = arb
		}

		statuses = append(statuses, langStatus{code: langCode, stats: stats})
	}

	if sourceARB == nil {
		fmt.Fprintf(os.Stderr, "Error: no catalog_en.arb or en.arb found in bundle\n")
		os.Exit(1)
	}

	// Print status table
	fmt.Printf("%-10s %8s %8s %8s %10s\n", "Language", "Total", "Done", "Empty", "Complete")
	fmt.Printf("%-10s %8s %8s %8s %10s\n", "--------", "-----", "----", "-----", "--------")

	for _, s := range statuses {
		bar := ""
		pct := s.stats.CompletenessPerc
		switch {
		case pct == 100:
			bar = "✓ Complete"
		case pct >= 75:
			bar = fmt.Sprintf("%.0f%%", pct)
		case pct >= 50:
			bar = fmt.Sprintf("%.0f%%", pct)
		case pct > 0:
			bar = fmt.Sprintf("%.0f%%", pct)
		default:
			bar = "Not started"
		}

		fmt.Printf("%-10s %8d %8d %8d %10s\n",
			s.code,
			s.stats.TotalMessages,
			s.stats.TranslatedCount,
			s.stats.EmptyCount,
			bar)
	}

	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  autotranslate arb <lang>      Translate empty entries for language")
	fmt.Println("  toki apply -t <lang>          Apply translations to markdown")
}

// sourceToTargetPath converts English source path to target language path
func sourceToTargetPath(sourcePath, targetLang string) string {
	// Load Hugo config to get target directory name
	config := translator.DefaultConfig()
	translator.TryLoadHugoConfig(config)

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
