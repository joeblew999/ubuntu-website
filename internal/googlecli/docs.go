package googlecli

import (
	"fmt"
	"os"
	"strings"

	"github.com/joeblew999/ubuntu-website/internal/google/gdocs"
)

func (c *cliContext) handleDocs(args []string) {
	if len(args) < 1 {
		c.printDocsUsage()
		return
	}

	cmd := args[0]
	cmdArgs := args[1:]
	jsonOutput := hasFlag(cmdArgs, "--json")
	cmdArgs = filterFlags(cmdArgs)

	config := gdocs.DefaultConfig()
	client, err := gdocs.NewAPIClient(config)
	if err != nil {
		c.exitError(fmt.Sprintf("Failed to create Docs client: %v", err))
	}

	switch cmd {
	case "get", "read":
		if len(cmdArgs) < 1 {
			c.exitError("DOCUMENT_ID required")
		}
		doc, err := client.Get(cmdArgs[0])
		if err != nil {
			c.exitError(err.Error())
		}
		if jsonOutput {
			c.outputJSON(doc)
		} else {
			fmt.Fprintf(c.stdout, "Title: %s\n", doc.Title)
			fmt.Fprintf(c.stdout, "ID: %s\n", doc.ID)
			fmt.Fprintln(c.stdout, "\n--- Content ---")
			fmt.Fprintln(c.stdout, doc.GetText())
		}

	case "create", "new":
		if len(cmdArgs) < 1 {
			c.exitError("TITLE required")
		}
		title := strings.Join(cmdArgs, " ")
		result, err := client.Create(title)
		if err != nil {
			c.exitError(err.Error())
		}
		if jsonOutput {
			c.outputJSON(result)
		} else {
			fmt.Fprintf(c.stdout, "Created: %s\n", result.Document.Title)
			fmt.Fprintf(c.stdout, "ID: %s\n", result.Document.ID)
			fmt.Fprintf(c.stdout, "Link: https://docs.google.com/document/d/%s/edit\n", result.Document.ID)
		}

	case "append", "add":
		if len(cmdArgs) < 2 {
			c.exitError("DOCUMENT_ID and TEXT required")
		}
		text := strings.Join(cmdArgs[1:], " ")
		result, err := client.AppendText(cmdArgs[0], text)
		if err != nil {
			c.exitError(err.Error())
		}
		if jsonOutput {
			c.outputJSON(result)
		} else {
			fmt.Fprintln(c.stdout, "Text appended successfully")
		}

	case "insert":
		if len(cmdArgs) < 3 {
			c.exitError("DOCUMENT_ID, INDEX, and TEXT required")
		}
		var index int
		if _, err := fmt.Sscanf(cmdArgs[1], "%d", &index); err != nil {
			c.exitError("INDEX must be a number")
		}
		text := strings.Join(cmdArgs[2:], " ")
		result, err := client.InsertText(cmdArgs[0], index, text)
		if err != nil {
			c.exitError(err.Error())
		}
		if jsonOutput {
			c.outputJSON(result)
		} else {
			fmt.Fprintln(c.stdout, "Text inserted successfully")
		}

	case "replace":
		if len(cmdArgs) < 3 {
			c.exitError("DOCUMENT_ID, FIND, and REPLACE required")
		}
		result, err := client.ReplaceText(cmdArgs[0], cmdArgs[1], cmdArgs[2])
		if err != nil {
			c.exitError(err.Error())
		}
		if jsonOutput {
			c.outputJSON(result)
		} else {
			fmt.Fprintf(c.stdout, "Replaced '%s' with '%s'\n", cmdArgs[1], cmdArgs[2])
		}

	default:
		fmt.Fprintf(c.stderr, "Unknown docs command: %s\n", cmd)
		c.printDocsUsage()
		os.Exit(1)
	}
}

func (c *cliContext) printDocsUsage() {
	fmt.Fprintln(c.stdout, `Usage: google docs <command> [arguments]

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
