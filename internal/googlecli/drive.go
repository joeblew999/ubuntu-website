package googlecli

import (
	"fmt"
	"os"

	"github.com/joeblew999/ubuntu-website/internal/google/gdrive"
)

func (c *cliContext) handleDrive(args []string) {
	if len(args) < 1 {
		c.printDriveUsage()
		return
	}

	cmd := args[0]
	cmdArgs := args[1:]
	jsonOutput := hasFlag(cmdArgs, "--json")
	cmdArgs = filterFlags(cmdArgs)

	config := gdrive.DefaultConfig()
	client, err := gdrive.NewAPIClient(config)
	if err != nil {
		c.exitError(fmt.Sprintf("Failed to create Drive client: %v", err))
	}

	switch cmd {
	case "list", "ls":
		folderID := "root"
		if len(cmdArgs) > 0 {
			folderID = cmdArgs[0]
		}
		result, err := client.List(folderID, 50)
		if err != nil {
			c.exitError(err.Error())
		}
		if jsonOutput {
			c.outputJSON(result)
		} else {
			if len(result.Files) == 0 {
				fmt.Fprintln(c.stdout, "No files found.")
				return
			}
			fmt.Fprintf(c.stdout, "Files in %s:\n\n", folderID)
			for _, f := range result.Files {
				icon := "📄"
				if f.IsFolder() {
					icon = "📁"
				}
				fmt.Fprintf(c.stdout, "  %s %s\n", icon, f.Name)
				fmt.Fprintf(c.stdout, "      ID: %s\n", f.ID)
			}
		}

	case "get", "info":
		if len(cmdArgs) < 1 {
			c.exitError("FILE_ID required")
		}
		file, err := client.Get(cmdArgs[0])
		if err != nil {
			c.exitError(err.Error())
		}
		if jsonOutput {
			c.outputJSON(file)
		} else {
			fmt.Fprintf(c.stdout, "Name: %s\n", file.Name)
			fmt.Fprintf(c.stdout, "ID: %s\n", file.ID)
			fmt.Fprintf(c.stdout, "Type: %s\n", file.MimeType)
			if file.Size > 0 {
				fmt.Fprintf(c.stdout, "Size: %d bytes\n", file.Size)
			}
			if file.WebViewLink != "" {
				fmt.Fprintf(c.stdout, "Link: %s\n", file.WebViewLink)
			}
		}

	case "download":
		if len(cmdArgs) < 1 {
			c.exitError("FILE_ID required")
		}
		result, err := client.Download(cmdArgs[0])
		if err != nil {
			c.exitError(err.Error())
		}
		if len(cmdArgs) > 1 {
			// Write to file
			if err := os.WriteFile(cmdArgs[1], result.Content, 0644); err != nil {
				c.exitError(fmt.Sprintf("Failed to write file: %v", err))
			}
			fmt.Fprintf(c.stdout, "Downloaded to %s (%d bytes)\n", cmdArgs[1], len(result.Content))
		} else {
			// Write to stdout
			c.stdout.Write(result.Content)
		}

	case "upload":
		if len(cmdArgs) < 1 {
			c.exitError("FILE required")
		}
		filePath := cmdArgs[0]
		content, err := os.ReadFile(filePath)
		if err != nil {
			c.exitError(fmt.Sprintf("Failed to read file: %v", err))
		}
		parentID := getFlagValue(args, "--parent=")

		// Get filename from path
		name := filePath
		for i := len(filePath) - 1; i >= 0; i-- {
			if filePath[i] == '/' || filePath[i] == '\\' {
				name = filePath[i+1:]
				break
			}
		}

		result, err := client.Upload(name, content, "", parentID)
		if err != nil {
			c.exitError(err.Error())
		}
		if jsonOutput {
			c.outputJSON(result)
		} else {
			fmt.Fprintf(c.stdout, "Uploaded: %s\n", result.File.Name)
			fmt.Fprintf(c.stdout, "ID: %s\n", result.File.ID)
			if result.File.WebViewLink != "" {
				fmt.Fprintf(c.stdout, "Link: %s\n", result.File.WebViewLink)
			}
		}

	case "mkdir":
		if len(cmdArgs) < 1 {
			c.exitError("NAME required")
		}
		parentID := getFlagValue(args, "--parent=")
		result, err := client.CreateFolder(cmdArgs[0], parentID)
		if err != nil {
			c.exitError(err.Error())
		}
		if jsonOutput {
			c.outputJSON(result)
		} else {
			fmt.Fprintf(c.stdout, "Created folder: %s\n", result.File.Name)
			fmt.Fprintf(c.stdout, "ID: %s\n", result.File.ID)
		}

	case "rm", "delete":
		if len(cmdArgs) < 1 {
			c.exitError("FILE_ID required")
		}
		if err := client.Delete(cmdArgs[0]); err != nil {
			c.exitError(err.Error())
		}
		fmt.Fprintln(c.stdout, "Deleted successfully")

	default:
		fmt.Fprintf(c.stderr, "Unknown drive command: %s\n", cmd)
		c.printDriveUsage()
		os.Exit(1)
	}
}

func (c *cliContext) printDriveUsage() {
	fmt.Fprintln(c.stdout, `Usage: google drive <command> [arguments]

Commands:
  list [FOLDER_ID]              List files (default: root)
  get FILE_ID                   Get file metadata
  download FILE_ID [OUTPUT]     Download file content
  upload FILE [--parent=ID]     Upload a file
  mkdir NAME [--parent=ID]      Create a folder
  rm FILE_ID                    Delete a file

Options:
  --json    Output as JSON

Examples:
  google drive list
  google drive list 1abc123def
  google drive get 1abc123def
  google drive download 1abc123def output.txt
  google drive upload myfile.txt
  google drive upload myfile.txt --parent=1abc123def
  google drive mkdir "My Folder"
  google drive rm 1abc123def`)
}
