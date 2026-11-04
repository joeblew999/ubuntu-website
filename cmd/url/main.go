// Command url opens the site's pages and dashboards in the browser, on any OS: the local dev server (mise run url:dev),
// production, the Cloudflare preview of main and of this commit, Google Search Console and Cloudflare Web Analytics.
// Each command is a file of its own and a task of .mise/tasks/url.toml (url:<command>).
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joeblew999/ubuntu-website/internal/cli"
)

const (
	cfAccountID = "7384af54e33b8a54ff240371ea368440"
	cfSiteTag   = "4c28a6bfb5514996914a603c999d5c79" // Cloudflare Web Analytics
)

func main() {
	if len(os.Args) < 2 {
		fail("usage: url dev | prod | preview | commit | search-console | cf-analytics")
	}
	switch cmd := os.Args[1]; cmd {
	case "dev":
		dev()
	case "prod":
		prod()
	case "preview":
		preview()
	case "commit":
		commit()
	case "search-console":
		searchConsole()
	case "cf-analytics":
		cfAnalytics()
	default:
		fail("no command " + strconv.Quote(cmd))
	}
}

// fail says what went wrong and exits 1.
func fail(msg string) {
	fmt.Println("✗ " + msg)
	os.Exit(1)
}

// open says the URL and opens it in the browser.
func open(url string) {
	fmt.Println(url)
	if err := cli.OpenBrowser(url); err != nil {
		fail(err.Error())
	}
}
