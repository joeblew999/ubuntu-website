package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/joeblew999/ubuntu-website/internal/google/gslides"
)

func handleSlides(args []string) {
	if len(args) < 1 {
		printSlidesUsage()
		return
	}

	cmd := args[0]
	cmdArgs := args[1:]
	jsonOutput := hasFlag(cmdArgs, "--json")
	cmdArgs = filterFlags(cmdArgs)

	config := gslides.DefaultConfig()
	client, err := gslides.NewAPIClient(config)
	if err != nil {
		exitError(fmt.Sprintf("Failed to create Slides client: %v", err))
	}

	switch cmd {
	case "get", "info":
		if len(cmdArgs) < 1 {
			exitError("PRESENTATION_ID required")
		}
		pres, err := client.Get(cmdArgs[0])
		if err != nil {
			exitError(err.Error())
		}
		if jsonOutput {
			outputJSON(pres)
		} else {
			fmt.Printf("Title: %s\n", pres.Title)
			fmt.Printf("ID: %s\n", pres.ID)
			fmt.Printf("Slides: %d\n", len(pres.Slides))
			fmt.Printf("Link: https://docs.google.com/presentation/d/%s/edit\n", pres.ID)
			if len(pres.Slides) > 0 {
				fmt.Println("\nSlide IDs:")
				for i, slide := range pres.Slides {
					fmt.Printf("  %d. %s\n", i+1, slide.ObjectID)
				}
			}
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
			fmt.Printf("Created: %s\n", result.Presentation.Title)
			fmt.Printf("ID: %s\n", result.Presentation.ID)
			fmt.Printf("Link: https://docs.google.com/presentation/d/%s/edit\n", result.Presentation.ID)
		}

	case "add-slide", "addslide":
		if len(cmdArgs) < 1 {
			exitError("PRESENTATION_ID required")
		}
		index := 0
		if len(cmdArgs) > 1 {
			fmt.Sscanf(cmdArgs[1], "%d", &index)
		}
		result, err := client.AddSlide(cmdArgs[0], index)
		if err != nil {
			exitError(err.Error())
		}
		if jsonOutput {
			outputJSON(result)
		} else {
			fmt.Printf("Added slide at position %d\n", index)
		}

	case "delete-slide", "rm-slide":
		if len(cmdArgs) < 2 {
			exitError("PRESENTATION_ID and SLIDE_ID required")
		}
		result, err := client.DeleteSlide(cmdArgs[0], cmdArgs[1])
		if err != nil {
			exitError(err.Error())
		}
		if jsonOutput {
			outputJSON(result)
		} else {
			fmt.Println("Slide deleted")
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown slides command: %s\n", cmd)
		printSlidesUsage()
		os.Exit(1)
	}
}

func printSlidesUsage() {
	fmt.Println(`Usage: google slides <command> [arguments]

Commands:
  get PRESENTATION_ID              Get presentation metadata
  create TITLE                     Create new presentation
  add-slide PRESENTATION_ID [INDEX]  Add a new slide
  delete-slide PRESENTATION_ID SLIDE_ID  Delete a slide

Options:
  --json    Output as JSON

Examples:
  google slides get 1abc123def
  google slides create "Quarterly Report"
  google slides add-slide 1abc123def
  google slides add-slide 1abc123def 2
  google slides delete-slide 1abc123def slide_id_here`)
}
