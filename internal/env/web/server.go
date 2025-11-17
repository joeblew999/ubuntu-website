package web

import (
	"fmt"
	"log"
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

	// Use service to load config
	svc := env.NewService(mockMode)
	cfg, err := svc.GetCurrentConfig()
	if err != nil {
		log.Printf("Error loading env: %v", err)
		return
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
	})

	// Register routes
	v.Page("/", func(c *via.Context) {
		homePage(c, cfg, mockMode)
	})

	v.Page("/cloudflare", func(c *via.Context) {
		cloudflarePage(c, cfg, mockMode)
	})

	v.Page("/claude", func(c *via.Context) {
		claudePage(c, cfg, mockMode)
	})

	v.Start()
}
