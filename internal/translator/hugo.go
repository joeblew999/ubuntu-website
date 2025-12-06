// hugo.go - Hugo-specific integration for the translator.
//
// This file provides optional Hugo integration. The core translator
// works with any markdown files and target languages. Hugo integration
// adds automatic language discovery from Hugo config files.
package translator

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// HugoConfig represents Hugo-specific configuration paths
type HugoConfig struct {
	LanguagesFile string // e.g., "config/_default/languages.toml"
}

// DefaultHugoConfig returns standard Hugo config paths
func DefaultHugoConfig() *HugoConfig {
	return &HugoConfig{
		LanguagesFile: "config/_default/languages.toml",
	}
}

// ParseHugoLanguages reads language config from Hugo's languages.toml
// Returns target languages, source language code, source content dir, and error.
//
// Hugo languages.toml format:
//
//	[en]
//	languageName = "English"
//	contentDir = "content/english"
//	weight = 1
//
//	[de]
//	languageName = "German"
//	contentDir = "content/german"
//	weight = 2
//
// Weight 1 is treated as the source language; all others are targets.
func ParseHugoLanguages(configPath string) ([]Language, string, string, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, "", "", err
	}
	defer file.Close()

	var langs []Language
	var sourceLang, sourceDir string
	var currentLang string
	var currentName, currentDir string
	var currentWeight int

	// Regex patterns for TOML parsing
	sectionRe := regexp.MustCompile(`^\[([a-z]{2}(?:-[a-z]{2})?)\]$`)
	nameRe := regexp.MustCompile(`^languageName\s*=\s*"([^"]+)"`)
	dirRe := regexp.MustCompile(`^contentDir\s*=\s*"content/([^"]+)"`)
	weightRe := regexp.MustCompile(`^weight\s*=\s*(\d+)`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// New language section
		if matches := sectionRe.FindStringSubmatch(line); matches != nil {
			// Save previous language
			if currentLang != "" && currentDir != "" {
				if currentWeight == 1 {
					// Weight 1 is the source language
					sourceLang = currentLang
					sourceDir = currentDir
				} else {
					langs = append(langs, Language{
						Code:    currentLang,
						Name:    currentName,
						DirName: currentDir,
					})
				}
			}
			currentLang = matches[1]
			currentName = ""
			currentDir = ""
			currentWeight = 0
			continue
		}

		// Language name
		if matches := nameRe.FindStringSubmatch(line); matches != nil {
			currentName = matches[1]
			continue
		}

		// Content directory
		if matches := dirRe.FindStringSubmatch(line); matches != nil {
			currentDir = matches[1]
			continue
		}

		// Weight
		if matches := weightRe.FindStringSubmatch(line); matches != nil {
			fmt.Sscanf(matches[1], "%d", &currentWeight)
			continue
		}
	}

	// Save last language
	if currentLang != "" && currentDir != "" {
		if currentWeight == 1 {
			sourceLang = currentLang
			sourceDir = currentDir
		} else {
			langs = append(langs, Language{
				Code:    currentLang,
				Name:    currentName,
				DirName: currentDir,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, "", "", err
	}

	return langs, sourceLang, sourceDir, nil
}

// TryLoadHugoConfig attempts to load language config from Hugo.
// Returns nil if Hugo config not found (caller can use defaults).
func TryLoadHugoConfig(config *Config) error {
	hugoConfig := DefaultHugoConfig()

	langs, sourceLang, sourceDir, err := ParseHugoLanguages(hugoConfig.LanguagesFile)
	if err != nil {
		// Not a Hugo project or config not found - that's okay
		return nil
	}

	// Apply Hugo config to translator config
	config.TargetLangs = langs
	if sourceLang != "" {
		config.SourceLang = sourceLang
		config.SourceDir = sourceDir
	}

	return nil
}

// ValidateHugoConfig checks if translator config matches Hugo config.
// Returns list of mismatches or nil if all good.
func ValidateHugoConfig(config *Config) []string {
	hugoConfig := DefaultHugoConfig()

	hugoLangs, hugoSourceLang, hugoSourceDir, err := ParseHugoLanguages(hugoConfig.LanguagesFile)
	if err != nil {
		return []string{fmt.Sprintf("Cannot read Hugo config: %v", err)}
	}

	var mismatches []string

	// Check source language
	if config.SourceLang != hugoSourceLang {
		mismatches = append(mismatches, fmt.Sprintf(
			"Source language mismatch: translator=%s, Hugo=%s",
			config.SourceLang, hugoSourceLang))
	}

	// Check source directory
	if config.SourceDir != hugoSourceDir {
		mismatches = append(mismatches, fmt.Sprintf(
			"Source directory mismatch: translator=%s, Hugo=%s",
			config.SourceDir, hugoSourceDir))
	}

	// Check target languages
	hugoLangMap := make(map[string]Language)
	for _, lang := range hugoLangs {
		hugoLangMap[lang.Code] = lang
	}

	configLangMap := make(map[string]Language)
	for _, lang := range config.TargetLangs {
		configLangMap[lang.Code] = lang
	}

	// Languages in translator but not Hugo
	for code := range configLangMap {
		if _, ok := hugoLangMap[code]; !ok {
			mismatches = append(mismatches, fmt.Sprintf(
				"Language '%s' in translator but not in Hugo config", code))
		}
	}

	// Languages in Hugo but not translator
	for code := range hugoLangMap {
		if _, ok := configLangMap[code]; !ok {
			mismatches = append(mismatches, fmt.Sprintf(
				"Language '%s' in Hugo config but not in translator", code))
		}
	}

	return mismatches
}

// IsHugoProject checks if current directory is a Hugo project
func IsHugoProject() bool {
	hugoConfig := DefaultHugoConfig()
	_, err := os.Stat(hugoConfig.LanguagesFile)
	return err == nil
}
