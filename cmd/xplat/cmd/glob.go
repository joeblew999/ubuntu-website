package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/spf13/cobra"
)

// GlobCmd expands glob patterns
var GlobCmd = &cobra.Command{
	Use:   "glob <pattern>",
	Short: "Expand glob pattern",
	Long: `Expand a glob pattern and print matching files.

Supports doublestar (**) patterns for recursive matching.
On Windows, matching is case-insensitive by default.

Patterns:
  *        - matches any sequence of characters (not including /)
  **       - matches any sequence including /
  ?        - matches any single character
  [abc]    - matches one of the characters
  {a,b}    - matches either 'a' or 'b'

Examples:
  xplat glob "taskfiles/*.yml"
  xplat glob "taskfiles/Taskfile.*.yml"
  xplat glob "**/*.go"
  xplat glob "src/**/*.{ts,tsx}"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pattern := args[0]

		// Build options - case insensitive on Windows
		var opts []doublestar.GlobOption
		if runtime.GOOS == "windows" {
			opts = append(opts, doublestar.WithCaseInsensitive())
		}

		// Use the OS filesystem
		matches, err := doublestar.Glob(os.DirFS("."), pattern, opts...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		for _, match := range matches {
			fmt.Println(match)
		}
	},
}
