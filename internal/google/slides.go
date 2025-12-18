package google

import (
	"fmt"
	"os"
	"strings"

	"github.com/joeblew999/ubuntu-website/internal/google/gslides"
)

func (c *cliContext) handleSlides(args []string) {
	if len(args) < 1 {
		c.printSlidesUsage()
		return
	}

	cmd := args[0]
	cmdArgs := args[1:]
	jsonOutput := hasFlag(cmdArgs, "--json")
	cmdArgs = filterFlags(cmdArgs)

	config := gslides.DefaultConfig()
	client, err := gslides.NewAPIClient(config)
	if err != nil {
		c.exitError(fmt.Sprintf("Failed to create Slides client: %v", err))
	}

	switch cmd {
	case "get", "info":
		if len(cmdArgs) < 1 {
			c.exitError("PRESENTATION_ID required")
		}
		pres, err := client.Get(cmdArgs[0])
		if err != nil {
			c.exitError(err.Error())
		}
		if jsonOutput {
			c.outputJSON(pres)
		} else {
			fmt.Fprintf(c.stdout, "Title: %s\n", pres.Title)
			fmt.Fprintf(c.stdout, "ID: %s\n", pres.ID)
			fmt.Fprintf(c.stdout, "Slides: %d\n", len(pres.Slides))
			fmt.Fprintf(c.stdout, "Link: https://docs.google.com/presentation/d/%s/edit\n", pres.ID)
			if len(pres.Slides) > 0 {
				fmt.Fprintln(c.stdout, "\nSlide IDs:")
				for i, slide := range pres.Slides {
					fmt.Fprintf(c.stdout, "  %d. %s\n", i+1, slide.ObjectID)
				}
			}
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
			fmt.Fprintf(c.stdout, "Created: %s\n", result.Presentation.Title)
			fmt.Fprintf(c.stdout, "ID: %s\n", result.Presentation.ID)
			fmt.Fprintf(c.stdout, "Link: https://docs.google.com/presentation/d/%s/edit\n", result.Presentation.ID)
		}

	case "add-slide", "addslide":
		if len(cmdArgs) < 1 {
			c.exitError("PRESENTATION_ID required")
		}
		index := 0
		if len(cmdArgs) > 1 {
			fmt.Sscanf(cmdArgs[1], "%d", &index)
		}
		result, err := client.AddSlide(cmdArgs[0], index)
		if err != nil {
			c.exitError(err.Error())
		}
		if jsonOutput {
			c.outputJSON(result)
		} else {
			fmt.Fprintf(c.stdout, "Added slide at position %d\n", index)
		}

	case "delete-slide", "rm-slide":
		if len(cmdArgs) < 2 {
			c.exitError("PRESENTATION_ID and SLIDE_ID required")
		}
		result, err := client.DeleteSlide(cmdArgs[0], cmdArgs[1])
		if err != nil {
			c.exitError(err.Error())
		}
		if jsonOutput {
			c.outputJSON(result)
		} else {
			fmt.Fprintln(c.stdout, "Slide deleted")
		}

	default:
		fmt.Fprintf(c.stderr, "Unknown slides command: %s\n", cmd)
		c.printSlidesUsage()
		os.Exit(1)
	}
}

func (c *cliContext) printSlidesUsage() {
	fmt.Fprintln(c.stdout, `Usage: google slides <command> [arguments]

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
