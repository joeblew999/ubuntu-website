package web

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/go-via/via"
	"github.com/go-via/via-plugin-picocss/picocss"
	"github.com/joeblew999/ubuntu-website/internal/env"
)

// ServeSetupGUI starts the web server for environment setup
func ServeSetupGUI() error {
	serveSetupGUIWithOptions(false)
	return nil
}

// ServeSetupGUIMock starts the web server in mock mode (no real validation)
func ServeSetupGUIMock() error {
	serveSetupGUIWithOptions(true)
	return nil
}

// cleanupPort kills any processes using the specified port
func cleanupPort(port string) {
	// Try to kill any existing processes on the port
	cmd := exec.Command("sh", "-c", fmt.Sprintf("lsof -ti:%s | xargs kill -9 2>/dev/null || true", port))
	_ = cmd.Run() // Ignore errors - port might not be in use
}

// serveSetupGUIWithOptions is the internal implementation using Via
func serveSetupGUIWithOptions(mockMode bool) {
	// Clean up port 3000 before starting
	cleanupPort("3000")

	// Use test env file in mock mode
	if mockMode {
		env.SetEnvFileForTesting(env.GetTestEnvFile())
		defer env.ResetEnvFile()
	}

	log.Printf("\n")
	title := "Environment Setup GUI"
	if mockMode {
		title += " (Mock Mode)"
	}
	log.Println(title)
	if mockMode {
		log.Println("Mock validation enabled - no real API calls")
	}
	log.Println("Opening in browser...")
	log.Printf("\n  %s\n\n", "http://localhost:3000")
	log.Println("Press Ctrl+C to stop")
	log.Println()

	v := via.New()
	v.Config(via.Options{
		DocumentTitle: "Environment Setup",
		Plugins:       []via.Plugin{picocss.Default},
		// DevMode enables the dataSPA Inspector debugging tool in the browser.
		// Set VIA_DEV_MODE=false to disable for production deployments.
		// Defaults to enabled when VIA_DEV_MODE is unset or any value other than "false".
		DevMode:       os.Getenv("VIA_DEV_MODE") != "false",
		LogLvl:        via.LogLevelWarn,  // Reduce noise from benign SSE race conditions
	})

	// Helper to load fresh config for each page request
	loadConfig := func() *env.EnvConfig {
		svc := env.NewService(mockMode)
		cfg, err := svc.GetCurrentConfig()
		if err != nil {
			log.Printf("Error loading config: %v", err)
			return &env.EnvConfig{}
		}
		return cfg
	}

	// Register routes - each loads fresh config
	v.Page("/", func(c *via.Context) {
		homePage(c, loadConfig(), mockMode)
	})

	// Cloudflare setup wizard - 4 steps
	v.Page("/cloudflare", func(c *via.Context) {
		cloudflarePage(c, loadConfig(), mockMode)
	})

	v.Page("/cloudflare/step2", func(c *via.Context) {
		cloudflareStep2Page(c, loadConfig(), mockMode)
	})

	v.Page("/cloudflare/step3", func(c *via.Context) {
		cloudflareStep3Page(c, loadConfig(), mockMode)
	})

	v.Page("/cloudflare/step4", func(c *via.Context) {
		cloudflareStep4Page(c, loadConfig(), mockMode)
	})

	v.Page("/claude", func(c *via.Context) {
		claudePage(c, loadConfig(), mockMode)
	})

	v.Start()
}
