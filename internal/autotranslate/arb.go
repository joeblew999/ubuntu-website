package autotranslate

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"
	"time"
)

// ARBFile represents an Application Resource Bundle file.
type ARBFile struct {
	Locale           string            `json:"-"`
	LastModified     time.Time         `json:"-"`
	Messages         map[string]string // message ID -> ICU message
	Metadata         map[string]any    // @msgID -> metadata
	CustomAttributes map[string]any    // @@x-... attributes
}

// LoadARB loads an ARB file from disk.
func LoadARB(path string) (*ARBFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}

	arb := &ARBFile{
		Messages:         make(map[string]string),
		Metadata:         make(map[string]any),
		CustomAttributes: make(map[string]any),
	}

	for key, value := range raw {
		if key == "@@locale" {
			var s string
			json.Unmarshal(value, &s)
			arb.Locale = s
		} else if key == "@@last_modified" {
			var s string
			json.Unmarshal(value, &s)
			arb.LastModified, _ = time.Parse(time.RFC3339, s)
		} else if strings.HasPrefix(key, "@@") {
			var v any
			json.Unmarshal(value, &v)
			arb.CustomAttributes[key] = v
		} else if strings.HasPrefix(key, "@") {
			var v any
			json.Unmarshal(value, &v)
			arb.Metadata[key] = v
		} else {
			var s string
			json.Unmarshal(value, &s)
			arb.Messages[key] = s
		}
	}

	return arb, nil
}

// SaveARB writes an ARB file to disk.
func SaveARB(path string, arb *ARBFile) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Build ordered entries
	type kv struct {
		Key   string
		Value any
	}

	var entries []kv

	// Metadata first
	entries = append(entries, kv{"@@locale", arb.Locale})
	entries = append(entries, kv{"@@last_modified", time.Now().Format(time.RFC3339)})

	// Custom attributes
	for _, k := range slices.Sorted(maps.Keys(arb.CustomAttributes)) {
		entries = append(entries, kv{k, arb.CustomAttributes[k]})
	}

	// Messages and their metadata
	for _, id := range slices.Sorted(maps.Keys(arb.Messages)) {
		entries = append(entries, kv{id, arb.Messages[id]})
		metaKey := "@" + id
		if meta, ok := arb.Metadata[metaKey]; ok {
			entries = append(entries, kv{metaKey, meta})
		}
	}

	// Write JSON with indentation
	f.WriteString("{\n")
	for i, entry := range entries {
		keyJSON, _ := json.Marshal(entry.Key)
		valueJSON, _ := json.MarshalIndent(entry.Value, "\t", "\t")

		// Handle multi-line values
		valueStr := string(valueJSON)
		if strings.Contains(valueStr, "\n") {
			// Already indented by MarshalIndent
		}

		f.WriteString("\t")
		f.Write(keyJSON)
		f.WriteString(": ")
		f.WriteString(valueStr)

		if i < len(entries)-1 {
			f.WriteString(",")
		}
		f.WriteString("\n")
	}
	f.WriteString("}\n")

	return nil
}

// ARBTranslator handles translation of ARB catalog entries.
type ARBTranslator struct {
	provider  Provider
	batchSize int
}

// NewARBTranslator creates a new ARB translator.
func NewARBTranslator(provider Provider) *ARBTranslator {
	return &ARBTranslator{
		provider:  provider,
		batchSize: 50, // Translate 50 messages at a time
	}
}

// TranslateARB translates empty entries in the target ARB using source ARB.
// Returns the number of entries translated.
func (t *ARBTranslator) TranslateARB(ctx context.Context, sourceARB, targetARB *ARBFile, targetLang string, verbose bool, progressFn func(done, total int)) (int, error) {
	// Find entries that need translation (empty in target, non-empty in source)
	var toTranslate []struct {
		ID         string
		SourceText string
	}

	for id, sourceText := range sourceARB.Messages {
		if sourceText == "" {
			continue
		}
		targetText, exists := targetARB.Messages[id]
		if !exists || targetText == "" {
			toTranslate = append(toTranslate, struct {
				ID         string
				SourceText string
			}{id, sourceText})
		}
	}

	if len(toTranslate) == 0 {
		return 0, nil
	}

	total := len(toTranslate)
	translated := 0

	// Process in batches
	for i := 0; i < len(toTranslate); i += t.batchSize {
		end := i + t.batchSize
		if end > len(toTranslate) {
			end = len(toTranslate)
		}
		batch := toTranslate[i:end]

		// Extract texts for this batch
		texts := make([]string, len(batch))
		for j, item := range batch {
			// Unescape ICU format for translation
			texts[j] = unescapeICU(item.SourceText)
		}

		// Translate batch
		translations, err := t.provider.TranslateBatch(ctx, texts, "en", targetLang)
		if err != nil {
			return translated, fmt.Errorf("translating batch %d-%d: %w", i, end, err)
		}

		// Update target ARB
		for j, item := range batch {
			if j < len(translations) && translations[j] != "" {
				// Re-escape for ICU format
				targetARB.Messages[item.ID] = escapeICU(translations[j])
				translated++
			}
		}

		if progressFn != nil {
			progressFn(i+len(batch), total)
		}

		if verbose {
			fmt.Printf("  Translated %d/%d entries\n", i+len(batch), total)
		}
	}

	return translated, nil
}

// unescapeICU converts ICU message format back to plain text for translation.
func unescapeICU(s string) string {
	// Remove wrapping quotes if present (for special characters)
	if len(s) >= 2 && s[0] == '\'' && s[len(s)-1] == '\'' {
		s = s[1 : len(s)-1]
	}
	// Unescape double single quotes
	s = strings.ReplaceAll(s, "''", "'")
	return s
}

// escapeICU escapes text for ICU message format.
func escapeICU(s string) string {
	// Check if escaping is needed
	hasSpecial := strings.ContainsAny(s, "{}#|")
	hasSingleQuote := strings.Contains(s, "'")

	if !hasSpecial && !hasSingleQuote {
		return s
	}

	// Escape single quotes first (double them)
	s = strings.ReplaceAll(s, "'", "''")

	// If there are special characters, wrap in quotes
	if hasSpecial {
		s = "'" + s + "'"
	}

	return s
}

// ARBStats holds statistics about an ARB file.
type ARBStats struct {
	TotalMessages    int
	TranslatedCount  int
	EmptyCount       int
	CompletenessPerc float64
}

// GetARBStats calculates statistics for an ARB file.
func GetARBStats(arb *ARBFile) ARBStats {
	stats := ARBStats{
		TotalMessages: len(arb.Messages),
	}

	for _, msg := range arb.Messages {
		if msg == "" {
			stats.EmptyCount++
		} else {
			stats.TranslatedCount++
		}
	}

	if stats.TotalMessages > 0 {
		stats.CompletenessPerc = float64(stats.TranslatedCount) / float64(stats.TotalMessages) * 100
	}

	return stats
}
