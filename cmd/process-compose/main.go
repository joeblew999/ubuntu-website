// Process-Compose wrapper
//
// This is a thin wrapper around the process-compose library.
// We wrap it so we can use xplat binary:install for consistent
// cross-platform installation across all project tools.
//
// Upstream: https://github.com/F1bonacc1/process-compose
package main

import (
	"github.com/f1bonacc1/process-compose/src/cmd"
)

func main() {
	cmd.Execute()
}
