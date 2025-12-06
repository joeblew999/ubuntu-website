package translator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// Checker handles translation status checking operations (no API key needed)
type Checker struct {
	config *Config
	git    *GitManager
}

// NewChecker creates a new Checker instance for status commands
func NewChecker() (*Checker, error) {
	config := DefaultConfig()

	git, err := NewGitManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create Git manager: %w", err)
	}

	return &Checker{
		config: config,
		git:    git,
	}, nil
}

// sourcePath returns the full path to the source content directory
func (c *Checker) sourcePath() string {
	return filepath.Join(c.config.ContentDir, c.config.SourceDir)
}

// Status shows what English files changed since last translation
// Returns exit code: 0 = no changes, 1 = changes found (in github-issue mode)
func (c *Checker) Status(githubIssue bool) int {
	newFiles := c.getNewFiles()
	uncommittedChanges := c.getUncommittedChanges()
	committedChanges := c.getCommittedChanges()

	hasChanges := len(newFiles) > 0 || len(uncommittedChanges) > 0 || len(committedChanges) > 0

	if githubIssue {
		if !hasChanges {
			return 0
		}
		fmt.Println("## Translation Status")
		fmt.Println()
		if len(newFiles) > 0 {
			fmt.Println("### New (untracked) files")
			for _, f := range newFiles {
				fmt.Printf("- `%s`\n", f)
			}
			fmt.Println()
		}
		if len(uncommittedChanges) > 0 {
			fmt.Println("### Uncommitted changes")
			for _, f := range uncommittedChanges {
				fmt.Printf("- `%s`\n", f)
			}
			fmt.Println()
		}
		if len(committedChanges) > 0 {
			fmt.Println("### Committed since last translation")
			for _, f := range committedChanges {
				fmt.Printf("- `%s`\n", f)
			}
			fmt.Println()
		}
		return 1
	}

	// Terminal output
	fmt.Println("========================================")
	fmt.Println("Translation Status")
	fmt.Println("========================================")
	fmt.Println()

	fmt.Println("=== New (untracked) files ===")
	if len(newFiles) > 0 {
		for _, f := range newFiles {
			fmt.Println(f)
		}
	} else {
		fmt.Println("(none)")
	}
	fmt.Println()

	fmt.Println("=== Uncommitted changes (modified) ===")
	if len(uncommittedChanges) > 0 {
		for _, f := range uncommittedChanges {
			fmt.Println(f)
		}
	} else {
		fmt.Println("(none)")
	}
	fmt.Println()

	fmt.Println("=== Committed since last translation ===")
	if len(committedChanges) > 0 {
		for _, f := range committedChanges {
			fmt.Println(f)
		}
	} else {
		if c.checkpointExists() {
			fmt.Println("(none)")
		} else {
			fmt.Println("(No checkpoint tag yet - run 'translate done' to set baseline)")
		}
	}
	fmt.Println()

	fmt.Println("========================================")
	fmt.Println("To translate, ask Claude Code:")
	fmt.Println("  'Translate the changed files to all languages'")
	fmt.Println()
	fmt.Println("After translating: translate done")
	fmt.Println("========================================")

	return 0
}

// Diff shows git diff for a specific English file since last translation
func (c *Checker) Diff(file string) int {
	sourcePath := c.sourcePath()

	// Handle both formats: "blog/post.md" or "content/english/blog/post.md"
	var enFile, relPath string
	if strings.HasPrefix(file, sourcePath+"/") {
		enFile = file
		relPath = strings.TrimPrefix(file, sourcePath+"/")
	} else {
		enFile = filepath.Join(sourcePath, file)
		relPath = file
	}

	if _, err := os.Stat(enFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "ERROR: File not found: %s\n", enFile)
		return 1
	}

	fmt.Println("========================================")
	fmt.Printf("Diff for: %s\n", relPath)
	fmt.Println("========================================")
	fmt.Println()

	// Check if file is new (not in last-translation)
	checkCmd := exec.Command("git", "show", c.config.CheckpointTag+":"+enFile)
	if err := checkCmd.Run(); err != nil {
		// File doesn't exist in checkpoint - it's new
		fmt.Println("STATUS: NEW FILE (did not exist at last translation checkpoint)")
		fmt.Println()
		fmt.Println("Full content:")
		fmt.Println("----------------------------------------")
		content, _ := os.ReadFile(enFile)
		fmt.Print(string(content))
		fmt.Println("----------------------------------------")
	} else {
		// Show diff since last translation (committed changes)
		diffCmd := exec.Command("git", "diff", c.config.CheckpointTag+"..HEAD", "--", enFile)
		committedDiff, _ := diffCmd.Output()

		// Also check for uncommitted changes (working directory vs HEAD)
		uncommittedCmd := exec.Command("git", "diff", "HEAD", "--", enFile)
		uncommittedDiff, _ := uncommittedCmd.Output()

		// And check for staged but uncommitted
		stagedCmd := exec.Command("git", "diff", "--cached", "--", enFile)
		stagedDiff, _ := stagedCmd.Output()

		hasCommitted := len(committedDiff) > 0
		hasUncommitted := len(uncommittedDiff) > 0 || len(stagedDiff) > 0

		if !hasCommitted && !hasUncommitted {
			fmt.Println("STATUS: NO CHANGES since last translation")
		} else {
			if hasCommitted {
				fmt.Println("STATUS: MODIFIED since last translation (committed)")
				fmt.Println()
				fmt.Println("Committed changes:")
				fmt.Println("----------------------------------------")
				fmt.Print(string(committedDiff))
				fmt.Println("----------------------------------------")
			}
			if hasUncommitted {
				if hasCommitted {
					fmt.Println()
				}
				fmt.Println("STATUS: UNCOMMITTED CHANGES (not yet committed)")
				fmt.Println()
				fmt.Println("Uncommitted changes:")
				fmt.Println("----------------------------------------")
				if len(stagedDiff) > 0 {
					fmt.Print(string(stagedDiff))
				}
				if len(uncommittedDiff) > 0 {
					fmt.Print(string(uncommittedDiff))
				}
				fmt.Println("----------------------------------------")
			}
		}
	}

	fmt.Println()
	fmt.Println("========================================")
	return 0
}

// Missing shows which languages are missing content files compared to English
func (c *Checker) Missing(githubIssue bool) int {
	englishFiles := c.getEnglishFiles()
	totalMissing := 0
	missingByLang := make(map[string][]string)

	for _, enFile := range englishFiles {
		relPath := strings.TrimPrefix(enFile, c.sourcePath()+string(os.PathSeparator))
		for _, lang := range c.config.TargetLangs {
			langFile := filepath.Join(c.config.ContentDir, lang.DirName, relPath)
			if _, err := os.Stat(langFile); os.IsNotExist(err) {
				missingByLang[lang.Code] = append(missingByLang[lang.Code], relPath)
				totalMissing++
			}
		}
	}

	if githubIssue {
		if totalMissing == 0 {
			return 0
		}
		fmt.Println("## Missing Translations")
		fmt.Println()
		for _, lang := range c.config.TargetLangs {
			if files, ok := missingByLang[lang.Code]; ok && len(files) > 0 {
				fmt.Printf("### %s (%d files)\n", lang.Name, len(files))
				for _, f := range files {
					fmt.Printf("- `%s`\n", f)
				}
				fmt.Println()
			}
		}
		return 1
	}

	// Terminal output
	fmt.Println("========================================")
	fmt.Println("Missing Content Files by Language")
	fmt.Println("========================================")
	fmt.Println()

	for _, lang := range c.config.TargetLangs {
		if files, ok := missingByLang[lang.Code]; ok && len(files) > 0 {
			fmt.Printf("MISSING: %s: Missing %d files\n", lang.Name, len(files))
			for _, f := range files {
				fmt.Printf("  - %s\n", f)
			}
			fmt.Println()
		} else {
			fmt.Printf("OK: %s: Complete\n", lang.Name)
		}
	}

	fmt.Println("========================================")
	return 0
}

// Stale shows target files that are smaller than English (may need updating)
func (c *Checker) Stale(githubIssue bool) int {
	englishFiles := c.getEnglishFiles()
	var staleFiles []string

	for _, enFile := range englishFiles {
		enInfo, err := os.Stat(enFile)
		if err != nil {
			continue
		}
		enSize := enInfo.Size()
		if enSize <= 500 {
			continue // Skip small files
		}

		relPath := strings.TrimPrefix(enFile, c.sourcePath()+string(os.PathSeparator))
		threshold := enSize / 2

		for _, lang := range c.config.TargetLangs {
			langFile := filepath.Join(c.config.ContentDir, lang.DirName, relPath)
			langInfo, err := os.Stat(langFile)
			if err != nil {
				continue // File doesn't exist
			}
			if langInfo.Size() < threshold {
				staleFiles = append(staleFiles, fmt.Sprintf("%s (English: %d bytes, Target: %d bytes)", langFile, enSize, langInfo.Size()))
			}
		}
	}

	if githubIssue {
		if len(staleFiles) == 0 {
			return 0
		}
		fmt.Println("## Potentially Stale Translations")
		fmt.Println()
		fmt.Println("These files are less than 50% the size of the English source:")
		fmt.Println()
		for _, f := range staleFiles {
			fmt.Printf("- `%s`\n", f)
		}
		return 1
	}

	// Terminal output
	fmt.Println("========================================")
	fmt.Println("Potentially Stale Translations")
	fmt.Println("(target file exists but is much smaller than English)")
	fmt.Println("========================================")
	fmt.Println()

	if len(staleFiles) == 0 {
		fmt.Println("OK: No stale translations found")
	} else {
		for _, f := range staleFiles {
			fmt.Printf("STALE: %s\n", f)
		}
		fmt.Println()
		fmt.Printf("Found %d potentially stale files\n", len(staleFiles))
		fmt.Println("Review and re-translate if needed")
	}
	fmt.Println("========================================")
	return 0
}

// Orphans shows files in target languages that no longer exist in English
func (c *Checker) Orphans(githubIssue bool) int {
	orphansByLang := make(map[string][]string)
	totalOrphans := 0

	for _, lang := range c.config.TargetLangs {
		langPath := filepath.Join(c.config.ContentDir, lang.DirName)
		filepath.Walk(langPath, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() || !strings.HasSuffix(path, ".md") {
				return nil
			}
			// Convert to English path
			relPath := strings.TrimPrefix(path, langPath+string(os.PathSeparator))
			enFile := filepath.Join(c.sourcePath(), relPath)
			if _, err := os.Stat(enFile); os.IsNotExist(err) {
				orphansByLang[lang.Code] = append(orphansByLang[lang.Code], path)
				totalOrphans++
			}
			return nil
		})
	}

	if githubIssue {
		if totalOrphans == 0 {
			return 0
		}
		fmt.Println("## Orphaned Translation Files")
		fmt.Println()
		fmt.Println("These files exist in target languages but not in English (should be deleted):")
		fmt.Println()
		for _, lang := range c.config.TargetLangs {
			if files, ok := orphansByLang[lang.Code]; ok && len(files) > 0 {
				fmt.Printf("### %s\n", lang.Name)
				for _, f := range files {
					fmt.Printf("- `%s`\n", f)
				}
				fmt.Println()
			}
		}
		return 1
	}

	// Terminal output
	fmt.Println("========================================")
	fmt.Println("Orphaned Files (exist in target but not in English)")
	fmt.Println("========================================")
	fmt.Println()

	for _, lang := range c.config.TargetLangs {
		if files, ok := orphansByLang[lang.Code]; ok && len(files) > 0 {
			fmt.Printf("ORPHANS: %s: %d orphaned files (DELETE THESE)\n", lang.Name, len(files))
			for _, f := range files {
				fmt.Printf("  - %s\n", f)
			}
			fmt.Println()
		} else {
			fmt.Printf("OK: %s: No orphans\n", lang.Name)
		}
	}

	fmt.Println("========================================")
	if totalOrphans > 0 {
		fmt.Println("Run 'translate clean' to delete all orphaned files")
	}
	fmt.Println("========================================")
	return 0
}

// Clean deletes orphaned files in target languages
func (c *Checker) Clean() int {
	deleted := 0

	for _, lang := range c.config.TargetLangs {
		langPath := filepath.Join(c.config.ContentDir, lang.DirName)
		filepath.Walk(langPath, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() || !strings.HasSuffix(path, ".md") {
				return nil
			}
			relPath := strings.TrimPrefix(path, langPath+string(os.PathSeparator))
			enFile := filepath.Join(c.sourcePath(), relPath)
			if _, err := os.Stat(enFile); os.IsNotExist(err) {
				fmt.Printf("Deleting: %s\n", path)
				os.Remove(path)
				deleted++
			}
			return nil
		})
	}

	fmt.Printf("OK: Deleted %d orphaned files\n", deleted)
	return 0
}

// Done marks current translations as complete (update checkpoint tag)
func (c *Checker) Done() int {
	cmd := exec.Command("git", "tag", "-f", c.config.CheckpointTag, "HEAD")
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error updating checkpoint: %v\n", err)
		return 1
	}
	fmt.Printf("OK: Translation checkpoint updated to current commit\n")
	return 0
}

// Next shows the next file to translate with progress
func (c *Checker) Next() int {
	englishFiles := c.getEnglishFiles()
	sort.Strings(englishFiles)

	// Count total missing
	totalMissing := 0
	for _, enFile := range englishFiles {
		relPath := strings.TrimPrefix(enFile, c.sourcePath()+string(os.PathSeparator))
		for _, lang := range c.config.TargetLangs {
			langFile := filepath.Join(c.config.ContentDir, lang.DirName, relPath)
			if _, err := os.Stat(langFile); os.IsNotExist(err) {
				totalMissing++
			}
		}
	}

	if totalMissing == 0 {
		fmt.Println("========================================")
		fmt.Println("All files translated!")
		fmt.Println("========================================")
		return 0
	}

	// Find first missing
	for _, enFile := range englishFiles {
		relPath := strings.TrimPrefix(enFile, c.sourcePath()+string(os.PathSeparator))
		var missingIn []string

		for _, lang := range c.config.TargetLangs {
			langFile := filepath.Join(c.config.ContentDir, lang.DirName, relPath)
			if _, err := os.Stat(langFile); os.IsNotExist(err) {
				missingIn = append(missingIn, lang.DirName)
			}
		}

		if len(missingIn) > 0 {
			// Calculate completed translations
			totalPossible := len(englishFiles) * len(c.config.TargetLangs)
			completed := totalPossible - totalMissing

			fmt.Println("========================================")
			fmt.Printf("Progress: %d/%d translations complete (%d remaining)\n", completed, totalPossible, totalMissing)
			fmt.Println()
			fmt.Println("Next file to translate:")
			fmt.Printf("  %s\n", relPath)
			fmt.Println()
			fmt.Printf("Missing in: %s\n", strings.Join(missingIn, " "))
			fmt.Println()
			fmt.Println("To translate, ask Claude Code:")
			fmt.Printf("  'Translate %s to all languages'\n", relPath)
			fmt.Println("========================================")
			return 0
		}
	}

	return 0
}

// Changed shows detailed changes for all English files since last translation
func (c *Checker) Changed() int {
	changedFiles := c.getCommittedChanges()

	fmt.Println("========================================")
	fmt.Println("Detailed Changes Since Last Translation")
	fmt.Println("========================================")
	fmt.Println()

	if len(changedFiles) == 0 {
		fmt.Println("No English files changed since last translation.")
		fmt.Println("========================================")
		return 0
	}

	fmt.Printf("Found %d changed file(s):\n", len(changedFiles))
	fmt.Println()

	for _, file := range changedFiles {
		relPath := strings.TrimPrefix(file, c.sourcePath()+"/")
		fmt.Printf("--- %s ---\n", relPath)

		// Show summary stats
		statCmd := exec.Command("git", "diff", "--stat", c.config.CheckpointTag+"..HEAD", "--", file)
		statOut, _ := statCmd.Output()
		lines := strings.Split(string(statOut), "\n")
		if len(lines) > 1 {
			fmt.Printf("  %s\n", strings.TrimSpace(lines[len(lines)-2]))
		}

		// Show first few diff lines
		diffCmd := exec.Command("git", "diff", c.config.CheckpointTag+"..HEAD", "--", file)
		diffOut, _ := diffCmd.Output()
		diffLines := []string{}
		for _, line := range strings.Split(string(diffOut), "\n") {
			if (strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++")) ||
				(strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---")) {
				diffLines = append(diffLines, line)
			}
		}

		if len(diffLines) > 0 {
			fmt.Println("  Preview:")
			showLines := diffLines
			if len(showLines) > 10 {
				showLines = showLines[:10]
			}
			for _, l := range showLines {
				fmt.Printf("    %s\n", l)
			}
			if len(diffLines) > 10 {
				fmt.Printf("    ... and %d more lines\n", len(diffLines)-10)
			}
		}
		fmt.Println()
	}

	fmt.Println("========================================")
	fmt.Println("To see full diff for a file:")
	fmt.Println("  translate diff <path>")
	fmt.Println("========================================")
	return 0
}

// Helper methods

func (c *Checker) getNewFiles() []string {
	cmd := exec.Command("git", "ls-files", "--others", "--exclude-standard", "--", c.sourcePath()+"/")
	output, _ := cmd.Output()
	var files []string
	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && strings.HasSuffix(line, ".md") {
			files = append(files, line)
		}
	}
	return files
}

func (c *Checker) getUncommittedChanges() []string {
	cmd := exec.Command("git", "diff", "--name-only", "--", c.sourcePath()+"/", "config/_default/menus.en.toml", "i18n/en.yaml")
	output, _ := cmd.Output()
	var files []string
	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && (strings.HasSuffix(line, ".md") || strings.HasSuffix(line, ".toml") || strings.HasSuffix(line, ".yaml")) {
			files = append(files, line)
		}
	}
	return files
}

func (c *Checker) getCommittedChanges() []string {
	if !c.checkpointExists() {
		return nil
	}
	cmd := exec.Command("git", "diff", "--name-only", c.config.CheckpointTag+"..HEAD", "--", c.sourcePath()+"/", "config/_default/menus.en.toml", "i18n/en.yaml")
	output, _ := cmd.Output()
	var files []string
	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && (strings.HasSuffix(line, ".md") || strings.HasSuffix(line, ".toml") || strings.HasSuffix(line, ".yaml")) {
			files = append(files, line)
		}
	}
	return files
}

func (c *Checker) checkpointExists() bool {
	cmd := exec.Command("git", "tag", "-l", c.config.CheckpointTag)
	output, _ := cmd.Output()
	return strings.TrimSpace(string(output)) != ""
}

func (c *Checker) getEnglishFiles() []string {
	var files []string
	filepath.Walk(c.sourcePath(), func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}
		files = append(files, path)
		return nil
	})
	return files
}

// Validate checks if the translator config matches Hugo's language config
func (c *Checker) Validate() int {
	fmt.Println("========================================")
	fmt.Println("Validating Translator Configuration")
	fmt.Println("========================================")
	fmt.Println()

	// Check if this is a Hugo project
	if !IsHugoProject() {
		fmt.Println("Mode: Standalone (no Hugo config found)")
		fmt.Println()
		fmt.Println("Current configuration:")
		fmt.Printf("  Source: %s → content/%s\n", c.config.SourceLang, c.config.SourceDir)
		for _, lang := range c.config.TargetLangs {
			fmt.Printf("  Target: %s (%s) → content/%s\n", lang.Code, lang.Name, lang.DirName)
		}
		fmt.Println()
		fmt.Println("========================================")
		fmt.Println("OK: Using default configuration")
		fmt.Println("========================================")
		return 0
	}

	fmt.Println("Mode: Hugo project detected")
	fmt.Println()

	// Show current config (auto-loaded from Hugo)
	fmt.Printf("Source: %s → content/%s\n", c.config.SourceLang, c.config.SourceDir)
	for _, lang := range c.config.TargetLangs {
		fmt.Printf("Target: %s (%s) → content/%s\n", lang.Code, lang.Name, lang.DirName)
	}
	fmt.Println()

	// Validate against Hugo config
	mismatches := ValidateHugoConfig(c.config)

	if len(mismatches) > 0 {
		fmt.Println("========================================")
		fmt.Printf("WARNING: %d mismatch(es) found\n", len(mismatches))
		fmt.Println("========================================")
		for _, m := range mismatches {
			fmt.Printf("  • %s\n", m)
		}
		fmt.Println()
		fmt.Println("This shouldn't happen - languages are auto-loaded from Hugo config.")
		fmt.Println("Check if config/_default/languages.toml changed after binary was built.")
		return 1
	}

	fmt.Println("========================================")
	fmt.Println("OK: Configuration loaded from Hugo")
	fmt.Println("========================================")
	return 0
}

// Langs shows all configured languages and detects stray content directories
func (c *Checker) Langs() int {
	fmt.Println("========================================")
	fmt.Println("Language Configuration")
	fmt.Println("========================================")
	fmt.Println()

	// Build set of known directories
	knownDirs := make(map[string]bool)
	knownDirs[c.config.SourceDir] = true
	for _, lang := range c.config.TargetLangs {
		knownDirs[lang.DirName] = true
	}

	// Show source language
	fmt.Printf("SOURCE: %s → content/%s/\n", c.config.SourceLang, c.config.SourceDir)
	fmt.Println()

	// Show target languages
	fmt.Println("TARGETS:")
	for _, lang := range c.config.TargetLangs {
		dirExists := "✓"
		langPath := filepath.Join(c.config.ContentDir, lang.DirName)
		if _, err := os.Stat(langPath); os.IsNotExist(err) {
			dirExists = "✗ (directory missing)"
		}
		fmt.Printf("  %s (%s) → content/%s/ %s\n", lang.Code, lang.Name, lang.DirName, dirExists)
	}
	fmt.Println()

	// Scan for stray directories (exist in content/ but not in config)
	strayDirs := []string{}
	entries, err := os.ReadDir(c.config.ContentDir)
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() && !knownDirs[entry.Name()] {
				strayDirs = append(strayDirs, entry.Name())
			}
		}
	}

	if len(strayDirs) > 0 {
		fmt.Println("WARNING: Stray directories (not in config):")
		for _, dir := range strayDirs {
			// Count files in stray directory
			count := 0
			filepath.Walk(filepath.Join(c.config.ContentDir, dir), func(path string, info os.FileInfo, err error) error {
				if err == nil && !info.IsDir() && strings.HasSuffix(path, ".md") {
					count++
				}
				return nil
			})
			fmt.Printf("  content/%s/ (%d .md files)\n", dir, count)
		}
		fmt.Println()
		fmt.Println("These directories may be from a removed language.")
		fmt.Println("If they should be deleted, remove them manually:")
		for _, dir := range strayDirs {
			fmt.Printf("  rm -rf content/%s/\n", dir)
		}
		fmt.Println()
		fmt.Println("========================================")
		fmt.Println("ACTION NEEDED: Stray directories found")
		fmt.Println("========================================")
		return 1
	}

	fmt.Println("========================================")
	fmt.Println("OK: All content directories are configured")
	fmt.Println("========================================")
	return 0
}
