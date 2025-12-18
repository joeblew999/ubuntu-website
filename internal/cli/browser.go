package cli

import (
	"os/exec"
	"runtime"
)

// OpenBrowser opens a URL in the default browser.
func OpenBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	default: // linux, freebsd, etc.
		cmd = "xdg-open"
		args = []string{url}
	}

	return exec.Command(cmd, args...).Start()
}

// Open opens a URL and prints a message.
func (c *Context) Open(url string) error {
	c.Printf("Opening %s...\n", url)
	return OpenBrowser(url)
}
