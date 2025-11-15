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
