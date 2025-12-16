// Package autotranslate provides automatic translation of markdown content
// using external translation APIs (DeepL, Google, OpenAI, etc.)
package autotranslate

import (
	"context"
	"fmt"
)

// Provider defines the interface for translation services.
// Implementations: DeepL, Google Cloud Translation, OpenAI, Claude, etc.
type Provider interface {
	// Name returns the provider name (e.g., "deepl", "google", "openai")
	Name() string

	// Translate translates text from source to target language.
	// sourceLang and targetLang use ISO 639-1 codes (e.g., "en", "de", "vi")
	Translate(ctx context.Context, text, sourceLang, targetLang string) (string, error)

	// TranslateBatch translates multiple texts in a single API call (more efficient).
	// Returns translations in the same order as input texts.
	TranslateBatch(ctx context.Context, texts []string, sourceLang, targetLang string) ([]string, error)

	// SupportedLanguages returns list of supported language codes
	SupportedLanguages() []string

	// SupportsLanguage checks if a language code is supported
	SupportsLanguage(langCode string) bool
}

// ProviderConfig holds common configuration for providers
type ProviderConfig struct {
	APIKey string
}

// Registry holds available translation providers
type Registry struct {
	providers map[string]Provider
}

// NewRegistry creates a new provider registry
func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]Provider),
	}
}

// Register adds a provider to the registry
func (r *Registry) Register(p Provider) {
	r.providers[p.Name()] = p
}

// Get returns a provider by name
func (r *Registry) Get(name string) (Provider, error) {
	p, ok := r.providers[name]
	if !ok {
		return nil, fmt.Errorf("provider '%s' not found", name)
	}
	return p, nil
}

// List returns all registered provider names
func (r *Registry) List() []string {
	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}
	return names
}

// DefaultRegistry is the global provider registry
var DefaultRegistry = NewRegistry()

// RegisterProvider registers a provider in the default registry
func RegisterProvider(p Provider) {
	DefaultRegistry.Register(p)
}

// GetProvider returns a provider from the default registry
func GetProvider(name string) (Provider, error) {
	return DefaultRegistry.Get(name)
}
