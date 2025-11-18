package web

import (
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
)

// pageContextLogger wraps log output to add human-readable page names to Via context errors
type pageContextLogger struct {
	output    io.Writer
	mu        sync.Mutex
	ctxRegex  *regexp.Regexp
	routeMap  map[string]string // maps context ID prefix to friendly page names
}

// newPageContextLogger creates a new logger that enhances Via error messages with page context
func newPageContextLogger() *pageContextLogger {
	return &pageContextLogger{
		output:   os.Stderr,
		ctxRegex: regexp.MustCompile(`via-ctx="([^"]+)"`),
		routeMap: map[string]string{
			"/_/":                 "Home",
			"/cloudflare_/":       "Cloudflare Setup - Step 1 (Token)",
			"/cloudflare/step2_/": "Cloudflare Setup - Step 2 (Account)",
			"/cloudflare/step3_/": "Cloudflare Setup - Step 3 (Domain)",
			"/cloudflare/step4_/": "Cloudflare Setup - Step 4 (Project)",
			"/claude_/":           "Claude AI Setup",
			"/deploy_/":           "Build & Deploy",
		},
	}
}

// Write implements io.Writer interface to intercept and enhance log messages
func (l *pageContextLogger) Write(p []byte) (n int, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	message := string(p)

	// Check if this is a Via error/warn log with context
	if strings.Contains(message, "via-ctx=") {
		// Extract the via-ctx value
		matches := l.ctxRegex.FindStringSubmatch(message)
		if len(matches) > 1 {
			ctxID := matches[1]

			// Try to find matching route prefix
			pageName := "Unknown Page"
			for routePrefix, friendlyName := range l.routeMap {
				if strings.HasPrefix(ctxID, routePrefix) {
					pageName = friendlyName
					break
				}
			}

			// Replace via-ctx with page-friendly format
			// From: [error] via-ctx="/_/dd297b1f" msg="..."
			// To:   [error] page="Home" via-ctx="/_/dd297b1f" msg="..."
			enhanced := strings.Replace(message,
				"via-ctx=\""+ctxID+"\"",
				"page=\""+pageName+"\" via-ctx=\""+ctxID+"\"",
				1)
			message = enhanced
		}
	}

	// Also handle errors without via-ctx but that come from Via
	// These appear as: [error] msg="failed to handle session close: ctx '/_/dd297b1f' not found"
	if strings.Contains(message, "[error]") && strings.Contains(message, "ctx '") && !strings.Contains(message, "via-ctx=") {
		ctxPattern := regexp.MustCompile(`ctx '([^']+)'`)
		matches := ctxPattern.FindStringSubmatch(message)
		if len(matches) > 1 {
			ctxID := matches[1]

			// Try to find matching route prefix
			pageName := "Unknown Page"
			for routePrefix, friendlyName := range l.routeMap {
				if strings.HasPrefix(ctxID, routePrefix) {
					pageName = friendlyName
					break
				}
			}

			// Insert page context after [error]
			enhanced := strings.Replace(message,
				"[error] ",
				"[error] page=\""+pageName+"\" ",
				1)
			message = enhanced
		}
	}

	return l.output.Write([]byte(message))
}

// SetupEnhancedLogging configures Go's standard logger to use our enhanced logger
func SetupEnhancedLogging() {
	logger := newPageContextLogger()
	log.SetOutput(logger)
	log.SetFlags(log.LstdFlags) // Keep timestamp
}
