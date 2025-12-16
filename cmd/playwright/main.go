// playwright - Reusable browser automation CLI
//
// A thin CLI wrapper around browser automation packages.
// Domain logic lives in internal/browser and internal/google/gmail.
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
//	-profile=DIR     Browser profile directory (reuses logins)
//	-browser=ENGINE  Browser engine: chromium, firefox, webkit
//	-chrome          Use system Chrome instead of bundled Chromium
//	-version         Show version
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/joeblew999/ubuntu-website/internal/browser"
	"github.com/joeblew999/ubuntu-website/internal/browser/sites"
	"github.com/joeblew999/ubuntu-website/internal/google/gmail"
)

var version = "dev"

func main() {
	// Global flags
	var (
		headless    = flag.Bool("headless", false, "Run browser in headless mode")
		timeout     = flag.Int("timeout", 120, "Timeout in seconds")
		port        = flag.Int("port", 8085, "Callback server port")
		profile     = flag.String("profile", "", "Browser profile directory (reuses logins)")
		useChrome   = flag.Bool("chrome", false, "Use system Chrome instead of Chromium")
		browserFlag = flag.String("browser", "chromium", "Browser engine: chromium (default), firefox, webkit")
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
		headless:  *headless,
		timeout:   *timeout,
		port:      *port,
		profile:   *profile,
		useChrome: *useChrome,
		browser:   *browserFlag,
	}

	if err := app.Run(args); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}

// App holds CLI configuration
type App struct {
	headless  bool
	timeout   int
	port      int
	profile   string
	useChrome bool
	browser   string
}

// browserEngine returns the browser engine from flag
func (a *App) browserEngine() browser.BrowserEngine {
	switch a.browser {
	case "firefox":
		return browser.BrowserFirefox
	case "webkit", "safari":
		return browser.BrowserWebKit
	default:
		return browser.BrowserChromium
	}
}

// browserConfig creates a browser config based on app flags
func (a *App) browserConfig() *browser.PlaywrightConfig {
	config := &browser.PlaywrightConfig{
		Engine:      a.browserEngine(),
		Headless:    a.headless,
		Timeout:     float64(a.timeout * 1000),
		UserDataDir: a.profile,
	}
	if a.useChrome && config.Engine == browser.BrowserChromium {
		config.Channel = "chrome"
	}
	return config
}

// automationConfig creates automation config for HIL flows
func (a *App) automationConfig() *browser.AutomationConfig {
	return &browser.AutomationConfig{
		Engine:    a.browserEngine(),
		Profile:   a.profile,
		Headless:  a.headless,
		Timeout:   time.Duration(a.timeout) * time.Second,
		Verbose:   true,
		UseChrome: a.useChrome,
	}
}

// Run dispatches to the appropriate command handler
func (a *App) Run(args []string) error {
	cmd := args[0]
	subArgs := args[1:]

	switch cmd {
	// Core browser commands
	case "oauth":
		return a.runOAuth(subArgs)
	case "open":
		return a.runOpen(subArgs)
	case "screenshot":
		return a.runScreenshot(subArgs)
	case "install":
		return a.runInstall()

	// Gmail commands - delegate to internal/google/gmail
	case "gmail-sendas":
		return a.runGmailSendAs(subArgs)
	case "gmail-settings":
		return a.runGmailSettings()
	case "gmail-logout":
		return a.runGmailLogout()
	case "gmail-switch":
		return a.runGmailSwitch()

	// Cloudflare commands - delegate to internal/browser/sites
	case "cloudflare":
		return a.runCloudflare(subArgs)
	case "cloudflare-email":
		return a.runCloudflareEmail(subArgs)
	case "cloudflare-email-add":
		return a.runCloudflareEmailAdd(subArgs)
	case "cloudflare-token":
		return a.runCloudflareToken(subArgs)

	default:
		return fmt.Errorf("unknown command: %s\nRun 'playwright' for usage", cmd)
	}
}

// =============================================================================
// Core Browser Commands
// =============================================================================

func (a *App) runOAuth(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("oauth requires a URL: playwright oauth <url>")
	}
	config := &browser.CLIOAuthConfig{
		URL:     args[0],
		Port:    a.port,
		Timeout: time.Duration(a.timeout) * time.Second,
		Browser: a.browserConfig(),
	}
	_, err := browser.RunCLIOAuthFlow(config)
	return err
}

func (a *App) runOpen(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("open requires a URL: playwright open <url>")
	}
	return browser.OpenURLInPlaywright(args[0], time.Duration(a.timeout)*time.Second, a.browserConfig())
}

func (a *App) runScreenshot(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("screenshot requires URL and filename: playwright screenshot <url> <file>")
	}
	return browser.TakeScreenshot(args[0], args[1], a.browserConfig())
}

func (a *App) runInstall() error {
	return browser.InstallBrowsers(a.browserEngine())
}

// =============================================================================
// Gmail Commands - Delegate to internal/google/gmail
// =============================================================================

func (a *App) runGmailSendAs(args []string) error {
	if len(args) < 5 {
		printGmailSendAsUsage()
		return nil
	}
	config := gmail.ParseSMTPConfig(args[0], args[1], args[2], args[3], args[4])
	fmt.Println("🔧 Gmail 'Send mail as' Configuration")
	fmt.Println("=====================================")
	fmt.Printf("  Name:     %s\n", config.Name)
	fmt.Printf("  Email:    %s\n", config.Email)
	fmt.Printf("  SMTP:     %s:%s\n", config.SMTPHost, config.SMTPPort)
	fmt.Printf("  Username: %s\n", config.SMTPUsername)
	fmt.Println()

	automation := gmail.NewSettingsAutomation(config)
	return automation.ConfigureSendAs()
}

func (a *App) runGmailSettings() error {
	fmt.Println("Opening Gmail Settings > Accounts...")
	fmt.Println()
	automation := gmail.NewSettingsAutomation(&gmail.SMTPConfig{})
	return automation.OpenSettingsPage()
}

func (a *App) runGmailLogout() error {
	fmt.Println("Logging out of Google accounts...")
	fmt.Println()
	automation := gmail.NewSettingsAutomation(&gmail.SMTPConfig{})
	return automation.Logout()
}

func (a *App) runGmailSwitch() error {
	fmt.Println("Opening Google Account chooser...")
	fmt.Println()
	automation := gmail.NewSettingsAutomation(&gmail.SMTPConfig{})
	return automation.SwitchAccount()
}

// =============================================================================
// Cloudflare Commands - Delegate to internal/browser/sites
// =============================================================================

func (a *App) runCloudflare(args []string) error {
	domain := "ubuntusoftware.net"
	page := ""
	if len(args) > 0 {
		page = args[0]
	}
	if len(args) > 1 {
		domain = args[1]
	}

	config := sites.DefaultCloudflareConfig(domain)
	if a.profile != "" {
		config.Profile = a.profile
	}
	cf := sites.NewCloudflareAutomation(config)
	return cf.OpenDashboard(page)
}

func (a *App) runCloudflareEmail(args []string) error {
	domain := "ubuntusoftware.net"
	if len(args) > 0 {
		domain = args[0]
	}

	config := sites.DefaultCloudflareConfig(domain)
	if a.profile != "" {
		config.Profile = a.profile
	}
	cf := sites.NewCloudflareAutomation(config)
	return cf.OpenEmailRouting()
}

func (a *App) runCloudflareEmailAdd(args []string) error {
	if len(args) < 2 {
		fmt.Println("Usage: playwright cloudflare-email-add <from-email> <to-email> [domain]")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  playwright cloudflare-email-add contact@ubuntusoftware.net gedw99@gmail.com")
		fmt.Println("  playwright cloudflare-email-add info@example.com me@gmail.com example.com")
		return nil
	}

	domain := "ubuntusoftware.net"
	if len(args) > 2 {
		domain = args[2]
	}

	config := sites.DefaultCloudflareConfig(domain)
	if a.profile != "" {
		config.Profile = a.profile
	}
	cf := sites.NewCloudflareAutomation(config)
	return cf.AddEmailRoutingRule(&sites.EmailRoutingRule{
		FromAddress: args[0],
		ToAddress:   args[1],
	})
}

func (a *App) runCloudflareToken(args []string) error {
	domain := "ubuntusoftware.net"
	envFile := ".env"
	if len(args) > 0 {
		envFile = args[0]
	}

	config := sites.DefaultCloudflareConfig(domain)
	if a.profile != "" {
		config.Profile = a.profile
	}
	cf := sites.NewCloudflareAutomation(config)
	return cf.SetupAPIToken(envFile)
}

// =============================================================================
// Usage
// =============================================================================

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
	fmt.Println("  gmail-sendas <args>      Configure Gmail 'Send mail as' with SMTP relay")
	fmt.Println("  gmail-settings           Open Gmail settings page (manual setup)")
	fmt.Println("  gmail-logout             Logout of Google accounts (clear session)")
	fmt.Println("  gmail-switch             Open Google account chooser (add/switch)")
	fmt.Println("  cloudflare [page]        Open Cloudflare dashboard (email|pages|dns)")
	fmt.Println("  cloudflare-email [domain] Open Cloudflare Email Routing")
	fmt.Println("  cloudflare-email-add <from> <to> Add email forwarding rule (HIL)")
	fmt.Println("  cloudflare-token         Setup Cloudflare API token (guided)")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -browser=ENGINE  Browser engine: chromium (default), firefox, webkit")
	fmt.Println("  -headless        Run browser in headless mode")
	fmt.Println("  -timeout=120     Timeout in seconds (default: 120)")
	fmt.Println("  -port=8085       Callback server port (default: 8085)")
	fmt.Println("  -profile=DIR     Browser profile directory (reuses logins)")
	fmt.Println("  -chrome          Use system Chrome instead of Chromium")
	fmt.Println("  -version         Show version")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  playwright oauth 'https://accounts.google.com/...'")
	fmt.Println("  playwright -headless screenshot https://example.com shot.png")
	fmt.Println("  playwright install")
	fmt.Println("  playwright gmail-sendas 'Name' 'email@domain.com' smtp2go 'user' 'pass'")
	fmt.Println("  playwright -profile=~/.my-profile cloudflare email")
}

func printGmailSendAsUsage() {
	fmt.Println("Usage: playwright gmail-sendas <name> <email> <provider> <user> <pass>")
	fmt.Println()
	fmt.Println("Pre-configured providers:")
	fmt.Println("  smtp2go - mail.smtp2go.com:587")
	fmt.Println("  brevo   - smtp-relay.brevo.com:587")
	fmt.Println("  resend  - smtp.resend.com:587")
	fmt.Println()
	fmt.Println("Or use custom SMTP host as provider.")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  playwright gmail-sendas 'Gerard Webb' 'gerard@domain.com' smtp2go 'user' 'pass'")
	fmt.Println("  playwright gmail-sendas 'Name' 'me@domain.com' 'smtp.example.com' 'user' 'pass'")
}
