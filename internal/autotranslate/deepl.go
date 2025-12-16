package autotranslate

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/bounoable/deepl"
)

// DeepL language code mappings to deepl.Language type
// Note: Vietnamese (vi) is not yet in the Go library but is supported by DeepL API
var deeplLangMap = map[string]deepl.Language{
	"en": deepl.English,
	"de": deepl.German,
	"fr": deepl.French,
	"es": deepl.Spanish,
	"it": deepl.Italian,
	"nl": deepl.Dutch,
	"pl": deepl.Polish,
	"pt": deepl.PortuguesePortugal,
	"ru": deepl.Russian,
	"ja": deepl.Japanese,
	"zh": deepl.ChineseSimplified,
	"bg": deepl.Bulgarian,
	"cs": deepl.Czech,
	"da": deepl.Danish,
	"el": deepl.Greek,
	"et": deepl.Estonian,
	"fi": deepl.Finnish,
	"hu": deepl.Hungarian,
	"lt": deepl.Lithuanian,
	"lv": deepl.Latvian,
	"ro": deepl.Romanian,
	"sk": deepl.Slovak,
	"sl": deepl.Slovenian,
	"sv": deepl.Swedish,
	"id": deepl.Indonesian,
	"tr": deepl.Turkish,
	"uk": deepl.Ukrainian,
	"ko": deepl.Korean,
	"nb": deepl.NorwegianBokmal,
	"ar": deepl.Arabic,
	// Vietnamese - use raw language code since not in Go library yet
	"vi": "VI",
}

// DeepLProvider implements Provider using the DeepL API
type DeepLProvider struct {
	client *deepl.Client
}

// API endpoints
const (
	DeepLProURL  = "https://api.deepl.com/v2"
	DeepLFreeURL = "https://api-free.deepl.com/v2"
)

// Usage holds DeepL API usage statistics
type Usage struct {
	CharacterCount int64 // Characters used this billing period
	CharacterLimit int64 // Character limit for this billing period
}

// NewDeepLProvider creates a new DeepL translation provider.
// If apiKey ends with ":fx", uses the free API endpoint automatically.
func NewDeepLProvider(apiKey string) (*DeepLProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("DeepL API key is required")
	}

	// Free tier keys end with ":fx" and require different endpoint
	var client *deepl.Client
	if strings.HasSuffix(apiKey, ":fx") {
		client = deepl.New(apiKey, deepl.BaseURL(DeepLFreeURL))
	} else {
		client = deepl.New(apiKey, deepl.BaseURL(DeepLProURL))
	}

	return &DeepLProvider{client: client}, nil
}

// Name returns the provider name
func (p *DeepLProvider) Name() string {
	return "deepl"
}

// Translate translates text using DeepL
func (p *DeepLProvider) Translate(ctx context.Context, text, sourceLang, targetLang string) (string, error) {
	// Convert language codes to DeepL format
	tgtLang, ok := deeplLangMap[strings.ToLower(targetLang)]
	if !ok {
		return "", fmt.Errorf("unsupported target language: %s", targetLang)
	}

	// Build options
	opts := []deepl.TranslateOption{
		deepl.PreserveFormatting(true),
	}

	// Source language is optional for DeepL (auto-detect)
	if srcLang, ok := deeplLangMap[strings.ToLower(sourceLang)]; ok {
		opts = append(opts, deepl.SourceLang(srcLang))
	}

	// Call DeepL API
	translated, _, err := p.client.Translate(ctx, text, tgtLang, opts...)
	if err != nil {
		return "", fmt.Errorf("DeepL API error: %w", err)
	}

	return translated, nil
}

// TranslateBatch translates multiple texts by calling Translate for each.
// Note: DeepL's Go client doesn't have a native batch API, so we translate sequentially.
// For better performance with many texts, consider chunking on the caller side.
func (p *DeepLProvider) TranslateBatch(ctx context.Context, texts []string, sourceLang, targetLang string) ([]string, error) {
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
func (p *DeepLProvider) SupportedLanguages() []string {
	langs := make([]string, 0, len(deeplLangMap))
	for code := range deeplLangMap {
		langs = append(langs, code)
	}
	return langs
}

// SupportsLanguage checks if a language code is supported
func (p *DeepLProvider) SupportsLanguage(langCode string) bool {
	_, ok := deeplLangMap[strings.ToLower(langCode)]
	return ok
}

// GetUsage retrieves current API usage statistics from DeepL
func (p *DeepLProvider) GetUsage(ctx context.Context) (*Usage, error) {
	// Build request to /usage endpoint
	url := p.client.BaseURL() + "/usage"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "DeepL-Auth-Key "+p.client.AuthKey())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("usage request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("usage request returned status %d", resp.StatusCode)
	}

	var result struct {
		CharacterCount int64 `json:"character_count"`
		CharacterLimit int64 `json:"character_limit"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode usage response: %w", err)
	}

	return &Usage{
		CharacterCount: result.CharacterCount,
		CharacterLimit: result.CharacterLimit,
	}, nil
}
