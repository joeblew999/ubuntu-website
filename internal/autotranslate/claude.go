package autotranslate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	claudeAPIURL     = "https://api.anthropic.com/v1/messages"
	claudeModel      = "claude-sonnet-4-20250514"
	claudeMaxTokens  = 4096
	claudeAPIVersion = "2023-06-01"
)

// Language name mappings for Claude prompts
var claudeLangNames = map[string]string{
	"en": "English",
	"de": "German",
	"fr": "French",
	"es": "Spanish",
	"it": "Italian",
	"nl": "Dutch",
	"pl": "Polish",
	"pt": "Portuguese",
	"ru": "Russian",
	"ja": "Japanese",
	"zh": "Chinese (Simplified)",
	"ko": "Korean",
	"vi": "Vietnamese",
	"id": "Indonesian",
	"tr": "Turkish",
	"uk": "Ukrainian",
	"cs": "Czech",
	"da": "Danish",
	"fi": "Finnish",
	"el": "Greek",
	"hu": "Hungarian",
	"lt": "Lithuanian",
	"lv": "Latvian",
	"nb": "Norwegian (Bokmål)",
	"ro": "Romanian",
	"sk": "Slovak",
	"sl": "Slovenian",
	"sv": "Swedish",
	"bg": "Bulgarian",
	"et": "Estonian",
	"ar": "Arabic",
	"he": "Hebrew",
	"hi": "Hindi",
	"th": "Thai",
}

// ClaudeProvider implements Provider using Claude API
type ClaudeProvider struct {
	apiKey     string
	httpClient *http.Client
}

// claudeRequest represents a request to Claude API
type claudeRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	Messages  []claudeMessage `json:"messages"`
}

// claudeMessage represents a message in the conversation
type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// claudeResponse represents a response from Claude API
type claudeResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

// NewClaudeProvider creates a new Claude translation provider
func NewClaudeProvider(apiKey string) (*ClaudeProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("Claude API key is required")
	}

	return &ClaudeProvider{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 120 * time.Second, // Longer timeout for batch translations
		},
	}, nil
}

// Name returns the provider name
func (p *ClaudeProvider) Name() string {
	return "claude"
}

// Translate translates text using Claude
func (p *ClaudeProvider) Translate(ctx context.Context, text, sourceLang, targetLang string) (string, error) {
	targetName := claudeLangNames[strings.ToLower(targetLang)]
	if targetName == "" {
		targetName = targetLang // Fallback to code if name not found
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

	return p.callAPI(ctx, prompt)
}

// TranslateBatch translates multiple texts efficiently
func (p *ClaudeProvider) TranslateBatch(ctx context.Context, texts []string, sourceLang, targetLang string) ([]string, error) {
	if len(texts) == 0 {
		return []string{}, nil
	}

	// For small batches, translate individually
	if len(texts) <= 3 {
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

	// For larger batches, use numbered format
	targetName := claudeLangNames[strings.ToLower(targetLang)]
	if targetName == "" {
		targetName = targetLang
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Translate each numbered text from English to %s.\n\n", targetName))
	builder.WriteString("CRITICAL INSTRUCTIONS:\n")
	builder.WriteString("1. Return translations in EXACTLY the same numbered format\n")
	builder.WriteString("2. Preserve ALL markdown formatting\n")
	builder.WriteString("3. DO NOT translate URLs, code blocks, shortcodes, or HTML tags\n")
	builder.WriteString("4. Return ONLY the numbered translations - no explanations\n\n")
	builder.WriteString("Texts to translate:\n\n")

	for i, text := range texts {
		builder.WriteString(fmt.Sprintf("[%d]\n%s\n\n", i+1, text))
	}

	response, err := p.callAPI(ctx, builder.String())
	if err != nil {
		return nil, err
	}

	// Parse numbered responses
	return p.parseNumberedResponse(response, len(texts))
}

// parseNumberedResponse extracts translations from numbered format
func (p *ClaudeProvider) parseNumberedResponse(response string, count int) ([]string, error) {
	results := make([]string, count)

	// Split by [N] markers
	for i := 0; i < count; i++ {
		startMarker := fmt.Sprintf("[%d]", i+1)
		endMarker := fmt.Sprintf("[%d]", i+2)

		startIdx := strings.Index(response, startMarker)
		if startIdx == -1 {
			return nil, fmt.Errorf("missing translation for item %d", i+1)
		}

		// Move past the marker
		startIdx += len(startMarker)

		var endIdx int
		if i == count-1 {
			// Last item - take rest of string
			endIdx = len(response)
		} else {
			endIdx = strings.Index(response, endMarker)
			if endIdx == -1 {
				endIdx = len(response)
			}
		}

		results[i] = strings.TrimSpace(response[startIdx:endIdx])
	}

	return results, nil
}

// SupportedLanguages returns list of supported language codes
func (p *ClaudeProvider) SupportedLanguages() []string {
	langs := make([]string, 0, len(claudeLangNames))
	for code := range claudeLangNames {
		langs = append(langs, code)
	}
	return langs
}

// SupportsLanguage checks if a language code is supported
func (p *ClaudeProvider) SupportsLanguage(langCode string) bool {
	_, ok := claudeLangNames[strings.ToLower(langCode)]
	return ok
}

// callAPI makes a request to Claude API
func (p *ClaudeProvider) callAPI(ctx context.Context, prompt string) (string, error) {
	req := claudeRequest{
		Model:     claudeModel,
		MaxTokens: claudeMaxTokens,
		Messages: []claudeMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", claudeAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.apiKey)
	httpReq.Header.Set("anthropic-version", claudeAPIVersion)

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Claude API error (status %d): %s", resp.StatusCode, string(body))
	}

	var claudeResp claudeResponse
	if err := json.Unmarshal(body, &claudeResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(claudeResp.Content) == 0 {
		return "", fmt.Errorf("no content in Claude response")
	}

	return claudeResp.Content[0].Text, nil
}
