package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/joeblew999/ubuntu-website/internal/google/gsheets"
)

func handleSheets(args []string) {
	if len(args) < 1 {
		printSheetsUsage()
		return
	}

	cmd := args[0]
	cmdArgs := args[1:]
	jsonOutput := hasFlag(cmdArgs, "--json")
	cmdArgs = filterFlags(cmdArgs)

	config := gsheets.DefaultConfig()
	client, err := gsheets.NewAPIClient(config)
	if err != nil {
		exitError(fmt.Sprintf("Failed to create Sheets client: %v", err))
	}

	switch cmd {
	case "info", "meta":
		if len(cmdArgs) < 1 {
			exitError("SPREADSHEET_ID required")
		}
		sheet, err := client.GetSpreadsheet(cmdArgs[0])
		if err != nil {
			exitError(err.Error())
		}
		if jsonOutput {
			outputJSON(sheet)
		} else {
			fmt.Printf("Title: %s\n", sheet.Title)
			fmt.Printf("ID: %s\n", sheet.ID)
			if sheet.Locale != "" {
				fmt.Printf("Locale: %s\n", sheet.Locale)
			}
			if sheet.TimeZone != "" {
				fmt.Printf("TimeZone: %s\n", sheet.TimeZone)
			}
			if len(sheet.SheetNames) > 0 {
				fmt.Printf("Sheets: %s\n", strings.Join(sheet.SheetNames, ", "))
			}
		}

	case "get", "read":
		if len(cmdArgs) < 2 {
			exitError("SPREADSHEET_ID and RANGE required")
		}
		result, err := client.GetValues(cmdArgs[0], cmdArgs[1])
		if err != nil {
			exitError(err.Error())
		}
		if jsonOutput {
			outputJSON(result)
		} else {
			fmt.Printf("Range: %s\n\n", result.Range)
			if len(result.Values) == 0 {
				fmt.Println("(empty)")
				return
			}
			for _, row := range result.Values {
				var cells []string
				for _, cell := range row {
					cells = append(cells, fmt.Sprintf("%v", cell))
				}
				fmt.Println(strings.Join(cells, "\t"))
			}
		}

	case "set", "update":
		if len(cmdArgs) < 3 {
			exitError("SPREADSHEET_ID, RANGE, and VALUES required")
		}
		// Parse values: comma-separated for single row, or multiple args for multiple cells
		var values [][]interface{}
		if len(cmdArgs) == 3 {
			// Single row, comma-separated
			cells := strings.Split(cmdArgs[2], ",")
			row := make([]interface{}, len(cells))
			for i, c := range cells {
				row[i] = strings.TrimSpace(c)
			}
			values = [][]interface{}{row}
		} else {
			// Multiple args = single row
			row := make([]interface{}, len(cmdArgs)-2)
			for i, c := range cmdArgs[2:] {
				row[i] = c
			}
			values = [][]interface{}{row}
		}

		result, err := client.UpdateValues(cmdArgs[0], cmdArgs[1], values)
		if err != nil {
			exitError(err.Error())
		}
		if jsonOutput {
			outputJSON(result)
		} else {
			fmt.Printf("Updated %d cells in %s\n", result.UpdatedCells, result.UpdatedRange)
		}

	case "append", "add":
		if len(cmdArgs) < 3 {
			exitError("SPREADSHEET_ID, RANGE, and VALUES required")
		}
		// Parse values same as set
		var values [][]interface{}
		if len(cmdArgs) == 3 {
			cells := strings.Split(cmdArgs[2], ",")
			row := make([]interface{}, len(cells))
			for i, c := range cells {
				row[i] = strings.TrimSpace(c)
			}
			values = [][]interface{}{row}
		} else {
			row := make([]interface{}, len(cmdArgs)-2)
			for i, c := range cmdArgs[2:] {
				row[i] = c
			}
			values = [][]interface{}{row}
		}

		result, err := client.AppendValues(cmdArgs[0], cmdArgs[1], values)
		if err != nil {
			exitError(err.Error())
		}
		if jsonOutput {
			outputJSON(result)
		} else {
			fmt.Printf("Appended %d rows to %s\n", result.UpdatedRows, result.UpdatedRange)
		}

	case "clear":
		if len(cmdArgs) < 2 {
			exitError("SPREADSHEET_ID and RANGE required")
		}
		if err := client.Clear(cmdArgs[0], cmdArgs[1]); err != nil {
			exitError(err.Error())
		}
		fmt.Printf("Cleared %s\n", cmdArgs[1])

	default:
		fmt.Fprintf(os.Stderr, "Unknown sheets command: %s\n", cmd)
		printSheetsUsage()
		os.Exit(1)
	}
}

func printSheetsUsage() {
	fmt.Println(`Usage: google sheets <command> [arguments]

Commands:
  info SPREADSHEET_ID              Get spreadsheet metadata
  get SPREADSHEET_ID RANGE         Get cell values
  set SPREADSHEET_ID RANGE VALUES  Update cells
  append SPREADSHEET_ID RANGE VALUES  Append rows
  clear SPREADSHEET_ID RANGE       Clear a range

Options:
  --json    Output as JSON

Examples:
  google sheets info 1abc123def
  google sheets get 1abc123def "Sheet1!A1:D10"
  google sheets set 1abc123def "Sheet1!A1" "Hello,World,Test"
  google sheets set 1abc123def "Sheet1!A1:C1" Hello World Test
  google sheets append 1abc123def "Sheet1!A:C" "New,Row,Data"
  google sheets clear 1abc123def "Sheet1!A1:D10"`)
}
