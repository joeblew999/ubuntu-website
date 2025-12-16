package autotranslate

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

// Placeholder markers for content that shouldn't be translated
const (
	placeholderPrefix = "[[NOTRANSLATE_"
	placeholderSuffix = "]]"
)

// MarkdownTranslator handles translation of Hugo markdown files
// while preserving front matter, shortcodes, code blocks, etc.
type MarkdownTranslator struct {
	provider Provider
}

// NewMarkdownTranslator creates a new markdown translator
func NewMarkdownTranslator(provider Provider) *MarkdownTranslator {
	return &MarkdownTranslator{provider: provider}
}

// TranslateFile translates a Hugo markdown file content
func (t *MarkdownTranslator) TranslateFile(ctx context.Context, content, sourceLang, targetLang string) (string, error) {
	// Split front matter from body
	frontMatter, body, err := splitFrontMatter(content)
	if err != nil {
		return "", err
	}

	// Translate front matter (only specific fields)
	translatedFrontMatter, err := t.translateFrontMatter(ctx, frontMatter, sourceLang, targetLang)
	if err != nil {
		return "", fmt.Errorf("translating front matter: %w", err)
	}

	// Translate body with placeholder protection
	translatedBody, err := t.translateBody(ctx, body, sourceLang, targetLang)
	if err != nil {
		return "", fmt.Errorf("translating body: %w", err)
	}

	// Reassemble
	return assembleFrontMatter(translatedFrontMatter, translatedBody), nil
}

// splitFrontMatter separates YAML front matter from markdown body
func splitFrontMatter(content string) (frontMatter, body string, err error) {
	content = strings.TrimSpace(content)

	// Check for front matter delimiter
	if !strings.HasPrefix(content, "---") {
		return "", content, nil
	}

	// Find closing delimiter
	rest := content[3:]
	idx := strings.Index(rest, "\n---")
	if idx == -1 {
		return "", content, nil
	}

	frontMatter = strings.TrimSpace(rest[:idx])
	body = strings.TrimSpace(rest[idx+4:])

	return frontMatter, body, nil
}

// assembleFrontMatter combines front matter and body
func assembleFrontMatter(frontMatter, body string) string {
	if frontMatter == "" {
		return body
	}
	return fmt.Sprintf("---\n%s\n---\n\n%s\n", frontMatter, body)
}

// Fields in front matter that should be translated
var translatableFrontMatterFields = map[string]bool{
	"title":       true,
	"meta_title":  true,
	"description": true,
	"excerpt":     true,
	"summary":     true,
}

// Fields that should never be translated
var preserveFrontMatterFields = map[string]bool{
	"date":       true,
	"draft":      true,
	"image":      true,
	"images":     true,
	"author":     true,
	"authors":    true,
	"slug":       true,
	"url":        true,
	"aliases":    true,
	"weight":     true,
	"categories": true,
	"tags":       true,
	"layout":     true,
	"type":       true,
}

// translateFrontMatter translates specific fields in YAML front matter
func (t *MarkdownTranslator) translateFrontMatter(ctx context.Context, frontMatter, sourceLang, targetLang string) (string, error) {
	if frontMatter == "" {
		return "", nil
	}

	lines := strings.Split(frontMatter, "\n")
	var result []string
	var textsToTranslate []string
	var lineIndices []int // Which lines need translation

	for i, line := range lines {
		// Skip empty lines and comments
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			result = append(result, line)
			continue
		}

		// Parse key: value
		colonIdx := strings.Index(line, ":")
		if colonIdx == -1 {
			result = append(result, line)
			continue
		}

		key := strings.TrimSpace(line[:colonIdx])
		value := strings.TrimSpace(line[colonIdx+1:])

		// Check if this field should be translated
		if translatableFrontMatterFields[key] && value != "" {
			// Remove quotes if present
			cleanValue := strings.Trim(value, "\"'")
			if cleanValue != "" {
				textsToTranslate = append(textsToTranslate, cleanValue)
				lineIndices = append(lineIndices, i)
			}
		}
		result = append(result, line)
	}

	// Batch translate all fields
	if len(textsToTranslate) > 0 {
		translations, err := t.provider.TranslateBatch(ctx, textsToTranslate, sourceLang, targetLang)
		if err != nil {
			return "", err
		}

		// Replace translated values
		for i, lineIdx := range lineIndices {
			line := result[lineIdx]
			colonIdx := strings.Index(line, ":")
			key := line[:colonIdx]
			indent := ""
			for _, ch := range line {
				if ch == ' ' || ch == '\t' {
					indent += string(ch)
				} else {
					break
				}
			}
			// Escape quotes in translated text
			escaped := strings.ReplaceAll(translations[i], "\"", "\\\"")
			result[lineIdx] = fmt.Sprintf("%s%s: \"%s\"", indent, strings.TrimSpace(key), escaped)
		}
	}

	return strings.Join(result, "\n"), nil
}

// Patterns for content that should not be translated
var (
	// Hugo shortcodes: {{< shortcode >}} or {{% shortcode %}}
	shortcodePattern = regexp.MustCompile(`\{\{[<%].*?[%>]\}\}`)

	// Code blocks: ```lang ... ```
	codeBlockPattern = regexp.MustCompile("(?s)```[a-z]*\\n.*?```")

	// Inline code: `code`
	inlineCodePattern = regexp.MustCompile("`[^`]+`")

	// URLs: [text](url) - protect the URL part
	urlPattern = regexp.MustCompile(`\]\([^)]+\)`)

	// HTML tags
	htmlTagPattern = regexp.MustCompile(`<[^>]+>`)

	// Image references: ![alt](path)
	imagePattern = regexp.MustCompile(`!\[[^\]]*\]\([^)]+\)`)

	// Reference links: [text][ref]
	refLinkPattern = regexp.MustCompile(`\[[^\]]+\]\[[^\]]*\]`)

	// Link definitions: [ref]: url
	linkDefPattern = regexp.MustCompile(`(?m)^\[[^\]]+\]:\s+.*$`)
)

// translateBody translates markdown body while preserving special content
func (t *MarkdownTranslator) translateBody(ctx context.Context, body, sourceLang, targetLang string) (string, error) {
	if body == "" {
		return "", nil
	}

	// Store protected content
	protected := make(map[string]string)
	counter := 0

	// Helper to create placeholder
	placeholder := func(content string) string {
		key := fmt.Sprintf("%s%d%s", placeholderPrefix, counter, placeholderSuffix)
		protected[key] = content
		counter++
		return key
	}

	// Protect content in order (code blocks first as they may contain other patterns)
	processed := body

	// 1. Code blocks (must be first - can contain anything)
	processed = codeBlockPattern.ReplaceAllStringFunc(processed, placeholder)

	// 2. Shortcodes
	processed = shortcodePattern.ReplaceAllStringFunc(processed, placeholder)

	// 3. Images (before URLs to avoid partial matches)
	processed = imagePattern.ReplaceAllStringFunc(processed, placeholder)

	// 4. Link definitions
	processed = linkDefPattern.ReplaceAllStringFunc(processed, placeholder)

	// 5. Reference links
	processed = refLinkPattern.ReplaceAllStringFunc(processed, placeholder)

	// 6. URLs in markdown links - protect just the URL part
	processed = urlPattern.ReplaceAllStringFunc(processed, placeholder)

	// 7. Inline code
	processed = inlineCodePattern.ReplaceAllStringFunc(processed, placeholder)

	// 8. HTML tags
	processed = htmlTagPattern.ReplaceAllStringFunc(processed, placeholder)

	// Translate the processed text
	translated, err := t.provider.Translate(ctx, processed, sourceLang, targetLang)
	if err != nil {
		return "", err
	}

	// Restore protected content
	for key, original := range protected {
		translated = strings.ReplaceAll(translated, key, original)
	}

	return translated, nil
}

// TranslateResult holds the result of a file translation
type TranslateResult struct {
	SourcePath string
	TargetPath string
	SourceLang string
	TargetLang string
	Success    bool
	Error      error
	CharCount  int // Characters translated (for quota tracking)
}
