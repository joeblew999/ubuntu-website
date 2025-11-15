package translator

import (
	"fmt"
	"os"
	"path/filepath"
)

// Translator handles translation operations
type Translator struct {
	apiKey string
	config *Config
	claude *ClaudeClient
	git    *GitManager
}

// Config holds translation configuration
type Config struct {
	SourceLang    string
	TargetLangs   []string
	ContentDir    string
	I18nDir       string
	CheckpointTag string
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		SourceLang:  "en",
		TargetLangs: []string{"de", "sv", "zh", "ja", "th"},
		ContentDir:  "content",
		I18nDir:     "i18n",
		CheckpointTag: "last-translation",
	}
}

// New creates a new Translator instance
func New(apiKey string) (*Translator, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	config := DefaultConfig()

	// Create Claude client
	claude, err := NewClaudeClient(apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create Claude client: %w", err)
	}

	// Create Git manager
	git, err := NewGitManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create Git manager: %w", err)
	}

	return &Translator{
		apiKey: apiKey,
		config: config,
		claude: claude,
		git:    git,
	}, nil
}

// Check shows which English files have changed since last translation
func (t *Translator) Check() error {
	fmt.Println("🔍 Checking for changes since last translation...")

	changes, err := t.git.GetChangedFiles(t.config.CheckpointTag, "content/english")
	if err != nil {
		return fmt.Errorf("failed to get changed files: %w", err)
	}

	if len(changes) == 0 {
		fmt.Println("✅ No changes detected. All content is up to date!")
		return nil
	}

	fmt.Printf("\n📝 Found %d changed files:\n\n", len(changes))
	for _, file := range changes {
		fmt.Printf("  - %s\n", file)
	}
	fmt.Println()

	return nil
}

// TranslateAll translates all changed English content to all target languages
func (t *Translator) TranslateAll() error {
	changes, err := t.git.GetChangedFiles(t.config.CheckpointTag, "content/english")
	if err != nil {
		return fmt.Errorf("failed to get changed files: %w", err)
	}

	if len(changes) == 0 {
		fmt.Println("✅ No changes detected. Nothing to translate.")
		return nil
	}

	fmt.Printf("📝 Translating %d files to %d languages...\n\n", len(changes), len(t.config.TargetLangs))

	for _, targetLang := range t.config.TargetLangs {
		if err := t.translateFiles(changes, targetLang); err != nil {
			return fmt.Errorf("failed to translate to %s: %w", targetLang, err)
		}
	}

	// Update checkpoint
	if err := t.git.UpdateCheckpoint(t.config.CheckpointTag, changes); err != nil {
		return fmt.Errorf("failed to update checkpoint: %w", err)
	}

	return nil
}

// TranslateLang translates changed English content to a specific language
func (t *Translator) TranslateLang(targetLang string) error {
	// Validate target language
	validLang := false
	for _, lang := range t.config.TargetLangs {
		if lang == targetLang {
			validLang = true
			break
		}
	}
	if !validLang {
		return fmt.Errorf("invalid target language: %s (valid: %v)", targetLang, t.config.TargetLangs)
	}

	changes, err := t.git.GetChangedFiles(t.config.CheckpointTag, "content/english")
	if err != nil {
		return fmt.Errorf("failed to get changed files: %w", err)
	}

	if len(changes) == 0 {
		fmt.Println("✅ No changes detected. Nothing to translate.")
		return nil
	}

	return t.translateFiles(changes, targetLang)
}

// translateFiles translates a list of files to a target language
func (t *Translator) translateFiles(files []string, targetLang string) error {
	langName := t.getLanguageName(targetLang)
	targetDir := t.getContentDir(targetLang)

	fmt.Printf("🌍 Translating to %s (%s)...\n", langName, targetLang)

	for i, file := range files {
		fmt.Printf("  [%d/%d] %s\n", i+1, len(files), filepath.Base(file))

		// Read source file
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", file, err)
		}

		// Parse markdown with front matter
		md, err := ParseMarkdown(content)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", file, err)
		}

		// Translate the content (not front matter or code blocks)
		translatedBody, err := t.claude.Translate(md.Body, targetLang, langName)
		if err != nil {
			return fmt.Errorf("failed to translate %s: %w", file, err)
		}

		// Reconstruct markdown with translated content
		md.Body = translatedBody
		output, err := md.Reconstruct()
		if err != nil {
			return fmt.Errorf("failed to reconstruct %s: %w", file, err)
		}

		// Write to target language directory
		targetFile := filepath.Join(targetDir, filepath.Base(filepath.Dir(file)), filepath.Base(file))
		targetFileDir := filepath.Dir(targetFile)

		if err := os.MkdirAll(targetFileDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", targetFileDir, err)
		}

		if err := os.WriteFile(targetFile, []byte(output), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", targetFile, err)
		}
	}

	return nil
}

// TranslateI18n translates i18n TOML files
func (t *Translator) TranslateI18n() error {
	sourceFile := filepath.Join(t.config.I18nDir, "en.yaml")

	// Check if source file exists
	if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
		return fmt.Errorf("source i18n file not found: %s", sourceFile)
	}

	for _, targetLang := range t.config.TargetLangs {
		langName := t.getLanguageName(targetLang)
		fmt.Printf("🌍 Translating i18n to %s (%s)...\n", langName, targetLang)

		// Read source i18n file
		content, err := os.ReadFile(sourceFile)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", sourceFile, err)
		}

		// Parse i18n file
		i18nData, err := ParseI18n(content)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", sourceFile, err)
		}

		// Translate values
		translatedData, err := t.claude.TranslateI18n(i18nData, targetLang, langName)
		if err != nil {
			return fmt.Errorf("failed to translate i18n to %s: %w", targetLang, err)
		}

		// Write translated i18n file
		targetFile := filepath.Join(t.config.I18nDir, fmt.Sprintf("%s.yaml", targetLang))
		output, err := ReconstructI18n(translatedData)
		if err != nil {
			return fmt.Errorf("failed to reconstruct i18n: %w", err)
		}

		if err := os.WriteFile(targetFile, []byte(output), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", targetFile, err)
		}
	}

	return nil
}

// getLanguageName returns the full language name for a language code
func (t *Translator) getLanguageName(langCode string) string {
	names := map[string]string{
		"de": "German",
		"sv": "Swedish",
		"zh": "Simplified Chinese",
		"ja": "Japanese",
		"th": "Thai",
	}
	if name, ok := names[langCode]; ok {
		return name
	}
	return langCode
}

// getContentDir returns the content directory for a language code
func (t *Translator) getContentDir(langCode string) string {
	dirs := map[string]string{
		"de": "content/german",
		"sv": "content/swedish",
		"zh": "content/chinese",
		"ja": "content/japanese",
		"th": "content/thai",
	}
	if dir, ok := dirs[langCode]; ok {
		return dir
	}
	return fmt.Sprintf("content/%s", langCode)
}
