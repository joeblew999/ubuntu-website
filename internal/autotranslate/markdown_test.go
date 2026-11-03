package autotranslate

import (
	"context"
	"regexp"
	"testing"
)

// echoProvider returns text unchanged, so a round trip must reproduce the input exactly.
type echoProvider struct{ mangle func(string) string }

func (echoProvider) Name() string { return "echo" }
func (p echoProvider) Translate(_ context.Context, text, _, _ string) (string, error) {
	if p.mangle != nil {
		return p.mangle(text), nil
	}
	return text, nil
}
func (p echoProvider) TranslateBatch(ctx context.Context, texts []string, src, dst string) ([]string, error) {
	out := make([]string, len(texts))
	for i, t := range texts {
		out[i], _ = p.Translate(ctx, t, src, dst)
	}
	return out, nil
}
func (echoProvider) SupportedLanguages() []string { return []string{"en", "vi"} }
func (echoProvider) SupportsLanguage(string) bool { return true }

func TestTranslateBodyRestoresNestedPlaceholders(t *testing.T) {
	// relref shortcodes inside link URLs get wrapped twice (shortcode, then URL); map-order
	// restoration used to leave [[NOTRANSLATE_n]] behind at random.
	body := "Explore the [fleet docs]({{< relref \"/fleet/architecture\" >}}) and\n" +
		"[Back to Partners →]({{< relref \"/partners\" >}}).\n\n" +
		"Plain [link](https://example.com/a), `inline` code, ![img](images/x.png),\n" +
		"<span class=\"x\">html</span>, {{< notice \"info\" >}}note{{< /notice >}}\n\n" +
		"```go\nfmt.Println(\"[x](y)\")\n```\n"

	tr := NewMarkdownTranslator(echoProvider{})
	for i := 0; i < 50; i++ { // map iteration order is random; repeat to catch order bugs
		got, err := tr.translateBody(context.Background(), body, "en", "vi")
		if err != nil {
			t.Fatalf("run %d: %v", i, err)
		}
		if got != body {
			t.Fatalf("run %d: round trip changed the body:\n got: %q\nwant: %q", i, got, body)
		}
	}
}

func TestTranslateBodyRejectsLostPlaceholders(t *testing.T) {
	// A provider that rewrites placeholders must not produce output with leftovers.
	renumber := regexp.MustCompile(`\[\[NOTRANSLATE_\d+\]\]`)
	tr := NewMarkdownTranslator(echoProvider{mangle: func(s string) string {
		return renumber.ReplaceAllString(s, placeholderPrefix+"99"+placeholderSuffix)
	}})
	if _, err := tr.translateBody(context.Background(), "See [docs]({{< relref \"/x\" >}}).", "en", "vi"); err == nil {
		t.Fatal("expected an error for an unrestorable placeholder, got nil")
	}
}
