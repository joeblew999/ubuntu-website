package main

import (
	"fmt"
	"os"

	"github.com/joeblew999/ubuntu-website/internal/env/web"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	var err error

	switch command {
	case "web-gui":
		err = web.ServeSetupGUI()
	case "web-gui-mock":
		err = web.ServeSetupGUIMock()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: go run cmd/env/main.go <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  web-gui             Open web GUI for environment setup")
	fmt.Println("  web-gui-mock        Open web GUI with mock validation (for testing)")
}
