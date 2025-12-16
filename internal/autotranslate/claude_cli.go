package autotranslate

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ClaudeCLIProvider implements Provider using the Claude CLI (uses logged-in session)
type ClaudeCLIProvider struct {
	cliBinary string
}

// NewClaudeCLIProvider creates a provider that uses the Claude CLI
func NewClaudeCLIProvider() (*ClaudeCLIProvider, error) {
	// Find claude binary - check multiple locations
	binary, err := exec.LookPath("claude")
	if err != nil {
		// Check common bun/npm global locations
		candidates := []string{
			os.ExpandEnv("$HOME/.bun/bin/claude"),
			os.ExpandEnv("$HOME/.local/bin/claude"),
			"/usr/local/bin/claude",
		}
		for _, candidate := range candidates {
			if _, statErr := os.Stat(candidate); statErr == nil {
				binary = candidate
				err = nil
				break
			}
		}
	}
	if err != nil {
		return nil, fmt.Errorf("claude CLI not found (install: bun add -g @anthropic-ai/claude-code)")
	}

	return &ClaudeCLIProvider{
		cliBinary: binary,
	}, nil
}

// Name returns the provider name
func (p *ClaudeCLIProvider) Name() string {
	return "claude-cli"
}

// Translate translates text using Claude CLI
func (p *ClaudeCLIProvider) Translate(ctx context.Context, text, sourceLang, targetLang string) (string, error) {
	targetName := claudeLangNames[strings.ToLower(targetLang)]
	if targetName == "" {
		targetName = targetLang
	}

	// Handle empty text
	if strings.TrimSpace(text) == "" {
		return "", nil
	}

	prompt := fmt.Sprintf(`Translate the following text from English to %s.

CRITICAL INSTRUCTIONS:
1. Return ONLY the translated text - no explanations, notes, or commentary
2. Preserve ALL markdown formatting exactly (headers, lists, bold, italic, links)
3. DO NOT translate:
   - URLs and links (keep them exactly as-is)
   - Code blocks (content between triple backticks)
   - Inline code (content between single backticks)
   - Hugo shortcodes (like {{< shortcode >}} and {{%% shortcode %%}})
   - HTML tags
4. Maintain the same paragraph structure and line breaks

Text to translate:
%s`, targetName, text)

	return p.callCLI(ctx, prompt)
}

// TranslateBatch translates multiple texts
func (p *ClaudeCLIProvider) TranslateBatch(ctx context.Context, texts []string, sourceLang, targetLang string) ([]string, error) {
	// For CLI, translate one at a time (CLI doesn't support batching well)
	results := make([]string, len(texts))
	for i, text := range texts {
		translated, err := p.Translate(ctx, text, sourceLang, targetLang)
		if err != nil {
			return nil, fmt.Errorf("translating text %d: %w", i, err)
		}
		results[i] = translated
	}
	return results, nil
}

// SupportedLanguages returns list of supported language codes
func (p *ClaudeCLIProvider) SupportedLanguages() []string {
	langs := make([]string, 0, len(claudeLangNames))
	for code := range claudeLangNames {
		langs = append(langs, code)
	}
	return langs
}

// SupportsLanguage checks if a language code is supported
func (p *ClaudeCLIProvider) SupportsLanguage(langCode string) bool {
	_, ok := claudeLangNames[strings.ToLower(langCode)]
	return ok
}

// callCLI executes a prompt via the Claude CLI
func (p *ClaudeCLIProvider) callCLI(ctx context.Context, prompt string) (string, error) {
	// Use claude CLI with --print flag for non-interactive output
	cmd := exec.CommandContext(ctx, p.cliBinary, "--print", "-p", prompt)

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("claude CLI error: %s", string(exitErr.Stderr))
		}
		return "", fmt.Errorf("failed to run claude CLI: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}
