// Package hugobuild proves the site still builds as production builds it, for the commands that change what it's
// built with (theme update, upgrade hugo).
package hugobuild

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Production is hugo -e production --gc --minify, into a temporary folder: its warnings and its summary printed, its
// first errors too if it fails.
func Production() error {
	tmp, err := os.MkdirTemp("", "hugobuild-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	log, err := exec.Command("hugo", "-e", "production", "--gc", "--minify", "-d", filepath.Join(tmp, "public")).CombinedOutput()
	lines := strings.Split(string(log), "\n")
	if err != nil {
		shown := 0
		for _, l := range lines {
			if (strings.HasPrefix(l, "ERROR") || strings.HasPrefix(l, "Error")) && shown < 5 {
				fmt.Println("    " + l)
				shown++
			}
		}
		return fmt.Errorf("the production build fails (hugo -e production)")
	}
	for _, l := range lines {
		if strings.HasPrefix(l, "WARN") || strings.Contains(l, "Pages ") || strings.Contains(l, "Total in") {
			fmt.Println("    " + l)
		}
	}
	return nil
}
