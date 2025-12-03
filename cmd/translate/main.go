// translate provides automated content translation using Claude API.
//
// WARNING: DO NOT RUN - This tool is currently broken/incomplete!
// Translation is currently done manually via Claude Code + Taskfile shell scripts.
//
// Current workflow (use this instead):
//
//	task translate:status   # See what English files changed
//	task translate:missing  # See which languages need files
//	# Then manually translate with Claude Code
//	task translate:done     # Mark translations complete
//
// Commands (when fixed):
//
//	go run cmd/translate/main.go -check              # Check changed files
//	go run cmd/translate/main.go -all                # Translate all changed
//	go run cmd/translate/main.go -lang de            # Translate to German
//	go run cmd/translate/main.go -i18n               # Translate i18n files
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joeblew999/ubuntu-website/internal/translator"
)

const version = "0.1.0"

func main() {
	// Define flags
	check := flag.Bool("check", false, "Check which English files have changed since last translation")
	all := flag.Bool("all", false, "Translate all changed English content to all languages")
	lang := flag.String("lang", "", "Translate to specific language (e.g., de, sv, zh, ja, th)")
	i18n := flag.Bool("i18n", false, "Translate i18n TOML files")
	ver := flag.Bool("version", false, "Print version and exit")

	flag.Parse()

	// Print version and exit
	if *ver {
		fmt.Printf("translate v%s\n", version)
		os.Exit(0)
	}

	// Get Claude API key from environment
	apiKey := os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" {
		log.Fatal("CLAUDE_API_KEY environment variable not set")
	}

	// Create translator instance
	t, err := translator.New(apiKey)
	if err != nil {
		log.Fatalf("Failed to create translator: %v", err)
	}

	// Execute commands
	switch {
	case *check:
		if err := t.Check(); err != nil {
			log.Fatalf("Check failed: %v", err)
		}

	case *all:
		fmt.Println("🔄 Translating all changed English content to all languages...")
		if err := t.TranslateAll(); err != nil {
			log.Fatalf("Translation failed: %v", err)
		}
		fmt.Println("✅ Translation complete!")

	case *lang != "":
		fmt.Printf("🔄 Translating to %s...\n", *lang)
		if err := t.TranslateLang(*lang); err != nil {
			log.Fatalf("Translation failed: %v", err)
		}
		fmt.Printf("✅ Translation to %s complete!\n", *lang)

	case *i18n:
		fmt.Println("🔄 Translating i18n files...")
		if err := t.TranslateI18n(); err != nil {
			log.Fatalf("Translation failed: %v", err)
		}
		fmt.Println("✅ i18n translation complete!")

	default:
		flag.Usage()
		os.Exit(1)
	}
}
