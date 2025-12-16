package main

import (
	"fmt"
	"os"

	"github.com/joeblew999/ubuntu-website/internal/google/gdrive"
)

func handleDrive(args []string) {
	if len(args) < 1 {
		printDriveUsage()
		return
	}

	cmd := args[0]
	cmdArgs := args[1:]
	jsonOutput := hasFlag(cmdArgs, "--json")
	cmdArgs = filterFlags(cmdArgs)

	config := gdrive.DefaultConfig()
	client, err := gdrive.NewAPIClient(config)
	if err != nil {
		exitError(fmt.Sprintf("Failed to create Drive client: %v", err))
	}

	switch cmd {
	case "list", "ls":
		folderID := "root"
		if len(cmdArgs) > 0 {
			folderID = cmdArgs[0]
		}
		result, err := client.List(folderID, 50)
		if err != nil {
			exitError(err.Error())
		}
		if jsonOutput {
			outputJSON(result)
		} else {
			if len(result.Files) == 0 {
				fmt.Println("No files found.")
				return
			}
			fmt.Printf("Files in %s:\n\n", folderID)
			for _, f := range result.Files {
				icon := "📄"
				if f.IsFolder() {
					icon = "📁"
				}
				fmt.Printf("  %s %s\n", icon, f.Name)
				fmt.Printf("      ID: %s\n", f.ID)
			}
		}

	case "get", "info":
		if len(cmdArgs) < 1 {
			exitError("FILE_ID required")
		}
		file, err := client.Get(cmdArgs[0])
		if err != nil {
			exitError(err.Error())
		}
		if jsonOutput {
			outputJSON(file)
		} else {
			fmt.Printf("Name: %s\n", file.Name)
			fmt.Printf("ID: %s\n", file.ID)
			fmt.Printf("Type: %s\n", file.MimeType)
			if file.Size > 0 {
				fmt.Printf("Size: %d bytes\n", file.Size)
			}
			if file.WebViewLink != "" {
				fmt.Printf("Link: %s\n", file.WebViewLink)
			}
		}

	case "download":
		if len(cmdArgs) < 1 {
			exitError("FILE_ID required")
		}
		result, err := client.Download(cmdArgs[0])
		if err != nil {
			exitError(err.Error())
		}
		if len(cmdArgs) > 1 {
			// Write to file
			if err := os.WriteFile(cmdArgs[1], result.Content, 0644); err != nil {
				exitError(fmt.Sprintf("Failed to write file: %v", err))
			}
			fmt.Printf("Downloaded to %s (%d bytes)\n", cmdArgs[1], len(result.Content))
		} else {
			// Write to stdout
			os.Stdout.Write(result.Content)
		}

	case "upload":
		if len(cmdArgs) < 1 {
			exitError("FILE required")
		}
		filePath := cmdArgs[0]
		content, err := os.ReadFile(filePath)
		if err != nil {
			exitError(fmt.Sprintf("Failed to read file: %v", err))
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
			exitError(err.Error())
		}
		if jsonOutput {
			outputJSON(result)
		} else {
			fmt.Printf("Uploaded: %s\n", result.File.Name)
			fmt.Printf("ID: %s\n", result.File.ID)
			if result.File.WebViewLink != "" {
				fmt.Printf("Link: %s\n", result.File.WebViewLink)
			}
		}

	case "mkdir":
		if len(cmdArgs) < 1 {
			exitError("NAME required")
		}
		parentID := getFlagValue(args, "--parent=")
		result, err := client.CreateFolder(cmdArgs[0], parentID)
		if err != nil {
			exitError(err.Error())
		}
		if jsonOutput {
			outputJSON(result)
		} else {
			fmt.Printf("Created folder: %s\n", result.File.Name)
			fmt.Printf("ID: %s\n", result.File.ID)
		}

	case "rm", "delete":
		if len(cmdArgs) < 1 {
			exitError("FILE_ID required")
		}
		if err := client.Delete(cmdArgs[0]); err != nil {
			exitError(err.Error())
		}
		fmt.Println("Deleted successfully")

	default:
		fmt.Fprintf(os.Stderr, "Unknown drive command: %s\n", cmd)
		printDriveUsage()
		os.Exit(1)
	}
}

func printDriveUsage() {
	fmt.Println(`Usage: google drive <command> [arguments]

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
