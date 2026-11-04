// Command google-mcp sets up and checks the Google MCP server Claude Code runs from .mcp.json
// (mise exec -- fnox exec -- google-mcp-server): the OAuth client in the keychain (fnox.toml), the accounts signed in
// (~/.google-mcp-accounts) and the Google Cloud screens. Each command is a file of its own and a task of mise.toml
// (google:<command>).
package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fail("usage: google-mcp secret <key> | console | rotate | auth | check")
	}
	switch cmd, args := os.Args[1], os.Args[2:]; cmd {
	case "secret":
		secret(args)
	case "console":
		console()
	case "rotate":
		rotate()
	case "auth":
		auth()
	case "check":
		check()
	default:
		fail("no command " + strconv.Quote(cmd))
	}
}

// fail says what went wrong and exits 1.
func fail(msg string) {
	fmt.Println("✗ " + msg)
	os.Exit(1)
}

// seconds is the whole number of seconds env says, or def.
func seconds(env string, def int) int {
	if n, err := strconv.Atoi(os.Getenv(env)); err == nil && n > 0 {
		return n
	}
	return def
}

// browse opens url, unless GOOGLE_NO_BROWSER is set.
func browse(url string) {
	if os.Getenv("GOOGLE_NO_BROWSER") == "" {
		openURL(url)
	}
}
