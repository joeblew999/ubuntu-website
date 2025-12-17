// Package playwrightcli provides the CLI entry point for the playwright browser automation tool.
package playwrightcli

import (
	"flag"
	"fmt"
	"io"
	"time"

	"github.com/joeblew999/ubuntu-website/internal/browser"
	"github.com/joeblew999/ubuntu-website/internal/browser/sites"
	"github.com/joeblew999/ubuntu-website/internal/google/gmail"
)

// Run executes the playwright CLI with the given arguments.
// Returns exit code: 0 for success, non-zero for errors.
func Run(args []string, version string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("playwright", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var (
		headless    = fs.Bool("headless", false, "Run browser in headless mode")
		timeout     = fs.Int("timeout", 120, "Timeout in seconds")
		port        = fs.Int("port", 8085, "Callback server port")
		profile     = fs.String("profile", "", "Browser profile directory (reuses logins)")
		useChrome   = fs.Bool("chrome", false, "Use system Chrome instead of Chromium")
		browserFlag = fs.String("browser", "chromium", "Browser engine: chromium (default), firefox, webkit")
		showVersion = fs.Bool("version", false, "Show version")
	)

	if err := fs.Parse(args[1:]); err != nil {
		return 1
	}

	if *showVersion {
		fmt.Fprintf(stdout, "playwright %s\n", version)
		return 0
	}

	cmdArgs := fs.Args()
	if len(cmdArgs) == 0 {
		printUsage(stdout)
		return 1
	}

	app := &cliApp{
		headless:  *headless,
		timeout:   *timeout,
		port:      *port,
		profile:   *profile,
		useChrome: *useChrome,
		browser:   *browserFlag,
		stdout:    stdout,
		stderr:    stderr,
	}

	if err := app.run(cmdArgs); err != nil {
		fmt.Fprintf(stderr, "ERROR: %v\n", err)
		return 1
	}
	return 0
}

// cliApp holds CLI configuration
type cliApp struct {
	headless  bool
	timeout   int
	port      int
	profile   string
	useChrome bool
	browser   string
	stdout    io.Writer
	stderr    io.Writer
}

// browserEngine returns the browser engine from flag
func (a *cliApp) browserEngine() browser.BrowserEngine {
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
func (a *cliApp) browserConfig() *browser.PlaywrightConfig {
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

// run dispatches to the appropriate command handler
func (a *cliApp) run(args []string) error {
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

func (a *cliApp) runOAuth(args []string) error {
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

func (a *cliApp) runOpen(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("open requires a URL: playwright open <url>")
	}
	return browser.OpenURLInPlaywright(args[0], time.Duration(a.timeout)*time.Second, a.browserConfig())
}

func (a *cliApp) runScreenshot(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("screenshot requires URL and filename: playwright screenshot <url> <file>")
	}
	return browser.TakeScreenshot(args[0], args[1], a.browserConfig())
}

func (a *cliApp) runInstall() error {
	return browser.InstallBrowsers(a.browserEngine())
}

// =============================================================================
// Gmail Commands - Delegate to internal/google/gmail
// =============================================================================

func (a *cliApp) runGmailSendAs(args []string) error {
	if len(args) < 5 {
		printGmailSendAsUsage(a.stdout)
		return nil
	}
	config := gmail.ParseSMTPConfig(args[0], args[1], args[2], args[3], args[4])
	fmt.Fprintln(a.stdout, "🔧 Gmail 'Send mail as' Configuration")
	fmt.Fprintln(a.stdout, "=====================================")
	fmt.Fprintf(a.stdout, "  Name:     %s\n", config.Name)
	fmt.Fprintf(a.stdout, "  Email:    %s\n", config.Email)
	fmt.Fprintf(a.stdout, "  SMTP:     %s:%s\n", config.SMTPHost, config.SMTPPort)
	fmt.Fprintf(a.stdout, "  Username: %s\n", config.SMTPUsername)
	fmt.Fprintln(a.stdout)

	automation := gmail.NewSettingsAutomation(config)
	return automation.ConfigureSendAs()
}

func (a *cliApp) runGmailSettings() error {
	fmt.Fprintln(a.stdout, "Opening Gmail Settings > Accounts...")
	fmt.Fprintln(a.stdout)
	automation := gmail.NewSettingsAutomation(&gmail.SMTPConfig{})
	return automation.OpenSettingsPage()
}

func (a *cliApp) runGmailLogout() error {
	fmt.Fprintln(a.stdout, "Logging out of Google accounts...")
	fmt.Fprintln(a.stdout)
	automation := gmail.NewSettingsAutomation(&gmail.SMTPConfig{})
	return automation.Logout()
}

func (a *cliApp) runGmailSwitch() error {
	fmt.Fprintln(a.stdout, "Opening Google Account chooser...")
	fmt.Fprintln(a.stdout)
	automation := gmail.NewSettingsAutomation(&gmail.SMTPConfig{})
	return automation.SwitchAccount()
}

// =============================================================================
// Cloudflare Commands - Delegate to internal/browser/sites
// =============================================================================

func (a *cliApp) runCloudflare(args []string) error {
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

func (a *cliApp) runCloudflareEmail(args []string) error {
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

func (a *cliApp) runCloudflareEmailAdd(args []string) error {
	if len(args) < 2 {
		fmt.Fprintln(a.stdout, "Usage: playwright cloudflare-email-add <from-email> <to-email> [domain]")
		fmt.Fprintln(a.stdout)
		fmt.Fprintln(a.stdout, "Examples:")
		fmt.Fprintln(a.stdout, "  playwright cloudflare-email-add contact@ubuntusoftware.net gedw99@gmail.com")
		fmt.Fprintln(a.stdout, "  playwright cloudflare-email-add info@example.com me@gmail.com example.com")
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

func (a *cliApp) runCloudflareToken(args []string) error {
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

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "playwright - Reusable browser automation CLI")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  playwright [flags] <command> [args]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Commands:")
	fmt.Fprintln(w, "  oauth <url>              Start OAuth flow with callback server")
	fmt.Fprintln(w, "  open <url>               Open URL in Playwright browser")
	fmt.Fprintln(w, "  screenshot <url> <file>  Take screenshot of URL")
	fmt.Fprintln(w, "  install                  Install Playwright browsers")
	fmt.Fprintln(w, "  gmail-sendas <args>      Configure Gmail 'Send mail as' with SMTP relay")
	fmt.Fprintln(w, "  gmail-settings           Open Gmail settings page (manual setup)")
	fmt.Fprintln(w, "  gmail-logout             Logout of Google accounts (clear session)")
	fmt.Fprintln(w, "  gmail-switch             Open Google account chooser (add/switch)")
	fmt.Fprintln(w, "  cloudflare [page]        Open Cloudflare dashboard (email|pages|dns)")
	fmt.Fprintln(w, "  cloudflare-email [domain] Open Cloudflare Email Routing")
	fmt.Fprintln(w, "  cloudflare-email-add <from> <to> Add email forwarding rule (HIL)")
	fmt.Fprintln(w, "  cloudflare-token         Setup Cloudflare API token (guided)")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Flags:")
	fmt.Fprintln(w, "  -browser=ENGINE  Browser engine: chromium (default), firefox, webkit")
	fmt.Fprintln(w, "  -headless        Run browser in headless mode")
	fmt.Fprintln(w, "  -timeout=120     Timeout in seconds (default: 120)")
	fmt.Fprintln(w, "  -port=8085       Callback server port (default: 8085)")
	fmt.Fprintln(w, "  -profile=DIR     Browser profile directory (reuses logins)")
	fmt.Fprintln(w, "  -chrome          Use system Chrome instead of Chromium")
	fmt.Fprintln(w, "  -version         Show version")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Examples:")
	fmt.Fprintln(w, "  playwright oauth 'https://accounts.google.com/...'")
	fmt.Fprintln(w, "  playwright -headless screenshot https://example.com shot.png")
	fmt.Fprintln(w, "  playwright install")
	fmt.Fprintln(w, "  playwright gmail-sendas 'Name' 'email@domain.com' smtp2go 'user' 'pass'")
	fmt.Fprintln(w, "  playwright -profile=~/.my-profile cloudflare email")
}

func printGmailSendAsUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage: playwright gmail-sendas <name> <email> <provider> <user> <pass>")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Pre-configured providers:")
	fmt.Fprintln(w, "  smtp2go - mail.smtp2go.com:587")
	fmt.Fprintln(w, "  brevo   - smtp-relay.brevo.com:587")
	fmt.Fprintln(w, "  resend  - smtp.resend.com:587")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Or use custom SMTP host as provider.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Examples:")
	fmt.Fprintln(w, "  playwright gmail-sendas 'Gerard Webb' 'gerard@domain.com' smtp2go 'user' 'pass'")
	fmt.Fprintln(w, "  playwright gmail-sendas 'Name' 'me@domain.com' 'smtp.example.com' 'user' 'pass'")
}
