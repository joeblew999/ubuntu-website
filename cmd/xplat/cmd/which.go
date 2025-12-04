package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// WhichCmd finds a binary in PATH
var WhichCmd = &cobra.Command{
	Use:   "which <binary>",
	Short: "Find binary in PATH",
	Long: `Find the path to an executable in PATH.

Works identically on macOS, Linux, and Windows.
Returns exit code 1 if the binary is not found.

Examples:
  xplat which go
  xplat which task
  xplat which node`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path, err := exec.LookPath(args[0])
		if err != nil {
			os.Exit(1)
		}
		fmt.Println(path)
	},
}
