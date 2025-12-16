package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/joeblew999/ubuntu-website/internal/google/gdocs"
)

func handleDocs(args []string) {
	if len(args) < 1 {
		printDocsUsage()
		return
	}

	cmd := args[0]
	cmdArgs := args[1:]
	jsonOutput := hasFlag(cmdArgs, "--json")
	cmdArgs = filterFlags(cmdArgs)

	config := gdocs.DefaultConfig()
	client, err := gdocs.NewAPIClient(config)
	if err != nil {
		exitError(fmt.Sprintf("Failed to create Docs client: %v", err))
	}

	switch cmd {
	case "get", "read":
		if len(cmdArgs) < 1 {
			exitError("DOCUMENT_ID required")
		}
		doc, err := client.Get(cmdArgs[0])
		if err != nil {
			exitError(err.Error())
		}
		if jsonOutput {
			outputJSON(doc)
		} else {
			fmt.Printf("Title: %s\n", doc.Title)
			fmt.Printf("ID: %s\n", doc.ID)
			fmt.Println("\n--- Content ---")
			fmt.Println(doc.GetText())
		}

	case "create", "new":
		if len(cmdArgs) < 1 {
			exitError("TITLE required")
		}
		title := strings.Join(cmdArgs, " ")
		result, err := client.Create(title)
		if err != nil {
			exitError(err.Error())
		}
		if jsonOutput {
			outputJSON(result)
		} else {
			fmt.Printf("Created: %s\n", result.Document.Title)
			fmt.Printf("ID: %s\n", result.Document.ID)
			fmt.Printf("Link: https://docs.google.com/document/d/%s/edit\n", result.Document.ID)
		}

	case "append", "add":
		if len(cmdArgs) < 2 {
			exitError("DOCUMENT_ID and TEXT required")
		}
		text := strings.Join(cmdArgs[1:], " ")
		result, err := client.AppendText(cmdArgs[0], text)
		if err != nil {
			exitError(err.Error())
		}
		if jsonOutput {
			outputJSON(result)
		} else {
			fmt.Println("Text appended successfully")
		}

	case "insert":
		if len(cmdArgs) < 3 {
			exitError("DOCUMENT_ID, INDEX, and TEXT required")
		}
		var index int
		if _, err := fmt.Sscanf(cmdArgs[1], "%d", &index); err != nil {
			exitError("INDEX must be a number")
		}
		text := strings.Join(cmdArgs[2:], " ")
		result, err := client.InsertText(cmdArgs[0], index, text)
		if err != nil {
			exitError(err.Error())
		}
		if jsonOutput {
			outputJSON(result)
		} else {
			fmt.Println("Text inserted successfully")
		}

	case "replace":
		if len(cmdArgs) < 3 {
			exitError("DOCUMENT_ID, FIND, and REPLACE required")
		}
		result, err := client.ReplaceText(cmdArgs[0], cmdArgs[1], cmdArgs[2])
		if err != nil {
			exitError(err.Error())
		}
		if jsonOutput {
			outputJSON(result)
		} else {
			fmt.Printf("Replaced '%s' with '%s'\n", cmdArgs[1], cmdArgs[2])
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown docs command: %s\n", cmd)
		printDocsUsage()
		os.Exit(1)
	}
}

func printDocsUsage() {
	fmt.Println(`Usage: google docs <command> [arguments]

Commands:
  get DOCUMENT_ID                  Get document content
  create TITLE                     Create new document
  append DOCUMENT_ID TEXT          Append text to end
  insert DOCUMENT_ID INDEX TEXT    Insert text at index
  replace DOCUMENT_ID FIND REPLACE Replace all occurrences

Options:
  --json    Output as JSON

Examples:
  google docs get 1abc123def
  google docs create "Meeting Notes"
  google docs append 1abc123def "New paragraph text"
  google docs insert 1abc123def 1 "Text at beginning"
  google docs replace 1abc123def "old text" "new text"`)
}
