// playwright - Reusable browser automation CLI
//
// This tool provides a reusable Playwright-based browser automation
// that can be used for OAuth flows, web scraping, and testing.
//
// Usage:
//
//	playwright oauth <url>             Start OAuth flow with callback server
//	playwright open <url>              Open URL in Playwright browser
//	playwright screenshot <url> <file> Take screenshot of URL
//	playwright install                 Install Playwright browsers
//
// Flags:
//
//	-headless        Run browser in headless mode
//	-timeout=120     Timeout in seconds (default: 120)
//	-port=8085       Callback server port (default: 8085)
//	-version         Show version
//
// The oauth command starts a local callback server and opens the URL
// in a Playwright-controlled browser. When the callback receives a
// response, it extracts the code/token and returns it.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

var version = "dev"

func main() {
	// Global flags
	var (
		headless    = flag.Bool("headless", false, "Run browser in headless mode")
		timeout     = flag.Int("timeout", 120, "Timeout in seconds")
		port        = flag.Int("port", 8085, "Callback server port")
		showVersion = flag.Bool("version", false, "Show version")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("playwright %s\n", version)
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	app := &App{
		headless: *headless,
		timeout:  *timeout,
		port:     *port,
	}

	if err := app.Run(args); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}

type App struct {
	headless bool
	timeout  int
	port     int
}

func (a *App) Run(args []string) error {
	cmd := args[0]
	subArgs := args[1:]

	switch cmd {
	case "oauth":
		if len(subArgs) < 1 {
			return fmt.Errorf("oauth requires a URL: playwright oauth <url>")
		}
		return a.runOAuth(subArgs[0])
	case "open":
		if len(subArgs) < 1 {
			return fmt.Errorf("open requires a URL: playwright open <url>")
		}
		return a.runOpen(subArgs[0])
	case "screenshot":
		if len(subArgs) < 2 {
			return fmt.Errorf("screenshot requires URL and filename: playwright screenshot <url> <file>")
		}
		return a.runScreenshot(subArgs[0], subArgs[1])
	case "install":
		return a.runInstall()
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}

// OAuthResult contains the result of an OAuth flow
type OAuthResult struct {
	Code  string            `json:"code,omitempty"`
	Token string            `json:"token,omitempty"`
	Error string            `json:"error,omitempty"`
	Query map[string]string `json:"query,omitempty"`
}

// runOAuth starts an OAuth flow with a callback server
func (a *App) runOAuth(authURL string) error {
	fmt.Println("Starting OAuth flow...")
	fmt.Printf("  URL: %s\n", authURL)
	fmt.Printf("  Port: %d\n", a.port)
	fmt.Printf("  Timeout: %ds\n", a.timeout)
	fmt.Println()

	// Channels for results
	resultChan := make(chan OAuthResult, 1)
	errChan := make(chan error, 1)

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.timeout)*time.Second)
	defer cancel()

	// Start callback server
	server := &http.Server{Addr: fmt.Sprintf(":%d", a.port)}
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		result := OAuthResult{
			Query: make(map[string]string),
		}

		// Extract all query parameters
		for key, values := range r.URL.Query() {
			if len(values) > 0 {
				result.Query[key] = values[0]
			}
		}

		// Check for common OAuth parameters
		if code := r.URL.Query().Get("code"); code != "" {
			result.Code = code
		}
		if token := r.URL.Query().Get("access_token"); token != "" {
			result.Token = token
		}
		if errMsg := r.URL.Query().Get("error"); errMsg != "" {
			result.Error = errMsg
		}

		// Send success response
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head><title>Success</title></head>
<body style="font-family: -apple-system, sans-serif; padding: 40px; text-align: center;">
<h1 style="color: #22c55e;">Authentication Complete</h1>
<p>You can close this window and return to the terminal.</p>
</body>
</html>`)

		select {
		case resultChan <- result:
		default:
		}
	})
	server.Handler = mux

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			select {
			case errChan <- fmt.Errorf("server error: %w", err):
			default:
			}
		}
	}()

	// Install Playwright if needed
	if err := playwright.Install(&playwright.RunOptions{
		Browsers: []string{"chromium"},
		Verbose:  false,
	}); err != nil {
		return fmt.Errorf("failed to install playwright: %w", err)
	}

	// Start Playwright
	pw, err := playwright.Run()
	if err != nil {
		return fmt.Errorf("could not start playwright: %w", err)
	}
	defer pw.Stop()

	// Launch browser
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(a.headless),
	})
	if err != nil {
		return fmt.Errorf("could not launch browser: %w", err)
	}
	defer browser.Close()

	// Create page and navigate
	page, err := browser.NewPage()
	if err != nil {
		return fmt.Errorf("could not create page: %w", err)
	}

	fmt.Println("Opening browser...")
	if _, err := page.Goto(authURL); err != nil {
		return fmt.Errorf("could not navigate: %w", err)
	}

	fmt.Println("Waiting for callback...")

	// Wait for result, error, or timeout
	select {
	case result := <-resultChan:
		server.Shutdown(ctx)
		browser.Close()

		if result.Error != "" {
			return fmt.Errorf("OAuth error: %s", result.Error)
		}

		// Output result as JSON for parsing by other tools
		output, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(output))
		return nil

	case err := <-errChan:
		browser.Close()
		server.Shutdown(ctx)
		return err

	case <-ctx.Done():
		browser.Close()
		server.Shutdown(context.Background())
		return fmt.Errorf("timeout after %d seconds", a.timeout)
	}
}

// runOpen opens a URL in Playwright browser (useful for debugging)
func (a *App) runOpen(url string) error {
	fmt.Printf("Opening %s...\n", url)

	// Install if needed
	if err := playwright.Install(&playwright.RunOptions{
		Browsers: []string{"chromium"},
		Verbose:  false,
	}); err != nil {
		return fmt.Errorf("failed to install playwright: %w", err)
	}

	pw, err := playwright.Run()
	if err != nil {
		return fmt.Errorf("could not start playwright: %w", err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(a.headless),
	})
	if err != nil {
		return fmt.Errorf("could not launch browser: %w", err)
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		return fmt.Errorf("could not create page: %w", err)
	}

	if _, err := page.Goto(url); err != nil {
		return fmt.Errorf("could not navigate: %w", err)
	}

	fmt.Println("Browser open. Press Ctrl+C to close.")

	// Wait for timeout or interrupt
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.timeout)*time.Second)
	defer cancel()
	<-ctx.Done()

	return nil
}

// runScreenshot takes a screenshot of a URL
func (a *App) runScreenshot(url, filename string) error {
	fmt.Printf("Taking screenshot of %s...\n", url)

	// Install if needed
	if err := playwright.Install(&playwright.RunOptions{
		Browsers: []string{"chromium"},
		Verbose:  false,
	}); err != nil {
		return fmt.Errorf("failed to install playwright: %w", err)
	}

	pw, err := playwright.Run()
	if err != nil {
		return fmt.Errorf("could not start playwright: %w", err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true), // Always headless for screenshots
	})
	if err != nil {
		return fmt.Errorf("could not launch browser: %w", err)
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		return fmt.Errorf("could not create page: %w", err)
	}

	if _, err := page.Goto(url); err != nil {
		return fmt.Errorf("could not navigate: %w", err)
	}

	// Wait for page to load
	page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateNetworkidle,
	})

	// Ensure filename has extension
	if !strings.HasSuffix(filename, ".png") && !strings.HasSuffix(filename, ".jpg") {
		filename = filename + ".png"
	}

	if _, err := page.Screenshot(playwright.PageScreenshotOptions{
		Path:     playwright.String(filename),
		FullPage: playwright.Bool(true),
	}); err != nil {
		return fmt.Errorf("could not take screenshot: %w", err)
	}

	fmt.Printf("Screenshot saved to %s\n", filename)
	return nil
}

// runInstall installs Playwright browsers
func (a *App) runInstall() error {
	fmt.Println("Installing Playwright browsers...")

	if err := playwright.Install(&playwright.RunOptions{
		Browsers: []string{"chromium"},
		Verbose:  true,
	}); err != nil {
		return fmt.Errorf("failed to install: %w", err)
	}

	fmt.Println("Installation complete!")
	return nil
}

func printUsage() {
	fmt.Println("playwright - Reusable browser automation CLI")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  playwright [flags] <command> [args]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  oauth <url>              Start OAuth flow with callback server")
	fmt.Println("  open <url>               Open URL in Playwright browser")
	fmt.Println("  screenshot <url> <file>  Take screenshot of URL")
	fmt.Println("  install                  Install Playwright browsers")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -headless        Run browser in headless mode")
	fmt.Println("  -timeout=120     Timeout in seconds (default: 120)")
	fmt.Println("  -port=8085       Callback server port (default: 8085)")
	fmt.Println("  -version         Show version")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  playwright oauth 'https://accounts.google.com/...'")
	fmt.Println("  playwright -headless screenshot https://example.com shot.png")
	fmt.Println("  playwright install")
}
